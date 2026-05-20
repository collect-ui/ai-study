SELECT
  a.*
FROM question_knowledge_rel a
WHERE a.question_id = {{.question_id}}
ORDER BY a.order_index
