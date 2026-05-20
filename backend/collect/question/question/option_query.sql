SELECT
  a.*
FROM question_option a
WHERE ifnull(a.is_delete, '0') = '0'
AND a.question_id = {{.question_id}}
ORDER BY a.option_order, a.option_key
