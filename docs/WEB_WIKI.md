# PQR Web Wiki (web/wiki.html)

This document explains the **web wiki/help page** shipped in `/web/`.

## Files
- `web/wiki.html` – Wiki/help page entry
- `web/index.css` – Shared landing/wiki styling (if used)
- `web/hud.css` / `web/dashboard.css` – Other shared styles (if referenced)

## How to run it locally
### Option A: Open directly (quickest)
1. Open `web/wiki.html` in a browser.

### Option B: Serve the `web/` folder (recommended)
Open:
- `http://localhost:<port>/web/wiki.html`

## Page behavior
The wiki page is generally used for:
- human-readable system notes
- references between docs and UI pages
- navigation to other assets under `/web/`

If the wiki includes links to dynamic API-backed information, those will be implemented in the HTML itself (or via referenced scripts). If links don’t work:
- ensure the relative paths in `wiki.html` are correct
- ensure you’re serving the full `web/` folder (not just the single file)

## Customization points
- Update content/sections in `web/wiki.html`
- Update styling in `web/index.css`
- Ensure link targets remain valid when hosted under a sub-path like `/web/`

## Related docs
- [docs/WEB_DASHBOARD.md](docs/WEB_DASHBOARD.md)
- [docs/WEB_HUD.md](docs/WEB_HUD.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
