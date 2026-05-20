#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null 2>&1 && pwd)"
HTTP_PORT="${WECHAT_IDE_HTTP_PORT:-3799}"

if ss -ltn 2>/dev/null | grep -q "127.0.0.1:$HTTP_PORT"; then
  echo "service: http://127.0.0.1:$HTTP_PORT"
else
  echo "service: not listening on 127.0.0.1:$HTTP_PORT"
  exit 1
fi

"$ROOT_DIR/.tools/bin/wechat-devtools-cli" islogin --port "$HTTP_PORT" --disable-gpu
