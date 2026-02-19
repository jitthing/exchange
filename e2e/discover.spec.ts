import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Trip Discovery Flow', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/discover');
  });

  test('page loads with optimizer form', async ({ page }) => {
    await expect(page.locator(sel.pageTitle)).toHaveText('Discover');
    await expect(page.locator(sel.discoverForm)).toBeVisible();
  });

  test('form has departure city, budget cap, max travel hours fields', async ({ page }) => {
    await expect(page.getByLabel('Departure city')).toBeVisible();
    await expect(page.getByLabel('Budget cap (€)')).toBeVisible();
    await expect(page.getByLabel('Max travel hours')).toBeVisible();
  });

  test('validation: empty departure city shows error', async ({ page }) => {
    await page.getByLabel('Departure city').clear();
    await page.locator(sel.discoverSubmit).click();
    await expect(page.locator(sel.fieldError).first()).toContainText('Departure city is required');
  });

  test('validation: budget cap <= 0 shows error', async ({ page }) => {
    await page.getByLabel('Budget cap (€)').fill('0');
    await page.locator(sel.discoverSubmit).click();
    await expect(page.locator(sel.fieldError).first()).toBeVisible();
  });

  test('can fill form and submit — results appear', async ({ page }) => {
    await page.getByLabel('Departure city').fill('Berlin');
    await page.getByLabel('Budget cap (€)').fill('300');
    await page.getByLabel('Max travel hours').fill('5');
    await page.locator(sel.discoverSubmit).click();
    await expect(page.locator(sel.tripResults).locator('h2').first()).toBeVisible({ timeout: 10_000 });
  });

  test('result cards show destination, cost, and reason tags', async ({ page }) => {
    await page.getByLabel('Departure city').fill('Berlin');
    await page.locator(sel.discoverSubmit).click();
    const results = page.locator(sel.tripResults);
    await expect(results.locator('h2').first()).toBeVisible({ timeout: 10_000 });
    // Cost: €xxx
    await expect(results.locator('text=/€\\d+/').first()).toBeVisible();
  });
});
