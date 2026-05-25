SELECT
  a.question_id,
  a.question_code,
  a.title,
  a.subject,
  COALESCE(subject.subject_name, subject_code.sys_code_text) AS subject_name,
  a.stage,
  stage_code.sys_code_text AS stage_name,
  a.grade,
  COALESCE(grade.grade_name, grade_code.sys_code_text) AS grade_name,
  semester_code.sys_code_text AS grade_semester_name,
  CASE
    WHEN ifnull(grade.grade_name, '') <> '' AND ifnull(semester_code.sys_code_text, '') <> '' THEN CONCAT(grade.grade_name, semester_code.sys_code_text)
    ELSE COALESCE(grade.grade_name, grade_code.sys_code_text)
  END AS grade_label,
  a.textbook_version,
  textbook_code.sys_code_text AS textbook_version_name,
  a.unit_id,
  a.unit_code,
  a.unit_name,
  a.question_type,
  type_code.sys_code_text AS question_type_name,
  a.question_category,
  category_code.sys_code_text AS question_category_name,
  a.difficulty,
  difficulty_code.sys_code_text AS difficulty_name,
  a.score,
  a.duration_seconds,
  a.sequence_no,
  a.stem_html,
  a.stem_text,
  a.analysis_media_url,
  a.analysis_media_name,
  a.analysis_media_type,
  a.option_count,
  a.blank_count,
  a.asset_count,
  a.source,
  a.status,
  status_code.sys_code_text AS status_name,
  a.version,
  a.publish_time,
  a.publish_user,
  a.remark,
  a.create_time,
  a.create_user,
  a.modify_time,
  a.modify_user
FROM question_item a
LEFT JOIN question_subject subject ON subject.subject_code = a.subject AND ifnull(subject.is_delete, '0') = '0'
LEFT JOIN question_grade grade ON grade.grade_code = a.grade AND ifnull(grade.is_delete, '0') = '0'
LEFT JOIN sys_code subject_code ON subject_code.sys_code_type = 'subject' AND subject_code.sys_code = a.subject
LEFT JOIN sys_code stage_code ON stage_code.sys_code_type = 'study_stage' AND stage_code.sys_code = a.stage
LEFT JOIN sys_code grade_code ON grade_code.sys_code_type = 'grade' AND grade_code.sys_code = a.grade
LEFT JOIN sys_code semester_code ON semester_code.sys_code_type = 'grade_semester' AND semester_code.sys_code = COALESCE(NULLIF(grade.semester, ''), 'upper')
LEFT JOIN sys_code textbook_code ON textbook_code.sys_code_type = 'textbook_version' AND textbook_code.sys_code = a.textbook_version
LEFT JOIN sys_code type_code ON type_code.sys_code_type = 'question_type' AND type_code.sys_code = a.question_type
LEFT JOIN sys_code category_code ON category_code.sys_code_type = 'question_category' AND category_code.sys_code = a.question_category
LEFT JOIN sys_code difficulty_code ON difficulty_code.sys_code_type = 'question_difficulty' AND difficulty_code.sys_code = a.difficulty
LEFT JOIN sys_code status_code ON status_code.sys_code_type = 'question_status' AND status_code.sys_code = a.status
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
ORDER BY a.sequence_no ASC, a.modify_time DESC, a.create_time DESC
{{ if .pagination }}
LIMIT {{.start}}, {{.size}}
{{ end }}
