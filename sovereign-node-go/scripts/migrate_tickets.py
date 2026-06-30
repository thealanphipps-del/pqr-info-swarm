import psycopg2

def migrate():
    conn_str = "postgresql://root@127.0.0.1:5196/antigravity?sslmode=disable"
    print("Connecting to CockroachDB...")
    conn = psycopg2.connect(conn_str)
    cur = conn.cursor()

    print("Altering tickets table...")
    alter_queries = [
        "ALTER TABLE tickets ADD COLUMN IF NOT EXISTS priority INT DEFAULT 20;",
        "ALTER TABLE tickets ADD COLUMN IF NOT EXISTS queue STRING DEFAULT 'General';",
        "ALTER TABLE tickets ADD COLUMN IF NOT EXISTS assigned_to STRING DEFAULT '';",
        "ALTER TABLE tickets ADD COLUMN IF NOT EXISTS is_sticky BOOL DEFAULT false;",
        "ALTER TABLE tickets ADD COLUMN IF NOT EXISTS referrer_code STRING DEFAULT '';"
    ]

    for q in alter_queries:
        try:
            cur.execute(q)
            print(f"Executed: {q}")
        except Exception as e:
            print(f"Error on query '{q}': {e}")

    conn.commit()
    cur.close()
    conn.close()
    print("Migration completed successfully!")

if __name__ == "__main__":
    migrate()
