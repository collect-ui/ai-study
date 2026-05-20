Run ID: `20260517-full-module-qa`

日期: `2026-05-17`

执行人: Codex

范围: 使用 `ai-study-miniprogram-qa` 技能对当前小程序做全量功能模块测试，并按模块输出测试截图。

静态检查:

- `find miniprogram -name '*.js' -print -exec node --check {} \;` 通过。
- `python3 -m json.tool miniprogram/app.json` 通过。
- `find miniprogram/pages miniprogram/components -name '*.json' -print -exec python3 -m json.tool {} >/dev/null \;` 通过。

自动化环境:

- DevTools HTTP: `127.0.0.1:3799`
- Automator: `127.0.0.1:9420`
- `miniprogram-automator` 已连接真实模拟器执行测试。

功能模块测试结果:

- 入口登录模块: 通知、手机号、验证码、获取验证码发送态、初中/初三/英语选择、底部开始测评、游客进入、登录进入首页均已截图。
- 首页模块: 首页展示、拍照录入反馈、AI 助教反馈、自主学习入口、短板测试入口均已截图。
- 测评模块: 单选、填空、判断、上一题、完成测评均已截图。
- 报告模块: 报告结果、空表单校验、有效咨询提交、查看错题、知识掌握入口、生成提升练习均已截图。
- 专项学习模块: 年级/单元/模式选择、认读入口、跟读入口均已截图。
- 认读模块: 初始词卡、播放、辅助入口、模糊反馈、认读报告、报告播放均已截图。
- 跟读模块: 初始句子、标准范读、录音中、录音完成、学习报告均已截图。
- 知识模块: 知识概览、顶部操作反馈、针对性练习入口均已截图。
- 错题模块: 列表、分类筛选、搜索、解析展开、进入复习均已截图。
- 我的模块: 概览、学习任务入口、知识入口、错题入口均已截图。
- 旧示例模块: `learn` 的听/下一个/学会了，`checkin` 的完成打卡/重置进度均已截图。

截图产物:

- App 截图: `wx_login/screenshots/selftest/qa-*.png`，共 60 张。
- 模拟器整窗截图: `wx_login/screenshots/selftest/qa-*-module-simulator.png`。
- 原型对比图: `wx_login/screenshots/prototype/qa-compare-*.png`，共 11 张。

预览:

- `./scripts/preview-qr.sh` 通过。
- 预览包大小: `108003` bytes。
- 二维码产物在 `wx_login/qr/`。

清理:

- 删除本轮调试图 `qa-debug-home-app.png`。
- 删除上一轮 `stitch-ai-*.png` 旧测试截图，避免干扰本轮模块证据。

结论:

- 当前全量功能模块自动化测试通过，未捕获小程序运行异常。
- 截图数量超过常规 20 张上限是因为本轮明确要求按全部功能模块输出测试图，已清理旧证据并保留本轮可追溯证据。
