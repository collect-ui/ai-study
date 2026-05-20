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
ON question_grade(stage, status, is_delete, order_index);

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

CREATE TABLE IF NOT EXISTS question_knowledge (
  knowledge_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL DEFAULT '',
  stage TEXT DEFAULT '',
  grade TEXT DEFAULT '',
  parent_id TEXT DEFAULT '',
  knowledge_code TEXT NOT NULL DEFAULT '',
  knowledge_name TEXT NOT NULL DEFAULT '',
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

INSERT OR IGNORE INTO sys_code (
  sys_code_id, sys_code_type, sys_code_type_name, sys_code, sys_code_text, order_index, icon
) VALUES
  ('question-subject-english', 'subject', '科目', 'english', '英语', 1, NULL),
  ('question-subject-math', 'subject', '科目', 'math', '数学', 2, NULL),
  ('question-subject-chinese', 'subject', '科目', 'chinese', '语文', 3, NULL),
  ('question-stage-primary', 'study_stage', '学段', 'primary', '小学', 1, NULL),
  ('question-stage-junior', 'study_stage', '学段', 'junior', '初中', 2, NULL),
  ('question-stage-senior', 'study_stage', '学段', 'senior', '高中', 3, NULL),
  ('question-grade-7', 'grade', '年级', 'grade_7', '七年级', 7, NULL),
  ('question-grade-8', 'grade', '年级', 'grade_8', '八年级', 8, NULL),
  ('question-grade-9', 'grade', '年级', 'grade_9', '九年级', 9, NULL),
  ('question-textbook-pep', 'textbook_version', '教材版本', 'pep', '人教版', 1, NULL),
  ('question-textbook-bsd', 'textbook_version', '教材版本', 'bsd', '北师大版', 2, NULL),
  ('question-type-single', 'question_type', '题型', 'single_choice', '单选题', 1, NULL),
  ('question-type-multiple', 'question_type', '题型', 'multiple_choice', '多选题', 2, NULL),
  ('question-type-judge', 'question_type', '题型', 'judge', '判断题', 3, NULL),
  ('question-type-blank', 'question_type', '题型', 'blank', '填空题', 4, NULL),
  ('question-type-short', 'question_type', '题型', 'short_answer', '简答题', 5, NULL),
  ('question-type-calc', 'question_type', '题型', 'calculation', '计算题', 6, NULL),
  ('question-category-normal', 'question_category', '题目属性', 'normal', '普通题型', 1, NULL),
  ('question-category-classic', 'question_category', '题目属性', 'classic', '经典题型', 2, NULL),
  ('question-category-exam', 'question_category', '题目属性', 'exam', '考试题型', 3, NULL),
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
  grade_id, stage, grade_code, grade_name, order_index, status, is_delete,
  create_time, create_user, modify_time, modify_user
) VALUES
  ('grade-primary-1', 'primary', 'grade_1', '一年级', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-2', 'primary', 'grade_2', '二年级', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-3', 'primary', 'grade_3', '三年级', 3, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-4', 'primary', 'grade_4', '四年级', 4, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-5', 'primary', 'grade_5', '五年级', 5, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-primary-6', 'primary', 'grade_6', '六年级', 6, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-7', 'junior', 'grade_7', '初一', 7, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-8', 'junior', 'grade_8', '初二', 8, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('grade-junior-9', 'junior', 'grade_9', '初三', 9, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

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

INSERT OR IGNORE INTO question_knowledge (
  knowledge_id, subject, stage, grade, parent_id, knowledge_code, knowledge_name,
  order_index, status, is_delete, create_time, create_user, modify_time, modify_user
) VALUES
  ('knowledge-english-vocabulary', 'english', 'junior', '', '', 'english_vocabulary', '词汇', 1, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('knowledge-english-grammar', 'english', 'junior', '', '', 'english_grammar', '语法', 2, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system'),
  ('knowledge-english-reading', 'english', 'junior', '', '', 'english_reading', '阅读理解', 3, 'enabled', '0', '2026-05-18 00:00:00', 'system', '2026-05-18 00:00:00', 'system');

DELETE FROM role_menu
WHERE sys_menu_id IN (
  'menu-question-bank',
  'menu-question-taxonomy',
  'menu-question-editor',
  'menu-question-review',
  'menu-question-import'
);

DELETE FROM sys_menu
WHERE sys_menu_id IN (
  'menu-question-bank',
  'menu-question-taxonomy',
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

INSERT OR IGNORE INTO role_menu (
  role_menu_id, role_id, sys_menu_id, belong_project
) VALUES
  ('admin-menu-question-root', 'admin', 'menu-question-root', 'ai-study-admin'),
  ('admin-menu-question-taxonomy', 'admin', 'menu-question-taxonomy', 'ai-study-admin');

COMMIT;
