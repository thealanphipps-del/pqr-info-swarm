import React, { useState, useEffect, useRef } from 'react';
import './SlingshotMixer.css';

type TrackInfo = {
  id: string;
  name: string;
  service: string;
  sequenceStart: string;
  sequenceEnd: string;
  volume: number;
  activity: number;
  baseFreq: number;
  hitlRequired: boolean;
  hitlAnomaly: string;
};

export const SlingshotMixer: React.FC = () => {
  const [crossfader, setCrossfader] = useState(50);
  const [powerOn, setPowerOn] = useState(false);
  const audioCtxRef = useRef<AudioContext | null>(null);
  
  const [tracks, setTracks] = useState<TrackInfo[]>([
    { id: 'trk-1', name: 'CH 1', service: 'SRRK Router', sequenceStart: '0x00A1', sequenceEnd: '0x09FF', volume: 80, activity: 0, baseFreq: 220, hitlRequired: false, hitlAnomaly: '' }, 
    { id: 'trk-2', name: 'CH 2', service: '5D Manifold', sequenceStart: '0x1B22', sequenceEnd: '0x2C4A', volume: 65, activity: 0, baseFreq: 329.63, hitlRequired: false, hitlAnomaly: '' }, 
    { id: 'trk-3', name: 'CH 3', service: 'Q-State Ledger', sequenceStart: '0xF400', sequenceEnd: '0xFFFF', volume: 90, activity: 0, baseFreq: 440, hitlRequired: false, hitlAnomaly: '' }, 
    { id: 'trk-4', name: 'CH 4', service: 'Midi Handshake', sequenceStart: '0x33A1', sequenceEnd: '0x33B9', volume: 40, activity: 0, baseFreq: 554.37, hitlRequired: false, hitlAnomaly: '' }, 
  ]);

  const togglePower = () => {
    if (!powerOn) {
      const AudioContextClass = window.AudioContext || (window as any).webkitAudioContext;
      audioCtxRef.current = new AudioContextClass();
    } else if (audioCtxRef.current) {
      audioCtxRef.current.close();
      audioCtxRef.current = null;
    }
    setPowerOn(!powerOn);
  };

  const playSynthBeep = (freq: number, volumeLevel: number, activityLevel: number, isHitl: boolean) => {
    if (!audioCtxRef.current || !powerOn) return;
    
    // HITL alarm is a dissonant tone
    if (isHitl) {
      freq = freq * 1.5; 
    } else if (activityLevel < 70) {
      return; // Skip normal play if activity is low
    }

    const osc = audioCtxRef.current.createOscillator();
    const gainNode = audioCtxRef.current.createGain();
    
    osc.type = isHitl ? 'sawtooth' : (crossfader < 50 ? 'sine' : 'square');
    
    osc.frequency.setValueAtTime(freq + (activityLevel / 2), audioCtxRef.current.currentTime);
    
    const vol = (volumeLevel / 100) * (activityLevel / 100) * (isHitl ? 0.3 : 0.1); 
    gainNode.gain.setValueAtTime(vol, audioCtxRef.current.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.001, audioCtxRef.current.currentTime + (isHitl ? 0.3 : 0.1));
    
    osc.connect(gainNode);
    gainNode.connect(audioCtxRef.current.destination);
    
    osc.start();
    osc.stop(audioCtxRef.current.currentTime + (isHitl ? 0.3 : 0.1));
  };

  useEffect(() => {
    const interval = setInterval(() => {
      setTracks(prev => prev.map(track => {
        if (track.hitlRequired) {
          // Flash max activity when frozen in HITL
          const pulse = Math.random() > 0.5 ? 100 : 0;
          if (powerOn && pulse > 0) playSynthBeep(track.baseFreq, track.volume, 100, true);
          return { ...track, activity: pulse };
        }

        // Random chance to trigger a HITL anomaly if power is on
        if (powerOn && Math.random() > 0.99) {
          const anomalies = ["CRDT COLLISION", "PHASE-STATE DRIFT", "SEQUENCE MISMATCH"];
          return { ...track, hitlRequired: true, hitlAnomaly: anomalies[Math.floor(Math.random()*anomalies.length)] };
        }

        const newActivity = Math.random() * 100 * (track.volume / 100);
        if (powerOn) {
          playSynthBeep(track.baseFreq, track.volume, newActivity, false);
        }
        return { ...track, activity: newActivity };
      }));
    }, 150);
    return () => clearInterval(interval);
  }, [powerOn, crossfader]);

  const handleVolumeChange = (id: string, newVolume: number) => {
    setTracks(prev => prev.map(t => t.id === id && !t.hitlRequired ? { ...t, volume: newVolume } : t));
  };

  const resolveHITL = (id: string, action: string) => {
    console.log(`[API MOCK] POST /hitl/resolve/${id} - Action: ${action}`);
    setTracks(prev => prev.map(t => t.id === id ? { ...t, hitlRequired: false, hitlAnomaly: '', activity: 0 } : t));
  };

  return (
    <div className="slingshot-mixer">
      <div className="mixer-header">
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          <h2>SLINGSHOT MERGE DESK</h2>
          <span style={{ fontSize: '0.7rem', color: '#0088ff', letterSpacing: '1.5px', textTransform: 'uppercase' }}>
            MIDI Transcoder Buffered // Mix Live
          </span>
        </div>
        <button 
          onClick={togglePower} 
          className={`status-badge ${powerOn ? 'active' : ''}`}
          style={{ cursor: 'pointer', background: powerOn ? 'rgba(0, 255, 213, 0.2)' : 'transparent', border: powerOn ? '1px solid #00ffd5' : '1px solid #555', color: powerOn ? '#00ffd5' : '#555' }}
        >
          {powerOn ? 'STATION SYNTHESIZED // PHASE LOCKED' : 'SYNTHESIZE MY STATION - CLICK TO PHASE LOCK'}
        </button>
      </div>

      <div className="mixer-deck">
        {tracks.map(track => (
          <div key={track.id} className={`channel-strip ${track.hitlRequired ? 'hitl-alarm' : ''}`} style={track.hitlRequired ? { border: '1px solid #ff2a2a', boxShadow: '0 0 15px rgba(255,42,42,0.4)' } : {}}>
            <div className="channel-display">
              <div className="service-name">{track.service}</div>
              <div className="seq-range">
                <span>{track.sequenceStart}</span> ➔ <span>{track.sequenceEnd}</span>
              </div>
            </div>

            {track.hitlRequired ? (
              <div style={{ display: 'flex', flexDirection: 'column', height: '100%', justifyContent: 'center', gap: '10px', padding: '10px' }}>
                <div style={{ color: '#ff2a2a', fontSize: '0.65rem', fontWeight: 'bold', textAlign: 'center', letterSpacing: '1px' }}>
                  HITL REQUIRED<br/>{track.hitlAnomaly}
                </div>
                <button 
                  onClick={() => resolveHITL(track.id, 'APPROVE')}
                  style={{ background: '#00ffd5', color: '#000', border: 'none', padding: '5px', fontSize: '0.6rem', fontWeight: 'bold', cursor: 'pointer', borderRadius: '3px' }}
                >
                  APPROVE REPLAY
                </button>
                <button 
                  onClick={() => resolveHITL(track.id, 'REJECT')}
                  style={{ background: 'transparent', color: '#ff2a2a', border: '1px solid #ff2a2a', padding: '5px', fontSize: '0.6rem', fontWeight: 'bold', cursor: 'pointer', borderRadius: '3px' }}
                >
                  REJECT MERGE
                </button>
              </div>
            ) : (
              <>
                <div className="eq-section">
                  <div style={{ fontSize: '0.55rem', color: '#0088ff', textAlign: 'center', marginBottom: '8px', letterSpacing: '1px' }}>EQUALIZE<br/>SUBSTRATUM STATE</div>
                  <div className="knob-container">
                    <div className="knob" style={{ transform: `rotate(${Math.random() * 270 - 135}deg)` }}></div>
                    <label>HI (Sync)</label>
                  </div>
                  <div className="knob-container">
                    <div className="knob" style={{ transform: `rotate(${Math.random() * 270 - 135}deg)` }}></div>
                    <label>MID (CRDT)</label>
                  </div>
                  <div className="knob-container">
                    <div className="knob" style={{ transform: `rotate(${Math.random() * 270 - 135}deg)` }}></div>
                    <label>LOW (Cache)</label>
                  </div>
                </div>

                <div className="vu-meter">
                  <div 
                    className="vu-level" 
                    style={{ 
                      height: `${track.activity}%`,
                      backgroundColor: track.activity > 85 ? '#ff2a2a' : track.activity > 60 ? '#f0db4f' : '#00ffd5'
                    }}
                  ></div>
                </div>

                <div className="fader-container">
                  <input 
                    type="range" 
                    min="0" 
                    max="100" 
                    value={track.volume} 
                    onChange={(e) => handleVolumeChange(track.id, parseInt(e.target.value))}
                    className="fader"
                  />
                </div>
              </>
            )}
            
            <div className="channel-label">{track.name}</div>
          </div>
        ))}
      </div>

      <div className="crossfader-section">
        <span className="deck-label">OFFLINE CACHE</span>
        <input 
          type="range" 
          min="0" 
          max="100" 
          value={crossfader} 
          onChange={(e) => setCrossfader(parseInt(e.target.value))}
          className="crossfader"
        />
        <span className="deck-label">LIVE PLATTER</span>
      </div>
    </div>
  );
};
