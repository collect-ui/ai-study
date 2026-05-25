PRAGMA foreign_keys = OFF;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS question_item (
  question_id TEXT PRIMARY KEY,
  question_code TEXT NOT NULL,
  title TEXT DEFAULT '',
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT NOT NULL DEFAULT '',
  grade TEXT NOT NULL DEFAULT '',
  textbook_version TEXT NOT NULL DEFAULT '',
  unit_id TEXT NOT NULL DEFAULT '',
  unit_code TEXT DEFAULT '',
  unit_name TEXT DEFAULT '',
  question_type TEXT NOT NULL DEFAULT 'single_choice',
  question_category TEXT NOT NULL DEFAULT 'normal',
  difficulty TEXT NOT NULL DEFAULT 'basic',
  score INTEGER NOT NULL DEFAULT 5,
  duration_seconds INTEGER NOT NULL DEFAULT 0,
  sequence_no INTEGER NOT NULL DEFAULT 0,
  stem_html TEXT NOT NULL DEFAULT '',
  stem_text TEXT NOT NULL DEFAULT '',
  analysis_html TEXT DEFAULT '',
  analysis_text TEXT DEFAULT '',
  analysis_media_url TEXT DEFAULT '',
  analysis_media_name TEXT DEFAULT '',
  analysis_media_type TEXT DEFAULT '',
  option_count INTEGER NOT NULL DEFAULT 0,
  blank_count INTEGER NOT NULL DEFAULT 0,
  asset_count INTEGER NOT NULL DEFAULT 0,
  content_hash TEXT DEFAULT '',
  source TEXT NOT NULL DEFAULT 'manual',
  status TEXT NOT NULL DEFAULT 'draft',
  version INTEGER NOT NULL DEFAULT 1,
  publish_time TEXT DEFAULT '',
  publish_user TEXT DEFAULT '',
  last_review_id TEXT DEFAULT '',
  remark TEXT DEFAULT '',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_item_query
ON question_item(subject, grade, unit_id, question_type, difficulty, status, is_delete);

CREATE INDEX IF NOT EXISTS idx_question_item_modify_time
ON question_item(modify_time);

CREATE UNIQUE INDEX IF NOT EXISTS idx_question_item_code
ON question_item(question_code);

CREATE TABLE IF NOT EXISTS question_option (
  option_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  option_key TEXT NOT NULL DEFAULT '',
  option_order INTEGER NOT NULL DEFAULT 1,
  content_mode TEXT NOT NULL DEFAULT 'plain',
  option_html TEXT NOT NULL DEFAULT '',
  option_text TEXT NOT NULL DEFAULT '',
  is_correct TEXT NOT NULL DEFAULT '0',
  asset_count INTEGER NOT NULL DEFAULT 0,
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_option_question
ON question_option(question_id, option_order, is_delete);

CREATE TABLE IF NOT EXISTS question_answer (
  answer_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  answer_type TEXT NOT NULL DEFAULT '',
  answer_value TEXT NOT NULL DEFAULT '[]',
  answer_text TEXT DEFAULT '',
  reference_text TEXT DEFAULT '',
  case_sensitive TEXT NOT NULL DEFAULT '0',
  allow_order_change TEXT NOT NULL DEFAULT '0',
  auto_grading_rule TEXT DEFAULT '',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_answer_question
ON question_answer(question_id, is_delete);

CREATE TABLE IF NOT EXISTS question_blank_answer (
  blank_answer_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  blank_index INTEGER NOT NULL DEFAULT 1,
  standard_answer TEXT NOT NULL DEFAULT '',
  alternative_answers TEXT DEFAULT '[]',
  score INTEGER NOT NULL DEFAULT 0,
  match_mode TEXT NOT NULL DEFAULT 'exact',
  case_sensitive TEXT NOT NULL DEFAULT '0',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_blank_answer_question
ON question_blank_answer(question_id, blank_index, is_delete);

CREATE TABLE IF NOT EXISTS question_scoring_point (
  scoring_point_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  point_index INTEGER NOT NULL DEFAULT 1,
  point_text TEXT NOT NULL DEFAULT '',
  score INTEGER NOT NULL DEFAULT 0,
  keywords TEXT DEFAULT '[]',
  is_required TEXT NOT NULL DEFAULT '0',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_scoring_point_question
ON question_scoring_point(question_id, point_index, is_delete);

CREATE TABLE IF NOT EXISTS question_grade (
  grade_id TEXT PRIMARY KEY,
  stage TEXT NOT NULL DEFAULT '',
  semester TEXT NOT NULL DEFAULT 'upper',
  grade_code TEXT NOT NULL DEFAULT '',
  grade_name TEXT NOT NULL DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_question_grade_code
ON question_grade(grade_code);

CREATE INDEX IF NOT EXISTS idx_question_grade_query
ON question_grade(stage, semester, status, is_delete, order_index);

CREATE TABLE IF NOT EXISTS question_subject (
  subject_id TEXT PRIMARY KEY,
  subject_code TEXT NOT NULL DEFAULT '',
  subject_name TEXT NOT NULL DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_question_subject_code
ON question_subject(subject_code);

CREATE INDEX IF NOT EXISTS idx_question_subject_query
ON question_subject(status, is_delete, order_index);

CREATE TABLE IF NOT EXISTS question_unit (
  unit_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT NOT NULL DEFAULT '',
  grade TEXT NOT NULL DEFAULT '',
  textbook_version TEXT NOT NULL DEFAULT '',
  parent_id TEXT DEFAULT '',
  unit_code TEXT NOT NULL DEFAULT '',
  unit_name TEXT NOT NULL DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_unit_query
ON question_unit(subject, stage, grade, textbook_version, status, is_delete);

CREATE TABLE IF NOT EXISTS question_section (
  section_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT NOT NULL DEFAULT '',
  grade TEXT NOT NULL DEFAULT '',
  unit_id TEXT NOT NULL DEFAULT '',
  section_code TEXT NOT NULL DEFAULT '',
  section_name TEXT NOT NULL DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_section_query
ON question_section(subject, stage, grade, unit_id, status, is_delete);

CREATE TABLE IF NOT EXISTS question_knowledge (
  knowledge_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT DEFAULT '',
  grade TEXT DEFAULT '',
  parent_id TEXT DEFAULT '',
  knowledge_code TEXT NOT NULL DEFAULT '',
  knowledge_name TEXT NOT NULL DEFAULT '',
  knowledge_type TEXT NOT NULL DEFAULT '',
  knowledge_category TEXT NOT NULL DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_query
ON question_knowledge(subject, stage, grade, status, is_delete);

CREATE TABLE IF NOT EXISTS question_knowledge_point (
  point_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT NOT NULL DEFAULT '',
  grade TEXT NOT NULL DEFAULT '',
  unit_id TEXT NOT NULL DEFAULT '',
  section_id TEXT DEFAULT '',
  point_code TEXT NOT NULL DEFAULT '',
  point_name TEXT NOT NULL DEFAULT '',
  content_detail TEXT DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_point_query
ON question_knowledge_point(subject, stage, grade, unit_id, section_id, status, is_delete);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_point_code
ON question_knowledge_point(point_code);

CREATE TABLE IF NOT EXISTS question_knowledge_rel (
  rel_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  knowledge_id TEXT NOT NULL,
  knowledge_name TEXT DEFAULT '',
  order_index INTEGER NOT NULL DEFAULT 1,
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_rel_question
ON question_knowledge_rel(question_id, order_index);

CREATE TABLE IF NOT EXISTS question_asset (
  asset_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  usage_type TEXT NOT NULL DEFAULT 'stem',
  usage_ref TEXT DEFAULT '',
  asset_url TEXT NOT NULL DEFAULT '',
  asset_name TEXT DEFAULT '',
  mime_type TEXT DEFAULT '',
  file_size INTEGER NOT NULL DEFAULT 0,
  sha256 TEXT DEFAULT '',
  status TEXT NOT NULL DEFAULT 'enabled',
  is_delete TEXT NOT NULL DEFAULT '0',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_asset_question
ON question_asset(question_id, usage_type, is_delete);

CREATE TABLE IF NOT EXISTS question_review_record (
  review_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  from_status TEXT NOT NULL DEFAULT '',
  to_status TEXT NOT NULL DEFAULT '',
  review_result TEXT NOT NULL DEFAULT '',
  review_comment TEXT DEFAULT '',
  review_user TEXT NOT NULL DEFAULT '',
  review_time TEXT NOT NULL DEFAULT '',
  is_delete TEXT NOT NULL DEFAULT '0'
);

CREATE INDEX IF NOT EXISTS idx_question_review_record_question
ON question_review_record(question_id, review_time, is_delete);

CREATE TABLE IF NOT EXISTS question_change_log (
  log_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL,
  op_type TEXT NOT NULL DEFAULT '',
  before_json TEXT DEFAULT '{}',
  after_json TEXT DEFAULT '{}',
  op_user TEXT NOT NULL DEFAULT '',
  op_time TEXT NOT NULL DEFAULT '',
  remark TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_change_log_question
ON question_change_log(question_id, op_time);

CREATE TABLE IF NOT EXISTS question_import_batch (
  batch_id TEXT PRIMARY KEY,
  file_name TEXT NOT NULL DEFAULT '',
  file_url TEXT DEFAULT '',
  subject TEXT DEFAULT '',
  stage TEXT DEFAULT '',
  grade TEXT DEFAULT '',
  textbook_version TEXT DEFAULT '',
  status TEXT NOT NULL DEFAULT 'parsing',
  total_count INTEGER NOT NULL DEFAULT 0,
  success_count INTEGER NOT NULL DEFAULT 0,
  fail_count INTEGER NOT NULL DEFAULT 0,
  error_summary TEXT DEFAULT '',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT '',
  modify_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_import_batch_query
ON question_import_batch(status, create_time);

CREATE TABLE IF NOT EXISTS question_import_row (
  row_id TEXT PRIMARY KEY,
  batch_id TEXT NOT NULL,
  row_index INTEGER NOT NULL DEFAULT 1,
  raw_json TEXT NOT NULL DEFAULT '{}',
  parsed_json TEXT DEFAULT '{}',
  validate_status TEXT NOT NULL DEFAULT 'pending',
  error_msg TEXT DEFAULT '',
  question_id TEXT DEFAULT '',
  create_time TEXT NOT NULL DEFAULT '',
  create_user TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_import_row_batch
ON question_import_row(batch_id, row_index);

CREATE TABLE IF NOT EXISTS question_source_document (
  source_doc_id TEXT PRIMARY KEY,
  import_batch_id TEXT NOT NULL DEFAULT '',
  file_name TEXT NOT NULL DEFAULT '',
  file_sha256 TEXT DEFAULT '',
  file_url TEXT DEFAULT '',
  page_count INTEGER NOT NULL DEFAULT 0,
  subject TEXT DEFAULT '',
  stage TEXT DEFAULT '',
  grade TEXT DEFAULT '',
  textbook_version TEXT DEFAULT '',
  parse_status TEXT NOT NULL DEFAULT 'parsed',
  import_status TEXT NOT NULL DEFAULT 'preview',
  create_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_document_file
ON question_source_document(file_name, import_batch_id, create_time);

CREATE TABLE IF NOT EXISTS question_source_page (
  source_page_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  page_no INTEGER NOT NULL DEFAULT 1,
  page_image_url TEXT DEFAULT '',
  width INTEGER NOT NULL DEFAULT 0,
  height INTEGER NOT NULL DEFAULT 0,
  extract_service TEXT NOT NULL DEFAULT 'question.ai_pdf_text',
  extract_params_json TEXT DEFAULT '{}',
  raw_text TEXT DEFAULT '',
  raw_html TEXT DEFAULT '',
  extract_meta_json TEXT DEFAULT '{}',
  page_hash TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_page_doc
ON question_source_page(source_doc_id, page_no);

CREATE TABLE IF NOT EXISTS question_source_block (
  source_block_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  page_no INTEGER NOT NULL DEFAULT 1,
  block_order INTEGER NOT NULL DEFAULT 1,
  bbox_json TEXT DEFAULT '{}',
  block_type TEXT NOT NULL DEFAULT 'text',
  raw_text TEXT DEFAULT '',
  normalized_text TEXT DEFAULT '',
  block_image_url TEXT DEFAULT '',
  semantic_type TEXT DEFAULT 'knowledge_summary',
  confidence REAL NOT NULL DEFAULT 1,
  content_hash TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_block_doc
ON question_source_block(source_doc_id, page_no, block_order);

CREATE TABLE IF NOT EXISTS question_source_fragment (
  source_fragment_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  source_page_id TEXT DEFAULT '',
  source_block_id TEXT DEFAULT '',
  page_no INTEGER NOT NULL DEFAULT 1,
  fragment_order INTEGER NOT NULL DEFAULT 1,
  fragment_type TEXT NOT NULL DEFAULT 'knowledge_point',
  raw_text TEXT DEFAULT '',
  raw_html TEXT DEFAULT '',
  normalized_text TEXT DEFAULT '',
  char_start INTEGER NOT NULL DEFAULT 0,
  char_end INTEGER NOT NULL DEFAULT 0,
  context_before TEXT DEFAULT '',
  context_after TEXT DEFAULT '',
  bbox_json TEXT DEFAULT '{}',
  fragment_hash TEXT DEFAULT '',
  confidence REAL NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'matched',
  create_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_fragment_doc
ON question_source_fragment(source_doc_id, page_no, fragment_order);

CREATE INDEX IF NOT EXISTS idx_question_source_fragment_hash
ON question_source_fragment(fragment_hash);

CREATE TABLE IF NOT EXISTS question_source_snapshot (
  snapshot_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  source_page_id TEXT DEFAULT '',
  source_block_id TEXT DEFAULT '',
  source_fragment_id TEXT DEFAULT '',
  question_id TEXT DEFAULT '',
  knowledge_id TEXT DEFAULT '',
  content_id TEXT DEFAULT '',
  extract_service TEXT NOT NULL DEFAULT 'question.ai_pdf_text',
  extract_params_json TEXT DEFAULT '{}',
  raw_text TEXT DEFAULT '',
  raw_html TEXT DEFAULT '',
  normalized_text TEXT DEFAULT '',
  ai_output_json TEXT DEFAULT '{}',
  validator_result_json TEXT DEFAULT '{}',
  status TEXT NOT NULL DEFAULT 'active',
  create_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_snapshot_doc
ON question_source_snapshot(source_doc_id, knowledge_id, question_id, content_id);

CREATE TABLE IF NOT EXISTS question_knowledge_content (
  content_id TEXT PRIMARY KEY,
  batch_id TEXT DEFAULT '',
  source_doc_id TEXT DEFAULT '',
  source_block_id TEXT DEFAULT '',
  unit_id TEXT DEFAULT '',
  knowledge_id TEXT NOT NULL,
  semantic_type TEXT DEFAULT 'knowledge_summary',
  section_title TEXT DEFAULT '',
  content_text TEXT NOT NULL DEFAULT '',
  content_html TEXT DEFAULT '',
  content_json TEXT DEFAULT '{}',
  content_hash TEXT DEFAULT '',
  asset_count INTEGER NOT NULL DEFAULT 0,
  order_index INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'published',
  create_time TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_content_query
ON question_knowledge_content(source_doc_id, unit_id, knowledge_id, status);

CREATE INDEX IF NOT EXISTS idx_question_knowledge_content_hash
ON question_knowledge_content(content_hash);

CREATE TABLE IF NOT EXISTS question_source_rel (
  rel_id TEXT PRIMARY KEY,
  question_id TEXT NOT NULL DEFAULT '',
  source_doc_id TEXT NOT NULL,
  source_page_no INTEGER NOT NULL DEFAULT 1,
  source_block_id TEXT DEFAULT '',
  source_fragment_id TEXT DEFAULT '',
  source_content_id TEXT DEFAULT '',
  knowledge_id TEXT DEFAULT '',
  relation_type TEXT NOT NULL DEFAULT 'generated_from',
  confidence REAL NOT NULL DEFAULT 1,
  create_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_rel_doc
ON question_source_rel(source_doc_id, question_id, knowledge_id);

CREATE TABLE IF NOT EXISTS question_source_field_rel (
  field_rel_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  source_page_id TEXT DEFAULT '',
  source_block_id TEXT DEFAULT '',
  source_fragment_id TEXT DEFAULT '',
  entity_type TEXT NOT NULL DEFAULT '',
  entity_id TEXT NOT NULL DEFAULT '',
  field_name TEXT NOT NULL DEFAULT '',
  field_part_order INTEGER NOT NULL DEFAULT 1,
  extracted_value TEXT DEFAULT '',
  normalized_value TEXT DEFAULT '',
  raw_quote TEXT DEFAULT '',
  context_before TEXT DEFAULT '',
  context_after TEXT DEFAULT '',
  char_start INTEGER NOT NULL DEFAULT 0,
  char_end INTEGER NOT NULL DEFAULT 0,
  bbox_json TEXT DEFAULT '{}',
  confidence REAL NOT NULL DEFAULT 1,
  match_status TEXT NOT NULL DEFAULT 'matched',
  review_status TEXT NOT NULL DEFAULT 'pending',
  create_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_source_field_rel_entity
ON question_source_field_rel(entity_type, entity_id, field_name);

CREATE INDEX IF NOT EXISTS idx_question_source_field_rel_doc
ON question_source_field_rel(source_doc_id, source_fragment_id, match_status);

CREATE TABLE IF NOT EXISTS question_pdf_parse_issue (
  issue_id TEXT PRIMARY KEY,
  source_doc_id TEXT NOT NULL,
  page_no INTEGER NOT NULL DEFAULT 1,
  source_block_id TEXT DEFAULT '',
  issue_type TEXT NOT NULL DEFAULT '',
  severity TEXT NOT NULL DEFAULT 'warning',
  raw_text TEXT DEFAULT '',
  crop_image_url TEXT DEFAULT '',
  ai_output_json TEXT DEFAULT '{}',
  expected_schema TEXT DEFAULT '',
  error_msg TEXT DEFAULT '',
  suggestion TEXT DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  create_time TEXT NOT NULL DEFAULT '',
  modify_time TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_question_pdf_parse_issue_doc
ON question_pdf_parse_issue(source_doc_id, status, severity);

UPDATE sys_code
SET sys_code_type = 'knowledge_type_english'
WHERE sys_code_type = 'knowledge_type';

UPDATE sys_code
SET sys_code_type = 'knowledge_category_english'
WHERE sys_code_type = 'knowledge_category';

UPDATE sys_code
SET sys_code_text = '经典题型', order_index = 1
WHERE sys_code_type = 'question_category' AND sys_code = 'classic';

UPDATE sys_code
SET sys_code_text = '普通题型', order_index = 2
WHERE sys_code_type = 'question_category' AND sys_code = 'normal';

UPDATE sys_code
SET sys_code_text = '考试真题', order_index = 3
WHERE sys_code_type = 'question_category' AND sys_code = 'exam';

INSERT OR IGNORE INTO sys_code (
  sys_code_id, sys_code_type, sys_code_type_name, sys_code, sys_code_text, order_index, icon
) VALUES
  ('question-subject-english', 'subject', '科目', 'english', '英语', 1, NULL),
  ('question-subject-math', 'subject', '科目', 'math', '数学', 2, NULL),
  ('question-subject-chinese', 'subject', '科目', 'chinese', '语文', 3, NULL),
  ('question-stage-primary', 'study_stage', '学段', 'primary', '小学', 1, NULL),
  ('question-stage-junior', 'study_stage', '学段', 'junior', '初中', 2, NULL),
  ('question-stage-senior', 'study_stage', '学段', 'senior', '高中', 3, NULL),
  ('question-grade-semester-upper', 'grade_semester', '学期', 'upper', '上学期', 1, NULL),
  ('question-grade-semester-lower', 'grade_semester', '学期', 'lower', '下学期', 2, NULL),
  ('question-grade-1-upper', 'grade', '年级', 'grade_1', '一年级上学期', 11, NULL),
  ('question-grade-1-lower', 'grade', '年级', 'grade_1_lower', '一年级下学期', 12, NULL),
  ('question-grade-2-upper', 'grade', '年级', 'grade_2', '二年级上学期', 21, NULL),
  ('question-grade-2-lower', 'grade', '年级', 'grade_2_lower', '二年级下学期', 22, NULL),
  ('question-grade-3-upper', 'grade', '年级', 'grade_3', '三年级上学期', 31, NULL),
  ('question-grade-3-lower', 'grade', '年级', 'grade_3_lower', '三年级下学期', 32, NULL),
  ('question-grade-4-upper', 'grade', '年级', 'grade_4', '四年级上学期', 41, NULL),
  ('question-grade-4-lower', 'grade', '年级', 'grade_4_lower', '四年级下学期', 42, NULL),
  ('question-grade-5-upper', 'grade', '年级', 'grade_5', '五年级上学期', 51, NULL),
  ('question-grade-5-lower', 'grade', '年级', 'grade_5_lower', '五年级下学期', 52, NULL),
  ('question-grade-6-upper', 'grade', '年级', 'grade_6', '六年级上学期', 61, NULL),
  ('question-grade-6-lower', 'grade', '年级', 'grade_6_lower', '六年级下学期', 62, NULL),
  ('question-grade-7', 'grade', '年级', 'grade_7', '七年级上学期', 71, NULL),
  ('question-grade-7-lower', 'grade', '年级', 'grade_7_lower', '七年级下学期', 72, NULL),
  ('question-grade-8', 'grade', '年级', 'grade_8', '八年级上学期', 81, NULL),
  ('question-grade-8-lower', 'grade', '年级', 'grade_8_lower', '八年级下学期', 82, NULL),
  ('question-grade-9', 'grade', '年级', 'grade_9', '九年级上学期', 91, NULL),
  ('question-grade-9-lower', 'grade', '年级', 'grade_9_lower', '九年级下学期', 92, NULL),
  ('question-grade-10-upper', 'grade', '年级', 'grade_10', '高一上学期', 101, NULL),
  ('question-grade-10-lower', 'grade', '年级', 'grade_10_lower', '高一下学期', 102, NULL),
  ('question-grade-11-upper', 'grade', '年级', 'grade_11', '高二上学期', 111, NULL),
  ('question-grade-11-lower', 'grade', '年级', 'grade_11_lower', '高二下学期', 112, NULL),
  ('question-grade-12-upper', 'grade', '年级', 'grade_12', '高三上学期', 121, NULL),
  ('question-grade-12-lower', 'grade', '年级', 'grade_12_lower', '高三下学期', 122, NULL),
  ('question-textbook-pep', 'textbook_version', '教材版本', 'pep', '人教版', 1, NULL),
  ('question-textbook-bsd', 'textbook_version', '教材版本', 'bsd', '北师大版', 2, NULL),
  ('question-type-single', 'question_type', '题型', 'single_choice', '单选题', 1, NULL),
  ('question-type-multiple', 'question_type', '题型', 'multiple_choice', '多选题', 2, NULL),
  ('question-type-judge', 'question_type', '题型', 'judge', '判断题', 3, NULL),
  ('question-type-blank', 'question_type', '题型', 'blank', '填空题', 4, NULL),
  ('question-type-short', 'question_type', '题型', 'short_answer', '简答题', 5, NULL),
  ('question-type-calc', 'question_type', '题型', 'calculation', '计算题', 6, NULL),
  ('question-category-classic', 'question_category', '题目属性', 'classic', '经典题型', 1, NULL),
  ('question-category-normal', 'question_category', '题目属性', 'normal', '普通题型', 2, NULL),
  ('question-category-exam', 'question_category', '题目属性', 'exam', '考试真题', 3, NULL),
  ('knowledge-type-vocabulary', 'knowledge_type_english', '目录类型', 'vocabulary', '词汇', 1, NULL),
  ('knowledge-type-grammar-rule', 'knowledge_type_english', '目录类型', 'grammar_rule', '语法规则', 2, NULL),
  ('knowledge-category-cn-to-en', 'knowledge_category_english', '目录具体分类', 'cn_to_en', '中译英', 1, 'vocabulary'),
  ('knowledge-category-en-to-cn', 'knowledge_category_english', '目录具体分类', 'en_to_cn', '英译中', 2, 'vocabulary'),
  ('knowledge-category-noun-singular-plural', 'knowledge_category_english', '目录具体分类', 'noun_singular_plural', '名词单复数', 1, 'grammar_rule'),
  ('knowledge-category-personal-pronoun', 'knowledge_category_english', '目录具体分类', 'personal_pronoun', '人称代词', 2, 'grammar_rule'),
  ('knowledge-category-verb-tense', 'knowledge_category_english', '目录具体分类', 'verb_tense', '动词时态', 3, 'grammar_rule'),
  ('question-difficulty-basic', 'question_difficulty', '题目难度', 'basic', '基础', 1, NULL),
  ('question-difficulty-medium', 'question_difficulty', '题目难度', 'medium', '进阶', 2, NULL),
  ('question-difficulty-hard', 'question_difficulty', '题目难度', 'hard', '困难', 3, NULL),
  ('question-status-draft', 'question_status', '题目状态', 'draft', '草稿', 1, NULL),
  ('question-status-reviewing', 'question_status', '题目状态', 'reviewing', '待审核', 2, NULL),
  ('question-status-published', 'question_status', '题目状态', 'published', '已发布', 3, NULL),
  ('question-status-offline', 'question_status', '题目状态', 'offline', '已下线', 4, NULL),
  ('question-status-rejected', 'question_status', '题目状态', 'rejected', '已退回', 5, NULL),
  ('question-content-plain', 'content_mode', '内容模式', 'plain', '普通文本', 1, NULL),
  ('question-content-rich', 'content_mode', '内容模式', 'rich', '富文本', 2, NULL),
  ('question-review-submit', 'review_result', '审核结果', 'submit', '提交审核', 1, NULL),
  ('question-review-publish', 'review_result', '审核结果', 'publish', '发布', 2, NULL),
  ('question-review-reject', 'review_result', '审核结果', 'reject', '退回', 3, NULL),
  ('question-review-offline', 'review_result', '审核结果', 'offline', '下线', 4, NULL),
  ('question-import-parsing', 'import_status', '导入状态', 'parsing', '解析中', 1, NULL),
  ('question-import-validated', 'import_status', '导入状态', 'validated', '已校验', 2, NULL),
  ('question-import-committed', 'import_status', '导入状态', 'committed', '已入库', 3, NULL),
  ('question-import-failed', 'import_status', '导入状态', 'failed', '失败', 4, NULL);

INSERT OR IGNORE INTO question_grade (
  grade_id, stage, semester, grade_code, grade_name, order_index, status, is_delete,
  create_time, create_user, modify_time, modify_user
) VALUES
  ('grade-primary-1', 'primary', 'upper', 'grade_1', '一年级', 11, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-1-lower', 'primary', 'lower', 'grade_1_lower', '一年级', 12, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-2', 'primary', 'upper', 'grade_2', '二年级', 21, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-2-lower', 'primary', 'lower', 'grade_2_lower', '二年级', 22, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-3', 'primary', 'upper', 'grade_3', '三年级', 31, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-3-lower', 'primary', 'lower', 'grade_3_lower', '三年级', 32, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-4', 'primary', 'upper', 'grade_4', '四年级', 41, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-4-lower', 'primary', 'lower', 'grade_4_lower', '四年级', 42, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-5', 'primary', 'upper', 'grade_5', '五年级', 51, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-5-lower', 'primary', 'lower', 'grade_5_lower', '五年级', 52, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-6', 'primary', 'upper', 'grade_6', '六年级', 61, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-6-lower', 'primary', 'lower', 'grade_6_lower', '六年级', 62, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-7', 'junior', 'upper', 'grade_7', '七年级', 71, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-7-lower', 'junior', 'lower', 'grade_7_lower', '七年级', 72, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-8', 'junior', 'upper', 'grade_8', '八年级', 81, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-8-lower', 'junior', 'lower', 'grade_8_lower', '八年级', 82, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-9', 'junior', 'upper', 'grade_9', '九年级', 91, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-9-lower', 'junior', 'lower', 'grade_9_lower', '九年级', 92, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-10', 'senior', 'upper', 'grade_10', '高一', 101, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-10-lower', 'senior', 'lower', 'grade_10_lower', '高一', 102, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-11', 'senior', 'upper', 'grade_11', '高二', 111, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-11-lower', 'senior', 'lower', 'grade_11_lower', '高二', 112, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-12', 'senior', 'upper', 'grade_12', '高三', 121, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-senior-12-lower', 'senior', 'lower', 'grade_12_lower', '高三', 122, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

UPDATE sys_code
SET sys_code_text = CASE sys_code
    WHEN 'grade_7' THEN '七年级上学期'
    WHEN 'grade_8' THEN '八年级上学期'
    WHEN 'grade_9' THEN '九年级上学期'
    ELSE sys_code_text
  END,
  order_index = CASE sys_code
    WHEN 'grade_7' THEN 71
    WHEN 'grade_8' THEN 81
    WHEN 'grade_9' THEN 91
    ELSE order_index
  END
WHERE sys_code_type = 'grade' AND sys_code IN ('grade_7', 'grade_8', 'grade_9');

UPDATE question_grade
SET stage = CASE grade_code
    WHEN 'grade_1' THEN 'primary'
    WHEN 'grade_2' THEN 'primary'
    WHEN 'grade_3' THEN 'primary'
    WHEN 'grade_4' THEN 'primary'
    WHEN 'grade_5' THEN 'primary'
    WHEN 'grade_6' THEN 'primary'
    WHEN 'grade_7' THEN 'junior'
    WHEN 'grade_8' THEN 'junior'
    WHEN 'grade_9' THEN 'junior'
    ELSE stage
  END,
  semester = 'upper',
  grade_name = CASE grade_code
    WHEN 'grade_1' THEN '一年级'
    WHEN 'grade_2' THEN '二年级'
    WHEN 'grade_3' THEN '三年级'
    WHEN 'grade_4' THEN '四年级'
    WHEN 'grade_5' THEN '五年级'
    WHEN 'grade_6' THEN '六年级'
    WHEN 'grade_7' THEN '七年级'
    WHEN 'grade_8' THEN '八年级'
    WHEN 'grade_9' THEN '九年级'
    ELSE grade_name
  END,
  order_index = CASE grade_code
    WHEN 'grade_1' THEN 11
    WHEN 'grade_2' THEN 21
    WHEN 'grade_3' THEN 31
    WHEN 'grade_4' THEN 41
    WHEN 'grade_5' THEN 51
    WHEN 'grade_6' THEN 61
    WHEN 'grade_7' THEN 71
    WHEN 'grade_8' THEN 81
    WHEN 'grade_9' THEN 91
    ELSE order_index
  END
WHERE grade_code IN ('grade_1', 'grade_2', 'grade_3', 'grade_4', 'grade_5', 'grade_6', 'grade_7', 'grade_8', 'grade_9');

INSERT OR IGNORE INTO question_subject (
  subject_id, subject_code, subject_name, order_index, status, is_delete,
  create_time, create_user, modify_time, modify_user
) VALUES
  ('subject-chinese', 'chinese', '语文', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('subject-math', 'math', '数学', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('subject-english', 'english', '英语', 3, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

INSERT OR IGNORE INTO question_unit (
  unit_id, subject, stage, grade, textbook_version, parent_id, unit_code, unit_name,
  order_index, status, is_delete, create_time, create_user, modify_time, modify_user
) VALUES
  ('unit-english-grade8-pep-u1', 'english', 'junior', 'grade_8', 'pep', '', 'unit_1', 'Unit 1 Friendship', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('unit-english-grade8-pep-u2', 'english', 'junior', 'grade_8', 'pep', '', 'unit_2', 'Unit 2 Daily Life', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

INSERT OR IGNORE INTO question_section (
  section_id, subject, stage, grade, unit_id, section_code, section_name,
  order_index, status, is_delete, create_time, create_user, modify_time, modify_user
) VALUES
  ('section-english-grade8-pep-u2-p1', 'english', 'junior', 'grade_8', 'unit-english-grade8-pep-u2', 'part_1', 'Part 1', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('section-english-grade8-pep-u2-p2', 'english', 'junior', 'grade_8', 'unit-english-grade8-pep-u2', 'part_2', 'Part 2', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

INSERT OR IGNORE INTO question_knowledge (
  knowledge_id, subject, stage, grade, parent_id, knowledge_code, knowledge_name,
  knowledge_type, knowledge_category, order_index, status, is_delete, create_time, create_user, modify_time, modify_user
) VALUES
  ('knowledge-english-vocabulary', 'english', 'junior', '', '', 'english_vocabulary', '词汇', 'vocabulary', 'cn_to_en', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('knowledge-english-grammar', 'english', 'junior', '', '', 'english_grammar', '语法', 'grammar_rule', 'verb_tense', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('knowledge-english-reading', 'english', 'junior', '', '', 'english_reading', '阅读理解', '', '', 3, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

INSERT OR IGNORE INTO question_knowledge_point (
  point_id, subject, stage, grade, unit_id, section_id, point_code, point_name,
  content_detail, order_index, status, is_delete, create_time, create_user, modify_time, modify_user
) VALUES
  ('point-english-grade8-u2-greeting', 'english', 'junior', 'grade_8', 'unit-english-grade8-pep-u2', 'section-english-grade8-pep-u2-p1', 'daily_greeting', '日常问候', '掌握日常问候、介绍和告别相关表达。', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('point-english-grade8-u2-time-expression', 'english', 'junior', 'grade_8', 'unit-english-grade8-pep-u2', 'section-english-grade8-pep-u2-p2', 'time_expression', '时间表达', '掌握频率、时间状语和日常作息描述。', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

DELETE FROM role_menu
WHERE sys_menu_id IN (
  'menu-question-bank',
  'menu-question-taxonomy',
  'menu-question-knowledge-point',
  'menu-question-editor',
  'menu-question-review',
  'menu-question-import'
);

DELETE FROM sys_menu
WHERE sys_menu_id IN (
  'menu-question-bank',
  'menu-question-taxonomy',
  'menu-question-knowledge-point',
  'menu-question-editor',
  'menu-question-review',
  'menu-question-import'
);

INSERT OR IGNORE INTO sys_menu (
  sys_menu_id, menu_type, menu_name, menu_code, icon, is_index, group_path,
  router_group, group_api, api, data, url, in_menu, is_common, parent_id,
  create_time, create_user, order_index, description, belong_project, type
) VALUES
  ('menu-question-root', '2', '题库管理', 'question_root', 'FaBookOpen', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.question_bank', 'question/question_bank.json', '/framework/question-bank',
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 200, 'Web 题库管理单页', 'ai-study-admin', NULL
  ),
  ('menu-question-taxonomy', '2', '年级科目维护', 'question_taxonomy', 'FaLayerGroup', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.question_taxonomy', 'question/question_taxonomy.json', '/framework/question-taxonomy',
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 210, '年级、科目、单元、知识点维护', 'ai-study-admin', NULL
  );

UPDATE sys_menu
SET menu_type = '2',
    menu_name = '题库管理',
    menu_code = 'question_root',
    icon = 'FaBookOpen',
    router_group = 'framework',
    api = 'post:/template_data/data?service=frontend.question_bank',
    data = 'question/question_bank.json',
    url = '/framework/question-bank',
    in_menu = '1',
    parent_id = '',
    order_index = 200,
    description = 'Web 题库管理单页'
WHERE sys_menu_id = 'menu-question-root';

INSERT OR IGNORE INTO sys_menu (
  sys_menu_id, menu_type, menu_name, menu_code, icon, is_index, group_path,
  router_group, group_api, api, data, url, in_menu, is_common, parent_id,
  create_time, create_user, order_index, description, belong_project, type
) VALUES
  ('menu-question-taxonomy', '2', '年级科目维护', 'question_taxonomy', 'FaLayerGroup', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.question_taxonomy', 'question/question_taxonomy.json', '/framework/question-taxonomy',
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 210, '年级、科目、单元、知识点维护', 'ai-study-admin', NULL
  );

UPDATE sys_menu
SET menu_type = '2',
    menu_name = '年级科目维护',
    menu_code = 'question_taxonomy',
    icon = 'FaLayerGroup',
    router_group = 'framework',
    api = 'post:/template_data/data?service=frontend.question_taxonomy',
    data = 'question/question_taxonomy.json',
    url = '/framework/question-taxonomy',
    in_menu = '1',
    parent_id = '',
    order_index = 210,
    description = '年级、科目、单元、知识点维护'
WHERE sys_menu_id = 'menu-question-taxonomy';

INSERT OR IGNORE INTO sys_menu (
  sys_menu_id, menu_type, menu_name, menu_code, icon, is_index, group_path,
  router_group, group_api, api, data, url, in_menu, is_common, parent_id,
  create_time, create_user, order_index, description, belong_project, type
) VALUES
  ('menu-question-knowledge-point', '2', '知识点维护', 'question_knowledge_point', 'FaBookReader', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.question_knowledge_point', 'question/question_knowledge_point.json', '/framework/question-knowledge-point',
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 220, '知识点内容维护', 'ai-study-admin', NULL
  );

UPDATE sys_menu
SET menu_type = '2',
    menu_name = '知识点维护',
    menu_code = 'question_knowledge_point',
    icon = 'FaBookReader',
    router_group = 'framework',
    api = 'post:/template_data/data?service=frontend.question_knowledge_point',
    data = 'question/question_knowledge_point.json',
    url = '/framework/question-knowledge-point',
    in_menu = '1',
    parent_id = '',
    order_index = 220,
    description = '知识点内容维护'
WHERE sys_menu_id = 'menu-question-knowledge-point';

INSERT OR IGNORE INTO role_menu (
  role_menu_id, role_id, sys_menu_id, belong_project
) VALUES
  ('admin-menu-question-root', 'admin', 'menu-question-root', 'ai-study-admin'),
  ('admin-menu-question-taxonomy', 'admin', 'menu-question-taxonomy', 'ai-study-admin'),
  ('admin-menu-question-knowledge-point', 'admin', 'menu-question-knowledge-point', 'ai-study-admin');

COMMIT;
