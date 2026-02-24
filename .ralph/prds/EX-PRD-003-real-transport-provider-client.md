# EX-PRD-003 — One real transport provider client (happy path)

## Task
Implement one real transport provider adapter and map response to normalized schema.

## Context
- Adapter from EX-PRD-002
- Config from EX-PRD-001

## Implementation
1. Add provider client with timeout + retry (small, bounded).
2. Map provider payload -> normalized `TransportOption`.
3. If provider fails, return structured error and keep mock fallback optional via flag.
4. Add unit test for mapping with sample fixture JSON.

## Acceptance checks
- `cd backend && go test ./...`
- With valid credentials, endpoint returns real provider data in normalized shape
- With provider failure, endpoint returns controlled error (no panic)

## Out of scope
- Multi-provider aggregation
- Caching layer
