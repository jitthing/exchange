/** Shared test selectors using data-testid */
export const sel = {
  // Bottom nav
  navHome: '[data-testid="nav-home"]',
  navCalendar: '[data-testid="nav-calendar"]',
  navDiscover: '[data-testid="nav-discover"]',
  navBudget: '[data-testid="nav-budget"]',
  navGroup: '[data-testid="nav-group"]',

  // Common
  pageTitle: '[data-testid="page-title"]',
  fieldError: '[data-testid="field-error"]',
  kpiCard: '[data-testid="kpi-card"]',

  // Discover
  discoverForm: '[data-testid="discover-form"]',
  discoverSubmit: '[data-testid="discover-submit"]',
  tripResults: '[data-testid="trip-results"]',

  // Budget
  budgetForm: '[data-testid="budget-form"]',
  budgetSubmit: '[data-testid="budget-submit"]',
  budgetEntries: '[data-testid="budget-entries"]',
} as const;
