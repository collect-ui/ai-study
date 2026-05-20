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
			name:        "单元",
			tabKey:      "unit",
			selection:   "unitSelection",
			rowID:       "unit_id",
			listKey:     "unit_id_list",
			service:     "question.unit_delete",
			serviceFile: "collect/question/unit/index.yml",
			filterKey:   "unit_id__in",
		},
		{
			name:        "知识点",
			panelTitle:  "知识点维护",
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

func TestQuestionTaxonomyDeleteServicesAcceptIDLists(t *testing.T) {
	specs := []batchDeleteSpec{
		{name: "年级", service: "question.grade_delete", serviceFile: "collect/question/grade/index.yml", listKey: "grade_id_list", filterKey: "grade_id__in"},
		{name: "科目", service: "question.subject_delete", serviceFile: "collect/question/subject/index.yml", listKey: "subject_id_list", filterKey: "subject_id__in"},
		{name: "单元", service: "question.unit_delete", serviceFile: "collect/question/unit/index.yml", listKey: "unit_id_list", filterKey: "unit_id__in"},
		{name: "知识点", service: "question.knowledge_delete", serviceFile: "collect/question/knowledge/index.yml", listKey: "knowledge_id_list", filterKey: "knowledge_id__in"},
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
