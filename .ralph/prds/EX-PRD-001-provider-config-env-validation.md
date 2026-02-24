# EX-PRD-001 — Provider config + env validation layer

## Task
Create a typed provider config module that validates required env vars at startup/runtime.

## Context
- `backend/cmd/server/main.go`
- `backend/internal/...` (new `config` package allowed)

## Implementation
1. Add `backend/internal/config/providers.go`.
2. Define typed config struct for transport provider credentials/base URL.
3. Add validation helper returning explicit errors for missing vars.
4. Wire config load in server bootstrap; fail gracefully with clear log message.

## Acceptance checks
- `cd backend && go test ./...`
- Start server with missing env var -> readable validation error
- Start server with env vars set -> boots normally

## Out of scope
- Real API call
- Frontend changes
