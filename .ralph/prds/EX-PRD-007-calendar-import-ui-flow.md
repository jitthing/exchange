# EX-PRD-007 — Calendar import UI flow

## Task
Add frontend upload flow to import ICS and show summary of imported events.

## Context
- `app/calendar/page.tsx`
- backend `/api/calendar/import`

## Implementation
1. Add file input + import button in Calendar page.
2. POST multipart ICS file to backend endpoint.
3. Display summary: imported count, skipped count, warnings.
4. Show inline error state for invalid file/server failures.

## Acceptance checks
- `npm run typecheck`
- `npm run lint`
- Manual test with valid + invalid ICS files

## Out of scope
- Fancy drag/drop uploader
