package plugins

import "testing"

func TestBuildPinyinCode(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "Chinese name", text: "词汇拼写测试", want: "cihuipinxieceshi"},
		{name: "Mixed content", text: "Unit 1 词汇-测试", want: "unit_1_cihui_ceshi"},
		{name: "Existing ascii code", text: "Grade 8", want: "grade_8"},
		{name: "Non mapped Chinese", text: "圆锥体积", want: "yuanzhuitiji"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPinyinCode(tt.text); got != tt.want {
				t.Fatalf("buildPinyinCode(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}
