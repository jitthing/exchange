import { API_BASE_URL } from '@/lib/config';
import { BudgetEntry, ConflictAlert, ForecastResult, TravelWindow, Trip, TripConstraint, TripOption } from '@/lib/types';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers || {})
    },
    cache: 'no-store'
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Request failed (${response.status})`);
  }

  return response.json() as Promise<T>;
}

export function getTravelWindows(from?: string, to?: string): Promise<{ windows: TravelWindow[] }> {
  const params = new URLSearchParams();
  if (from) params.set('from', from);
  if (to) params.set('to', to);
  const query = params.toString();
  return request<{ windows: TravelWindow[] }>(`/api/travel-windows${query ? `?${query}` : ''}`);
}

export function optimizeTrips(payload: TripConstraint): Promise<{ options: TripOption[] }> {
  return request<{ options: TripOption[] }>('/api/trips/optimize', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function getBudgetEntries(userId = 'demo-user'): Promise<{ entries: BudgetEntry[] }> {
  return request<{ entries: BudgetEntry[] }>(`/api/budget/entries?userId=${encodeURIComponent(userId)}`);
}

export function createBudgetEntry(payload: Omit<BudgetEntry, 'id'>): Promise<BudgetEntry> {
  return request<BudgetEntry>('/api/budget/entries', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function getForecast(tripId?: string, userId = 'demo-user'): Promise<ForecastResult> {
  const params = new URLSearchParams({ userId });
  if (tripId) params.set('tripId', tripId);
  return request<ForecastResult>(`/api/budget/forecast?${params.toString()}`);
}

export function getTrip(tripId: string): Promise<Trip> {
  return request<Trip>(`/api/trips/${tripId}`);
}

export function shareTrip(tripId: string, memberIds: string[]): Promise<Trip> {
  return request<Trip>(`/api/trips/${tripId}/share`, {
    method: 'POST',
    body: JSON.stringify({ memberIds })
  });
}

export function getTransport(from: string, to: string): Promise<{ options: TripOption['transportOptions'] }> {
  const params = new URLSearchParams({ from, to });
  return request<{ options: TripOption['transportOptions'] }>(`/api/search/transport?${params.toString()}`);
}

export function getStays(city: string): Promise<{ options: TripOption['stayOptions'] }> {
  return request<{ options: TripOption['stayOptions'] }>(`/api/search/stays?city=${encodeURIComponent(city)}`);
}

export function evaluateConflicts(windowId: string): Promise<{ alerts: ConflictAlert[] }> {
  return request<{ alerts: ConflictAlert[] }>('/api/conflicts/evaluate', {
    method: 'POST',
    body: JSON.stringify({ windowId })
  });
}
