import { Page } from '@playwright/test';

/** Navigate to a page and wait for it to be ready */
export async function navigateTo(page: Page, path: string) {
  await page.goto(path);
  await page.waitForLoadState('networkidle');
}

/** Click a bottom nav tab and wait for navigation */
export async function clickNav(page: Page, testId: string) {
  await page.click(`[data-testid="${testId}"]`);
  await page.waitForLoadState('networkidle');
}
