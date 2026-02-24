# EX-PRD-012 — Offline fallback page + UX messaging

## Task
Create friendly offline fallback experience for unavailable pages/actions.

## Context
- app routes + service worker fallback

## Implementation
1. Add `/offline` page with actionable guidance.
2. Route failed navigations to offline page when network unavailable.
3. Add non-intrusive banner/toast when app detects offline mode.

## Acceptance checks
- `npm run typecheck`
- `npm run lint`
- Offline mode shows fallback page + clear messaging

## Out of scope
- Full offline data editing/sync queue
