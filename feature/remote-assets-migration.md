# 小程序远程图片与附件迁移设计执行文档

更新时间：2026-05-17 12:20 CST

## 1. 当前状态

- 远程主机：`202.140.140.117`
- 域名：`collect-ui.top`，当前解析到 `202.140.140.117`
- 远程根目录：`/data/file`
- 远程资源目录：`/data/file/ai-study/assets`
- 远程资源 URL 前缀：`https://collect-ui.top/ai-study/assets`
- 远程 manifest：`/data/file/ai-study/asset-manifest.sha256`
- 本地 manifest：`feature/remote-assets/asset-manifest.json`
- 远程 HTTPS 服务：`ai-study-file.service`
- 远程服务脚本：`/data/file/bin/remote-file-server.py`
- 远程证书：`/data/file/certs/collect-ui.top.crt`
- 证书状态：2026-05-17 已切换为 Let's Encrypt 受信任证书，过期时间为 2026-08-15 11:16:50 CST
- ACME 续期：远程 root crontab 已安装 `/root/.acme.sh/acme.sh --cron`，当前建议续期窗口从 2026-07-15 开始
- 资源缓存版本：`20260517-ca1`
- 小程序包内状态：`miniprogram/assets` 已删除，`miniprogram` 内当前没有图片、音视频、Office、PDF、zip 等附件类文件

注意：微信小程序正式版仍需要在微信公众平台后台配置 `https://collect-ui.top` 为合法下载/资源域名。

## 2. 目标

1. 小程序包只保留源码、配置和逻辑文件，图片和附件长期放在远程 `/data/file`。
2. 小程序源码统一通过 `miniprogram/config/assets.js` 引用远程资源 URL。
3. 远程静态资源服务可重复部署，支持 HTTPS 443。
4. 后续新增图片或附件时，通过脚本同步到远程并更新 manifest，避免重新把资源放回小程序包。
5. 文档、脚本、日志能让下一次接手时直接继续执行和验收。

## 3. 方案

### 3.1 资源路径规范

小程序内使用：

```js
const { ASSETS, assetUrl } = require("../../config/assets");
```

基础地址固定在：

```text
https://collect-ui.top/ai-study/assets
```

新增资源建议保持远程相对路径稳定，例如：

```text
prototype/logo.png
course/unit-1/cover.webp
attachments/demo.pdf
```

对应访问地址：

```text
https://collect-ui.top/ai-study/assets/prototype/logo.png
https://collect-ui.top/ai-study/assets/course/unit-1/cover.webp
https://collect-ui.top/ai-study/assets/attachments/demo.pdf
```

### 3.2 远程服务

远程 80 端口已被 `autocheck.service` 占用，本次没有改动它。图片服务使用独立 systemd 服务：

```text
ai-study-file.service -> python3 /data/file/bin/remote-file-server.py
listen: 0.0.0.0:443
root: /data/file
```

服务特点：

- 支持 `GET`、`HEAD`、`OPTIONS`
- 对 `/ai-study/assets/` 增加长期缓存头
- 增加基础 CORS 头，便于后续附件下载和调试
- 禁用目录列表，并屏蔽 `/bin`、`/certs`、隐藏文件路径
- 将 TLS 握手放到工作线程处理，listen backlog 调整为 128，避免慢连接占满 443
- 使用 `/data/file/certs/collect-ui.top.crt` 和 `.key`

### 3.3 本地脚本

部署或重装远程 HTTPS 服务：

```bash
bash scripts/deploy-remote-file-service.sh
```

签发或续签受信任 Let's Encrypt 证书：

```bash
bash scripts/issue-trusted-cert.sh
```

强制重新签发：

```bash
AI_STUDY_ACME_FORCE=1 bash scripts/issue-trusted-cert.sh
```

`scripts/deploy-remote-file-service.sh` 只在证书文件缺失时生成自签兜底。真机预览和正式环境必须使用受信任 CA 证书。

同步新增资源：

```bash
mkdir -p feature/remote-assets/staging
# 把要同步的图片或附件按远程相对目录放入 staging
bash scripts/sync-remote-assets.sh
```

如果资源源目录不在默认 staging：

```bash
AI_STUDY_LOCAL_ASSET_DIR=/path/to/assets bash scripts/sync-remote-assets.sh
```

SSH 密码来源遵循同步 skill，不写入仓库：

```bash
AUTO_CHECK_SERVER_PASSWORD=... bash scripts/sync-remote-assets.sh
# 或使用已存在的 TEST_SERVER_PASSWORD
```

## 4. 开发记录

### 2026-05-17

已完成：

1. 读取项目 QA skill：`/data/project/ai-study/.codex/skills/ai-study-miniprogram-qa/SKILL.md`
2. 读取远程同步 skill：`/data/project/auto-check/.codex/skills/autocheck-server-sync/SKILL.md`
3. 确认用户给的 `/data/auto-check/...` 本机不存在，实际路径为 `/data/project/auto-check/...`
4. 确认 `collect-ui.top` 解析到 `202.140.140.117`
5. 确认远程 `/data/file` 存在
6. 确认远程 80 被 `/data/dist/bin` 占用，443 空闲
7. 同步当前 `miniprogram/assets` 的 13 个图片文件到 `/data/file/ai-study/assets`
8. 创建远程 manifest：`/data/file/ai-study/asset-manifest.sha256`
9. 部署并启动 `ai-study-file.service`
10. 修改小程序图片引用为远程 HTTPS URL：
    - `miniprogram/config/assets.js`
    - `miniprogram/pages/home/home.js`
    - `miniprogram/pages/home/home.wxml`
    - `miniprogram/components/app-top-bar/app-top-bar.js`
    - `miniprogram/components/app-top-bar/app-top-bar.wxml`
11. 删除 `miniprogram/assets`
12. 新增脚本：
    - `scripts/sync-remote-assets.sh`
    - `scripts/deploy-remote-file-service.sh`
    - `scripts/remote-file-server.py`
13. 对远程静态服务做稳定性与安全收敛：
    - 禁用目录列表
    - 屏蔽 `/bin` 和 `/certs`
    - 将 TLS 握手移出主 accept 循环
    - 将 listen backlog 从默认值提升到 128
14. 2026-05-17 12:15 使用 ACME TLS-ALPN-01 签发 Let's Encrypt 证书，并替换远程证书文件。
15. 2026-05-17 12:16 为小程序资源 URL 增加 `?v=20260517-ca1`，避免真机缓存旧的自签失败请求。
16. 2026-05-17 12:20 重新生成当前状态预览二维码。

## 5. 测试记录

已执行并通过：

```bash
bash -n scripts/sync-remote-assets.sh
bash -n scripts/deploy-remote-file-service.sh
python3 -m py_compile scripts/remote-file-server.py
find miniprogram -name '*.js' -print0 | xargs -0 -n1 node --check
find miniprogram -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool
```

远程服务验证：

```bash
curl --noproxy '*' -k -I https://collect-ui.top/ai-study/assets/prototype/logo.png
```

结果要点：

```text
HTTP/1.0 200 OK
Content-type: image/png
Content-Length: 712
Cache-Control: public, max-age=31536000, immutable
```

远程服务状态：

```bash
systemctl is-active ai-study-file.service
```

结果：

```text
active
```

443 监听：

```text
0.0.0.0:443 users:(("python3",pid=4004575,fd=3)) backlog=128
```

屏蔽服务目录验证：

```bash
curl --noproxy '*' -k -I https://collect-ui.top/certs/
```

结果：

```text
HTTP/1.0 404 File not found
```

证书验证：

```text
subject=CN = collect-ui.top
issuer=C = US, O = Let's Encrypt, CN = E8
notAfter=Aug 15 03:16:50 2026 GMT
Verify return code: 0 (ok)
SAN: DNS:collect-ui.top, IP Address:202.140.140.117
```

不带 `-k` 的 HTTPS 资源验证：

```bash
curl --noproxy '*' -I https://collect-ui.top/ai-study/assets/prototype/logo.png?v=20260517-ca1
```

结果：

```text
HTTP/1.0 200 OK
Content-type: image/png
Content-Length: 712
```

包内附件检查：

```bash
find miniprogram -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.gif' -o -iname '*.webp' -o -iname '*.svg' -o -iname '*.mp3' -o -iname '*.mp4' -o -iname '*.pdf' -o -iname '*.doc' -o -iname '*.docx' -o -iname '*.xls' -o -iname '*.xlsx' -o -iname '*.zip' \) -print
```

结果：无输出。

微信开发者工具烟测尝试：

```bash
bash scripts/start-devtools.sh
```

结果：

```text
WeChat DevTools HTTP service did not start on 127.0.0.1:3799.
See /data/project/ai-study/.runtime/wechat-devtools.log
```

日志要点：

```text
Opening in existing browser session.
```

结论：默认 `3799` 端口未开放，但当前 DevTools 实际 IDE 端口为 `16794`、automator 端口为 `9420`，后续已用这两个端口补齐真实模拟器自测。

2026-05-17 12:17 已补齐开发者工具自动化自测，使用当前 IDE 端口 `16794`、automator 端口 `9420`：

```text
homeImageCount=6
dashboardImageCount=2
profileImageCount=2
exceptionCount=0
```

首页远程图片 URL：

```text
https://collect-ui.top/ai-study/assets/prototype/logo.png?v=20260517-ca1
https://collect-ui.top/ai-study/assets/prototype/bell.png?v=20260517-ca1
https://collect-ui.top/ai-study/assets/prototype/book.png?v=20260517-ca1
https://collect-ui.top/ai-study/assets/prototype/sigma.png?v=20260517-ca1
https://collect-ui.top/ai-study/assets/prototype/translate.png?v=20260517-ca1
https://collect-ui.top/ai-study/assets/prototype/sparkle.png?v=20260517-ca1
```

截图证据：

```text
wx_login/screenshots/selftest/remote-assets-home-login-app.png
wx_login/screenshots/selftest/remote-assets-dashboard-app.png
wx_login/screenshots/selftest/remote-assets-profile-app.png
wx_login/screenshots/selftest/remote-assets-learn-actions-app.png
wx_login/screenshots/selftest/remote-assets-checkin-done-app.png
wx_login/screenshots/selftest/remote-assets-checkin-reset-app.png
```

全页面路由烟测：

```text
13/13 pages ok
failedRoutes=[]
exceptionCount=0
```

learn/checkin 功能烟测：

```text
learnCurrentIndex=1
learnCurrentWord=review
checkinCheckedToday=false
checkinStreak=0
```

远程日志确认图片请求全部 200：

```text
GET /ai-study/assets/prototype/logo.png?v=20260517-ca1 200
GET /ai-study/assets/prototype/bell.png?v=20260517-ca1 200
GET /ai-study/assets/prototype/book.png?v=20260517-ca1 200
GET /ai-study/assets/prototype/sigma.png?v=20260517-ca1 200
GET /ai-study/assets/prototype/translate.png?v=20260517-ca1 200
GET /ai-study/assets/prototype/sparkle.png?v=20260517-ca1 200
```

当前预览包：

```text
TOTAL 102.7 KB / 105119 bytes
QR: wx_login/qr/ai-study-preview-qr-real.png
```

## 6. 验收标准

必须全部满足：

1. `curl --noproxy '*' -I https://collect-ui.top/ai-study/assets/prototype/logo.png?v=20260517-ca1` 返回 `200 OK`，不能依赖 `-k`。
2. `systemctl is-active ai-study-file.service` 返回 `active`。
3. `/data/file/ai-study/assets` 中存在同步后的资源文件。
4. `feature/remote-assets/asset-manifest.json` 与远程 `/data/file/ai-study/asset-manifest.sha256` 内容一致。
5. `miniprogram` 内无图片、音视频、PDF、Office、zip 等附件类文件。
6. 小程序所有图片引用从 `miniprogram/config/assets.js` 生成，不再硬编码 `/assets/...`。
7. 微信开发者工具中预览页面时，首页 logo、通知图标、科目图标和提示图标能显示。
8. 正式发布前，在微信公众平台后台配置 `collect-ui.top` 合法域名。

## 7. 后续同步流程

新增资源时：

1. 不要把图片或附件放入 `miniprogram`。
2. 将资源按远程相对路径放入 `feature/remote-assets/staging`，例如 `feature/remote-assets/staging/course/unit-1/cover.webp`。
3. 执行：

```bash
bash scripts/sync-remote-assets.sh
```

4. 如果代码要引用新资源，在 `miniprogram/config/assets.js` 增加 key，或直接用 `assetUrl("course/unit-1/cover.webp")`。
5. 执行静态检查和远程 URL 校验。
6. 同步成功后，可以清空 `feature/remote-assets/staging`，不要把 staging 资源提交到仓库。

## 8. 风险与待办

- Let's Encrypt 证书有效期约 90 天，已安装 acme.sh cron 自动续期；续期需要 443 可被 TLS-ALPN-01 临时占用。
- 微信正式环境需要后台配置合法域名；否则真机或线上版本可能拒绝加载远程资源。
- 如果资源数量变多，需要补充 CDN、版本号路径或文件名 hash，避免长期缓存导致旧资源不更新。
- 如果未来附件涉及权限控制，不能继续使用完全公开静态目录，需要改为鉴权下载或签名 URL。
