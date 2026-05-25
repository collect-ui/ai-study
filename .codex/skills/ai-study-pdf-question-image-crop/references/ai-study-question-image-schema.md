# AI Study Question Image Schema Notes

## Fields

- `question_item.stem_html`: rich HTML shown as the question stem. Put `<img src="/files/...">` here.
- `question_item.stem_text`: searchable plain text. Keep the prompt and omit HTML.
- `question_asset`: attachment registry for question media.
  - `usage_type`: use `stem` for stem images.
  - `usage_ref`: use `stem_html` when the image appears in `stem_html`.
  - `asset_url`: public `/files/...` URL.
  - `asset_name`, `mime_type`, `file_size`, `sha256`: store crop metadata.

## Services

- Save/update question:
  - `question.question_choice_save`
  - `question.question_update`
- Save/query assets:
  - `question.asset_save`
  - `question.asset_query`
- Verify details:
  - `question.question_choice_detail`

## Mobile HTML

Prefer small crops in a wrapping grid:

```html
<div>观察图片，看看它们分别像什么字母，并将其写在四线三格内。</div>
<div style="display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:8px;max-width:520px;">
  <div style="font-size:14px;"><strong>1.</strong><br><img src="/files/..." style="width:100%;max-width:220px;height:auto;"></div>
</div>
```

Avoid one very wide crop for phone screens unless the whole row must be preserved.

## Current Import Risks

- OCR text does not contain image pixels.
- Full-page OCR on two-column answer pages may interleave columns.
- Existing single-choice save checks may require A-D in some flows; image-to-letter tasks are safer as `blank` questions with image crops in `stem_html`.
