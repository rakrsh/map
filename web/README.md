Dev frontend for Map Navigation Engine

This minimal frontend uses Leaflet and the `tile-service` dev PNG endpoint to display a simple map.

Quick start:

1. Ensure `tile-service` is running (via `docker compose up` or `cd services/tile-service && go run ./cmd/server`).
2. Serve the `web/` directory locally (recommended via `make web-dev`):

```bash
make web-dev
```

3. Open `http://localhost:3000` in your browser.

Notes:
- This frontend is intentionally minimal for development demos. It uses CDN-hosted Leaflet and loads tiles from `http://localhost:8082/tiles/{z}/{x}/{y}.png`.
- For production, replace dev tiles with real vector tiles and host static assets using a proper web server.
