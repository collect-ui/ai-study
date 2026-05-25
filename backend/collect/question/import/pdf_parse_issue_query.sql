SELECT
  i.issue_id,
  i.source_doc_id,
  d.file_name,
  i.page_no,
  i.source_block_id,
  i.issue_type,
  i.severity,
  i.raw_text,
  i.crop_image_url,
  i.ai_output_json,
  i.expected_schema,
  i.error_msg,
  i.suggestion,
  i.status,
  i.create_time,
  i.modify_time
FROM question_pdf_parse_issue i
LEFT JOIN question_source_document d ON d.source_doc_id = i.source_doc_id
WHERE 1 = 1
{{ if .source_doc_id }}
AND i.source_doc_id = {{.source_doc_id}}
{{ end }}
{{ if .file_name }}
AND d.file_name LIKE {{.file_name}}
{{ end }}
{{ if .issue_type }}
AND i.issue_type = {{.issue_type}}
{{ end }}
{{ if .severity }}
AND i.severity = {{.severity}}
{{ end }}
{{ if .status }}
AND i.status = {{.status}}
{{ end }}
ORDER BY i.create_time DESC, i.page_no
{{ if .size }}
LIMIT {{.size}} OFFSET {{.start}}
{{ end }}
