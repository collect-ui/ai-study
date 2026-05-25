# PDF 知识题目统一导入还原设计

## 需求目的

本功能面向教师，不是后台技术人员的批处理工具。教师需要在现有知识点维护页面上传 PDF，一键生成可维护、可追溯、可复核的知识点。

核心目的：

- 教师在知识点维护页完成 PDF 导入，不进入单独菜单，不跳到独立导入页面。
- PDF 导入后自动生成 `Unit -> 知识点 -> 知识正文`，减少手工录入。
- 知识点必须保留 PDF 原始来源，能按文件名、页码、片段和字段取值回看。
- AI 只能辅助抽取和整理，不能把无法确认的内容静默写入正式知识库。
- 题目型或混合型 PDF 可以识别题目，但知识点导入的主目标仍是自动生成知识点。
- 失败项必须可见、可复核、可重新解析，方便教师和教研人员验收。

## 页面方案

入口固定在现有页面：

`http://192.168.232.130:8026/collect-ui/#/collect-ui/framework/question-taxonomy`

页面位置：

- 页面：`question-taxonomy`。
- 区域：右侧 `知识点维护` 面板。
- 按钮区：与现有 `新增知识点` 同一组顶部操作区。
- 按钮：`PDF导入知识点`。
- 展示方式：点击后打开 `PDF知识点导入` 弹窗或抽屉。
- 菜单约束：不新增单独菜单，不新增独立导入页面。

教师使用路径：

```text
打开知识点维护页
  -> 点击 PDF导入知识点
  -> 上传 PDF
  -> 选择或确认学科/年级/教材/单元
  -> 点击 一键生成知识点
  -> 查看生成结果
  -> 查看来源对比和失败复核
  -> 确认入库或修正后入库
```

页面必须展示：

- 当前 PDF 文件名、页数、抽取字符数。
- 生成的单元、知识点、知识正文数量。
- 自动入库数量、待复核数量、跳过重复数量。
- 每个知识点对应的来源页码、来源片段、字段来源状态。
- `按本文件过滤来源`，只看当前 PDF 生成的知识点和来源片段。
- `来源对比`，左侧看 PDF 原始片段，右侧看整理后的知识点内容。
- `失败复核`，展示 AI 无法确认、来源缺失、跨页边界不清、疑似题目误判等内容。

## 测试验收计划

验收分为服务验收、数据验收、页面验收和教师流程验收。验收目标不是“接口返回成功”，而是教师能在现有知识点维护页完成一键导入，并能核对 PDF 原文和生成知识点是否一致。

服务验收：

- 上传或指定 PDF 后，必须调用 `question.ai_pdf_text` 作为唯一 PDF 文本入口。
- `question.pdf_knowledge_import_one_click` 能返回导入批次、来源文档、生成数量、待复核数量和验收状态。
- 纯知识点 PDF 能生成知识点，不自动生成题目。
- 知识点 + 题目混合 PDF 能生成知识点，并把真实题目放入草稿或待确认队列。
- 解析失败、来源缺失、AI 推断答案或字段不可追溯时，必须生成失败复核项。

数据验收：

- 每次导入都有 `source_doc_id` 和 `import_batch_id`。
- 每个知识点正文都有 `question_source_fragment` 来源片段。
- 每个关键字段都有 `question_source_field_rel.raw_quote`。
- 能按 `file_name` 查到本文件生成的知识点、来源片段、字段来源和失败项。
- 重复内容要能说明是跳过、合并还是生成新版本。

页面验收：

- 打开 `http://192.168.232.130:8026/collect-ui/#/collect-ui/framework/question-taxonomy`。
- 页面不出现新的知识点导入菜单。
- `知识点维护` 顶部能看到 `PDF导入知识点` 按钮。
- 点击按钮后出现 `PDF知识点导入` 弹窗或抽屉。
- 弹窗能上传 PDF，并显示学科、年级、教材、单元、导入模式和题目处理选项。
- 点击 `一键生成知识点` 后能看到进度、结果摘要和失败复核入口。
- 导入成功后知识点列表自动刷新。
- 点击 `按本文件过滤来源` 后，只显示当前 PDF 相关知识点和来源片段。
- 点击某个知识点的 `来源对比` 后，能看到 PDF 原始片段、前文、后文和整理后内容。
- 失败复核项能定位到原文块、片段和页图。
- 页面验收必须采集 `console.error`、`pageerror`、关键接口失败和资源加载失败。

教师流程验收：

- 教师只需要在一个页面完成上传、生成、查看、复核和确认。
- 教师能判断“这个知识点来自 PDF 哪一段”。
- 教师能按文件名过滤本次导入结果。
- 教师能看到哪些内容已经入库，哪些内容需要复核。
- 教师不需要理解数据库表，也不需要进入题库页才能完成知识点导入。

## 1. 统一结论

知识点总结 PDF、全册知识考点 PDF、易错题 PDF 都走同一套导入、解析、落库、展示、还原设计。

统一入口：

```json
{
  "service": "question.ai_pdf_text",
  "file_path": "/data/project/ai-study/题目/9英语易错题-秒懂比较级句型.pdf",
  "max_chars": 120000
}
```

统一原则：

- 不在前端或新服务里直接调用外部 OCR。
- PDF 文本入口只用 `question.ai_pdf_text`，腾讯识别、文本层读取、fallback 都由后端封装。
- AI 只消费 `question.ai_pdf_text` 返回的 `raw_text`、`raw_html`、`file_name`、`source_chars` 和抽取参数。
- AI 解析不是无脑识别：必须由大模型返回明确 JSON 结构，区分 Unit、知识点、知识正文、候选题、真实题目、答案解析和失败项。
- 后端 Go 代码不得通过关键词、题号、标题、正则、字符串截取等方式判断某段文本是题目、答案、解析、Unit、知识点或知识正文。
- 后端 Go 代码只负责 PDF 文本入口、大模型调用、JSON 外形校验、来源保存、业务表落库和错误返回；结构理解必须统一交给大模型提示词。
- 原始信息必须保存，知识点界面和题目界面都能看到来源原文。
- 来源不能只停留在整篇 PDF：题目、选项、答案、解析、知识点正文都必须能追到具体页、块、片段和字段取值来源。
- 教师侧必须有明确入口：在知识点维护页提供 PDF 一键导入，上传后自动生成知识点，并展示导入结果和失败复核项。
- 导入和导出必须支持可逆：原 PDF、页面视觉、结构语义三层可还原。

## 1.1 AI 结构化解析边界

PDF 解析结构必须由大模型返回，不能由 Go 代码补规则解决。后端禁止实现以下行为：

- 禁止用 `strings.Contains`、正则、标题词表、题号统计、选项字母统计判断知识点、题目、答案或解析。
- 禁止在 Go 中维护 `Unit 1-6`、`字母考点/单词考点/语法考点`、`答案及解析`、`参考答案` 等业务词表来切分结构。
- 禁止在 Go 中根据题号、选项、空格、问号、答案区标题估算题目数量或候选题数量。
- 禁止在 Go 中把答案区编号和题干区编号做匹配，答案、解析、选项、知识正文的归属必须来自大模型 JSON。
- 禁止在 Go 中把某段原文截取后直接作为知识点标题或正文；Go 只能保存大模型返回的字段，并用大模型返回的 `source_quote` 做来源追溯。

允许 Go 做的事情只限于：

- 调用 `question.ai_pdf_text` 得到 `raw_text/raw_html`。
- 把 `raw_text`、文件信息和默认上下文传给大模型。
- 校验大模型返回是否为合法 JSON、是否包含必需数组和字段类型。
- 根据大模型返回的 ID、名称、正文、`source_quote`、`issue` 写入数据库。
- 用 `source_quote` 在原文中定位字符区间；定位失败时只标记为待复核，不自行推断结构。
- 保存大模型原始返回 `ai_output_json`，便于后续重放、对比和修正提示词。

如果大模型返回结构不正确，处理方式是改提示词、重跑大模型或进入失败复核；不能新增 Go 关键词规则兜底。

### 1.2 大模型返回 JSON 契约

`question.pdf_knowledge_import_one_click` 的核心输入是 `question.ai_pdf_text.raw_text`，核心输出必须是如下 JSON 对象：

```json
{
  "units": [
    {
      "unit_id": "可为空，后端可按字段生成稳定 ID",
      "unit_code": "unit_1",
      "unit_name": "Unit 1 Making friends",
      "subject": "english",
      "stage": "primary",
      "grade": "grade_3",
      "textbook_version": "pep",
      "order_index": 1,
      "source_quote": "原文中可定位该单元的短句或标题",
      "knowledge": [
        {
          "knowledge_id": "可为空，后端可按字段生成稳定 ID",
          "knowledge_code": "greeting_sentence_pattern",
          "knowledge_name": "问候与自我介绍句型",
          "semantic_type": "sentence_pattern",
          "order_index": 1,
          "source_quote": "原文中可定位该知识点的短句",
          "contents": [
            {
              "section_title": "句型结构",
              "semantic_type": "sentence_pattern",
              "content_text": "整理后的知识正文",
              "content_html": "",
              "source_quote": "正文对应的原文片段，必须来自 raw_text",
              "source_quotes": [],
              "page_no": 1,
              "order_index": 1,
              "confidence": 0.95
            }
          ]
        }
      ]
    }
  ],
  "questions": [
    {
      "question_draft": true,
      "question_type": "single_choice",
      "question_category": "grammar_choice",
      "stem_text": "题干文本",
      "options": [
        {"key": "A", "text": "选项原文"}
      ],
      "answer_key": "只允许来自原文答案区或解析区；缺失则为空",
      "analysis_text": "只允许来自原文解析；缺失则为空",
      "knowledge_refs": ["greeting_sentence_pattern"],
      "source_quotes": {
        "stem_text": "题干原文片段",
        "option_a": "A 选项原文片段",
        "answer_key": "答案原文片段",
        "analysis_text": "解析原文片段"
      },
      "review_status": "pending"
    }
  ],
  "issues": [
    {
      "issue_type": "missing_source|low_confidence|candidate_question|missing_answer|structure_uncertain",
      "severity": "info|warning|error",
      "raw_text": "需要复核的原文或大模型说明",
      "source_quote": "可定位的原文片段",
      "expected_schema": "期望的大模型结构",
      "error_msg": "为什么不能直接落库",
      "suggestion": "给教师或教研人员的处理建议",
      "status": "pending",
      "page_no": 1
    }
  ],
  "question_draft_total": 0,
  "acceptance_status": "pass|warning|pending_review|error",
  "summary": "本次解析摘要"
}
```

字段约束：

- `units[].knowledge[].contents[].content_text` 是知识正文最终展示内容，由大模型组织。
- `source_quote` 是来源追溯锚点，必须来自 `raw_text`；后端只负责定位和保存，不用它反向判断结构。
- `questions` 中的答案和解析同样由大模型返回；后端不得用代码从 `raw_text` 中匹配答案区。
- `question_draft_total` 必须由大模型返回；后端不得统计题号或选项数量生成。
- `issues` 必须由大模型返回或由 JSON 校验失败生成；不得因 Go 关键词识别生成候选题。

## 2. 样本矩阵

这里的样本不是三套设计，而是统一管线要覆盖的三类 PDF。界面展示时也不应只显示“文件、类型、结论”三列，而应显示“抽取结果、知识点结果、题目结果、待处理问题”。

### 2.1 知识点总结型

文件：

`/data/project/ai-study/题目/三年级（上）英语 知识点总结《人教版PEP》.pdf`

展示信息：

- 来源类型：知识点总结型。
- 抽取结果：通过 `question.ai_pdf_text` 得到 `raw_text/raw_html`，页图只用于预览和还原。
- 知识点结果：按 Unit、Part、章节生成知识点和知识正文。
- 题目结果：默认不生成题目。
- 待处理问题：图片裁切、章节边界、文本抽取错字。

落库策略：

- 写入 `question_unit`、`question_knowledge`、`question_knowledge_content`、`question_knowledge_asset`。
- 写入 `question_source_snapshot`，保证知识点界面能看到来源原文。
- 写入 `question_source_fragment` 和 `question_source_field_rel`，保证知识点正文每个取值能追到 PDF 片段。

### 2.2 全册考点汇总型

文件：

`/data/project/ai-study/题目/三年级上册英语全册知识考点汇总人教版.pdf`

展示信息：

- 来源类型：全册考点汇总型。
- 抽取结果：纯图片型或文本层为空时仍统一走 `question.ai_pdf_text`。
- 知识点结果：按 Unit 1-6 和 `字母考点/单词考点/短语考点/句型考点/常用表达考点/语法考点` 生成知识点。
- 题目结果：编号考点不是题目，不写入 `question_item`。
- 待处理问题：跨页续接、跨单元边界、重复标题、重复内容、文本抽取错字。

落库策略：

- 先落知识库和来源原文。
- 每个知识条保存来源片段和字段级来源。
- AI 误判成题目的内容写入 `question_pdf_parse_issue`，不静默丢弃。

### 2.3 知识点 + 题目混合型

文件：

`/data/project/ai-study/题目/9英语易错题-秒懂比较级句型.pdf`

展示信息：

- 来源类型：知识点 + 题目混合型。
- 抽取结果：`question.ai_pdf_text` 返回知识讲解、题干、选项、答案、解析。
- 知识点结果：比较级常用句型、同级比较句型、比较级表达最高级。
- 题目结果：6 道选择题，答案和解析从原文回填。
- 待处理问题：题目与答案解析的对应关系、缺失答案、选项边界。

落库策略：

- 知识点写入知识库。
- 题目写入题库草稿或导入结果。
- 每道题必须写入 `question_knowledge_rel`、`question_source_rel`、`question_source_snapshot` 和字段级来源关系。
- 题目 1 这类条目必须能回答：来源文件、页码、片段号、题干原文、选项原文、答案解析原文，以及命中片段前后的正文。

`9英语易错题-秒懂比较级句型.pdf` 使用 `question.ai_pdf_text` 实测成功：

- `source_chars=2684`
- `raw_text` 包含“比较级常用句型”知识点。
- `raw_text` 包含 6 道题的题干、选项、答案、解析片段。

这说明统一设计必须同时支持“知识内容”和“题目内容”，不能做两条割裂链路。前端样本展示也应按上面的卡片/分区结构呈现，避免长路径表格导致信息不可读。

## 3. 统一架构

PDF 的统一中间结构：

```text
source_document
  -> source_page
    -> source_block
      -> source_fragment
        -> semantic_item
          -> knowledge_content
          -> candidate_question
          -> question_item
          -> parse_issue
```

业务映射：

```text
question_unit
  -> question_knowledge
    -> question_knowledge_content
      -> question_knowledge_asset

question_item
  -> question_option / question_answer / question_blank_answer
  -> question_knowledge_rel
  -> question_source_rel
  -> question_source_field_rel

question_source_snapshot
  -> raw_text/raw_html/ai_output/validator_result
```

来源粒度：

- `source_document`：文件级，只回答“来自哪个 PDF”。
- `source_page`：页面级，只回答“来自第几页”。
- `source_block`：版面块级，回答“来自页面上的哪个区域”。
- `source_fragment`：语义片段级，回答“题目 1、某个知识条、答案解析来自哪段原文”。
- `question_source_field_rel`：字段取值级，回答“题干、A 选项、答案、解析、知识点正文分别取自哪个片段”。

## 4. 统一导入流程

```text
PDF 文件
  -> question.ai_pdf_text
  -> 保存 raw_text/raw_html/source_chars/调用参数
  -> 渲染页图用于原文定位和视觉还原
  -> 大模型读取 raw_text 并返回统一 JSON 结构
  -> 按大模型 JSON 保存来源片段和字段取值来源
  -> JSON 外形和必填字段校验
  -> 失败项落库
  -> 人工预览确认
  -> 业务表落库
```

教师一键导入流程：

```text
教师打开知识点维护页
  -> 点击 PDF导入知识点
  -> 选择学科/年级/教材/单元，默认带入当前筛选条件
  -> 上传 PDF
  -> 点击 一键生成知识点
  -> 后端调用 question.ai_pdf_text
  -> 大模型返回 Unit/知识点/正文/题目草稿/失败项 JSON
  -> 后端校验 JSON 外形、保存来源、落库高置信知识点
  -> 返回导入结果、来源片段、失败复核项
  -> 知识点列表自动刷新
```

一键导入不是无审计导入：

- `auto_commit=true` 时，只允许高置信、来源完整、字段来源可追溯的知识点自动进入正式知识库。
- `auto_commit=false` 时，只生成预览批次，教师确认后再提交。
- 来源缺失、跨页边界不清、AI 推断、字段无法追溯的内容进入 `question_pdf_parse_issue`，状态为 `pending_review`。
- 混合 PDF 中识别出的真实题目默认进入题目草稿或待确认队列，不随知识点一键直接进正式题库。

### 4.1 `question.ai_pdf_text`

输入：

- `file_path`
- `max_chars=120000`
- 上传文件时使用 upload 调用同一服务。

输出必须保存：

- `raw_text`
- `raw_html`
- `file_name`
- `source_chars`
- `max_chars`
- `extract_service=question.ai_pdf_text`
- `extract_params_json`
- `raw_text_sha256`
- `raw_html_sha256`

任何 PDF 解析服务都不得绕过它直接读取 PDF 或直接调腾讯 OCR。

取值来源要求：

- `question.ai_pdf_text.raw_text/raw_html` 是唯一原始文本底稿。
- `source_block.raw_text` 从原始底稿切出，不能由 AI 改写。
- `source_fragment.raw_text` 从 `source_block.raw_text` 切出，保存字符区间和上下文。
- 业务字段的最终值可以规范化，但必须在 `question_source_field_rel` 记录原始命中片段。
- 如果一个字段由多个片段合并得到，保存多条字段来源关系，并记录 `field_part_order`。

### 4.2 页图

页图只负责：

- 原文定位。
- 图片资源裁切。
- 视觉还原 PDF。
- 失败项截图证据。

页图不是文本主入口。

## 5. 统一语义类型

AI 输出的 `semantic_item.semantic_type` 使用统一枚举：

知识类：

- `letter_point`：字母、发音、例词；不直接入题库。
- `vocabulary_point`：单词和释义；不直接入题库。
- `phrase_point`：短语和释义；不直接入题库。
- `sentence_pattern`：句型结构和例句；不直接入题库。
- `daily_expression`：常用表达、交际用语；不直接入题库。
- `grammar_point`：语法规则和例句；不直接入题库。
- `knowledge_summary`：综合知识总结；不直接入题库。

题目类：

- `normal_question`：普通题目；校验后可保存为草稿或导入题目。
- `reading_choice`：阅读选择父题；校验后可保存为草稿或导入题目。
- `cloze_choice`：完形父题；校验后可保存为草稿或导入题目。
- `candidate_question`：大模型返回的候选题；不直接入题库，必须先经过人工复核或 JSON 字段完整性校验。

资源类：

- `illustration`：插画或版式图片；不直接入题库，作为知识点或题目资源关联。

## 6. 题目判断规则

### 6.1 不是题目的内容

这类编号内容是知识点，不是题目：

```text
1、句型结构: Hello/Hi! I'm/My name is+姓名。
2、句型结构: Nice to meet you.
3、Hey, Sarah! We can share.
```

落库：

- `question_knowledge_content.semantic_type=sentence_pattern`
- `question_item` 不新增
- `question_source_snapshot` 保存来源原文

### 6.2 是题目的内容

满足以下条件才可进入题库：

- 有明确作答任务、题干、空格、选项或问句。
- 选择题有可分辨选项。
- 答案来自原文答案区或明确解析区，不能由 AI 自行解题。
- 若答案缺失，题目可以入草稿，但 `answer_key` 留空，状态标记为待补全。

### 6.3 `9英语易错题` 类型

`9英语易错题-秒懂比较级句型.pdf` 应拆为：

- 知识点：比较级常用句型、同级比较句型、比较级表达最高级。
- 题目：6 道选择题。
- 答案解析：从原文每日易错题区回填。
- 来源原文：每道题保存题干、选项、答案解析附近原文。

## 7. 失败项设计

AI 失败项必须落库，不只写日志。

`question_pdf_parse_issue`

- `issue_id`
- `source_doc_id`
- `page_no`
- `source_block_id`
- `issue_type`
- `severity`: `info`、`warning`、`error`
- `raw_text`
- `crop_image_url`
- `ai_output_json`
- `expected_schema`
- `error_msg`
- `suggestion`
- `status`: `pending`、`fixed`、`ignored`、`reparsed`
- `create_time`
- `modify_time`

典型失败项：

- `not_a_question`：AI 把知识点编号误判成题。
  处理：不入题库，映射为知识点。
- `missing_answer`：题目存在但答案缺失。
  处理：入草稿或失败队列，答案留空。
- `text_extract_error`：`We can` 识别成 `Wecan` 等。
  处理：进入修正项。
- `cross_unit_boundary`：一页跨两个 Unit。
  处理：切块后分别落单元。
- `duplicate_content`：同内容重复出现。
  处理：按 hash 去重。
- `low_confidence_question`：题目边界不稳定。
  处理：待人工确认。

“失败题目”也要保留原文：

```json
{
  "candidate_id": "docsha-p001-cq001",
  "raw_text": "1、句型结构: Hello/Hi! I'm/My name is+姓名。",
  "ai_question_type": "blank",
  "final_status": "rejected",
  "reject_reason": "编号考点，不是可作答题目",
  "mapped_knowledge_id": "u1_sentence_pattern"
}
```

## 8. 统一数据表

### 8.1 来源文档

`question_source_document`

- `source_doc_id`
- `import_batch_id`
- `file_name`
- `file_sha256`
- `file_url`
- `page_count`
- `subject`
- `stage`
- `grade`
- `textbook_version`
- `parse_status`
- `import_status`
- `create_time`

### 8.2 来源页面

`question_source_page`

- `source_page_id`
- `source_doc_id`
- `page_no`
- `page_image_url`
- `width`
- `height`
- `extract_service`
- `extract_params_json`
- `raw_text`
- `raw_html`
- `extract_meta_json`
- `page_hash`

### 8.3 来源块

`question_source_block`

- `source_block_id`
- `source_doc_id`
- `page_no`
- `block_order`
- `bbox_json`
- `block_type`
- `raw_text`
- `normalized_text`
- `block_image_url`
- `semantic_type`
- `confidence`
- `content_hash`

### 8.4 来源片段

`question_source_fragment`

用于保存可审计的最小语义片段。片段不是整篇文档，也不是整页文本，而是题目、答案解析、知识条、例句、对话、小练习等可独立比对的原文范围。

- `source_fragment_id`
- `source_doc_id`
- `source_page_id`
- `source_block_id`
- `page_no`
- `fragment_order`
- `fragment_type`: `knowledge_point`、`example_sentence`、`dialogue`、`candidate_question`、`question_stem`、`question_option`、`answer`、`analysis`、`mixed`、`unknown`
- `raw_text`
- `raw_html`
- `normalized_text`
- `char_start`
- `char_end`
- `context_before`
- `context_after`
- `bbox_json`
- `fragment_hash`
- `confidence`
- `status`
- `create_time`

规则：

- `raw_text` 保存命中的原始片段。
- `context_before/context_after` 只保存相邻正文窗口，不保存整篇 PDF，默认各 200 到 500 字。
- `char_start/char_end` 相对于 `question.ai_pdf_text.raw_text` 或页面 raw_text 的稳定字符位置。
- `bbox_json` 有坐标时必须保存，方便在页图上高亮。
- 同一题目可以关联多个片段，例如题干片段、选项片段、答案片段、解析片段。
- 无法确定片段边界时，先落 `unknown` 或 `mixed`，并生成失败项等待复核。

示例：题目 1 的来源不是“整篇 PDF”，而是片段集合：

```json
{
  "question_no": "1",
  "source_file_name": "9英语易错题-秒懂比较级句型.pdf",
  "stem_fragment_id": "docsha-p003-b002-f001",
  "option_fragment_id": "docsha-p003-b002-f002",
  "answer_fragment_id": "docsha-p004-b001-f001",
  "analysis_fragment_id": "docsha-p004-b001-f002"
}
```

### 8.5 来源快照

`question_source_snapshot`

用于保存知识点和题目的原始信息。

- `snapshot_id`
- `source_doc_id`
- `source_page_id`
- `source_block_id`
- `source_fragment_id`
- `question_id`
- `knowledge_id`
- `content_id`
- `extract_service`
- `extract_params_json`
- `raw_text`
- `raw_html`
- `normalized_text`
- `ai_output_json`
- `validator_result_json`
- `status`
- `create_time`

规则：

- `raw_text/raw_html` 永不被 AI 清洗覆盖。
- 结构化题目字段写入 `question_item` 等现有题库表。
- 结构化知识字段写入 `question_knowledge_content`。
- 原始上下文统一从 `question_source_snapshot` 回看。
- 需要比对具体字段时，统一查 `question_source_field_rel`，不要只看整条 snapshot。

### 8.6 知识正文

`question_knowledge_content`

- `content_id`
- `batch_id`
- `source_doc_id`
- `source_block_id`
- `unit_id`
- `knowledge_id`
- `semantic_type`
- `section_title`
- `content_text`
- `content_html`
- `content_json`
- `content_hash`
- `asset_count`
- `order_index`
- `status`
- `create_time`
- `modify_time`

### 8.7 题目来源关系

`question_source_rel`

- `rel_id`
- `question_id`
- `source_doc_id`
- `source_page_no`
- `source_block_id`
- `source_fragment_id`
- `source_content_id`
- `knowledge_id`
- `relation_type`: `generated_from`、`imported_from`、`manually_linked`
- `confidence`
- `create_time`

规则：

- `question_source_rel` 是题目整体来源关系。
- 题目整体来源不能替代字段级来源。
- 如果题干、选项、答案、解析来自不同片段，必须同时写多条 `question_source_field_rel`。

### 8.8 字段级来源关系

`question_source_field_rel`

用于记录每个结构化字段的取值来自哪个 PDF 片段。它是人工比对和验收的核心表。

- `field_rel_id`
- `source_doc_id`
- `source_page_id`
- `source_block_id`
- `source_fragment_id`
- `entity_type`: `knowledge_content`、`question_item`、`question_option`、`question_answer`、`question_analysis`、`parse_issue`
- `entity_id`
- `field_name`: `content_text`、`stem`、`option_a`、`option_b`、`option_c`、`option_d`、`answer_key`、`analysis_text`
- `field_part_order`
- `extracted_value`
- `normalized_value`
- `raw_quote`
- `context_before`
- `context_after`
- `char_start`
- `char_end`
- `bbox_json`
- `confidence`
- `match_status`: `matched`、`partial`、`missing`、`inferred`、`conflict`
- `review_status`: `pending`、`accepted`、`fixed`、`ignored`
- `create_time`

规则：

- `raw_quote` 必须是原始 PDF 片段中的连续文本，不能写 AI 改写后的内容。
- `extracted_value` 是从原文抽取出来的候选值。
- `normalized_value` 是写入业务表前的规范化值。
- `match_status=inferred` 代表 AI 推断，默认不允许正式落库。
- `match_status=missing/conflict` 必须生成 `question_pdf_parse_issue`。
- 一个业务字段来自多个片段时，用多条记录和 `field_part_order` 表达合并顺序。
- 删除或重新解析业务字段时，不删除旧来源，改写 `review_status` 或生成新版本。

题目 1 字段来源示例：

```json
{
  "question_id": "q_001",
  "source_file_name": "9英语易错题-秒懂比较级句型.pdf",
  "fields": [
    {
      "field_name": "stem",
      "source_fragment_id": "docsha-p003-b002-f001",
      "raw_quote": "1. ...",
      "context_before": "比较级常用句型 ...",
      "context_after": "A. ... B. ... C. ..."
    },
    {
      "field_name": "answer_key",
      "source_fragment_id": "docsha-p004-b001-f001",
      "raw_quote": "答案：C",
      "context_before": "每日易错题答案及解析",
      "context_after": "解析：..."
    },
    {
      "field_name": "analysis_text",
      "source_fragment_id": "docsha-p004-b001-f002",
      "raw_quote": "解析：...",
      "context_before": "答案：C",
      "context_after": "2. ..."
    }
  ]
}
```

## 9. 单元与知识点统一编码

三年级上册 PEP 英语统一上下文：

- `subject=english`
- `stage=primary`
- `grade=grade_3`
- `textbook_version=pep`

单元：

- Unit 1
  编码：`unit-english-grade3-pep-u1`
  名称：`Unit 1 Making friends`
- Unit 2
  编码：`unit-english-grade3-pep-u2`
  名称：`Unit 2 Different families`
- Unit 3
  编码：`unit-english-grade3-pep-u3`
  名称：`Unit 3 Amazing animals`
- Unit 4
  编码：`unit-english-grade3-pep-u4`
  名称：`Unit 4 Plants around us`
- Unit 5
  编码：`unit-english-grade3-pep-u5`
  名称：`Unit 5 The colourful world`
- Unit 6
  编码：`unit-english-grade3-pep-u6`
  名称：`Unit 6 Useful numbers`

知识点编码：

- 全册考点汇总：`u1_letter_point`、`u1_vocabulary_point`、`u1_phrase_point`、`u1_sentence_pattern`、`u1_daily_expression`、`u1_grammar_point`
- Part 型知识总结：`u1_part_a_four_skill_words`、`u1_part_a_sentence_pattern`、`u1_part_b_reading`、`u1_part_c_project`
- 易错题专题：`comparative_sentence_pattern`、`comparative_common_mistakes`

## 10. 知识点界面与教师导入入口

`question-taxonomy` 的知识点维护页必须能看到来源原文。

教师导入入口：

- 不新增单独菜单，不新增独立页面。
- 固定入口页面：`http://192.168.232.130:8026/collect-ui/#/collect-ui/framework/question-taxonomy`。
- 对应前端路由：`/framework/question-taxonomy`。
- 页面位置：右侧 `知识点维护` 面板顶部按钮区。
- 按钮位置：与现有 `新增知识点` 同一组 `topRight` 操作区，放在 `新增知识点` 左侧或右侧。
- 按钮文案：`PDF导入知识点`。
- 按钮图标：`CloudUploadOutlined` 或 `FilePdfOutlined`。
- 打开方式：弹窗或抽屉 `PDF知识点导入`，不跳出当前知识点维护页。

入口边界：

- 只能在现有 `question-taxonomy` 页面上加入口。
- 不在菜单中新增 `PDF导入`、`知识点导入` 或其他独立菜单项。
- `question-bank` 已有题目 AI 导入，适合题目 PDF。
- 知识点 PDF 的目标是生成 `question_unit/question_knowledge/question_knowledge_content`，主入口必须在 `question-taxonomy`。
- 混合 PDF 解析出的题目可以在导入结果里给出 `查看题目草稿` 或 `转到题库导入结果`，但不替代知识点入口。

教师弹窗字段：

- `PDF文件`：上传或选择已上传文件。
- `学科`：默认取当前筛选，示例 `english`。
- `年级`：默认取当前筛选，示例 `grade_3`。
- `教材版本`：默认 `pep`。
- 不提供 `目标单元` 输入；Unit 必须由大模型根据 PDF 内容返回。
- `导入模式`：`一键生成并入库`、`只生成预览`。
- `题目处理`：`只生成知识点`、`题目进草稿`。
- `来源保留`：默认开启，不允许关闭。

教师操作按钮：

- `一键生成知识点`：调用统一导入服务，自动生成知识点。
- `只预览`：只生成批次和来源片段，不写正式知识库。
- `查看失败复核`：打开当前导入批次的失败项列表。
- `查看来源对比`：打开来源片段和整理后知识点的左右对比。

导入完成后展示：

- 文件名、页数、抽取字符数。
- 生成单元数、知识点数、知识正文数。
- 自动入库数量、待复核数量、跳过重复数量。
- 每个知识点的来源页码、来源片段、字段来源状态。
- 如果混合 PDF 识别出题目，显示题目草稿数量和待补全数量。
- 提供 `刷新知识点列表`、`按本文件过滤来源`、`下载验收报告`。

来源查询能力：

- 支持按 `file_name`、`source_doc_id`、导入批次、学科、年级、教材版本过滤。
- 支持按 Unit、知识点、题目、失败项、字段名过滤。
- 搜索结果展示到片段级，不展示整篇文档全文。
- 每条结果必须显示文件名、页码、块号、片段号、片段类型、匹配字段、匹配状态。
- 点击片段时展示页图高亮、命中原文、前文、后文、AI 输出和校验结果。

知识点详情增加 Tab：

- `整理后内容`：展示 `question_knowledge_content.content_html`
- `来源原文`：展示关联 `question_source_fragment.raw_html/raw_text`，提供字段来源和相邻上下文
- `字段来源`：展示 `content_text/content_json` 各字段取自哪个 PDF 片段
- `题目关联`：展示该知识点关联的题目、状态、答案完整度
- `解析失败项`：展示待处理失败项，可定位页图和原文块

题目详情也要显示：

- 来源 PDF 文件名。
- 页码/块号。
- 来源片段号。
- 原始题干上下文。
- 选项、答案、解析各自的来源片段。
- 每个结构化字段的 `raw_quote`、`context_before`、`context_after`。
- AI 输出 JSON。
- 校验结果。

题目来源对比视图：

- 左侧：PDF 片段预览，只显示命中片段和相邻正文，不显示整篇文档。
- 右侧：结构化题目字段，包括题干、选项、答案、解析、知识点关联。
- 高亮：字段值对应的 `raw_quote`。
- 差异：展示 `normalized_value` 与 `raw_quote` 的差异，例如 OCR 空格、标点、大小写修正。
- 操作：接受、修正字段、标记 AI 推断、加入失败项、重新解析当前片段。

示例查询：

```json
{
  "service": "question.pdf_source_trace_query",
  "file_name": "9英语易错题-秒懂比较级句型.pdf",
  "entity_type": "question_item",
  "question_no": "1",
  "include_context": true,
  "context_chars": 300
}
```

示例返回：

```json
{
  "file_name": "9英语易错题-秒懂比较级句型.pdf",
  "question_id": "q_001",
  "source_fragments": [
    {
      "field_name": "stem",
      "page_no": 3,
      "source_block_id": "docsha-p003-b002",
      "source_fragment_id": "docsha-p003-b002-f001",
      "context_before": "比较级常用句型 ...",
      "raw_quote": "1. ...",
      "context_after": "A. ... B. ... C. ..."
    }
  ]
}
```

## 11. 双向还原

### 11.1 原文件还原

保存原 PDF 或远程 URL，并记录 `file_sha256`。导出时直接下载原文件，做到字节级可逆。

### 11.2 视觉还原

保存页图和页面尺寸。导出时按页图重新封装 PDF，做到视觉基本一致。

### 11.3 结构还原

根据结构化数据导出可编辑 PDF：

```text
Unit
  -> 知识点
    -> 整理后内容
    -> 来源图片
    -> 关联题目
    -> 答案解析
```

结构还原不保证和原 PDF 完全同版式，但必须能再导入并比对结构 hash。

## 12. Round Trip 校验

```text
原始 PDF
  -> question.ai_pdf_text
  -> AI 解析
  -> 落库
  -> 导出结构 PDF
  -> 再导入
  -> 比对 hash
```

比对项：

- `source_doc_id/file_sha256`
- Unit 编码和数量。
- 知识点编码和数量。
- `content_hash`
- 题目 `content_hash` 或题干/选项/答案 hash。
- 图片 sha256。
- 失败项状态。
- `question_source_snapshot.raw_text` 不丢失。
- `question_source_fragment.fragment_hash` 不丢失。
- `question_source_field_rel.raw_quote/normalized_value` 对应关系不丢失。

## 13. 服务设计

新增统一服务：

- `question.pdf_knowledge_import_one_click`
  教师端一键导入入口服务。接收 PDF、学科、年级、教材和导入模式，不接收目标单元；Unit 必须由大模型从 PDF 内容中返回。服务内部调用 `question.ai_pdf_text` 和大模型结构化提示词，只做 JSON 外形校验、来源保存、落库和验收报告生成。
- `question.pdf_import_preview`
  调用 `question.ai_pdf_text`，生成统一预览。
- `question.pdf_import_commit`
  将确认后的知识点、题目、失败项落库。
- `question.pdf_source_query`
  查询来源文档、页面、块、片段，支持按文件名和导入批次过滤。
- `question.pdf_snapshot_query`
  查询知识点或题目的来源原文。
- `question.pdf_source_trace_query`
  查询题目、知识点或具体字段的取值来源，返回命中片段和前后正文。
- `question.pdf_field_source_save`
  保存人工修正后的字段来源关系。
- `question.pdf_parse_issue_query`
  查询解析失败项。
- `question.pdf_export_original`
  导出原 PDF。
- `question.pdf_export_visual`
  导出页图还原 PDF。
- `question.pdf_export_structured`
  导出结构化 PDF。

`question.pdf_knowledge_import_one_click` 输入：

```json
{
  "file_path": "/data/project/ai-study/题目/三年级上册英语全册知识考点汇总人教版.pdf",
  "subject": "english",
  "stage": "primary",
  "grade": "grade_3",
  "textbook_version": "pep",
  "unit_id": "",
  "auto_commit": true,
  "question_policy": "draft",
  "max_chars": 120000
}
```

`question.pdf_knowledge_import_one_click` 输出：

```json
{
  "import_batch_id": "...",
  "source_doc_id": "...",
  "file_name": "...",
  "unit_total": 6,
  "knowledge_total": 42,
  "content_total": 126,
  "auto_committed_total": 120,
  "pending_review_total": 6,
  "question_draft_total": 0,
  "acceptance_status": "warning"
}
```

保留现有服务：

- `question.ai_pdf_text`：唯一 PDF 文本入口。
- `question.ai_parse`：题目型 PDF 的题目 JSON 解析能力，可作为统一服务内部的题目解析子能力。
- `question.question_choice_save`：保存结构化题目。
- `question.knowledge_*`：保存知识点。

## 14. 实施分期

M1：统一来源层

- 新增 `question_source_document/page/block/fragment/snapshot/field_rel`。
- `question.pdf_import_preview` 必须先调 `question.ai_pdf_text`。
- 保存 `raw_text/raw_html`。
- 支持按文件名查询来源文档、片段和字段来源。

M2：知识点导入

- 在现有页面 `http://192.168.232.130:8026/collect-ui/#/collect-ui/framework/question-taxonomy` 的 `知识点维护` 顶部操作区新增 `PDF导入知识点`。
- 不新增单独菜单，不新增独立导入页面。
- 新增教师弹窗 `PDF知识点导入`，支持上传 PDF 和一键生成知识点。
- 新增 `question.pdf_knowledge_import_one_click`，串联读取、解析、校验、落库和报告。
- 三年级上册 Unit 和知识点 upsert。
- 保存 `question_knowledge_content`。
- 知识点页面显示来源原文。
- 知识点页面显示 `content_text/content_json` 的字段级来源。
- 导入成功后自动刷新知识点列表，并支持按本次文件名过滤来源。

M3：题目导入

- 支持易错题 PDF 拆出题目。
- 保存题目、答案、解析、知识点关系、来源关系。
- 保存题干、选项、答案、解析的字段级来源关系。
- 缺答案题和 AI 失败题进入失败项。

M4：可逆导出

- 原 PDF 下载。
- 视觉 PDF。
- 结构化 PDF。
- Round Trip 比对。

## 15. 验证与测试验收计划

测试验收执行顺序：

- 服务验收：先验证 `question.ai_pdf_text` 和 `question.pdf_knowledge_import_one_click` 能正确返回原文、批次、生成数量和失败项。
- 数据验收：再验证来源文档、来源片段、字段级来源、知识点正文和失败项是否完整落库。
- 页面验收：最后在 `question-taxonomy` 现有页面完成教师导入流程，采集页面错误、接口错误和页面截图证据。
- 回归验收：再次打开知识点维护页，按文件名过滤本次导入结果，并检查来源对比、失败复核、知识点列表刷新是否稳定。

文本入口验证：

```json
{
  "service": "question.ai_pdf_text",
  "file_path": "/data/project/ai-study/题目/9英语易错题-秒懂比较级句型.pdf",
  "max_chars": 120000
}
```

PDF 与知识点一致性验收：

验收目标不是“PDF 读成功”就结束，而是确认 `question.ai_pdf_text` 读出的原文、AI 解析结果、落库后的知识点三者一致。

比对链路：

```text
question.ai_pdf_text.raw_text/raw_html
  -> question_source_block
  -> question_source_fragment
  -> semantic_item
  -> question_knowledge_content.content_text/content_json
  -> question_source_snapshot.raw_text/raw_html
  -> question_source_field_rel.raw_quote/normalized_value
```

验收维度：

- `来源覆盖`：PDF 中识别出的 Unit、章节标题、知识栏目都能在知识点树中找到对应节点。
- `内容一致`：知识点 `content_text/content_json` 中的单词、短语、句型、语法规则来自 PDF 原文，不能凭空新增。
- `原文可追溯`：每个知识点内容都能反查 `question_source_snapshot.raw_text/raw_html`。
- `页块可定位`：每个知识点内容能定位到 `source_doc_id/source_page_id/source_block_id`。
- `片段可定位`：每个知识点、题目、答案解析都能定位到 `source_fragment_id`。
- `字段可追溯`：题干、选项、答案、解析、知识点正文必须能通过 `question_source_field_rel` 找到 `raw_quote`。
- `文件可过滤`：界面按 PDF 文件名过滤后，只展示该文件对应的片段、题目、知识点和失败项。
- `上下文可比对`：界面能展示命中片段的 `context_before/raw_quote/context_after`，不要求打开整篇文档才能验证。
- `题目不误入`：编号考点不得进入 `question_item`；真正题目必须有题干、选项、答案解析来源。
- `缺失可见`：PDF 中读到了但无法映射到知识点或题目的内容，必须进入 `question_pdf_parse_issue`，不能静默丢失。
- `重复可解释`：重复标题、重复知识块通过 hash 去重，并在验收报告里说明保留哪一版。

样本验收标准：

- 知识点总结型
  必须一致：Unit、Part、四会/三会单词、句型、Letters and sounds、阅读/项目活动。
  不应出现：不应自动生成题目。
- 全册考点汇总型
  必须一致：Unit 1-6、字母/单词/短语/句型/常用表达/语法考点。
  不应出现：编号考点不应进题库。
- 知识点 + 题目混合型
  必须一致：比较级知识点、6 道题、选项、答案、解析、来源片段、字段来源。
  不应出现：答案不能由 AI 自行推断。

建议生成验收报告 `question.pdf_import_acceptance_report`：

```json
{
  "source_doc_id": "...",
  "file_name": "...",
  "raw_text_chars": 2684,
  "unit_total": 1,
  "knowledge_total": 2,
  "question_total": 6,
  "matched_block_total": 12,
  "matched_fragment_total": 18,
  "field_source_total": 32,
  "unmatched_block_total": 0,
  "unmatched_field_total": 0,
  "issue_total": 0,
  "checks": [
    {"key": "source_coverage", "status": "pass", "message": "所有来源块已映射"},
    {"key": "knowledge_consistency", "status": "pass", "message": "知识点内容均可追溯到 PDF 原文"},
    {"key": "question_consistency", "status": "pass", "message": "题目答案解析均来自原文"},
    {"key": "field_traceability", "status": "pass", "message": "题干、选项、答案、解析均有片段级来源"}
  ]
}
```

验收失败处理：

- `warning`：允许提交，但必须展示在导入结果里，例如少量文本抽取错字。
- `error`：禁止确认落库，例如 Unit 缺失、知识点无来源原文、题目答案由 AI 推断、题目字段没有片段来源。
- `pending_review`：允许保存为待处理批次，不进入正式知识库或题库。

数据验证：

```sql
select count(*) from question_source_snapshot
where source_doc_id = ? and length(raw_text) > 0;

select count(*) from question_knowledge
where subject='english' and stage='primary' and grade='grade_3';

select count(*) from question_item
where source in ('pdf_import', 'knowledge_pdf_generated');

select count(*) from question_pdf_parse_issue
where source_doc_id = ? and status = 'pending';

select count(*) from question_knowledge_content c
left join question_source_snapshot s on s.content_id = c.content_id
where c.source_doc_id = ? and (s.snapshot_id is null or length(s.raw_text) = 0);

select count(*) from question_source_block b
left join question_source_snapshot s on s.source_block_id = b.source_block_id
left join question_pdf_parse_issue i on i.source_block_id = b.source_block_id
where b.source_doc_id = ? and s.snapshot_id is null and i.issue_id is null;

select count(*) from question_source_fragment
where source_doc_id = ? and length(raw_text) > 0;

select count(*) from question_source_rel sr
left join question_source_field_rel r
  on r.entity_type = 'question_item'
 and r.entity_id = sr.question_id
 and r.field_name = 'stem'
where sr.source_doc_id = ? and r.field_rel_id is null;

select count(*) from question_source_field_rel
where source_doc_id = ? and match_status in ('missing', 'inferred', 'conflict');

select count(*) from question_source_field_rel r
join question_source_document d on d.source_doc_id = r.source_doc_id
where d.file_name = ?;
```

页面验收：

验收页面：

`http://192.168.232.130:8026/collect-ui/#/collect-ui/framework/question-taxonomy`

入口验收：

- 不出现新的知识点导入菜单。
- 不出现独立的知识点导入页面。
- `PDF导入知识点` 只出现在 `question-taxonomy` 的 `知识点维护` 顶部操作区。
- `PDF导入知识点` 与 `新增知识点` 属于同一组操作。
- 点击 `PDF导入知识点` 后打开 `PDF知识点导入` 弹窗或抽屉，不离开当前页面。

导入弹窗验收：

- 能上传 PDF。
- 能显示并修改学科、年级、教材版本。
- 弹窗不出现目标单元字段，导入请求不发送 `unit_id`。
- 默认值能从当前知识点筛选条件带入。
- 能选择 `一键生成并入库` 或 `只生成预览`。
- 能选择题目处理策略：`只生成知识点` 或 `题目进草稿`。
- `来源保留` 默认开启，且不能关闭。
- 点击 `一键生成知识点` 后出现处理中状态，避免教师重复提交。

导入结果验收：

- 页面展示文件名、页数、抽取字符数。
- 页面展示生成单元数、知识点数、知识正文数。
- 页面展示自动入库数量、待复核数量、跳过重复数量。
- 知识点列表自动刷新。
- 能点击 `按本文件过滤来源`，只显示当前 PDF 生成的知识点和来源片段。
- 能点击 `下载验收报告` 或查看验收结果。
- 混合 PDF 识别出题目时，只展示题目草稿或待确认数量，不把题目静默写入正式题库。

来源对比验收：

- `question-taxonomy` 能看到知识点整理后内容。
- `question-taxonomy` 能看到来源原文。
- `question-taxonomy` 能按 PDF 文件名过滤知识点来源。
- `question-taxonomy` 能展示知识点字段来源片段和前后正文。
- `question-taxonomy` 能看到 PDF 与知识点一致性验收结果。
- 点击知识点来源后，左侧显示 PDF 原始片段，右侧显示整理后的知识点内容。
- 原始片段必须包含 `context_before/raw_quote/context_after`。
- 字段来源能标出 `content_text/content_json` 每个取值来自哪个片段。
- 题目关联列表能显示来源题干片段。
- 题目详情能展示题干、选项、答案、解析分别来自哪个 PDF 片段。
- 题目详情能左右对比结构化字段和 `context_before/raw_quote/context_after`。
- 失败项能定位到原文块、来源片段和页图。

页面错误验收：

- 浏览器控制台不能出现新的 `console.error`。
- 页面不能出现 `pageerror`。
- `question.ai_pdf_text`、`question.pdf_knowledge_import_one_click`、来源查询和失败项查询接口不能失败。
- PDF 上传、导入、来源过滤、来源对比、失败复核过程中不能出现关键资源加载失败。
- 移动或窄屏下弹窗内容不能互相遮挡，按钮文字不能溢出。

可逆验证：

- 原 PDF 导出 sha256 等于导入时 `file_sha256`。
- 视觉 PDF 页数一致。
- 结构化 PDF 再导入后 Unit、知识点、题目、来源快照、来源片段、字段来源 hash 一致。
