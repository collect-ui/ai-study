SELECT
  a.*
FROM question_import_batch a
WHERE 1 = 1
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
ORDER BY a.create_time DESC
