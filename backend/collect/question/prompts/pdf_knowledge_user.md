请把下面 PDF 抽取文本转换成统一 JSON。

硬性要求：

- 只返回合法 JSON 对象，不能有 markdown。
- 必须返回 units、questions、issues、question_draft_total、acceptance_status、summary。
- units 必须形成树：Unit -> knowledge -> contents。
- 不要把标题孤立成知识点；知识点必须有可展示的 content_text。
- 必须导入整份 PDF 中所有可展示知识点；不要只返回摘要、不要只返回前几条、不要把多个知识点压缩成一条。
- 对列表或表格，逐行识别：每个词汇、短语、句型、语法点、易错点都要进入 knowledge/contents。
- 不要把知识点编号误当题目。是否是题目由你在 JSON 中判断，后端不会再做代码判断。
- 不要把答案、解析、题目归属交给后端处理。答案和解析必须由你在 questions 中直接返回。
- 如果来源不确定，放入 issues，不要让后端靠关键词兜底。
- 每个 content_text、题干、选项、答案、解析都尽量返回 source_quote 或 source_quotes。
- source_quote 可以是最接近的原文行或标题锚点；不要因为 source_quote 不够长而漏掉知识点。

推荐结构：

{
  "units": [
    {
      "unit_id": "",
      "unit_code": "unit_1",
      "unit_name": "Unit 1 ...",
      "subject": "english",
      "stage": "primary",
      "grade": "grade_3",
      "textbook_version": "pep",
      "order_index": 1,
      "source_quote": "原文中的单元标题或锚点",
      "knowledge": [
        {
          "knowledge_id": "",
          "knowledge_code": "stable_code",
          "knowledge_name": "知识点名称",
          "semantic_type": "sentence_pattern",
          "order_index": 1,
          "source_quote": "原文中的知识点锚点",
          "contents": [
            {
              "section_title": "正文小标题",
              "semantic_type": "sentence_pattern",
              "content_text": "整理后的知识正文",
              "content_html": "",
              "source_quote": "raw_text 中对应原文片段",
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
  "questions": [],
  "issues": [],
  "question_draft_total": 0,
  "acceptance_status": "pass",
  "summary": "解析摘要"
}

PDF 抽取文本：

{{.RawText}}
