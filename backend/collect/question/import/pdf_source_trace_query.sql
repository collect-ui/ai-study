SELECT
  r.field_rel_id,
  r.source_doc_id,
  d.file_name,
  r.source_page_id,
  r.source_block_id,
  r.source_fragment_id,
  f.page_no,
  f.fragment_order,
  f.fragment_type,
  r.entity_type,
  r.entity_id,
  r.field_name,
  r.field_part_order,
  r.extracted_value,
  r.normalized_value,
  r.raw_quote,
  r.context_before,
  r.context_after,
  r.char_start,
  r.char_end,
  r.confidence,
  r.match_status,
  r.review_status,
  r.create_time
FROM question_source_field_rel r
LEFT JOIN question_source_document d ON d.source_doc_id = r.source_doc_id
LEFT JOIN question_source_fragment f ON f.source_fragment_id = r.source_fragment_id
WHERE 1 = 1
{{ if .source_doc_id }}
AND r.source_doc_id = {{.source_doc_id}}
{{ end }}
{{ if .file_name }}
AND d.file_name LIKE {{.file_name}}
{{ end }}
{{ if .entity_type }}
AND r.entity_type = {{.entity_type}}
{{ end }}
{{ if .entity_id }}
AND r.entity_id = {{.entity_id}}
{{ end }}
{{ if .field_name }}
AND r.field_name = {{.field_name}}
{{ end }}
{{ if .match_status }}
AND r.match_status = {{.match_status}}
{{ end }}
ORDER BY d.create_time DESC, f.page_no, f.fragment_order, r.field_part_order
{{ if .size }}
LIMIT {{.size}} OFFSET {{.start}}
{{ end }}
