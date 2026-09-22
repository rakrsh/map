// Copyright (c) 2026 Ravi Sharma

import { test, expect } from '@playwright/test';

test.describe('Microservices API Contracts & Health Checks', () => {
  const geocodingBase = process.env.GEOCODING_URL || 'http://localhost:8000';
  const routingBase = process.env.ROUTING_URL || 'http://localhost:8081';
  const tileBase = process.env.TILE_URL || 'http://localhost:8082';

  test('all services expose healthy status', async ({ request }) => {
    const geoHealth = await request.get(`${geocodingBase}/health`);
    expect(geoHealth.ok()).toBeTruthy();
    expect(await geoHealth.json()).toEqual({ status: 'ok' });

    const routeHealth = await request.get(`${routingBase}/health`);
    expect(routeHealth.ok()).toBeTruthy();
    expect(await routeHealth.json()).toEqual({ status: 'ok' });

    const tileHealth = await request.get(`${tileBase}/health`);
    expect(tileHealth.ok()).toBeTruthy();
    expect(await tileHealth.json()).toEqual({ status: 'ok' });
  });

  test('geocoding API handles reverse lookup validation', async ({ request }) => {
    const res = await request.get(`${geocodingBase}/api/v1/reverse`, {
      params: { lat: '37.7749', lon: '-122.4194' },
    });
    expect(res.ok()).toBeTruthy();
    const data = await res.json();
    expect(data.latitude).toBe(37.7749);
    expect(data.longitude).toBe(-122.4194);
  });

  test('tile service rejects unsupported HTTP methods', async ({ request }) => {
    const res = await request.post(`${tileBase}/health`);
    expect(res.status()).toBe(405);
  });
});
