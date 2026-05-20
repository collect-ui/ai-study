Run ID: `20260517-student-perspective-selftest`

日期: `2026-05-17`

执行人: Codex

范围: 站在学生使用者角度，对当前小程序进行真实 DevTools + miniprogram-automator 自测。

环境:

- DevTools HTTP: 实际连接端口 `16794`。
- Automator: `127.0.0.1:9420`。
- Xvfb: `DISPLAY=:99`。
- 说明: 固定 `--port 3799` 在本轮环境会让 CLI 误等旧端口，因此预览二维码使用 `WECHAT_DEVTOOLS_PORT=16794` 生成。

静态检查:

- `node --check miniprogram/pages/home/home.js` 通过。
- `python3 -m json.tool miniprogram/app.json` 通过。
- `python3 -m json.tool miniprogram/pages/home/home.json` 通过。
- `find miniprogram -name '*.js' -print -exec node --check {} \;` 通过。
- `find miniprogram/pages miniprogram/components -name '*.json' -print -exec python3 -m json.tool {} \;` 通过。

学生路径自测结果:

- 入口: 通知、手机号错误提示、验证码发送态、初中/初三/英语选择、底部开始测评均通过。
- 首页: 拍照录入、在线答疑、AI 助教反馈均可触发；短板测试可进入测评；自主学习可进入专项学习设置。
- 测评: 空题提交提示、单选、填空、判断、连续答题、完成测评均通过。
- 报告: 本轮答题结果为 `100` 分，`6` 题正确、`0` 题错误；咨询空表单校验和有效提交均通过。
- 学习: 认读模式播放、辅助入口、模糊/认识反馈、学习报告均通过；跟读模式标准范读、录音中、评分结果、学习报告均通过。
- 知识/错题: 知识概览、能力卡进入错题、分类筛选、搜索、解析展开、进入复习均通过。
- 我的: 个人概览和学习任务入口通过。
- 旧示例: `learn` 的听/下一个/学会了、`checkin` 的完成打卡通过。

学生体验观察:

- 未发现阻断学生完成主路径的问题。入口到测评、报告、专项学习、错题复习的路径是连续的。
- `拍照录入`、`在线答疑`、`AI 助教` 当前是“待接入”反馈。作为学生会理解为功能入口已经存在，但无法真正使用；建议后续明确标注试用/即将开放，或进入一个说明页。
- 报告为满分时仍展示“阅读速度和长难句拆分仍有提升空间”的固定建议。学生可能觉得和 `100` 分不完全匹配；建议按分数或错题分布动态调整文案。
- 咨询表单使用“姓名/联系电话”。如果目标用户是学生，建议文案区分“学生姓名”和“家长联系电话”，降低未成年人填写个人电话的歧义。

截图产物:

- 学生路径 App/模拟器截图: `wx_login/screenshots/selftest/student-*.png`，共 37 张，其中 `*-simulator.png` 4 张。
- 保留上一轮全模块截图: `wx_login/screenshots/selftest/qa-*.png`，共 60 张。
- 原型对比图: `wx_login/screenshots/prototype/qa-compare-*.png`，共 11 张。

预览:

- `WECHAT_DEVTOOLS_PORT=16794 ./scripts/preview-qr.sh` 通过。
- 预览包大小: `108003` bytes。
- 二维码产物在 `wx_login/qr/`。

结论:

- 学生视角自测通过，自动化过程未捕获小程序运行异常。
- 当前主要风险不是流程阻断，而是部分入口仍为占位功能、报告建议文案未按满分场景细化。
