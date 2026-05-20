package plugins

import (
	"fmt"
	"strings"
	"testing"
)

func TestNormalizeExtractedTextCleansReplacementRunsAndControls(t *testing.T) {
	raw := "答案及解析����������������（七）\x00\x16Peter"
	text := normalizeExtractedText(raw, 0)
	if text != "答案及解析\n（七）\nPeter" {
		t.Fatalf("unexpected normalized text: %q", text)
	}
	if replacementRuneCount(raw) != 16 {
		t.Fatalf("unexpected replacement count: %d", replacementRuneCount(raw))
	}
}

func TestNormalizeExtractedTextRestoresCollapsedPDFQuestionBreaks(t *testing.T) {
	raw := "小升初英语复习题三十套（含详细解析）（一）一、语法巩固1. What _____ useful dictionary it is!A. a B. an C. the D. /2. Mr. Green has little time today, _____?A. have he B. hasn’t he C. does he D. doesn’t he答案及解析1． A 解析：a后面跟以辅音开头的词。2． C 解析：反意疑问句。二、完形精练John sent for a doctor."
	text := normalizeExtractedText(raw, 0)
	for _, want := range []string{
		"（一）\n一、语法巩固\n1. What",
		"\nA. a\nB. an\nC. the\nD. /\n2. Mr. Green",
		"\n答案及解析\n1． A 解析",
		"\n2． C 解析",
		"\n二、完形精练\nJohn sent",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("normalized text missing %q:\n%s", want, text)
		}
	}
}

func TestNormalizeExtractedTextRestoresDialogueAndGluedQuestionBreaks(t *testing.T) {
	raw := "D. famous at8. --How long have you _______the dictionary?D. keep9.--_______may I keep this book?G. Is it a lovely day?W: Hi, Jack.M: ____2_____W: Yes, it is."
	text := normalizeExtractedText(raw, 0)
	for _, want := range []string{
		"D. famous at\n8. --How long",
		"D. keep\n9.--_______may",
		"G. Is it a lovely day?\nW: Hi, Jack.",
		"M: ____2_____\nW: Yes, it is.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("normalized text missing %q:\n%s", want, text)
		}
	}
}

func TestExtractedTextToHTMLEscapesContent(t *testing.T) {
	html := extractedTextToHTML("A < B\n\nC & D")
	want := "<div>A &lt; B</div><div><br></div><div>C &amp; D</div>"
	if html != want {
		t.Fatalf("unexpected html:\nwant %q\n got %q", want, html)
	}
}

func TestFormatQuestionAIJSONParseErrorIncludesRawResponse(t *testing.T) {
	err := formatQuestionAIJSONParseError("第 1/31 段 JSON 解析失败", "{\"questions\":[", fmt.Errorf("AI 返回不是合法 JSON"))
	for _, want := range []string{
		"第 1/31 段 JSON 解析失败: AI 返回不是合法 JSON",
		"原始返回字符数: 14",
		"原始返回:\n{\"questions\":[",
	} {
		if !strings.Contains(err, want) {
			t.Fatalf("error missing %q:\n%s", want, err)
		}
	}
}

func TestParseQuestionRowsRepairsMissingClosingDelimiters(t *testing.T) {
	rows, err := parseQuestionRows(`{"questions":[{"stem_text":"What?","answer_key":"A"}`)
	if err != nil {
		t.Fatalf("parseQuestionRows returned error: %v", err)
	}
	if len(rows) != 1 || rows[0]["stem_text"] != "What?" || rows[0]["answer_key"] != "A" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestParseQuestionRowsRepairsOutsideStringParenthesis(t *testing.T) {
	rows, err := parseQuestionRows(`{"questions":[{"stem_text":"Read (carefully).","choice_items":[{"sub_no":"1","question_text":"When?","answer_key":"Six months old.","analysis":"") ,{"sub_no":"2","question_text":"Why?","answer_key":"Because she was ill.","analysis":""}]}]}`)
	if err != nil {
		t.Fatalf("parseQuestionRows returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d: %#v", len(rows), rows)
	}
	if rows[0]["stem_text"] != "Read (carefully)." {
		t.Fatalf("parentheses inside strings should be preserved: %#v", rows[0]["stem_text"])
	}
	items, ok := rows[0]["choice_items"].([]interface{})
	if !ok || len(items) != 2 {
		t.Fatalf("expected repaired choice_items, got %#v", rows[0]["choice_items"])
	}
}

func TestQuestionAINumberMarkerCount(t *testing.T) {
	text := "1. One 2．Two 100. Hundred 101. skip A. option"
	if got := questionAINumberMarkerCount(text); got != 3 {
		t.Fatalf("expected 3 markers, got %d", got)
	}
}

func TestQuestionAINumberMarkerCountSkipsAnswerSection(t *testing.T) {
	text := strings.Join([]string{
		"一、语法巩固",
		"1. What _____ useful dictionary it is!",
		"2. Mr. Green has little time today, _____?",
		"答案及解析",
		"1． A 解析：a后面跟以辅音开头的词。",
		"2． C 解析：反意疑问句。",
	}, "\n")
	if got := questionAINumberMarkerCount(text); got != 2 {
		t.Fatalf("expected 2 question markers before answer section, got %d", got)
	}
}

func TestQuestionAIInstructionsRequireDialogueSplit(t *testing.T) {
	instructions := questionAIInstructions(map[string]interface{}{"subject": "english"})
	for _, want := range []string{
		"补全对话",
		"不要合并成 reading_choice",
		"每个空单独生成一条 question",
		"D. famous at8.",
		"标题、书名、卷号、章节名、栏目名不是题目",
		"小升初英语复习题三十套",
		"答案区编号单独生成 question",
		"最终 questions 数应为 6，不是 12",
	} {
		if !strings.Contains(instructions, want) {
			t.Fatalf("instructions missing %q", want)
		}
	}
}

func TestQuestionAIUserPromptRepeatsCoverageRules(t *testing.T) {
	prompt := questionAIUserPrompt("二、补全对话\n____1____\n____2____", 1, 2)
	for _, want := range []string{
		"完整性硬性要求",
		"必须输出 5 条独立 question",
		"禁止合并成 1 条 reading_choice",
		"D. keep9.",
		"小升初英语复习题三十套",
		"最终必须输出 6 条 question，不是 12 条",
		"stem_text 为“选出不同的词”的总题",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q", want)
		}
	}
}

func TestReadQuestionAIResponsesStreamFromDeltas(t *testing.T) {
	stream := strings.Join([]string{
		"event: response.output_text.delta",
		`data: {"type":"response.output_text.delta","delta":"{\"questions\":[{\"stem_text\":\"A\"}"}`,
		"",
		"event: response.output_text.delta",
		`data: {"type":"response.output_text.delta","delta":"]}"}`,
		"",
		"data: [DONE]",
		"",
	}, "\n")
	got, err := readQuestionAIResponsesStream(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("readQuestionAIResponsesStream returned error: %v", err)
	}
	if got != `{"questions":[{"stem_text":"A"}]}` {
		t.Fatalf("unexpected stream text: %q", got)
	}
}

func TestReadQuestionAIResponsesStreamUsesCompletedResponse(t *testing.T) {
	stream := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"partial"}`,
		"",
		`data: {"type":"response.completed","response":{"output":[{"content":[{"text":"{\"questions\":[{\"stem_text\":\"done\"}]}"}]}]}}`,
		"",
	}, "\n")
	got, err := readQuestionAIResponsesStream(strings.NewReader(stream))
	if err != nil {
		t.Fatalf("readQuestionAIResponsesStream returned error: %v", err)
	}
	if got != `{"questions":[{"stem_text":"done"}]}` {
		t.Fatalf("unexpected completed stream text: %q", got)
	}
}

func TestReadQuestionAIResponsesStreamError(t *testing.T) {
	stream := strings.Join([]string{
		"event: error",
		`data: {"error":{"message":"stream denied"}}`,
		"",
	}, "\n")
	_, err := readQuestionAIResponsesStream(strings.NewReader(stream))
	if err == nil || !strings.Contains(err.Error(), "stream denied") {
		t.Fatalf("expected stream error, got %v", err)
	}
}

func TestQuestionAIChunkConcurrencyByProvider(t *testing.T) {
	if got := questionAIChunkConcurrency(map[string]interface{}{}, "deepseek", 31); got != 10 {
		t.Fatalf("deepseek concurrency = %d, want 10", got)
	}
	if got := questionAIChunkConcurrency(map[string]interface{}{}, "codex", 31); got != 2 {
		t.Fatalf("codex concurrency = %d, want 2", got)
	}
	if got := questionAIChunkConcurrency(map[string]interface{}{}, "codex", 1); got != 1 {
		t.Fatalf("single chunk concurrency = %d, want 1", got)
	}
	if got := questionAIChunkConcurrency(map[string]interface{}{"concurrency": 3}, "deepseek", 31); got != 3 {
		t.Fatalf("explicit concurrency = %d, want 3", got)
	}
}

func TestQuestionAICodexModelForChatGPTAccessToken(t *testing.T) {
	credential := questionAICredential{Mode: "chatgpt_access_token"}
	if got := questionAICodexModelForCredential(map[string]interface{}{}, credential); got != defaultQuestionAIChatGPTModel {
		t.Fatalf("chatgpt access token model = %q, want %q", got, defaultQuestionAIChatGPTModel)
	}
	if got := questionAICodexModelForCredential(map[string]interface{}{"model": "custom-model"}, credential); got != "custom-model" {
		t.Fatalf("explicit model should be preserved, got %q", got)
	}
}

func TestSplitQuestionAITextByPaper(t *testing.T) {
	raw := "（一）\n一、单项选择\n1. A first question\n答案及解析\n1. A 解析：one\n（ 二 ）\nthis line should stay with previous because heading has spaces\n（二）\n一、单项选择\n1. A second question\n答案及解析\n1. B 解析：two"
	chunks, err := splitQuestionAIText(raw, 80, 10)
	if err != nil {
		t.Fatalf("splitQuestionAIText returned error: %v", err)
	}
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
	}
	if chunks[0].Index != 1 || chunks[1].Index != 2 {
		t.Fatalf("unexpected chunk indexes: %#v", chunks)
	}
	if chunks[0].SourceChars == 0 || chunks[1].SourceChars == 0 {
		t.Fatalf("expected source char counts: %#v", chunks)
	}
}

func TestSplitQuestionAITextByInlinePaperHeading(t *testing.T) {
	raw := "标题（一）题1.A答案1.A（二）题2.B答案2.B"
	chunks, err := splitQuestionAIText(raw, 18, 10)
	if err != nil {
		t.Fatalf("splitQuestionAIText returned error: %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d: %#v", len(chunks), chunks)
	}
	if !strings.Contains(chunks[0].Text, "（一）") || strings.Contains(chunks[0].Text, "（二）") {
		t.Fatalf("first chunk should contain only first paper: %#v", chunks[0])
	}
	if !strings.HasPrefix(chunks[1].Text, "（二）") {
		t.Fatalf("second chunk should start at second paper: %#v", chunks[1])
	}
}

func TestMergeQuestionRowsDedup(t *testing.T) {
	merged := mergeQuestionRows([][]map[string]interface{}{
		{
			{"question_code": "Q001", "stem_text": "first"},
			{"stem_text": "Repeated stem"},
		},
		{
			{"question_code": "Q001", "stem_text": "first duplicate"},
			{"stem_text": "Repeated   stem"},
			{"stem_text": "second"},
		},
	})
	if len(merged) != 3 {
		t.Fatalf("expected 3 merged rows, got %d: %#v", len(merged), merged)
	}
}

func TestParseQuestionRowsObject(t *testing.T) {
	rows, err := parseQuestionRows(`{"questions":[{"stem_text":"What is it?","answer_key":"A"}]}`)
	if err != nil {
		t.Fatalf("parseQuestionRows returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["answer_key"] != "A" {
		t.Fatalf("unexpected row: %#v", rows[0])
	}
}

func TestNormalizeReadingChoiceKeepsSubQuestions(t *testing.T) {
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
	rows := normalizeQuestionAIImportRows([]map[string]interface{}{
		{
			"title":             "阅读理解验收",
			"question_category": "阅读理解",
			"stem_text":         "Read the passage and answer questions.",
			"choice_items": []map[string]interface{}{
				{
					"sub_no":        "1",
					"question_text": "Where is Tom?",
					"options": map[string]interface{}{
						"A": "At school",
						"B": "At home",
						"C": "In a shop",
						"D": "In a park",
					},
					"answer":        "B",
					"analysis_text": "The passage says he is at home.",
				},
			},
		},
	}, defaults)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["question_category"] != "reading_choice" {
		t.Fatalf("unexpected category: %#v", row["question_category"])
	}
	if row["question_type"] != "single_choice" {
		t.Fatalf("unexpected type: %#v", row["question_type"])
	}
	if row["answer_text"] != "1:B" {
		t.Fatalf("unexpected answer_text: %#v", row["answer_text"])
	}
	if row["answer_key"] != "1:B" {
		t.Fatalf("unexpected answer_key: %#v", row["answer_key"])
	}
	if row["blank_count"] != 1 {
		t.Fatalf("unexpected blank_count: %#v", row["blank_count"])
	}
	ref, ok := row["reference_text"].(string)
	if !ok || ref == "" || ref == "[]" {
		t.Fatalf("reference_text was not preserved: %#v", row["reference_text"])
	}
}
