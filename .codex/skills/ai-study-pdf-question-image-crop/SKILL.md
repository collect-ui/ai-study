---
name: ai-study-pdf-question-image-crop
description: Extract and crop image-based question regions from AI Study PDFs, especially scanned/OCR worksheets, then attach mobile-friendly crops to question stem_html and question_asset. Use when importing AI Study question PDFs where OCR text omits pictures, diagrams, image options, letter-shape pictures, or any visual prompt that must appear inside the saved exam question.
---

# AI Study PDF Question Image Crop

## Core Rule

Do not treat OCR text as enough when a question depends on a picture. OCR provides text only; preserve visuals by rendering the PDF page, cropping the question image area, saving the crop under the backend file directory, and injecting `<img>` tags into `stem_html`.

For AI Study, use these local conventions:

- Project root: `/data/project/ai-study`
- Backend file root: `/data/project/ai-study/backend/file_data/files`
- Public URL prefix: `/files`
- Question fields: `stem_html`, `stem_text`, `asset_count`
- Asset table/service: `question_asset` via `question.asset_save`
- Typical asset usage: `usage_type=stem`, `usage_ref=stem_html`

## Workflow

1. Split large PDFs first. Do not send a full 15MB scanned PDF through OCR or AI at once.
2. Render only the needed page or small fragment at 180 DPI unless the source is blurry.
3. Inspect the rendered page image and identify visual question regions.
4. Crop mobile-friendly images:
   - Prefer one crop per subquestion/image option for phones.
   - Keep a full-question crop only as debug evidence or when the layout cannot be split safely.
   - Add a small white border if the crop touches text or page edges.
5. Save crops under `backend/file_data/files/question-stem/YYYY-MM-DD/<source-key>/`.
6. Use `/files/question-stem/YYYY-MM-DD/<source-key>/<file>` URLs in `stem_html`.
7. Save or update the question so `stem_html` contains the images and the text prompt remains searchable in `stem_text`.
8. Register each crop with `question.asset_save` when the question ID is known.
9. Verify with:
   - `curl --noproxy '*' -I http://127.0.0.1:8026/files/...`
   - `question.question_choice_detail` contains `<img`
   - `question.asset_query` returns the saved assets

## Script

Use `scripts/crop_question_images.py` for repeatable cropping.

Example:

```bash
python3 /data/project/ai-study/.codex/skills/ai-study-pdf-question-image-crop/scripts/crop_question_images.py \
  --image test-results/grade3-english-prep-fragments/image-debug/page-0003.jpg \
  --out-dir backend/file_data/files/question-stem/2026-05-21/grade3-unit1-parta-mobile \
  --url-prefix /files/question-stem/2026-05-21/grade3-unit1-parta-mobile \
  --prompt '观察图片，看看它们分别像什么字母，并将其写在四线三格内。' \
  --crops '[{"name":"q1-1.jpg","label":"1","box":[230,405,410,590]},{"name":"q1-2.jpg","label":"2","box":[505,395,725,590]}]' \
  --report test-results/grade3-english-prep-fragments/reports/mobile-crops.json
```

The script outputs JSON with `assets` and `stem_html`. Use those values when saving or updating the question.

## AI Parse Integration

When `question.ai_parse` returns rows like `第1题（图片）`, post-process them before saving:

- Keep `stem_text` textual, e.g. `观察图片，看看它像什么字母。第1题`
- Set `stem_html` to the text plus the crop `<img>`
- Set `question_type=blank` for image-to-letter answers unless the UI supports image choices
- Set `blank_answers` from the answer section
- Keep `source` and `remark` traceable, e.g. `ai_import_image_crop`

Do not invent images from OCR descriptions. Crop from the rendered PDF page.

## Reference

Read `references/ai-study-question-image-schema.md` when touching save services, asset rows, or mobile display details.
