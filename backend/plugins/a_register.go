package plugins

import templateService "github.com/collect-ui/collect/src/collect/service_imp"

func GetRegisterList() []templateService.ModuleResult {
	return []templateService.ModuleResult{
		&SchemaTransfer{},
		&ToLocalFile{},
		&Pinyin{},
		&QuestionPDFTextService{},
		&QuestionAIParseService{},
	}
}
