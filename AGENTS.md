# AGENTS.md

## AI Study 远程资源规则

项目根目录：`/data/project/ai-study`

处理小程序 UI、配置或资源时，先读取项目 QA skill：

```bash
sed -n '1,220p' /data/project/ai-study/.codex/skills/ai-study-miniprogram-qa/SKILL.md
```

处理远程服务器同步时，读取实际存在的远程同步 skill：

```bash
sed -n '1,240p' /data/project/auto-check/.codex/skills/autocheck-server-sync/SKILL.md
```

不要在文档、代码或日志里写入 SSH 密码。使用环境变量 `AUTO_CHECK_SERVER_PASSWORD`，或沿用环境里已有的 `TEST_SERVER_PASSWORD`。

## 管理后台前端开发与打包

AI Study 管理后台由三个目录协同：

- `/data/project/collect-ui`：低代码组件库。新增或修改 schema tag、组件样式、动作时在这里改；如果 AI Study 后台要用新 tag，确认它已包含在 `src/index.min.tsx` 的 minimal 白名单中。
- `/data/project/sport-ui`：后台前端静态壳。AI Study 使用专用入口，不复用 auto-check 入口：
  - `index-ai-study.html`
  - `src/main.ai-study.tsx`
  - `data/ai-study/app.json`
  - `vite.config.ai-study.min.js`
- `/data/project/ai-study`：本项目后端与页面 schema。页面配置在 `backend/collect/frontend/page_data/data`，静态部署目标是 `backend/frontend/collect-ui`。

本地联调：

```bash
cd /data/project/ai-study/backend
go run .
```

```bash
cd /data/project/sport-ui
npm ls --depth=0 collect-ui
ls -ld node_modules/collect-ui && readlink -f node_modules/collect-ui
AI_STUDY_BACKEND_URL=http://127.0.0.1:8026 npm run dev:ai-study
```

构建 AI Study 后台静态壳：

```bash
cd /data/project/collect-ui
COLLECT_UI_ENTRY=./src/index.min.tsx npm run build

cd /data/project/sport-ui
npm run build:ai-study:min
```

构建并部署到本项目：

```bash
cd /data/project/sport-ui
bash scripts/deploy_ai_study_collect_ui.sh
```

`scripts/deploy_ai_study_collect_ui.sh` 默认会：

- 使用 `/data/project/collect-ui` 构建 minimal `collect-ui`。
- 使用 `vite.config.ai-study.min.js` 输出 `build-ai-study`。
- 覆盖部署到 `/data/project/ai-study/backend/frontend/collect-ui`。
- 将 `index-ai-study.html` 复制为线上入口 `index.html`。
- 访问 `http://127.0.0.1:8026` 做登录页与 AI Study 菜单路由的 Playwright smoke test。

如果只想打包不验证：

```bash
VERIFY_AFTER_DEPLOY=0 bash scripts/deploy_ai_study_collect_ui.sh
```

不要手工编辑 `backend/frontend/collect-ui/assets` 内的 hash 产物；需要变更时从 `collect-ui` 或 `sport-ui` 重新构建。AI Study 专用入口变更只改 `ai-study` 命名文件，不要修改 `index-autocheck.html`、`src/main.autocheck.tsx` 或 `vite.config.autocheck.min.js`。

## 后端工程开发规范

本节吸收自 `/data/project/sport/AGENTS.md`，并按 AI Study 当前目录、端口和脚本名调整。

### 项目快照

- 后端目录：`/data/project/ai-study/backend`
- 语言：Go，HTTP 框架：Gin
- 架构：低代码后端，主要由 YAML/JSON 配置、SQL 文件和 Go 插件驱动
- 入口：`backend/main.go`
- 关键目录：
  - `backend/model/`：模型与表注册
  - `backend/plugins/`：自定义低代码/plugin handler
  - `backend/collect/`：服务定义、SQL、页面元数据
  - `backend/conf/`：运行配置
  - `backend/frontend/`：已部署静态前端资产
  - `backend/database/`：本地运行数据，默认不作为源码修改对象

### 本地运行

优先使用项目脚本启停后台，不要随手启动长期后台 `go run main.go &` 进程。

```bash
cd /data/project/ai-study/backend
./linux-start-dev.sh
```

停止服务：

```bash
cd /data/project/ai-study/backend
./shutdown.sh
```

兼容脚本：

- `./linux-start-dev.sh`：本地启动，日志写入 `run-dev.log`，PID 写入 `run-dev.pid`
- `./shutdown.sh`：按端口停止服务
- `./linux-shutdown.sh`：按 PID 文件停止服务

重要运行行为：

- 本项目不默认做 Go 代码热重载。
- 修改 Go 代码或 `backend/collect/`、`backend/conf/` 下的运行配置后，必须重启服务。
- 推荐重启流程：

```bash
cd /data/project/ai-study/backend
./shutdown.sh
./linux-start-dev.sh
```

必做启动校验，端口为 `8026`：

```bash
cd /data/project/ai-study/backend
./shutdown.sh
ss -ltnp | rg ':8026' || true
./linux-start-dev.sh
ss -ltnp | rg ':8026'
curl --noproxy '*' -sS -m 5 -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8026/collect-ui/
```

注意：环境可能设置了 `http_proxy/https_proxy`，检查本地 `127.0.0.1:8026` 时使用 `--noproxy '*'`，避免误判。

### 构建与测试

后端构建：

```bash
cd /data/project/ai-study/backend
go build ./...
```

后端测试：

```bash
cd /data/project/ai-study/backend
go test ./...
```

格式化与静态检查：

```bash
cd /data/project/ai-study/backend
go fmt ./...
go vet ./...
```

如果本地安装了可选工具，也可以执行：

```bash
staticcheck ./...
golangci-lint run
```

典型 Go 代码变更推荐验证顺序：

1. `go fmt ./...`
2. `go test ./...`
3. `go vet ./...`

仅修改 `backend/collect/` 或 `backend/conf/` 配置时，至少执行：

1. `go test ./...`
2. 重启后台并访问真实页面或接口验证运行时装配

### 代码风格

- Go import 使用标准分组：标准库、第三方、本项目。
- 保留邻近代码已有别名与命名风格。
- Go 文件必须使用 `gofmt`，不要手工排版。
- 文件默认保持 ASCII；已有中文业务标签、注释或配置文案的文件可以继续使用中文。
- 包名使用小写；导出标识符使用 PascalCase，非导出标识符使用 camelCase。
- 多词文件名沿用现有 snake_case 风格，尤其是 plugin/handler 相关文件。
- 添加模型时同步更新对应注册文件；添加插件时同步更新 `plugins/a_register.go`。
- 处理错误时及时检查并返回上下文信息；请求路径代码不要新增不必要的 panic。
- 控制流优先早返回，保持 handler/plugin 的 `Result` 逻辑线性易读。
- 注释保持短而有用；新增导出函数或类型时补充必要说明。

### 低代码配置

- 许多行为由 `backend/collect/` 下 YAML/JSON/SQL 驱动，优先改配置，不要先写硬编码分支。
- 修改页面、服务、store、form、action、fragment 时，保持原有 key、缩进风格和 schema 形状。
- 重命名 service key、字段或页面路由前，先跨 `backend/collect/`、`backend/conf/` 和插件代码查引用。
- 如果低代码配置已能表达行为，不要新增 Go/TypeScript 分支；确实需要代码时，范围要小，并能说明配置为何不够。
- 除 model 定义、表注册等模型层变更外，不要私自新增或修改 Go 业务代码。确实需要写 Go 时，必须先向用户说明要写什么、为什么必须写、为什么不能用低代码配置或已有服务替代，并获得确认后再实施。

### 数据库与生成产物

- `backend/database/` 是本地运行数据，默认不要作为源码修改提交。
- 不要提交本地数据库、IDE 文件、二进制、压缩包、临时日志或测试产物，除非任务明确要求。
- 不要手工编辑 `backend/frontend/collect-ui/assets` 内 hash 产物；需要变更时从 `collect-ui` 或 `sport-ui` 重新构建部署。
- 生成文件后，确认它是预期源码或部署产物，避免把环境输出混入代码变更。

### 工作约定

- 改一个子系统前先读邻近文件，沿用现有模式。
- 低代码模块优先通过配置扩展行为。
- 行为依赖运行时装配时，修改后要重启后台并做真实 URL/hash 路由验证。
- 浏览器验收要收集 `console.error`、`pageerror`、关键资源失败和接口失败。
- 如果构建或运行异常，先确认相关兄弟目录是否同步：
  - `/data/project/collect-ui`
  - `/data/project/sport-ui`

## 图片和附件长期策略

- 不要把图片、音视频、PDF、Office、zip 等附件放入 `miniprogram`。
- 小程序资源 URL 统一从 `miniprogram/config/assets.js` 生成。
- 远程资源根目录：`/data/file/ai-study/assets`
- 线上 URL 前缀：`https://collect-ui.top/ai-study/assets`
- 本地执行文档：`feature/remote-assets-migration.md`
- 本地 manifest：`feature/remote-assets/asset-manifest.json`
- 远程 manifest：`/data/file/ai-study/asset-manifest.sha256`

新增资源同步流程：

```bash
mkdir -p feature/remote-assets/staging
# 将资源按远程相对路径放入 feature/remote-assets/staging
bash scripts/sync-remote-assets.sh
```

如果资源在别的目录：

```bash
AI_STUDY_LOCAL_ASSET_DIR=/path/to/assets bash scripts/sync-remote-assets.sh
```

远程 HTTPS 服务部署或修复：

```bash
bash scripts/deploy-remote-file-service.sh
```

受信任 HTTPS 证书签发或续签：

```bash
bash scripts/issue-trusted-cert.sh
```

强制重新签发 Let's Encrypt 证书：

```bash
AI_STUDY_ACME_FORCE=1 bash scripts/issue-trusted-cert.sh
```

当前证书由 Let's Encrypt 签发，acme.sh 在远程 root crontab 中自动续期。不要用自签证书作为真机预览或正式环境证书；`scripts/deploy-remote-file-service.sh` 只在证书文件缺失时生成自签兜底。

## 验证要求

资源迁移或新增后至少执行：

```bash
bash -n scripts/sync-remote-assets.sh
bash -n scripts/deploy-remote-file-service.sh
find miniprogram -name '*.js' -print0 | xargs -0 -n1 node --check
find miniprogram -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool
curl --noproxy '*' -k -I https://collect-ui.top/ai-study/assets/prototype/logo.png
```

确认 `miniprogram` 内没有附件类文件：

```bash
find miniprogram -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.gif' -o -iname '*.webp' -o -iname '*.svg' -o -iname '*.mp3' -o -iname '*.mp4' -o -iname '*.pdf' -o -iname '*.doc' -o -iname '*.docx' -o -iname '*.xls' -o -iname '*.xlsx' -o -iname '*.zip' \) -print
```

生产发布前确认微信公众平台后台已配置 `collect-ui.top` 合法域名。
