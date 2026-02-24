# EX-PRD-006 — ICS parser VEVENT normalization

## Task
Parse ICS VEVENT blocks and normalize to internal `AcademicEvent` objects.

## Context
- New parser package under `backend/internal/calendar` allowed
- Domain types in `backend/internal/domain/types.go`

## Implementation
1. Parse DTSTART/DTEND/SUMMARY (timezone-safe basic handling).
2. Infer event type: class|deadline|exam|holiday via keyword rules.
3. Return normalized events list + skipped-count + warnings.
4. Add parser unit tests using 2 fixture ICS samples.

## Acceptance checks
- `cd backend && go test ./...`
- Parser test covers malformed VEVENT and still returns partial success

## Out of scope
- Recurrence rules perfection
- Calendar provider syncing
