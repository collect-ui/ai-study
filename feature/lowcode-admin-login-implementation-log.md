# AI Study 小程序后台低代码迁入实现日志

日期：2026-05-18

## 1. 实现范围

本次按 `feature/lowcode-admin-login-design.md` 完成第一阶段后台实现：

- 后台目录：`backend/`
- 运行端口：`8026`
- 数据库：`backend/database/ai_study_admin.db`
- 初始化脚本：`backend/scripts/init-ai-study-admin.sql`
- 运行日志：`backend/logs/server.log`

保留能力：

- `system.login`
- `system.current_user`
- `system.logout`
- `hrm.user_list`
- `hrm.create_user`
- `hrm.update_user_by_user_id`
- `hrm.update_user_by_username`
- `hrm.role_query`
- `hrm.create_role`
- `hrm.edit_role`
- `hrm.user_role_add`
- `hrm.user_role_replace`
- `system.menu_query`
- `system.role_menu_query`
- `system.role_menu_save_bulk_by_menu`
- `system.sys_code_*`
- `system.sys_btn_*`
- `system.schema_page_*`
- `frontend.*` 基础页面数据

## 2. 研发记录

1. 后台入口增加运行日志：
   - `backend/main.go`
   - 输出到控制台和 `backend/logs/server.log`
   - 记录服务名、HTTP 方法、状态码、耗时、客户端 IP
   - 不记录请求体，避免密码进入日志

2. 增加 query service 兼容：
   - 支持 `POST /template_data/data?service=system.login`
   - 自动把 query 中的 `service` 合并进 JSON body
   - 兼容低代码前端 `api: post:/template_data/data?service=...`

3. 调整登录链路：
   - `system.login` 调用 `hrm.user_list` 时传 `with_password: true`
   - 登录校验用户名、密码、用户状态
   - 成功后写入 `username`、`nick`、`userid`、`user_id`

4. 调整用户列表安全：
   - `hrm.user_list` 默认不返回密码哈希
   - 普通用户列表中 `password` 返回空字符串
   - 仅登录链路通过 `with_password=true` 读取密码哈希

5. 限制用户编辑字段：
   - `hrm.update_user_by_user_id` 增加 `update_fields`
   - 避免编辑用户资料时误更新密码字段

6. 初始化 AI Study 后台基础数据：
   - 默认用户：`admin / 123456`
   - 默认角色：`admin`
   - 默认菜单：首页、系统管理、用户管理、角色管理、菜单管理、码表管理
   - 所有菜单归属：`ai-study-admin`

7. 清理旧业务残留：
   - 删除未注册的 LDAP 服务目录
   - 删除旧旅行社、转写、检测、MySQL 探测脚本
   - 删除旧二进制、PDF、DOCX、图标、历史 pid 和测试残留
   - 删除旧敏感示例文档内容

## 3. 测试记录

### Go 编译测试

```bash
cd /data/project/ai-study/backend
go test ./...
go build ./...
go vet ./...
```

结果：全部通过。

### JSON 静态检查

```bash
cd /data/project/ai-study
find backend/collect/frontend/page_data/data -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool
find miniprogram -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool
find miniprogram -name '*.js' -print0 | xargs -0 -n1 node --check
```

结果：全部通过。

### 数据初始化验证

```bash
sqlite3 backend/database/ai_study_admin.db ".read backend/scripts/init-ai-study-admin.sql"
sqlite3 backend/database/ai_study_admin.db "select count(*) from role; select count(*) from user_account; select count(*) from role_menu;"
```

结果：

```text
1
1
7
```

### 登录成功

```bash
curl --noproxy '*' -s -c /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"123456"}'
```

结果：

```json
{"count":0,"success":true,"code":"0","msg":"成功","data":{"nick":"管理员","user_id":"admin-user","userid":"admin-user","username":"admin"}}
```

### 密码错误

```bash
curl --noproxy '*' -s \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"wrong-password"}'
```

结果：返回 `success=false`，错误信息为密码错误。

### 当前用户

```bash
curl --noproxy '*' -s -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.current_user' \
  -d '{}'
```

结果：返回当前管理员用户。

### 用户列表

```bash
curl --noproxy '*' -s -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=hrm.user_list' \
  -d '{"count":false}'
```

结果：返回管理员用户，`password` 字段为空字符串。

### 角色和菜单

```bash
curl --noproxy '*' -s -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=hrm.role_query' \
  -d '{"count":false}'

curl --noproxy '*' -s -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.menu_query' \
  -d '{"to_tree":true,"with_role":true}'
```

结果：返回 `admin` 角色和 AI Study 后台基础菜单。

### 退出登录

```bash
curl --noproxy '*' -s -c /tmp/ai-study-admin.cookie -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.logout' \
  -d '{}'

curl --noproxy '*' -s -i -b /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.current_user' \
  -d '{}'
```

结果：退出后 `system.current_user` 返回 `请登录！！！`。

### 静态管理端

```bash
curl --noproxy '*' -s -I 'http://127.0.0.1:8026/collect-ui'
curl --noproxy '*' -s -I -X GET 'http://127.0.0.1:8026/'
```

结果：

- `/collect-ui` 返回 `200 OK`
- `/` 返回 `301 Moved Permanently`，跳转到 `/collect-ui`

## 4. 安全验证

```bash
rg -n "123456|wrong-password|\"password\"|密码错误" backend/logs/server.log
rg -n "log:\s*true" backend/collect
```

结果：

- `backend/logs/server.log` 未写入明文密码、请求体或密码错误详情
- `backend/collect` 中无 `log: true`

```bash
rg -n "ldap|DB_PASSWORD_PATTERN|sport_level|wechat|travel|tencent|旅行|腾讯|游客|auto-check|webshell|autodesk|jira|detect_escape|sales_cs|customer_assign|DB_DSN_PATTERN|corpsecret|deepseek|coze|BEARER_TOKEN_PATTERN" backend/collect backend/conf backend/main.go backend/model backend/plugins backend/frontend/data backend/collect/frontend/page_data/data backend/scripts backend/AGENTS.md backend/readme.md
```

结果：无旧业务和敏感配置残留。

## 5. 小程序影响验证

本次未修改 `miniprogram`。

```bash
find miniprogram -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.gif' -o -iname '*.webp' -o -iname '*.svg' -o -iname '*.mp3' -o -iname '*.mp4' -o -iname '*.pdf' -o -iname '*.doc' -o -iname '*.docx' -o -iname '*.xls' -o -iname '*.xlsx' -o -iname '*.zip' \) -print
```

结果：无输出，`miniprogram` 内没有附件类文件。

## 6. 当前运行状态

后台已启动：

```text
http://127.0.0.1:8026/collect-ui
```

运行日志：

```text
backend/logs/server.log
```
