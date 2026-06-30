import React, { useState, useEffect, useRef } from 'react';
import { 
  Terminal as TerminalIcon, 
  Database, 
  Shield, 
  Activity, 
  Users, 
  Send, 
  AlertTriangle, 
  Clock, 
  RefreshCw, 
  Cpu, 
  Search, 
  Layers3,
  Network,
  Lock,
  Wifi,
  FileText
} from 'lucide-react';
import './App.css';

// Type Definitions
interface Agent {
  agent_id: string;
  name: string;
  parent_agent_id: string | null;
  layer_level: number;
  specialty: string;
  subspecialty: string;
  _links?: any;
}

interface LedgerBlock {
  block_index: number;
  previous_hash: string;
  timestamp: string;
  agent_id: string;
  mutation_payload: string;
  consensus_votes: string;
  minority_report: string;
  block_hash: string;
}

interface PqrTicket {
  ticket_id: string | number;
  Queue: string;
  Subject: string;
  Status: string;
  Owner: string;
  Creator: string;
  Priority: string;
  TimeEstimated: string | number;
  TimeWorked: string | number;
  TimeLeft: string | number;
  Created: string;
  Resolved: string | null;
  LastUpdated: string;
  LastUpdatedBy: string;
  agent_id: string | null;
  layer_level: number | null;
  specialty: string | null;
  task_description: string | null;
}

interface MeshNode {
  id: string;
  ip: string;
  role: string;
  status: string;
  address_5d?: {
    node_id: string;
    role_id: string;
    lineage_id: string;
    block_id: string;
    thread_id: string;
  };
}

function App() {
  // Tabs State
  const [activeTab, setActiveTab] = useState<'dashboard' | 'agents' | 'ledger' | 'topology'>('dashboard');

  // Core Data State
  const [systemStatus, setSystemStatus] = useState<string>('Initializing Forensic Hub Sync...\nSystem online.\nRunlevel: 7 (STARBIRTH)');
  const [agents, setAgents] = useState<Agent[]>([]);
  const [ledger, setLedger] = useState<LedgerBlock[]>([]);
  const [pqrTickets, setPqrTickets] = useState<PqrTicket[]>([]);
  const [topology, setTopology] = useState<MeshNode[]>([]);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [memorySegments, setMemorySegments] = useState<any[]>([]);
  const [selectedSegment, setSelectedSegment] = useState<any | null>(null);
  const [selectedTicket, setSelectedTicket] = useState<any | null>(null);

  // Form / Action States
  const [command, setCommand] = useState<string>('');
  const [terminalLogs, setTerminalLogs] = useState<string[]>([
    '[SYSTEM] Sovereign v10.0 Singularity Kernel loaded successfully.',
    '[BRIDGE] Autonomous Gemini Hooks active on native workspace.',
    '[BRIDGE] gRPC Control Bus dialer initialized on 127.0.0.1:1111.',
    '[BRIDGE] HighSpeed RAM Bus dialer active on 127.0.0.1:11111.',
    'Ready for secure input...'
  ]);
  const [logFilter, setLogFilter] = useState<string>('');
  
  // New Agent Form
  const [newAgentId, setNewAgentId] = useState('');
  const [newAgentName, setNewAgentName] = useState('');
  const [newAgentParent, setNewAgentParent] = useState('');
  const [newAgentLayer, setNewAgentLayer] = useState(7);
  const [newAgentSpecialty, setNewAgentSpecialty] = useState('');
  const [newAgentSubspecialty, setNewAgentSubspecialty] = useState('');
  const [agentFormMsg, setAgentFormMsg] = useState({ text: '', type: '' });

  // Time Travel Form
  const [targetBlockIndex, setTargetBlockIndex] = useState<number>(0);
  const [timeMutationKey, setTimeMutationKey] = useState<string>('');
  const [timeMutationVal, setTimeMutationVal] = useState<string>('');
  const [timeMutationReason, setTimeMutationReason] = useState<string>('');
  const [timeMutationMsg, setTimeMutationMsg] = useState({ text: '', type: '' });

  // Terminal Auto Scroll
  const terminalEndRef = useRef<HTMLDivElement>(null);

  // Fetch all data
  const fetchData = async () => {
    setIsRefreshing(true);
    try {
      // 1. System Status
      const statusRes = await fetch('/api/status');
      if (statusRes.ok) {
        const text = await statusRes.text();
        setSystemStatus(text);
      }

      // 2. Agents
      const agentsRes = await fetch('/api/v2/agents');
      if (agentsRes.ok) {
        const json = await agentsRes.json();
        setAgents(json.data || []);
      }

      // 3. Ledger
      const ledgerRes = await fetch('/api/v2/ledger');
      if (ledgerRes.ok) {
        const json = await ledgerRes.json();
        setLedger(json.data || []);
      }

      // 4. Tickets (PQR Federated Ticketing)
      const ticketsRes = await fetch('/api/v2/tickets');
      if (ticketsRes.ok) {
        const json = await ticketsRes.json();
        setPqrTickets(json.data || []);
      }

      // 5. Mesh Topology
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
      }
    } catch (err) {
      console.error('Error fetching dashboard data:', err);
    } finally {
      setIsRefreshing(false);
    }
  };

  // Initial and periodic fetch
  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 8000);
    return () => clearInterval(interval);
  }, []);

  // Auto-scroll terminal on new logs
  useEffect(() => {
    if (terminalEndRef.current) {
      terminalEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [terminalLogs]);

  // Execute Bridge Command
  const handleExecuteCommand = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    const cmdTrimmed = command.trim();
    if (!cmdTrimmed) return;

    // Echo command in terminal
    setTerminalLogs(prev => [
      ...prev,
      `$ ${cmdTrimmed}`,
      `[EXEC] Routing execution directive to shell bridge...`
    ]);
    setCommand('');

    try {
      const response = await fetch(`/api/bridge?cmd=${encodeURIComponent(cmdTrimmed)}`);
      const text = await response.text();
      // Format text output lines
      const outputLines = text.split('\n').filter(line => line.length > 0);
      setTerminalLogs(prev => [...prev, ...outputLines]);
    } catch (err) {
      setTerminalLogs(prev => [...prev, `[ERROR] Failed to contact shell bridge: ${err}`]);
    }
  };

  const handleTicketClick = async (ticketId: string | number) => {
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
  const handleSpawnAgent = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAgentId || !newAgentName) {
      setAgentFormMsg({ text: 'Agent ID and Name are required.', type: 'error' });
      return;
    }

    try {
      const res = await fetch('/api/v2/agents', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          agent_id: newAgentId,
          name: newAgentName,
          parent_agent_id: newAgentParent || 'None',
          layer_level: newAgentLayer,
          specialty: newAgentSpecialty,
          subspecialty: newAgentSubspecialty
        })
      });
      const json = await res.json();
      if (json.success) {
        setAgentFormMsg({ text: json.message, type: 'success' });
        // Clear Form
        setNewAgentId('');
        setNewAgentName('');
        setNewAgentParent('');
        setNewAgentSpecialty('');
        setNewAgentSubspecialty('');
        // Refresh Agents list
        fetchData();
      } else {
        setAgentFormMsg({ text: json.error?.message || 'Spawn failed.', type: 'error' });
      }
    } catch (err) {
      setAgentFormMsg({ text: `Network failure: ${err}`, type: 'error' });
    }
  };

  // Retroactive Time Travel Override handler
  const handleTimeMutation = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!timeMutationKey || !timeMutationVal) {
      setTimeMutationMsg({ text: 'Mutation target Key and proposed Value are required.', type: 'error' });
      return;
    }

    try {
      const res = await fetch('/api/v2/timemachine', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          block_index: targetBlockIndex,
          key: timeMutationKey,
          value: timeMutationVal,
          reason: timeMutationReason
        })
      });
      const json = await res.json();
      if (json.data?.success) {
        setTimeMutationMsg({ 
          text: `Success: ${json.data.message}. Refactored logs: ${json.data.refactor_logs?.join(', ') || 'none'}. New Chain validation: ${json.data.new_chain_validation_status}`, 
          type: 'success' 
        });
        setTimeMutationKey('');
        setTimeMutationVal('');
        setTimeMutationReason('');
        // Refresh
        fetchData();
      } else {
        setTimeMutationMsg({ text: json.error?.message || 'Time Travel Override failed.', type: 'error' });
      }
    } catch (err) {
      setTimeMutationMsg({ text: `Network failure: ${err}`, type: 'error' });
    }
  };

  // Render logs filtered by Loki Query Bar
  const filteredLogs = terminalLogs.filter(logLine => {
    if (!logFilter) return true;
    return logLine.toLowerCase().includes(logFilter.toLowerCase());
  });

  return (
    <div className="app-container">
      {/* Top Header Bar */}
      <header className="header-bar">
        <div className="header-brand">
          <Network className="glow-icon accent-blue pulse" size={18} />
          <h1>SOVEREIGN <span>V10.0 SINGULARITY</span></h1>
          <span className="badge success">Swarm Active</span>
        </div>
        
        <div className="header-metrics">
          <div className="metric-ticker">
            <Wifi className="metric-icon success" size={13} />
            <span className="metric-label">Consensus Bus</span>
            <span className="metric-value text-green glow-text">ONLINE</span>
          </div>
          
          <div className="metric-ticker">
            <Lock className="metric-icon orange" size={13} />
            <span className="metric-label">ASSET LOCK</span>
            <span className="metric-value text-orange glow-text">$814.68</span>
          </div>

          <button 
            className={`refresh-btn ${isRefreshing ? 'rotating' : ''}`}
            onClick={fetchData}
            title="Force telemetry refresh"
          >
            <RefreshCw size={14} className={isRefreshing ? 'rotate' : ''} />
          </button>
        </div>
      </header>

      {/* Main Page Layout */}
      <div className="workspace">
        {/* Left Grafana Sidebar */}
        <aside className="sidebar">
          <nav className="sidebar-nav">
            <button 
              className={`nav-item ${activeTab === 'dashboard' ? 'active' : ''}`}
              onClick={() => setActiveTab('dashboard')}
              title="Telemetry Dashboard"
            >
              <Activity size={18} />
              <span>Console</span>
            </button>
            
            <button 
              className={`nav-item ${activeTab === 'agents' ? 'active' : ''}`}
              onClick={() => setActiveTab('agents')}
              title="7-Layer Swarm Map"
            >
              <Users size={18} />
              <span>Swarm</span>
            </button>
            
            <button 
              className={`nav-item ${activeTab === 'ledger' ? 'active' : ''}`}
              onClick={() => setActiveTab('ledger')}
              title="Relational Ledger Audit"
            >
              <Database size={18} />
              <span>Ledger</span>
            </button>
            
            <button 
              className={`nav-item ${activeTab === 'topology' ? 'active' : ''}`}
              onClick={() => setActiveTab('topology')}
              title="Mesh Topology & PQR Tunnels"
            >
              <Layers3 size={18} />
              <span>Network</span>
            </button>
          </nav>

          <div className="sidebar-footer">
            <Shield size={16} className="text-muted" />
            <div className="sidebar-ver">v10.0.0-FE</div>
          </div>
        </aside>

        {/* Dashboard Panels Viewport */}
        <main className="content-viewport">
          
          {/* TAB 1: DASHBOARD CONSOLE */}
          {activeTab === 'dashboard' && (
            <div className="dashboard-grid">
              
              {/* Row 1: System Telemetry & Status Panel (Span 4 of 12) */}
              <div className="panel accent-blue col-span-4 min-h-300">
                <div className="panel-header">
                  <div className="panel-title">
                    <Cpu size={14} className="text-blue" />
                    <span>SYSTEM HARDWARE TELEMETRY</span>
                  </div>
                  <span className="badge success">LIVE</span>
                </div>
                <div className="panel-content">
                  <div className="status-terminal">
                    <pre>{systemStatus}</pre>
                  </div>
                  <div className="telemetry-bar-list mt-12">
                    <div className="bar-item">
                      <div className="bar-label">Swarm Runlevel Status</div>
                      <div className="bar-container">
                        <div className="bar-fill bg-blue" style={{ width: '85%' }}></div>
                      </div>
                      <div className="bar-value">85% Vitality</div>
                    </div>
                    <div className="bar-item">
                      <div className="bar-label">Consensus Agree Ratio</div>
                      <div className="bar-container">
                        <div className="bar-fill bg-green" style={{ width: '92%' }}></div>
                      </div>
                      <div className="bar-value">92% (Consensus)</div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Row 1: PQR Ticket Stream Panel (Span 8 of 12) */}
              <div className="panel accent-orange col-span-8 min-h-300">
                <div className="panel-header">
                  <div className="panel-title">
                    <FileText size={14} className="text-orange" />
                    <span>PQR REGISTER-EVENT TICKETING (REST 2.0 BRIDGE)</span>
                  </div>
                  <span className="badge info">{pqrTickets.length} ACTIVE</span>
                </div>
                <div className="panel-content">
                  {pqrTickets.length === 0 ? (
                    <div className="empty-state">
                      <AlertTriangle className="text-muted" size={24} />
                      <p>No active PQR ticketing logs discovered.</p>
                    </div>
                  ) : (
                    <div className="ticket-table-wrapper">
                      <table className="custom-table">
                        <thead>
                          <tr>
                            <th>ID</th>
                            <th>Subject</th>
                            <th>Queue</th>
                            <th>Priority</th>
                            <th>Status</th>
                            <th>Assigned Specialty</th>
                          </tr>
                        </thead>
                        <tbody>
                          {pqrTickets.slice(0, 6).map((ticket) => (
                            <tr key={ticket.ticket_id} className="hover-row cursor-pointer" onClick={() => handleTicketClick(ticket.ticket_id)}>
                              <td className="text-orange">#{ticket.ticket_id}</td>
                              <td>{ticket.Subject}</td>
                              <td><span className="badge info">{ticket.Queue}</span></td>
                              <td>{ticket.Priority}</td>
                              <td>
                                <span className={`badge ${
                                  ticket.Status === 'resolved' || ticket.Status === 'Open' ? 'success' : 'warn'
                                }`}>
                                  {ticket.Status}
                                </span>
                              </td>
                              <td className="text-muted">{ticket.specialty || 'General Swarm'}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              </div>

              {/* Row 1.5: Shared Memory Heatmap (Span 12 of 12) */}
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

              {/* Row 2: Loki Log Terminal & Command Input (Span 12 of 12) */}
              <div className="panel col-span-12 accent-purple">
                <div className="panel-header">
                  <div className="panel-title">
                    <TerminalIcon size={14} className="text-purple" />
                    <span>LOKI LOG VITALITY TERMINAL (NATIVE BASH BRIDGE)</span>
                  </div>
                  
                  {/* Loki Query Bar */}
                  <div className="loki-query-bar">
                    <Search size={12} className="text-muted" />
                    <input 
                      type="text" 
                      placeholder='Loki Filter: e.g. {level="error"} or search query...' 
                      value={logFilter}
                      onChange={(e) => setLogFilter(e.target.value)}
                    />
                    {logFilter && (
                      <button className="clear-filter-btn" onClick={() => setLogFilter('')}>
                        Clear
                      </button>
                    )}
                  </div>
                </div>
                
                <div className="panel-content terminal-viewport">
                  <div className="terminal-lines">
                    {filteredLogs.map((log, index) => {
                      let colorClass = 'log-info';
                      if (log.startsWith('$')) colorClass = 'log-cmd glow-text';
                      else if (log.includes('[ERROR]') || log.includes('CORRUPTED') || log.includes('failed')) colorClass = 'log-err';
                      else if (log.includes('[SUCCESS]') || log.includes('SECURE') || log.includes('ACK')) colorClass = 'log-success';
                      else if (log.includes('[EXEC]')) colorClass = 'log-warn';

                      return (
                        <div key={index} className={`terminal-line ${colorClass}`}>
                          {log}
                        </div>
                      );
                    })}
                    <div ref={terminalEndRef} />
                  </div>
                </div>

                <div className="panel-footer-input">
                  <form onSubmit={handleExecuteCommand} className="terminal-input-form">
                    <span className="prompt-prefix">$</span>
                    <input 
                      type="text" 
                      placeholder="Execute terminal bash command or mesh directive..." 
                      value={command}
                      onChange={(e) => setCommand(e.target.value)}
                    />
                    <button type="submit" className="primary">
                      <Send size={12} />
                      <span>EXECUTE</span>
                    </button>
                  </form>
                </div>
              </div>

            </div>
          )}

          {/* TAB 2: 7-LAYER SWARM AGENTS */}
          {activeTab === 'agents' && (
            <div className="dashboard-grid">
              
              {/* Agent Relational Map (Span 8 of 12) */}
              <div className="panel col-span-8 accent-blue min-h-500">
                <div className="panel-header">
                  <div className="panel-title">
                    <Users size={14} className="text-blue" />
                    <span>7-LAYER RELATIONAL GRAPH AGENTS ({agents.length} SPAWNED)</span>
                  </div>
                </div>
                <div className="panel-content">
                  <div className="agents-tree-view">
                    <table className="custom-table">
                      <thead>
                        <tr>
                          <th>Layer</th>
                          <th>Agent ID</th>
                          <th>Emergent Persona</th>
                          <th>Parent Anchor</th>
                          <th>Specialty Domain</th>
                          <th>Subspecialty</th>
                        </tr>
                      </thead>
                      <tbody>
                        {agents.map((agent) => (
                          <tr key={agent.agent_id} className="hover-row">
                            <td>
                              <span className="badge success">L{agent.layer_level}</span>
                            </td>
                            <td className="text-blue font-mono">{agent.agent_id}</td>
                            <td className="text-bright font-medium">{agent.name}</td>
                            <td className="text-purple font-mono">{agent.parent_agent_id || 'Genesis Node'}</td>
                            <td>{agent.specialty}</td>
                            <td className="text-muted">{agent.subspecialty}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              {/* Agent Spawn Form (Span 4 of 12) */}
              <div className="panel col-span-4 accent-purple min-h-500">
                <div className="panel-header">
                  <div className="panel-title">
                    <Layers3 size={14} className="text-purple" />
                    <span>SPAWN &amp; REGISTER SWARM AGENT</span>
                  </div>
                </div>
                <div className="panel-content">
                  <form onSubmit={handleSpawnAgent} className="custom-form">
                    <div className="form-group">
                      <label>Agent Identifier ID (Unique shortcode)</label>
                      <input 
                        type="text" 
                        placeholder="e.g. 42.mh" 
                        value={newAgentId} 
                        onChange={(e) => setNewAgentId(e.target.value)} 
                        required 
                      />
                    </div>
                    
                    <div className="form-group">
                      <label>Emergent Persona Name</label>
                      <input 
                        type="text" 
                        placeholder="e.g. Sentry Watcher" 
                        value={newAgentName} 
                        onChange={(e) => setNewAgentName(e.target.value)} 
                        required 
                      />
                    </div>

                    <div className="form-group">
                      <label>Parent Anchor Node ID</label>
                      <input 
                        type="text" 
                        placeholder="e.g. 39.mh (optional)" 
                        value={newAgentParent} 
                        onChange={(e) => setNewAgentParent(e.target.value)} 
                      />
                    </div>

                    <div className="form-group">
                      <label>Swarm Layer Assignment (1 to 7)</label>
                      <select 
                        value={newAgentLayer} 
                        onChange={(e) => setNewAgentLayer(Number(e.target.value))}
                      >
                        <option value={1}>Layer 1 - Physical Interface</option>
                        <option value={2}>Layer 2 - Network Tunnels</option>
                        <option value={3}>Layer 3 - Consensus Engine</option>
                        <option value={4}>Layer 4 - Process Swaps</option>
                        <option value={5}>Layer 5 - Multi-Swarm Exec</option>
                        <option value={6}>Layer 6 - Neural Interface</option>
                        <option value={7}>Layer 7 - Cognitive Persona</option>
                      </select>
                    </div>

                    <div className="form-group">
                      <label>Core Specialty Domain</label>
                      <input 
                        type="text" 
                        placeholder="e.g. Security Audit" 
                        value={newAgentSpecialty} 
                        onChange={(e) => setNewAgentSpecialty(e.target.value)} 
                      />
                    </div>

                    <div className="form-group">
                      <label>Subspecialty Domain</label>
                      <input 
                        type="text" 
                        placeholder="e.g. File Ownership checks" 
                        value={newAgentSubspecialty} 
                        onChange={(e) => setNewAgentSubspecialty(e.target.value)} 
                      />
                    </div>

                    <button type="submit" className="primary mt-12 w-full">
                      <span>SPAWN INTER-AGENT SYSTEM</span>
                    </button>

                    {agentFormMsg.text && (
                      <div className={`form-message mt-12 ${agentFormMsg.type}`}>
                        {agentFormMsg.text}
                      </div>
                    )}
                  </form>
                </div>
              </div>

            </div>
          )}

          {/* TAB 3: MASTER RELATIONAL LEDGER & TIME MACHINE */}
          {activeTab === 'ledger' && (
            <div className="dashboard-grid">
              
              {/* Relational Relayed Relays Ledger Blocks (Span 8 of 12) */}
              <div className="panel col-span-8 accent-green min-h-500">
                <div className="panel-header">
                  <div className="panel-title">
                    <Database size={14} className="text-green" />
                    <span>IMMUTABLE RELATIONAL BLOCKCHAIN LEDGER</span>
                  </div>
                </div>
                <div className="panel-content">
                  <div className="ledger-timeline-wrapper">
                    {ledger.length === 0 ? (
                      <div className="empty-state">
                        <AlertTriangle size={24} className="text-muted" />
                        <p>No ledger blocks synchronised.</p>
                      </div>
                    ) : (
                      <div className="timeline-cards">
                        {ledger.map((block) => (
                          <div key={block.block_index} className="timeline-card" onClick={() => setTargetBlockIndex(block.block_index)}>
                            <div className="timeline-card-header">
                              <span className="block-index">Block #{block.block_index}</span>
                              <span className="block-time"><Clock size={11} /> {block.timestamp}</span>
                            </div>
                            <div className="timeline-card-body">
                              <div className="payload-line">
                                <span className="label">Proposer:</span>
                                <span className="val font-mono">{block.agent_id}</span>
                              </div>
                              <div className="payload-line">
                                <span className="label">Payload:</span>
                                <span className="val">{block.mutation_payload}</span>
                              </div>
                              <div className="payload-line">
                                <span className="label">Consensus:</span>
                                <span className="val text-green">{block.consensus_votes}</span>
                              </div>
                              {block.minority_report && (
                                <div className="payload-line minority-report-alert">
                                  <span className="label text-orange">Minority Report:</span>
                                  <span className="val text-orange">{block.minority_report}</span>
                                </div>
                              )}
                              <div className="payload-line hash-line font-mono text-muted">
                                <span className="label">Hash:</span>
                                <span className="val font-mono">{block.block_hash}</span>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Time Travel Override Console (Span 4 of 12) */}
              <div className="panel col-span-4 accent-red min-h-500">
                <div className="panel-header">
                  <div className="panel-title">
                    <Clock size={14} className="text-red" />
                    <span>JETWEB RETROACTIVE TIME MACHINE</span>
                  </div>
                </div>
                <div className="panel-content">
                  <div className="time-machine-notice mb-16">
                    <AlertTriangle className="text-red animate-pulse" size={16} />
                    <p>
                      <strong>WARNING:</strong> Triggering a retroactive Time Travel override refactors the entire relational state tree and recalculates downstream hashes under 4/5 consensus.
                    </p>
                  </div>
                  
                  <form onSubmit={handleTimeMutation} className="custom-form">
                    <div className="form-group">
                      <label>Target Ledger Block Index</label>
                      <input 
                        type="number" 
                        value={targetBlockIndex} 
                        onChange={(e) => setTargetBlockIndex(Number(e.target.value))} 
                        required 
                      />
                      <span className="input-hint">Select a block index from the ledger list to patch.</span>
                    </div>

                    <div className="form-group">
                      <label>Retroactive Target Key Mutation</label>
                      <input 
                        type="text" 
                        placeholder="e.g. system/runlevel" 
                        value={timeMutationKey} 
                        onChange={(e) => setTimeMutationKey(e.target.value)} 
                        required 
                      />
                    </div>

                    <div className="form-group">
                      <label>New Proposed Retroactive Value</label>
                      <input 
                        type="text" 
                        placeholder="e.g. 7 (Starbirth mode)" 
                        value={timeMutationVal} 
                        onChange={(e) => setTimeMutationVal(e.target.value)} 
                        required 
                      />
                    </div>

                    <div className="form-group">
                      <label>Override Retroactive Justification Reason</label>
                      <textarea 
                        placeholder="State deviation justification for consensus validation" 
                        value={timeMutationReason} 
                        onChange={(e) => setTimeMutationReason(e.target.value)}
                        rows={3}
                      />
                    </div>

                    <button type="submit" className="danger mt-12 w-full">
                      <span>ENGAGE RETROACTIVE MUTATION OVERRIDE</span>
                    </button>

                    {timeMutationMsg.text && (
                      <div className={`form-message mt-12 ${timeMutationMsg.type}`}>
                        {timeMutationMsg.text}
                      </div>
                    )}
                  </form>
                </div>
              </div>

            </div>
          )}

          {/* TAB 4: MESH TOPOLOGY & PQR NETWORK */}
          {activeTab === 'topology' && (
            <div className="dashboard-grid">
              
              {/* Nodes Topology (Span 6 of 12) */}
              <div className="panel col-span-6 accent-blue min-h-400">
                <div className="panel-header">
                  <div className="panel-title">
                    <Network size={14} className="text-blue" />
                    <span>SOVEREIGN SWARM ACTIVE MESH TOPOLOGY</span>
                  </div>
                </div>
                <div className="panel-content">
                  <div className="topology-nodes">
                    {topology.map((node) => (
                      <div key={node.id} className="node-card">
                        <div className="node-icon-wrapper">
                          <Network size={20} className={node.status === 'ONLINE' ? 'text-green animate-pulse' : 'text-muted'} />
                        </div>
                        <div className="node-details">
                          <div className="node-name">{node.id}</div>
                          <div className="node-ip text-muted">{node.ip}</div>
                          {node.address_5d && (
                            <div className="node-address5d font-mono text-muted mt-8" style={{ fontSize: '11px', borderTop: '1px solid var(--border)', paddingTop: '4px' }}>
                              <div>Coord: {node.address_5d.node_id} | {node.address_5d.role_id}</div>
                              <div>Lineage: {node.address_5d.lineage_id}</div>
                              <div>Block/Thread: {node.address_5d.block_id} / {node.address_5d.thread_id}</div>
                            </div>
                          )}
                        </div>
                        <div className="node-meta">
                          <span className="badge info">{node.role}</span>
                          <span className="node-status text-green glow-text">{node.status}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Tunnels & Certs Info (Span 6 of 12) */}
              <div className="panel col-span-6 accent-yellow min-h-400">
                <div className="panel-header">
                  <div className="panel-title">
                    <Shield size={14} className="text-yellow" />
                    <span>SECURITY CREDENTIALS &amp; ACTIVE PORT TUNNELS</span>
                  </div>
                </div>
                <div className="panel-content text-lines">
                  <h3>🔒 SSL Active Infrastructure</h3>
                  <p className="text-muted">
                    Full SSL TLS secure tunnels initialized on domain <strong>pqr.info</strong>. Certificates are dynamically mapped onto:
                  </p>
                  <ul className="custom-bullets">
                    <li><strong>Full Chain:</strong> <code>certs/pqr.info.fullchain.pem</code></li>
                    <li><strong>Private Key:</strong> <code>certs/pqr.info.key</code></li>
                    <li><strong>Google Cloud Key:</strong> <code>gcp-key.json</code> (Sovereign Service Account)</li>
                  </ul>

                  <h3 className="mt-16">🛰️ Mesh Port Bindings</h3>
                  <div className="port-matrix mt-8">
                    <div className="port-row">
                      <span className="port font-mono">1111</span>
                      <span className="service">gRPC Agent Sync Bus</span>
                      <span className="status-indicator online">ONLINE</span>
                    </div>
                    <div className="port-row">
                      <span className="port font-mono">11111</span>
                      <span className="service">HighSpeed Shared Memory Bus</span>
                      <span className="status-indicator online">ONLINE</span>
                    </div>
                    <div className="port-row">
                      <span className="port font-mono">1113</span>
                      <span className="service">Dedicated Tool-Use Native gRPC</span>
                      <span className="status-indicator online">ONLINE</span>
                    </div>
                    <div className="port-row">
                      <span className="port font-mono">8080</span>
                      <span className="service">Command Bridge (omnibus-gsh)</span>
                      <span className="status-indicator online">ONLINE</span>
                    </div>
                  </div>
                </div>
              </div>

            </div>
          )}

          {selectedTicket && (
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
                      {`$ ${selectedTicket.go_client.command}\n`}
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
}

export default App;
