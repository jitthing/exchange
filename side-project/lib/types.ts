export type Severity = 'info' | 'warning' | 'high-risk';

export type AcademicEventType = 'class' | 'deadline' | 'exam' | 'holiday';

export type TravelWindow = {
  id: string;
  startDate: string;
  endDate: string;
  score: number;
  conflicts: string[];
};

export type AcademicEvent = {
  id: string;
  type: AcademicEventType;
  title: string;
  start: string;
  end: string;
  priority: number;
};

export type TripConstraint = {
  budgetCap: number;
  maxTravelHours: number;
  partySize: number;
  style: 'city' | 'nature' | 'nightlife' | 'culture';
  windowId: string;
  departureCity: string;
};

export type TransportOption = {
  provider: string;
  mode: 'train' | 'bus' | 'flight';
  durationHours: number;
  price: number;
  deeplink: string;
};

export type StayOption = {
  provider: string;
  kind: 'hostel' | 'budget-hotel' | 'apartment';
  nightlyPrice: number;
  rating: number;
  deeplink: string;
};

export type ConflictAlert = {
  severity: Severity;
  reason: string;
  relatedEventId: string;
};

export type TripOption = {
  id: string;
  destination: string;
  reasonTags: string[];
  totalEstimatedCost: number;
  transportOptions: TransportOption[];
  stayOptions: StayOption[];
  riskLevel: Severity;
};

export type Trip = {
  id: string;
  ownerId: string;
  destination: string;
  windowId: string;
  members: string[];
  itinerary: string[];
  estimatedCost: number;
};

export type BudgetCategory = 'living' | 'travel';

export type BudgetEntry = {
  id: string;
  userId: string;
  category: BudgetCategory;
  amount: number;
  currency: string;
  date: string;
  tripId?: string;
  note?: string;
};

export type ForecastResult = {
  projectedMonthlySpend: number;
  remainingBudget: number;
  affordability: 'green' | 'amber' | 'red';
};
