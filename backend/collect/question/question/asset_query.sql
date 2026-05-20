SELECT
  a.*
FROM question_asset a
WHERE ifnull(a.is_delete, '0') = '0'
AND a.question_id = {{.question_id}}
ORDER BY a.usage_type, a.create_time
