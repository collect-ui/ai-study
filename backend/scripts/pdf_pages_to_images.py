#!/usr/bin/env python3

import argparse
import json
from pathlib import Path

import pypdfium2 as pdfium


def render_pdf_pages(pdf_path: Path, out_dir: Path, dpi: int, max_pages: int, image_format: str, page_number: int = 0) -> dict:
    out_dir.mkdir(parents=True, exist_ok=True)
    document = pdfium.PdfDocument(str(pdf_path))
    page_count = len(document)
    if page_number > 0:
        if page_number > page_count:
            document.close()
            return {
                "page_count": page_count,
                "rendered_count": 0,
                "pages": [],
            }
        page_indexes = [page_number - 1]
    else:
        limit = page_count if max_pages <= 0 else min(page_count, max_pages)
        page_indexes = list(range(limit))
    scale = max(36, dpi) / 72.0
    pages = []

    for index in page_indexes:
        page = document[index]
        bitmap = None
        image = None
        try:
            bitmap = page.render(scale=scale)
            image = bitmap.to_pil()
            if image.mode not in {"RGB", "L"}:
                image = image.convert("RGB")

            ext = "jpg" if image_format in {"jpg", "jpeg"} else "png"
            out_path = out_dir / f"page-{index + 1:04d}.{ext}"
            if ext == "jpg":
                if image.mode != "RGB":
                    image = image.convert("RGB")
                image.save(out_path, format="JPEG", quality=90, optimize=True)
            else:
                image.save(out_path, format="PNG", optimize=True)
            pages.append({
                "page": index + 1,
                "path": str(out_path),
                "width": image.width,
                "height": image.height,
            })
        finally:
            if image is not None:
                image.close()
            if bitmap is not None:
                bitmap.close()
            page.close()

    document.close()
    return {
        "page_count": page_count,
        "rendered_count": len(pages),
        "pages": pages,
    }


def count_pdf_pages(pdf_path: Path) -> dict:
    document = pdfium.PdfDocument(str(pdf_path))
    try:
        return {
            "page_count": len(document),
            "rendered_count": 0,
            "pages": [],
        }
    finally:
        document.close()


def main() -> int:
    parser = argparse.ArgumentParser(description="Render PDF pages to images for OCR.")
    parser.add_argument("--pdf", required=True)
    parser.add_argument("--out-dir", required=True)
    parser.add_argument("--dpi", type=int, default=180)
    parser.add_argument("--max-pages", type=int, default=80)
    parser.add_argument("--page", type=int, default=0)
    parser.add_argument("--format", choices=["jpg", "jpeg", "png"], default="jpg")
    parser.add_argument("--count-only", action="store_true")
    args = parser.parse_args()

    if args.count_only:
        result = count_pdf_pages(pdf_path=Path(args.pdf))
    else:
        result = render_pdf_pages(
            pdf_path=Path(args.pdf),
            out_dir=Path(args.out_dir),
            dpi=args.dpi,
            max_pages=args.max_pages,
            page_number=args.page,
            image_format=args.format,
        )
    print(json.dumps(result, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
