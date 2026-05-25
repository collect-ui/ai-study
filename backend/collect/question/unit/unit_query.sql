SELECT
  a.*,
  subject.subject_name,
  grade.grade_name,
  grade.semester AS grade_semester,
  semester_code.sys_code_text AS grade_semester_name,
  CASE
    WHEN ifnull(grade.grade_name, '') <> '' AND ifnull(semester_code.sys_code_text, '') <> '' THEN CONCAT(grade.grade_name, semester_code.sys_code_text)
    ELSE COALESCE(grade.grade_name, a.grade)
  END AS grade_label,
  grade.stage AS grade_stage,
  stage_code.sys_code_text AS stage_name,
  textbook_code.sys_code_text AS textbook_version_name
FROM question_unit a
LEFT JOIN question_subject subject ON subject.subject_code = a.subject AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = a.grade AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(grade.semester, ''), 'upper')
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = a.stage
LEFT JOIN sys_code textbook_code ON textbook_code.sys_code_type = 'textbook_version' AND textbook_code.sys_code = a.textbook_version
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
{{ if .textbook_version }}
AND a.textbook_version = {{.textbook_version}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (
  a.unit_code LIKE {{.keyword}}
  OR a.unit_name LIKE {{.keyword}}
  OR subject.subject_name LIKE {{.keyword}}
  OR grade.grade_name LIKE {{.keyword}}
  OR semester_code.sys_code_text LIKE {{.keyword}}
)
{{ end }}
ORDER BY COALESCE(grade.order_index, 9999), a.order_index, a.unit_code
