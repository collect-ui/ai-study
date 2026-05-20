# AGENTS.md

Guidance for coding agents working in `/data/project/ai-study/backend`.

This backend is the AI Study admin service. It is copied from the collect-ui low-code backend architecture and trimmed to the foundation needed by the mini program admin console.

## Scope

- Keep login, current user, logout, user, role, user-role, menu, role-menu, sys-code, sys-btn, schema-page, schema-page-field, schema-page-data, and frontend page-data services.
- Do not reintroduce old reference-project business modules.
- Apart from model definitions and table registration, do not add or modify Go business code without approval. If Go code is required, first explain what will be written, why it is necessary, and why existing low-code config or services cannot replace it.
- Do not store production passwords, tokens, or SSH secrets in code, docs, logs, or config.
- Runtime request logs are written to `logs/server.log`.

## Commands

```bash
go test ./...
go build ./...
go run .
```

Reset the local SQLite seed data:

```bash
sqlite3 database/ai_study_admin.db ".read scripts/init-ai-study-admin.sql"
```

Default local admin account:

```text
username: admin
password: 123456
```

Change the default password before production use.
