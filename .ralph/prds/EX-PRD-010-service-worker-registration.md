# EX-PRD-010 — Service worker registration (safe baseline)

## Task
Register service worker in production build only, with safe fallback if unsupported.

## Context
- frontend app bootstrap/layout

## Implementation
1. Add SW registration script/client hook.
2. Gate to production + browser support check.
3. Add verbose console logs only in development.

## Acceptance checks
- `npm run build`
- In prod preview, SW registers successfully
- In unsupported context, app still works normally

## Out of scope
- Cache rules implementation
