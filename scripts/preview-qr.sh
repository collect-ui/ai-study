#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PORT="${WECHAT_DEVTOOLS_PORT:-3799}"
OUTPUT_DIR="$ROOT/wx_login/qr"
QR_OUTPUT="$OUTPUT_DIR/ai-study-preview-qr.png"
REAL_PNG_OUTPUT="$OUTPUT_DIR/ai-study-preview-qr-real.png"
JPEG_OUTPUT="$OUTPUT_DIR/ai-study-preview-qr.jpg"
INFO_OUTPUT="$OUTPUT_DIR/ai-study-preview-info.json"
PREVIEW_PROJECT="$ROOT/miniprogram"

mkdir -p "$OUTPUT_DIR"

"$ROOT/.tools/bin/wechat-devtools-cli" preview \
  --project "$PREVIEW_PROJECT" \
  --port "$PORT" \
  --disable-gpu \
  --qr-format image \
  --qr-output "$QR_OUTPUT" \
  --qr-size 430 \
  --info-output "$INFO_OUTPUT"

python3 - "$QR_OUTPUT" "$REAL_PNG_OUTPUT" "$JPEG_OUTPUT" <<'PY'
from PIL import Image
import sys

src, png_output, jpeg_output = sys.argv[1:4]
image = Image.open(src).convert("RGB")
white = Image.new("RGB", image.size, "white")
white.paste(image)
white.save(png_output, "PNG")
white.save(jpeg_output, "JPEG", quality=95)
PY

echo "QR: $QR_OUTPUT"
echo "Real PNG QR: $REAL_PNG_OUTPUT"
echo "JPEG QR: $JPEG_OUTPUT"
echo "Info: $INFO_OUTPUT"
