package main

import (
	"os"
	"strings"
	"testing"
)

func TestQuestionBankKnowledgeSelectUsesKnowledgePointDataAsMultiSelect(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_bank.json")
	rootObj, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question_bank root is not an object")
	}

	initAction, ok := rootObj["initAction"].([]any)
	if !ok {
		t.Fatalf("question_bank initAction missing")
	}
	knowledgeListAction, ok := findObjectWithString(initAction, "group", "question-knowledge-list")
	if !ok {
		t.Fatalf("question-knowledge-list init action not found")
	}
	if got, want := stringField(knowledgeListAction, "api"), "post:/template_data/data?service=question.knowledge_point_query"; got != want {
		t.Fatalf("question-knowledge-list api = %q, want %q", got, want)
	}
	adapt, ok := knowledgeListAction["adapt"].(map[string]any)
	if !ok {
		t.Fatalf("question-knowledge-list adapt missing")
	}
	adaptExpr, _ := adapt["questionKnowledgeList"].(string)
	for _, want := range []string{"point_id", "point_name", "knowledge_id", "knowledge_name"} {
		if !strings.Contains(adaptExpr, want) {
			t.Fatalf("question-knowledge-list adapt should map knowledge point fields, missing %q in %q", want, adaptExpr)
		}
	}
	data, ok := knowledgeListAction["data"].(map[string]any)
	if !ok {
		t.Fatalf("question-knowledge-list data missing")
	}
	if got, want := data["unit_id"], "${questionForm.unit_id || ''}"; got != want {
		t.Fatalf("question-knowledge-list unit_id = %#v, want %q", got, want)
	}
	if got := data["pagination"]; got != false {
		t.Fatalf("question-knowledge-list pagination = %#v, want false", got)
	}
	if got := data["count"]; got != false {
		t.Fatalf("question-knowledge-list count = %#v, want false", got)
	}
	searchKnowledgeListAction, ok := findObjectWithString(initAction, "group", "knowledge-list")
	if !ok {
		t.Fatalf("knowledge-list init action not found")
	}
	if got, want := stringField(searchKnowledgeListAction, "api"), "post:/template_data/data?service=question.knowledge_point_query"; got != want {
		t.Fatalf("knowledge-list api = %q, want %q", got, want)
	}

	selectObj, ok := findQuestionKnowledgeMultiSelect(root)
	if !ok {
		t.Fatalf("question editor knowledge select with mode=multiple not found")
	}
	if got := stringField(selectObj, "options"); got != "${questionKnowledgeList}" {
		t.Fatalf("editor knowledge select options = %q, want ${questionKnowledgeList}", got)
	}
	if got := stringField(selectObj, "maxTagCount"); got != "responsive" {
		t.Fatalf("editor knowledge select maxTagCount = %q, want responsive", got)
	}
	if !objectTreeContainsString(selectObj["action"], "knowledge_id", "Array.isArray(value)") {
		t.Fatalf("editor knowledge select action should normalize selected value into an array")
	}
	if !objectTreeContainsString(selectObj["action"], "knowledge_name", "questionKnowledgeList") {
		t.Fatalf("editor knowledge select action should derive names from questionKnowledgeList")
	}
}

func TestQuestionBankTopToolbarHidesFiltersDuplicatedInSidebar(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_bank.json")
	toolbar, ok := findObjectWithString(root, "className", "qb-toolbar-search")
	if !ok {
		t.Fatalf("top toolbar search container not found")
	}

	for _, name := range []string{"question_type", "difficulty", "unit_id", "knowledge_id"} {
		if item, ok := findFormItemByName(toolbar, name); ok {
			t.Fatalf("top toolbar should not render duplicate form item %s: %#v", name, item)
		}
	}
}

func TestQuestionBankQuestionCategoryOptions(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_bank.json")
	items := collectFormItemsByName(root, "question_category")
	if len(items) < 2 {
		t.Fatalf("question_category form items = %d, want at least search and editor", len(items))
	}

	var searchItem map[string]any
	var editorItem map[string]any
	for _, item := range items {
		if stringField(item, "label") != "属性" {
			continue
		}
		values := selectOptionValues(t, item)
		if containsString(values, "") {
			searchItem = item
		} else {
			editorItem = item
		}
	}
	if searchItem == nil || editorItem == nil {
		t.Fatalf("question_category search/editor 属性 form items not found")
	}

	assertSelectOptions(t, searchItem, []string{"全部属性", "经典题型", "普通题型", "考试真题"}, []string{"", "classic", "normal", "exam"})
	assertSelectOptions(t, editorItem, []string{"经典题型", "普通题型", "考试真题"}, []string{"classic", "normal", "exam"})
}

func TestQuestionBankSaveButtonLoadingDoesNotRewriteSearchForm(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_bank.json")
	rootObj, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question_bank root is not an object")
	}
	initStore, ok := rootObj["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question_bank initStore missing")
	}
	if busy, ok := initStore["questionSaveBusy"].(bool); !ok || busy {
		t.Fatalf("questionSaveBusy initial value = %#v, want false", initStore["questionSaveBusy"])
	}

	buttons := collectButtonsWithChildren(root, "保存题目")
	if len(buttons) != 1 {
		t.Fatalf("expected 1 save button, got %d", len(buttons))
	}
	button := buttons[0]
	if got, want := stringField(button, "loading"), "${questionSaveBusy === true}"; got != want {
		t.Fatalf("save button loading = %q, want %q", got, want)
	}
	if got, want := stringField(button, "disabled"), "${questionSaveBusy === true}"; got != want {
		t.Fatalf("save button disabled = %q, want %q", got, want)
	}

	for _, service := range []string{"question.question_choice_save", "question.question_choice_update"} {
		ajax, ok := findAjaxByService(button["action"], service)
		if !ok {
			t.Fatalf("save ajax for %s not found", service)
		}
		assertAjaxTogglesSaveBusy(t, ajax, service)
	}

	visitObjects(button["action"], func(obj map[string]any) {
		if stringField(obj, "tag") == "update-form" && stringField(obj, "formName") == "questionSearchForm" {
			t.Fatalf("save action should not update questionSearchForm: %#v", obj)
		}
		if stringField(obj, "tag") == "update-store" {
			value, ok := obj["value"].(map[string]any)
			if !ok {
				return
			}
			if _, hasSearchForm := value["searchForm"]; hasSearchForm {
				t.Fatalf("save action should not rewrite searchForm: %#v", obj)
			}
		}
	})
}

func TestQuestionChoiceServicesSaveMultipleKnowledgeRelations(t *testing.T) {
	bulk := readServiceDefinition(t, "collect/question/knowledge/index.yml", "knowledge_rel_bulk_save")
	if bulk.Module != "bulk_create" {
		t.Fatalf("knowledge_rel_bulk_save module = %q, want bulk_create", bulk.Module)
	}
	if bulk.Table != "question_knowledge_rel" {
		t.Fatalf("knowledge_rel_bulk_save table = %q, want question_knowledge_rel", bulk.Table)
	}
	if bulk.ModelField != "[knowledge_rels]" {
		t.Fatalf("knowledge_rel_bulk_save model_field = %#v, want [knowledge_rels]", bulk.ModelField)
	}

	for _, serviceKey := range []string{"question_choice_save", "question_choice_update"} {
		t.Run(serviceKey, func(t *testing.T) {
			def := readServiceDefinition(t, "collect/question/question/index.yml", serviceKey)
			if _, ok := def.Params["knowledge_rels"]; !ok {
				t.Fatalf("%s params missing knowledge_rels", serviceKey)
			}
			if !handlerParamsContain(def.HandlerParams, "update_array", "foreach", "[knowledge_rels]") {
				t.Fatalf("%s should build knowledge_rels rows before saving", serviceKey)
			}
			if !handlerParamsContainService(def.HandlerParams, "question.knowledge_rel_bulk_save") {
				t.Fatalf("%s should call question.knowledge_rel_bulk_save", serviceKey)
			}
			if !handlerParamsContainService(def.HandlerParams, "question.knowledge_rel_save") {
				t.Fatalf("%s should keep single knowledge_id fallback for legacy callers", serviceKey)
			}
		})
	}

	update := readServiceDefinition(t, "collect/question/question/index.yml", "question_choice_update")
	if !handlerParamsContainService(update.HandlerParams, "question.knowledge_rel_delete_by_question") {
		t.Fatalf("question_choice_update should delete old knowledge relations before writing new ones")
	}
}

func TestQuestionDetailReturnsKnowledgeIDsAsArray(t *testing.T) {
	for _, name := range []string{
		"collect/question/question/question_detail.sql",
		"collect/question/question/question_choice_detail.sql",
	} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(name)
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			sql := string(data)
			for _, want := range []string{
				"JSON_ARRAYAGG(rel.knowledge_id)",
				"AS knowledge_id",
				"GROUP_CONCAT(rel.knowledge_name SEPARATOR '、')",
				"AS knowledge_name",
			} {
				if !strings.Contains(sql, want) {
					t.Fatalf("%s missing %q", name, want)
				}
			}
		})
	}
}

func TestQuestionQuerySearchesTaxonomyKnowledgeNames(t *testing.T) {
	for _, name := range []string{
		"collect/question/question/question_query.sql",
		"collect/question/question/question_query_count.sql",
	} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(name)
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			sql := string(data)
			for _, want := range []string{
				"LEFT JOIN question_knowledge knowledge",
				"COALESCE(NULLIF(kr.knowledge_name, ''), knowledge.knowledge_name) LIKE {{.keyword}}",
				"kr.knowledge_id = {{.knowledge_id}}",
			} {
				if !strings.Contains(sql, want) {
					t.Fatalf("%s missing %q", name, want)
				}
			}
		})
	}
}

func findQuestionKnowledgeMultiSelect(root any) (map[string]any, bool) {
	var found map[string]any
	visitObjects(root, func(obj map[string]any) {
		if found != nil {
			return
		}
		if stringField(obj, "tag") == "select" &&
			stringField(obj, "mode") == "multiple" &&
			stringField(obj, "options") == "${questionKnowledgeList}" {
			found = obj
		}
	})
	return found, found != nil
}

func findFormItemByName(root any, name string) (map[string]any, bool) {
	var found map[string]any
	visitObjects(root, func(obj map[string]any) {
		if found != nil {
			return
		}
		if stringField(obj, "tag") == "form-item" && stringField(obj, "name") == name {
			found = obj
		}
	})
	return found, found != nil
}

func collectFormItemsByName(root any, name string) []map[string]any {
	var found []map[string]any
	visitObjects(root, func(obj map[string]any) {
		if stringField(obj, "tag") == "form-item" && stringField(obj, "name") == name {
			found = append(found, obj)
		}
	})
	return found
}

func assertSelectOptions(t *testing.T, formItem map[string]any, wantLabels []string, wantValues []string) {
	t.Helper()
	selectObj, ok := findObjectWithString(formItem["children"], "tag", "select")
	if !ok {
		t.Fatalf("select not found for form item %#v", formItem)
	}
	options, ok := selectObj["options"].([]any)
	if !ok {
		t.Fatalf("select options missing: %#v", selectObj["options"])
	}
	if len(options) != len(wantLabels) {
		t.Fatalf("option count = %d, want %d", len(options), len(wantLabels))
	}
	for i, option := range options {
		row, ok := option.(map[string]any)
		if !ok {
			t.Fatalf("option %d is not object: %#v", i, option)
		}
		if got := stringField(row, "label"); got != wantLabels[i] {
			t.Fatalf("option %d label = %q, want %q", i, got, wantLabels[i])
		}
		if got := stringField(row, "value"); got != wantValues[i] {
			t.Fatalf("option %d value = %q, want %q", i, got, wantValues[i])
		}
	}
}

func selectOptionValues(t *testing.T, formItem map[string]any) []string {
	t.Helper()
	selectObj, ok := findObjectWithString(formItem["children"], "tag", "select")
	if !ok {
		t.Fatalf("select not found for form item %#v", formItem)
	}
	options, ok := selectObj["options"].([]any)
	if !ok {
		t.Fatalf("select options missing: %#v", selectObj["options"])
	}
	values := make([]string, 0, len(options))
	for _, option := range options {
		row, ok := option.(map[string]any)
		if !ok {
			t.Fatalf("option is not object: %#v", option)
		}
		values = append(values, stringField(row, "value"))
	}
	return values
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func findAjaxByService(root any, service string) (map[string]any, bool) {
	var found map[string]any
	visitObjects(root, func(obj map[string]any) {
		if found != nil || stringField(obj, "tag") != "ajax" {
			return
		}
		if strings.Contains(stringField(obj, "api"), "service="+service) {
			found = obj
		}
	})
	return found, found != nil
}

func assertAjaxTogglesSaveBusy(t *testing.T, ajax map[string]any, service string) {
	t.Helper()
	start, ok := ajax["start"].(map[string]any)
	if !ok || start["questionSaveBusy"] != true {
		t.Fatalf("%s ajax start should set questionSaveBusy=true, got %#v", service, ajax["start"])
	}
	end, ok := ajax["end"].(map[string]any)
	if !ok || end["questionSaveBusy"] != false {
		t.Fatalf("%s ajax end should set questionSaveBusy=false, got %#v", service, ajax["end"])
	}
}

func objectTreeContainsString(root any, field string, part string) bool {
	found := false
	visitObjects(root, func(obj map[string]any) {
		if found {
			return
		}
		if strings.Contains(stringField(obj, field), part) {
			found = true
		}
	})
	return found
}

func handlerParamsContain(params []map[string]any, key string, field string, value string) bool {
	for _, param := range params {
		if stringField(param, "key") == key && stringField(param, field) == value {
			return true
		}
	}
	return false
}

func handlerParamsContainService(params []map[string]any, serviceName string) bool {
	for _, param := range params {
		service, ok := param["service"].(map[string]any)
		if !ok {
			continue
		}
		if stringField(service, "service") == serviceName {
			return true
		}
	}
	return false
}
