package base

import "reflect"

func GetTable() (map[string]interface{}, map[string][]string) {
	modelMap := make(map[string]interface{})
	primaryKeyMap := make(map[string][]string)

	register := func(model interface {
		TableName() string
		PrimaryKey() []string
	}) {
		modelValue := reflect.ValueOf(model)
		if modelValue.Kind() == reflect.Ptr {
			modelMap[model.TableName()] = modelValue.Elem().Interface()
		} else {
			modelMap[model.TableName()] = model
		}
		primaryKeyMap[model.TableName()] = model.PrimaryKey()
	}

	register(&Role{})

	register(&UserRoleIDList{})

	register(&SysCode{})

	register(&UserAccount{})

	register(&SysMenu{})

	register(&RoleMenu{})

	register(&SchemaPageData{})

	register(&SchemaPageField{})

	register(&SchemaPage{})

	register(&UserChangeHistory{})

	register(&SysBtn{})

	register(&BtnRoleIDList{})

	register(&QuestionItem{})

	register(&QuestionOption{})

	register(&QuestionAnswer{})

	register(&QuestionBlankAnswer{})

	register(&QuestionScoringPoint{})

	register(&QuestionGrade{})

	register(&QuestionSubject{})

	register(&QuestionUnit{})

	register(&QuestionKnowledge{})

	register(&QuestionKnowledgeRel{})

	register(&QuestionAsset{})

	register(&QuestionReviewRecord{})

	register(&QuestionChangeLog{})

	register(&QuestionImportBatch{})

	register(&QuestionImportRow{})

	return modelMap, primaryKeyMap
}
