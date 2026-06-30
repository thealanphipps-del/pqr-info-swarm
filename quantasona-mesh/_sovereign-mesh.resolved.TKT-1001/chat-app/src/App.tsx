import React, { useState } from 'react';
import './index.css';

interface Message {
  role: 'user' | 'bot';
  content: string;
}

const App: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);

  const sendMessage = async () => {
    if (!input.trim()) return;
    const userMsg: Message = { role: 'user', content: input };
    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setLoading(true);
    try {
      const response = await fetch('/api/v1/infer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ prompt: input, model: 'gemma-2b', provider: 'LM_STUDIO' })
      });
      const data = await response.json();
      const botMsg: Message = { role: 'bot', content: data.text || 'No response' };
      setMessages(prev => [...prev, botMsg]);
    } catch (e) {
      const errMsg: Message = { role: 'bot', content: 'Error contacting backend' };
      setMessages(prev => [...prev, errMsg]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className="app-container">
      <header>Gemma Chat</header>
      <div className="chat-window" id="chat-window">
        {messages.map((msg, idx) => (
          <div key={idx} className={`message ${msg.role}`}>{msg.content}</div>
        ))}
      </div>
      <div className="input-area">
        <input
          type="text"
          placeholder="Type a message…"
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          disabled={loading}
        />
        <button className="send" onClick={sendMessage} disabled={loading}>Send</button>
      </div>
    </div>
  );
};

export default App;
