import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Calendar & Conflicts', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/calendar');
  });

  test('page loads with travel window selector', async ({ page }) => {
    await expect(page.locator(sel.pageTitle)).toHaveText('Calendar');
    await expect(page.getByLabel('Travel window')).toBeVisible();
  });

  test('travel windows dropdown is populated', async ({ page }) => {
    const select = page.getByLabel('Travel window');
    await expect(select).toBeVisible({ timeout: 10_000 });
    const options = select.locator('option');
    const count = await options.count();
    expect(count).toBeGreaterThan(0);
  });

  test('conflict alerts section is displayed', async ({ page }) => {
    await expect(page.getByText('Conflict Alerts')).toBeVisible({ timeout: 10_000 });
    // Either shows alerts with severity or "No conflicts"
    const hasConflicts = await page.getByText('No conflicts').isVisible().catch(() => false);
    if (!hasConflicts) {
      // There should be at least one alert card
      await expect(page.locator('text=/warning|high-risk|info/i').first()).toBeVisible();
    }
  });
});
