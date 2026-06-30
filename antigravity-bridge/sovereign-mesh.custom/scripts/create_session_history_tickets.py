#!/usr/bin/env python3
"""Create tickets for each entry in MUDD session history.

The MUDD session file (mudd_session.json) stores a list of visited rooms in the
"history" key. This script reads that file, extracts each room identifier and
creates a ticket in the RT/PQR SQLite database. The ticket subject includes the
room ID and a short description.
"""
import os
import json
import sqlite3
import sys
from datetime import datetime, timezone

# Paths – match the ones used by the rest of the project
REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
SESSION_FILE = os.path.join(REPO_ROOT, "mudd_session.json")
DB_PATH = os.path.join(REPO_ROOT, "agent_pedigree.db")

def load_history(session_path):
    with open(session_path, "r", encoding="utf-8") as f:
        data = json.load(f)
    return data.get("history", [])

def insert_ticket(conn, subject):
    now = datetime.now(timezone.utc).isoformat(timespec="seconds")
    cur = conn.cursor()
    cur.execute(
        """
        INSERT INTO tickets (
            Queue, Subject, Status, Owner, Creator, Priority, 
            TimeEstimated, TimeWorked, TimeLeft, Created
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
        ("Session-History", subject, "new", "Nobody", "system", 3, 0, 0, 0, now),
    )
    return cur.lastrowid

def main():
    if not os.path.isfile(SESSION_FILE):
        print(f"Session file not found at {SESSION_FILE}")
        return

    if not os.path.isfile(DB_PATH):
        print(f"Database not found at {DB_PATH}")
        return

    history = load_history(SESSION_FILE)
    if not history:
        print("No history found in session file.")
        return

    conn = sqlite3.connect(DB_PATH)
    created = 0
    
    # Process each room in history
    for i, room_id in enumerate(history, 1):
        subject = f"[Session:{i}] Visited Node {room_id}"
        
        # Check if already ticketed
        cur = conn.cursor()
        cur.execute("SELECT ticket_id FROM tickets WHERE Subject = ?", (subject,))
        if cur.fetchone():
            continue
            
        try:
            ticket_id = insert_ticket(conn, subject)
            created += 1
            print(f"Created ticket #{ticket_id}: {subject}")
        except sqlite3.IntegrityError as e:
            print(f"Failed to insert ticket for {room_id}: {e}")
    conn.commit()
    conn.close()
    print(f"\nDone. {created} tickets created from session history.")

if __name__ == "__main__":
    main()
