import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Responsive Design', () => {
  test('renders correctly at mobile viewport (375px)', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await navigateTo(page, '/');
    await expect(page.locator(sel.pageTitle)).toBeVisible();
    await expect(page.locator(sel.navHome)).toBeVisible();
    // Check no horizontal overflow
    const body = page.locator('body');
    const box = await body.boundingBox();
    expect(box!.width).toBeLessThanOrEqual(375);
  });

  test('renders correctly at desktop viewport (1280px)', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 720 });
    await navigateTo(page, '/');
    await expect(page.locator(sel.pageTitle)).toBeVisible();
  });

  test('bottom nav is visible on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await navigateTo(page, '/');
    await expect(page.locator(sel.navHome)).toBeVisible();
    await expect(page.locator(sel.navCalendar)).toBeVisible();
    await expect(page.locator(sel.navDiscover)).toBeVisible();
    await expect(page.locator(sel.navBudget)).toBeVisible();
    await expect(page.locator(sel.navGroup)).toBeVisible();
  });
});
