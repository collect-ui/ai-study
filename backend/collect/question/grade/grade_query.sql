SELECT
  a.grade_id,
  a.stage,
  COALESCE(NULLIF(a.semester, ''), 'upper') AS semester,
  a.grade_code,
  a.grade_name,
  a.order_index,
  a.status,
  a.is_delete,
  a.create_time,
  a.create_user,
  a.modify_time,
  a.modify_user,
  stage_code.sys_code_text AS stage_name,
  semester_code.sys_code_text AS semester_name,
  CASE
    WHEN ifnull(semester_code.sys_code_text, '') <> '' THEN CONCAT(a.grade_name, semester_code.sys_code_text)
    ELSE a.grade_name
  END AS grade_label
FROM question_grade a
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = a.stage
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(a.semester, ''), 'upper')
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .stage }}
AND a.stage = {{.stage}}
{{ end }}
{{ if .semester }}
AND COALESCE(NULLIF(a.semester, ''), 'upper') = {{.semester}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (a.grade_code LIKE {{.keyword}} OR a.grade_name LIKE {{.keyword}})
{{ end }}
ORDER BY a.order_index, a.grade_code
