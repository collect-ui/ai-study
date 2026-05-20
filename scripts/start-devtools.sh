#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null 2>&1 && pwd)"
TOOLS_DIR="$ROOT_DIR/.tools"
DEVTOOLS_DIR="$TOOLS_DIR/WeChat_Dev_Tools_v2.01.2510290-2_x86_64_linux"
GTK_ROOT="$TOOLS_DIR/gtk-root"
RUNTIME_DIR="$ROOT_DIR/.runtime"
DATA_DIR="${XDG_CONFIG_HOME:-$HOME/.config}"
APP_NAME="wechat-devtools"

DISPLAY_NUM="${DISPLAY_NUM:-99}"
export DISPLAY="${DISPLAY:-:$DISPLAY_NUM}"
HTTP_PORT="${WECHAT_IDE_HTTP_PORT:-3799}"
REMOTE_PORT="${WECHAT_REMOTE_PORT:-3798}"
XVFB_SCREEN="${XVFB_SCREEN:-1280x800x24}"

mkdir -p "$RUNTIME_DIR"

export LD_LIBRARY_PATH="$GTK_ROOT/usr/lib/x86_64-linux-gnu:${LD_LIBRARY_PATH:-}"
export XDG_DATA_DIRS="$GTK_ROOT/usr/share:${XDG_DATA_DIRS:-/usr/local/share:/usr/share}"
export GIO_EXTRA_MODULES="$GTK_ROOT/usr/lib/x86_64-linux-gnu/gio/modules${GIO_EXTRA_MODULES:+:$GIO_EXTRA_MODULES}"
export GSETTINGS_SCHEMA_DIR="$GTK_ROOT/usr/share/glib-2.0/schemas"
export WECHAT_DEVTOOLS_DIR="$DEVTOOLS_DIR/nwjs"
export APPDATA="$DATA_DIR/$APP_NAME"
export USERPROFILE="$APPDATA"
export PATH="$DEVTOOLS_DIR/node/bin:$DEVTOOLS_DIR/nwjs:$PATH"
mkdir -p "$APPDATA/Default"

if [ -L "$APPDATA/SingletonLock" ]; then
  lock_target="$(readlink "$APPDATA/SingletonLock" || true)"
  lock_pid="${lock_target##*-}"
  if [ -n "$lock_pid" ] && ! kill -0 "$lock_pid" >/dev/null 2>&1; then
    backup_dir="$RUNTIME_DIR/singleton-backup/$(date +%Y%m%d%H%M%S)"
    mkdir -p "$backup_dir"
    for singleton_file in "$APPDATA"/Singleton*; do
      [ -e "$singleton_file" ] && mv "$singleton_file" "$backup_dir/"
    done
  fi
fi

CLI="$TOOLS_DIR/bin/wechat-devtools-cli"
NW="$DEVTOOLS_DIR/nwjs/nw"
APP="$DEVTOOLS_DIR/package.nw"
PLUGIN="$APP/js/ideplugin"

if [ ! -x "$NW" ]; then
  echo "Cannot find WeChat DevTools runtime: $NW" >&2
  exit 1
fi

if [ ! -S "/tmp/.X11-unix/X$DISPLAY_NUM" ]; then
  if ! command -v Xvfb >/dev/null 2>&1; then
    echo "Xvfb is required but was not found." >&2
    exit 1
  fi

  Xvfb "$DISPLAY" -screen 0 "$XVFB_SCREEN" -nolisten tcp >"$RUNTIME_DIR/xvfb.log" 2>&1 &
  echo "$!" >"$RUNTIME_DIR/xvfb.pid"
  sleep 1
fi

if ! ss -ltn 2>/dev/null | grep -q "127.0.0.1:$HTTP_PORT"; then
  launch_args=(
    "$APP"
    "--load-extension=$PLUGIN"
    "--custom-devtools-frontend=file://$PLUGIN/inspector"
    --cli
    --remote-port "$REMOTE_PORT"
    --ide-http-port "$HTTP_PORT"
    --disable-gpu
    --enable-service-port
    --lang zh
  )

  if command -v setsid >/dev/null 2>&1; then
    setsid "$NW" "${launch_args[@]}" >"$RUNTIME_DIR/wechat-devtools.log" 2>&1 < /dev/null &
  else
    nohup "$NW" "${launch_args[@]}" >"$RUNTIME_DIR/wechat-devtools.log" 2>&1 < /dev/null &
  fi
  echo "$!" >"$RUNTIME_DIR/wechat-devtools.pid"

  for _ in $(seq 1 60); do
    if ss -ltn 2>/dev/null | grep -q "127.0.0.1:$HTTP_PORT"; then
      break
    fi
    sleep 1
  done
fi

if ! ss -ltn 2>/dev/null | grep -q "127.0.0.1:$HTTP_PORT"; then
  echo "WeChat DevTools HTTP service did not start on 127.0.0.1:$HTTP_PORT." >&2
  echo "See $RUNTIME_DIR/wechat-devtools.log" >&2
  exit 1
fi

"$CLI" open --project "$ROOT_DIR" --port "$HTTP_PORT" --disable-gpu --trust-project
"$CLI" auto --project "$ROOT_DIR" --auto-port "$HTTP_PORT" --trust-project >/dev/null 2>&1 || true

sleep 5
if ! ss -ltn 2>/dev/null | grep -q "127.0.0.1:$HTTP_PORT"; then
  echo "WeChat DevTools exited after project open." >&2
  echo "See $RUNTIME_DIR/wechat-devtools.log" >&2
  exit 1
fi

echo "WeChat DevTools is listening on http://127.0.0.1:$HTTP_PORT"
