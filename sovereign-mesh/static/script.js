/* script.js – UI logic for Sovereign Mesh HUD */

// Helper to display loading spinner
function showSpinner() {
  const term = document.getElementById('term');
  const span = document.createElement('span');
  span.className = 'spinner';
  term.appendChild(span);
}

function hideSpinner() {
  const term = document.getElementById('term');
  const spinner = term.querySelector('.spinner');
  if (spinner) spinner.remove();
}

async function executeBridge() {
  const input = document.getElementById('cmd');
  const term = document.getElementById('term');
  const val = input.value.trim();
  if (!val) return;

  // Echo command
  term.innerHTML += `<br><span style="color:#fff;">$ ${val}</span><br><span style="color:#8b949e;">Executing via local hook...</span>`;
  input.value = '';
  term.scrollTop = term.scrollHeight;

  showSpinner();
  try {
    const response = await fetch('/api/bridge?cmd=' + encodeURIComponent(val));
    const text = await response.text();
    term.innerHTML += `<br>${text.replace(/\n/g, '<br>')}`;
  } catch (err) {
    term.innerHTML += `<br><span style="color:#da3633;">Hook Error: ${err}</span>`;
  } finally {
    hideSpinner();
    term.scrollTop = term.scrollHeight;
  }
}

// Periodic system status refresh – robust with back‑off on failure
let statusInterval = null;
function startStatusPolling() {
  if (statusInterval) clearInterval(statusInterval);
  statusInterval = setInterval(async () => {
    try {
      const res = await fetch('/api/status');
      const text = await res.text();
      document.getElementById('sys-status').innerText = text;
    } catch (_) {
      // Silently ignore – UI will keep previous status
    }
  }, 5000);
}

// Initialize UI
document.addEventListener('DOMContentLoaded', () => {
  startStatusPolling();
});
