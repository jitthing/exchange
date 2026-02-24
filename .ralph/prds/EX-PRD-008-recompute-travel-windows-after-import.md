# EX-PRD-008 — Recompute travel windows after import

## Task
Trigger travel-window recomputation after successful ICS import and expose result counts.

## Context
- calendar import endpoint + travel-window logic

## Implementation
1. Add post-import step to recompute windows for affected date range.
2. Return metadata: `windowsCreated`, `windowsUpdated`.
3. Log recompute duration for observability.
4. Add backend test covering import -> recompute call path.

## Acceptance checks
- `cd backend && go test ./...`
- Import response includes recompute metadata

## Out of scope
- Optimizer ranking changes
