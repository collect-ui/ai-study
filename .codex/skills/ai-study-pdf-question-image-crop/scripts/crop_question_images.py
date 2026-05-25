#!/usr/bin/env python3
import argparse
import hashlib
import json
from pathlib import Path

from PIL import Image


def load_crops(raw):
    if raw.startswith("@"):
        raw = Path(raw[1:]).read_text(encoding="utf-8")
    crops = json.loads(raw)
    if not isinstance(crops, list):
        raise ValueError("crops must be a JSON array")
    for index, crop in enumerate(crops, 1):
        if not isinstance(crop, dict):
            raise ValueError(f"crop {index} must be an object")
        if "box" not in crop or len(crop["box"]) != 4:
            raise ValueError(f"crop {index} must contain box=[left,top,right,bottom]")
    return crops


def render_pdf_page(pdf_path, page_no, dpi):
    import pypdfium2 as pdfium

    document = pdfium.PdfDocument(str(pdf_path))
    try:
        page = document[page_no - 1]
        try:
            bitmap = page.render(scale=max(36, dpi) / 72.0)
            try:
                image = bitmap.to_pil()
                if image.mode != "RGB":
                    image = image.convert("RGB")
                return image
            finally:
                bitmap.close()
        finally:
            page.close()
    finally:
        document.close()


def image_for_input(args):
    if args.image:
        image = Image.open(args.image)
        if image.mode != "RGB":
            image = image.convert("RGB")
        return image
    if args.pdf:
        return render_pdf_page(Path(args.pdf), args.page, args.dpi)
    raise ValueError("Provide --image or --pdf")


def apply_border(image, border):
    if border <= 0:
        return image
    out = Image.new("RGB", (image.width + border * 2, image.height + border * 2), "white")
    out.paste(image, (border, border))
    return out


def resize_max_width(image, max_width):
    if max_width <= 0 or image.width <= max_width:
        return image
    height = round(image.height * (max_width / image.width))
    return image.resize((max_width, height), Image.LANCZOS)


def file_meta(path, url, crop):
    data = path.read_bytes()
    return {
        "name": path.name,
        "path": str(path),
        "url": url,
        "size": len(data),
        "sha256": hashlib.sha256(data).hexdigest(),
        "mime_type": "image/jpeg",
        "label": str(crop.get("label", "")),
        "alt": str(crop.get("alt", crop.get("label", path.stem))),
        "box": crop["box"],
    }


def build_stem_html(prompt, assets, columns, max_img_width):
    cells = []
    for asset in assets:
        label = asset["label"]
        prefix = f"<strong>{label}.</strong><br>" if label else ""
        cells.append(
            '<div style="font-size:14px;line-height:1.4;">'
            f'{prefix}<img src="{asset["url"]}" alt="{asset["alt"]}" '
            f'style="width:100%;max-width:{max_img_width}px;height:auto;" />'
            "</div>"
        )
    return (
        f"<div>{prompt}</div>"
        f'<div style="display:grid;grid-template-columns:repeat({columns},minmax(0,1fr));'
        'gap:8px;max-width:520px;">'
        + "".join(cells)
        + "</div>"
    )


def main():
    parser = argparse.ArgumentParser(description="Crop mobile-friendly question images for AI Study.")
    parser.add_argument("--pdf", default="")
    parser.add_argument("--page", type=int, default=1)
    parser.add_argument("--dpi", type=int, default=180)
    parser.add_argument("--image", default="")
    parser.add_argument("--out-dir", required=True)
    parser.add_argument("--url-prefix", required=True)
    parser.add_argument("--crops", required=True, help="JSON array or @path to JSON array")
    parser.add_argument("--prompt", default="")
    parser.add_argument("--border", type=int, default=12)
    parser.add_argument("--max-width", type=int, default=260)
    parser.add_argument("--columns", type=int, default=2)
    parser.add_argument("--quality", type=int, default=90)
    parser.add_argument("--report", default="")
    args = parser.parse_args()

    crops = load_crops(args.crops)
    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    image = image_for_input(args)
    assets = []
    url_prefix = args.url_prefix.rstrip("/")

    for index, crop in enumerate(crops, 1):
        box = tuple(int(v) for v in crop["box"])
        name = crop.get("name") or f"crop-{index:02d}.jpg"
        if not name.lower().endswith((".jpg", ".jpeg")):
            name += ".jpg"
        cropped = image.crop(box)
        cropped = apply_border(cropped, args.border)
        cropped = resize_max_width(cropped, args.max_width)
        out_path = out_dir / name
        cropped.save(out_path, "JPEG", quality=args.quality, optimize=True)
        assets.append(file_meta(out_path, f"{url_prefix}/{name}", crop))

    result = {
        "source_image": args.image,
        "source_pdf": args.pdf,
        "page": args.page,
        "page_size": list(image.size),
        "assets": assets,
        "stem_html": build_stem_html(args.prompt, assets, max(1, args.columns), args.max_width),
    }

    text = json.dumps(result, ensure_ascii=False, indent=2)
    if args.report:
        Path(args.report).write_text(text + "\n", encoding="utf-8")
    print(text)


if __name__ == "__main__":
    main()
