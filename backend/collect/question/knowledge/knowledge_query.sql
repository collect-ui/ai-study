SELECT
  a.knowledge_id,
  COALESCE(NULLIF(a.subject, ''), unit.subject) AS subject,
  COALESCE(NULLIF(a.stage, ''), unit.stage) AS stage,
  COALESCE(NULLIF(a.grade, ''), unit.grade) AS grade,
  a.parent_id,
  a.knowledge_code,
  a.knowledge_name,
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
  stage_code.sys_code_text AS stage_name
FROM question_knowledge a
LEFT JOIN question_unit unit ON unit.unit_id = a.parent_id AND ifnull(unit.is_delete, '0') = '0'
LEFT JOIN question_subject subject ON subject.subject_code = COALESCE(NULLIF(a.subject, ''), unit.subject) AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = COALESCE(NULLIF(a.grade, ''), unit.grade) AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = COALESCE(NULLIF(a.stage, ''), unit.stage)
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
)
{{ end }}
ORDER BY a.order_index, a.knowledge_code
