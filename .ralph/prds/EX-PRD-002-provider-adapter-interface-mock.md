# EX-PRD-002 — Provider adapter interface + mock adapter

## Task
Introduce adapter interface for transport search and add an in-memory mock implementation.

## Context
- `backend/internal/domain/types.go`
- `backend/internal/httpapi/router.go`

## Implementation
1. Define `TransportProviderAdapter` interface.
2. Define normalized `TransportOption` type (price, duration, currency, deeplink, provider).
3. Implement mock adapter returning deterministic sample results.
4. Expose a backend endpoint (or internal handler path) that uses adapter output.

## Acceptance checks
- `cd backend && go test ./...`
- Endpoint returns normalized JSON array with >= 2 items from mock adapter

## Out of scope
- Real provider auth/network calls
