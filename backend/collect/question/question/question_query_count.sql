SELECT COUNT(1)
FROM question_item a
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .keyword }}
AND (
  a.question_code LIKE {{.keyword}}
  OR a.title LIKE {{.keyword}}
  OR a.stem_html LIKE {{.keyword}}
  OR a.stem_text LIKE {{.keyword}}
  OR a.analysis_text LIKE {{.keyword}}
  OR EXISTS (
    SELECT 1
    FROM question_knowledge_rel kr
    LEFT JOIN question_knowledge knowledge ON knowledge.knowledge_id = kr.knowledge_id AND ifnull(knowledge.is_delete, '0') = '0'
    WHERE kr.question_id = a.question_id
      AND COALESCE(NULLIF(kr.knowledge_name, ''), knowledge.knowledge_name) LIKE {{.keyword}}
  )
)
{{ end }}
{{ if .subject }}
AND a.subject = {{.subject}}
{{ end }}
{{ if .stage }}
AND a.stage = {{.stage}}
{{ end }}
{{ if .grade }}
AND a.grade = {{.grade}}
{{ end }}
{{ if .textbook_version }}
AND a.textbook_version = {{.textbook_version}}
{{ end }}
{{ if .unit_id }}
AND a.unit_id = {{.unit_id}}
{{ end }}
{{ if .knowledge_id }}
AND EXISTS (
  SELECT 1
  FROM question_knowledge_rel kr
  WHERE kr.question_id = a.question_id
    AND kr.knowledge_id = {{.knowledge_id}}
)
{{ end }}
{{ if .question_type }}
AND a.question_type = {{.question_type}}
{{ end }}
{{ if .question_category }}
AND a.question_category = {{.question_category}}
{{ end }}
{{ if .difficulty }}
AND a.difficulty = {{.difficulty}}
{{ end }}
{{ if .difficulty_list }}
AND a.difficulty IN ({{.difficulty_list}})
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
