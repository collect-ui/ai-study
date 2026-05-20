package plugins

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	texttemplate "text/template"
	"time"

	common "github.com/collect-ui/collect/src/collect/common"
	config "github.com/collect-ui/collect/src/collect/config"
	collectFilters "github.com/collect-ui/collect/src/collect/filters"
	templateService "github.com/collect-ui/collect/src/collect/service_imp"
	"github.com/demdxx/gocast"
	"github.com/ledongthuc/pdf"
)

const (
	defaultQuestionAIDeepSeekBaseURL = "https://api.deepseek.com"
	defaultQuestionAIDeepSeekModel   = "deepseek-chat"
	defaultQuestionAICodexModel      = "gpt-5.4-mini"
	defaultQuestionAIChatGPTModel    = "gpt-5.5"
	defaultQuestionAIMaxChars        = 24000
	defaultQuestionAISourceMaxChars  = 120000
	defaultQuestionAIChunkChars      = 6000
	defaultQuestionAIMaxChunks       = 40
	defaultQuestionAIDeepSeekWorkers = 10
	defaultQuestionAICodexWorkers    = 2
	defaultQuestionAISystemPrompt    = "collect/question/prompts/ai_parse_system.md"
	defaultQuestionAIUserPrompt      = "collect/question/prompts/ai_parse_user.md"
)

var questionAIPaperHeadingRE = regexp.MustCompile(`^（[一二三四五六七八九十百]+）$`)
var questionAIPaperHeadingInlineRE = regexp.MustCompile(`（[一二三四五六七八九十百]+）`)
var questionAINumberMarkerRE = regexp.MustCompile(`(^|[^0-9])([1-9]|[1-9][0-9]|100)[\.．]`)
var questionAIAnswerHeadingRE = regexp.MustCompile(`(?m)(^|\n)\s*(答案[及与]?解析|答案解析|参考答案)\s*(\n|$)`)
var questionPDFAnswerHeadingBeforeRE = regexp.MustCompile(`([^\n])(答案[及与]?解析|答案解析|参考答案)`)
var questionPDFAnswerHeadingAfterRE = regexp.MustCompile(`(答案[及与]?解析|答案解析|参考答案)([^\n])`)
var questionPDFPaperHeadingBeforeRE = regexp.MustCompile(`([^\n])([（(][一二三四五六七八九十百]+[）)])`)
var questionPDFPaperHeadingAfterRE = regexp.MustCompile(`([（(][一二三四五六七八九十百]+[）)])([^\n])`)
var questionPDFSectionHeadingRE = regexp.MustCompile(`([^\n])([一二三四五六七八九十百]+、)`)
var questionPDFSectionTitleAfterRE = regexp.MustCompile(`([一二三四五六七八九十百]+、[^\n0-9A-Za-z]{2,20})([0-9A-Za-z_])`)
var questionPDFNumberMarkerRE = regexp.MustCompile(`([^\n\d])([1-9][0-9]{0,2}[\.．]\s*(?:[A-Z_]|[-–—]|[\p{Han}]|[（(]))`)
var questionPDFBareNumberOptionRE = regexp.MustCompile(`([^\n\d])([1-9][0-9]{0,2})\s+([A-H][\.．])`)
var questionPDFOptionMarkerRE = regexp.MustCompile(`([^\n])([A-H][\.．]\s*)`)
var questionPDFDialogueSpeakerRE = regexp.MustCompile(`([^\n])([A-Z]:\s*)`)

type QuestionPDFTextService struct {
	templateService.BaseHandler
}

type QuestionAIParseService struct {
	templateService.BaseHandler
}

type questionCodexAuthFile struct {
	AuthMode     string  `json:"auth_mode"`
	OPENAIAPIKey *string `json:"OPENAI_API_KEY"`
	Tokens       struct {
		AccessToken string `json:"access_token"`
		AccountID   string `json:"account_id"`
	} `json:"tokens"`
}

type questionAICredential struct {
	Token     string
	AccountID string
	Mode      string
	Source    string
}

type questionAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type questionAIChatRequest struct {
	Model       string                  `json:"model"`
	Messages    []questionAIChatMessage `json:"messages"`
	Temperature float64                 `json:"temperature"`
}

type questionAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type questionAIResponsesResponse struct {
	ID         string `json:"id"`
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type questionAIResponsesStreamState struct {
	deltas    strings.Builder
	doneText  string
	finalText string
}

type questionAITextChunk struct {
	Index       int
	Text        string
	SourceChars int
}

type questionAIPromptTemplateData struct {
	DefaultsJSON string
	RawText      string
	ChunkIndex   int
	ChunkTotal   int
}

func questionAIResponseText(result questionAIResponsesResponse) string {
	if text := strings.TrimSpace(result.OutputText); text != "" {
		return text
	}
	var builder strings.Builder
	for _, output := range result.Output {
		for _, content := range output.Content {
			if text := strings.TrimSpace(content.Text); text != "" {
				if builder.Len() > 0 {
					builder.WriteString("\n")
				}
				builder.WriteString(text)
			}
		}
	}
	return strings.TrimSpace(builder.String())
}

func questionAIStreamMapString(payload map[string]interface{}, key string) string {
	value, ok := payload[key]
	if !ok {
		return ""
	}
	return strings.TrimSpace(gocast.ToString(value))
}

func questionAIStreamErrorMessage(payload map[string]interface{}) string {
	if errorValue, ok := payload["error"]; ok && errorValue != nil {
		switch typed := errorValue.(type) {
		case string:
			return strings.TrimSpace(typed)
		case map[string]interface{}:
			for _, key := range []string{"message", "detail", "code", "type"} {
				if msg := strings.TrimSpace(gocast.ToString(typed[key])); msg != "" {
					return msg
				}
			}
		default:
			if msg := strings.TrimSpace(gocast.ToString(typed)); msg != "" {
				return msg
			}
		}
	}
	if responseValue, ok := payload["response"].(map[string]interface{}); ok {
		if msg := questionAIStreamErrorMessage(responseValue); msg != "" {
			return msg
		}
	}
	for _, key := range []string{"message", "detail"} {
		if msg := questionAIStreamMapString(payload, key); msg != "" {
			return msg
		}
	}
	return ""
}

func (state *questionAIResponsesStreamState) consume(eventType string, data string) error {
	data = strings.TrimSpace(data)
	if data == "" || data == "[DONE]" {
		return nil
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return fmt.Errorf("Codex SSE 数据不是合法 JSON: %w, data=%s", err, limitRunes(data, 400))
	}
	payloadType := questionAIStreamMapString(payload, "type")
	if payloadType == "" {
		payloadType = strings.TrimSpace(eventType)
	}
	if msg := questionAIStreamErrorMessage(payload); msg != "" && (payloadType == "error" || strings.Contains(payloadType, ".failed")) {
		return fmt.Errorf("Codex SSE 返回错误: %s", msg)
	}
	if delta, ok := payload["delta"].(string); ok && delta != "" && strings.Contains(payloadType, "delta") {
		state.deltas.WriteString(delta)
	}
	if text, ok := payload["text"].(string); ok && text != "" && strings.Contains(payloadType, "output_text.done") {
		state.doneText = text
	}
	if text := questionAIStreamMapString(payload, "output_text"); text != "" {
		state.finalText = text
	}
	if responseValue, ok := payload["response"]; ok && responseValue != nil {
		responseBytes, _ := json.Marshal(responseValue)
		var response questionAIResponsesResponse
		if err := json.Unmarshal(responseBytes, &response); err == nil {
			if response.Error != nil && strings.TrimSpace(response.Error.Message) != "" {
				return fmt.Errorf("Codex SSE 返回错误: %s", strings.TrimSpace(response.Error.Message))
			}
			if text := questionAIResponseText(response); text != "" {
				state.finalText = text
			}
		}
	}
	if state.finalText == "" {
		responseBytes, _ := json.Marshal(payload)
		var response questionAIResponsesResponse
		if err := json.Unmarshal(responseBytes, &response); err == nil {
			if response.Error != nil && strings.TrimSpace(response.Error.Message) != "" {
				return fmt.Errorf("Codex SSE 返回错误: %s", strings.TrimSpace(response.Error.Message))
			}
			if text := questionAIResponseText(response); text != "" {
				state.finalText = text
			}
		}
	}
	return nil
}

func (state *questionAIResponsesStreamState) text() string {
	if text := strings.TrimSpace(state.finalText); text != "" {
		return text
	}
	if text := strings.TrimSpace(state.doneText); text != "" {
		return text
	}
	return strings.TrimSpace(state.deltas.String())
}

func readQuestionAIResponsesStream(body io.Reader) (string, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	state := &questionAIResponsesStreamState{}
	eventType := ""
	dataLines := make([]string, 0, 4)
	dispatch := func() error {
		if len(dataLines) == 0 {
			eventType = ""
			return nil
		}
		data := strings.Join(dataLines, "\n")
		dataLines = dataLines[:0]
		currentEventType := eventType
		eventType = ""
		return state.consume(currentEventType, data)
	}
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" {
			if err := dispatch(); err != nil {
				return "", err
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field := line
		value := ""
		if index := strings.Index(line, ":"); index >= 0 {
			field = line[:index]
			value = line[index+1:]
			if strings.HasPrefix(value, " ") {
				value = value[1:]
			}
		}
		switch field {
		case "event":
			eventType = value
		case "data":
			dataLines = append(dataLines, value)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Codex SSE 读取失败: %w", err)
	}
	if err := dispatch(); err != nil {
		return "", err
	}
	if text := state.text(); text != "" {
		return text, nil
	}
	return "", fmt.Errorf("Codex SSE 返回为空")
}

type questionAIChunkSummary struct {
	Index           int  `json:"index"`
	SourceChars     int  `json:"source_chars"`
	QuestionMarkers int  `json:"question_markers"`
	RowCount        int  `json:"row_count"`
	FromCache       bool `json:"from_cache"`
}

type questionAIChunkResult struct {
	Index    int
	Rows     []map[string]interface{}
	Summary  questionAIChunkSummary
	Model    string
	Source   string
	Provider string
	Err      error
}

type questionAIChunkCacheEntry struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Source    string `json:"source"`
	FixedText string `json:"fixed_text"`
	CreatedAt string `json:"created_at"`
}

func appKeyAny(keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(collectFilters.GetKey(key))
		if value != "" {
			return value
		}
	}
	return ""
}

func intFromParams(params map[string]interface{}, key string, fallback int) int {
	value := gocast.ToInt(params[key])
	if value > 0 {
		return value
	}
	return fallback
}

func intParam(params map[string]interface{}, key string) (int, bool) {
	if _, ok := params[key]; !ok {
		return 0, false
	}
	value := gocast.ToInt(params[key])
	return value, value > 0
}

func limitRunes(text string, maxChars int) string {
	next, _ := limitRunesWithTruncated(text, maxChars)
	return next
}

func limitRunesWithTruncated(text string, maxChars int) (string, bool) {
	if maxChars <= 0 {
		return text, false
	}
	runes := []rune(text)
	if len(runes) <= maxChars {
		return text, false
	}
	return string(runes[:maxChars]), true
}

func runeCount(text string) int {
	return len([]rune(text))
}

func questionAINumberMarkerCount(text string) int {
	text = questionAITextBeforeAnswerSection(text)
	return len(questionAINumberMarkerRE.FindAllStringIndex(text, -1))
}

func questionAITextBeforeAnswerSection(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	match := questionAIAnswerHeadingRE.FindStringIndex(text)
	if match == nil {
		return text
	}
	return strings.TrimSpace(text[:match[0]])
}

func replacementRuneCount(text string) int {
	count := 0
	for _, r := range text {
		if r == '\uFFFD' {
			count++
		}
	}
	return count
}

func cleanExtractedPDFText(text string) string {
	var out strings.Builder
	lastWasReplacement := false
	lastWasControlSpace := false
	for _, r := range text {
		switch {
		case r == '\uFFFD':
			if !lastWasReplacement {
				out.WriteRune('\n')
			}
			lastWasReplacement = true
			lastWasControlSpace = false
		case r == '\r':
			out.WriteRune('\n')
			lastWasReplacement = false
			lastWasControlSpace = false
		case r == '\n':
			out.WriteRune('\n')
			lastWasReplacement = false
			lastWasControlSpace = false
		case r == '\t':
			if !lastWasControlSpace {
				out.WriteRune(' ')
			}
			lastWasReplacement = false
			lastWasControlSpace = true
		case r < 0x20 || r == 0x7f:
			if !lastWasControlSpace {
				out.WriteRune(' ')
			}
			lastWasReplacement = false
			lastWasControlSpace = true
		default:
			out.WriteRune(r)
			lastWasReplacement = false
			lastWasControlSpace = false
		}
	}
	return out.String()
}

func restoreQuestionPDFLineBreaks(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	replacements := []struct {
		re   *regexp.Regexp
		with string
	}{
		{questionPDFAnswerHeadingBeforeRE, "$1\n$2"},
		{questionPDFAnswerHeadingAfterRE, "$1\n$2"},
		{questionPDFPaperHeadingBeforeRE, "$1\n$2"},
		{questionPDFPaperHeadingAfterRE, "$1\n$2"},
		{questionPDFSectionHeadingRE, "$1\n$2"},
		{questionPDFSectionTitleAfterRE, "$1\n$2"},
		{questionPDFNumberMarkerRE, "$1\n$2"},
		{questionPDFBareNumberOptionRE, "$1\n$2\n$3"},
		{questionPDFOptionMarkerRE, "$1\n$2"},
		{questionPDFDialogueSpeakerRE, "$1\n$2"},
	}
	for _, replacement := range replacements {
		text = replacement.re.ReplaceAllString(text, replacement.with)
	}
	return text
}

func normalizeExtractedText(text string, maxChars int) string {
	text = cleanExtractedPDFText(text)
	text = restoreQuestionPDFLineBreaks(text)
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	lastBlank := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if !lastBlank {
				out = append(out, "")
			}
			lastBlank = true
			continue
		}
		out = append(out, line)
		lastBlank = false
	}
	return limitRunes(strings.TrimSpace(strings.Join(out, "\n")), maxChars)
}

func extractedTextToHTML(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	lines := strings.Split(text, "\n")
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			parts = append(parts, "<div><br></div>")
			continue
		}
		parts = append(parts, "<div>"+html.EscapeString(line)+"</div>")
	}
	return strings.Join(parts, "")
}

func isAllowedQuestionImportPath(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absPath = filepath.Clean(absPath)
	allowedRoots := []string{
		"/data/project/ai-study",
	}
	cwd, err := os.Getwd()
	if err == nil {
		if absCwd, absErr := filepath.Abs(cwd); absErr == nil {
			allowedRoots = append(allowedRoots, filepath.Clean(absCwd))
		}
	}
	for _, root := range allowedRoots {
		root = filepath.Clean(root)
		if absPath == root || strings.HasPrefix(absPath, root+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func extractTextFromPDFPath(path string, maxChars int) (string, error) {
	file, reader, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	plainReader, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(plainReader)
	if err != nil {
		return "", err
	}
	text := normalizeExtractedText(string(data), maxChars)
	if text == "" {
		return "", fmt.Errorf("PDF 未抽取到可用文本")
	}
	return text, nil
}

func extractTextFromLocalPath(path string, maxChars int) (string, error) {
	if !isAllowedQuestionImportPath(path) {
		return "", fmt.Errorf("不允许读取该文件路径")
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".pdf" {
		return extractTextFromPDFPath(path, maxChars)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := normalizeExtractedText(string(data), maxChars)
	if text == "" {
		return "", fmt.Errorf("文件未读取到可用文本")
	}
	return text, nil
}

func extractTextFromUploadedFile(ts *templateService.TemplateService, maxChars int) (string, string, error) {
	if ts.File == nil {
		return "", "", fmt.Errorf("上传文件不能为空")
	}
	ext := ".pdf"
	fileName := "upload.pdf"
	if ts.FileHeader != nil && strings.TrimSpace(ts.FileHeader.Filename) != "" {
		fileName = filepath.Base(ts.FileHeader.Filename)
		if nextExt := filepath.Ext(fileName); nextExt != "" {
			ext = nextExt
		}
	}
	tmpFile, err := os.CreateTemp("", "ai-study-question-import-*"+ext)
	if err != nil {
		return "", "", err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	if _, err := io.Copy(tmpFile, ts.File); err != nil {
		tmpFile.Close()
		return "", "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", "", err
	}
	text, err := extractTextFromPDFPath(tmpPath, maxChars)
	if err != nil {
		return "", "", err
	}
	return text, fileName, nil
}

func (s *QuestionPDFTextService) Result(template *config.Template, ts *templateService.TemplateService) *common.Result {
	params := template.GetParams()
	maxChars := defaultQuestionAISourceMaxChars
	if configMax := appKeyAny("question_ai_source_max_chars"); configMax != "" {
		if n, err := strconv.Atoi(configMax); err == nil && n > 0 {
			maxChars = n
		}
	}
	if n, ok := intParam(params, "max_chars"); ok {
		maxChars = n
	}

	var (
		text     string
		fileName string
		err      error
	)
	if ts.File != nil {
		text, fileName, err = extractTextFromUploadedFile(ts, maxChars)
	} else {
		filePath := strings.TrimSpace(gocast.ToString(params["file_path"]))
		if filePath == "" {
			filePath = appKeyAny("question_ai_default_pdf_path")
		}
		if filePath == "" {
			return common.NotOk("file_path 不能为空")
		}
		text, err = extractTextFromLocalPath(filePath, maxChars)
		fileName = filepath.Base(filePath)
	}
	if err != nil {
		return common.NotOk(err.Error())
	}
	return common.Ok(map[string]interface{}{
		"raw_text":     text,
		"raw_html":     extractedTextToHTML(text),
		"file_name":    fileName,
		"source_chars": runeCount(text),
		"max_chars":    maxChars,
	}, "PDF 文本抽取成功")
}

func questionAIPromptPath(configKey string, defaultPath string) string {
	if path := appKeyAny(configKey); path != "" {
		return path
	}
	return defaultPath
}

func readQuestionAIPromptTemplate(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("提示词模板路径不能为空")
	}
	candidates := []string{path}
	if !filepath.IsAbs(path) {
		candidates = append(candidates,
			filepath.Join("backend", path),
			filepath.Join("/data/project/ai-study/backend", path),
		)
		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates,
				filepath.Join(cwd, path),
				filepath.Join(cwd, "backend", path),
			)
			if parent := filepath.Dir(cwd); parent != "" && parent != "." {
				candidates = append(candidates, filepath.Join(parent, "backend", path))
			}
		}
	}
	seen := map[string]bool{}
	for _, candidate := range candidates {
		candidate = filepath.Clean(candidate)
		if seen[candidate] {
			continue
		}
		seen[candidate] = true
		data, err := os.ReadFile(candidate)
		if err == nil {
			text := strings.TrimSpace(string(data))
			if text == "" {
				return "", fmt.Errorf("提示词模板为空: %s", candidate)
			}
			return text, nil
		}
	}
	return "", fmt.Errorf("提示词模板不存在: %s", path)
}

func renderQuestionAIPromptTemplate(path string, data questionAIPromptTemplateData) string {
	text, err := readQuestionAIPromptTemplate(path)
	if err != nil {
		return fmt.Sprintf("提示词模板读取失败: %s", err.Error())
	}
	tmpl, err := texttemplate.New(filepath.Base(path)).Option("missingkey=error").Parse(text)
	if err != nil {
		return fmt.Sprintf("提示词模板解析失败: %s", err.Error())
	}
	var builder strings.Builder
	if err := tmpl.Execute(&builder, data); err != nil {
		return fmt.Sprintf("提示词模板渲染失败: %s", err.Error())
	}
	return strings.TrimSpace(builder.String())
}

func questionAIInstructions(defaults map[string]interface{}) string {
	defaultsJSON, _ := json.Marshal(defaults)
	return renderQuestionAIPromptTemplate(
		questionAIPromptPath("question_ai_system_prompt_path", defaultQuestionAISystemPrompt),
		questionAIPromptTemplateData{DefaultsJSON: string(defaultsJSON)},
	)
}

func questionAIUserPrompt(rawText string, chunkIndex int, chunkTotal int) string {
	return renderQuestionAIPromptTemplate(
		questionAIPromptPath("question_ai_user_prompt_path", defaultQuestionAIUserPrompt),
		questionAIPromptTemplateData{
			RawText:    rawText,
			ChunkIndex: chunkIndex,
			ChunkTotal: chunkTotal,
		},
	)
}

func normalizeAIJSONString(raw string) string {
	text := strings.TrimSpace(raw)
	text = strings.TrimPrefix(text, "\ufeff")
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) > 1 {
			lines = lines[1:]
			if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
				lines = lines[:len(lines)-1]
			}
			text = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}
	if json.Valid([]byte(text)) {
		return text
	}
	if repaired, ok := repairTruncatedJSON(text); ok {
		return repaired
	}
	if repaired, ok := repairJSONOutsideStringParens(text); ok {
		return repaired
	}
	firstObj := strings.Index(text, "{")
	firstArr := strings.Index(text, "[")
	start := firstObj
	if firstArr >= 0 && (start < 0 || firstArr < start) {
		start = firstArr
	}
	if start < 0 {
		return text
	}
	endObj := strings.LastIndex(text, "}")
	endArr := strings.LastIndex(text, "]")
	end := endObj
	if endArr > end {
		end = endArr
	}
	if end >= start {
		candidate := strings.TrimSpace(text[start : end+1])
		if json.Valid([]byte(candidate)) {
			return candidate
		}
		if repaired, ok := repairTruncatedJSON(candidate); ok {
			return repaired
		}
		if repaired, ok := repairJSONOutsideStringParens(candidate); ok {
			return repaired
		}
	}
	return text
}

func repairTruncatedJSON(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", false
	}
	stack := make([]rune, 0)
	inString := false
	escaped := false
	for _, r := range text {
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if r == '\\' {
				escaped = true
				continue
			}
			if r == '"' {
				inString = false
			}
			continue
		}
		switch r {
		case '"':
			inString = true
		case '{':
			stack = append(stack, '}')
		case '[':
			stack = append(stack, ']')
		case '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != r {
				return "", false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if inString || len(stack) == 0 {
		return "", false
	}
	var builder strings.Builder
	builder.WriteString(text)
	for i := len(stack) - 1; i >= 0; i-- {
		builder.WriteRune(stack[i])
	}
	repaired := builder.String()
	if !json.Valid([]byte(repaired)) {
		return "", false
	}
	return repaired, true
}

func repairJSONOutsideStringParens(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", false
	}
	var builder strings.Builder
	builder.Grow(len(text))
	stack := make([]rune, 0)
	inString := false
	escaped := false
	changed := false
	for _, r := range text {
		if inString {
			builder.WriteRune(r)
			if escaped {
				escaped = false
				continue
			}
			if r == '\\' {
				escaped = true
				continue
			}
			if r == '"' {
				inString = false
			}
			continue
		}
		switch r {
		case '"':
			inString = true
			builder.WriteRune(r)
		case '{':
			stack = append(stack, '}')
			builder.WriteRune(r)
		case '[':
			stack = append(stack, ']')
			builder.WriteRune(r)
		case '}', ']':
			if len(stack) > 0 && stack[len(stack)-1] == r {
				stack = stack[:len(stack)-1]
			}
			builder.WriteRune(r)
		case ')':
			changed = true
			if len(stack) > 0 && stack[len(stack)-1] == '}' {
				stack = stack[:len(stack)-1]
				builder.WriteRune('}')
			}
		case '(':
			changed = true
		default:
			builder.WriteRune(r)
		}
	}
	if inString || !changed {
		return "", false
	}
	repaired := strings.TrimSpace(builder.String())
	if json.Valid([]byte(repaired)) {
		return repaired, true
	}
	if completed, ok := repairTruncatedJSON(repaired); ok {
		return completed, true
	}
	return "", false
}

func parseQuestionRows(fixedText string) ([]map[string]interface{}, error) {
	text := normalizeAIJSONString(fixedText)
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(text), &arr); err == nil {
		return arr, nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(text), &obj); err != nil {
		return nil, fmt.Errorf("AI 返回不是合法 JSON: %w", err)
	}
	rawQuestions, ok := obj["questions"]
	if !ok {
		return nil, fmt.Errorf("JSON 中缺少 questions 字段")
	}
	bytesData, _ := json.Marshal(rawQuestions)
	if err := json.Unmarshal(bytesData, &arr); err != nil {
		return nil, fmt.Errorf("questions 不是题目数组: %w", err)
	}
	return arr, nil
}

func formatQuestionAIJSONParseError(prefix string, rawText string, err error) string {
	rawText = strings.TrimSpace(rawText)
	var builder strings.Builder
	builder.WriteString(prefix)
	builder.WriteString(": ")
	builder.WriteString(err.Error())
	builder.WriteString(fmt.Sprintf("\n\n原始返回字符数: %d", runeCount(rawText)))
	if rawText == "" {
		builder.WriteString("\n\n原始返回为空")
		return builder.String()
	}
	builder.WriteString("\n\n原始返回:\n")
	builder.WriteString(rawText)
	return builder.String()
}

func splitQuestionAIText(rawText string, chunkChars int, maxChunks int) ([]questionAITextChunk, error) {
	rawText = strings.TrimSpace(rawText)
	if rawText == "" {
		return nil, nil
	}
	if chunkChars <= 0 {
		chunkChars = defaultQuestionAIChunkChars
	}
	if maxChunks <= 0 {
		maxChunks = defaultQuestionAIMaxChunks
	}
	units := splitQuestionAIUnits(rawText)
	if len(units) == 0 {
		units = []string{rawText}
	}
	parts := make([]string, 0, len(units))
	for _, unit := range units {
		parts = append(parts, splitQuestionAILongUnit(unit, chunkChars)...)
	}
	chunks := make([]questionAITextChunk, 0, len(parts))
	addChunk := func(text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		chunks = append(chunks, questionAITextChunk{
			Index:       len(chunks) + 1,
			Text:        text,
			SourceChars: runeCount(text),
		})
	}
	for _, part := range parts {
		addChunk(part)
	}
	if len(chunks) > maxChunks {
		return nil, fmt.Errorf("PDF 分段数量 %d 超过 max_chunks=%d，请调大 chunk_chars 或 max_chunks", len(chunks), maxChunks)
	}
	return chunks, nil
}

func splitQuestionAIUnits(rawText string) []string {
	rawText = strings.TrimSpace(rawText)
	if rawText == "" {
		return nil
	}
	if units := splitQuestionAIUnitsByInlineHeading(rawText); len(units) > 1 {
		return units
	}
	lines := strings.Split(rawText, "\n")
	units := make([]string, 0)
	current := make([]string, 0)
	flush := func() {
		text := strings.TrimSpace(strings.Join(current, "\n"))
		if text != "" {
			units = append(units, text)
		}
		current = current[:0]
	}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if questionAIPaperHeadingRE.MatchString(trimmed) && len(current) > 0 {
			flush()
		}
		current = append(current, line)
	}
	flush()
	return units
}

func splitQuestionAIUnitsByInlineHeading(rawText string) []string {
	matches := questionAIPaperHeadingInlineRE.FindAllStringIndex(rawText, -1)
	if len(matches) <= 1 {
		return nil
	}
	units := make([]string, 0, len(matches))
	for index, match := range matches {
		start := match[0]
		if index == 0 {
			start = 0
		}
		end := len(rawText)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		text := strings.TrimSpace(rawText[start:end])
		if text != "" {
			units = append(units, text)
		}
	}
	return units
}

func splitQuestionAILongUnit(unit string, chunkChars int) []string {
	if runeCount(unit) <= chunkChars {
		return []string{unit}
	}
	lines := strings.Split(unit, "\n")
	parts := make([]string, 0, 2)
	var current strings.Builder
	currentChars := 0
	flush := func() {
		text := strings.TrimSpace(current.String())
		if text != "" {
			parts = append(parts, text)
		}
		current.Reset()
		currentChars = 0
	}
	for _, line := range lines {
		lineChars := runeCount(line)
		if lineChars > chunkChars {
			flush()
			runes := []rune(line)
			for start := 0; start < len(runes); start += chunkChars {
				end := start + chunkChars
				if end > len(runes) {
					end = len(runes)
				}
				parts = append(parts, string(runes[start:end]))
			}
			continue
		}
		separatorChars := 0
		if currentChars > 0 {
			separatorChars = 1
		}
		if currentChars > 0 && currentChars+separatorChars+lineChars > chunkChars {
			flush()
		}
		if currentChars > 0 {
			current.WriteString("\n")
			currentChars++
		}
		current.WriteString(line)
		currentChars += lineChars
	}
	flush()
	return parts
}

func mergeQuestionRows(rowGroups [][]map[string]interface{}) []map[string]interface{} {
	rows := make([]map[string]interface{}, 0)
	seen := map[string]bool{}
	for _, group := range rowGroups {
		for _, row := range group {
			key := questionRowDedupKey(row)
			if key != "" {
				if seen[key] {
					continue
				}
				seen[key] = true
			}
			rows = append(rows, row)
		}
	}
	return rows
}

func questionRowDedupKey(row map[string]interface{}) string {
	for _, key := range []string{"question_code", "code"} {
		value := strings.TrimSpace(gocast.ToString(row[key]))
		if value != "" {
			return "code:" + strings.ToLower(value)
		}
	}
	for _, key := range []string{"stem_text", "stem", "question", "content", "title"} {
		value := strings.TrimSpace(gocast.ToString(row[key]))
		if value != "" {
			value = strings.ToLower(strings.Join(strings.Fields(value), " "))
			return "stem:" + value
		}
	}
	return ""
}

func normalizeQuestionAIImportRows(rows []map[string]interface{}, defaults map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rows))
	for index, row := range rows {
		out = append(out, normalizeQuestionAIImportRow(row, defaults, index+1))
	}
	return out
}

func normalizeQuestionAIImportRow(row map[string]interface{}, defaults map[string]interface{}, index int) map[string]interface{} {
	if row == nil {
		row = map[string]interface{}{}
	}
	next := map[string]interface{}{}
	for key, value := range row {
		next[key] = value
	}

	category := normalizeQuestionCategory(firstAIString(next, "question_category", "category", "questionCategory"))
	if category == "" {
		category = strings.TrimSpace(gocast.ToString(defaults["question_category"]))
	}
	choiceItems := normalizeAIChoiceItems(firstAIValue(next, "choice_items", "sub_questions", "subQuestions", "sub_items", "items"), category)
	if category == "" && len(choiceItems) > 0 {
		category = "reading_choice"
	}
	if category == "" {
		category = "normal"
	}
	groupedChoice := category == "reading_choice" || category == "cloze_choice"

	qType := normalizeQuestionType(firstAIString(next, "question_type", "type", "questionType"))
	if qType == "" {
		qType = strings.TrimSpace(gocast.ToString(defaults["question_type"]))
	}
	if groupedChoice {
		qType = "single_choice"
	} else if qType == "" || qType == "normal" {
		qType = inferQuestionAIType(next)
	}
	if qType == "" {
		qType = "single_choice"
	}

	stem := firstAIString(next, "stem_text", "stem", "question", "content")
	answerKey := normalizeQuestionAIAnswer(firstAIString(next, "answer_key", "answer", "answer_text"), qType)
	optionA := optionAIText(next, "A", 0)
	optionB := optionAIText(next, "B", 1)
	optionC := optionAIText(next, "C", 2)
	optionD := optionAIText(next, "D", 3)
	blankAnswers := normalizeAIBlankAnswers(firstAIValue(next, "blank_answers", "blankAnswers"), firstAIString(next, "answer_key", "answer", "answer_text"))

	if groupedChoice {
		answerKey = groupedChoiceAnswerText(choiceItems, category)
		next["choice_items"] = choiceItems
		next["reference_text"] = mustJSON(choiceItems)
		next["answer_key"] = answerKey
		next["answer_text"] = answerKey
		next["answer_value"] = mustJSON(groupedChoiceAnswers(choiceItems))
		next["option_count"] = 0
		next["blank_count"] = len(choiceItems)
		next["option_a_text"] = ""
		next["option_b_text"] = ""
		next["option_c_text"] = ""
		next["option_d_text"] = ""
		next["option_a_html"] = ""
		next["option_b_html"] = ""
		next["option_c_html"] = ""
		next["option_d_html"] = ""
	} else {
		next["answer_text"] = answerKey
		next["answer_value"] = mustJSON(answerAIList(answerKey, qType, blankAnswers))
		next["reference_text"] = firstAIString(next, "reference_text")
		if optionA != "" {
			next["option_a_text"] = optionA
			next["option_a_html"] = firstAIString(next, "option_a_html")
			if firstAIString(next, "option_a_html") == "" {
				next["option_a_html"] = optionA
			}
		}
		if optionB != "" {
			next["option_b_text"] = optionB
			next["option_b_html"] = firstAIString(next, "option_b_html")
			if firstAIString(next, "option_b_html") == "" {
				next["option_b_html"] = optionB
			}
		}
		if optionC != "" {
			next["option_c_text"] = optionC
			next["option_c_html"] = firstAIString(next, "option_c_html")
			if firstAIString(next, "option_c_html") == "" {
				next["option_c_html"] = optionC
			}
		}
		if optionD != "" {
			next["option_d_text"] = optionD
			next["option_d_html"] = firstAIString(next, "option_d_html")
			if firstAIString(next, "option_d_html") == "" {
				next["option_d_html"] = optionD
			}
		}
		next["option_count"] = optionCountForAI(qType, optionA, optionB, optionC, optionD)
		next["blank_answers"] = blankAnswers
		if qType == "blank" {
			next["blank_count"] = len(blankAnswers)
			if len(blankAnswers) == 0 {
				next["blank_count"] = 1
			}
		} else {
			next["blank_count"] = 0
		}
	}

	next["question_type"] = qType
	next["question_category"] = category
	next["stem_text"] = stem
	if firstAIString(next, "stem_html") == "" {
		next["stem_html"] = stem
	}
	if firstAIString(next, "analysis_text", "analysis", "explanation") != "" {
		next["analysis_text"] = firstAIString(next, "analysis_text", "analysis", "explanation")
	}
	if firstAIString(next, "analysis_html") == "" {
		next["analysis_html"] = firstAIString(next, "analysis_text")
	}
	if firstAIString(next, "title") == "" {
		if stem != "" {
			runes := []rune(stem)
			if len(runes) > 40 {
				next["title"] = string(runes[:40])
			} else {
				next["title"] = stem
			}
		} else {
			next["title"] = fmt.Sprintf("AI导入题目%d", index)
		}
	}
	for _, key := range []string{"subject", "stage", "grade", "textbook_version", "question_category", "difficulty"} {
		if firstAIString(next, key) == "" {
			next[key] = defaults[key]
		}
	}
	if gocast.ToInt(next["score"]) <= 0 {
		next["score"] = defaults["score"]
	}
	if gocast.ToInt(next["sequence_no"]) <= 0 {
		next["sequence_no"] = index
	}
	next["option_a_correct"] = boolAIString(strings.Contains(","+strings.ToUpper(answerKey)+",", ",A,"))
	next["option_b_correct"] = boolAIString(strings.Contains(","+strings.ToUpper(answerKey)+",", ",B,"))
	next["option_c_correct"] = boolAIString(strings.Contains(","+strings.ToUpper(answerKey)+",", ",C,"))
	next["option_d_correct"] = boolAIString(strings.Contains(","+strings.ToUpper(answerKey)+",", ",D,"))
	return next
}

func firstAIValue(row map[string]interface{}, keys ...string) interface{} {
	for _, key := range keys {
		if value, ok := row[key]; ok && value != nil {
			if text, ok := value.(string); ok && strings.TrimSpace(text) == "" {
				continue
			}
			return value
		}
	}
	return nil
}

func firstAIString(row map[string]interface{}, keys ...string) string {
	value := firstAIValue(row, keys...)
	if value == nil {
		return ""
	}
	return strings.TrimSpace(gocast.ToString(value))
}

func normalizeQuestionCategory(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "-", "_")
	switch value {
	case "阅读理解", "reading", "reading_comprehension", "readingchoice":
		return "reading_choice"
	case "完形填空", "完形", "cloze", "clozechoice":
		return "cloze_choice"
	case "单项选择", "选择题", "语法选择", "语法/单项选择", "grammar", "grammarchoice":
		return "grammar_choice"
	case "词形填空", "用所给词填空", "填空", "fill", "fillword":
		return "fill_word"
	case "判断正误", "判断题", "true_false", "judge", "judgetf":
		return "judge_tf"
	case "阅读回答", "阅读简答", "short_reading", "readingshortanswer":
		return "reading_short_answer"
	case "补全对话", "dialogue", "dialoguecompletion":
		return "dialogue_completion"
	case "normal", "grammar_choice", "vocabulary_choice", "fill_word", "reading_short_answer", "judge_tf", "dialogue_completion", "iq":
		return value
	}
	return value
}

func normalizeQuestionType(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "-", "_")
	switch value {
	case "单选", "单项选择", "选择题", "single", "singlechoice":
		return "single_choice"
	case "多选", "multiple", "multiplechoice":
		return "multiple_choice"
	case "填空", "blank", "fill_blank":
		return "blank"
	case "判断", "judge", "true_false":
		return "judge"
	case "简答", "问答", "short", "shortanswer":
		return "short_answer"
	}
	return value
}

func inferQuestionAIType(row map[string]interface{}) string {
	if len(normalizeAIBlankAnswers(firstAIValue(row, "blank_answers", "blankAnswers"), firstAIString(row, "answer_key", "answer", "answer_text"))) > 0 {
		return "blank"
	}
	answer := strings.ToLower(firstAIString(row, "answer_key", "answer", "answer_text"))
	if answer == "true" || answer == "false" || answer == "对" || answer == "错" {
		return "judge"
	}
	if optionAIText(row, "A", 0) != "" && optionAIText(row, "B", 1) != "" {
		return "single_choice"
	}
	return "short_answer"
}

func optionAIText(row map[string]interface{}, key string, index int) string {
	lower := strings.ToLower(key)
	if value := firstAIString(row, "option_"+lower+"_text", "option_"+lower, "option"+key, key, lower); value != "" {
		return value
	}
	options := firstAIValue(row, "options", "choices")
	switch value := options.(type) {
	case map[string]interface{}:
		return strings.TrimSpace(gocast.ToString(firstAIValue(value, key, lower)))
	case []interface{}:
		if index < len(value) {
			return optionAIArrayText(value[index])
		}
	case []map[string]interface{}:
		for _, item := range value {
			if strings.EqualFold(firstAIString(item, "key", "label", "option_key"), key) {
				return firstAIString(item, "text", "value", "content", "option_text")
			}
		}
		if index < len(value) {
			return optionAIArrayText(value[index])
		}
	}
	return ""
}

func optionAIArrayText(value interface{}) string {
	switch item := value.(type) {
	case string:
		return strings.TrimSpace(item)
	case map[string]interface{}:
		return firstAIString(item, "text", "value", "content", "option_text")
	default:
		return strings.TrimSpace(gocast.ToString(value))
	}
}

func normalizeAIChoiceItems(value interface{}, category string) []map[string]interface{} {
	rawItems := aiMapSlice(value)
	out := make([]map[string]interface{}, 0, len(rawItems))
	for index, raw := range rawItems {
		if raw == nil {
			continue
		}
		item := map[string]interface{}{}
		if category == "cloze_choice" {
			item["__rowId"] = firstAIString(raw, "__rowId")
			if item["__rowId"] == "" {
				item["__rowId"] = fmt.Sprintf("choice_%d", index+1)
			}
			item["blank_no"] = firstAIString(raw, "blank_no", "blank_index", "no", "sub_no", "index")
		} else {
			item["__rowId"] = firstAIString(raw, "__rowId")
			if item["__rowId"] == "" {
				item["__rowId"] = fmt.Sprintf("reading_%d", index+1)
			}
			item["sub_no"] = firstAIString(raw, "sub_no", "no", "blank_no", "index")
			item["question_text"] = firstAIString(raw, "question_text", "stem_text", "stem", "question", "content", "title")
		}
		item["option_a"] = optionAIText(raw, "A", 0)
		item["option_b"] = optionAIText(raw, "B", 1)
		item["option_c"] = optionAIText(raw, "C", 2)
		item["option_d"] = optionAIText(raw, "D", 3)
		item["answer_key"] = normalizeQuestionAIAnswer(firstAIString(raw, "answer_key", "answer", "answer_text"), "single_choice")
		item["analysis"] = firstAIString(raw, "analysis", "analysis_text", "explanation")
		if item["sub_no"] == "" && item["blank_no"] == "" {
			if category == "cloze_choice" {
				item["blank_no"] = strconv.Itoa(index + 1)
			} else {
				item["sub_no"] = strconv.Itoa(index + 1)
			}
		}
		out = append(out, item)
	}
	return out
}

func aiMapSlice(value interface{}) []map[string]interface{} {
	if value == nil {
		return nil
	}
	if text, ok := value.(string); ok {
		text = strings.TrimSpace(text)
		if text == "" {
			return nil
		}
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(text), &arr); err == nil {
			return arr
		}
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil
	}
	return arr
}

func normalizeAIBlankAnswers(value interface{}, answerText string) []map[string]interface{} {
	rows := aiMapSlice(value)
	if len(rows) == 0 && strings.TrimSpace(answerText) != "" {
		parts := strings.FieldsFunc(answerText, func(r rune) bool {
			return r == ';' || r == '；'
		})
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				rows = append(rows, map[string]interface{}{"standard_answer": part})
			}
		}
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for index, row := range rows {
		answer := firstAIString(row, "standard_answer", "answer", "value")
		if answer == "" {
			continue
		}
		out = append(out, map[string]interface{}{
			"blank_index":         gocast.ToInt(firstAIString(row, "blank_index", "index")),
			"standard_answer":     answer,
			"alternative_answers": firstOrDefault(firstAIString(row, "alternative_answers"), "[]"),
			"score":               gocast.ToInt(firstAIString(row, "score")),
			"match_mode":          firstOrDefault(firstAIString(row, "match_mode"), "exact"),
			"case_sensitive":      firstOrDefault(firstAIString(row, "case_sensitive"), "0"),
		})
		if gocast.ToInt(out[len(out)-1]["blank_index"]) <= 0 {
			out[len(out)-1]["blank_index"] = index + 1
		}
	}
	return out
}

func normalizeQuestionAIAnswer(value string, qType string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "，", ","))
	if qType == "judge" {
		lower := strings.ToLower(value)
		if lower == "false" || value == "错" || value == "否" || lower == "f" || lower == "no" {
			return "false"
		}
		return "true"
	}
	if qType == "single_choice" {
		letters := answerLetters(value)
		if len(letters) > 0 {
			return letters[0]
		}
	}
	if qType == "multiple_choice" {
		letters := answerLetters(value)
		if len(letters) > 0 {
			return strings.Join(letters, ",")
		}
	}
	return value
}

func answerLetters(value string) []string {
	matches := regexp.MustCompile(`[A-H]`).FindAllString(strings.ToUpper(value), -1)
	seen := map[string]bool{}
	letters := make([]string, 0, len(matches))
	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			letters = append(letters, match)
		}
	}
	return letters
}

func answerAIList(answerKey string, qType string, blankAnswers []map[string]interface{}) []string {
	switch qType {
	case "multiple_choice":
		return answerLetters(answerKey)
	case "single_choice":
		if answerKey == "" {
			return nil
		}
		return []string{answerKey}
	case "judge":
		return []string{answerKey}
	case "blank":
		values := make([]string, 0, len(blankAnswers))
		for _, row := range blankAnswers {
			if value := firstAIString(row, "standard_answer"); value != "" {
				values = append(values, value)
			}
		}
		return values
	default:
		if answerKey == "" {
			return nil
		}
		return []string{answerKey}
	}
}

func groupedChoiceAnswerText(items []map[string]interface{}, category string) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		no := firstAIString(item, "sub_no", "blank_no")
		if category == "cloze_choice" {
			no = firstAIString(item, "blank_no", "sub_no")
		}
		answer := firstAIString(item, "answer_key")
		if no != "" && answer != "" {
			parts = append(parts, no+":"+answer)
		}
	}
	return strings.Join(parts, "；")
}

func groupedChoiceAnswers(items []map[string]interface{}) []string {
	answers := make([]string, 0, len(items))
	for _, item := range items {
		if answer := firstAIString(item, "answer_key"); answer != "" {
			answers = append(answers, answer)
		}
	}
	return answers
}

func optionCountForAI(qType string, optionA string, optionB string, optionC string, optionD string) int {
	if qType == "judge" {
		return 2
	}
	if qType != "single_choice" && qType != "multiple_choice" {
		return 0
	}
	count := 0
	for _, value := range []string{optionA, optionB, optionC, optionD} {
		if strings.TrimSpace(value) != "" {
			count++
		}
	}
	return count
}

func firstOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func boolAIString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func mustJSON(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func buildAIParseDefaults(params map[string]interface{}) map[string]interface{} {
	defaults := map[string]interface{}{
		"subject":           "english",
		"stage":             "junior",
		"grade":             "grade_7",
		"textbook_version":  "pep",
		"question_type":     "single_choice",
		"question_category": "normal",
		"difficulty":        "basic",
		"score":             5,
	}
	for _, key := range []string{"subject", "stage", "grade", "textbook_version", "unit_id", "unit_code", "unit_name", "question_type", "question_category", "difficulty"} {
		if value := strings.TrimSpace(gocast.ToString(params[key])); value != "" {
			defaults[key] = value
		}
	}
	if score := gocast.ToInt(params["score"]); score > 0 {
		defaults["score"] = score
	}
	return defaults
}

func callQuestionAIProvider(params map[string]interface{}, provider string, systemPrompt string, userPrompt string) (string, string, string, string, error) {
	model := strings.TrimSpace(gocast.ToString(params["model"]))
	source := ""
	var err error
	switch provider {
	case "deepseek":
		var fixedText string
		fixedText, model, err = callDeepSeekQuestionAI(params, systemPrompt, userPrompt)
		source = "deepseek"
		return fixedText, model, source, provider, err
	case "codex", "codex_auth", "auth_json":
		var fixedText string
		fixedText, model, source, err = callCodexQuestionAI(params, systemPrompt, userPrompt)
		return fixedText, model, source, "codex", err
	default:
		return "", "", "", provider, fmt.Errorf("不支持的 AI 解析模式: %s", provider)
	}
}

func questionAIEffectiveModel(params map[string]interface{}, provider string) string {
	model := strings.TrimSpace(gocast.ToString(params["model"]))
	if model != "" {
		return model
	}
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "deepseek":
		model = appKeyAny("question_ai_deepseek_model", "deepseek_model")
		if model == "" {
			model = defaultQuestionAIDeepSeekModel
		}
	default:
		model = questionAICodexModelForCredential(params, resolveQuestionCodexCredential())
	}
	return model
}

func questionAICodexModelForCredential(params map[string]interface{}, credential questionAICredential) string {
	if model := strings.TrimSpace(gocast.ToString(params["model"])); model != "" {
		return model
	}
	configuredModel := appKeyAny("question_ai_codex_model", "codex_model")
	if credential.Mode == "chatgpt_access_token" {
		if model := appKeyAny("question_ai_codex_chatgpt_model", "codex_chatgpt_model"); model != "" {
			return model
		}
		if configuredModel == "" || configuredModel == defaultQuestionAICodexModel {
			return defaultQuestionAIChatGPTModel
		}
	}
	if configuredModel != "" {
		return configuredModel
	}
	return defaultQuestionAICodexModel
}

func questionAIChunkConcurrency(params map[string]interface{}, provider string, chunkCount int) int {
	if chunkCount <= 1 {
		return 1
	}
	if n, ok := intParam(params, "concurrency"); ok {
		if n > chunkCount {
			return chunkCount
		}
		return n
	}
	concurrency := defaultQuestionAICodexWorkers
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "deepseek":
		concurrency = defaultQuestionAIDeepSeekWorkers
	}
	if concurrency > chunkCount {
		concurrency = chunkCount
	}
	if concurrency < 1 {
		concurrency = 1
	}
	return concurrency
}

func questionAIChunkCacheDir() string {
	dir := appKeyAny("question_ai_cache_dir")
	if dir == "" {
		dir = "./database/question_ai_cache"
	}
	return dir
}

func questionAIChunkCacheKey(provider string, model string, systemPrompt string, userPrompt string) string {
	payload, _ := json.Marshal(map[string]string{
		"version":       "question-ai-import-v3",
		"provider":      strings.ToLower(strings.TrimSpace(provider)),
		"model":         strings.TrimSpace(model),
		"system_prompt": systemPrompt,
		"user_prompt":   userPrompt,
	})
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func questionAIChunkCachePath(key string) string {
	return filepath.Join(questionAIChunkCacheDir(), key+".json")
}

func readQuestionAIChunkCache(key string) (questionAIChunkCacheEntry, bool) {
	var entry questionAIChunkCacheEntry
	data, err := os.ReadFile(questionAIChunkCachePath(key))
	if err != nil {
		return entry, false
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return questionAIChunkCacheEntry{}, false
	}
	if strings.TrimSpace(entry.FixedText) == "" {
		return questionAIChunkCacheEntry{}, false
	}
	return entry, true
}

func writeQuestionAIChunkCache(key string, entry questionAIChunkCacheEntry) {
	dir := questionAIChunkCacheDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	path := questionAIChunkCachePath(key)
	tmp, err := os.CreateTemp(dir, ".question-ai-cache-*")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	encoder := json.NewEncoder(tmp)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(entry)
	closeErr := tmp.Close()
	if err != nil || closeErr != nil {
		_ = os.Remove(tmpPath)
		return
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
	}
}

func callDeepSeekQuestionAI(params map[string]interface{}, systemPrompt string, userPrompt string) (string, string, error) {
	apiKey := strings.TrimSpace(gocast.ToString(params["api_key"]))
	if apiKey == "" {
		apiKey = appKeyAny("question_ai_deepseek_api_key", "deepseek_api_key")
	}
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	}
	if apiKey == "" {
		return "", "", fmt.Errorf("DeepSeek API Key 未配置")
	}
	baseURL := strings.TrimSpace(gocast.ToString(params["base_url"]))
	if baseURL == "" {
		baseURL = appKeyAny("question_ai_deepseek_base_url", "deepseek_base_url")
	}
	if baseURL == "" {
		baseURL = defaultQuestionAIDeepSeekBaseURL
	}
	model := strings.TrimSpace(gocast.ToString(params["model"]))
	if model == "" {
		model = appKeyAny("question_ai_deepseek_model", "deepseek_model")
	}
	if model == "" {
		model = defaultQuestionAIDeepSeekModel
	}
	requestData := questionAIChatRequest{
		Model: model,
		Messages: []questionAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.1,
	}
	body, err := json.Marshal(requestData)
	if err != nil {
		return "", model, err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", model, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return "", model, fmt.Errorf("DeepSeek 请求失败: %w", err)
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	var result questionAIChatResponse
	_ = json.Unmarshal(rawBody, &result)
	if resp.StatusCode >= 300 {
		if result.Error != nil && strings.TrimSpace(result.Error.Message) != "" {
			return "", model, fmt.Errorf("DeepSeek 状态异常: %d, %s", resp.StatusCode, strings.TrimSpace(result.Error.Message))
		}
		return "", model, fmt.Errorf("DeepSeek 状态异常: %d", resp.StatusCode)
	}
	if result.Error != nil && strings.TrimSpace(result.Error.Message) != "" {
		return "", model, fmt.Errorf("DeepSeek 返回错误: %s", strings.TrimSpace(result.Error.Message))
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", model, fmt.Errorf("DeepSeek 返回为空")
	}
	return result.Choices[0].Message.Content, model, nil
}

func loadQuestionCodexAuthFile() (*questionCodexAuthFile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	authPath := filepath.Join(homeDir, ".codex", "auth.json")
	data, err := os.ReadFile(authPath)
	if err != nil {
		return nil, err
	}
	var authData questionCodexAuthFile
	if err := json.Unmarshal(data, &authData); err != nil {
		return nil, err
	}
	return &authData, nil
}

func resolveQuestionCodexCredential() questionAICredential {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	accountID := strings.TrimSpace(os.Getenv("OPENAI_ACCOUNT_ID"))
	if apiKey != "" {
		return questionAICredential{Token: apiKey, AccountID: accountID, Mode: "platform_api_key", Source: "env_openai_api_key"}
	}
	authData, err := loadQuestionCodexAuthFile()
	if err != nil || authData == nil {
		return questionAICredential{}
	}
	if authData.OPENAIAPIKey != nil && strings.TrimSpace(*authData.OPENAIAPIKey) != "" {
		return questionAICredential{
			Token:     strings.TrimSpace(*authData.OPENAIAPIKey),
			AccountID: strings.TrimSpace(authData.Tokens.AccountID),
			Mode:      "platform_api_key",
			Source:    "codex_auth_openai_api_key",
		}
	}
	if strings.TrimSpace(authData.Tokens.AccessToken) != "" {
		return questionAICredential{
			Token:     strings.TrimSpace(authData.Tokens.AccessToken),
			AccountID: strings.TrimSpace(authData.Tokens.AccountID),
			Mode:      "chatgpt_access_token",
			Source:    "codex_auth_access_token",
		}
	}
	return questionAICredential{}
}

func questionCodexBaseURL(credential questionAICredential) string {
	if baseURL := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	if credential.Mode == "chatgpt_access_token" {
		return "https://chatgpt.com/backend-api/codex"
	}
	return "https://api.openai.com/v1"
}

func callCodexQuestionAI(params map[string]interface{}, systemPrompt string, userPrompt string) (string, string, string, error) {
	credential := resolveQuestionCodexCredential()
	if strings.TrimSpace(credential.Token) == "" {
		return "", "", "", fmt.Errorf("Codex auth.json 或 OPENAI_API_KEY 未配置")
	}
	model := questionAICodexModelForCredential(params, credential)
	baseURL := questionCodexBaseURL(credential)
	var requestBody map[string]interface{}
	if credential.Mode == "chatgpt_access_token" {
		requestBody = map[string]interface{}{
			"model": model,
			"input": []map[string]interface{}{
				{
					"type": "message",
					"role": "user",
					"content": []map[string]interface{}{
						{"type": "input_text", "text": userPrompt},
					},
				},
			},
			"instructions":        systemPrompt,
			"tools":               []interface{}{},
			"tool_choice":         "auto",
			"parallel_tool_calls": false,
			"store":               false,
			"stream":              true,
			"include":             []interface{}{},
			"client_metadata": map[string]string{
				"x-codex-installation-id": "ai-study-question-import",
			},
		}
	} else {
		requestBody = map[string]interface{}{
			"model":        model,
			"instructions": systemPrompt,
			"input":        userPrompt,
			"store":        false,
		}
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", model, credential.Source, err
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+"/responses", bytes.NewReader(body))
	if err != nil {
		return "", model, credential.Source, err
	}
	req.Header.Set("Authorization", "Bearer "+credential.Token)
	req.Header.Set("Content-Type", "application/json")
	if credential.Mode == "chatgpt_access_token" {
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("originator", "codex_cli_rs")
		req.Header.Set("User-Agent", "codex_cli_rs/0.0.0 (Linux 0.0.0; x86_64) ai-study-question-import")
		if credential.AccountID != "" {
			req.Header.Set("ChatGPT-Account-ID", credential.AccountID)
		}
	} else if credential.AccountID != "" {
		req.Header.Set("OpenAI-Account-ID", credential.AccountID)
	}
	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return "", model, credential.Source, fmt.Errorf("Codex 请求失败: mode=%s, source=%s, err=%w", credential.Mode, credential.Source, err)
	}
	defer resp.Body.Close()
	rawBody, _ := io.ReadAll(resp.Body)
	var result questionAIResponsesResponse
	_ = json.Unmarshal(rawBody, &result)
	if resp.StatusCode >= 300 {
		if result.Error != nil && strings.TrimSpace(result.Error.Message) != "" {
			return "", model, credential.Source, fmt.Errorf("Codex 状态异常: %d, mode=%s, source=%s, message=%s", resp.StatusCode, credential.Mode, credential.Source, strings.TrimSpace(result.Error.Message))
		}
		return "", model, credential.Source, fmt.Errorf("Codex 状态异常: %d, mode=%s, source=%s, body=%s", resp.StatusCode, credential.Mode, credential.Source, limitRunes(strings.TrimSpace(string(rawBody)), 400))
	}
	if credential.Mode == "chatgpt_access_token" {
		outputText, err := readQuestionAIResponsesStream(bytes.NewReader(rawBody))
		if err != nil {
			return "", model, credential.Source, err
		}
		return outputText, model, credential.Source, nil
	}
	if result.Error != nil && strings.TrimSpace(result.Error.Message) != "" {
		return "", model, credential.Source, fmt.Errorf("Codex 返回错误: %s", strings.TrimSpace(result.Error.Message))
	}
	outputText := questionAIResponseText(result)
	if strings.TrimSpace(outputText) == "" {
		return "", model, credential.Source, fmt.Errorf("Codex 返回为空")
	}
	return outputText, model, credential.Source, nil
}

func parseQuestionAIChunk(params map[string]interface{}, provider string, systemPrompt string, chunk questionAITextChunk, chunkTotal int, enableCache bool) questionAIChunkResult {
	result := questionAIChunkResult{Index: chunk.Index}
	userPrompt := questionAIUserPrompt(chunk.Text, chunk.Index, chunkTotal)
	effectiveModel := questionAIEffectiveModel(params, provider)
	cacheKey := questionAIChunkCacheKey(provider, effectiveModel, systemPrompt, userPrompt)
	fromCache := false
	var chunkFixedText, nextModel, nextSource, nextProvider string
	if enableCache {
		if cacheEntry, ok := readQuestionAIChunkCache(cacheKey); ok {
			chunkFixedText = cacheEntry.FixedText
			nextModel = cacheEntry.Model
			nextSource = cacheEntry.Source
			nextProvider = cacheEntry.Provider
			fromCache = true
		}
	}
	if !fromCache {
		var callErr error
		chunkFixedText, nextModel, nextSource, nextProvider, callErr = callQuestionAIProvider(params, provider, systemPrompt, userPrompt)
		if callErr != nil {
			result.Err = fmt.Errorf("第 %d/%d 段解析失败: %s", chunk.Index, chunkTotal, callErr.Error())
			return result
		}
	}
	chunkRows, parseErr := parseQuestionRows(chunkFixedText)
	if parseErr != nil {
		result.Err = fmt.Errorf("%s", formatQuestionAIJSONParseError(
			fmt.Sprintf("第 %d/%d 段 JSON 解析失败", chunk.Index, chunkTotal),
			chunkFixedText,
			parseErr,
		))
		return result
	}
	if enableCache && !fromCache {
		writeQuestionAIChunkCache(cacheKey, questionAIChunkCacheEntry{
			Provider:  nextProvider,
			Model:     nextModel,
			Source:    nextSource,
			FixedText: chunkFixedText,
			CreatedAt: time.Now().Format(time.RFC3339),
		})
	}
	result.Rows = chunkRows
	result.Model = nextModel
	result.Source = nextSource
	result.Provider = nextProvider
	result.Summary = questionAIChunkSummary{
		Index:           chunk.Index,
		SourceChars:     chunk.SourceChars,
		QuestionMarkers: questionAINumberMarkerCount(chunk.Text),
		RowCount:        len(chunkRows),
		FromCache:       fromCache,
	}
	return result
}

func (s *QuestionAIParseService) Result(template *config.Template, ts *templateService.TemplateService) *common.Result {
	params := template.GetParams()
	maxChars := intFromParams(params, "max_chars", defaultQuestionAIMaxChars)
	chunkChars := intFromParams(params, "chunk_chars", defaultQuestionAIChunkChars)
	maxChunks := intFromParams(params, "max_chunks", defaultQuestionAIMaxChunks)
	if _, ok := intParam(params, "chunk_chars"); !ok {
		if configChunkChars := appKeyAny("question_ai_chunk_chars"); configChunkChars != "" {
			if n, err := strconv.Atoi(configChunkChars); err == nil && n > 0 {
				chunkChars = n
			}
		}
	}
	if _, ok := intParam(params, "max_chunks"); !ok {
		if configMaxChunks := appKeyAny("question_ai_max_chunks"); configMaxChunks != "" {
			if n, err := strconv.Atoi(configMaxChunks); err == nil && n > 0 {
				maxChunks = n
			}
		}
	}
	sourceMaxChars := defaultQuestionAISourceMaxChars
	if configMax := appKeyAny("question_ai_source_max_chars"); configMax != "" {
		if n, err := strconv.Atoi(configMax); err == nil && n > 0 {
			sourceMaxChars = n
		}
	}
	if n, ok := intParam(params, "source_max_chars"); ok {
		sourceMaxChars = n
	}
	rawText := strings.TrimSpace(gocast.ToString(params["raw_text"]))
	if rawText == "" && ts.File != nil {
		text, _, err := extractTextFromUploadedFile(ts, sourceMaxChars)
		if err != nil {
			return common.NotOk(err.Error())
		}
		rawText = text
	}
	if rawText == "" {
		filePath := strings.TrimSpace(gocast.ToString(params["file_path"]))
		if filePath != "" {
			text, err := extractTextFromLocalPath(filePath, sourceMaxChars)
			if err != nil {
				return common.NotOk(err.Error())
			}
			rawText = text
		}
	}
	rawText, sourceTruncated := limitRunesWithTruncated(rawText, sourceMaxChars)

	defaults := buildAIParseDefaults(params)
	mockResponse := strings.TrimSpace(gocast.ToString(params["mock_response"]))
	fixedText := strings.TrimSpace(gocast.ToString(params["fixed_text"]))
	provider := strings.TrimSpace(gocast.ToString(params["provider"]))
	model := strings.TrimSpace(gocast.ToString(params["model"]))
	parseMode := strings.ToLower(strings.TrimSpace(gocast.ToString(params["parse_mode"])))
	enableCache := gocast.ToBool(params["enable_cache"])
	source := ""
	chunkCount := 0
	chunkConcurrency := 0
	chunkSuccessCount := 0
	chunkFailedCount := 0
	chunkSkippedCount := 0
	chunkErrors := []map[string]interface{}{}
	chunkSummaries := []questionAIChunkSummary{}
	var err error
	if mockResponse != "" {
		fixedText = mockResponse
		source = "mock_response"
	} else if fixedText == "" {
		if rawText == "" {
			return common.NotOk("raw_text 不能为空")
		}
		if provider == "" {
			provider = "codex"
		}
		systemPrompt := questionAIInstructions(defaults)
		if parseMode == "" {
			parseMode = "auto"
		}
		useChunked := parseMode == "chunked" || (parseMode == "auto" && runeCount(rawText) > maxChars)
		if !useChunked {
			promptText := rawText
			if runeCount(promptText) > maxChars {
				promptText = limitRunes(promptText, maxChars)
				sourceTruncated = true
			}
			userPrompt := questionAIUserPrompt(promptText, 1, 1)
			fixedText, model, source, provider, err = callQuestionAIProvider(params, provider, systemPrompt, userPrompt)
			chunkCount = 1
			chunkSummaries = append(chunkSummaries, questionAIChunkSummary{
				Index:           1,
				SourceChars:     runeCount(promptText),
				QuestionMarkers: questionAINumberMarkerCount(promptText),
			})
			if err != nil {
				return common.NotOk(err.Error())
			}
		} else {
			chunks, splitErr := splitQuestionAIText(rawText, chunkChars, maxChunks)
			if splitErr != nil {
				return common.NotOk(splitErr.Error())
			}
			rowGroups := make([][]map[string]interface{}, 0, len(chunks))
			concurrency := questionAIChunkConcurrency(params, provider, len(chunks))
			chunkConcurrency = concurrency
			results := make([]questionAIChunkResult, len(chunks))
			jobs := make(chan int)
			var wg sync.WaitGroup
			for worker := 0; worker < concurrency; worker++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for index := range jobs {
						result := parseQuestionAIChunk(params, provider, systemPrompt, chunks[index], len(chunks), enableCache)
						results[index] = result
					}
				}()
			}
			for index := range chunks {
				jobs <- index
			}
			close(jobs)
			wg.Wait()
			for _, result := range results {
				if result.Index == 0 {
					chunkFailedCount++
					chunkSkippedCount++
					chunkErrors = append(chunkErrors, map[string]interface{}{
						"index": 0,
						"msg":   "分段未执行",
					})
					continue
				}
				if result.Err != nil {
					chunkFailedCount++
					chunkSkippedCount++
					chunkSummaries = append(chunkSummaries, questionAIChunkSummary{
						Index:           result.Index,
						SourceChars:     chunks[result.Index-1].SourceChars,
						QuestionMarkers: questionAINumberMarkerCount(chunks[result.Index-1].Text),
						RowCount:        0,
						FromCache:       result.Summary.FromCache,
					})
					chunkErrors = append(chunkErrors, map[string]interface{}{
						"index": result.Index,
						"msg":   result.Err.Error(),
					})
					continue
				}
				chunkSuccessCount++
				if strings.TrimSpace(result.Model) != "" {
					model = result.Model
				}
				if strings.TrimSpace(result.Source) != "" {
					source = result.Source
				}
				if strings.TrimSpace(result.Provider) != "" {
					provider = result.Provider
				}
				rowGroups = append(rowGroups, result.Rows)
				chunkSummaries = append(chunkSummaries, result.Summary)
			}
			sort.SliceStable(chunkSummaries, func(i, j int) bool {
				return chunkSummaries[i].Index < chunkSummaries[j].Index
			})
			sort.SliceStable(chunkErrors, func(i, j int) bool {
				return gocast.ToInt(chunkErrors[i]["index"]) < gocast.ToInt(chunkErrors[j]["index"])
			})
			rows := mergeQuestionRows(rowGroups)
			fixedBytes, _ := json.Marshal(map[string]interface{}{"questions": rows})
			fixedText = string(fixedBytes)
			chunkCount = len(chunks)
		}
	}
	fixedText = normalizeAIJSONString(fixedText)
	rows, err := parseQuestionRows(fixedText)
	if err != nil {
		return common.NotOk(formatQuestionAIJSONParseError("JSON 解析失败", fixedText, err))
	}
	rows = normalizeQuestionAIImportRows(rows, defaults)
	fixedBytes, _ := json.Marshal(map[string]interface{}{"questions": rows})
	fixedText = string(fixedBytes)
	if chunkCount == 1 && len(chunkSummaries) == 1 && chunkSummaries[0].RowCount == 0 {
		chunkSummaries[0].RowCount = len(rows)
	}
	return common.Ok(map[string]interface{}{
		"fixed_text":          fixedText,
		"rows":                rows,
		"provider":            provider,
		"model":               model,
		"source":              source,
		"source_chars":        runeCount(rawText),
		"source_max_chars":    sourceMaxChars,
		"source_truncated":    sourceTruncated,
		"chunk_chars":         chunkChars,
		"chunk_count":         chunkCount,
		"chunk_concurrency":   chunkConcurrency,
		"chunk_success_count": chunkSuccessCount,
		"chunk_failed_count":  chunkFailedCount,
		"chunk_skipped_count": chunkSkippedCount,
		"chunk_errors":        chunkErrors,
		"chunk_summaries":     chunkSummaries,
		"row_count":           len(rows),
	}, "AI 解析完成")
}
