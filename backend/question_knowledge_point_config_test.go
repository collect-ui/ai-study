package main

import (
	"os"
	"strings"
	"testing"
)

func TestQuestionKnowledgePointModuleIsWired(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_knowledge_point.json")

	initAction, ok := root.(map[string]any)["initAction"].([]any)
	if !ok {
		t.Fatalf("question_knowledge_point initAction missing")
	}
	listAction, ok := findObjectWithString(initAction, "group", "point-list")
	if !ok {
		t.Fatalf("point-list init action not found")
	}
	if got, want := stringField(listAction, "api"), "post:/template_data/data?service=question.knowledge_point_query"; got != want {
		t.Fatalf("point-list api = %q, want %q", got, want)
	}
	if _, ok := findObjectWithString(initAction, "key", "ai-study:question-knowledge-point:search-form"); !ok {
		t.Fatalf("knowledge point search form should be restored from localstore")
	}
	if _, ok := findObjectWithString(initAction, "key", "ai-study:question-knowledge-point:sider-collapsed"); !ok {
		t.Fatalf("knowledge point sider collapsed state should be restored from localstore")
	}
	initStore, ok := root.(map[string]any)["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("initStore missing")
	}
	searchForm, ok := initStore["searchForm"].(map[string]any)
	if !ok {
		t.Fatalf("searchForm initStore missing")
	}
	if got, want := stringField(searchForm, "subject"), "english"; got != want {
		t.Fatalf("searchForm subject = %q, want %q", got, want)
	}

	panelGroup, ok := findObjectWithString(root, "className", "kp-panel-scroll-root")
	if !ok {
		t.Fatalf("knowledge point panel group not found")
	}
	if got, want := stringField(panelGroup, "direction"), "horizontal"; got != want {
		t.Fatalf("panel group direction = %q, want %q", got, want)
	}

	if _, ok := findObjectWithString(root, "tag", "sider"); !ok {
		t.Fatalf("knowledge point sider not found")
	}
	if _, ok := findObjectWithString(root, "key", "ai-study:question-knowledge-point:sider-collapsed"); !ok {
		t.Fatalf("knowledge point sider collapse action should write localstore")
	}
	if _, ok := findObjectWithString(root, "className", "kp-toolbar-keyword"); !ok {
		t.Fatalf("knowledge point toolbar keyword input should use compact sizing class")
	}
	if !objectTreeContainsString(root, "css", ".kp-shell-layout > .ant-layout-sider > .ant-layout-sider-trigger") {
		t.Fatalf("knowledge point sider should use the same bottom trigger style as question bank")
	}
	if !objectTreeContainsString(root, "css", ".layout-content-inner > div > div") {
		t.Fatalf("knowledge point page should constrain the low-code wrapper height so the sider trigger remains visible")
	}

	listview, ok := findObjectWithString(root, "keyField", "point_id")
	if !ok {
		t.Fatalf("knowledge point card listview not found")
	}
	if got, want := stringField(listview, "itemData"), "${dataList}"; got != want {
		t.Fatalf("listview itemData = %q, want %q", got, want)
	}
	itemAttr, ok := listview["itemAttr"].(map[string]any)
	if !ok {
		t.Fatalf("listview itemAttr missing")
	}
	if got, want := stringField(itemAttr, "className"), "kp-card"; got != want {
		t.Fatalf("listview item class = %q, want %q", got, want)
	}

	deleteButton, ok := findObjectWithString(itemAttr, "icon", "DeleteOutlined")
	if !ok {
		t.Fatalf("card delete button not found")
	}
	ajax, ok := findObjectWithString(deleteButton["action"], "api", "post:/template_data/data?service=question.knowledge_point_delete")
	if !ok {
		t.Fatalf("card delete should call question.knowledge_point_delete")
	}
	data, ok := ajax["data"].(map[string]any)
	if !ok {
		t.Fatalf("delete ajax data missing")
	}
	if expr, ok := data["point_id_list"].(string); !ok || !strings.Contains(expr, "row.point_id") {
		t.Fatalf("point_id_list expression = %q, want row.point_id", expr)
	}

	form, ok := findObjectWithString(root, "name", "knowledgePointForm")
	if !ok {
		t.Fatalf("knowledge point editor form not found")
	}
	for _, label := range []string{"科目", "年级", "单元", "小节", "名称", "编码", "详细内容"} {
		if _, ok := findObjectWithString(form, "label", label); !ok {
			t.Fatalf("editor form missing label %q", label)
		}
	}
	contentItem, ok := findObjectWithString(form, "name", "content_detail")
	if !ok {
		t.Fatalf("editor form missing content_detail field")
	}
	knowledgeRichText, ok := findObjectWithString(contentItem, "tag", "ai-study-knowledge-rich-text")
	if !ok {
		t.Fatalf("editor form should use AI Study knowledge rich-text extension for content detail")
	}
	if _, ok := knowledgeRichText["knowledgePointApi"].(string); !ok {
		t.Fatalf("knowledge rich-text should declare fixed knowledge point query api")
	}
	if _, ok := findObjectWithString(form, "name", "mnemonic_method"); ok {
		t.Fatalf("mnemonic method should be associated inside content_detail rich text, not a separate form field")
	}
	if _, ok := findObjectWithString(form, "name", "antonyms"); ok {
		t.Fatalf("antonyms should be associated inside content_detail rich text, not a separate form field")
	}
	if _, ok := findObjectWithString(root, "tag", "rich-text-render"); !ok {
		t.Fatalf("knowledge point cards should render content detail through rich-text-render")
	}
	for _, service := range []string{"question.knowledge_point_save", "question.knowledge_point_update"} {
		if _, ok := findObjectWithString(root, "api", "post:/template_data/data?service="+service); !ok {
			t.Fatalf("editor save action should call %s", service)
		}
	}

	for _, name := range []string{"keyword", "subject", "grade", "unit_id", "section_id", "status"} {
		if _, ok := findObjectWithString(root, "name", name); !ok {
			t.Fatalf("page missing filter/editor field %s", name)
		}
	}
}

func TestQuestionKnowledgePointServicesAndMigration(t *testing.T) {
	serviceList, err := os.ReadFile("collect/question/service.yml")
	if err != nil {
		t.Fatalf("read question service.yml: %v", err)
	}
	if !strings.Contains(string(serviceList), "knowledge_point/index.yml") {
		t.Fatalf("question service.yml should include knowledge_point/index.yml")
	}

	save := readServiceDefinition(t, "collect/question/knowledge_point/index.yml", "knowledge_point_save")
	if save.Module != "model_save" || save.Table != "question_knowledge_point" {
		t.Fatalf("knowledge_point_save = module %q table %q, want model_save question_knowledge_point", save.Module, save.Table)
	}
	for _, key := range []string{"point_id", "subject", "stage", "grade", "unit_id", "section_id", "section_name", "point_code", "point_name", "content_detail"} {
		if _, ok := save.Params[key]; !ok {
			t.Fatalf("knowledge_point_save params missing %s", key)
		}
	}

	update := readServiceDefinition(t, "collect/question/knowledge_point/index.yml", "knowledge_point_update")
	for _, key := range []string{"subject", "stage", "grade", "unit_id", "section_id", "point_code", "point_name", "content_detail"} {
		if !stringSliceContains(update.UpdateFields, key) {
			t.Fatalf("knowledge_point_update update_fields missing %s", key)
		}
	}
	del := readServiceDefinition(t, "collect/question/knowledge_point/index.yml", "knowledge_point_delete")
	if del.Module != "model_delete" || del.Filter["point_id__in"] != "[point_id_list]" {
		t.Fatalf("knowledge_point_delete should delete by point_id_list, got module=%q filter=%#v", del.Module, del.Filter)
	}

	sqlData, err := os.ReadFile("collect/question/knowledge_point/knowledge_point_query.sql")
	if err != nil {
		t.Fatalf("read knowledge_point_query.sql: %v", err)
	}
	sql := string(sqlData)
	for _, want := range []string{"question_knowledge_point", "question_section", "section_id", "content_detail", "a.point_id = {{.point_id}}", "LIMIT {{.start}}, {{.size}}"} {
		if !strings.Contains(sql, want) {
			t.Fatalf("knowledge_point_query.sql missing %q", want)
		}
	}
	for _, forbidden := range []string{"mnemonic_method", "antonyms"} {
		if strings.Contains(sql, forbidden) {
			t.Fatalf("knowledge_point_query.sql should not use separate field %q", forbidden)
		}
	}

	migrationData, err := os.ReadFile("scripts/migrate-question-web.sql")
	if err != nil {
		t.Fatalf("read migrate-question-web.sql: %v", err)
	}
	migration := string(migrationData)
	for _, want := range []string{"CREATE TABLE IF NOT EXISTS question_knowledge_point", "idx_question_knowledge_point_query", "menu-question-knowledge-point", "frontend.question_knowledge_point"} {
		if !strings.Contains(migration, want) {
			t.Fatalf("migrate-question-web.sql missing %q", want)
		}
	}
	for _, forbidden := range []string{"mnemonic_method", "antonyms"} {
		if strings.Contains(migration, forbidden) {
			t.Fatalf("migrate-question-web.sql should not add separate field %q", forbidden)
		}
	}
}
