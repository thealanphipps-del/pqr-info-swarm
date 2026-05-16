// dashboard.js - PQR Swarm Dashboard Logic

const API_BASE = '/REST/2.0';

document.addEventListener('DOMContentLoaded', () => {
    fetchTickets();
    fetchMetrics();
    setupEventListeners();
    
    // Refresh every 10 seconds
    setInterval(() => {
        fetchTickets();
        fetchMetrics();
    }, 10000);
});

async function fetchMetrics() {
    try {
        const resp = await fetch(`${API_BASE}/metrics/tokens`);
        const data = await resp.json();
        
        const bar = document.getElementById('token-bar');
        const text = document.getElementById('token-text');
        
        const used = Math.round(data.tokens_used);
        const quota = Math.round(data.token_quota);
        const percent = data.usage_percentage || 0;
        
        bar.style.width = `${percent}%`;
        text.innerText = `${used} / ${quota} used`;
    } catch (err) {
        console.error('Failed to fetch metrics:', err);
    }
}

function setupEventListeners() {
    const newTicketBtn = document.getElementById('btn-new-ticket');
    const modal = document.getElementById('modal-new-ticket');
    const cancelBtn = document.getElementById('btn-cancel-ticket');
    const submitBtn = document.getElementById('btn-submit-ticket');

    newTicketBtn.addEventListener('click', () => modal.style.display = 'flex');
    cancelBtn.addEventListener('click', () => modal.style.display = 'none');
    
    submitBtn.addEventListener('click', async () => {
        const subject = document.getElementById('input-subject').value;
        const content = document.getElementById('input-content').value;
        const priority = parseInt(document.getElementById('input-priority').value);

        if (!subject) return alert('Subject is required');

        try {
            const resp = await fetch(`${API_BASE}/ticket`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    Subject: subject,
                    Text: content,
                    Layer: priority + 1, // Mapping priority to layer
                    Queue: 'General',
                    AgentID: 'DASHBOARD-USER'
                })
            });

            if (resp.ok) {
                modal.style.display = 'none';
                document.getElementById('input-subject').value = '';
                document.getElementById('input-content').value = '';
                fetchTickets();
            }
        } catch (err) {
            console.error('Failed to create ticket:', err);
        }
    });
}

async function fetchTickets() {
    try {
        const resp = await fetch(`${API_BASE}/tickets`);
        const tickets = await resp.json();
        
        renderBoard(tickets);
    } catch (err) {
        console.error('Failed to fetch tickets:', err);
    }
}

function renderBoard(tickets) {
    const columns = {
        backlog: document.getElementById('list-backlog'),
        todo: document.getElementById('list-todo'),
        inprogress: document.getElementById('list-inprogress'),
        inreview: document.getElementById('list-inreview'),
        done: document.getElementById('list-done')
    };

    const counts = {
        backlog: document.getElementById('count-backlog'),
        todo: document.getElementById('count-todo'),
        inprogress: document.getElementById('count-inprogress'),
        inreview: document.getElementById('count-inreview'),
        done: document.getElementById('count-done')
    };

    // Clear lists
    Object.values(columns).forEach(col => col.innerHTML = '');
    
    const stats = { backlog: 0, todo: 0, inprogress: 0, inreview: 0, done: 0 };

    tickets.forEach(ticket => {
        const colKey = mapStatusToColumn(ticket.status);
        if (columns[colKey]) {
            const card = createTicketCard(ticket);
            columns[colKey].appendChild(card);
            stats[colKey]++;
        }
    });

    // Update counts
    Object.keys(counts).forEach(key => counts[key].innerText = stats[key]);
}

function mapStatusToColumn(status) {
    switch (status) {
        case 'COMPLETED': return 'done';
        case 'PENDING': return 'todo';
        case 'STALLED': return 'backlog';
        case 'HEALING': return 'inprogress'; // hypothetical status
        default: return 'todo';
    }
}

function createTicketCard(ticket) {
    const card = document.createElement('div');
    card.className = 'ticket-card';
    
    const priority = ticket.layer <= 2 ? 'Low' : (ticket.layer <= 4 ? 'Medium' : 'High');
    const priorityClass = `priority-${priority.toLowerCase()}`;
    
    // Shorten ID
    const shortId = ticket.id.substring(0, 8);
    
    // Intent snippet or content
    const subject = ticket.intent?.subject || ticket.subject || "Untitled Issue";
    
    card.innerHTML = `
        <div class="ticket-id">PQR-${shortId}</div>
        <div class="ticket-subject">${subject}</div>
        <div class="ticket-footer">
            <span class="priority-tag ${priorityClass}">${priority}</span>
            <div class="agent-avatar" title="${ticket.creator}">${ticket.creator.substring(0,1).toUpperCase()}</div>
        </div>
    `;
    
    card.addEventListener('click', () => {
        window.location.href = `/REST/2.0/ticket/${ticket.id}`; // For now just link to raw JSON or a detail view
    });

    return card;
}
