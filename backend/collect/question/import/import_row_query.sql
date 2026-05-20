SELECT
  a.*
FROM question_import_row a
WHERE a.batch_id = {{.batch_id}}
{{ if .validate_status }}
AND a.validate_status = {{.validate_status}}
{{ end }}
ORDER BY a.row_index
