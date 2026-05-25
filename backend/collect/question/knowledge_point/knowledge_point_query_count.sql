SELECT COUNT(1)
FROM question_knowledge_point a
LEFT JOIN question_unit unit ON unit.unit_id = a.unit_id AND ifnull(unit.is_delete, '0') = '0'
LEFT JOIN question_section section ON section.section_id = ifnull(a.section_id, '') AND ifnull(section.is_delete, '0') = '0'
LEFT JOIN question_subject subject ON subject.subject_code = a.subject AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = a.grade AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(grade.semester, ''), 'upper')
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .subject }}
AND a.subject = {{.subject}}
{{ end }}
{{ if .stage }}
AND a.stage = {{.stage}}
{{ end }}
{{ if .grade }}
AND a.grade = {{.grade}}
{{ end }}
{{ if .unit_id }}
AND a.unit_id = {{.unit_id}}
{{ end }}
{{ if .section_id }}
AND ifnull(a.section_id, '') = {{.section_id}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (
  a.point_code LIKE {{.keyword}}
  OR a.point_name LIKE {{.keyword}}
  OR a.content_detail LIKE {{.keyword}}
  OR unit.unit_code LIKE {{.keyword}}
  OR unit.unit_name LIKE {{.keyword}}
  OR section.section_code LIKE {{.keyword}}
  OR section.section_name LIKE {{.keyword}}
  OR subject.subject_name LIKE {{.keyword}}
  OR grade.grade_name LIKE {{.keyword}}
  OR semester_code.sys_code_text LIKE {{.keyword}}
)
{{ end }}
