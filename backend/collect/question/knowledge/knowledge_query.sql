SELECT
  a.knowledge_id,
  COALESCE(NULLIF(a.subject, ''), unit.subject) AS subject,
  COALESCE(NULLIF(a.stage, ''), unit.stage) AS stage,
  COALESCE(NULLIF(a.grade, ''), unit.grade) AS grade,
  a.parent_id,
  a.knowledge_code,
  a.knowledge_name,
  ifnull(a.knowledge_type, '') AS knowledge_type,
  ifnull(a.knowledge_category, '') AS knowledge_category,
  ifnull(type_code.sys_code_text, '') AS knowledge_type_name,
  ifnull(category_code.sys_code_text, '') AS knowledge_category_name,
  a.order_index,
  a.status,
  a.is_delete,
  a.create_time,
  a.create_user,
  a.modify_time,
  a.modify_user,
  a.parent_id AS unit_id,
  unit.unit_code,
  unit.unit_name,
  subject.subject_name,
  grade.grade_name,
  grade.semester AS grade_semester,
  semester_code.sys_code_text AS grade_semester_name,
  CASE
    WHEN ifnull(grade.grade_name, '') <> '' AND ifnull(semester_code.sys_code_text, '') <> '' THEN CONCAT(grade.grade_name, semester_code.sys_code_text)
    ELSE COALESCE(grade.grade_name, COALESCE(NULLIF(a.grade, ''), unit.grade))
  END AS grade_label,
  stage_code.sys_code_text AS stage_name,
  ifnull(content_summary.content_count, 0) AS content_count,
  ifnull(content_summary.content_detail, '') AS content_detail,
  ifnull(source_summary.pdf_source_image_url, '') AS pdf_source_image_url,
  ifnull(source_summary.pdf_source_excerpt, '') AS pdf_source_excerpt,
  ifnull(source_summary.source_doc_id, '') AS source_doc_id,
  ifnull(source_summary.source_file_name, '') AS source_file_name
FROM question_knowledge a
LEFT JOIN question_unit unit ON unit.unit_id = a.parent_id AND ifnull(unit.is_delete, '0') = '0'
LEFT JOIN question_subject subject ON subject.subject_code = COALESCE(NULLIF(a.subject, ''), unit.subject) AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = COALESCE(NULLIF(a.grade, ''), unit.grade) AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(grade.semester, ''), 'upper')
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = COALESCE(NULLIF(a.stage, ''), unit.stage)
LEFT JOIN sys_code type_code ON type_code.sys_code_type = CONCAT('knowledge_type_', COALESCE(NULLIF(a.subject, ''), unit.subject, 'english')) AND type_code.sys_code = ifnull(a.knowledge_type, '')
LEFT JOIN sys_code category_code ON category_code.sys_code_type = CONCAT('knowledge_category_', COALESCE(NULLIF(a.subject, ''), unit.subject, 'english')) AND category_code.sys_code = ifnull(a.knowledge_category, '') AND ifnull(category_code.icon, '') = ifnull(a.knowledge_type, '')
LEFT JOIN (
  SELECT
    c.knowledge_id,
    COUNT(*) AS content_count,
    GROUP_CONCAT(
      CONCAT(
        ifnull(NULLIF(c.section_title, ''), '知识正文'),
        '：',
        LEFT(ifnull(c.content_text, ''), 500)
      )
      ORDER BY c.order_index
      SEPARATOR '\n\n'
    ) AS content_detail
  FROM question_knowledge_content c
  WHERE ifnull(c.status, '') <> 'deleted'
  GROUP BY c.knowledge_id
) content_summary ON content_summary.knowledge_id = a.knowledge_id
LEFT JOIN (
  SELECT
    c.knowledge_id,
    MIN(c.source_doc_id) AS source_doc_id,
    MIN(d.file_name) AS source_file_name,
    MIN(NULLIF(b.block_image_url, '')) AS pdf_source_image_url,
    GROUP_CONCAT(
      DISTINCT CONCAT(
        'P',
        ifnull(f.page_no, 1),
        ' ',
        LEFT(ifnull(NULLIF(r.raw_quote, ''), ifnull(f.raw_text, '')), 500)
      )
      SEPARATOR '\n---\n'
    ) AS pdf_source_excerpt
  FROM question_knowledge_content c
  LEFT JOIN question_source_field_rel r ON r.entity_type = 'knowledge_content' AND r.entity_id = c.content_id
  LEFT JOIN question_source_fragment f ON f.source_fragment_id = r.source_fragment_id
  LEFT JOIN question_source_block b ON b.source_block_id = r.source_block_id
  LEFT JOIN question_source_document d ON d.source_doc_id = c.source_doc_id
  WHERE ifnull(c.status, '') <> 'deleted'
  GROUP BY c.knowledge_id
) source_summary ON source_summary.knowledge_id = a.knowledge_id
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .subject }}
AND COALESCE(NULLIF(a.subject, ''), unit.subject) = {{.subject}}
{{ end }}
{{ if .stage }}
AND COALESCE(NULLIF(a.stage, ''), unit.stage) = {{.stage}}
{{ end }}
{{ if .grade }}
AND COALESCE(NULLIF(a.grade, ''), unit.grade) = {{.grade}}
{{ end }}
{{ if .parent_id }}
AND a.parent_id = {{.parent_id}}
{{ end }}
{{ if .unit_id }}
AND a.parent_id = {{.unit_id}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (
  a.knowledge_code LIKE {{.keyword}}
  OR a.knowledge_name LIKE {{.keyword}}
  OR unit.unit_code LIKE {{.keyword}}
  OR unit.unit_name LIKE {{.keyword}}
  OR subject.subject_name LIKE {{.keyword}}
  OR grade.grade_name LIKE {{.keyword}}
  OR semester_code.sys_code_text LIKE {{.keyword}}
  OR type_code.sys_code_text LIKE {{.keyword}}
  OR category_code.sys_code_text LIKE {{.keyword}}
  OR content_summary.content_detail LIKE {{.keyword}}
  OR source_summary.pdf_source_excerpt LIKE {{.keyword}}
)
{{ end }}
ORDER BY a.order_index, a.knowledge_code
