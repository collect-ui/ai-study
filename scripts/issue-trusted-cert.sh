#!/usr/bin/env bash
set -euo pipefail

REMOTE_HOST="${AI_STUDY_REMOTE_HOST:-202.140.140.117}"
REMOTE_USER="${AI_STUDY_REMOTE_USER:-root}"
REMOTE_ROOT="${AI_STUDY_REMOTE_ROOT:-/data/file}"
SERVICE_NAME="${AI_STUDY_FILE_SERVICE_NAME:-ai-study-file.service}"
DOMAIN="${AI_STUDY_FILE_DOMAIN:-collect-ui.top}"
EMAIL="${AI_STUDY_ACME_EMAIL:-admin@collect-ui.top}"
KNOWN_HOSTS="${AI_STUDY_KNOWN_HOSTS:-/tmp/ai-study-known-hosts}"
FORCE="${AI_STUDY_ACME_FORCE:-0}"

SSH_BASE=(ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile="$KNOWN_HOSTS")

if [[ -n "${AUTO_CHECK_SERVER_PASSWORD:-}" || -n "${TEST_SERVER_PASSWORD:-}" ]]; then
  export SSHPASS="${AUTO_CHECK_SERVER_PASSWORD:-${TEST_SERVER_PASSWORD:-}}"
  SSH_BASE=(sshpass -e "${SSH_BASE[@]}")
fi

issue_args=(--issue --alpn -d "$DOMAIN" --server letsencrypt)
if [[ "$FORCE" == "1" ]]; then
  issue_args+=(--force)
fi

REMOTE="$REMOTE_USER@$REMOTE_HOST"
remote_script="$(cat <<REMOTE_SCRIPT
set -euo pipefail
if [ ! -x /root/.acme.sh/acme.sh ]; then
  curl https://get.acme.sh | sh -s email=$EMAIL
fi
/root/.acme.sh/acme.sh --set-default-ca --server letsencrypt
systemctl stop "$SERVICE_NAME"
restore_service() {
  systemctl start "$SERVICE_NAME" || true
}
trap restore_service EXIT
/root/.acme.sh/acme.sh ${issue_args[*]}
/root/.acme.sh/acme.sh --install-cert -d "$DOMAIN" \\
  --key-file "$REMOTE_ROOT/certs/$DOMAIN.key" \\
  --fullchain-file "$REMOTE_ROOT/certs/$DOMAIN.crt" \\
  --reloadcmd "systemctl restart $SERVICE_NAME"
systemctl start "$SERVICE_NAME" || true
systemctl restart "$SERVICE_NAME"
trap - EXIT
systemctl is-active "$SERVICE_NAME"
openssl x509 -in "$REMOTE_ROOT/certs/$DOMAIN.crt" -noout -issuer -subject -dates
REMOTE_SCRIPT
)"

printf '%s\n' "$remote_script" | "${SSH_BASE[@]}" "$REMOTE" "bash -s"
