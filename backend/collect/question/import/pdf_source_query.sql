SELECT
  d.source_doc_id,
  d.import_batch_id,
  d.file_name,
  d.file_sha256,
  d.file_url,
  d.page_count,
  d.subject,
  d.stage,
  d.grade,
  d.textbook_version,
  d.parse_status,
  d.import_status,
  d.create_time,
  f.source_fragment_id,
  f.source_page_id,
  f.source_block_id,
  f.page_no,
  f.fragment_order,
  f.fragment_type,
  f.raw_text,
  f.raw_html,
  f.context_before,
  f.context_after,
  f.fragment_hash,
  f.confidence,
  f.status AS fragment_status
FROM question_source_document d
LEFT JOIN question_source_fragment f ON f.source_doc_id = d.source_doc_id
WHERE 1 = 1
{{ if .source_doc_id }}
AND d.source_doc_id = {{.source_doc_id}}
{{ end }}
{{ if .import_batch_id }}
AND d.import_batch_id = {{.import_batch_id}}
{{ end }}
{{ if .file_name }}
AND d.file_name LIKE {{.file_name}}
{{ end }}
{{ if .subject }}
AND d.subject = {{.subject}}
{{ end }}
{{ if .grade }}
AND d.grade = {{.grade}}
{{ end }}
{{ if .textbook_version }}
AND d.textbook_version = {{.textbook_version}}
{{ end }}
{{ if .fragment_type }}
AND f.fragment_type = {{.fragment_type}}
{{ end }}
ORDER BY d.create_time DESC, f.page_no, f.fragment_order
{{ if .size }}
LIMIT {{.size}} OFFSET {{.start}}
{{ end }}
