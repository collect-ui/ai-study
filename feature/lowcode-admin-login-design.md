# AI Study 小程序后台低代码迁入与登录接入设计

## 1. 背景

本项目需要为 AI Study 小程序重新立一个后台管理项目。后台架构沿用既有低代码项目 `/data/project/auto-check`，用户输入中提到的 `/data/project/auto_check` 当前本机未找到，实际参考路径为 `/data/project/auto-check`。

当前 `/data/project/ai-study` 是小程序项目，后台独立放在 `/data/project/ai-study/backend`，不把图片、音视频、PDF、Office、zip 等附件放入 `miniprogram`。后台静态管理端资源继续放在 `backend/frontend`，小程序远程资源仍按 `miniprogram/config/assets.js` 和远程资源规则处理。

## 2. 目标

1. 从 `/data/project/auto-check` 整体拷贝低代码后台能力到本项目，再按 AI Study 后台裁剪。
2. 新后台项目保留用户、角色、菜单和低代码页面配置能力。
3. 第一阶段先打通管理后台登录链路：账号密码登录、会话写入、当前用户读取、退出登录。
4. 后续业务模块基于保留的低代码配置继续扩展，不继承 auto-check 的旅行社、运动、文档、WebShell 等业务能力。

## 3. 参考架构理解

auto-check 后台是 Go + Gin + collect-ui 低代码服务：

- 入口：`main.go`
- 配置：`conf/application.properties`
- 低代码总路由：`collect/service_router.yml`
- 统一接口：`POST /template_data/data?service=<service_name>`
- 静态管理端：`/collect-ui` 对应 `frontend/collect-ui`
- 静态数据资源：`/data` 对应 `frontend/data`
- WebSocket：`GET /template_data/ws/:token`
- 会话：`gin-contrib/sessions`，cookie 名为 `session_id`

低代码配置的核心结构：

- `collect/service_router.yml` 注册一级服务域和模块处理器。
- `collect/system/service.yml` 注册系统服务，如登录、菜单、码表、页面配置。
- `collect/hrm/service.yml` 注册用户、角色、用户角色关系等服务。
- 每个服务目录下的 `index.yml` 定义参数、校验、SQL 文件、模型保存/更新/删除、结果处理器。
- 模型定义在 `model/base`，由 `templateService.SetDatabaseModel(&model.TableData{})` 注入。

## 4. 新项目后台边界

新后台目录为：

```text
/data/project/ai-study/backend
```

第一阶段保留：

- `main.go`
- `go.mod`、`go.sum`
- `conf/application.properties`
- `collect/service_router.yml`
- `collect/system/login`
- `collect/system/menu`
- `collect/system/role_menu`
- `collect/system/sys_code`
- `collect/system/sys_btn`
- `collect/system/schema_page`
- `collect/system/schema_page_field`
- `collect/system/schema_page_data`
- `collect/hrm/user`
- `collect/hrm/role`
- `collect/hrm/user_role`
- `collect/hrm/user_flow`
- `collect/frontend`
- `frontend/collect-ui`
- `frontend/data`
- `model/base` 中用户、角色、菜单、码表、页面配置需要的模型
- 必需的低代码外部处理器，优先只保留页面 schema 转换等后台必须能力

第一阶段删除或不注册：

- `travel_*`
- `detect_escape`
- `tencent_key`
- `sales_cs_user`
- `customer_assign`
- `sport`
- `autodesk`
- `doc`
- `webshell`
- `jira`
- `wechat`
- `work_task`
- SSH、SFTP、Shell、文档生成、外部同步等非基础后台插件
- auto-check 原项目中的历史报告、测试产物、打包产物和业务 SQL

## 5. 配置设计

`backend/conf/application.properties` 使用 AI Study 独立配置：

```properties
collect_file_path=./collect/service_router.yml
system_model=my
driverName=sqlite3
dataSourceName=./database/ai_study_admin.db
otherDataSource=
must_login=true
user_id_key=user_id
current_project_code=ai-study-admin
server_port=8026
is_https=false
domain=collect-ui.top
dirList=/data,./frontend/data,false;/collect-ui,./frontend/collect-ui,true
project=ai-study-admin
basic_auth_service=system.login
```

要求：

- 不复制 auto-check 的生产数据库密码、企业微信密钥、聊天服务 token 等敏感配置。
- 本地优先使用 SQLite，文件为 `backend/database/ai_study_admin.db`。
- `must_login=true`，只有 `system.login`、必要的公开健康检查接口允许 `must_login: false`。
- `project/current_project_code` 统一为 `ai-study-admin`，用于菜单归属和页面配置隔离。

## 6. 登录链路设计

### 6.1 登录接口

低代码服务：

```text
system.login
```

HTTP 调用：

```bash
curl -i -c /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"123456"}'
```

处理流程：

1. 校验 `username` 和 `password` 非空。
2. 对入参密码执行 `md5` 模板处理。
3. 调用 `hrm.user_list`，按 `username` 查询本地 `user_account`。
4. 校验用户存在、密码哈希一致、账号未软删除、状态可用。
5. 写入 session 字段：`username`、`nick`、`userid`、`user_id`。
6. 调用 `system.current_user` 并返回当前用户。

### 6.2 当前用户接口

低代码服务：

```text
system.current_user
```

HTTP 调用：

```bash
curl -i -b /tmp/ai-study-admin.cookie \
  'http://127.0.0.1:8026/template_data/data?service=system.current_user'
```

返回字段：

- `userid`
- `user_id`
- `username`
- `nick`

### 6.3 退出接口

低代码服务：

```text
system.logout
```

处理流程：

1. 删除 session 中的 `username`、`nick`、`userid`、`user_id`。
2. 前端清理登录态并跳转登录页。

### 6.4 默认初始化数据

本地数据库需要最小种子数据：

- `user_account`：默认管理员 `admin`，默认密码仅用于本地初始化，落库保存 MD5。
- `role`：默认管理员角色。
- `user_role_id_list`：管理员用户和管理员角色关系。
- `sys_menu`：用户管理、角色管理、菜单管理、页面配置。
- `role_menu`：管理员角色绑定全部基础菜单。
- `sys_code`：账号状态、性别等用户列表依赖码表。

生产环境必须在首次登录后修改默认密码。

## 7. 菜单与权限设计

保留三层权限关系：

```text
user_account -> user_role_id_list -> role -> role_menu -> sys_menu
```

菜单服务：

- `system.menu_query`：查询菜单，可按 `belong_project=ai-study-admin` 隔离。
- `system.menu_save`：新增菜单并可同步角色菜单关系。
- `system.menu_update`：更新菜单并对比角色授权变化。
- `system.role_menu_query`：查询角色拥有菜单。
- `system.role_menu_save_bulk_by_menu`：保存角色菜单授权。

用户角色服务：

- `hrm.user_list`
- `hrm.create_user`
- `hrm.update_user_by_user_id`
- `hrm.delete_user_by_user_id_list`
- `hrm.role_query`
- `hrm.create_role`
- `hrm.edit_role`
- `hrm.user_role_add`
- `hrm.user_role_replace`

要求：

- 菜单查询只返回当前项目菜单。
- 用户编辑时可以替换角色。
- 删除用户使用软删除 `is_delete=1`。
- 登录态只放最小用户信息，角色和菜单由接口查询，不直接塞入 session。

## 8. 研发计划

### 阶段 1：迁入与裁剪

1. 从 `/data/project/auto-check` 整体复制后台到 `/data/project/ai-study/backend`。
2. 清理 auto-check 原有 `.git`、报告、打包产物、业务测试产物。
3. 裁剪 `collect/service_router.yml`，只注册 `system`、`hrm`、`frontend`。
4. 裁剪 `collect/system/service.yml` 和 `collect/hrm/service.yml`，只保留基础后台服务。
5. 清理 `model/register.go`，只注册基础模型，去掉旅行社、腾讯 Key、聊天记录等业务 schema 初始化。
6. 清理 `plugins` 注册，只保留低代码基础功能需要的处理器。
7. 修改 `conf/application.properties` 为 AI Study 独立配置。

### 阶段 2：登录接入

1. 确认 `system.login`、`system.current_user`、`system.logout` 可用。
2. 确认 `hrm.user_list` SQL 可按用户名查询，并过滤 `is_delete=0`、不可用状态。
3. 初始化 SQLite 基础表和默认管理员。
4. 管理端登录页调用 `system.login`。
5. 登录成功后进入后台首页，拉取 `system.current_user` 和菜单。
6. 未登录访问受保护服务时返回统一未登录错误。

### 阶段 3：用户、角色、菜单管理

1. 用户管理支持新增、编辑、软删除、重置密码、角色替换。
2. 角色管理支持新增、编辑、删除。
3. 菜单管理支持新增、编辑、排序、授权角色。
4. 页面配置继续使用低代码 schema 表，后续业务页面基于它扩展。

### 阶段 4：小程序业务后台扩展

1. 在保留后台基础上新增 AI Study 业务模型。
2. 新增业务服务目录时只注册 AI Study 相关模块。
3. 小程序端如需调用后台登录态，使用合法 HTTPS 域名并在请求中维护 cookie。
4. 小程序资源继续走远程资产规则，不放入 `miniprogram`。

## 9. 测试计划

### 9.1 静态检查

```bash
cd /data/project/ai-study/backend
go test ./...
go build ./...
```

```bash
cd /data/project/ai-study
find backend/collect -name '*.yml' -print
find backend/collect -name '*.sql' -print
find miniprogram -name '*.js' -print0 | xargs -0 -n1 node --check
find miniprogram -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool
```

### 9.2 登录接口测试

启动后台：

```bash
cd /data/project/ai-study/backend
go run .
```

登录成功：

```bash
curl -s -i -c /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"123456"}'
```

读取当前用户：

```bash
curl -s -i -b /tmp/ai-study-admin.cookie \
  'http://127.0.0.1:8026/template_data/data?service=system.current_user'
```

密码错误：

```bash
curl -s -i \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"wrong-password"}'
```

退出后再读当前用户：

```bash
curl -s -i -b /tmp/ai-study-admin.cookie \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.logout'

curl -s -i -b /tmp/ai-study-admin.cookie \
  'http://127.0.0.1:8026/template_data/data?service=system.current_user'
```

### 9.3 用户角色菜单测试

1. `hrm.user_list` 能返回管理员用户和 `role_id_list`。
2. `hrm.role_query` 能返回管理员角色。
3. `system.menu_query` 能返回 `ai-study-admin` 项目菜单。
4. `system.role_menu_query` 能返回管理员角色菜单授权。
5. 新增用户后调用 `hrm.user_role_replace`，角色关系能被替换。

### 9.4 管理端 UI 测试

1. 打开 `http://127.0.0.1:8026/collect-ui`。
2. 未登录时进入登录页。
3. 输入管理员账号密码可登录。
4. 登录后显示用户、角色、菜单相关导航。
5. 刷新页面后 session 仍有效。
6. 退出登录后回到登录页。

### 9.5 小程序影响测试

本阶段原则上不改 `miniprogram`。如果后续接入小程序端后台接口，需要执行项目 QA workflow：

```bash
node --check miniprogram/pages/home/home.js
python3 -m json.tool miniprogram/app.json
python3 -m json.tool miniprogram/pages/home/home.json
```

涉及 UI 的改动还需要启动微信开发者工具、运行 automator、截图并生成预览二维码。

## 10. 验证清单

必须通过：

- `backend` 能启动在 `8026`。
- `POST /template_data/data?service=system.login` 正确写入 cookie。
- `system.current_user` 能通过 cookie 返回当前用户。
- 错误密码不能登录。
- `system.logout` 后会话不可继续使用。
- 未登录访问受保护服务时被拦截。
- `hrm.user_list`、`hrm.role_query`、`system.menu_query`、`system.role_menu_query` 可用。
- `collect/service_router.yml` 不再注册 auto-check 业务服务。
- `conf/application.properties` 不包含生产数据库密码、企业微信密钥、聊天 token 等敏感信息。
- `miniprogram` 内没有新增附件类文件。

可用命令：

```bash
cd /data/project/ai-study
find miniprogram -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.gif' -o -iname '*.webp' -o -iname '*.svg' -o -iname '*.mp3' -o -iname '*.mp4' -o -iname '*.pdf' -o -iname '*.doc' -o -iname '*.docx' -o -iname '*.xls' -o -iname '*.xlsx' -o -iname '*.zip' \) -print
```

## 11. 验收标准

1. 后台代码位于 `/data/project/ai-study/backend`，能独立启动。
2. 登录、当前用户、退出登录三个服务可用。
3. 用户、角色、菜单基础服务可用。
4. auto-check 无关业务模块不再对外注册。
5. 配置中没有旧项目敏感信息。
6. 设计、研发、测试、验证记录保存在 `feature/`。

## 12. 风险与处理

- 低代码前端可能仍残留 auto-check 品牌文案：在阶段 1 裁剪后统一搜索 `auto-check`、旅行社、运动、腾讯 Key 等关键词并替换或删除。
- SQLite 和 MySQL 的 SQL 函数差异：基础 SQL 优先使用兼容写法；如必须支持 MySQL，后续补充双数据库验证。
- 默认密码风险：默认管理员只用于初始化，生产发布前必须改密。
- 旧插件残留风险：只在 `plugins.GetRegisterList()` 注册确实需要的处理器，避免暴露 SSH/Shell/SFTP 等高风险能力。
- 小程序 cookie 维护和合法域名限制：真机环境必须确认微信公众平台后台已配置 `collect-ui.top` 合法域名。
