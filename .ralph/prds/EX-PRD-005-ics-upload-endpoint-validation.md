# EX-PRD-005 — ICS upload endpoint scaffold + validation

## Task
Create backend endpoint to receive ICS file and validate size/type/basic format.

## Context
- `backend/internal/httpapi/router.go`

## Implementation
1. Add POST endpoint `/api/calendar/import` (multipart/form-data).
2. Validate file exists, size limit (e.g., <= 2MB), and `.ics` mime/name check.
3. Return structured JSON validation errors.
4. Stub parser call (implemented in EX-PRD-006).

## Acceptance checks
- `cd backend && go test ./...`
- Invalid file -> 4xx + clear error JSON
- Valid-looking `.ics` file -> parser stub success response

## Out of scope
- Persisting events
