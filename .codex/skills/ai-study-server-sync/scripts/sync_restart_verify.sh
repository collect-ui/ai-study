#!/usr/bin/env bash
set -euo pipefail

HOST="${AI_STUDY_SERVER_HOST:-202.140.140.117}"
USER_NAME="${AI_STUDY_SERVER_USER:-root}"
PASSWORD="${AUTO_CHECK_SERVER_PASSWORD:-${TEST_SERVER_PASSWORD:-}}"
REMOTE_DIR="${AI_STUDY_SERVER_REMOTE_DIR:-/data/ai-study}"
LOCAL_BACKEND_DIR="${AI_STUDY_SERVER_LOCAL_BACKEND_DIR:-/data/project/ai-study/backend}"
REMOTE_BACKUP_DIR="${AI_STUDY_SERVER_BACKUP_DIR:-/data/ai-study-backup}"
SERVICE_NAME="${AI_STUDY_SERVER_SERVICE_NAME:-ai-study.service}"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}"
VERIFY_PATH="${AI_STUDY_SERVER_VERIFY_PATH:-/collect-ui/}"
SEED_CONF_IF_MISSING="${AI_STUDY_SERVER_SEED_CONF_IF_MISSING:-1}"
DEFAULT_PORT="${AI_STUDY_SERVER_DEFAULT_PORT:-8026}"

SYNC_ONLY=0
VERIFY_ONLY=0
KEEPALIVE_CHECK=0

usage() {
  cat <<'USAGE'
Usage: sync_restart_verify.sh [--sync-only] [--verify-only] [--keepalive-check] [--no-seed-conf]

Options:
  --sync-only        Build and sync runtime files, then exit without restarting.
  --verify-only      Verify the existing remote service without building or syncing.
  --keepalive-check  Kill the service main PID and verify systemd restarts it.
  --no-seed-conf     Do not copy local conf/application.properties when remote is missing.
USAGE
}

for arg in "$@"; do
  case "$arg" in
    --sync-only) SYNC_ONLY=1 ;;
    --verify-only) VERIFY_ONLY=1 ;;
    --keepalive-check) KEEPALIVE_CHECK=1 ;;
    --no-seed-conf) SEED_CONF_IF_MISSING=0 ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown arg: $arg" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if [[ "${SYNC_ONLY}" -eq 1 && "${VERIFY_ONLY}" -eq 1 ]]; then
  echo "Cannot use --sync-only with --verify-only" >&2
  exit 2
fi

if [[ -z "${PASSWORD}" ]]; then
  echo "Missing password: set AUTO_CHECK_SERVER_PASSWORD (or TEST_SERVER_PASSWORD)." >&2
  exit 2
fi

if [[ ! -d "${LOCAL_BACKEND_DIR}" ]]; then
  echo "Local backend directory not found: ${LOCAL_BACKEND_DIR}" >&2
  exit 2
fi

require_cmd() {
  local name="$1"
  if ! command -v "${name}" >/dev/null 2>&1; then
    echo "${name} is required but not found." >&2
    exit 2
  fi
}

require_cmd sshpass
require_cmd rsync
require_cmd go

if [[ "${VERIFY_PATH}" != /* ]]; then
  VERIFY_PATH="/${VERIFY_PATH}"
fi

run_remote() {
  local cmd="$1"
  SSHPASS="${PASSWORD}" sshpass -e ssh \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    "${USER_NAME}@${HOST}" "${cmd}"
}

rsync_file_to_remote() {
  local src="$1"
  local dst="$2"
  SSHPASS="${PASSWORD}" sshpass -e rsync -az \
    -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
    "${src}" "${USER_NAME}@${HOST}:${dst}"
}

PACKAGE_DIR=""
cleanup() {
  if [[ -n "${PACKAGE_DIR}" && -d "${PACKAGE_DIR}" ]]; then
    rm -rf "${PACKAGE_DIR}"
  fi
}
trap cleanup EXIT

prepare_local_package() {
  PACKAGE_DIR="$(mktemp -d "${TMPDIR:-/tmp}/ai-study-backend-package.XXXXXX")"
  echo "Building local linux binary -> ${PACKAGE_DIR}/bin"
  (
    cd "${LOCAL_BACKEND_DIR}"
    GOOS=linux GOARCH=amd64 go build -o "${PACKAGE_DIR}/bin" main.go
  )
  chmod +x "${PACKAGE_DIR}/bin"

  local item
  for item in collect frontend static conf startup.sh shutdown.sh certbot_renew.sh; do
    if [[ -e "${LOCAL_BACKEND_DIR}/${item}" ]]; then
      cp -a "${LOCAL_BACKEND_DIR}/${item}" "${PACKAGE_DIR}/"
    fi
  done

  local required
  for required in collect frontend conf startup.sh shutdown.sh; do
    if [[ ! -e "${PACKAGE_DIR}/${required}" ]]; then
      echo "Required runtime item missing from package: ${required}" >&2
      exit 1
    fi
  done

  chmod +x "${PACKAGE_DIR}/startup.sh" "${PACKAGE_DIR}/shutdown.sh" 2>/dev/null || true
  echo "PACKAGE_DIR=${PACKAGE_DIR}"
}

backup_remote_state() {
  echo "Backing up remote config/database into ${REMOTE_BACKUP_DIR}"
  local cmd
  cmd=$(cat <<'CMD'
set -e
ts="$(date +%Y%m%d-%H%M%S)"
mkdir -p "__REMOTE_BACKUP_DIR__"
if [ -f "__REMOTE_DIR__/conf/application.properties" ]; then
  mkdir -p "__REMOTE_BACKUP_DIR__/conf-${ts}"
  cp -a "__REMOTE_DIR__/conf/application.properties" "__REMOTE_BACKUP_DIR__/conf-${ts}/application.properties"
  echo "CONFIG_BACKUP=__REMOTE_BACKUP_DIR__/conf-${ts}/application.properties"
else
  echo "CONFIG_BACKUP_SKIPPED=no_remote_conf"
fi
if [ -d "__REMOTE_DIR__/database" ]; then
  mkdir -p "__REMOTE_BACKUP_DIR__/database-${ts}"
  cp -a "__REMOTE_DIR__/database/." "__REMOTE_BACKUP_DIR__/database-${ts}/" 2>/dev/null || true
  echo "DATABASE_BACKUP=__REMOTE_BACKUP_DIR__/database-${ts}"
else
  echo "DATABASE_BACKUP_SKIPPED=no_remote_database_dir"
fi
CMD
)
  cmd="${cmd//__REMOTE_BACKUP_DIR__/${REMOTE_BACKUP_DIR}}"
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  run_remote "${cmd}"
}

remote_conf_exists() {
  local cmd
  cmd=$(cat <<'CMD'
if [ -f "__REMOTE_DIR__/conf/application.properties" ]; then
  echo "YES"
else
  echo "NO"
fi
CMD
)
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  run_remote "${cmd}" | tail -n 1 | tr -d '\r'
}

remote_database_has_entries() {
  local cmd
  cmd=$(cat <<'CMD'
if [ -d "__REMOTE_DIR__/database" ] && [ -n "$(find "__REMOTE_DIR__/database" -mindepth 1 -print -quit 2>/dev/null)" ]; then
  echo "YES"
else
  echo "NO"
fi
CMD
)
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  run_remote "${cmd}" | tail -n 1 | tr -d '\r'
}

sync_package_to_remote() {
  echo "Syncing package to ${USER_NAME}@${HOST}:${REMOTE_DIR}"
  run_remote "mkdir -p '${REMOTE_DIR}'"
  SSHPASS="${PASSWORD}" sshpass -e rsync -az --delete \
    -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
    --exclude "conf/application.properties" \
    --exclude "database/" \
    --exclude "file_data/" \
    --exclude "logs/" \
    --exclude "run-dev.log" \
    --exclude "run-dev.pid" \
    --exclude "run.log" \
    --exclude "*.tar.gz" \
    --exclude "*.zip" \
    "${PACKAGE_DIR}/" "${USER_NAME}@${HOST}:${REMOTE_DIR}/"
}

seed_remote_conf_if_missing() {
  local exists="$1"
  local local_conf="${LOCAL_BACKEND_DIR}/conf/application.properties"
  if [[ "${exists}" == "YES" ]]; then
    echo "Remote conf/application.properties exists. Skip seeding local config."
    return
  fi
  if [[ "${SEED_CONF_IF_MISSING}" != "1" ]]; then
    echo "Remote conf/application.properties is missing and config seeding is disabled." >&2
    exit 1
  fi
  if [[ ! -f "${local_conf}" ]]; then
    echo "Local config not found: ${local_conf}" >&2
    exit 1
  fi
  echo "Remote conf/application.properties is missing. Seeding local config once."
  run_remote "mkdir -p '${REMOTE_DIR}/conf'"
  rsync_file_to_remote "${local_conf}" "${REMOTE_DIR}/conf/application.properties"
}

check_application_properties_diff() {
  local local_conf="${LOCAL_BACKEND_DIR}/conf/application.properties"
  local local_keys remote_keys missing_keys
  local_keys="$(mktemp)"
  remote_keys="$(mktemp)"
  missing_keys="$(mktemp)"

  if [[ ! -f "${local_conf}" ]]; then
    echo "CONFIG_CHECK_SKIPPED=local_conf_missing:${local_conf}"
    rm -f "${local_keys}" "${remote_keys}" "${missing_keys}"
    return
  fi

  sed -n "s/^[[:space:]]*\\([A-Za-z0-9_.-][A-Za-z0-9_.-]*\\)[[:space:]]*=.*/\\1/p" "${local_conf}" \
    | sort -u >"${local_keys}"

  run_remote "if [ -f '${REMOTE_DIR}/conf/application.properties' ]; then sed -n \"s/^[[:space:]]*\\([A-Za-z0-9_.-][A-Za-z0-9_.-]*\\)[[:space:]]*=.*/\\1/p\" '${REMOTE_DIR}/conf/application.properties' | sort -u; fi" \
    >"${remote_keys}"

  if [[ ! -s "${remote_keys}" ]]; then
    echo "CONFIG_CHECK_NOTICE=remote_conf_missing_or_empty:${REMOTE_DIR}/conf/application.properties"
    rm -f "${local_keys}" "${remote_keys}" "${missing_keys}"
    return
  fi

  comm -23 "${local_keys}" "${remote_keys}" >"${missing_keys}" || true
  if [[ -s "${missing_keys}" ]]; then
    echo "CONFIG_CHECK_NOTICE=detected new local keys not in remote conf/application.properties"
    sed 's/^/CONFIG_NEW_KEY=/' "${missing_keys}"
    echo "CONFIG_CHECK_NOTICE=please review and update remote config manually."
  else
    echo "CONFIG_CHECK_OK=no new local keys compared with remote conf/application.properties"
  fi

  rm -f "${local_keys}" "${remote_keys}" "${missing_keys}"
}

copy_local_database_if_needed() {
  local remote_has_entries="$1"
  local local_db="${LOCAL_BACKEND_DIR}/database"
  if [[ "${remote_has_entries}" == "YES" ]]; then
    echo "Remote database already has entries. Skip copying local database."
    return
  fi
  if [[ ! -d "${local_db}" ]]; then
    echo "Local database directory not found. Skip copying."
    return
  fi
  if [[ -z "$(find "${local_db}" -mindepth 1 -print -quit 2>/dev/null)" ]]; then
    echo "Local database has no entries. Skip copying."
    return
  fi
  echo "Remote database is empty. Copying local database for first deployment."
  run_remote "mkdir -p '${REMOTE_DIR}/database'"
  rsync_file_to_remote "${local_db}/" "${REMOTE_DIR}/database/"
}

ensure_remote_service() {
  local cmd
  cmd=$(cat <<'CMD'
set -e
cat > "__SERVICE_FILE__" <<'UNIT'
[Unit]
Description=AI Study Admin Backend
After=network.target

[Service]
Type=forking
PIDFile=__REMOTE_DIR__/run-dev.pid
ExecStart=__REMOTE_DIR__/startup.sh
ExecStop=__REMOTE_DIR__/shutdown.sh
Restart=always
RestartSec=3
TimeoutStopSec=20
KillMode=control-group
User=root
Group=root
WorkingDirectory=__REMOTE_DIR__

[Install]
WantedBy=multi-user.target
UNIT

chmod 644 "__SERVICE_FILE__"
systemctl daemon-reload
systemctl enable "__SERVICE_NAME__"
mkdir -p "__REMOTE_DIR__/database" "__REMOTE_DIR__/file_data/files" "__REMOTE_DIR__/logs"
CMD
)
  cmd="${cmd//__SERVICE_FILE__/${SERVICE_FILE}}"
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  cmd="${cmd//__SERVICE_NAME__/${SERVICE_NAME}}"
  run_remote "${cmd}"
}

ensure_remote_executable_permissions() {
  local cmd
  cmd=$(cat <<'CMD'
set -e
chmod +x "__REMOTE_DIR__/startup.sh" "__REMOTE_DIR__/shutdown.sh" 2>/dev/null || true
chmod +x "__REMOTE_DIR__/bin" 2>/dev/null || true
CMD
)
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  run_remote "${cmd}"
}

restart_remote_service() {
  run_remote "systemctl restart '${SERVICE_NAME}'"
}

read_remote_port() {
  local cmd
  cmd=$(cat <<'CMD'
port="$(sed -n 's/^[[:space:]]*server_port=\([0-9][0-9]*\).*/\1/p' "__REMOTE_DIR__/conf/application.properties" | head -n 1)"
if [ -z "${port}" ]; then
  port="__DEFAULT_PORT__"
fi
echo "${port}"
CMD
)
  cmd="${cmd//__REMOTE_DIR__/${REMOTE_DIR}}"
  cmd="${cmd//__DEFAULT_PORT__/${DEFAULT_PORT}}"
  run_remote "${cmd}" | tail -n 1 | tr -d '\r'
}

verify_remote_service_status() {
  run_remote "systemctl is-enabled '${SERVICE_NAME}'"
  local cmd
  cmd=$(cat <<'CMD'
set -e
for _ in $(seq 1 30); do
  state="$(systemctl is-active "__SERVICE_NAME__" || true)"
  if [ "${state}" = "active" ]; then
    echo "active"
    exit 0
  fi
  sleep 1
done
echo "service is not active after retries: __SERVICE_NAME__" >&2
systemctl --no-pager -l status "__SERVICE_NAME__" | sed -n '1,80p' >&2 || true
exit 1
CMD
)
  cmd="${cmd//__SERVICE_NAME__/${SERVICE_NAME}}"
  run_remote "${cmd}"
  run_remote "systemctl --no-pager -l status '${SERVICE_NAME}' | sed -n '1,40p'"
}

verify_remote_http() {
  local port="$1"
  local cmd
  cmd=$(cat <<'CMD'
set -e
port="__PORT__"
verify_path="__VERIFY_PATH__"
for _ in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:${port}${verify_path}" >/tmp/ai_study_probe.out 2>/tmp/ai_study_probe.err; then
    echo "PROBE_PATH=${verify_path}"
    head -c 400 /tmp/ai_study_probe.out || true
    echo
    rm -f /tmp/ai_study_probe.out /tmp/ai_study_probe.err
    exit 0
  fi
  if curl -fsS "http://127.0.0.1:${port}/" >/tmp/ai_study_probe.out 2>/tmp/ai_study_probe.err; then
    echo "PROBE_PATH=/"
    head -c 400 /tmp/ai_study_probe.out || true
    echo
    rm -f /tmp/ai_study_probe.out /tmp/ai_study_probe.err
    exit 0
  fi
  sleep 2
done
echo "HTTP check failed on port ${port}, path ${verify_path}" >&2
cat /tmp/ai_study_probe.err >&2 2>/dev/null || true
rm -f /tmp/ai_study_probe.out /tmp/ai_study_probe.err
exit 1
CMD
)
  cmd="${cmd//__PORT__/${port}}"
  cmd="${cmd//__VERIFY_PATH__/${VERIFY_PATH}}"
  run_remote "${cmd}"
}

verify_remote_keepalive() {
  local cmd
  cmd=$(cat <<'CMD'
set -e
before="$(systemctl show "__SERVICE_NAME__" -p MainPID --value)"
if [ -z "${before}" ] || [ "${before}" = "0" ]; then
  echo "service has no main pid before keepalive check" >&2
  exit 1
fi
kill -9 "${before}" || true
after=""
for _ in $(seq 1 30); do
  state="$(systemctl is-active "__SERVICE_NAME__" || true)"
  after="$(systemctl show "__SERVICE_NAME__" -p MainPID --value)"
  if [ "${state}" = "active" ] && [ -n "${after}" ] && [ "${after}" != "0" ] && [ "${after}" != "${before}" ]; then
    echo "KEEPALIVE_OK before=${before} after=${after}"
    exit 0
  fi
  sleep 1
done
echo "keepalive check failed: before=${before} after=${after}" >&2
systemctl --no-pager -l status "__SERVICE_NAME__" | sed -n '1,80p' >&2 || true
exit 1
CMD
)
  cmd="${cmd//__SERVICE_NAME__/${SERVICE_NAME}}"
  run_remote "${cmd}"
}

run_verification() {
  verify_remote_service_status
  local port
  port="$(read_remote_port)"
  echo "Remote port: ${port}"
  verify_remote_http "${port}"
  if [[ "${KEEPALIVE_CHECK}" -eq 1 ]]; then
    verify_remote_keepalive
    verify_remote_service_status
    verify_remote_http "${port}"
  fi
}

if [[ "${VERIFY_ONLY}" -eq 1 ]]; then
  echo "Mode: verify-only"
  run_verification
  exit 0
fi

prepare_local_package
backup_remote_state
remote_conf_before="$(remote_conf_exists)"
remote_db_has_entries="$(remote_database_has_entries)"
echo "REMOTE_CONF_EXISTS=${remote_conf_before}"
echo "REMOTE_DATABASE_HAS_ENTRIES=${remote_db_has_entries}"

sync_package_to_remote
seed_remote_conf_if_missing "${remote_conf_before}"
check_application_properties_diff
copy_local_database_if_needed "${remote_db_has_entries}"
ensure_remote_service
ensure_remote_executable_permissions

if [[ "${SYNC_ONLY}" -eq 1 ]]; then
  echo "Mode: sync-only complete. Service was not restarted."
  exit 0
fi

restart_remote_service
run_verification
