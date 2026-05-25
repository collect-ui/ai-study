SELECT
  a.point_id,
  a.subject,
  a.stage,
  a.grade,
  a.unit_id,
  ifnull(a.section_id, '') AS section_id,
  a.point_code,
  a.point_name,
  ifnull(a.content_detail, '') AS content_detail,
  a.order_index,
  a.status,
  a.is_delete,
  a.create_time,
  a.create_user,
  a.modify_time,
  a.modify_user,
  unit.unit_code,
  unit.unit_name,
  unit.textbook_version,
  textbook_code.sys_code_text AS textbook_version_name,
  section.section_code,
  section.section_name,
  subject.subject_name,
  grade.grade_name,
  grade.semester AS grade_semester,
  semester_code.sys_code_text AS grade_semester_name,
  CASE
    WHEN ifnull(grade.grade_name, '') <> '' AND ifnull(semester_code.sys_code_text, '') <> '' THEN CONCAT(grade.grade_name, semester_code.sys_code_text)
    ELSE COALESCE(grade.grade_name, a.grade)
  END AS grade_label,
  stage_code.sys_code_text AS stage_name
FROM question_knowledge_point a
LEFT JOIN question_unit unit ON unit.unit_id = a.unit_id AND ifnull(unit.is_delete, '0') = '0'
LEFT JOIN question_section section ON section.section_id = ifnull(a.section_id, '') AND ifnull(section.is_delete, '0') = '0'
LEFT JOIN question_subject subject ON subject.subject_code = a.subject AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = a.grade AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(grade.semester, ''), 'upper')
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = a.stage
LEFT JOIN sys_code textbook_code ON textbook_code.sys_code_type = 'textbook_version' AND textbook_code.sys_code = unit.textbook_version
WHERE ifnull(a.is_delete, '0') = '0'
{{ if .subject }}
AND a.subject = {{.subject}}
{{ end }}
{{ if .stage }}
AND a.stage = {{.stage}}
{{ end }}
{{ if .grade }}
AND a.grade = {{.grade}}
{{ end }}
{{ if .unit_id }}
AND a.unit_id = {{.unit_id}}
{{ end }}
{{ if .section_id }}
AND ifnull(a.section_id, '') = {{.section_id}}
{{ end }}
{{ if .status }}
AND a.status = {{.status}}
{{ end }}
{{ if .keyword }}
AND (
  a.point_code LIKE {{.keyword}}
  OR a.point_name LIKE {{.keyword}}
  OR a.content_detail LIKE {{.keyword}}
  OR unit.unit_code LIKE {{.keyword}}
  OR unit.unit_name LIKE {{.keyword}}
  OR section.section_code LIKE {{.keyword}}
  OR section.section_name LIKE {{.keyword}}
  OR subject.subject_name LIKE {{.keyword}}
  OR grade.grade_name LIKE {{.keyword}}
  OR semester_code.sys_code_text LIKE {{.keyword}}
)
{{ end }}
ORDER BY COALESCE(grade.order_index, 9999), unit.order_index, section.order_index, a.order_index, a.point_code
{{ if .pagination }}
LIMIT {{.start}}, {{.size}}
{{ end }}
