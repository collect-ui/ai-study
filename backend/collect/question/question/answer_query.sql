SELECT
  a.*
FROM question_answer a
WHERE ifnull(a.is_delete, '0') = '0'
AND a.question_id = {{.question_id}}
ORDER BY a.create_time DESC
