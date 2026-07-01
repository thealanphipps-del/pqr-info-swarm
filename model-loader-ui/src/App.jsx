import React, { useState, useEffect } from 'react';
import './index.css';

function App() {
  const [activeModel, setActiveModel] = useState('gemma-2-27b-it');
  const [logs, setLogs] = useState([
    { time: '00:00:00', prefix: 'sys', text: 'Initializing Sovereign Mesh node...', type: 'info' },
    { time: '00:00:02', prefix: 'pqr', text: 'Connecting to Cloudflare Tunnel (pqr-sovereign-node)...', type: 'info' },
    { time: '00:00:05', prefix: 'auth', text: 'Vault root token verified successfully.', type: 'success' },
  ]);

  const models = [
    { id: 'gemma-2-27b-it', name: 'Gemma 2 (27B)', desc: 'General reasoning & instruction' },
    { id: 'gemma-2-9b-it', name: 'Gemma 2 (9B)', desc: 'Fast, efficient edge inference' },
    { id: 'pqr-quant-v1', name: 'PQR Quant (v1)', desc: 'Specialized 5D Q-State processor' }
  ];

  useEffect(() => {
    const timer = setInterval(() => {
      const now = new Date().toLocaleTimeString('en-US', { hour12: false });
      setLogs(prev => {
        const newLogs = [...prev, {
          time: now,
          prefix: 'mesh',
          text: `Heartbeat synchronized across ${Math.floor(Math.random() * 20) + 40} nodes.`,
          type: 'info'
        }];
        if (newLogs.length > 20) newLogs.shift();
        return newLogs;
      });
    }, 5000);
    return () => clearInterval(timer);
  }, []);

  const handleLoadModel = () => {
    const now = new Date().toLocaleTimeString('en-US', { hour12: false });
    setLogs(prev => [...prev, {
      time: now,
      prefix: 'loader',
      text: `Allocating VRAM for ${activeModel}. Preparing weights...`,
      type: 'warning'
    }]);
    
    setTimeout(() => {
      const finishTime = new Date().toLocaleTimeString('en-US', { hour12: false });
      setLogs(prev => [...prev, {
        time: finishTime,
        prefix: 'loader',
        text: `${activeModel} loaded successfully into PQR Mesh!`,
        type: 'success'
      }]);
    }, 2500);
  };

  return (
    <>
      <div className="scanline-overlay"></div>
      <div className="app-container">
        
        {/* Sidebar */}
        <div className="sidebar glass-panel">
          <div className="logo-area">
            <div className="logo-icon">PQR</div>
            <div>
              <div className="title">NODE</div>
              <div className="subtitle">Model Loader v2.0</div>
            </div>
          </div>
          
          <div style={{ marginTop: '2rem' }}>
            <div className="section-title">Available Models</div>
            <div className="model-list">
              {models.map(m => (
                <div 
                  key={m.id} 
                  className={`model-item ${activeModel === m.id ? 'active' : ''}`}
                  onClick={() => setActiveModel(m.id)}
                >
                  <div className="model-info">
                    <h4>{m.name}</h4>
                    <p>{m.desc}</p>
                  </div>
                  {activeModel === m.id && (
                    <div style={{ width: '8px', height: '8px', background: 'var(--neon-blue)', borderRadius: '50%', boxShadow: '0 0 10px var(--neon-blue)' }} />
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="main-content glass-panel">
          <div className="header">
            <div style={{ display: 'flex', alignItems: 'baseline', gap: '1rem' }}>
              <h2 style={{ fontSize: '1.2rem', color: 'white', letterSpacing: '1px' }}>Sovereign Mesh Terminal</h2>
              <span style={{ fontFamily: 'JetBrains Mono, monospace', color: '#666', fontSize: '0.8rem' }}>us-east1-b // 10.142.0.4</span>
            </div>
            <div className="status-badge">
              <div className="status-dot"></div>
              TUNNEL CONNECTED
            </div>
          </div>

          <div className="metrics-grid">
            <div className="metric-card">
              <div className="metric-value">47</div>
              <div className="metric-label">Active Nodes</div>
            </div>
            <div className="metric-card">
              <div className="metric-value">12ms</div>
              <div className="metric-label">Mesh Latency</div>
            </div>
            <div className="metric-card">
              <div className="metric-value">99.9%</div>
              <div className="metric-label">PQR Sync Rate</div>
            </div>
          </div>

          <div className="terminal-view" id="terminal">
            {logs.map((log, i) => (
              <div key={i} className="terminal-line">
                <span className="terminal-time">[{log.time}]</span>
                <span className="terminal-prefix">{log.prefix}:~#</span>
                <span className={`terminal-text ${log.type}`}>{log.text}</span>
              </div>
            ))}
          </div>

          <div className="action-bar">
            <button className="btn btn-secondary">Configure Mesh</button>
            <button className="btn btn-primary" onClick={handleLoadModel}>
              Deploy {models.find(m => m.id === activeModel)?.name}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
