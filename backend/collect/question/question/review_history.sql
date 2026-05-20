SELECT
  a.*,
  u.nick AS review_user_name
FROM question_review_record a
LEFT JOIN user_account u ON u.user_id = a.review_user
WHERE ifnull(a.is_delete, '0') = '0'
AND a.question_id = {{.question_id}}
ORDER BY a.review_time DESC
