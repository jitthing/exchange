import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Home Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/');
  });

  test('shows greeting and title', async ({ page }) => {
    await expect(page.getByText('Good to see you')).toBeVisible();
    await expect(page.locator(sel.pageTitle)).toHaveText('Travel Planner');
  });

  test('KPI cards are visible', async ({ page }) => {
    const cards = page.locator(sel.kpiCard);
    await expect(cards).toHaveCount(2);
    await expect(cards.first()).toContainText('Projected');
    await expect(cards.last()).toContainText('Remaining');
  });

  test('travel window highlight card is shown', async ({ page }) => {
    await expect(page.getByText('Next travel window')).toBeVisible();
  });

  test('quick action buttons are present and clickable', async ({ page }) => {
    const findTrips = page.getByText('🧭 Find Trips');
    const calendar = page.getByText('📅 Calendar');
    await expect(findTrips).toBeVisible();
    await expect(calendar).toBeVisible();
    await findTrips.click();
    await expect(page).toHaveURL(/\/discover/);
  });
});
