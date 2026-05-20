SELECT
  a.*
FROM question_subject a
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (a.subject_code LIKE {{.keyword}} OR a.subject_name LIKE {{.keyword}})
{{ end }}
ORDER BY a.order_index, a.subject_code
