PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS question_import_row;
DROP TABLE IF EXISTS question_import_batch;
DROP TABLE IF EXISTS question_change_log;
DROP TABLE IF EXISTS question_review_record;
DROP TABLE IF EXISTS question_asset;
DROP TABLE IF EXISTS question_knowledge_rel;
DROP TABLE IF EXISTS question_knowledge;
DROP TABLE IF EXISTS question_unit;
DROP TABLE IF EXISTS question_subject;
DROP TABLE IF EXISTS question_grade;
DROP TABLE IF EXISTS question_scoring_point;
DROP TABLE IF EXISTS question_blank_answer;
DROP TABLE IF EXISTS question_answer;
DROP TABLE IF EXISTS question_option;
DROP TABLE IF EXISTS question_item;

BEGIN TRANSACTION;

DELETE FROM btn_role_id_list;
DELETE FROM sys_btn;
DELETE FROM role_menu;
DELETE FROM user_role_id_list;
DELETE FROM role;
DELETE FROM user_change_history;
DELETE FROM user_account;
DELETE FROM sys_menu;
DELETE FROM sys_code;
DELETE FROM schema_page_data;
DELETE FROM schema_page_field;
DELETE FROM schema_page;

INSERT INTO sys_code (
  sys_code_id, sys_code_type, sys_code_type_name, sys_code, sys_code_text, order_index, icon
) VALUES
  ('code-menu-type-dir', 'menu_type', '菜单类型', '1', '目录', 1, NULL),
  ('code-menu-type-page', 'menu_type', '菜单类型', '2', '菜单', 2, NULL),
  ('code-sex-man', 'sex', '性别', 'man', '男', 1, NULL),
  ('code-sex-woman', 'sex', '性别', 'woman', '女', 2, NULL),
  ('code-user-status-normal', 'user_job_status', '用户状态', 'normal', '正常', 1, NULL),
  ('code-user-status-trial', 'user_job_status', '用户状态', 'trial', '试用', 2, NULL),
  ('code-user-status-leave', 'user_job_status', '用户状态', 'leave', '停用', 3, NULL);

INSERT INTO role (
  role_id, role_name, order_index, role_code
) VALUES
  ('admin', '管理员', 1, 'admin');

INSERT INTO user_account (
  user_id, nick, username, password, status, email, phone, is_delete,
  create_time, create_user, modify_time, modify_user, work_code, sex
) VALUES
  (
    'admin-user',
    '管理员',
    'admin',
    'e10adc3949ba59abbe56e057f20f883e',
    'normal',
    '',
    '',
    '0',
    '2026-05-18 00:00:00',
    'system',
    '2026-05-18 00:00:00',
    'system',
    'admin',
    'man'
  );

INSERT INTO user_role_id_list (
  role_list_id, user_id, role_id, create_time, modify_time, order_weight, user_group_id, user_role_type, belong_project
) VALUES
  ('admin-user-role-admin', 'admin-user', 'admin', '2026-05-18 00:00:00', '2026-05-18 00:00:00', 1, '', '', 'ai-study-admin');

INSERT INTO sys_menu (
  sys_menu_id, menu_type, menu_name, menu_code, icon, is_index, group_path,
  router_group, group_api, api, data, url, in_menu, is_common, parent_id,
  create_time, create_user, order_index, description, belong_project, type
) VALUES
  (
    'menu-login', '2', '登录', 'login', NULL, '0', NULL,
    NULL, NULL, 'post:/template_data/data?service=frontend.login', 'framework/login.json', '/login',
    '0', '1', '',
    '2026-05-18 00:00:00', 'system', 10, '后台登录页', 'ai-study-admin', NULL
  ),
  (
    'menu-framework', 'framework', '框架', 'framework', NULL, '0', '/framework',
    NULL, NULL, 'post:/template_data/data?service=frontend.framework', 'framework/framework.json', '/framework',
    '0', '0', '',
    '2026-05-18 00:00:00', 'system', 20, '后台框架路由', 'ai-study-admin', NULL
  ),
  (
    'menu-home', '2', '首页', 'first', 'FaTachometerAlt', '1', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.home', 'framework/home.json', '/framework',
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 30, '后台首页', 'ai-study-admin', NULL
  ),
  (
    'menu-system', '1', '系统管理', 'system_admin', 'FaCog', '0', NULL,
    NULL, NULL, NULL, NULL, NULL,
    '1', '0', '',
    '2026-05-18 00:00:00', 'system', 100, '系统管理目录', 'ai-study-admin', NULL
  ),
  (
    'menu-user', '2', '用户管理', 'user', 'FaUsers', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.user', 'system/user.json', '/framework/user',
    '1', '0', 'menu-system',
    '2026-05-18 00:00:00', 'system', 110, '用户管理', 'ai-study-admin', NULL
  ),
  (
    'menu-role', '2', '角色管理', 'role', 'FaUserShield', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.role', 'system/role.json', '/framework/role',
    '1', '0', 'menu-system',
    '2026-05-18 00:00:00', 'system', 120, '角色管理', 'ai-study-admin', NULL
  ),
  (
    'menu-menu-manage', '2', '菜单管理', 'menu_manage', 'FaListUl', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.menu_manage', 'system/menu_manage.json', '/framework/menu_manage',
    '1', '0', 'menu-system',
    '2026-05-18 00:00:00', 'system', 130, '菜单管理', 'ai-study-admin', NULL
  ),
  (
    'menu-sys-code', '2', '码表管理', 'sys_code', 'FaCode', '0', NULL,
    'framework', NULL, 'post:/template_data/data?service=frontend.sys_code', 'system/sys_code.json', '/framework/sys_code',
    '1', '0', 'menu-system',
    '2026-05-18 00:00:00', 'system', 140, '码表管理', 'ai-study-admin', NULL
  );

INSERT INTO role_menu (
  role_menu_id, role_id, sys_menu_id, belong_project
) VALUES
  ('admin-menu-framework', 'admin', 'menu-framework', 'ai-study-admin'),
  ('admin-menu-home', 'admin', 'menu-home', 'ai-study-admin'),
  ('admin-menu-system', 'admin', 'menu-system', 'ai-study-admin'),
  ('admin-menu-user', 'admin', 'menu-user', 'ai-study-admin'),
  ('admin-menu-role', 'admin', 'menu-role', 'ai-study-admin'),
  ('admin-menu-menu-manage', 'admin', 'menu-menu-manage', 'ai-study-admin'),
  ('admin-menu-sys-code', 'admin', 'menu-sys-code', 'ai-study-admin');

COMMIT;

.read scripts/migrate-question-web.sql
