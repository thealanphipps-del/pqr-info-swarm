#!/data/data/com.termux/files/usr/bin/bash
# HUD INITIALIZATION FOR ALAN PHIPPS (GARLAND)
# MONITORING VITALITY_SLOPE VS FLOOR (98.2%)

termux-gui-view --create-activity << 'JSON'
{
  "layout": {
    "type": "LinearLayout",
    "orientation": "vertical",
    "children": [
      {"type": "TextView", "text": "SENTINEL HUD - ACTIVE", "id": "status_header"},
      {"type": "TextView", "text": "VITALITY: PENDING", "id": "vitality_val"},
      {"type": "Button", "text": "FATALITY PURGE (EMERGENCY)", "id": "purge_btn"},
      {"type": "Button", "text": "SYNC 39.MH", "id": "sync_btn"}
    ]
  }
}
JSON
