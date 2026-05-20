# AI Study

一个微信小程序项目，当前首页按 `原型/screen.png` 还原为 AI 学习测评入口。项目根目录就是微信开发者工具可打开的项目目录，实际小程序源码放在 `miniprogram/`。

## 启动项目

```bash
./scripts/start-devtools.sh
```

脚本会做三件事：

1. 使用 `:99` 的 Xvfb 虚拟显示器。
2. 用项目内置的微信开发者工具启动 IDE HTTP 服务，默认端口是 `3799`。
3. 打开当前项目并信任项目。

检查登录和服务状态：

```bash
./scripts/status-devtools.sh
```

重新生成项目首页成功截图：

```bash
./scripts/capture-start-success.sh
```

截图输出到 `wx_login/screenshots/preview/project-start-success.png`。

生成微信预览二维码：

```bash
./scripts/preview-qr.sh
```

二维码输出到 `wx_login/qr/ai-study-preview-qr.png`，预览信息输出到 `wx_login/qr/ai-study-preview-info.json`。
微信开发者工具 Linux CLI 生成的图片扩展名是 `.png`，实际内容可能是 JPEG；脚本会额外生成真正 PNG 白底图 `wx_login/qr/ai-study-preview-qr-real.png`，扫码优先用这张。

`wx_login/` 只作为本地调试产物目录，按类别归档：

- `wx_login/qr/`：预览二维码和预览信息。
- `wx_login/screenshots/prototype/`：原型对比截图、裁剪图和测量结果。
- `wx_login/screenshots/preview/`：打开项目、预览过程截图。
- `wx_login/screenshots/selftest/`：自动化自测截图。

## 项目结构

```text
project.config.json          微信开发者工具项目配置
miniprogram/app.json         页面和窗口配置
miniprogram/app.js           全局初始化逻辑
miniprogram/app.wxss         全局样式
miniprogram/utils/           示例学习数据
miniprogram/pages/home/      首页
miniprogram/pages/learn/     学习页
miniprogram/pages/checkin/   打卡页
scripts/                     本地启动、状态检查和截图脚本
.codex/skills/               当前项目专用 Codex skill
```

当前项目的开发测试与验证流程定义在 `.codex/skills/ai-study-miniprogram-qa/SKILL.md`，覆盖原型对比、模拟器功能测试、截图留证、二维码生成和 `wx_login/` 产物清理规范。

当前 `project.config.json` 使用正式小程序 AppID：`wxd23932c04171118f`。基础库固定为 `3.8.12`，避免每次启动都依赖 `latest`。

## 页面说明

- 首页展示登录入口、手机号/验证码输入、学段切换、年级选择、科目选择和开始测评按钮。
- 学习页展示单词、音标、释义和例句，点击“学会了”会累加本地学习进度。
- 打卡页展示累计数据，支持完成今日打卡和重置示例进度。

数据暂时保存在微信小程序本地缓存里，没有接入后端服务。
