# AI Study SQLite 到 MySQL 迁移测试验收计划

## 1. 目标

将 `/data/project/ai-study/backend/database/ai_study_admin.db` 当前 SQLite 主库迁移到 MySQL，MySQL 连接实例、账号和密码与 `/data/project/auto-check` 当前主库配置保持一致。

迁移完成后，必须完成数据库实例重建、数据迁移、后台启动、接口验证、页面验证和回滚验证。当前后台实际菜单为 6 个功能页面加首页，全部纳入验收：

| 类型 | 菜单 | 路由 | 页面配置 |
| --- | --- | --- | --- |
| 首页 | 首页 | `/collect-ui/#/collect-ui/framework` | `framework/home.json` |
| 系统 | 用户管理 | `/collect-ui/#/collect-ui/framework/user` | `system/user.json` |
| 系统 | 角色管理 | `/collect-ui/#/collect-ui/framework/role` | `system/role.json` |
| 系统 | 菜单管理 | `/collect-ui/#/collect-ui/framework/menu_manage` | `system/menu_manage.json` |
| 系统 | 码表管理 | `/collect-ui/#/collect-ui/framework/sys_code` | `system/sys_code.json` |
| 题库 | 题库管理 | `/collect-ui/#/collect-ui/framework/question-bank` | `question/question_bank.json` |
| 题库 | 年级科目维护 | `/collect-ui/#/collect-ui/framework/question-taxonomy` | `question/question_taxonomy.json` |

登录页 `/collect-ui/#/collect-ui/login` 作为前置验证，不算业务菜单，但必须验证登录成功和登录后菜单加载。

## 2. 约束

1. 不在文档、脚本、日志、提交信息中写入 MySQL 明文密码。
2. MySQL 账号、密码、host、port 从 `/data/project/auto-check/conf/application.properties` 的当前主库配置读取，或使用运行环境中的安全变量；执行过程禁止 `set -x`。
3. AI Study 目标库使用独立库名 `ai_study`；如已有同名库，本轮验收必须先备份再重建。
4. 迁移前必须保留 SQLite 文件备份和 MySQL dump 备份；失败时能回到 SQLite 配置。
5. 页面验收必须采集 `console.error`、`pageerror`、关键资源失败、接口 HTTP 失败和接口业务失败。

## 3. 验收环境

| 项 | 值 |
| --- | --- |
| 项目目录 | `/data/project/ai-study` |
| 后端目录 | `/data/project/ai-study/backend` |
| 原 SQLite | `backend/database/ai_study_admin.db` |
| 目标 MySQL 库 | `ai_study` |
| 后台端口 | `8026` |
| 本地访问 | `http://127.0.0.1:8026/collect-ui/` |
| 登录账号 | 沿用当前 SQLite 管理员账号 |
| MySQL 凭据 | 与 `/data/project/auto-check` 当前主库配置一致，不落明文 |

## 4. 迁移前检查

### 4.1 代码与配置检查

```bash
cd /data/project/ai-study/backend

go test ./...

sqlite3 database/ai_study_admin.db "PRAGMA integrity_check;"
sqlite3 database/ai_study_admin.db ".tables"
sqlite3 database/ai_study_admin.db "select count(*) from user_account; select count(*) from sys_menu; select count(*) from question_item;"
```

通过标准：

- `go test ./...` 通过。
- SQLite `integrity_check` 返回 `ok`。
- `user_account`、`sys_menu`、`role`、`role_menu`、`sys_code`、`question_*` 核心表均存在。
- 当前管理员账号存在且可登录。

### 4.2 MySQL 兼容性预检

本项目当前 SQL 中已发现需要重点确认的 MySQL 兼容点：

- `backend/scripts/init-ai-study-admin.sql` 和 `backend/scripts/migrate-question-web.sql` 包含 SQLite 专用语句：`PRAGMA`、`.read`、`BEGIN TRANSACTION`。
- `question_detail.sql`、`question_choice_detail.sql` 使用 SQLite JSON 聚合函数 `json_group_array/json_object`，MySQL 需验证或改写为 MySQL 等价写法。
- `question_stats.sql` 使用 `date('now', 'localtime')`，MySQL 需验证或改写为 `curdate()`/`date(now())`。
- `schema_page_field/get_page_schema_field.sql` 使用 `CAST(... AS INTEGER)`，MySQL 建议改为 `CAST(... AS SIGNED)`。

通过标准：

- 所有页面涉及的 SQL 在 MySQL 下无语法错误。
- 低代码服务返回结构与 SQLite 版本一致，前端无需因为数据库切换修改页面 schema。

## 5. 数据库实例重建

### 5.1 备份

```bash
cd /data/project/ai-study
mkdir -p reports/mysql_migration

cp backend/database/ai_study_admin.db \
  reports/mysql_migration/ai_study_admin.sqlite.before_mysql_$(date +%Y%m%d-%H%M%S).db
```

如果目标 MySQL 库已经存在，先做 dump 备份。密码通过环境变量或安全配置传入，不在命令中明文展开。

```bash
mysqldump --defaults-extra-file="$AI_STUDY_MYSQL_CNF" \
  --single-transaction --routines --triggers ai_study \
  > reports/mysql_migration/ai_study.before_recreate_$(date +%Y%m%d-%H%M%S).sql
```

### 5.2 重建目标库

```sql
DROP DATABASE IF EXISTS ai_study;
CREATE DATABASE ai_study
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
```

通过标准：

- 目标库是全新空库。
- 字符集为 `utf8mb4`。
- 使用的 MySQL 用户、密码、host、port 与 auto-check 当前主库一致。

## 6. 数据迁移

建议分两步执行：先建表，再导入数据。

### 6.1 建表

优先使用 Go model/GORM 或专用 MySQL schema 初始化脚本生成表结构。不要直接把 SQLite 初始化脚本原样导入 MySQL。

通过标准：

- `show tables` 包含所有注册模型表：
  `user_account`、`role`、`role_menu`、`user_role_id_list`、`sys_menu`、`sys_code`、`schema_page*`、`sys_btn`、`btn_role_id_list`、`question_item`、`question_option`、`question_answer`、`question_blank_answer`、`question_scoring_point`、`question_grade`、`question_subject`、`question_unit`、`question_knowledge`、`question_knowledge_rel`、`question_asset`、`question_review_record`、`question_change_log`、`question_import_batch`、`question_import_row`。
- 主键、唯一索引、常用查询索引存在。

### 6.2 数据导入

将 SQLite 全量数据导入 MySQL，至少覆盖以下数据域：

- 系统登录与权限：`user_account`、`role`、`user_role_id_list`、`sys_menu`、`role_menu`。
- 系统配置：`sys_code`、`sys_btn`、`btn_role_id_list`、`schema_page*`。
- 题库数据：全部 `question_*` 表。

通过标准：

- SQLite 与 MySQL 的核心表行数一致。
- 管理员账号、角色、菜单权限一致。
- 题库题目数量、选项数量、答案数量、年级/科目/单元/知识点数量一致。
- 随机抽样 10 道题，题干、选项、答案、解析、知识点关联一致。

## 7. 配置切换与启动

### 7.1 配置切换

修改 `/data/project/ai-study/backend/conf/application.properties`：

```properties
driverName=mysql
dataSourceName=<与 auto-check 同源凭据，仅替换目标库为 ai_study>
```

要求：

- 不提交明文密码。
- 如果需要落地本机运行配置，使用不进入版本控制的本地配置文件或部署环境注入。
- 保留切回 SQLite 的备份配置。

### 7.2 启动验证

```bash
cd /data/project/ai-study/backend
./shutdown.sh
ss -ltnp | rg ':8026' || true
./linux-start-dev.sh
ss -ltnp | rg ':8026'
curl --noproxy '*' -sS -m 5 -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8026/collect-ui/
```

通过标准：

- 8026 端口监听正常。
- `/collect-ui/` 返回 `200` 或可接受的前端入口重定向。
- `logs/server.log` 中没有数据库连接失败、SQL 语法错误、panic。

## 8. 接口验收

### 8.1 登录和菜单接口

必须验证：

- `system.login` 登录成功。
- `system.menu_query` 返回当前用户可见菜单。
- 菜单数量和层级与 SQLite 版本一致。
- 登录后 session 生效，刷新页面仍可访问菜单。

通过标准：

- 登录接口 `success=true`。
- 菜单包含：首页、系统管理/用户管理、角色管理、菜单管理、码表管理、题库管理、年级科目维护。
- 不出现 500、SQL 错误、空菜单。

### 8.2 功能页面接口

每个菜单至少验证初始化接口和主表查询接口：

| 页面 | 必验服务 |
| --- | --- |
| 首页 | `frontend.home` |
| 用户管理 | `frontend.user`、`hrm.user_list`、`hrm.role_query` |
| 角色管理 | `frontend.role`、`hrm.role_query`、`system.role_menu_query` |
| 菜单管理 | `frontend.menu_manage`、`system.menu_query`、`system.menu_save/menu_update/menu_delete` 的校验路径 |
| 码表管理 | `frontend.sys_code`、`system.sys_code_list`、`system.get_sys_code` |
| 题库管理 | `frontend.question_bank`、`question.question_query`、`question.question_stats`、`question.question_detail`、`question.question_choice_detail`、题目保存/编辑/删除校验路径 |
| 年级科目维护 | `frontend.question_taxonomy`、`question.grade_query`、`question.subject_query`、`question.unit_query`、`question.knowledge_query` |

通过标准：

- 所有页面初始化服务 HTTP 状态为 200。
- 所有主查询服务业务 `success=true`。
- 空参触发的业务校验失败可以接受，但必须归类为预期校验失败，不能是 SQL/连接/类型错误。
- 写入类服务至少执行一次新增、修改、删除或软删除闭环，并清理测试数据。

## 9. 页面验收

使用 Playwright 登录后逐个打开全部菜单。每页采集：

- 页面截图。
- `console.error`。
- `pageerror`。
- `requestfailed`。
- `/template_data/data?service=` HTTP 状态。
- 低代码接口返回的 `success/msg`。

逐页验收点：

| 页面 | 验收动作 | 通过标准 |
| --- | --- | --- |
| 首页 | 打开首页，确认卡片/表格可见 | 页面无白屏，无接口失败 |
| 用户管理 | 查询用户，打开新增/编辑弹窗，执行测试用户新增-修改-删除闭环 | 表格有数据，测试数据闭环成功 |
| 角色管理 | 查询角色，查看角色用户/菜单授权区域 | 角色列表可见，授权相关接口无 SQL 错误 |
| 菜单管理 | 展开菜单树/表格，查看菜单详情 | 菜单层级完整，6 个功能菜单和首页都存在 |
| 码表管理 | 查询码表，按关键字筛选，打开新增/编辑弹窗 | 码表列表和筛选正常 |
| 题库管理 | 查询题目，查看详情，执行一条测试题新增-编辑-删除闭环 | 题干、选项、答案、解析正常展示 |
| 年级科目维护 | 切换年级/科目/单元/知识点 tab，执行测试数据新增-修改-删除闭环 | 四类维护列表和树结构正常 |

通过标准：

- `console.error=0`。
- `pageerror=0`。
- 关键资源失败为 0。
- 接口 HTTP 失败为 0。
- 非预期业务失败为 0。
- 每个页面截图落盘。

## 10. 数据一致性验收

迁移后执行 SQLite 与 MySQL 对账：

| 对账项 | 通过标准 |
| --- | --- |
| 表清单 | MySQL 不缺核心业务表 |
| 行数 | 核心表行数与 SQLite 一致 |
| 管理员 | `admin` 可登录，角色为管理员 |
| 菜单 | 可见菜单层级一致 |
| 码表 | 常用码表类型和值一致 |
| 题库 | 题目、选项、答案、填空答案、知识点关联数量一致 |
| 富文本/JSON | 题干、解析、答案 JSON、导入批次 JSON 未乱码、未截断 |
| 时间字段 | 创建/修改时间可查询、排序正常 |

## 11. 回滚验收

回滚步骤：

1. 恢复 `backend/conf/application.properties` 为 SQLite 配置。
2. 恢复迁移前 SQLite 备份文件。
3. 重启 AI Study 后台。
4. 重新验证登录和全部菜单页面。

通过标准：

- SQLite 配置下后台恢复可用。
- 登录、菜单、题库查询正常。
- 回滚过程中不破坏 MySQL 备份数据。

## 12. 最终通过门槛

本次迁移验收只有在以下条件全部满足时通过：

1. MySQL 目标库已重新创建，并使用与 auto-check 一致的 MySQL 实例和账号体系。
2. SQLite 到 MySQL 的核心表行数和抽样数据一致。
3. 后端 `go test ./...` 通过。
4. AI Study 后台 8026 启动正常，日志无数据库连接/SQL/panic 错误。
5. 登录页、首页和 6 个功能菜单全部验证通过。
6. 页面级浏览器验收中 `console.error/pageerror/requestfailed/接口 HTTP 失败/非预期业务失败` 均为 0。
7. 写入类页面完成测试数据新增、修改、删除或软删除闭环，并完成测试数据清理。
8. 回滚路径已演练或至少完成可执行性检查。

## 13. 验收产物

每轮验收需落盘：

- SQLite 迁移前备份。
- 如目标 MySQL 库原先存在，对应 MySQL dump 备份。
- 数据库重建和导入日志。
- 表行数对账报告。
- 接口巡检 JSON/Markdown 报告。
- Playwright 页面验收 JSON 报告。
- 7 个页面截图。
- 失败项清单和修复记录。
