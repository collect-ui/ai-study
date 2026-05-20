SELECT
  COUNT(1) AS total_count,
  COALESCE(SUM(CASE WHEN date(a.create_time) = curdate() THEN 1 ELSE 0 END), 0) AS today_new_count,
  COALESCE(SUM(CASE WHEN a.status = 'reviewing' THEN 1 ELSE 0 END), 0) AS reviewing_count,
  COALESCE(SUM(CASE WHEN a.status = 'published' THEN 1 ELSE 0 END), 0) AS published_count
FROM question_item a
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .subject }}
AND a.subject = {{.subject}}
{{ end }}
{{ if .grade }}
AND a.grade = {{.grade}}
{{ end }}
{{ if .stage }}
AND a.stage = {{.stage}}
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
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
