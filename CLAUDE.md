# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**All project conventions, commands, and architecture notes live in [AGENTS.md](AGENTS.md) — that file is the single source of truth.** Read it before making changes, and update it (not this file) when conventions change.

Quick orientation:

- Echo **v5** + GORM + PostgreSQL, manual DI in `internal/di/container.go`, layering `handler` → `service` → `repository`.
- Run with `go run cmd/main.go` (needs `.env` with `DATABASE_URL` + `JWT_SECRET`); local services via `docker compose up -d --wait`.
- Test with `go test ./...` (no DB needed); lint with `golangci-lint run`.
- Responses go through `pkg/response`; pagination via `handler.ParsePaginationParams`.

