import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo } from './helpers/setup';

test.describe('Group Planning', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/group');
  });

  test('page loads', async ({ page }) => {
    await expect(page.locator(sel.pageTitle)).toHaveText('Group');
  });

  test('can load a shared trip', async ({ page }) => {
    await page.getByRole('button', { name: 'Load Shared Trip' }).click();
    await expect(page.locator('h2').first()).toBeVisible({ timeout: 10_000 });
  });

  test('trip details show destination and cost', async ({ page }) => {
    await page.getByRole('button', { name: 'Load Shared Trip' }).click();
    await expect(page.locator('text=/€\\d+/')).toBeVisible({ timeout: 10_000 });
  });

  test('members list is displayed after loading', async ({ page }) => {
    await page.getByRole('button', { name: 'Load Shared Trip' }).click();
    await expect(page.getByText('Members', { exact: true })).toBeVisible({ timeout: 10_000 });
  });

  test('can add a member', async ({ page }) => {
    await page.getByRole('button', { name: 'Load Shared Trip' }).click();
    await expect(page.getByText('Members', { exact: true })).toBeVisible({ timeout: 10_000 });
    const initialMembers = await page.locator('text=/👤/').count();
    const uniqueMember = `friend-${Date.now()}`;
    await page.getByPlaceholder('Enter member ID').fill(uniqueMember);
    await page.getByRole('button', { name: 'Add' }).click();
    await expect(page.getByText(uniqueMember)).toBeVisible({ timeout: 10_000 });
  });
});
