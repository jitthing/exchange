# EX-PRD-004 — Discover page wiring to backend search API

## Task
Wire Discover page search action to backend transport API with loading + error + empty states.

## Context
- `app/discover/page.tsx`
- backend transport endpoint from EX-PRD-002/003

## Implementation
1. Add client fetch helper for transport search.
2. Trigger request from Discover form submit.
3. Render loading skeleton, error banner, empty state, success list.
4. Keep current UI structure; minimal visual change.

## Acceptance checks
- `npm run typecheck`
- `npm run lint`
- Manual test on Discover page: loading/error/success states all observable

## Out of scope
- Advanced filtering/sorting UI
