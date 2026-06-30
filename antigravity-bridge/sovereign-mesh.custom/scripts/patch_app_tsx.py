#!/usr/bin/env python3
import sys

filepath = "/home/aellok/sovereign-mesh/frontend/src/App.tsx"

with open(filepath, "r", encoding="utf-8") as f:
    content = f.read()

# 1. Add states
target_states = """  const [topology, setTopology] = useState<MeshNode[]>([]);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);"""

replacement_states = """  const [topology, setTopology] = useState<MeshNode[]>([]);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [memorySegments, setMemorySegments] = useState<any[]>([]);
  const [selectedSegment, setSelectedSegment] = useState<any | null>(null);
  const [selectedTicket, setSelectedTicket] = useState<any | null>(null);"""

if target_states in content:
    content = content.replace(target_states, replacement_states)
    print("Added states.")
else:
    print("Failed to find states target!")

# 2. Add fetchData calls
target_fetch = """      // 5. Mesh Topology
      const topoRes = await fetch('/api/v2/mesh/topology');
      if (topoRes.ok) {
        const json = await topoRes.json();
        setTopology(json.nodes || []);
      }"""

replacement_fetch = """      // 5. Mesh Topology
      const topoRes = await fetch('/api/v2/mesh/topology');
      if (topoRes.ok) {
        const json = await topoRes.json();
        setTopology(json.nodes || []);
      }

      // 6. Memory segments
      const memRes = await fetch('/api/v2/memory/segments');
      if (memRes.ok) {
        const json = await memRes.json();
        setMemorySegments(json.data || []);
      }"""

if target_fetch in content:
    content = content.replace(target_fetch, replacement_fetch)
    print("Added fetch.")
else:
    print("Failed to find fetch target!")

# 3. Add ticket row click handler
target_click = """  // Spawn New Agent Form handler
  const handleSpawnAgent = async (e: React.FormEvent) => {"""

replacement_click = """  const handleTicketClick = async (ticketId: string | number) => {
    try {
      const res = await fetch(`/api/v2/tickets/${ticketId}`);
      if (res.ok) {
        const json = await res.json();
        setSelectedTicket(json);
      }
    } catch (err) {
      console.error("Failed to load ticket details:", err);
    }
  };

  // Spawn New Agent Form handler
  const handleSpawnAgent = async (e: React.FormEvent) => {"""

if target_click in content:
    content = content.replace(target_click, replacement_click)
    print("Added ticket click handler.")
else:
    print("Failed to find spawn agent target!")

# 4. Make ticket rows clickable
target_tr = """                          {pqrTickets.slice(0, 6).map((ticket) => (
                            <tr key={ticket.ticket_id}>"""

replacement_tr = """                          {pqrTickets.slice(0, 6).map((ticket) => (
                            <tr key={ticket.ticket_id} className="hover-row cursor-pointer" onClick={() => handleTicketClick(ticket.ticket_id)}>"""

if target_tr in content:
    content = content.replace(target_tr, replacement_tr)
    print("Made rows clickable.")
else:
    print("Failed to find ticket row target!")

# 5. Insert heatmap panel in Dashboard Console (below Row 1 ticket stream panel)
target_loki = """              {/* Row 2: Loki Log Terminal & Command Input (Span 12 of 12) */}"""

replacement_heatmap = """              {/* Row 1.5: Shared Memory Heatmap (Span 12 of 12) */}
              <div className="panel col-span-12 accent-orange">
                <div className="panel-header">
                  <div className="panel-title">
                    <Layers3 size={14} className="text-orange" />
                    <span>MMAP SHARED MEMORY BUS (64 SEGMENTS ACTIVE STATE TRACKING)</span>
                  </div>
                </div>
                <div className="panel-content">
                  <p className="text-muted mb-12">
                    Real-time state tracking of the high-speed paged RAM bus. Click on a segment cell to resolve its 5D Coordinate Address and active pipeline metrics.
                  </p>
                  <div className="memory-bus-heatmap">
                    {memorySegments.map((seg) => (
                      <button
                        key={seg.offset}
                        className={`heatmap-cell ${seg.status.toLowerCase()} ${selectedSegment?.offset === seg.offset ? 'selected' : ''}`}
                        onClick={() => setSelectedSegment(seg)}
                        title={`Segment ${seg.offset}: ${seg.status}`}
                      >
                        {seg.offset}
                      </button>
                    ))}
                  </div>
                  
                  {selectedSegment && (
                    <div className="segment-detail-overlay mt-12 panel accent-blue">
                      <div className="panel-header">
                        <div className="panel-title">
                          <span>SEGMENT #{selectedSegment.offset} ACTIVE STATE DETAILED REPORT</span>
                        </div>
                        <button className="clear-filter-btn" onClick={() => setSelectedSegment(null)}>Close</button>
                      </div>
                      <div className="panel-content text-lines grid grid-cols-2 gap-16">
                        <div>
                          <h4>🌌 Mapped 5D Address Coordinates</h4>
                          <ul className="custom-bullets font-mono">
                            <li><strong>Dimension 1 (NodeID):</strong> {selectedSegment.address_5d.node_id}</li>
                            <li><strong>Dimension 2 (RoleID):</strong> {selectedSegment.address_5d.role_id}</li>
                            <li><strong>Dimension 3 (LineageID):</strong> {selectedSegment.address_5d.lineage_id}</li>
                            <li><strong>Dimension 4 (BlockID):</strong> {selectedSegment.address_5d.block_id}</li>
                            <li><strong>Dimension 5 (ThreadID):</strong> {selectedSegment.address_5d.thread_id}</li>
                          </ul>
                        </div>
                        <div>
                          <h4>🛰️ Spatial Hypergrid Coordinates</h4>
                          <ul className="custom-bullets font-mono">
                            <li><strong>X Coordinate:</strong> {selectedSegment.coordinates.x} / 1024</li>
                            <li><strong>Y Coordinate:</strong> {selectedSegment.coordinates.y} / 1024</li>
                            <li><strong>Z Coordinate:</strong> {selectedSegment.coordinates.z} / 1024</li>
                            <li><strong>Ψ Coordinate (hyperplane):</strong> {selectedSegment.coordinates.psi} / 27</li>
                            <li><strong>T (System nanoseconds):</strong> {selectedSegment.coordinates.t}</li>
                          </ul>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </div>

              {/* Row 2: Loki Log Terminal & Command Input (Span 12 of 12) */}"""

if target_loki in content:
    content = content.replace(target_loki, replacement_heatmap)
    print("Inserted heatmap panel.")
else:
    print("Failed to find loki row target!")

# 6. Render Modal for ticket details (placed right before final closing tags)
target_end = """        </main>
      </div>
    </div>
  );
}"""

replacement_modal = """          {selectedTicket && (
            <div className="modal-overlay">
              <div className="modal-container panel accent-orange">
                <div className="panel-header">
                  <div className="panel-title">
                    <FileText size={14} className="text-orange" />
                    <span>PQR TICKET #{selectedTicket.data.ticket_id} - RT REST 2.0 GO CLIENT RESOLUTION</span>
                  </div>
                  <button className="clear-filter-btn" onClick={() => setSelectedTicket(null)}>CLOSE CLIENT</button>
                </div>
                <div className="panel-content modal-content">
                  <div className="grid grid-cols-2 gap-16 text-lines">
                    <div>
                      <h3>📋 Ticket Metadata</h3>
                      <p><strong>Subject:</strong> {selectedTicket.data.Subject}</p>
                      <p><strong>Queue:</strong> {selectedTicket.data.Queue}</p>
                      <p><strong>Priority:</strong> {selectedTicket.data.Priority}</p>
                      <p><strong>Status:</strong> <span className="badge info">{selectedTicket.data.Status}</span></p>
                      <p><strong>Creator:</strong> {selectedTicket.data.Creator}</p>
                      <p><strong>Created:</strong> {selectedTicket.data.Created}</p>
                      <p><strong>Last Updated:</strong> {selectedTicket.data.LastUpdated}</p>
                    </div>
                    <div>
                      <h3>🤖 Assigned Specialist</h3>
                      <p><strong>Agent Assigned:</strong> <code className="text-blue">{selectedTicket.data.agent_id || "Unassigned"}</code></p>
                      <p><strong>Specialty Domain:</strong> <code className="text-purple">{selectedTicket.data.specialty || "General Swarm"}</code></p>
                      <p><strong>Task Description:</strong></p>
                      <pre className="pre-wrap-box">{selectedTicket.data.task_description || "No description provided."}</pre>
                    </div>
                  </div>
                  
                  <div className="go-client-shell mt-16">
                    <h3>📟 Go RT REST 2.0 Client Output (stdout/stderr)</h3>
                    <pre className="terminal-stdout-box">
                      {`$ ${selectedTicket.go_client.command}\\n`}
                      {selectedTicket.go_client.raw_output}
                    </pre>
                  </div>
                </div>
              </div>
            </div>
          )}

        </main>
      </div>
    </div>
  );
}"""

if target_end in content:
    content = content.replace(target_end, replacement_modal)
    print("Inserted modal.")
else:
    print("Failed to find end tags target!")

# 7. Update Topology tab nodes to display 5D addressing if available
target_topo = """                        <div className="node-details">
                          <div className="node-name">{node.id}</div>
                          <div className="node-ip text-muted">{node.ip}</div>
                        </div>"""

replacement_topo = """                        <div className="node-details">
                          <div className="node-name">{node.id}</div>
                          <div className="node-ip text-muted">{node.ip}</div>
                          {node.address_5d && (
                            <div className="node-address5d font-mono text-muted mt-8" style={{ fontSize: '11px', borderTop: '1px solid var(--border)', paddingTop: '4px' }}>
                              <div>Coord: {node.address_5d.node_id} | {node.address_5d.role_id}</div>
                              <div>Lineage: {node.address_5d.lineage_id}</div>
                              <div>Block/Thread: {node.address_5d.block_id} / {node.address_5d.thread_id}</div>
                            </div>
                          )}
                        </div>"""

if target_topo in content:
    content = content.replace(target_topo, replacement_topo)
    print("Added 5D addressing on node cards.")
else:
    print("Failed to find node topology target!")

with open(filepath, "w", encoding="utf-8") as f:
    f.write(content)
print("Saved App.tsx.")
