SELECT
  a.*,
  stage_code.sys_code_text AS stage_name
FROM question_grade a
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = a.stage
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .stage }}
AND a.stage = {{.stage}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (a.grade_code LIKE {{.keyword}} OR a.grade_name LIKE {{.keyword}})
{{ end }}
ORDER BY a.order_index, a.grade_code
