SELECT
  a.section_id,
  a.subject,
  a.stage,
  a.grade,
  a.unit_id,
  a.section_code,
  a.section_name,
  a.order_index,
  a.status
FROM question_section a
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
AND a.section_name = {{.section_name}}
ORDER BY a.order_index, a.section_code
LIMIT 1
