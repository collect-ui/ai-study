#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

PID_FILE="run-dev.pid"
LOG_FILE="run-dev.log"

if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE" 2>/dev/null || true)
    if [ -n "$OLD_PID" ] && kill -0 "$OLD_PID" 2>/dev/null; then
        echo "Already running, PID: $OLD_PID"
        exit 0
    fi
    rm -f "$PID_FILE"
fi

if [ -x "./bin" ]; then
    CMD="./bin"
elif [ -x "./dist/bin" ]; then
    CMD="./dist/bin"
elif command -v go >/dev/null 2>&1; then
    CMD="go run main.go"
else
    echo "Error: no executable found (./bin, ./dist/bin) and go command not available"
    exit 1
fi

if [ "$CMD" = "go run main.go" ]; then
    nohup sh -c "$CMD" >"$LOG_FILE" 2>&1 &
else
    nohup sh -c "$CMD" >"$LOG_FILE" 2>&1 &
fi

NEW_PID=$!
echo "$NEW_PID" >"$PID_FILE"

echo "Started: PID=$NEW_PID"
echo "Log: $SCRIPT_DIR/$LOG_FILE"
echo "PID file: $SCRIPT_DIR/$PID_FILE"
