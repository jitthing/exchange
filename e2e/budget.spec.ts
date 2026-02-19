import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Budget Tracking Flow', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/budget');
  });

  test('page loads with forecast section', async ({ page }) => {
    await expect(page.locator(sel.pageTitle)).toHaveText('Budget');
    await expect(page.getByText('Projected')).toBeVisible({ timeout: 10_000 });
    await expect(page.getByText('Remaining')).toBeVisible();
    await expect(page.getByText('Status')).toBeVisible();
  });

  test('add entry form is present with fields', async ({ page }) => {
    await expect(page.locator(sel.budgetForm)).toBeVisible();
    await expect(page.getByLabel('Category')).toBeVisible();
    await expect(page.getByLabel('Amount (€)')).toBeVisible();
    await expect(page.getByLabel('Date')).toBeVisible();
  });

  test('validation: amount <= 0 shows error', async ({ page }) => {
    await page.getByLabel('Amount (€)').fill('0');
    await page.locator(sel.budgetSubmit).click();
    await expect(page.locator(sel.fieldError).first()).toBeVisible();
  });

  test('can add a budget entry', async ({ page }) => {
    await page.getByLabel('Amount (€)').fill('50');
    await page.locator(sel.budgetSubmit).click();
    // Entry should appear in the entries list
    await expect(page.locator(sel.budgetEntries)).toContainText('€50', { timeout: 10_000 });
  });
});
