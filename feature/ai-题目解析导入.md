# 目的
将pdf 的题目一键导入，转换成考试题目，减少教师的工作量
# 文件
/data/project/ai-study/题目/【小升初】英语复习题三十套全国通用（含详细解析）.pdf
# 需求
- 目前题目支持手工录入，我希望增加ai 解析并且导入
- 将pdf 解析成一个固定字符串，弹框前端确认，左边 一堆字符串输入框文本，右边一个表格对应题目
- 然后前端解析字符串，生成一条一条的记录
- 确认之后一把入口
# 地址
http://192.168.232.130:8026/collect-ui#/collect-ui/framework/question-bank

# 模型支持
- 支持deepseek v4 聊天模型，我配置文件配置key
- 支持codex auth.json 模式
  - /data/project/sport/plugins/module_agent_run.go 这个里面有处理ai auth.json 可以参考这个
- 都是直接chat聊天返回固定格式
- 前端页面支持选择2种模式解析


# 测试要求
测试要求：

- 使用无头浏览器打开目标页面，按用户真实路径完成操作。
- 记录 console error、pageerror、requestfailed。
- 保存 JSON 报告和关键截图。
- 失败时先根据截图和 DOM 证据修复，再重复验证直到通过。

需要再本文后面追加设计文档、验收、测试计划，以及运行过程中执行日志

---

# 设计文档

## 范围

本次实现覆盖 AI Study 管理后台题库页：

- 在题库管理页增加“AI解析导入”入口。
- 弹框左侧支持读取 PDF 文本、编辑 AI 固定 JSON 字符串。
- 弹框右侧预览解析后的题目表格。
- 确认后复用现有 `question.question_choice_save` 服务逐条入库。
- 支持后端 DeepSeek Chat 和 Codex `~/.codex/auth.json` 两种 AI 模式。

## 后端

- 新增插件：
  - `QuestionPDFTextService`：抽取 PDF 文本，服务为 `question.ai_pdf_text`。
  - `QuestionAIParseService`：调用 AI 或解析 `mock_response/fixed_text`，服务为 `question.ai_parse`。
- 新增配置项：
  - `question_ai_default_pdf_path`
  - `question_ai_source_max_chars`
  - `question_ai_deepseek_api_key`
  - `question_ai_deepseek_base_url`
  - `question_ai_deepseek_model`
  - `question_ai_codex_model`
- DeepSeek 密钥不写入文档和代码，运行时放到配置文件或 `DEEPSEEK_API_KEY`。
- Codex 模式优先读 `OPENAI_API_KEY` 环境变量，再读 `~/.codex/auth.json`。

## 前端

- 题库页 `question_bank.json` 使用低代码配置实现入口按钮、导入弹框、表单、上传、预览表格和动作流。
- 页面支持手动填写 DeepSeek Key，并通过 `question.ai_parse` 的 `api_key/base_url` 参数传给后端。
- 已移除原自定义组件 `question-ai-import`，不再加入 minimal 白名单。
- 固定 JSON 格式：

```json
{"questions":[{"stem_text":"...","option_a_text":"...","option_b_text":"...","option_c_text":"...","option_d_text":"...","answer_key":"A","analysis_text":"..."}]}
```

## 入库策略

- `question.ai_parse` 负责把 AI 输出或固定字符串解析为题目数组，低代码页面负责补齐入库字段与预览状态。
- 表格会标记缺题干、缺选项、缺答案等问题。
- “确认导入”只导入有效行，并通过低代码后端 `question.ai_import_save` 批量调用 `question.question_choice_save`。
- 导入记录使用 `source=ai_import`，手工录入仍保持 `source=manual`。

# 验收

- 题库页显示“AI解析导入”按钮。
- 弹框能读取示例 PDF 文本。
- 固定 JSON 字符串能解析成表格记录。
- 有效记录能通过确认导入写入题库。
- 导入后题库列表刷新。
- 浏览器验收记录 console error、pageerror、requestfailed、HTTP 失败。
- 生成 JSON 报告和关键截图。

# 测试计划

1. 后端静态校验：`go test ./...`、`go build ./...`、`go vet ./...`。
2. 前端构建：`COLLECT_UI_ENTRY=./src/index.min.tsx npm run build`。
3. AI Study 静态壳构建部署：`bash scripts/deploy_ai_study_collect_ui.sh`。
4. 浏览器路径验收：
   - 登录后台。
   - 打开题库管理页。
   - 打开 AI 导入弹框。
   - 读取示例 PDF。
   - 粘贴固定 JSON 字符串。
   - 解析成右侧表格。
   - 确认导入。
   - 查询导入题目并清理验收数据。

# 执行日志

执行时间：2026-05-20

- 已读取项目 QA skill：`/data/project/ai-study/.codex/skills/ai-study-miniprogram-qa/SKILL.md`
- 已实现后端插件：`backend/plugins/module_question_ai_import.go`
- 已注册后端插件：`backend/plugins/a_register.go`
- 已新增服务配置：`backend/collect/question/import/index.yml`
- 已新增模块注册：`backend/collect/service_router.yml`
- 已新增 AI 配置占位：`backend/conf/application.properties`
- 已调整题目保存配置支持 `source=ai_import`：`backend/collect/question/question/index.yml`
- 已改为低代码导入弹框：`backend/collect/frontend/page_data/data/question/question_bank.json`
- 已新增低代码批量导入服务：`question.ai_import_save`
- 已为 DeepSeek 手动 key 补充服务参数：`question.ai_parse.api_key/base_url`
- 已移除 collect-ui 自定义组件：`/data/project/collect-ui/src/components/question-ai-import/question-ai-import.tsx`
- 已从 minimal 白名单移除 `question-ai-import`
- 已接入题库页面：`backend/collect/frontend/page_data/data/question/question_bank.json`

验证结果：

- `go test ./...`：通过
- `go build ./...`：通过
- `go vet ./...`：通过
- `python3 -m json.tool backend/collect/frontend/page_data/data/question/question_bank.json`：通过
- `COLLECT_UI_ENTRY=./src/index.min.tsx npm run build`：通过
- `bash scripts/deploy_ai_study_collect_ui.sh`：通过，6 个页面 smoke 全部通过
- 后端重启与访问校验：
  - `ss -ltnp | rg ':8026'`：端口监听正常
  - `curl --noproxy '*' http://127.0.0.1:8026/collect-ui/`：HTTP 200

浏览器验收报告：

- JSON：`feature/ai-question-import-verify/report.json`
- 截图：
  - `feature/ai-question-import-verify/dialog-open.png`
  - `feature/ai-question-import-verify/parsed-table.png`
  - `feature/ai-question-import-verify/imported-result.png`
- 报告结论：`ok=true`
- console error：0
- pageerror：0
- requestfailed：0
- HTTP 失败：0
- PDF 抽取：示例 PDF 抽取到 24000 字符
- 导入验收：固定 JSON 解析 1 道题，导入成功后已按 `question_id` 清理验收数据

## 2026-05-20 大 PDF 分段解析优化

问题：示例 PDF 体量较大，原链路默认只读取/传入前 24000 字符，AI 单段解析容易只得到约 75 条题目，无法覆盖整本资料。

调整：

- `question.ai_pdf_text` 默认读取上限从 24000 字符调整为 120000 字符。
- `question.ai_parse` 新增 `parse_mode`、`source_max_chars`、`chunk_chars`、`max_chunks` 参数。
- 默认解析方式改为 `chunked`，按约 12000 字符分段调用 AI，再汇总 `questions` 并按 `question_code` 或题干去重。
- 题库页 AI 导入弹框增加“解析方式”，默认“分段解析”。
- 页面读取 PDF 后提示全文字符数和后续分段策略，AI 解析完成后显示题目数、分段数、原文字符数和截断状态。
- 补充 Go 单元测试覆盖文本分段、汇总去重、固定 JSON 解析。

验证：

- `python3 -m json.tool backend/collect/frontend/page_data/data/question/question_bank.json`：通过
- `backend/collect/question/import/index.yml` YAML 解析：通过
- `go test ./...`：通过
- `go vet ./...`：通过
- `go build ./...`：通过
- 后台重启后 `/collect-ui/` HTTP 200
- 浏览器验收脚本：`node feature/ai-question-import-verify/verify-ai-question-import.mjs`

浏览器验收结果：

- 报告：`feature/ai-question-import-verify/report.json`
- 截图：
  - `feature/ai-question-import-verify/dialog-open.png`
  - `feature/ai-question-import-verify/parsed-table.png`
  - `feature/ai-question-import-verify/imported-result.png`
- `ok=true`
- 示例 PDF 读取长度：81674 字符，确认已突破 24000 字符截断
- console error：0
- pageerror：0
- requestfailed：0
- HTTP 失败：0
- 固定 JSON 解析 1 道题，导入后已回查并清理验收数据

补充核对：

- `test-results/full-pdf-import/import-summary.json` 当前全量离线解析记录为 31 套、246 条。
- 本地数据库当前 `source='pdf_import'` 有 246 条有效题目。

## 2026-05-20 题型与字段对齐修复

问题：AI 导入预览阶段会把阅读理解/完形填空的 `choice_items` 覆盖成空 `reference_text`，导致入库后只能看到大题题干，看不全小题、选项、答案和解析；普通选择、阅读理解等题型字段也缺少统一标准化。

调整：

- 后端 `question.ai_parse` 增加题目行标准化：
  - 中文题型/分类会归一到 `reading_choice`、`cloze_choice`、`grammar_choice`、`blank`、`judge` 等内部编码。
  - 阅读理解/完形填空按“一篇材料一条大题”返回，`choice_items` 保留小题、选项、答案、解析。
  - 自动生成 `reference_text`、`answer_text`、`answer_value`、`blank_count`、`option_count`。
  - 普通选择题补齐 `option_a_text` 到 `option_d_text`、`answer_value` 和正确选项标记。
- 低代码导入弹框预览映射修复：
  - `reading_choice`/`cloze_choice` 不再按普通选择题校验顶层 A-D。
  - 预览行保留 `choice_items`，确认导入时写入 `question_answer.reference_text`。
  - 保留 `question_code`、`remark` 等来源字段。
- `question.question_choice_save` 补传 `question_code`、`remark` 到题目保存服务。

验证：

- 新增 Go 单元测试覆盖阅读理解 `choice_items` 保留。
- 浏览器验收固定 JSON 同时导入：
  - 普通单选 1 条
  - 阅读理解 1 条，包含 1 个小题
- 导入后通过 `question.question_choice_detail` 回读确认：
  - `question_category=reading_choice`
  - `question_type=single_choice`
  - `choice_items[0].question_text=Where is Tom?`
  - `choice_items[0].option_b=At home`
  - `choice_items[0].answer_key=B`
- 验收数据已清理，本地数据库无残留验收题。

最新验证：

- `python3 -m json.tool backend/collect/frontend/page_data/data/question/question_bank.json`：通过
- `node --check feature/ai-question-import-verify/verify-ai-question-import.mjs`：通过
- `go test ./...`：通过
- `go vet ./...`：通过
- `go build ./...`：通过
- 浏览器验收报告：`feature/ai-question-import-verify/report.json`，`ok=true`
