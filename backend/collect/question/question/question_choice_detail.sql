SELECT
  a.question_id,
  a.question_code,
  a.title,
  a.subject,
  a.stage,
  a.grade,
  a.textbook_version,
  a.unit_id,
  a.unit_code,
  a.unit_name,
  COALESCE((
    SELECT JSON_ARRAYAGG(rel.knowledge_id)
    FROM (
      SELECT kr.knowledge_id
      FROM question_knowledge_rel kr
      WHERE kr.question_id = a.question_id
      ORDER BY kr.order_index
    ) rel
  ), '[]') AS knowledge_id,
  COALESCE((
    SELECT GROUP_CONCAT(rel.knowledge_name SEPARATOR '、')
    FROM (
      SELECT COALESCE(NULLIF(kr.knowledge_name, ''), knowledge.knowledge_name) AS knowledge_name
      FROM question_knowledge_rel kr
      LEFT JOIN question_knowledge knowledge ON knowledge.knowledge_id = kr.knowledge_id AND ifnull(knowledge.is_delete, '0') = '0'
      WHERE kr.question_id = a.question_id
      ORDER BY kr.order_index
    ) rel
  ), '') AS knowledge_name,
  a.question_type,
  a.question_category,
  a.difficulty,
  a.score,
  a.duration_seconds,
  a.sequence_no,
  a.stem_html,
  a.stem_text,
  a.analysis_html,
  a.analysis_text,
  a.analysis_media_url,
  a.analysis_media_name,
  a.analysis_media_type,
  a.option_count,
  a.blank_count,
  a.source,
  a.remark,
  a.status,
  a.version,
  (
    SELECT o.option_html
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'A'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_a_html,
  (
    SELECT o.option_text
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'A'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_a_text,
  (
    SELECT o.option_html
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'B'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_b_html,
  (
    SELECT o.option_text
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'B'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_b_text,
  (
    SELECT o.option_html
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'C'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_c_html,
  (
    SELECT o.option_text
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'C'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_c_text,
  (
    SELECT o.option_html
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'D'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_d_html,
  (
    SELECT o.option_text
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'D'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ) AS option_d_text,
  (
    SELECT ans.answer_text
    FROM question_answer ans
    WHERE ans.question_id = a.question_id
      AND ifnull(ans.is_delete, '0') = '0'
    ORDER BY ans.modify_time DESC, ans.create_time DESC
    LIMIT 1
  ) AS answer_key,
  (
    SELECT ans.answer_text
    FROM question_answer ans
    WHERE ans.question_id = a.question_id
      AND ifnull(ans.is_delete, '0') = '0'
    ORDER BY ans.modify_time DESC, ans.create_time DESC
    LIMIT 1
  ) AS answer_text,
  COALESCE((
    SELECT ans.reference_text
    FROM question_answer ans
    WHERE ans.question_id = a.question_id
      AND ifnull(ans.is_delete, '0') = '0'
    ORDER BY ans.modify_time DESC, ans.create_time DESC
    LIMIT 1
  ), '') AS choice_items,
  COALESCE((
    SELECT JSON_ARRAYAGG(JSON_OBJECT(
      '__rowId', ba.blank_answer_id,
      'blank_answer_id', ba.blank_answer_id,
      'blank_index', ba.blank_index,
      'standard_answer', ba.standard_answer,
      'alternative_answers', ba.alternative_answers,
      'score', ba.score,
      'match_mode', ba.match_mode,
      'case_sensitive', ba.case_sensitive
    ))
    FROM (
      SELECT *
      FROM question_blank_answer
      WHERE question_id = a.question_id
        AND ifnull(is_delete, '0') = '0'
      ORDER BY blank_index
    ) ba
  ), '[]') AS blank_answers,
  COALESCE((
    SELECT o.content_mode
    FROM question_option o
    WHERE o.question_id = a.question_id
      AND o.option_key = 'A'
      AND ifnull(o.is_delete, '0') = '0'
    ORDER BY o.modify_time DESC, o.create_time DESC
    LIMIT 1
  ), 'plain') AS content_mode
FROM question_item a
WHERE a.question_id = {{.question_id}}
  AND ifnull(a.is_delete, '0') = '0'
