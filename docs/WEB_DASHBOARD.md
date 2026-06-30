# PQR Web Dashboard (web/dashboard.html)

This document explains the **front-end web dashboard** shipped in `/web/`.

## Files
- `web/dashboard.html` – Dashboard page entry
- `web/dashboard.js` – Client-side logic
- `web/dashboard.css` – Dashboard styling
- `web/assets/` – Images and other static assets

## How to run it locally

### Option A: Open directly (quickest)
1. Open `web/dashboard.html` in a browser.
2. If the browser blocks file-based requests, use Option B.

### Option B: Serve the `web/` folder (recommended)
Run a static server in the repo root, then open:
- `http://localhost:<port>/web/dashboard.html`

(Use any static server; examples below are left to your environment.)

## Page behavior
The dashboard is designed to:
- Display current system state (where the UI is wired to your API endpoints)
- Provide navigation to the HUD and wiki pages
- Surface key “health” style signals via client-side fetch calls

> Note: The exact API URLs used are determined by `web/dashboard.js`. If your backend runs at a different origin, update the base URL in `web/dashboard.js` (look for configuration variables).

## Customization points
1. **Base API URL**
   - Search in `web/dashboard.js` for the variable that sets the API origin.
2. **UI endpoints**
   - The UI will fetch health/status and any other endpoints referenced in `dashboard.js`.
3. **Styling**
   - Modify `web/dashboard.css` for layout/typography. The UI is meant to be simple and readable.

## Known limitations
- If you open HTML files directly from disk, some browsers restrict `fetch()`/CORS behavior.
- For cross-origin setups, ensure the PQR API server allows the web origin, or proxy requests.

## Related docs
- [docs/WEB_HUD.md](docs/WEB_HUD.md)
- [docs/WEB_WIKI.md](docs/WEB_WIKI.md)
- [docs/SOVEREIGN_REST_API.md](docs/SOVEREIGN_REST_API.md)
