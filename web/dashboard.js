// dashboard.js - PQR Swarm Dashboard Logic v2.0
const API_BASE = '/REST/2.0';

let state = {
    tickets: [],
    metrics: { used: 0, quota: 1000000, percent: 0 },
    view: 'board' // 'board' or 'lineage'
};

document.addEventListener('DOMContentLoaded', async () => {
    console.log('PQR Dashboard v2.0 Initializing...');
    
    // Initialize Lucide icons
    if (window.lucide) lucide.createIcons();

    setupEventListeners();
    
    // Initial fetch
    await refreshData();
    
    // Refresh loop
    setInterval(refreshData, 5000);
});

async function refreshData() {
    try {
        const [ticketsResp, metricsResp] = await Promise.all([
            fetch(`${API_BASE}/tickets`),
            fetch(`${API_BASE}/metrics/tokens`)
        ]);

        if (ticketsResp.ok) {
            state.tickets = await ticketsResp.json();
            renderBoard();
        }

        if (metricsResp.ok) {
            const data = await metricsResp.json();
            state.metrics = {
                used: data.tokens_used || 0,
                quota: data.token_quota || 1000000,
                percent: data.usage_percentage || 0
            };
            renderMetrics();
        }
    } catch (err) {
        console.error('Data refresh failed:', err);
    }
}

function renderMetrics() {
    const bar = document.getElementById('token-bar');
    const text = document.getElementById('token-text');
    if (bar) bar.style.width = `${state.metrics.percent}%`;
    if (text) text.innerText = `${Math.round(state.metrics.used).toLocaleString()} / ${Math.round(state.metrics.quota).toLocaleString()} tokens`;
}

function setupEventListeners() {
    // New Ticket Modal
    const newTicketBtn = document.getElementById('btn-new-ticket');
    const modal = document.getElementById('modal-new-ticket');
    const cancelBtn = document.getElementById('btn-cancel-ticket');
    const submitBtn = document.getElementById('btn-submit-ticket');

    if (newTicketBtn) newTicketBtn.onclick = () => modal.style.display = 'flex';
    if (cancelBtn) cancelBtn.onclick = () => modal.style.display = 'none';
    
    if (submitBtn) {
        submitBtn.onclick = async () => {
            const subject = document.getElementById('input-subject')?.value;
            const content = document.getElementById('input-content')?.value;
            const priority = parseInt(document.getElementById('input-priority')?.value || "2");

            if (!subject) return alert('Subject is required');

            submitBtn.disabled = true;
            submitBtn.innerText = 'Syncing...';

            try {
                const resp = await fetch(`${API_BASE}/ticket`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        Subject: subject,
                        Text: content,
                        Layer: priority,
                        Queue: 'General',
                        AgentID: 'DASHBOARD-USER'
                    })
                });

                if (resp.ok) {
                    modal.style.display = 'none';
                    document.getElementById('input-subject').value = '';
                    document.getElementById('input-content').value = '';
                    await refreshData();
                }
            } catch (err) {
                console.error(err);
            } finally {
                submitBtn.disabled = false;
                submitBtn.innerText = 'Submit to Swarm';
            }
        };
    }

    // View Switching
    document.querySelectorAll('.nav-item').forEach(item => {
        item.onclick = (e) => {
            e.preventDefault();
            document.querySelectorAll('.nav-item').forEach(i => i.classList.remove('active'));
            item.classList.add('active');
            // Logic for different views could go here
        };
    });
}

function renderBoard() {
    const columns = {
        backlog: { el: document.getElementById('list-backlog'), count: document.getElementById('count-backlog'), tickets: [] },
        todo: { el: document.getElementById('list-todo'), count: document.getElementById('count-todo'), tickets: [] },
        inprogress: { el: document.getElementById('list-inprogress'), count: document.getElementById('count-inprogress'), tickets: [] },
        inreview: { el: document.getElementById('list-inreview'), count: document.getElementById('count-inreview'), tickets: [] },
        done: { el: document.getElementById('list-done'), count: document.getElementById('count-done'), tickets: [] }
    };

    // Distribute tickets
    state.tickets.forEach(ticket => {
        const colKey = mapStatusToColumn(ticket.status);
        if (columns[colKey]) columns[colKey].tickets.push(ticket);
    });

    // Render columns
    Object.keys(columns).forEach(key => {
        const col = columns[key];
        if (!col.el) return;
        
        col.el.innerHTML = '';
        col.count.innerText = col.tickets.length;
        
        col.tickets.forEach(ticket => {
            col.el.appendChild(createTicketCard(ticket));
        });
    });
}

function mapStatusToColumn(status) {
    if (!status) return 'todo';
    const s = status.toUpperCase();
    if (s === 'COMPLETED' || s === 'DONE') return 'done';
    if (s === 'STALLED') return 'backlog';
    if (s === 'IN_PROGRESS' || s === 'HEALING') return 'inprogress';
    if (s === 'PENDING') return 'todo';
    return 'todo';
}

function createTicketCard(ticket) {
    const card = document.createElement('div');
    card.className = 'ticket-card';
    
    const layer = ticket.layer || ticket.layer_id || 2;
    const priority = layer >= 7 ? 'High' : (layer >= 4 ? 'Medium' : 'Low');
    const priorityClass = `priority-${priority.toLowerCase()}`;
    
    const shortId = (ticket.id || '0000').substring(0, 8);
    const subject = ticket.intent?.subject || ticket.subject || "Untitled Sovereign Directive";
    const creator = ticket.creator || ticket.creator_agent_id || "System";
    const avatar = creator.substring(0, 1).toUpperCase();
    
    card.innerHTML = `
        <div class="ticket-id">PQR-${shortId}</div>
        <div class="ticket-subject">${subject}</div>
        <div class="ticket-footer">
            <span class="priority-tag ${priorityClass}">
                <i data-lucide="bar-chart-2" style="width:10px;"></i>
                ${priority}
            </span>
            <div class="agent-avatar" title="Creator: ${creator}">${avatar}</div>
        </div>
    `;
    
    if (window.lucide) {
        setTimeout(() => lucide.createIcons({ props: { "stroke-width": 2 }, root: card }), 0);
    }

    card.onclick = () => openTicketDetail(ticket);

    return card;
}

function openTicketDetail(ticket) {
    // For now, let's just log it or alert. 
    // In a future update we can build a detailed modal.
    console.log('Opening Ticket Detail:', ticket);
    alert(`TICKET DETAIL\nID: ${ticket.id}\nStatus: ${ticket.status}\nLayer: ${ticket.layer_id}\n\nContent: ${ticket.content || 'N/A'}`);
}
