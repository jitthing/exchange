import { test, expect } from '@playwright/test';
import { navigateTo } from './helpers/setup';

test.describe('Trip Detail Page', () => {
  test('loads for a valid trip', async ({ page }) => {
    await navigateTo(page, '/trips/trip-1');
    await expect(page.locator('[data-testid="page-title"]')).toBeVisible();
  });

  test('shows destination and estimated cost', async ({ page }) => {
    await navigateTo(page, '/trips/trip-1');
    await expect(page.getByText('Estimated Cost')).toBeVisible();
    await expect(page.locator('text=/€\\d+/')).toBeVisible();
  });

  test('itinerary items are listed', async ({ page }) => {
    await navigateTo(page, '/trips/trip-1');
    await expect(page.getByRole('heading', { name: 'Itinerary' })).toBeVisible();
  });

  test('members are shown', async ({ page }) => {
    await navigateTo(page, '/trips/trip-1');
    await expect(page.getByText('Members')).toBeVisible();
    await expect(page.locator('text=/👤/').first()).toBeVisible();
  });
});
