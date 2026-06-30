// app.js – premium vanilla‑JS chat client
// Connects to the Sovereign multi‑model backend at http://localhost:4000/api/chat/send
// Handles user input, displays messages with smooth animations, and auto‑scrolls.

const API_URL = "http://localhost/api/chat/send"; // backend now on port 80

const chatWindow = document.getElementById("chat-window");
const chatForm = document.getElementById("chat-form");
const messageInput = document.getElementById("message-input");

// Helper to create message bubbles
function createMessageBubble(content, role) {
  const msgEl = document.createElement("div");
  msgEl.classList.add("message", role);

  const bubble = document.createElement("div");
  bubble.classList.add("bubble");
  bubble.textContent = content;
  msgEl.appendChild(bubble);
  return msgEl;
}

// Append a message to the chat window and keep it scrolled to bottom
function appendMessage(content, role) {
  const bubble = createMessageBubble(content, role);
  chatWindow.appendChild(bubble);
  // Slight delay for the fade‑in animation to be visible
  setTimeout(() => {
    chatWindow.scrollTop = chatWindow.scrollHeight;
  }, 50);
}

// Disable the form while waiting for a response
function setFormEnabled(enabled) {
  messageInput.disabled = !enabled;
  document.getElementById("send-button").disabled = !enabled;
}

chatForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const userMessage = messageInput.value.trim();
  if (!userMessage) return;

  // Show user message immediately
  appendMessage(userMessage, "user");
  messageInput.value = "";
  setFormEnabled(false);

  // Build payload – minimal fields required by backend
  const payload = {
    threadId: null, // let backend create a new thread if needed
    userId: "anonymous",
    message: userMessage,
    participants: [],
    mode: "single"
  };

  try {
    const resp = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    const data = await resp.json();
    // Backend returns an object – assume it contains a `reply` field with the assistant text
    const reply = data.reply || data.message || JSON.stringify(data);
    appendMessage(reply, "assistant");
  } catch (err) {
    console.error(err);
    appendMessage(`⚠️ Error: ${err.message}`, "assistant");
  } finally {
    setFormEnabled(true);
    messageInput.focus();
  }
});

// Auto‑grow textarea height
messageInput.addEventListener("input", () => {
  messageInput.style.height = "auto";
  messageInput.style.height = `${messageInput.scrollHeight}px`;
});

// Initial focus for better UX
messageInput.focus();
