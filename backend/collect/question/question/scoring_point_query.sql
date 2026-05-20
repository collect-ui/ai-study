SELECT
  a.*
FROM question_scoring_point a
WHERE ifnull(a.is_delete, '0') = '0'
AND a.question_id = {{.question_id}}
ORDER BY a.point_index
