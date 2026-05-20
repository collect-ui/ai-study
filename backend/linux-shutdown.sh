#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PID_FILE="run-dev.pid"
LOG_FILE="run-dev.log"

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
  echo "PID file not found: $PID_FILE"
  echo "Checking for processes on port 8016..."
  
  # 检查8016端口是否有进程
  if command -v lsof >/dev/null 2>&1; then
    PIDS=$(lsof -ti:8016 2>/dev/null || true)
    if [ -n "$PIDS" ]; then
      echo "Found processes on port 8016: $PIDS"
      echo -n "Kill these processes? (y/N): "
      read -r response
      if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "$PIDS" | xargs kill -9 2>/dev/null || true
        echo "Processes killed."
      else
        echo "Aborted."
        exit 0
      fi
    else
      echo "No processes found on port 8016."
    fi
  else
    echo "lsof command not available. Checking with ss/netstat..."
    if command -v ss >/dev/null 2>&1; then
      PIDS=$(ss -tlnp | grep :8016 | awk '{print $NF}' | cut -d= -f2 | cut -d, -f1 | sort -u)
    elif command -v netstat >/dev/null 2>&1; then
      PIDS=$(netstat -tlnp 2>/dev/null | grep :8016 | awk '{print $NF}' | cut -d/ -f1 | sort -u)
    fi
    
    if [ -n "$PIDS" ]; then
      echo "Found processes on port 8016: $PIDS"
      echo -n "Kill these processes? (y/N): "
      read -r response
      if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "$PIDS" | xargs kill -9 2>/dev/null || true
        echo "Processes killed."
      else
        echo "Aborted."
        exit 0
      fi
    else
      echo "No processes found on port 8016."
    fi
  fi
  exit 0
fi

# 读取PID
PID="$(cat "$PID_FILE" 2>/dev/null || true)"

if [ -z "$PID" ]; then
  echo "PID file is empty: $PID_FILE"
  rm -f "$PID_FILE"
  exit 0
fi

# 检查进程是否存在
if kill -0 "$PID" 2>/dev/null; then
  echo "Stopping process (PID: $PID)..."
  
  # 先尝试正常停止
  kill "$PID" 2>/dev/null || true
  
  # 等待最多10秒
  for i in {1..10}; do
    if ! kill -0 "$PID" 2>/dev/null; then
      echo "Process stopped gracefully."
      break
    fi
    sleep 1
    echo -n "."
  done
  
  # 如果进程还在，强制停止
  if kill -0 "$PID" 2>/dev/null; then
    echo -e "\nProcess not responding, forcing kill..."
    kill -9 "$PID" 2>/dev/null || true
    sleep 1
  fi
  
  # 再次检查
  if kill -0 "$PID" 2>/dev/null; then
    echo "Warning: Process may still be running (PID: $PID)"
  else
    echo "Process stopped."
  fi
else
  echo "Process not running (PID: $PID)"
fi

# 清理PID文件
rm -f "$PID_FILE"
echo "Cleaned up PID file: $PID_FILE"

# 可选：清理日志文件
if [ -f "$LOG_FILE" ]; then
  echo "Log file: $LOG_FILE (size: $(du -h "$LOG_FILE" | cut -f1))"
  echo -n "Clear log file? (y/N): "
  read -r response
  if [[ "$response" =~ ^[Yy]$ ]]; then
    > "$LOG_FILE"
    echo "Log file cleared."
  fi
fi

echo "Shutdown completed."