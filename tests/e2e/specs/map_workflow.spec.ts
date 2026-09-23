// Copyright (c) 2026 Ravi Sharma

import { test, expect } from '@playwright/test';

test.describe('Map Navigation End-to-End Workflow', () => {
  const geocodingBase = process.env.GEOCODING_URL || 'http://localhost:8000';
  const routingBase = process.env.ROUTING_URL || 'http://localhost:8081';
  const tileBase = process.env.TILE_URL || 'http://localhost:8082';

  test('complete map user navigation workflow: search -> route -> tile retrieval', async ({ request, page }) => {
    // 1. Geocoding Search
    const searchRes = await request.get(`${geocodingBase}/api/v1/search`, {
      params: { query: 'Market Street, San Francisco' },
    });
    expect(searchRes.ok()).toBeTruthy();
    const searchData = await searchRes.json();
    expect(searchData.query).toBe('Market Street, San Francisco');
    expect(searchData.status).toBe('pending_implementation');

    // 2. Routing Path Request
    const routeRes = await request.post(`${routingBase}/routes`, {
      data: {
        origin: [37.7749, -122.4194],
        destination: [37.7833, -122.4167],
        preference: 'low_stress',
      },
    });
    expect(routeRes.ok()).toBeTruthy();
    const routeData = await routeRes.json();
    expect(routeData.message).toBe('route request received');
    expect(routeData.status).toBe('pending_implementation');

    // 3. Tile Service PNG tile retrieval (dev frontend)
    const pngRes = await request.get(`${tileBase}/tiles/12/655/1583.png`);
    expect(pngRes.ok()).toBeTruthy();
    expect(pngRes.headers()['content-type']).toContain('image/png');

    // 4. Load the real frontend and verify a Leaflet tile image is present
    await page.goto(process.env.FRONTEND_URL || 'http://localhost:3000');
    // Wait for the tile image network requests to be attempted
    await page.waitForTimeout(500);

    // Check that at least one tile <img> element is present in the Leaflet pane
    const tileImg = page.locator('.leaflet-tile');
    await expect(tileImg.first()).toBeVisible();

    // Verify marker popup or marker exists (marker uses default Leaflet icon)
    const marker = page.locator('.leaflet-marker-icon');
    await expect(marker.first()).toBeVisible();
  });
});
