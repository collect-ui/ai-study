#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REMOTE_HOST="${AI_STUDY_REMOTE_HOST:-202.140.140.117}"
REMOTE_USER="${AI_STUDY_REMOTE_USER:-root}"
REMOTE_ROOT="${AI_STUDY_REMOTE_ROOT:-/data/file}"
SERVICE_NAME="${AI_STUDY_FILE_SERVICE_NAME:-ai-study-file.service}"
DOMAIN="${AI_STUDY_FILE_DOMAIN:-collect-ui.top}"
PORT="${AI_STUDY_FILE_PORT:-443}"
SERVER_SCRIPT="$PROJECT_ROOT/scripts/remote-file-server.py"
KNOWN_HOSTS="${AI_STUDY_KNOWN_HOSTS:-/tmp/ai-study-known-hosts}"
RENEW_CERT="${AI_STUDY_RENEW_CERT:-0}"
CERT_DAYS="${AI_STUDY_CERT_DAYS:-365}"
CERT_MODE="${AI_STUDY_CERT_MODE:-preserve}"

SSH_BASE=(ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile="$KNOWN_HOSTS")
RSYNC_BASE=(rsync -az)

if [[ -n "${AUTO_CHECK_SERVER_PASSWORD:-}" || -n "${TEST_SERVER_PASSWORD:-}" ]]; then
  export SSHPASS="${AUTO_CHECK_SERVER_PASSWORD:-${TEST_SERVER_PASSWORD:-}}"
  SSH_BASE=(sshpass -e "${SSH_BASE[@]}")
  RSYNC_BASE=(sshpass -e "${RSYNC_BASE[@]}")
fi

if [[ ! -f "$SERVER_SCRIPT" ]]; then
  echo "server script not found: $SERVER_SCRIPT" >&2
  exit 1
fi

REMOTE="$REMOTE_USER@$REMOTE_HOST"
"${SSH_BASE[@]}" "$REMOTE" "mkdir -p '$REMOTE_ROOT/bin' '$REMOTE_ROOT/certs' '$REMOTE_ROOT/ai-study/assets'"
"${RSYNC_BASE[@]}" "$SERVER_SCRIPT" "$REMOTE:$REMOTE_ROOT/bin/remote-file-server.py"

remote_script="$(cat <<REMOTE_SCRIPT
set -euo pipefail
CERT="$REMOTE_ROOT/certs/$DOMAIN.crt"
KEY="$REMOTE_ROOT/certs/$DOMAIN.key"

chmod 755 "$REMOTE_ROOT/bin/remote-file-server.py"

if [[ ! -f "\$CERT" || ! -f "\$KEY" ]]; then
  openssl req -x509 -newkey rsa:2048 -sha256 -days "$CERT_DAYS" -nodes \\
    -subj "/CN=$DOMAIN" \\
    -addext "subjectAltName=DNS:$DOMAIN,IP:$REMOTE_HOST" \\
    -keyout "\$KEY" \\
    -out "\$CERT"
  chmod 600 "\$KEY"
elif [[ "$RENEW_CERT" == "1" && "$CERT_MODE" == "self-signed" ]]; then
  openssl req -x509 -newkey rsa:2048 -sha256 -days "$CERT_DAYS" -nodes \\
    -subj "/CN=$DOMAIN" \\
    -addext "subjectAltName=DNS:$DOMAIN,IP:$REMOTE_HOST" \\
    -keyout "\$KEY" \\
    -out "\$CERT"
  chmod 600 "\$KEY"
elif [[ "$RENEW_CERT" == "1" ]]; then
  echo "existing cert preserved; set AI_STUDY_CERT_MODE=self-signed only for temporary self-signed replacement"
fi

cat > "/etc/systemd/system/$SERVICE_NAME" <<'UNIT'
[Unit]
Description=AI Study static file HTTPS service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
Environment=AI_STUDY_FILE_ROOT=$REMOTE_ROOT
Environment=AI_STUDY_FILE_HOST=0.0.0.0
Environment=AI_STUDY_FILE_PORT=$PORT
Environment=AI_STUDY_FILE_CERT=$REMOTE_ROOT/certs/$DOMAIN.crt
Environment=AI_STUDY_FILE_KEY=$REMOTE_ROOT/certs/$DOMAIN.key
WorkingDirectory=$REMOTE_ROOT
ExecStart=/usr/bin/python3 $REMOTE_ROOT/bin/remote-file-server.py
Restart=always
RestartSec=3
User=root

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"
systemctl --no-pager --full status "$SERVICE_NAME" | sed -n '1,16p'
REMOTE_SCRIPT
)"

printf '%s\n' "$remote_script" | "${SSH_BASE[@]}" "$REMOTE" "bash -s"

echo "deployed $SERVICE_NAME on $REMOTE_HOST:$PORT"
