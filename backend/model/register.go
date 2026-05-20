package model

import (
	templateService "github.com/collect-ui/collect/src/collect/service_imp"
	utils "github.com/collect-ui/collect/src/collect/utils"

	"ai-study-admin/model/base"
)

var modelMap map[string]interface{}
var primaryKeyMap map[string][]string

type TableData struct {
	templateService.DatabaseModel
}

func init() {
	modelMap = make(map[string]interface{})
	primaryKeyMap = make(map[string][]string)

	baseTableMap, basePkMap := base.GetTable()
	for k, v := range baseTableMap {
		modelMap[k] = v
	}
	for k, v := range basePkMap {
		primaryKeyMap[k] = v
	}
}

func (*TableData) GetModel(tableName string) interface{} {
	return modelMap[tableName]
}

func (*TableData) CloneModel(tableName string) interface{} {
	return utils.Copy(modelMap[tableName])
}

func (*TableData) GetPrimaryKey(tableName string) []string {
	return primaryKeyMap[tableName]
}
