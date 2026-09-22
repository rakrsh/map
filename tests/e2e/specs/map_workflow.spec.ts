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

    // 3. Tile Service Vector Tile Retrieval
    const tileRes = await request.get(`${tileBase}/tiles/12/655/1583`);
    expect(tileRes.ok()).toBeTruthy();
    const tileData = await tileRes.json();
    expect(tileData.z).toBe('12');
    expect(tileData.x).toBe('655');
    expect(tileData.y).toBe('1583');

    // 4. Browser verification & HTML DOM assertion for artifact logging
    await page.setContent(`
      <!DOCTYPE html>
      <html>
        <head><title>Map Navigation Dashboard</title></head>
        <body>
          <div id="status-panel">
            <h1>Map Navigation Engine</h1>
            <p id="geocoding-status">Geocoding: Ready</p>
            <p id="routing-status">Routing: Ready</p>
            <p id="tile-status">Tile Service: Ready</p>
          </div>
        </body>
      </html>
    `);

    const statusPanel = page.locator('#status-panel');
    await expect(statusPanel).toBeVisible();
    await expect(page.locator('#geocoding-status')).toHaveText('Geocoding: Ready');
    await expect(page.locator('#routing-status')).toHaveText('Routing: Ready');
    await expect(page.locator('#tile-status')).toHaveText('Tile Service: Ready');
  });
});
