#!/usr/bin/env python3
import curses
import os
import sys
import json
import subprocess

# --- AESTHETIC CONSTANTS ---
C_CYAN = 1
C_GOLD = 2
C_GREEN = 3
C_RED = 4
C_MAGENTA = 5
C_BLUE = 6

# --- MESH TOPOLOGY ---
NODES = [
    {"id": "0.MH", "name": "The Nuremberg Crypt", "role": "ANCHOR", "ip": "46.224.84.64"},
    {"id": "38.MH", "name": "The Forge Chamber", "role": "FORGE", "ip": "62.238.2.240"},
    {"id": "39.MH", "name": "The Helsinki Sentry", "role": "SENTRY", "ip": "204.168.138.60"},
    {"id": "40.MH", "name": "The Capicant Lab", "role": "CAPICANT", "ip": "10.128.0.2"},
    {"id": "50.MH", "name": "The Operations Hub", "role": "OPS", "ip": "34.42.122.68"},
    {"id": "201.MH", "name": "The Edge Relay Vault", "role": "EDGE", "ip": "89.167.91.81"},
    {"id": "7.MH", "name": "The Tokyo Phased Node", "role": "PHASED", "ip": "35.243.111.42"},
    {"id": "8.MH", "name": "The Mumbai Phased Node", "role": "PHASED", "ip": "34.93.120.15"},
    {"id": "9.MH", "name": "The Garland Phased Node", "role": "PHASED", "ip": "34.123.50.201"},
    {"id": "AURORA", "name": "The Local Training Node", "role": "LOCAL", "ip": "127.0.0.1"}
]

def draw_menu(stdscr, selected_idx):
    h, w = stdscr.getmaxyx()
    stdscr.clear()
    
    # Title
    title = " 🌐 SOVEREIGN MESH: INTERACTIVE NAVIGATOR "
    stdscr.attron(curses.color_pair(C_CYAN) | curses.A_BOLD | curses.A_REVERSE)
    stdscr.addstr(0, max(0, w//2 - len(title)//2), title)
    stdscr.attroff(curses.color_pair(C_CYAN) | curses.A_BOLD | curses.A_REVERSE)
    
    # Status Line
    status = " Use UP/DOWN to navigate nodes | ENTER to select | 'q' to exit "
    stdscr.addstr(h-2, max(0, w//2 - len(status)//2), status, curses.A_DIM)

    # Draw Node List
    menu_start_y = 3
    for idx, node in enumerate(NODES):
        x = max(2, w//4)
        y = menu_start_y + idx
        
        indicator = " ▶ " if idx == selected_idx else "   "
        style = curses.color_pair(C_GOLD) | curses.A_BOLD if idx == selected_idx else curses.A_NORMAL
        
        line = f"{indicator}{node['id']:<8} | {node['role']:<10} | {node['name']}"
        
        if idx == selected_idx:
            stdscr.attron(curses.A_REVERSE)
            stdscr.addstr(y, x, line)
            stdscr.attroff(curses.A_REVERSE)
        else:
            stdscr.addstr(y, x, line, style)

    # Draw Details Box for selected node
    sel = NODES[selected_idx]
    details_y = menu_start_y
    details_x = max(2, w//2 + 5)
    
    stdscr.attron(curses.color_pair(C_CYAN) | curses.A_BOLD)
    stdscr.addstr(details_y, details_x, f"📡 NODE: {sel['id']}")
    stdscr.attroff(curses.color_pair(C_CYAN) | curses.A_BOLD)
    
    stdscr.addstr(details_y+2, details_x, f"🏷️  NAME: {sel['name']}")
    stdscr.addstr(details_y+3, details_x, f"🛠️  ROLE: {sel['role']}")
    stdscr.addstr(details_y+4, details_x, f"🌐 IP:   {sel['ip']}")
    
    stdscr.attron(curses.color_pair(C_GREEN))
    stdscr.addstr(details_y+6, details_x, "[ ONLINE ]")
    stdscr.attroff(curses.color_pair(C_GREEN))

    stdscr.refresh()

def main(stdscr):
    # Initialize Colors
    curses.start_color()
    curses.init_pair(C_CYAN, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(C_GOLD, curses.COLOR_YELLOW, curses.COLOR_BLACK)
    curses.init_pair(C_GREEN, curses.COLOR_GREEN, curses.COLOR_BLACK)
    curses.init_pair(C_RED, curses.COLOR_RED, curses.COLOR_BLACK)
    
    curses.curs_set(0) # Hide cursor
    selected_idx = 0
    
    while True:
        draw_menu(stdscr, selected_idx)
        
        key = stdscr.getch()
        
        if key == curses.KEY_UP:
            selected_idx = (selected_idx - 1) % len(NODES)
        elif key == curses.KEY_DOWN:
            selected_idx = (selected_idx + 1) % len(NODES)
        elif key in [curses.KEY_ENTER, ord('\n')]:
            # Selection Action
            sel = NODES[selected_idx]
            curses.endwin()
            print(f"\n🚀 Teleporting mind segment to {sel['id']} ({sel['name']})...")
            # In a real scenario, we'd trigger the teleport protocol here
            os.system(f"bash mesh_control.sh exec 'echo Signal synchronized at {sel['id']}'")
            input("\nPress ENTER to return to Navigator...")
        elif key == ord('q'):
            break

if __name__ == "__main__":
    try:
        curses.wrapper(main)
    except curses.error as e:
        print(f"Error: Terminal too small or unsupported. ({e})")
        sys.exit(1)
