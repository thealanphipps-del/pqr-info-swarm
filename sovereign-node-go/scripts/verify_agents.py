import sqlite3
c = sqlite3.connect('/home/aellok/sovereign-mesh/agent_pedigree.db')
rows = c.execute('SELECT agent_id, name, layer_level FROM agents ORDER BY layer_level').fetchall()
print(f"\nTotal agents: {len(rows)}")
print("\nLayer 6 Shared Resource Bindings:")
for r in rows:
    if r[2] == 6:
        print(f"  ✅ L{r[2]}: {r[0]} — {r[1]}")
c.close()
