SELECT
  s.snapshot_id,
  s.source_doc_id,
  d.file_name,
  s.source_page_id,
  s.source_block_id,
  s.source_fragment_id,
  s.question_id,
  s.knowledge_id,
  s.content_id,
  s.extract_service,
  s.extract_params_json,
  s.raw_text,
  s.raw_html,
  s.normalized_text,
  s.ai_output_json,
  s.validator_result_json,
  s.status,
  s.create_time
FROM question_source_snapshot s
LEFT JOIN question_source_document d ON d.source_doc_id = s.source_doc_id
WHERE 1 = 1
{{ if .source_doc_id }}
AND s.source_doc_id = {{.source_doc_id}}
{{ end }}
{{ if .file_name }}
AND d.file_name LIKE {{.file_name}}
{{ end }}
{{ if .knowledge_id }}
AND s.knowledge_id = {{.knowledge_id}}
{{ end }}
{{ if .question_id }}
AND s.question_id = {{.question_id}}
{{ end }}
{{ if .content_id }}
AND s.content_id = {{.content_id}}
{{ end }}
ORDER BY s.create_time DESC
{{ if .size }}
LIMIT {{.size}} OFFSET {{.start}}
{{ end }}
