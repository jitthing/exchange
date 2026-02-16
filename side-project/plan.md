# Exchange Travel Planner (Mobile-First PWA) — Decision-Complete Build Plan

  ## Summary

  Build a mobile-first responsive web app with PWA install support for
  exchange students in Europe (Schengen-focused), initially for you + close
  friends.
  MVP focuses on one tight workflow: find travel windows -> plan weekend trips
  -> track budget impact -> avoid study conflicts.
  Architecture is optimized for a solo builder and fast shipping, while
  keeping clear extension points for social, visa, and alerts features.

  ## Product Scope

  ### In Scope (MVP)

  1. Account + profile (email/social login, cloud sync).
  2. Academic calendar ingestion and travel-window detection.
  3. Weekend Trip Optimizer (constraint-based destination suggestions).
  4. Budget tracker with split categories (travel vs living) and trip-level
     forecasting.
  5. Transport/accommodation comparison via direct provider APIs (Europe-
     first).
  6. Study-Travel balance warnings (assignment/exam conflict checks).
  7. Basic shared trip plan for close friends (view/edit itinerary +
     expenses).
  8. Mobile-first UX + offline-friendly basics (cached key trip details).

  ### Out of Scope (Post-MVP)

  1. Full booking/checkout in-app.
  2. Global provider coverage beyond Europe.
  3. Complex social graph/feed features.
  4. AI chatbot assistant.
  5. Advanced “hidden gems” curation engine at scale.

  ## Users and Jobs-to-be-Done

  1. Exchange student plans short trips around classes and exams.
  2. Student verifies affordability before committing.
  3. Small friend group coordinates itinerary + expense split quickly on
     phone.

  ## Core Feature Specification (MVP)

  ### 1) Academic Calendar Sync + Travel Windows

  1. Inputs:

  - Manual date entry.
  - ICS import (university calendar).
  - Assignment/exam deadlines (manual + optional calendar import).

  2. Logic:

  - Detect long weekends, exam gaps, holiday windows.
  - Mark “safe”, “warning”, “blocked” travel windows.

  3. Output:

  - Window cards with date range, conflict score, confidence.

  ### 2) Weekend Trip Optimizer

  1. Inputs:

  - Budget cap.
  - Departure city.
  - Max travel time.
  - Date window.
  - Trip style (city/nature/nightlife/culture).

  2. Output:

  - Ranked destination list with total estimated cost (transport + stay +
    local daily budget).
  - Reason tags: “cheapest”, “shortest transit”, “best fit for 2 nights”.

  ### 3) Budget Tracker + Forecast

  1. Separate ledgers:

  - Living expenses.
  - Travel expenses.

  2. Trip forecast:

  - Predicted spend before booking.
  - Post-trip reconciliation.

  3. Guardrails:

  - Monthly budget remaining.
  - “If you take this trip, remaining this month = X”.

  ### 4) Transport + Accommodation Finder (Direct APIs)

  1. Transport: train/bus/budget flights (provider set limited to Europe
     launch).
  2. Accommodation: hostel/budget hotel/apartment providers.
  3. Comparison layer:

  - Normalize price, duration, baggage/cancellation metadata where available.

  4. Action:

  - Deep-link out to provider booking page.

  ### 5) Study-Travel Balance Planner

  1. Rule engine:

  - Warn when departure/return overlaps assignment deadlines or exams.
  - Warn on compressed recovery time before key academic events.

  2. Severity:

  - Info / Warning / High-risk flags shown in optimizer and trip detail view.

  ### 6) Lightweight Group Planning

  1. Shared trip object with invited friends.
  2. Editable itinerary blocks.
  3. Expense split summary (equal split in MVP).
  4. Voting deferred to post-MVP.

  ## Information Architecture and Screens

  1. Home (next safe window, budget snapshot, suggested trips).
  2. Calendar (academic + travel overlay).
  3. Discover (optimizer inputs + results).
  4. Trip Detail (transport, stay, itinerary, conflicts, costs).
  5. Budget (monthly overview + trip impact).
  6. Group (shared trip + expense split).
  7. Settings (currency, profile, sync sources, notification prefs).

  ## UX and Mobile Requirements

  1. Mobile-first breakpoints (360px+ baseline).
  2. Thumb-friendly bottom navigation.
  3. PWA:

  - Install prompt.
  - Offline cache for upcoming trips + essential contacts.

  4. Performance target:

  - <2.5s meaningful load on mid-tier mobile over 4G.

  5. Accessibility:

  - WCAG AA contrast, keyboard nav on web, semantic structure.

  ## Technical Architecture

  ### Frontend

  1. Next.js (App Router) + TypeScript.
  2. Tailwind CSS + component primitives.
  3. React Query for server state.
  4. Local-first optimistic UI for budget/itinerary edits; background sync.

  ### Backend

  1. Next.js API routes (initial monolith backend).
  2. PostgreSQL + Prisma ORM.
  3. Redis (optional, phase 2) for API response caching.
  4. Auth: NextAuth (Google + email magic link).

  ### Data Integrations (Europe-first)

  1. Provider adapter layer:

  - TransportProviderAdapter
  - AccommodationProviderAdapter

  2. Each adapter maps external schema -> normalized internal search result.
  3. Rate-limit + retry policy + graceful fallback to partial results.

  ### Notifications

  1. Scheduled checks for:

  - Price drops (if supported by provider endpoints).
  - Study conflict changes.

  2. Channels:

  - Push (PWA) and email digest.

  ## Public APIs / Interfaces / Types (Important)

  ### Core API Endpoints

  1. POST /api/calendar/import
  2. GET /api/travel-windows?from=&to=
  3. POST /api/trips/optimize
  4. GET /api/trips/:tripId
  5. POST /api/trips/:tripId/share
  6. POST /api/budget/entries
  7. GET /api/budget/forecast?tripId=
  8. GET /api/search/transport
  9. GET /api/search/stays
  10. POST /api/conflicts/evaluate

  ### Key Domain Types

  1. TravelWindow { startDate, endDate, score, conflicts[] }
  2. TripConstraint { budgetCap, maxTravelHours, partySize, style, windowId }
  3. TripOption { destination, totalEstimatedCost, transportOptions[],
     stayOptions[], riskLevel }
  4. BudgetEntry { category: living|travel, amount, currency, date, tripId? }
  5. ForecastResult { projectedMonthlySpend, remainingBudget, affordability:
     green|amber|red }
  6. AcademicEvent { type: class|deadline|exam|holiday, start, end, priority }
  7. ConflictAlert { severity, reason, relatedEventId }

  ## Data Model (Initial)

  1. users
  2. profiles
  3. academic_events
  4. travel_windows
  5. trips
  6. trip_members
  7. trip_itinerary_items
  8. budget_entries
  9. trip_forecasts
  10. provider_search_cache
  11. conflict_alerts

  ## Delivery Plan (Solo, Fast, Private Use)

  ### Phase 0 (Week 1): Foundation

  1. Product skeleton, auth, DB schema, PWA shell.
  2. Basic nav + responsive layout.
  3. Seed scripts + dev/staging env.

  ### Phase 1 (Weeks 2-3): Calendar + Budget Core

  1. Calendar import/manual entry.
  2. Travel-window detection.
  3. Budget ledger + categories + monthly summary.

  ### Phase 2 (Weeks 4-5): Optimizer + Provider Search

  1. Constraint form and ranking engine v1.
  2. Transport/stay adapters (minimum 1-2 providers each).
  3. Trip detail with estimated total cost.

  ### Phase 3 (Weeks 6-7): Study Conflict + Group Planning

  1. Conflict warning engine.
  2. Shared trip + expense split.
  3. Notification scaffolding.

  ### Phase 4 (Week 8): Hardening for Real Use

  1. Performance tuning, offline cache, error handling.
  2. End-to-end flow testing on mobile devices.
  3. Private launch for you + friends.

  ## Testing and Acceptance Criteria

  ### Functional Tests

  1. Calendar import correctly creates travel windows for known semester
     fixtures.
  2. Optimizer returns ranked destinations respecting budget/time constraints.
  3. Forecast updates correctly after adding/removing trip costs.
  4. Conflict engine flags overlaps with deadlines/exams.
  5. Shared trip edits are visible to all members.

  ### Integration Tests

  1. Provider adapter normalization handles missing/partial fields.
  2. API failure from one provider does not break full results page.
  3. Auth/session persists across mobile PWA reloads.

  ### UX/Performance Tests

  1. Lighthouse mobile performance >= 80 on key pages.
  2. Main planner actions complete within 3 taps from home.
  3. Offline mode loads upcoming trip summary and emergency notes.

  ### Acceptance Criteria (MVP Done)

  1. User can import academic dates, find a valid weekend window, choose a
     destination option, and see budget impact in one session on mobile.
  2. User gets warning when selected trip conflicts with important academic
     events.
  3. User can share one trip with friends and track split costs.

  ## Risks and Mitigations

  1. Provider API reliability/commercial limits.

  - Mitigation: adapter abstraction, cached fallback, progressive provider
    rollout.

  2. Data normalization inconsistency.

  - Mitigation: strict internal schema + contract tests per adapter.

  3. Scope creep.

  - Mitigation: freeze MVP to planning/budget/calendar + minimal group
    sharing.

  ## Post-MVP Roadmap (Prioritized)

  1. Schengen visa-day tracker + rule checker.
  2. 48-hour city blueprints and hidden gems curation.
  3. Spontaneous trip alerts (price-drop + nearby departure).
  4. Country unlock map.
  5. Offline city packs (maps/phrases/emergency bundle).

  ## Assumptions and Defaults

  1. Primary users are exchange students in Europe.
  2. Product is private/small-scale at launch (you + close friends), no formal
     pilot process.
  3. Platform is responsive web + PWA, not native mobile initially.
  4. Monetization is deferred; no paywall in MVP.
  5. Build is solo-optimized with an 8-week delivery target.
  6. Direct provider APIs are used where feasible; deep-link booking only (no
     in-app transactions).