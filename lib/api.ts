import { API_BASE_URL } from '@/lib/config';
import { supabase } from '@/lib/supabase';
import { BudgetEntry, ConflictAlert, ForecastResult, TravelWindow, Trip, TripConstraint, TripOption } from '@/lib/types';

async function getAuthHeaders(): Promise<Record<string, string>> {
  if (!supabase) return {};
  const { data: { session } } = await supabase.auth.getSession();
  if (session?.access_token) {
    return { Authorization: `Bearer ${session.access_token}` };
  }
  return {};
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const authHeaders = await getAuthHeaders();
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...authHeaders,
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

export function getBudgetEntries(): Promise<{ entries: BudgetEntry[] }> {
  return request<{ entries: BudgetEntry[] }>('/api/budget/entries');
}

export function createBudgetEntry(payload: Omit<BudgetEntry, 'id'>): Promise<BudgetEntry> {
  return request<BudgetEntry>('/api/budget/entries', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function getForecast(tripId?: string): Promise<ForecastResult> {
  const params = new URLSearchParams();
  if (tripId) params.set('tripId', tripId);
  const query = params.toString();
  return request<ForecastResult>(`/api/budget/forecast${query ? `?${query}` : ''}`);
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
