# PQR HUD (web/hud.html)

This document explains the **HUD overlay / companion page** shipped in `/web/`.

## Files
- `web/hud.html` – HUD page entry
- `web/hud.css` – HUD styling
- `web/index.html` / `web/index.css` – Shared landing styles (if used)
- `web/dashboard.html` – The main dashboard page (often linked)

## How to run it locally
### Option A: Open directly (quickest)
1. Open `web/hud.html` in a browser.
2. If you need API calls and CORS blocks them, use Option B.

### Option B: Serve the `web/` folder (recommended)
Open:
- `http://localhost:<port>/web/hud.html`

## Page behavior
The HUD is intended for a compact, “always visible” view of:
- system status / health signals
- short operational indicators
- quick navigation links (depending on wiring in `web/hud.html` / `web/hud.css`)

The HUD’s actual data sources (which REST/HTTP endpoints it calls) are defined in the HTML/JS it includes. If the HUD doesn’t display anything, the most likely causes are:
- backend not reachable from the browser origin
- base API URL misconfigured
- CORS restrictions

## Customization points
1. **API base URL / endpoint wiring**
   - Search for fetch/XHR calls inside `web/hud.html`.
2. **Layout**
   - Adjust sizes and positioning in `web/hud.css`.
3. **Add new signals**
   - Add a DOM element in `hud.html`
   - Update the JS logic to populate it

## Related docs
- [docs/WEB_DASHBOARD.md](docs/WEB_DASHBOARD.md)
- [docs/WEB_WIKI.md](docs/WEB_WIKI.md)
