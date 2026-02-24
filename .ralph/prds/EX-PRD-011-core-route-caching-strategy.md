# EX-PRD-011 — Cache strategy for core app shell routes

## Task
Implement conservative caching for essential shell/static assets and key route shells.

## Context
- service worker from EX-PRD-010

## Implementation
1. Cache static assets with cache-first strategy.
2. Cache route shells with network-first + offline fallback.
3. Keep cache versioning and cleanup on activate.

## Acceptance checks
- `npm run build`
- Manual offline test: app shell for Home/Calendar/Budget opens offline

## Out of scope
- Dynamic API response caching
