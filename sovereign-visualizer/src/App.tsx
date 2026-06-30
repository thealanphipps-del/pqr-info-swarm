import { useEffect, useRef, useState } from 'react';
import './index.css';

type Address5D = { x: number; y: number; z: number; phi: number; lambda: number };
type TesseractVertex = { v0: number; v1: number; v2: number; v3: number };

type AgentState = {
  id: string;
  address5d: Address5D;
  vertex: TesseractVertex;
  epoch: number;
  neighbors: string[];
};

type AnnouncePayload = {
  type: string;
  asn: number;
  ipv6Prefix: string;
  coordinate: Address5D;
  hostname?: string;
  onion?: string;
  timestamp: number;
  status?: string;
};

type CircuitHop = {
  agentId?: string;
  onion?: string;
  hostname?: string;
  ipv6?: string;
  asn?: number;
  coordinate: Address5D;
  vertex: TesseractVertex;
};

type EventCircuit = {
  id: string;
  epoch: number;
  hops: CircuitHop[];
};

type HealthMetrics = {
  epoch: number;
  circuitDensity: number;
  continuityVelocity: number;
  agentChurn: number;
};

type KernelEvent = {
  type: string;
  timestamp: number;
  payload: any;
};

type LineageRecord = {
  clusterId: string;
  epoch: number;
  event: string;
  agents: Record<string, AgentState>;
};

type DepthMode = "2d" | "3d" | "4d";

function Metric({ label, value, maxVal = 1 }: { label: string; value: number, maxVal?: number }) {
  const pct = Math.min(1, value / maxVal);
  return (
    <div className="metric">
      <span>{label}</span>
      <div className="bar">
        <div className="fill" style={{ width: `${pct * 100}%` }} />
      </div>
      <span className="value">{value.toFixed(2)}</span>
    </div>
  );
}

export default function App() {
  const [epoch, setEpoch] = useState(0);
  const [agents, setAgents] = useState<Record<string, AgentState>>({});
  const [asnRegistry, setAsnRegistry] = useState<Record<string, AnnouncePayload>>({});
  const [vertices, setVertices] = useState<Record<string, TesseractVertex>>({});
  const [circuits, setCircuits] = useState<EventCircuit[]>([]);
  const [continuityEvents, setContinuityEvents] = useState<Array<{ epoch: number; agentCount: number; ts: number }>>([]);
  const [lineageEvents, setLineageEvents] = useState<Array<LineageRecord>>([]);
  const [health, setHealth] = useState<HealthMetrics | null>(null);

  const [depthMode, setDepthMode] = useState<DepthMode>("4d");
  const [showLineage, setShowLineage] = useState(false);
  
  const [isScrubbing, setIsScrubbing] = useState(false);
  const [scrubIndex, setScrubIndex] = useState(0);

  const historyRef = useRef<Array<{ epoch: number, vertices: Record<string, TesseractVertex>, asnRegistry: Record<string, AnnouncePayload> }>>([]);
  const asnRegistryRef = useRef<Record<string, AnnouncePayload>>({});

  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [connected, setConnected] = useState(false);

  const phaseRef = useRef(0);
  const unfoldPhaseRef = useRef(0);
  const requestRef = useRef<number>();
  const mouseRef = useRef<{x: number, y: number} | null>(null);

  useEffect(() => {
    const ws = new WebSocket("ws://127.0.0.3:1112/ws");

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);

    ws.onmessage = (ev) => {
      const evt = JSON.parse(ev.data) as KernelEvent;
      switch (evt.type) {
        case "continuity":
          setContinuityEvents(prev => [{
            epoch: evt.payload.epoch,
            agentCount: evt.payload.agentCount,
            ts: evt.timestamp
          }, ...prev].slice(0, 50));
          break;
        case "lineage":
          setLineageEvents(prev => [evt.payload, ...prev].slice(0, 50));
          break;
        case "circuit":
          setCircuits(prev => [evt.payload, ...prev].slice(0, 5));
          break;
        case "health":
          setHealth(evt.payload as HealthMetrics);
          break;
        case "tesseract":
          setEpoch(evt.payload.epoch);
          setVertices(evt.payload.vertices);
          
          historyRef.current.push({
            epoch: evt.payload.epoch,
            vertices: evt.payload.vertices,
            asnRegistry: { ...asnRegistryRef.current }
          });
          if (historyRef.current.length > 500) historyRef.current.shift();
          
          break;
        case "agent_join":
          setAgents(prev => ({
            ...prev,
            [evt.payload.id]: evt.payload
          }));
          break;
        case "5d_announce":
          asnRegistryRef.current[evt.payload.ipv6Prefix] = evt.payload;
          setAsnRegistry(prev => ({
            ...prev,
            [evt.payload.ipv6Prefix]: evt.payload
          }));
          break;
        case "5d_withdraw":
          delete asnRegistryRef.current[evt.payload.ipv6Prefix];
          setAsnRegistry(prev => {
            const next = { ...prev };
            delete next[evt.payload.ipv6Prefix];
            return next;
          });
          break;
        case "epoch":
          setEpoch(evt.payload.epoch);
          break;
      }
    };

    return () => ws.close();
  }, []);

  const renderCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    phaseRef.current = (phaseRef.current + 0.01) % 1;
    unfoldPhaseRef.current = (unfoldPhaseRef.current + 0.005) % 1;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = '#050505';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;
    const SCALE_FACTOR = 8;
    
    let activeVertices = vertices;
    let activeAsnRegistry = asnRegistry;
    let activeEpoch = epoch;

    if (isScrubbing && historyRef.current.length > 0) {
      const snap = historyRef.current[Math.min(scrubIndex, historyRef.current.length - 1)];
      if (snap) {
        activeVertices = snap.vertices;
        activeAsnRegistry = snap.asnRegistry;
        activeEpoch = snap.epoch;
      }
    }

    const project2D = (v: TesseractVertex) => {
      return { x: centerX + v.v0 * SCALE_FACTOR, y: centerY + v.v1 * SCALE_FACTOR, z: 0 };
    };

    const project3D = (v: TesseractVertex) => {
      const scale = 1 / (1 + Math.abs(v.v3) * 0.05);
      return {
        x: centerX + v.v0 * scale * SCALE_FACTOR,
        y: centerY + v.v1 * scale * SCALE_FACTOR,
        z: v.v2 * scale
      };
    };

    const projectAnimated = (v: TesseractVertex, mode: DepthMode, t: number) => {
      if (mode === "2d") return project2D(v);
      if (mode === "3d") return project3D(v);

      const p2 = project2D(v);
      const p3 = project3D(v);
      
      const bounce = Math.sin(t * Math.PI * 2);
      const normalized = (bounce + 1) / 2; 

      const mix = (a: number, b: number) => a * (1 - normalized) + b * normalized;
      return {
        x: mix(p2.x, p3.x),
        y: mix(p2.y, p3.y),
        z: p3.z
      };
    };

    // Draw ASN Regions
    Object.values(activeAsnRegistry).forEach(asn => {
      const v = activeVertices[asn.ipv6Prefix] || { v0: 10, v1: 10, v2: 0, v3: 0 }; // fallback mock
      const p = projectAnimated(v, depthMode, unfoldPhaseRef.current);
      
      ctx.save();
      
      // Draw massive boundary zone for the prefix
      ctx.beginPath();
      ctx.arc(p.x, p.y, 40, 0, 2 * Math.PI);
      ctx.fillStyle = 'rgba(255, 100, 0, 0.05)';
      ctx.fill();
      ctx.strokeStyle = 'rgba(255, 100, 0, 0.3)';
      ctx.setLineDash([5, 5]);
      ctx.stroke();

      // Draw the BGP Edge Router inside it
      ctx.fillStyle = "rgba(255, 150, 0, 0.9)";
      ctx.shadowColor = "rgba(255, 150, 0, 0.8)";
      ctx.shadowBlur = 15;
      ctx.beginPath();
      ctx.rect(p.x - 6, p.y - 6, 12, 12);
      ctx.fill();
      
      ctx.fillStyle = '#ff9600';
      ctx.font = 'bold 10px monospace';
      ctx.fillText(`AS${asn.asn}`, p.x - 15, p.y - 15);
      ctx.fillStyle = '#9ca3af';
      ctx.font = '9px monospace';
      ctx.fillText(asn.ipv6Prefix, p.x + 15, p.y + 4);

      ctx.restore();
    });

    const drawnEdges = new Set<string>();
    Object.entries(activeVertices).forEach(([id, vertex]) => {
      const agent = agents[id];
      if (!agent || !agent.neighbors) return;
      
      const p1 = projectAnimated(vertex, depthMode, unfoldPhaseRef.current);
      
      agent.neighbors.forEach(nid => {
        const nVertex = activeVertices[nid];
        if (!nVertex) return;
        
        const edgeId = [id, nid].sort().join('-');
        if (drawnEdges.has(edgeId)) return;
        drawnEdges.add(edgeId);

        const p2 = projectAnimated(nVertex, depthMode, unfoldPhaseRef.current);

        const dist = Math.sqrt(Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2));
        const alpha = Math.max(0.05, 0.4 - (dist / 400));

        ctx.beginPath();
        ctx.moveTo(p1.x, p1.y);
        ctx.lineTo(p2.x, p2.y);
        ctx.strokeStyle = `rgba(139, 92, 246, ${alpha})`;
        ctx.lineWidth = 1;
        ctx.stroke();
      });
    });

    circuits.forEach(circuit => {
      ctx.save();
      const phase = phaseRef.current;
      
      const externalHopsCount = circuit.hops.filter(h => h.asn).length;
      const isBackbone = externalHopsCount > 1;
      
      if (isBackbone) {
        ctx.strokeStyle = `rgba(255, 0, 170, ${0.5 + 0.5 * Math.sin(phase * Math.PI)})`;
        ctx.shadowColor = "rgba(255, 0, 170, 0.9)";
        ctx.lineWidth = 3 + 2 * Math.sin(phase * Math.PI);
      } else {
        ctx.strokeStyle = `rgba(0, 255, 200, ${0.4 + 0.6 * Math.sin(phase * Math.PI)})`;
        ctx.shadowColor = "rgba(0, 255, 200, 0.8)";
        ctx.lineWidth = 2 + 2 * Math.sin(phase * Math.PI);
      }
      ctx.shadowBlur = 12;

      ctx.beginPath();
      circuit.hops.forEach((hop, idx) => {
        const p = projectAnimated(hop.vertex, depthMode, unfoldPhaseRef.current);
        if (idx === 0) ctx.moveTo(p.x, p.y);
        else ctx.lineTo(p.x, p.y);
      });
      ctx.stroke();
      ctx.restore();
    });

    // Draw Nodes (Agents)
    Object.entries(activeVertices).forEach(([id, vertex]) => {
      if (activeAsnRegistry[id]) return;

      const p = projectAnimated(vertex, depthMode, unfoldPhaseRef.current);
      const depth = p.z ?? 0;
      const radius = Math.max(2, 4 + depth * 0.5);
      const hue = 180 + depth * 20; 
      
      ctx.beginPath();
      ctx.arc(p.x, p.y, radius * 2.5, 0, 2 * Math.PI);
      ctx.fillStyle = `hsla(${hue}, 80%, 60%, 0.2)`;
      ctx.fill();

      ctx.beginPath();
      ctx.arc(p.x, p.y, radius, 0, 2 * Math.PI);
      ctx.fillStyle = `hsla(${hue}, 80%, 60%, 0.9)`;
      ctx.fill();
      
      ctx.fillStyle = '#9ca3af';
      ctx.font = '10px monospace';
      ctx.fillText(id, p.x + radius + 4, p.y + 3);
    });

    // Draw Tooltips
    const mouse = mouseRef.current;
    if (mouse) {
      let hovered: { title: string, lines: string[], x: number, y: number } | null = null;
      
      // Check ASNs
      Object.values(activeAsnRegistry).forEach(asn => {
        const v = activeVertices[asn.ipv6Prefix] || { v0: 10, v1: 10, v2: 0, v3: 0 };
        const p = projectAnimated(v, depthMode, unfoldPhaseRef.current);
        if (Math.hypot(p.x - mouse.x, p.y - mouse.y) < 20) {
          hovered = {
            title: `BGP Autonomous System ${asn.asn}`,
            lines: [
              `IPv6 Prefix: ${asn.ipv6Prefix}`,
              `Hostname: ${asn.hostname || 'N/A'}`,
              `Coordinate: [${asn.coordinate.x}, ${asn.coordinate.y}, ${asn.coordinate.z}, ${asn.coordinate.phi}, ${asn.coordinate.lambda}]`
            ],
            x: mouse.x,
            y: mouse.y
          };
        }
      });

      // Check Agents
      if (!hovered) {
        Object.entries(activeVertices).forEach(([id, vertex]) => {
          if (activeAsnRegistry[id]) return;
          const p = projectAnimated(vertex, depthMode, unfoldPhaseRef.current);
          if (Math.hypot(p.x - mouse.x, p.y - mouse.y) < 15) {
            const agent = agents[id];
            hovered = {
              title: `Agent Node: ${id.slice(0,8)}...`,
              lines: [
                `Full ID: ${id}`,
                agent ? `Epoch: ${agent.epoch}` : 'Status: Unknown',
                agent ? `5D: [${agent.address5d.x}, ${agent.address5d.y}, ${agent.address5d.z}, ${agent.address5d.phi}, ${agent.address5d.lambda}]` : ''
              ].filter(Boolean),
              x: mouse.x,
              y: mouse.y
            };
          }
        });
      }

      if (hovered) {
        ctx.fillStyle = 'rgba(10, 15, 25, 0.95)';
        ctx.strokeStyle = '#00ffc8';
        ctx.lineWidth = 1;
        
        const w = 260;
        const h = 30 + hovered.lines.length * 15;
        const tx = Math.min(hovered.x + 15, canvas.width - w - 10);
        const ty = Math.min(hovered.y + 15, canvas.height - h - 10);
        
        ctx.fillRect(tx, ty, w, h);
        ctx.strokeRect(tx, ty, w, h);
        
        ctx.fillStyle = '#00ffc8';
        ctx.font = 'bold 11px monospace';
        ctx.fillText(hovered.title, tx + 10, ty + 20);
        
        ctx.fillStyle = '#9ca3af';
        ctx.font = '10px monospace';
        hovered.lines.forEach((line, i) => {
          ctx.fillText(line, tx + 10, ty + 35 + i * 15);
        });
      }
    }

    requestRef.current = requestAnimationFrame(renderCanvas);
  };

  useEffect(() => {
    requestRef.current = requestAnimationFrame(renderCanvas);
    return () => {
      if (requestRef.current) cancelAnimationFrame(requestRef.current);
    };
  }, [vertices, agents, asnRegistry, circuits, depthMode]);

  return (
    <div className="layout">
      {/* Top Bar */}
      <header className="topbar">
        <div className="brand">Sovereign Mesh <span className="badge" style={{ background: isScrubbing ? '#ff00aa' : '#00ffc8' }}>{isScrubbing ? 'REPLAY' : 'LIVE'}</span></div>
        <div className="stats">
          <div className="stat">
            <span className="label">EPOCH</span>
            <span className="value">{isScrubbing ? activeEpoch : epoch}</span>
          </div>
          <div className="stat">
            <span className="label">BGP ZONES (ASN)</span>
            <span className="value" style={{ color: '#ff9600' }}>{Object.keys(asnRegistry).length}</span>
          </div>
          <div className="stat">
            <span className="label">ACTIVE CIRCUITS</span>
            <span className="value">{circuits.length}</span>
          </div>
          <div className="stat">
            <span className="label">CONNECTION</span>
            <span className={`value ${connected ? 'success' : 'error'}`}>
              {connected ? 'CONNECTED' : 'DISCONNECTED'}
            </span>
          </div>
          <button 
            className={`btn ${showLineage ? 'active' : ''}`} 
            style={{ marginLeft: '1rem', border: '1px solid #00ffc8' }}
            onClick={() => setShowLineage(!showLineage)}
          >
            LINEAGE EXPLORER
          </button>
        </div>
      </header>

      <div className="main">
        {/* Left Panel */}
        <aside className="panel left-panel" style={{ display: 'flex', flexDirection: 'column' }}>
          
          <div style={{ flex: 1, overflowY: 'auto' }}>
            <h2>Agent Roster</h2>
            <div className="agent-list">
              {Object.values(agents).map(agent => (
                <div key={agent.id} className="agent-card">
                  <div className="agent-header">
                    <span className="agent-id">{agent.id}</span>
                    <span className="pulse-dot" style={{ backgroundColor: 'hsla(180, 80%, 60%, 0.9)' }}></span>
                  </div>
                  <div className="agent-coords">
                    5D: [{agent.address5d.x}, {agent.address5d.y}, {agent.address5d.z}, {agent.address5d.phi}, {agent.address5d.lambda}]
                  </div>
                </div>
              ))}
              {Object.keys(agents).length === 0 && <div className="empty">No agents registered</div>}
            </div>
          </div>

          <div style={{ flex: 1, overflowY: 'auto', borderTop: '1px solid rgba(255,255,255,0.1)' }}>
            <h2>Sovereign AS Fabric (BGP)</h2>
            <div className="agent-list">
              {Object.values(asnRegistry).map(asn => (
                <div key={asn.ipv6Prefix} className="agent-card" style={{ borderLeft: '3px solid #ff9600' }}>
                  <div className="agent-header">
                    <span className="agent-id" style={{ color: '#ff9600' }}>AS{asn.asn}</span>
                  </div>
                  <div className="agent-coords">
                    PREFIX: {asn.ipv6Prefix}<br/>
                    ROUTER: {asn.hostname || 'unknown'}<br/>
                  </div>
                </div>
              ))}
              {Object.keys(asnRegistry).length === 0 && <div className="empty">No BGP ASN Peers Discovered</div>}
            </div>
          </div>
        </aside>

        {/* Center Canvas */}
        <section className="canvas-container">
          <div className="canvas-wrapper">
            <canvas 
              ref={canvasRef} 
              width={800} 
              height={600} 
              onMouseMove={(e) => {
                const rect = canvasRef.current?.getBoundingClientRect();
                if (rect) {
                  mouseRef.current = {
                    x: e.clientX - rect.left,
                    y: e.clientY - rect.top
                  };
                }
              }}
              onMouseLeave={() => mouseRef.current = null}
            />
            <div className="canvas-overlay">
              <span>TESSERACT EVOLUTION: {depthMode.toUpperCase()}</span>
            </div>
            
            {health && (
              <div className="hud">
                <h3 className="hud-title">Mesh Vital Signs</h3>
                {health.continuityVelocity < 1.0 && <div style={{color: '#ff0055', fontSize: '10px', marginBottom: '5px'}}>⚠️ WARNING: Network Partitioning Detected</div>}
                {health.circuitDensity > 50 && <div style={{color: '#ff9600', fontSize: '10px', marginBottom: '5px'}}>⚠️ WARNING: High Circuit Congestion</div>}
                <Metric label="Circuit Density" value={health.circuitDensity} />
                <Metric label="Continuity Velocity" value={health.continuityVelocity} maxVal={5} />
                <Metric label="Agent Churn (join/s)" value={health.agentChurn} maxVal={1} />
              </div>
            )}
          </div>
        </section>

        {/* Right Panel */}
        <aside className="panel right-panel">
          <h2>Event Stream</h2>
          <div className="event-list">
            {circuits.map((ev, i) => {
                const asnHops = ev.hops.filter(h => h.asn);
                const isBackbone = asnHops.length > 1;
                const dest = asnHops.length > 0 ? `AS${asnHops.map(h => h.asn).join(' -> AS')}` : "Internal";
                return (
                  <div key={`circ-${i}`} className="event-card" style={{ borderLeft: `3px solid ${isBackbone ? '#ff00aa' : '#00ffc8'}` }}>
                    <div className="event-time">EPOCH {ev.epoch}</div>
                    <div className="event-body">
                      <strong>{isBackbone ? "Backbone Transit Route" : "SRRP Circuit Formed"}</strong>
                      <br/>
                      <span style={{fontSize: '0.75rem', color: '#9ca3af'}}>Hops: {ev.hops.length} | Path: {dest}</span>
                    </div>
                  </div>
                );
            })}
            {continuityEvents.map((ev, i) => (
              <div key={`cont-${i}`} className="event-card type-continuity">
                <div className="event-time">{new Date(ev.ts).toLocaleTimeString()}</div>
                <div className="event-body">
                  <strong>Continuity</strong> emitted at Epoch {ev.epoch}
                </div>
              </div>
            ))}
            {lineageEvents.map((ev, i) => (
              <div key={`lin-${i}`} className="event-card type-lineage">
                <div className="event-time">EPOCH {ev.epoch}</div>
                <div className="event-body">
                  <strong>Lineage:</strong> {ev.event}
                </div>
              </div>
            ))}
          </div>
        </aside>
      </div>

      {/* Bottom Bar */}
      <footer className="bottombar" style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
        {historyRef.current.length > 0 && (
          <div className="scrubber" style={{ width: '100%', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <span style={{color: '#00ffc8', fontSize: '10px', whiteSpace: 'nowrap'}}>EPOCH REPLAY</span>
            <input 
              type="range" 
              min="0" 
              max={historyRef.current.length - 1} 
              value={isScrubbing ? scrubIndex : historyRef.current.length - 1}
              onChange={(e) => {
                setIsScrubbing(true);
                setScrubIndex(parseInt(e.target.value));
              }}
              style={{ flex: 1, accentColor: '#ff00aa' }}
            />
            <button className="btn" style={{ padding: '2px 8px', fontSize: '10px' }} onClick={() => setIsScrubbing(false)}>
              {isScrubbing ? 'RESUME LIVE' : 'LIVE'}
            </button>
          </div>
        )}
        <div className="controls">
          <button 
            className={`btn ${depthMode === '2d' ? 'active' : ''}`}
            onClick={() => setDepthMode('2d')}
          >
            2D Projection
          </button>
          <button 
            className={`btn ${depthMode === '3d' ? 'active' : ''}`}
            onClick={() => setDepthMode('3d')}
          >
            3D Depth
          </button>
          <button 
            className={`btn ${depthMode === '4d' ? 'active' : ''}`}
            onClick={() => setDepthMode('4d')}
          >
            4D Unfolding
          </button>
        </div>
      </footer>

      {showLineage && (
        <div className="lineage-overlay">
          <div className="lineage-header">
            <h2>HyperFrame Lineage Chains</h2>
            <button className="btn" onClick={() => setShowLineage(false)}>Close</button>
          </div>
          <div className="lineage-timeline">
            {lineageEvents.length === 0 ? <div className="empty">No lineage events recorded yet.</div> : null}
            {lineageEvents.map((ev, idx) => (
              <div key={`${ev.epoch}-${idx}`} className="lineage-node">
                <div className="lineage-node-marker"></div>
                <div className="lineage-node-content">
                  <div className="lineage-epoch">Epoch {ev.epoch}</div>
                  <div className="lineage-event">{ev.event}</div>
                  <div className="lineage-cluster">Cluster ID: <code>{ev.clusterId || 'N/A'}</code></div>
                  <div className="lineage-agents">
                    <strong>Agents involved:</strong> {Object.keys(ev.agents || {}).length}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
