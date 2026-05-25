package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type batchDeleteSpec struct {
	name        string
	tabKey      string
	panelTitle  string
	selection   string
	rowID       string
	listKey     string
	service     string
	serviceFile string
	filterKey   string
}

func TestQuestionTaxonomyBatchDeleteButtonsAreScoped(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	specs := []batchDeleteSpec{
		{
			name:        "年级",
			tabKey:      "grade",
			selection:   "gradeSelection",
			rowID:       "grade_id",
			listKey:     "grade_id_list",
			service:     "question.grade_delete",
			serviceFile: "collect/question/grade/index.yml",
			filterKey:   "grade_id__in",
		},
		{
			name:        "科目",
			tabKey:      "subject",
			selection:   "subjectSelection",
			rowID:       "subject_id",
			listKey:     "subject_id_list",
			service:     "question.subject_delete",
			serviceFile: "collect/question/subject/index.yml",
			filterKey:   "subject_id__in",
		},
		{
			name:        "课本目录结构",
			panelTitle:  "课本目录结构",
			selection:   "knowledgeSelection",
			rowID:       "knowledge_id",
			listKey:     "knowledge_id_list",
			service:     "question.knowledge_delete",
			serviceFile: "collect/question/knowledge/index.yml",
			filterKey:   "knowledge_id__in",
		},
	}

	allBulkButtons := collectButtonsWithChildren(root, "批量删除")
	if len(allBulkButtons) != len(specs) {
		t.Fatalf("expected %d bulk delete buttons, got %d", len(specs), len(allBulkButtons))
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			container := taxonomyContainer(t, root, spec)
			if countButtonsWithChildren(container["searchToolBar"], "批量删除") != 0 {
				t.Fatalf("%s search toolbar should not contain bulk delete buttons", spec.name)
			}

			buttons := collectButtonsWithChildren(container["topRight"], "批量删除")
			if len(buttons) != 1 {
				t.Fatalf("%s topRight expected 1 bulk delete button, got %d", spec.name, len(buttons))
			}
			assertBatchDeleteButton(t, buttons[0], spec)

			table, ok := findObjectWithString(container["children"], "tag", "table")
			if !ok {
				t.Fatalf("%s table not found", spec.name)
			}
			if got := stringField(table, "selection"); got != "${"+spec.selection+"}" {
				t.Fatalf("%s table selection = %q, want ${%s}", spec.name, got, spec.selection)
			}
			if got := stringField(table, "rowSelection"); got != "multiple" {
				t.Fatalf("%s table rowSelection = %q, want multiple", spec.name, got)
			}
		})
	}
}

func TestQuestionTaxonomyUnitTabUsesSubjectGradeTree(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	rootMap, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy root should be an object")
	}
	initStore, ok := rootMap["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy initStore missing")
	}
	if _, exists := initStore["selectedUnitNode"]; !exists {
		t.Fatalf("initStore missing selectedUnitNode")
	}

	container := taxonomyContainer(t, root, batchDeleteSpec{name: "单元", tabKey: "unit"})
	if countButtonsWithChildren(container["topRight"], "新增单元") != 1 {
		t.Fatalf("unit tab should render one 新增单元 button")
	}
	editButtons := collectButtonsWithChildren(container["topRight"], "编辑单元")
	if len(editButtons) != 1 {
		t.Fatalf("unit tab should render one 编辑单元 button, got %d", len(editButtons))
	}
	if got := stringField(editButtons[0], "disabled"); !strings.Contains(got, "selectedUnitNode") || !strings.Contains(got, "unit") {
		t.Fatalf("编辑单元 disabled expression = %q, want selected unit guard", got)
	}

	deleteButtons := collectButtonsWithChildren(container["topRight"], "删除单元")
	if len(deleteButtons) != 1 {
		t.Fatalf("unit tab should render one 删除单元 button, got %d", len(deleteButtons))
	}
	deleteButton := deleteButtons[0]
	if got := stringField(deleteButton, "disabled"); !strings.Contains(got, "selectedUnitNode") || !strings.Contains(got, "unit") {
		t.Fatalf("删除单元 disabled expression = %q, want selected unit guard", got)
	}
	ajax, ok := findObjectWithString(deleteButton["action"], "api", "post:/template_data/data?service=question.unit_delete")
	if !ok {
		t.Fatalf("删除单元 should call question.unit_delete")
	}
	data, ok := ajax["data"].(map[string]any)
	if !ok {
		t.Fatalf("删除单元 ajax data missing")
	}
	unitIDExpr, ok := data["unit_id_list"].(string)
	if !ok || !strings.Contains(unitIDExpr, "selectedUnitNode") || !strings.Contains(unitIDExpr, "unit_id") {
		t.Fatalf("unit_id_list expression = %q, want selectedUnitNode unit_id", unitIDExpr)
	}

	if _, ok := findObjectWithString(container["children"], "tag", "table"); ok {
		t.Fatalf("unit tab should use a tree instead of a table")
	}
	tree, ok := findObjectWithString(container["children"], "tag", "tree")
	if !ok {
		t.Fatalf("unit tab tree not found")
	}
	if got := stringField(tree, "treeSelected"); got != "selectedUnitNode" {
		t.Fatalf("unit treeSelected = %q, want selectedUnitNode", got)
	}
	if got := stringField(tree, "selectedKeys"); got != "" {
		t.Fatalf("unit selectedKeys = %q, want empty so Ant Tree keeps its internal selected state", got)
	}
	treeData := stringField(tree, "treeData")
	for _, want := range []string{"subjectList", "gradeList", "unitList", "grade_label", "grade_semester_name", "subject_", "grade_", "unit_"} {
		if !strings.Contains(treeData, want) {
			t.Fatalf("unit treeData missing %q", want)
		}
	}
	selectAction, ok := findObjectWithString(tree["selectAction"], "tag", "update-store")
	if !ok {
		t.Fatalf("unit tree should update store on select")
	}
	selectValue, ok := selectAction["value"].(map[string]any)
	if !ok {
		t.Fatalf("unit tree select update-store value missing")
	}
	selectedUnitExpr, ok := selectValue["selectedUnitNode"].(string)
	if !ok || !strings.Contains(selectedUnitExpr, "selected ? node") {
		t.Fatalf("selectedUnitNode expression = %q, want selected node sync", selectedUnitExpr)
	}
	unitSelectionExpr, ok := selectValue["unitSelection"].(string)
	if !ok || !strings.Contains(unitSelectionExpr, "node.type === 'unit'") {
		t.Fatalf("unitSelection expression = %q, want unit-only selection", unitSelectionExpr)
	}
}

func TestQuestionTaxonomyUnitDialogSplitsUnitNumberAndTitle(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	rootMap, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy root should be an object")
	}
	initStore, ok := rootMap["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy initStore missing")
	}
	options, ok := initStore["unitNumberOptions"].([]any)
	if !ok || len(options) < 20 {
		t.Fatalf("unitNumberOptions should include Unit 1-20 mappings, got %#v", initStore["unitNumberOptions"])
	}

	dialog, ok := findObjectWithString(root, "open", "${unitDialogVisible}")
	if !ok {
		t.Fatalf("unit dialog not found")
	}
	unitNoItem, ok := findObjectWithString(dialog, "name", "unit_no")
	if !ok {
		t.Fatalf("unit dialog missing unit_no field")
	}
	if got := stringField(unitNoItem, "label"); got != "单元" {
		t.Fatalf("unit_no label = %q, want 单元", got)
	}
	unitNoSelect, ok := findObjectWithString(unitNoItem["children"], "tag", "select")
	if !ok {
		t.Fatalf("unit_no should render a select")
	}
	unitNoOptions := stringField(unitNoSelect, "options")
	for _, want := range []string{"unitNumberOptions", "english_label", "cn_label"} {
		if !strings.Contains(unitNoOptions, want) {
			t.Fatalf("unit_no options missing %q: %s", want, unitNoOptions)
		}
	}
	if _, ok := findObjectWithString(dialog, "name", "unit_title"); !ok {
		t.Fatalf("unit dialog missing unit_title field")
	}
	unitNameItem, ok := findObjectWithString(dialog, "name", "unit_name")
	if !ok {
		t.Fatalf("unit dialog should keep hidden unit_name for persistence")
	}
	if hidden, ok := unitNameItem["hidden"].(bool); !ok || !hidden {
		t.Fatalf("unit_name should be hidden, got %#v", unitNameItem["hidden"])
	}

	dialogJSON, err := json.Marshal(dialog)
	if err != nil {
		t.Fatalf("marshal unit dialog: %v", err)
	}
	for _, want := range []string{"unit_title", "unit_name", "Unit ", "cn_label", "unitAutoCodeLast"} {
		if !strings.Contains(string(dialogJSON), want) {
			t.Fatalf("unit dialog config missing %q", want)
		}
	}
}

func TestQuestionTaxonomyAutoCodesAllowManualEditsOnAdd(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	data, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal question taxonomy: %v", err)
	}
	config := string(data)
	for _, want := range []string{
		"gradePinyinLastCode",
		"subjectPinyinLastCode",
		"sectionPinyinLastCode",
		"knowledgePinyinLastCode",
		"String(grade_code || '').trim() === (gradePinyinLastCode || '')",
		"String(subject_code || '').trim() === (subjectPinyinLastCode || '')",
		"String(section_code || '').trim() === (sectionPinyinLastCode || '')",
		"String(knowledge_code || '').trim() === (knowledgePinyinLastCode || '')",
	} {
		if !strings.Contains(config, want) {
			t.Fatalf("auto-code manual edit guard missing %q", want)
		}
	}
	if strings.Contains(config, "unitPinyinTypingName") || strings.Contains(config, "unitPinyinResult") {
		t.Fatalf("unit dialog should not use name-driven pinyin auto-fill that can overwrite manual code edits")
	}
}

func TestQuestionTaxonomySectionTabUsesIndependentCRUD(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	rootMap, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy root should be an object")
	}
	initStore, ok := rootMap["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy initStore missing")
	}
	for _, key := range []string{"sectionSearchForm", "selectedSectionNode", "sectionDialogVisible", "sectionOp", "editingSectionID", "sectionSelection", "sectionList", "sectionForm"} {
		if _, exists := initStore[key]; !exists {
			t.Fatalf("initStore missing %s", key)
		}
	}
	loader, ok := findObjectWithString(rootMap["initAction"], "group", "section-list")
	if !ok {
		t.Fatalf("initAction missing section-list loader")
	}
	if got := stringField(loader, "api"); got != "post:/template_data/data?service=question.section_query" {
		t.Fatalf("section-list api = %q, want question.section_query", got)
	}
	if got := stringField(loader, "appendFields"); got != "${sectionSearchForm}" {
		t.Fatalf("section-list appendFields = %q, want ${sectionSearchForm}", got)
	}

	container := taxonomyContainer(t, root, batchDeleteSpec{name: "小节", tabKey: "section"})
	if countButtonsWithChildren(container["topRight"], "新增小节") != 1 {
		t.Fatalf("section tab should render one 新增小节 button")
	}
	editButtons := collectButtonsWithChildren(container["topRight"], "编辑小节")
	if len(editButtons) != 1 {
		t.Fatalf("section tab should render one 编辑小节 button, got %d", len(editButtons))
	}
	if got := stringField(editButtons[0], "disabled"); !strings.Contains(got, "selectedSectionNode") || !strings.Contains(got, "section") {
		t.Fatalf("编辑小节 disabled expression = %q, want selected section guard", got)
	}

	deleteButtons := collectButtonsWithChildren(container["topRight"], "删除小节")
	if len(deleteButtons) != 1 {
		t.Fatalf("section tab should render one 删除小节 button, got %d", len(deleteButtons))
	}
	deleteButton := deleteButtons[0]
	if got := stringField(deleteButton, "disabled"); !strings.Contains(got, "selectedSectionNode") || !strings.Contains(got, "section") {
		t.Fatalf("删除小节 disabled expression = %q, want selected section guard", got)
	}
	ajax, ok := findObjectWithString(deleteButton["action"], "api", "post:/template_data/data?service=question.section_delete")
	if !ok {
		t.Fatalf("删除小节 should call question.section_delete")
	}
	data, ok := ajax["data"].(map[string]any)
	if !ok {
		t.Fatalf("删除小节 ajax data missing")
	}
	sectionIDExpr, ok := data["section_id_list"].(string)
	if !ok || !strings.Contains(sectionIDExpr, "selectedSectionNode") || !strings.Contains(sectionIDExpr, "section_id") {
		t.Fatalf("section_id_list expression = %q, want selectedSectionNode section_id", sectionIDExpr)
	}

	tree, ok := findObjectWithString(container["children"], "tag", "tree")
	if !ok {
		t.Fatalf("section tab tree not found")
	}
	if got := stringField(tree, "treeSelected"); got != "selectedSectionNode" {
		t.Fatalf("section treeSelected = %q, want selectedSectionNode", got)
	}
	treeData := stringField(tree, "treeData")
	for _, want := range []string{"subjectList", "gradeList", "allUnitList", "sectionList", "grade_label", "unit_", "section_", "section_id", "section_name"} {
		if !strings.Contains(treeData, want) {
			t.Fatalf("section treeData missing %q", want)
		}
	}
	if !strings.Contains(treeData, "ProfileOutlined") {
		t.Fatalf("section tree should render section leaf icons")
	}
	if _, ok := findObjectWithString(tree["dblAction"], "formName", "sectionForm"); !ok {
		t.Fatalf("section tree should open sectionForm on double click")
	}

	dialog, ok := findObjectWithString(root, "open", "${sectionDialogVisible}")
	if !ok {
		t.Fatalf("section dialog not found")
	}
	for _, label := range []string{"年级", "科目", "单元", "名称", "编码", "排序"} {
		if _, ok := findObjectWithString(dialog, "label", label); !ok {
			t.Fatalf("section dialog missing label %q", label)
		}
	}
	for _, service := range []string{"question.section_save", "question.section_update"} {
		if _, ok := findObjectWithString(dialog, "api", "post:/template_data/data?service="+service); !ok {
			t.Fatalf("section dialog should call %s", service)
		}
	}

	save := readServiceDefinition(t, "collect/question/section/index.yml", "section_save")
	if save.Module != "model_save" || save.Table != "question_section" {
		t.Fatalf("section_save = module %q table %q, want model_save question_section", save.Module, save.Table)
	}
	for _, key := range []string{"section_id", "subject", "stage", "grade", "unit_id", "section_code", "section_name"} {
		if _, ok := save.Params[key]; !ok {
			t.Fatalf("section_save params missing %s", key)
		}
	}
	update := readServiceDefinition(t, "collect/question/section/index.yml", "section_update")
	for _, key := range []string{"subject", "stage", "grade", "unit_id", "section_code", "section_name"} {
		if !stringSliceContains(update.UpdateFields, key) {
			t.Fatalf("section_update update_fields missing %s", key)
		}
	}

	sqlData, err := os.ReadFile("collect/question/section/section_query.sql")
	if err != nil {
		t.Fatalf("read section_query.sql: %v", err)
	}
	sql := string(sqlData)
	for _, want := range []string{"question_section", "question_unit", "grade_label", "grade_semester_name", "section_code", "section_name"} {
		if !strings.Contains(sql, want) {
			t.Fatalf("section_query.sql missing %q", want)
		}
	}

	migrationData, err := os.ReadFile("scripts/migrate-question-web.sql")
	if err != nil {
		t.Fatalf("read migrate-question-web.sql: %v", err)
	}
	migration := string(migrationData)
	for _, want := range []string{"CREATE TABLE IF NOT EXISTS question_section", "idx_question_section_query", "'part_1'", "'Part 2'"} {
		if !strings.Contains(migration, want) {
			t.Fatalf("migrate-question-web.sql missing %q", want)
		}
	}
}

func TestQuestionTaxonomyGradesIncludeSemesters(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	container := taxonomyContainer(t, root, batchDeleteSpec{name: "单元", tabKey: "unit"})
	gradeSelect, ok := findObjectWithString(container["searchToolBar"], "name", "grade")
	if !ok {
		t.Fatalf("unit search grade form item not found")
	}
	selectObj, ok := findObjectWithString(gradeSelect["children"], "tag", "select")
	if !ok {
		t.Fatalf("unit grade select not found")
	}
	fieldNameMap, ok := selectObj["fieldNames"].(map[string]any)
	if !ok || stringField(fieldNameMap, "label") != "grade_label" {
		t.Fatalf("unit grade select label field = %#v, want grade_label", selectObj["fieldNames"])
	}

	for _, file := range []string{
		"collect/question/unit/unit_query.sql",
		"collect/question/section/section_query.sql",
		"collect/question/knowledge/knowledge_query.sql",
		"collect/question/question/question_query.sql",
	} {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		sql := string(data)
		for _, want := range []string{"grade_label", "grade_semester_name"} {
			if !strings.Contains(sql, want) {
				t.Fatalf("%s missing %q", file, want)
			}
		}
	}
	gradeSQL, err := os.ReadFile("collect/question/grade/grade_query.sql")
	if err != nil {
		t.Fatalf("read grade_query.sql: %v", err)
	}
	if !strings.Contains(string(gradeSQL), "grade_label") || !strings.Contains(string(gradeSQL), "semester_name") {
		t.Fatalf("grade_query.sql should expose semester_name and grade_label")
	}

	migrationData, err := os.ReadFile("scripts/migrate-question-web.sql")
	if err != nil {
		t.Fatalf("read migrate-question-web.sql: %v", err)
	}
	migration := string(migrationData)
	for _, want := range []string{"'grade_1_lower'", "'grade_7_lower'", "'grade_12_lower'", "'高三下学期'"} {
		if !strings.Contains(migration, want) {
			t.Fatalf("migrate-question-web.sql missing %q", want)
		}
	}
}

func TestQuestionTaxonomyDeleteServicesAcceptIDLists(t *testing.T) {
	specs := []batchDeleteSpec{
		{name: "年级", service: "question.grade_delete", serviceFile: "collect/question/grade/index.yml", listKey: "grade_id_list", filterKey: "grade_id__in"},
		{name: "科目", service: "question.subject_delete", serviceFile: "collect/question/subject/index.yml", listKey: "subject_id_list", filterKey: "subject_id__in"},
		{name: "单元", service: "question.unit_delete", serviceFile: "collect/question/unit/index.yml", listKey: "unit_id_list", filterKey: "unit_id__in"},
		{name: "小节", service: "question.section_delete", serviceFile: "collect/question/section/index.yml", listKey: "section_id_list", filterKey: "section_id__in"},
		{name: "课本目录结构", service: "question.knowledge_delete", serviceFile: "collect/question/knowledge/index.yml", listKey: "knowledge_id_list", filterKey: "knowledge_id__in"},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			service := strings.TrimPrefix(spec.service, "question.")
			def := readServiceDefinition(t, spec.serviceFile, service)
			if def.Module != "model_delete" {
				t.Fatalf("%s module = %q, want model_delete", service, def.Module)
			}
			if _, ok := def.Params[spec.listKey]; !ok {
				t.Fatalf("%s params missing %s", service, spec.listKey)
			}
			if got, want := def.Filter[spec.filterKey], "["+spec.listKey+"]"; got != want {
				t.Fatalf("%s filter %s = %q, want %q", service, spec.filterKey, got, want)
			}
		})
	}
}

func TestQuestionTaxonomyPDFKnowledgeImportEntryIsScopedToKnowledgePanel(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	container := taxonomyContainer(t, root, batchDeleteSpec{name: "课本目录结构", panelTitle: "课本目录结构"})
	rootMap, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy root should be an object")
	}
	initStore, ok := rootMap["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy initStore missing")
	}
	if form, ok := initStore["pdfKnowledgeImportForm"].(map[string]any); ok {
		if _, exists := form["unit_id"]; exists {
			t.Fatalf("PDF import initStore should not include unit_id")
		}
	}

	if countButtonsWithChildren(root, "PDF导入目录") != 1 {
		t.Fatalf("question-taxonomy should render exactly one PDF导入目录 button")
	}
	buttons := collectButtonsWithChildren(container["topRight"], "PDF导入目录")
	if len(buttons) != 1 {
		t.Fatalf("课本目录结构 topRight expected PDF导入目录 button, got %d", len(buttons))
	}
	if got := stringField(buttons[0], "icon"); got != "FilePdfOutlined" {
		t.Fatalf("PDF import icon = %q, want FilePdfOutlined", got)
	}
	if _, ok := findObjectWithString(buttons[0]["action"], "formName", "pdfKnowledgeImportForm"); !ok {
		t.Fatalf("PDF import button should initialize pdfKnowledgeImportForm")
	}
	updateStore, ok := findObjectWithString(buttons[0]["action"], "tag", "update-store")
	if !ok {
		t.Fatalf("PDF import button update-store action not found")
	}
	updateValue, ok := updateStore["value"].(map[string]any)
	if !ok {
		t.Fatalf("PDF import button update-store value missing")
	}
	formInit, _ := updateValue["pdfKnowledgeImportForm"].(string)
	if strings.Contains(formInit, "unit_id") {
		t.Fatalf("PDF import form should not prefill unit_id: %s", formInit)
	}
	dialog, ok := findObjectWithString(root, "title", "PDF目录导入")
	if !ok {
		t.Fatalf("PDF目录导入 dialog not found")
	}
	if got := stringField(dialog, "open"); got != "${pdfKnowledgeImportVisible}" {
		t.Fatalf("PDF dialog open binding = %q, want ${pdfKnowledgeImportVisible}", got)
	}
	if _, ok := findObjectWithString(dialog, "api", "post:/template_data/data?service=question.pdf_knowledge_import_one_click"); !ok {
		t.Fatalf("PDF dialog should call question.pdf_knowledge_import_one_click")
	}
	if _, ok := findObjectWithString(dialog, "label", "目标单元"); ok {
		t.Fatalf("PDF dialog should not render a target unit field; unit must be detected from the PDF")
	}
	for _, label := range []string{"PDF文件", "教材", "导入模式", "题目处理", "来源保留", "最大字符"} {
		if _, ok := findObjectWithString(dialog, "label", label); ok {
			t.Fatalf("PDF dialog should not render advanced field %q in the simplified import flow", label)
		}
	}
	for _, text := range []string{"一键生成知识点", "只预览", "上传PDF并生成", "按本文件过滤来源", "查看失败复核"} {
		if _, ok := findObjectWithString(dialog, "children", text); ok {
			t.Fatalf("PDF dialog should not render old action %q in the simplified import flow", text)
		}
	}
	if _, ok := findObjectWithString(dialog, "label", "PDF路径"); !ok {
		t.Fatalf("PDF dialog should render the temporary PDF path text input")
	}
	if _, ok := findObjectWithString(dialog, "children", "一键导入全部"); !ok {
		t.Fatalf("PDF dialog should render the one-click full import button")
	}
	if _, ok := findObjectWithString(dialog, "tag", "upload"); ok {
		t.Fatalf("PDF dialog should use a temporary text input instead of upload")
	}
	visitObjects(dialog, func(obj map[string]any) {
		if stringField(obj, "api") != "post:/template_data/data?service=question.pdf_knowledge_import_one_click" {
			return
		}
		data, ok := obj["data"].(map[string]any)
		if !ok {
			return
		}
		if _, exists := data["unit_id"]; exists {
			t.Fatalf("PDF import request should not send unit_id from the form")
		}
		if _, exists := data["file_path"]; !exists {
			t.Fatalf("PDF import request should send file_path from the temporary text input")
		}
		if got := data["auto_commit"]; got != true {
			t.Fatalf("PDF import request auto_commit = %v, want true", got)
		}
		if got := data["question_policy"]; got != "draft" {
			t.Fatalf("PDF import request question_policy = %v, want draft", got)
		}
	})
}

func TestQuestionPDFKnowledgeImportServicesAreConfigured(t *testing.T) {
	importFile := "collect/question/import/index.yml"
	oneClick := readServiceDefinition(t, importFile, "pdf_knowledge_import_one_click")
	if oneClick.Module != "question_pdf_knowledge_import" {
		t.Fatalf("pdf_knowledge_import_one_click module = %q, want question_pdf_knowledge_import", oneClick.Module)
	}
	for _, key := range []string{"file_path", "subject", "stage", "grade", "textbook_version", "auto_commit", "question_policy", "max_chars", "provider"} {
		if _, ok := oneClick.Params[key]; !ok {
			t.Fatalf("pdf_knowledge_import_one_click params missing %s", key)
		}
	}
	if _, ok := oneClick.Params["unit_id"]; ok {
		t.Fatalf("pdf_knowledge_import_one_click should not accept unit_id; Unit must come from AI JSON")
	}
	preview := readServiceDefinition(t, importFile, "pdf_import_preview")
	if preview.Module != "question_pdf_knowledge_import" {
		t.Fatalf("pdf_import_preview module = %q, want question_pdf_knowledge_import", preview.Module)
	}
	if _, ok := preview.Params["unit_id"]; ok {
		t.Fatalf("pdf_import_preview should not accept unit_id; Unit must come from AI JSON")
	}
	for _, key := range []string{"pdf_source_query", "pdf_snapshot_query", "pdf_source_trace_query", "pdf_parse_issue_query", "pdf_field_source_save"} {
		def := readServiceDefinition(t, importFile, key)
		if def.Module == "" {
			t.Fatalf("%s module should be configured", key)
		}
	}
}

func TestQuestionKnowledgeListShowsContentWithoutPDFSourceColumn(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	list, ok := findObjectWithString(root, "title", "目录结构列表")
	if !ok {
		t.Fatalf("目录结构列表 layout not found")
	}
	table, ok := findObjectWithString(list["children"], "tag", "table")
	if !ok {
		t.Fatalf("目录结构列表 table not found")
	}
	col, ok := findObjectWithString(table["columnDefs"], "field", "content_detail")
	if !ok {
		t.Fatalf("catalog table missing content_detail column")
	}
	if got := stringField(col, "headerName"); got != "明细" {
		t.Fatalf("content_detail header = %q, want 明细", got)
	}
	if _, ok := findObjectWithString(table["columnDefs"], "field", "pdf_source_excerpt"); ok {
		t.Fatalf("catalog table should not render PDF截取 column")
	}
	nameCol, ok := findObjectWithString(table["columnDefs"], "field", "knowledge_name")
	if !ok {
		t.Fatalf("catalog table missing knowledge_name column")
	}
	if got := stringField(nameCol, "headerName"); got != "名称" {
		t.Fatalf("knowledge_name header = %q, want 名称", got)
	}
	if countObjectsWithString(root, "tag", "panel-resize") != 2 {
		t.Fatalf("question-taxonomy should render two panel resize handles")
	}
	for _, id := range []string{"taxonomy-master-panel", "taxonomy-knowledge-panel", "taxonomy-knowledge-unit-tree-panel", "taxonomy-knowledge-list-panel"} {
		panel, ok := findObjectWithString(root, "id", id)
		if !ok {
			t.Fatalf("%s panel not found", id)
		}
		if got := numberField(panel, "minSize"); got != 0 {
			t.Fatalf("%s minSize = %v, want 0", id, got)
		}
	}

	data, err := os.ReadFile("collect/question/knowledge/knowledge_query.sql")
	if err != nil {
		t.Fatalf("read knowledge_query.sql: %v", err)
	}
	sql := string(data)
	for _, want := range []string{
		"question_knowledge_content",
		"question_source_field_rel",
		"question_source_block",
		"AS content_detail",
		"AS pdf_source_image_url",
		"AS pdf_source_excerpt",
		"OR content_summary.content_detail LIKE {{.keyword}}",
		"OR source_summary.pdf_source_excerpt LIKE {{.keyword}}",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("knowledge_query.sql missing %q", want)
		}
	}
}

func TestQuestionTaxonomyKnowledgeTypeCategoryConfig(t *testing.T) {
	root := readJSONFile(t, "collect/frontend/page_data/data/question/question_taxonomy.json")
	rootMap, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy root should be an object")
	}
	initStore, ok := rootMap["initStore"].(map[string]any)
	if !ok {
		t.Fatalf("question-taxonomy initStore missing")
	}
	for _, key := range []string{"knowledgeTypeManageVisible", "knowledgeTypeEditVisible", "knowledgeTypeSubject", "selectedKnowledgeTypeCode", "knowledgeTypeList", "knowledgeCategoryList", "knowledgeTypeForm"} {
		if _, exists := initStore[key]; !exists {
			t.Fatalf("initStore missing %s", key)
		}
	}
	knowledgeForm, ok := initStore["knowledgeForm"].(map[string]any)
	if !ok {
		t.Fatalf("knowledgeForm missing")
	}
	for _, key := range []string{"knowledge_type", "knowledge_category", "content_detail"} {
		if _, exists := knowledgeForm[key]; !exists {
			t.Fatalf("knowledgeForm missing %s", key)
		}
	}

	if countButtonsWithChildren(root, "类型维护") != 1 {
		t.Fatalf("question-taxonomy should render exactly one 类型维护 button")
	}
	manageDialog, ok := findObjectWithString(root, "title", "类型/具体分类维护")
	if !ok {
		t.Fatalf("类型/具体分类维护 dialog not found")
	}
	for _, text := range []string{"新增类型", "全部", "新增分类"} {
		if countButtonsWithChildren(manageDialog, text) != 1 {
			t.Fatalf("type manage dialog should render %s button", text)
		}
	}
	if _, ok := findObjectWithString(manageDialog, "rowClickAction", ""); !ok {
		// rowClickAction is an array, so fall back to checking the serialized dialog below.
		manageJSON, err := json.Marshal(manageDialog)
		if err != nil {
			t.Fatalf("marshal type manage dialog: %v", err)
		}
		for _, want := range []string{"rowClickAction", "selectedKnowledgeTypeCode", "knowledge_type_", "knowledge_category_"} {
			if !strings.Contains(string(manageJSON), want) {
				t.Fatalf("type manage dialog missing %q", want)
			}
		}
	}
	editDialog, ok := findObjectWithString(root, "open", "${knowledgeTypeEditVisible}")
	if !ok {
		t.Fatalf("knowledge type edit dialog not found")
	}
	for _, label := range []string{"所属类型", "名称", "编码", "排序"} {
		if _, ok := findObjectWithString(editDialog, "label", label); !ok {
			t.Fatalf("knowledge type edit dialog missing label %q", label)
		}
	}
	for _, service := range []string{"system.get_sys_code", "system.sys_code_save", "system.sys_code_update", "system.sys_code_delete"} {
		if _, ok := findObjectWithString(root, "api", "post:/template_data/data?service="+service); !ok {
			t.Fatalf("question-taxonomy should call %s", service)
		}
	}

	knowledgeDialog, ok := findObjectWithString(root, "open", "${knowledgeDialogVisible}")
	if !ok {
		t.Fatalf("knowledge dialog not found")
	}
	for _, label := range []string{"类型", "具体分类", "明细"} {
		if _, ok := findObjectWithString(knowledgeDialog, "label", label); !ok {
			t.Fatalf("knowledge dialog missing label %q", label)
		}
	}
	list, ok := findObjectWithString(root, "title", "目录结构列表")
	if !ok {
		t.Fatalf("目录结构列表 layout not found")
	}
	table, ok := findObjectWithString(list["children"], "tag", "table")
	if !ok {
		t.Fatalf("目录结构列表 table not found")
	}
	for _, field := range []string{"knowledge_type_name", "knowledge_category_name", "content_detail"} {
		if _, ok := findObjectWithString(table["columnDefs"], "field", field); !ok {
			t.Fatalf("catalog table missing %s column", field)
		}
	}

	save := readServiceDefinition(t, "collect/question/knowledge/index.yml", "knowledge_save")
	update := readServiceDefinition(t, "collect/question/knowledge/index.yml", "knowledge_update")
	for _, key := range []string{"knowledge_type", "knowledge_category"} {
		if _, ok := save.Params[key]; !ok {
			t.Fatalf("knowledge_save params missing %s", key)
		}
		if !stringSliceContains(update.UpdateFields, key) {
			t.Fatalf("knowledge_update update_fields missing %s", key)
		}
	}

	sqlData, err := os.ReadFile("collect/question/knowledge/knowledge_query.sql")
	if err != nil {
		t.Fatalf("read knowledge_query.sql: %v", err)
	}
	sql := string(sqlData)
	for _, want := range []string{"a.knowledge_type", "a.knowledge_category", "knowledge_type_name", "knowledge_category_name", "knowledge_type_", "knowledge_category_"} {
		if !strings.Contains(sql, want) {
			t.Fatalf("knowledge_query.sql missing %q", want)
		}
	}

	migrationData, err := os.ReadFile("scripts/migrate-question-web.sql")
	if err != nil {
		t.Fatalf("read migrate-question-web.sql: %v", err)
	}
	migration := string(migrationData)
	for _, want := range []string{"knowledge_type TEXT", "knowledge_category TEXT", "SET sys_code_type = 'knowledge_type_english'", "SET sys_code_type = 'knowledge_category_english'", "'knowledge-type-vocabulary', 'knowledge_type_english'", "'knowledge-category-verb-tense', 'knowledge_category_english'"} {
		if !strings.Contains(migration, want) {
			t.Fatalf("migrate-question-web.sql missing %q", want)
		}
	}
}

func countObjectsWithString(root any, key string, value string) int {
	total := 0
	visitObjects(root, func(obj map[string]any) {
		if stringField(obj, key) == value {
			total++
		}
	})
	return total
}

func numberField(obj map[string]any, key string) float64 {
	value, ok := obj[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	default:
		return 0
	}
}

func taxonomyContainer(t *testing.T, root any, spec batchDeleteSpec) map[string]any {
	t.Helper()
	var scope map[string]any
	var ok bool
	if spec.tabKey != "" {
		scope, ok = findObjectWithString(root, "key", spec.tabKey)
		if !ok {
			t.Fatalf("%s tab not found", spec.name)
		}
		scope, ok = findObjectWithString(scope["children"], "tag", "layout-fit")
		if !ok {
			t.Fatalf("%s layout-fit not found", spec.name)
		}
		return scope
	}

	if spec.panelTitle != "" {
		scope, ok = findObjectWithString(root, "id", "taxonomy-knowledge-panel")
		if ok {
			scope, ok = findObjectWithString(scope["children"], "tag", "layout-fit")
			if ok && stringField(scope, "title") == spec.panelTitle {
				return scope
			}
		}
	}

	scope, ok = findObjectWithString(root, "title", spec.panelTitle)
	if !ok {
		t.Fatalf("%s panel not found", spec.name)
	}
	return scope
}

func assertBatchDeleteButton(t *testing.T, button map[string]any, spec batchDeleteSpec) {
	t.Helper()
	if got := stringField(button, "disabled"); !strings.Contains(got, spec.selection) {
		t.Fatalf("%s disabled expression %q does not reference %s", spec.name, got, spec.selection)
	}
	if got := stringField(button, "icon"); got != "DeleteOutlined" {
		t.Fatalf("%s icon = %q, want DeleteOutlined", spec.name, got)
	}

	ajax, ok := findObjectWithString(button["action"], "tag", "ajax")
	if !ok {
		t.Fatalf("%s ajax action not found", spec.name)
	}
	if got, want := stringField(ajax, "api"), "post:/template_data/data?service="+spec.service; got != want {
		t.Fatalf("%s ajax api = %q, want %q", spec.name, got, want)
	}
	data, ok := ajax["data"].(map[string]any)
	if !ok {
		t.Fatalf("%s ajax data missing", spec.name)
	}
	value, ok := data[spec.listKey].(string)
	if !ok {
		t.Fatalf("%s ajax data missing %s", spec.name, spec.listKey)
	}
	if !strings.Contains(value, spec.selection) || !strings.Contains(value, spec.rowID) {
		t.Fatalf("%s %s expression %q should reference %s and %s", spec.name, spec.listKey, value, spec.selection, spec.rowID)
	}

	updateStore, ok := findObjectWithString(button["action"], "tag", "update-store")
	if !ok {
		t.Fatalf("%s update-store action not found", spec.name)
	}
	storeValue, ok := updateStore["value"].(map[string]any)
	if !ok {
		t.Fatalf("%s update-store value missing", spec.name)
	}
	selectionValue, ok := storeValue[spec.selection].([]any)
	if !ok || len(selectionValue) != 0 {
		t.Fatalf("%s update-store should clear %s", spec.name, spec.selection)
	}
}

func readJSONFile(t *testing.T, name string) any {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	return root
}

type serviceFile struct {
	Service []serviceDefinition `yaml:"service"`
}

type serviceDefinition struct {
	Key           string            `yaml:"key"`
	Module        string            `yaml:"module"`
	Table         string            `yaml:"table"`
	ModelField    any               `yaml:"model_field"`
	Params        map[string]any    `yaml:"params"`
	UpdateFields  []string          `yaml:"update_fields"`
	Filter        map[string]string `yaml:"filter"`
	HandlerParams []map[string]any  `yaml:"handler_params"`
}

func readServiceDefinition(t *testing.T, name string, key string) serviceDefinition {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	var cfg serviceFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	for _, service := range cfg.Service {
		if service.Key == key {
			if service.Params == nil {
				service.Params = map[string]any{}
			}
			if service.Filter == nil {
				service.Filter = map[string]string{}
			}
			return service
		}
	}
	t.Fatalf("service %s not found in %s", key, name)
	return serviceDefinition{}
}

func collectButtonsWithChildren(root any, label string) []map[string]any {
	var out []map[string]any
	visitObjects(root, func(obj map[string]any) {
		if stringField(obj, "tag") == "button" && stringField(obj, "children") == label {
			out = append(out, obj)
		}
	})
	return out
}

func countButtonsWithChildren(root any, label string) int {
	return len(collectButtonsWithChildren(root, label))
}

func findObjectWithString(root any, field string, value string) (map[string]any, bool) {
	var found map[string]any
	visitObjects(root, func(obj map[string]any) {
		if found != nil {
			return
		}
		if stringField(obj, field) == value {
			found = obj
		}
	})
	return found, found != nil
}

func visitObjects(root any, visit func(map[string]any)) {
	switch node := root.(type) {
	case map[string]any:
		visit(node)
		for _, child := range node {
			visitObjects(child, visit)
		}
	case []any:
		for _, child := range node {
			visitObjects(child, visit)
		}
	}
}

func stringField(obj map[string]any, field string) string {
	value, _ := obj[field].(string)
	return value
}

func stringSliceContains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
