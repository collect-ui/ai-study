# AI Study Admin Backend

AI Study 后台基于 collect-ui 低代码服务架构，统一入口为：

```text
POST /template_data/data?service=<service_name>
```

本项目当前保留登录、用户、角色、菜单、码表、按钮权限和页面配置基础能力。运行日志写入 `logs/server.log`，不会记录请求体，避免把登录密码写入日志。

## 本地启动

```bash
go run .
```

默认端口来自 `conf/application.properties`：

```text
server_port=8026
```

管理端地址：

```text
http://127.0.0.1:8026/collect-ui
```

默认本地账号：

```text
admin / 123456
```

## 初始化数据

```bash
sqlite3 database/ai_study_admin.db ".read scripts/init-ai-study-admin.sql"
```

初始化内容：

- 管理员用户 `admin`
- 管理员角色 `admin`
- 用户、角色、菜单、码表菜单
- `ai-study-admin` 项目菜单授权

## 验证

```bash
go test ./...
go build ./...
```

登录接口：

```bash
curl -i -c /tmp/ai-study-admin.cookie \
  -H 'Content-Type: application/json' \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.login' \
  -d '{"username":"admin","password":"123456"}'
```

当前用户接口：

```bash
curl -i -b /tmp/ai-study-admin.cookie \
  -X POST 'http://127.0.0.1:8026/template_data/data?service=system.current_user' \
  -d '{}'
```
