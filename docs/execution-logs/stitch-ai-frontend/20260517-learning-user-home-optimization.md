Run ID: `20260517-learning-user-home-optimization`

日期: `2026-05-17`

执行人: Codex

目标:

- 站在初中生学习使用者角度优化首页。
- 首页必须能看出学习进度、该学习什么知识、课程有哪些。
- 移除首页“为你推荐”，并尽量保持和原型一致的浅色背景、白色卡片、蓝色主行动风格。

代码变更:

- `miniprogram/utils/mock-data.js`: 为首页补充 `grade`、`subject`、`progress`、`knowledgePlan`、`courses`，并把占位功能标为“即将开放”。
- `miniprogram/pages/dashboard/dashboard.wxml`: 移除推荐区，新增学习进度、当前知识计划、课程列表。
- `miniprogram/pages/dashboard/dashboard.wxss`: 新增学习进度、知识计划、课程列表样式，沿用原有卡片体系。
- `miniprogram/pages/dashboard/dashboard.js`: 新增知识计划和课程点击流程，课程可进入学习设置或测评。

静态检查:

- `node --check miniprogram/pages/dashboard/dashboard.js` 通过。
- `node --check miniprogram/utils/mock-data.js` 通过。
- `python3 -m json.tool miniprogram/app.json` 通过。
- `python3 -m json.tool miniprogram/pages/dashboard/dashboard.json` 通过。
- `find miniprogram -name '*.js' -print -exec node --check {} \;` 通过。
- `find miniprogram/pages miniprogram/components -name '*.json' -print -exec python3 -m json.tool {} \;` 通过。
- `rg -n "为你推荐|recommend|recommendations|openRecommendation" miniprogram` 无命中。

自动化自测:

- Round 1: 新首页结构与跳转验证通过。
  - 首页渲染 `3` 个知识项、`3` 个课程项。
  - `为你推荐` 未出现在页面文本中。
  - 知识计划: 练习进入测评、词汇进入认读设置、跟读进入跟读设置均通过。
  - 课程列表: Unit 2 进入学习设置、语法专项进入测评均通过。
  - 工具入口: 短板测试进入测评、自主学习进入学习设置均通过。
- Round 2: 初中生端到端学习路径通过。
  - 入口选择初三英语后进入首页。
  - 首页课程 Unit 1 进入认读并生成学习报告。
  - 语法专项进入测评并生成报告，结果 `100` 分，`6` 正确，`0` 错误。
- Round 3: 非首页页面烟测通过。
  - 知识、错题、我的、旧 `learn`、旧 `checkin` 页面均可打开并完成关键动作。
- 三轮自测 `exceptionCount` 均为 `0`。

截图产物:

- 本轮优化截图: `wx_login/screenshots/selftest/optimized-*.png`，共 `23` 张。
- 关键截图:
  - `optimized-dashboard-top-app.png`
  - `optimized-dashboard-knowledge-courses-app.png`
  - `optimized-dashboard-simulator.png`
  - `optimized-flow-dashboard-from-entry-app.png`
  - `optimized-flow-course-unit1-study-report-app.png`
  - `optimized-flow-grammar-report-app.png`

预览:

- `WECHAT_DEVTOOLS_PORT=16794 ./scripts/preview-qr.sh` 通过。
- 预览包大小: `111429` bytes。
- 二维码产物在 `wx_login/qr/`。

结论:

- 首页已从“推荐内容”改为“学习进度 + 当前该学知识 + 课程列表”的学习者视角。
- 当前流程可让初中生明确知道自己学到哪里、下一步该补什么、可进入哪些课程。
