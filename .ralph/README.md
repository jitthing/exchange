# Ralph Loop PRDs (Exchange)

Use these PRDs as atomic loop items. Target: one PRD = one small PR (30–90 mins).

## Execution rules
- Pick next unchecked PRD from `PRD-QUEUE.md`
- Do only that PRD
- Run required checks
- Update `PRD-QUEUE.md` status + notes
- Open PR with title: `[RALPH] <PRD-ID> <short title>`

## Global constraints
- Keep API/provider keys server-side only
- No UI redesign unless explicitly requested by PRD
- Do not break existing pages: Home, Calendar, Discover, Trip Detail, Budget, Group
- Keep fallback/demo behavior if real API is unavailable
