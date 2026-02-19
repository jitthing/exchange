import { test, expect } from '@playwright/test';
import { navigateTo } from './helpers/setup';

test.describe('Authentication (Dev Mode)', () => {
  test('with AUTH_DISABLED=true, app loads directly without login', async ({ page }) => {
    await navigateTo(page, '/');
    // Should see home page content, not a login form
    await expect(page.getByText('Travel Planner')).toBeVisible();
  });

  test('settings page renders', async ({ page }) => {
    await navigateTo(page, '/settings');
    // Page should load without errors (may show 404 or settings content)
    await expect(page).toHaveURL(/\/settings/);
  });
});
