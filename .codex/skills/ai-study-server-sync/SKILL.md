---
name: ai-study-server-sync
description: Deploy and verify the AI Study backend/admin service to the remote server under /data/ai-study. Use when Codex needs to build the Go backend, sync backend runtime files, preserve remote conf/application.properties, database, and file uploads, manage ai-study.service, restart the service, or verify the live /collect-ui admin page on the target host.
---

# AI Study Server Sync

Use this skill for repeatable backend/admin deployment to the target server.

## Fixed Environment

- Default host: `202.140.140.117`
- Default user: `root`
- Remote project dir: `/data/ai-study`
- Local backend dir: `/data/project/ai-study/backend`
- Target managed service: `/etc/systemd/system/ai-study.service`
- Default backend port: `8026`

## Environment Variables

- `AI_STUDY_SERVER_HOST` - optional, defaults to `202.140.140.117`
- `AI_STUDY_SERVER_USER` - optional, defaults to `root`
- `AUTO_CHECK_SERVER_PASSWORD` - preferred password variable
- `TEST_SERVER_PASSWORD` - fallback password variable
- `AI_STUDY_SERVER_REMOTE_DIR` - optional, defaults to `/data/ai-study`
- `AI_STUDY_SERVER_LOCAL_BACKEND_DIR` - optional, defaults to `/data/project/ai-study/backend`
- `AI_STUDY_SERVER_BACKUP_DIR` - optional, defaults to `/data/ai-study-backup`
- `AI_STUDY_SERVER_SERVICE_NAME` - optional, defaults to `ai-study.service`
- `AI_STUDY_SERVER_VERIFY_PATH` - optional, defaults to `/collect-ui/`
- `AI_STUDY_SERVER_SEED_CONF_IF_MISSING` - optional, defaults to `1`

Never write SSH passwords into docs, code, logs, or command history examples. Pass them through `AUTO_CHECK_SERVER_PASSWORD` or the existing `TEST_SERVER_PASSWORD`.

## Workflow

1. If backend Go code changed, run the normal backend checks first:
```bash
cd /data/project/ai-study/backend
go fmt ./...
go test ./...
go vet ./...
```

2. If the admin frontend changed, build/deploy the AI Study collect-ui static shell into `backend/frontend/collect-ui` first:
```bash
cd /data/project/sport-ui
bash scripts/deploy_ai_study_collect_ui.sh
```

3. Export the password variable without putting the value in files:
```bash
read -r -s -p "AUTO_CHECK_SERVER_PASSWORD: " AUTO_CHECK_SERVER_PASSWORD
export AUTO_CHECK_SERVER_PASSWORD
printf '\n'
```

4. Run full backend deployment:
```bash
cd /data/project/ai-study
bash .codex/skills/ai-study-server-sync/scripts/sync_restart_verify.sh
```

The script will:

- build a Linux amd64 backend binary locally
- prepare a minimal runtime package from `backend/`
- sync it to remote `/data/ai-study`
- preserve remote `conf/application.properties`, `database/`, `file_data/`, logs, and PID files
- seed `conf/application.properties` only when the remote file is missing and `AI_STUDY_SERVER_SEED_CONF_IF_MISSING=1`
- backup remote `conf/application.properties` and `database/` before deploying
- copy local `database/` only when the remote database directory is empty
- create/update and enable `ai-study.service`
- restart the service and verify systemd status plus the local `/collect-ui/` page on the remote host

## Commands

- Full sync + restart + verify:
```bash
bash .codex/skills/ai-study-server-sync/scripts/sync_restart_verify.sh
```

- Sync package and service file, but do not restart:
```bash
bash .codex/skills/ai-study-server-sync/scripts/sync_restart_verify.sh --sync-only
```

- Verify existing remote service only:
```bash
bash .codex/skills/ai-study-server-sync/scripts/sync_restart_verify.sh --verify-only
```

- Include an explicit keepalive restart check after health verification:
```bash
bash .codex/skills/ai-study-server-sync/scripts/sync_restart_verify.sh --keepalive-check
```

## Notes

- Do not hand-edit `backend/frontend/collect-ui/assets` hash files. Rebuild from `/data/project/collect-ui` and `/data/project/sport-ui`, then deploy this backend package.
- The script intentionally does not overwrite an existing remote `conf/application.properties`; it only prints missing key names for manual review.
- Remote backup snapshots are stored under `/data/ai-study-backup` by default.
- The deploy target should be treated as owned by this service. Files under `/data/ai-study` that are not in the package and not explicitly excluded may be removed by rsync.
