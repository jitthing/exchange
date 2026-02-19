import { test, expect } from '@playwright/test';
import { sel } from './helpers/selectors';
import { navigateTo, clickNav } from './helpers/setup';

test.describe('Core Navigation', () => {
  test('app loads and shows home page', async ({ page }) => {
    await navigateTo(page, '/');
    await expect(page.locator(sel.pageTitle)).toHaveText('Travel Planner');
  });

  test('bottom nav is visible with 5 tabs', async ({ page }) => {
    await navigateTo(page, '/');
    await expect(page.locator(sel.navHome)).toBeVisible();
    await expect(page.locator(sel.navCalendar)).toBeVisible();
    await expect(page.locator(sel.navDiscover)).toBeVisible();
    await expect(page.locator(sel.navBudget)).toBeVisible();
    await expect(page.locator(sel.navGroup)).toBeVisible();
  });

  test('can navigate to Calendar via bottom nav', async ({ page }) => {
    await navigateTo(page, '/');
    await clickNav(page, 'nav-calendar');
    await expect(page).toHaveURL(/\/calendar/);
    await expect(page.locator(sel.pageTitle)).toHaveText('Calendar');
  });

  test('can navigate to Discover via bottom nav', async ({ page }) => {
    await navigateTo(page, '/');
    await clickNav(page, 'nav-discover');
    await expect(page).toHaveURL(/\/discover/);
    await expect(page.locator(sel.pageTitle)).toHaveText('Discover');
  });

  test('can navigate to Budget via bottom nav', async ({ page }) => {
    await navigateTo(page, '/');
    await clickNav(page, 'nav-budget');
    await expect(page).toHaveURL(/\/budget/);
    await expect(page.locator(sel.pageTitle)).toHaveText('Budget');
  });

  test('can navigate to Group via bottom nav', async ({ page }) => {
    await navigateTo(page, '/');
    await clickNav(page, 'nav-group');
    await expect(page).toHaveURL(/\/group/);
    await expect(page.locator(sel.pageTitle)).toHaveText('Group');
  });
});
