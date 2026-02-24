# EX-PRD-009 — PWA manifest + installability baseline

## Task
Set up/verify web app manifest and metadata for install prompt eligibility.

## Context
- Next.js app router metadata files

## Implementation
1. Add/verify `manifest.webmanifest` with icons/name/theme/start URL.
2. Ensure metadata references manifest correctly.
3. Verify basic install criteria in Chrome DevTools.

## Acceptance checks
- `npm run build`
- Lighthouse (PWA) no manifest-related critical failure

## Out of scope
- Full offline support (next PRDs)
