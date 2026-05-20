package plugins

import (
	"fmt"
	"strings"
	"unicode"

	common "github.com/collect-ui/collect/src/collect/common"
	config "github.com/collect-ui/collect/src/collect/config"
	templateService "github.com/collect-ui/collect/src/collect/service_imp"
	utils "github.com/collect-ui/collect/src/collect/utils"
	"github.com/mozillazg/go-pinyin"
)

type Pinyin struct {
	templateService.BaseHandler
}

func pinyinSourceText(handlerParam *config.HandlerParam, params map[string]interface{}) string {
	for _, field := range []string{handlerParam.Field, handlerParam.Value} {
		if strings.TrimSpace(field) == "" {
			continue
		}
		value := utils.RenderVar(field, params)
		if value == nil {
			continue
		}
		return strings.TrimSpace(fmt.Sprint(value))
	}
	for _, key := range []string{"text", "name", "knowledge_name", "value"} {
		value, ok := params[key]
		if ok && value != nil {
			return strings.TrimSpace(fmt.Sprint(value))
		}
	}
	return ""
}

func normalizePinyinToken(token string) string {
	token = strings.NewReplacer("ü", "v", "Ü", "v", "u:", "v", "U:", "v").Replace(token)
	var b strings.Builder
	for _, r := range token {
		if r >= 'a' && r <= 'z' {
			b.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			b.WriteRune(unicode.ToLower(r))
		} else if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func buildPinyinCode(text string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	args.Heteronym = false

	var b strings.Builder
	lastSep := false
	appendSep := func() {
		if b.Len() == 0 || lastSep {
			return
		}
		b.WriteByte('_')
		lastSep = true
	}
	appendToken := func(token string) {
		token = normalizePinyinToken(token)
		if token == "" {
			appendSep()
			return
		}
		b.WriteString(token)
		lastSep = false
	}

	for _, r := range strings.TrimSpace(text) {
		switch {
		case unicode.Is(unicode.Han, r):
			syllables := pinyin.Pinyin(string(r), args)
			if len(syllables) > 0 && len(syllables[0]) > 0 {
				appendToken(syllables[0][0])
			} else {
				appendSep()
			}
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastSep = false
		case r >= 'A' && r <= 'Z':
			b.WriteRune(unicode.ToLower(r))
			lastSep = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastSep = false
		default:
			appendSep()
		}
	}

	return strings.Trim(b.String(), "_")
}

func (p *Pinyin) HandlerData(template *config.Template, handlerParam *config.HandlerParam, ts *templateService.TemplateService) *common.Result {
	text := pinyinSourceText(handlerParam, template.GetParams())
	code := buildPinyinCode(text)
	data := map[string]interface{}{
		"text":   text,
		"pinyin": code,
		"code":   code,
	}
	return common.Ok(data, "pinyin converted")
}
