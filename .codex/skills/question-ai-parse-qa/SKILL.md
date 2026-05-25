---
name: question-ai-parse-qa
description: "Use for AI Study backend question.ai_parse QA in /data/project/ai-study: adding or splitting question/PDF text fixtures, comparing Codex CLI and DeepSeek parse outputs, tuning backend/collect/question/prompts, and rerunning verification reports."
---

# Question AI Parse QA

## Scope

Use this skill for AI Study question import parsing work, especially when the user provides new题目/PDF text fragments and wants Codex CLI and DeepSeek results to match.

Primary files:

- Prompt templates: `backend/collect/question/prompts/ai_parse_system.md`, `backend/collect/question/prompts/ai_parse_user.md`
- Compare script: `backend/scripts/compare_question_ai_parse.py`
- Manual fixtures: `backend/testdata/question_ai_fragments/`
- Real PDF fixtures: `backend/testdata/question_ai_pdf_fragments/`
- Reports: `reports/`

Do not write API keys into files or logs. Use `DEEPSEEK_API_KEY` from the environment.

## Fixture Workflow

1. Put each new题目片段 in a small `.txt` file under the relevant `backend/testdata/...` directory.
2. Add or update the directory `manifest.json`.
3. Keep manifest expectations simple unless the user asks otherwise:
   - `expected_count`
   - existing `expected_choice_items` only when the first output row is a grouped reading/cloze parent
   - `params` with intended default `question_type` and `question_category`
4. Do not add auxiliary expected-answer/category JSON fields unless explicitly requested.
5. If a fixture has a truncated answer section, first look for the complete answer section in the source text/PDF extraction and complete the fixture. Prefer fixing fixture completeness over forcing prompts to handle incomplete evidence.

Useful source search:

```bash
rg -n "关键题干|答案及解析|参考答案" reports/*.txt backend/testdata
```

## Backend Prerequisites

Confirm the backend is running before comparisons:

```bash
cd /data/project/ai-study/backend
ss -ltnp | rg ':8026'
curl --noproxy '*' -sS -m 5 -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8026/collect-ui/
```

If it is down or requests start closing unexpectedly:

```bash
cd /data/project/ai-study/backend
./shutdown.sh
./linux-start-dev.sh
```

Use project scripts, not ad hoc long-running `go run ... &`.

## Compare Commands

Run focused high-risk cases first:

```bash
cd /data/project/ai-study
python3 backend/scripts/compare_question_ai_parse.py \
  --manifest backend/testdata/question_ai_pdf_fragments/manifest_selected.json \
  --providers codex,deepseek \
  --repeat 1 \
  --out reports/question_ai_pdf_selected_codex_deepseek_compare.json \
  --summary
```

Run one or more cases with `--case-id`:

```bash
python3 backend/scripts/compare_question_ai_parse.py \
  --manifest backend/testdata/question_ai_pdf_fragments/manifest.json \
  --case-id pdf_01_header_single_choice,pdf_02_cloze_choice \
  --providers codex,deepseek \
  --out reports/question_ai_pdf_batch_compare.json \
  --summary
```

For full validation, prefer batches if the backend or provider is unstable:

- Batch 1: `pdf_01` to `pdf_05`
- Batch 2: `pdf_06` to `pdf_10`

The target is `OK` for both providers and `diff: same`. The diff must cover question structure, answers, and analysis/explanation text (`analysis_text` and grouped `choice_items.analysis`).
Analysis diff is normalized for whitespace, punctuation, and obvious OCR noise; do not treat exact prose formatting as a failure when both providers extracted the same explanation.

## Diff Triage

When a report says `diff: different`, inspect row summaries:

```bash
python3 - <<'PY'
import json
path = "reports/your_report.json"
data = json.load(open(path, encoding="utf-8"))
for case in data["cases"]:
    run = case["runs"][0]
    print("\\nCASE", case["id"])
    c = run["providers"].get("codex", {}).get("rows", [])
    d = run["providers"].get("deepseek", {}).get("rows", [])
    for i, (a, b) in enumerate(zip(c, d), 1):
        diffs = [k for k in ["question_type", "question_category", "answer_key", "analysis_text", "choice_item_count", "choice_item_answers", "choice_item_analyses"] if a.get(k) != b.get(k)]
        if diffs:
            print(i, diffs, {k: a.get(k) for k in diffs}, {k: b.get(k) for k in diffs})
PY
```

Common fixes:

- Missing row: add a completeness/quantity rule to prompts, or complete a truncated fixture.
- Merged ordinary questions: add a rule that each numbered item is one `question`.
- Reading/cloze split into many parents: add a rule that one passage is one parent with `choice_items`.
- Inferred answers when answer section is missing: strengthen answer-source rules; answer keys must come from explicit answer lines.
- Stem formatting only: do not overfit prompts unless content is materially missing. The compare script should compare structure, answers, and analysis/explanation, not exact stem prose.
- Missing or mismatched analysis: make sure the fixture includes complete “解析/答案及解析” text, then prompt both providers to copy or summarize the explicit explanation instead of inventing one.
- If analysis wording diverges, first tighten prompts so `analysis_text` uses the source explanation/key sentence rather than provider-written explanations; then rerun the focused case.
- Keep analysis comparison practical for OCR sources: require the same extracted explanation content, but tolerate whitespace, punctuation, and obvious OCR spelling noise.

## Prompt Adjustment Rules

Prefer editing prompts before code:

- System prompt: durable schema, grouping, answer-source, and quantity rules.
- User prompt: concise hard constraints repeated near the raw text.
- Keep examples short and directly tied to observed failure patterns.
- Avoid provider-specific wording unless both providers repeatedly diverge on the same rule.
- After every prompt change, rerun focused cases before full batches.

Do not add Go business logic for parser behavior unless prompt/config cannot express the rule; if Go is needed, explain why first.

## Verification Checklist

After fixture, script, or prompt changes:

```bash
cd /data/project/ai-study
python3 -m py_compile backend/scripts/compare_question_ai_parse.py
python3 -c "import json; [json.load(open(p, encoding='utf-8')) for p in ['backend/testdata/question_ai_pdf_fragments/manifest.json','backend/testdata/question_ai_pdf_fragments/manifest_selected.json']]"

cd /data/project/ai-study/backend
go test ./plugins -count=1 -run 'TestQuestionAIInstructionsRequireDialogueSplit|TestQuestionAIUserPromptRepeatsCoverageRules'
go test ./... -count=1
```

Then run Codex/DeepSeek comparison and cite the report paths.

## Keep This Skill Current

When a new recurring parsing failure is discovered:

1. Add the smallest reusable rule to `ai_parse_system.md` or `ai_parse_user.md`.
2. Add or adjust a focused fixture for that pattern.
3. Rerun focused comparison until Codex and DeepSeek are `diff: same`.
4. Add a concise note to this `SKILL.md` only if the workflow itself changes or the pattern is likely to recur.
