// dashboard.js - PQR Swarm Dashboard Logic v2.0
const API_BASE = '/REST/2.0';

let state = {
    tickets: [],
    metrics: { used: 0, quota: 1000000, percent: 0 },
    view: 'board', // 'board' or 'list'
    filter: null
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
            
            const span = item.querySelector('span');
            if (span) {
                const viewName = span.innerText.toLowerCase();
                state.view = (viewName === 'board' || viewName === 'inbox') ? 'board' : 'list';
                if (viewName === 'inbox') {
                    state.filter = 'PENDING';
                } else {
                    state.filter = null;
                }
                renderBoard();
            }
        };
    });

    // Filter button
    const filterBtn = document.querySelector('.topbar-actions .btn:first-child');
    if (filterBtn) {
        filterBtn.onclick = () => {
            const val = prompt("Filter by status (PENDING, COMPLETED, etc) or priority (High, Medium, Low):");
            if (val) {
                state.filter = val.toUpperCase();
                renderBoard();
            } else {
                state.filter = null;
                renderBoard();
            }
        };
    }
}

function renderBoard() {
    const boardEl = document.querySelector('.board');
    
    if (state.view === 'list') {
        renderListView();
        return;
    }

    // Restore board structure if it was replaced by list view
    if (!document.getElementById('list-backlog')) {
        boardEl.innerHTML = `
            <div class="column" id="col-backlog">
                <div class="column-header"><div class="column-title"><i data-lucide="circle" style="width:14px; color:var(--status-backlog);"></i>Backlog<span class="column-badge" id="count-backlog">0</span></div></div>
                <div class="card-list" id="list-backlog"></div>
            </div>
            <div class="column" id="col-todo">
                <div class="column-header"><div class="column-title"><i data-lucide="circle" style="width:14px; color:var(--status-todo);"></i>Todo<span class="column-badge" id="count-todo">0</span></div></div>
                <div class="card-list" id="list-todo"></div>
            </div>
            <div class="column" id="col-inprogress">
                <div class="column-header"><div class="column-title"><i data-lucide="play-circle" style="width:14px; color:var(--status-inprogress);"></i>In Progress<span class="column-badge" id="count-inprogress">0</span></div></div>
                <div class="card-list" id="list-inprogress"></div>
            </div>
            <div class="column" id="col-inreview">
                <div class="column-header"><div class="column-title"><i data-lucide="check-circle-2" style="width:14px; color:var(--status-inreview);"></i>In Review<span class="column-badge" id="count-inreview">0</span></div></div>
                <div class="card-list" id="list-inreview"></div>
            </div>
            <div class="column" id="col-done">
                <div class="column-header"><div class="column-title"><i data-lucide="check-circle" style="width:14px; color:var(--status-done);"></i>Done<span class="column-badge" id="count-done">0</span></div></div>
                <div class="card-list" id="list-done"></div>
            </div>
        `;
        if (window.lucide) lucide.createIcons();
    }

    const columns = {
        backlog: { el: document.getElementById('list-backlog'), count: document.getElementById('count-backlog'), tickets: [] },
        todo: { el: document.getElementById('list-todo'), count: document.getElementById('count-todo'), tickets: [] },
        inprogress: { el: document.getElementById('list-inprogress'), count: document.getElementById('count-inprogress'), tickets: [] },
        inreview: { el: document.getElementById('list-inreview'), count: document.getElementById('count-inreview'), tickets: [] },
        done: { el: document.getElementById('list-done'), count: document.getElementById('count-done'), tickets: [] }
    };

    // Filter and Distribute tickets
    state.tickets.forEach(ticket => {
        const layer = ticket.layer || ticket.layer_id || 2;
        const priority = layer >= 7 ? 'HIGH' : (layer >= 4 ? 'MEDIUM' : 'LOW');
        
        if (state.filter) {
            const s = ticket.status?.toUpperCase();
            if (state.filter !== s && state.filter !== priority) return;
        }

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

function renderListView() {
    const boardEl = document.querySelector('.board');
    boardEl.innerHTML = `
        <div style="width:100%; padding: 2rem; overflow-y:auto;">
            <table style="width:100%; border-collapse: collapse; background:white; border-radius:8px; overflow:hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
                <thead>
                    <tr style="background:#f8f9fc; border-bottom: 2px solid #eaecf4; text-align:left;">
                        <th style="padding:1rem;">ID</th>
                        <th style="padding:1rem;">Subject</th>
                        <th style="padding:1rem;">Status</th>
                        <th style="padding:1rem;">Priority</th>
                        <th style="padding:1rem;">Created</th>
                    </tr>
                </thead>
                <tbody id="list-view-body"></tbody>
            </table>
        </div>
    `;
    
    const body = document.getElementById('list-view-body');
    state.tickets.forEach(ticket => {
        const row = document.createElement('tr');
        row.style.borderBottom = "1px solid #eaecf4";
        row.style.cursor = "pointer";
        row.onclick = () => openTicketDetail(ticket);
        
        const layer = ticket.layer || ticket.layer_id || 2;
        const priority = layer >= 7 ? 'High' : (layer >= 4 ? 'Medium' : 'Low');

        row.innerHTML = `
            <td style="padding:1rem; font-weight:600; color:var(--sidebar-active);">PQR-${ticket.id.substring(0,8)}</td>
            <td style="padding:1rem;">${ticket.intent?.subject || ticket.subject || "Untitled Sovereign Directive"}</td>
            <td style="padding:1rem;"><span class="status-tag status-${ticket.status?.toLowerCase()}">${ticket.status}</span></td>
            <td style="padding:1rem;">${priority}</td>
            <td style="padding:1rem; color:#888;">${new Date(ticket.created_at).toLocaleString()}</td>
        `;
        body.appendChild(row);
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

async function resolveTicket(ticketId) {
    if (!confirm('Resolve this ticket?')) return;
    
    try {
        const resp = await fetch(`${API_BASE}/ticket/${ticketId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                Status: 'COMPLETED',
                Creator: 'DASHBOARD-USER'
            })
        });
        if (resp.ok) {
            await refreshData();
        }
    } catch (err) {
        console.error(err);
    }
}

function openTicketDetail(ticket) {
    const detail = `ID: ${ticket.id}\nStatus: ${ticket.status}\nCreator: ${ticket.creator || ticket.creator_agent_id}\n\nContent: ${ticket.content || 'No detailed content available.'}`;
    if (confirm(`TICKET DETAIL\n\n${detail}\n\nClick OK to resolve this ticket, or Cancel to close.`)) {
        resolveTicket(ticket.id);
    }
}
