#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REMOTE_HOST="${AI_STUDY_REMOTE_HOST:-202.140.140.117}"
REMOTE_USER="${AI_STUDY_REMOTE_USER:-root}"
REMOTE_ROOT="${AI_STUDY_REMOTE_ROOT:-/data/file}"
REMOTE_ASSET_DIR="${AI_STUDY_REMOTE_ASSET_DIR:-$REMOTE_ROOT/ai-study/assets}"
LOCAL_ASSET_DIR="${AI_STUDY_LOCAL_ASSET_DIR:-$PROJECT_ROOT/feature/remote-assets/staging}"
MANIFEST_PATH="${AI_STUDY_MANIFEST_PATH:-$PROJECT_ROOT/feature/remote-assets/asset-manifest.json}"
KNOWN_HOSTS="${AI_STUDY_KNOWN_HOSTS:-/tmp/ai-study-known-hosts}"

SSH_BASE=(ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile="$KNOWN_HOSTS")
RSYNC_BASE=(rsync -az --delete)

if [[ -n "${AUTO_CHECK_SERVER_PASSWORD:-}" || -n "${TEST_SERVER_PASSWORD:-}" ]]; then
  export SSHPASS="${AUTO_CHECK_SERVER_PASSWORD:-${TEST_SERVER_PASSWORD:-}}"
  SSH_BASE=(sshpass -e "${SSH_BASE[@]}")
  RSYNC_BASE=(sshpass -e "${RSYNC_BASE[@]}")
fi

if [[ ! -d "$LOCAL_ASSET_DIR" ]]; then
  echo "local asset dir not found: $LOCAL_ASSET_DIR" >&2
  echo "create feature/remote-assets/staging or set AI_STUDY_LOCAL_ASSET_DIR=/path/to/assets" >&2
  exit 1
fi

if ! find "$LOCAL_ASSET_DIR" -type f | grep -q .; then
  echo "local asset dir has no files: $LOCAL_ASSET_DIR" >&2
  exit 1
fi

mkdir -p "$(dirname "$MANIFEST_PATH")"
(
  cd "$LOCAL_ASSET_DIR"
  find . -type f -print0 |
    sort -z |
    xargs -0 sha256sum |
    sed 's#  ./#  #'
) > "$MANIFEST_PATH"

"${SSH_BASE[@]}" "$REMOTE_USER@$REMOTE_HOST" "mkdir -p '$REMOTE_ASSET_DIR'"
"${RSYNC_BASE[@]}" "$LOCAL_ASSET_DIR"/ "$REMOTE_USER@$REMOTE_HOST:$REMOTE_ASSET_DIR"/
"${RSYNC_BASE[@]}" "$MANIFEST_PATH" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_ROOT/ai-study/asset-manifest.sha256"

echo "synced assets to $REMOTE_USER@$REMOTE_HOST:$REMOTE_ASSET_DIR"
echo "manifest: $MANIFEST_PATH"
