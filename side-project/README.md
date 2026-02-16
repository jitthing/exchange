# Exchange Travel Planner

Mobile-first web app for exchange trip planning with a Go backend API.

## Stack
- Frontend: Next.js 14 + TypeScript + Tailwind CSS
- Backend: Go (net/http), in-memory store for MVP demo

## Run
1. Start backend:
   ```bash
   cd backend
   go run ./cmd/server
   ```
   Backend runs on `http://localhost:8080`.

2. Start frontend in another terminal:
   ```bash
   npm install
   npm run dev
   ```
   Frontend runs on `http://localhost:3000`.

3. Optional API URL override:
   ```bash
   export NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
   ```

## Implemented MVP Areas
- Weekend Trip Optimizer (`POST /api/trips/optimize`)
- Academic Travel Windows (`GET /api/travel-windows`)
- Budget Tracker + Forecast (`/api/budget/entries`, `/api/budget/forecast`)
- Transport/Stay search adapters (`/api/search/transport`, `/api/search/stays`)
- Study-travel conflict checks (`POST /api/conflicts/evaluate`)
- Group trip share flow (`GET /api/trips/:id`, `POST /api/trips/:id/share`)
- Mobile-friendly screens for Home, Calendar, Discover, Budget, Group, Settings, Trip Detail
- PWA manifest and install metadata

## Notes
- Current backend storage is in-memory and resets on restart.
- Integrations use deterministic mock provider data and deep links.
