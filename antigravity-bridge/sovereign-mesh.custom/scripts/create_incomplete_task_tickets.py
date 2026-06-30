#!/usr/bin/env python3
"""Script to generate tickets for all incomplete tasks (unchecked checklist items).

It scans the entire sovereign_mesh repository for lines containing an unchecked
checkbox "- [ ]" and creates a ticket in the RT/PQR SQLite database for each
found task.
"""

import os
import re
import sqlite3
import sys
from datetime import datetime, timezone

# Repository root (assumed to be the directory containing this script)
REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))

# Database path
DB_PATH = os.path.join(REPO_ROOT, "agent_pedigree.db")

# Pattern for an unchecked task line
TASK_PATTERN = re.compile(r"^- \[ \]\s*(.+)")

def find_incomplete_tasks(root: str):
    """Yield (file_path, line_number, task_text) for each unchecked task.

    The search is recursive, ignoring binary files and the .git directory.
    """
    for dirpath, dirnames, filenames in os.walk(root):
        # Skip specific large or sensitive directories, and meta-agent frameworks
        dirnames[:] = [d for d in dirnames if d not in ('.git', '.agent', '.github', 'node_modules', 'venv', '.venv')]
        for fname in filenames:
            # Broad filter for code, config, and documentation files
            extensions = (
                '.py', '.md', '.txt', '.sh', '.js', '.go', '.cpp', '.c', '.java', 
                '.json', '.yaml', '.yml', '.ts', '.tsx', '.rs', '.toml', '.ps1', 
                '.proto', '.sql', '.h', '.hpp', '.css', '.html', '.xml', '.mjs',
                '.sum', '.mod', '.lock'
            )
            if fname.endswith(extensions):
                fpath = os.path.join(dirpath, fname)
                try:
                    with open(fpath, "r", encoding="utf-8") as f:
                        for i, line in enumerate(f, start=1):
                            m = TASK_PATTERN.match(line.strip())
                            if m:
                                yield fpath, i, m.group(1).strip()
                except (UnicodeDecodeError, OSError):
                    # Skip files we cannot read as text
                    pass

def insert_ticket(conn, subject):
    now = datetime.now(timezone.utc).isoformat(timespec='seconds')
    cur = conn.cursor()
    cur.execute(
        """
        INSERT INTO tickets (
            Queue, Subject, Status, Owner, Creator, Priority, 
            TimeEstimated, TimeWorked, TimeLeft, Created
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
        ("Incomplete-Tasks", subject, "new", "Nobody", "system", 3, 0, 0, 0, now)
    )
    return cur.lastrowid

def main():
    # Ensure the database exists
    if not os.path.exists(DB_PATH):
        print(f"Database not found at {DB_PATH}. Exiting.")
        sys.exit(1)

    conn = sqlite3.connect(DB_PATH)
    created = 0
    
    print(f"Scanning {REPO_ROOT} for incomplete tasks...")
    for fpath, line_num, task_text in find_incomplete_tasks(REPO_ROOT):
        # Format subject for traceability
        rel_path = os.path.relpath(fpath, REPO_ROOT)
        subject = f"[{rel_path}:{line_num}] {task_text}"
        
        # Check if ticket already exists
        cur = conn.cursor()
        cur.execute("SELECT ticket_id FROM tickets WHERE Subject = ?", (subject,))
        if cur.fetchone():
            continue
            
        try:
            ticket_id = insert_ticket(conn, subject)
            created += 1
            print(f"Created ticket #{ticket_id}: {subject}")
        except sqlite3.IntegrityError as e:
            print(f"Failed to insert ticket for {subject}: {e}")
            
    conn.commit()
    conn.close()
    print(f"\nDone. {created} tickets created for incomplete tasks.")

if __name__ == "__main__":
    main()
