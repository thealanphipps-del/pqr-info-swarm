/* ============================================================
   Sovereign HUD v8.5 — REST 2.0 Client + GEDCOM Genealogy
   ============================================================ */

document.addEventListener('DOMContentLoaded', () => {
    // ── DOM refs ──────────────────────────────────────────────
    const hexGrid = document.getElementById('hexGrid');
    const hexContainer = document.getElementById('hexContainer');
    const logStream = document.getElementById('logStream');
    const ticketList = document.getElementById('ticketList');
    const ticketCount = document.getElementById('ticketCount');
    const headerCount = document.getElementById('headerTicketCount');
    const resetViewBtn = document.getElementById('resetView');
    const clearLogsBtn = document.getElementById('clearLogs');
    const refreshTreeBtn = document.getElementById('refreshTree');
    const gedcomSvg = document.getElementById('gedcomSvg');
    const gedcomContainer = document.getElementById('gedcomContainer');

    // ── Zoom/Pan State ────────────────────────────────────────
    let scale = 1, tx = 0, ty = 0, dragging = false, sx, sy;

    // ── Tab Switcher ──────────────────────────────────────────
    document.querySelectorAll('.tab').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
            btn.classList.add('active');
            document.getElementById('tab-' + btn.dataset.tab).classList.add('active');
            if (btn.dataset.tab === 'gedcom') loadGedcomTree();
        });
    });

    // ── Generate 64 Hex Gates ─────────────────────────────────
    for (let i = 1; i <= 64; i++) {
        const hex = document.createElement('div');
        hex.className = 'hex';
        hex.id = `gate-${String(i).padStart(2, '0')}`;
        hex.innerHTML = `<div class="gate-label">GATE</div>${String(i).padStart(2, '0')}`;
        hex.addEventListener('click', e => {
            e.stopPropagation();
            hex.classList.toggle('active');
            addLog(`[SWARM] Gate ${i} ${hex.classList.contains('active') ? 'activated' : 'deactivated'}.`);
        });
        hexGrid.appendChild(hex);
    }

    // Swarm pulse animation
    setInterval(() => {
        const i = Math.floor(Math.random() * 64) + 1;
        const el = document.getElementById(`gate-${String(i).padStart(2, '0')}`);
        if (!el) return;
        el.classList.add('pulsing');
        setTimeout(() => el.classList.remove('pulsing'), 1400);
    }, 400);

    // ── Hex Zoom/Pan ──────────────────────────────────────────
    hexContainer.addEventListener('wheel', e => {
        e.preventDefault();
        scale = Math.min(3, Math.max(0.2, scale * (e.deltaY > 0 ? 0.9 : 1.1)));
        applyTransform();
    }, { passive: false });

    hexContainer.addEventListener('mousedown', e => {
        dragging = true; sx = e.clientX - tx; sy = e.clientY - ty;
    });
    window.addEventListener('mousemove', e => {
        if (!dragging) return;
        tx = e.clientX - sx; ty = e.clientY - sy; applyTransform();
    });
    window.addEventListener('mouseup', () => dragging = false);
    resetViewBtn.addEventListener('click', () => {
        scale = 1; tx = 0; ty = 0; applyTransform();
    });
    function applyTransform() {
        hexGrid.style.transform = `translate(${tx}px,${ty}px) scale(${scale})`;
    }

    // ── SSE — Real-time Ticket Stream ─────────────────────────
    const es = new EventSource('/api/tickets/stream');
    es.onmessage = e => {
        const tickets = JSON.parse(e.data) || [];
        renderTicketList(tickets);
    };
    es.addEventListener('error', () => {
        addLog('[SSE] Stream error — reconnecting…');
    });

    function renderTicketList(tickets) {
        ticketList.innerHTML = '';
        const n = tickets.length;
        ticketCount.textContent = n;
        headerCount.textContent = n;

        (tickets || []).forEach(t => {
            const item = document.createElement('div');
            item.className = `ticket-item priority-${t.priority || 20}`;
            const badge = badgeFor(t.title);
            item.innerHTML = `
                <div class="ticket-header">
                    <span class="ticket-id">${badge} <code>${t.id.substring(0, 8)}</code></span>
                    <span class="ticket-status status-${(t.status || 'PENDING').toUpperCase()}">${t.status}</span>
                </div>
                <div class="ticket-title">${t.title || 'Untitled'}</div>
                <div class="ticket-meta">Layer ${t.layer} · ${t.queue} · ${t.assigned_to || 'unassigned'} · ${new Date(t.created_at).toLocaleTimeString()}</div>`;
            item.addEventListener('click', () => openEditModal(t));
            ticketList.appendChild(item);
        });
    }

    function badgeFor(title) {
        if (!title) return '🎫';
        if (title.includes('ASSIGN') || title.includes('IDENTITY')) return '🆔';
        if (title.includes('BOOT')) return '🚀';
        if (title.includes('GENESIS')) return '🌌';
        if (title.includes('PROMOTE') || title.includes('GOLD')) return '⭐';
        return '🎫';
    }

    function statusColor(s) {
        switch ((s || '').toUpperCase()) {
            case 'NEW': return '#00f3ff';
            case 'OPEN': return '#00ff88';
            case 'IN_PROGRESS': return '#3498db';
            case 'STALLED': return '#f1c40f';
            case 'STUCK': return '#e74c3c';
            case 'RESOLVED': return '#2ecc71';
            case 'PROMOTED': return '#00ff88';
            case 'PENDING': return '#d4af37';
            case 'IMMUTABLE': return '#9b59b6';
            case 'DELETED': return '#555555';
            default: return 'rgba(255,255,255,0.4)';
        }
    }

    // ── GEDCOM Tree (Reingold-Tilford + SVG Pan/Zoom) ─────────
    const NW = 210, NH = 66, VGAP = 100, HGAP = 24;

    // SVG viewport pan/zoom state
    let treeScale = 1, treeTX = 0, treeTY = 0;
    let treeDragging = false, treeSX, treeSY;

    // Collapsed node IDs
    const collapsed = new Set();

    gedcomContainer.addEventListener('wheel', e => {
        e.preventDefault();
        const delta = e.deltaY > 0 ? 0.85 : 1.18;
        const newScale = Math.min(4, Math.max(0.1, treeScale * delta));
        // Zoom toward cursor
        const rect = gedcomContainer.getBoundingClientRect();
        const mx = e.clientX - rect.left, my = e.clientY - rect.top;
        treeTX = mx - (mx - treeTX) * (newScale / treeScale);
        treeTY = my - (my - treeTY) * (newScale / treeScale);
        treeScale = newScale;
        applyTreeTransform();
    }, { passive: false });

    gedcomContainer.addEventListener('mousedown', e => {
        // Only drag on container background, not on nodes
        if (e.target.closest('g')) return;
        treeDragging = true;
        treeSX = e.clientX - treeTX;
        treeSY = e.clientY - treeTY;
        gedcomContainer.style.cursor = 'grabbing';
    });
    window.addEventListener('mousemove', e => {
        if (!treeDragging) return;
        treeTX = e.clientX - treeSX;
        treeTY = e.clientY - treeSY;
        applyTreeTransform();
    });
    window.addEventListener('mouseup', () => {
        treeDragging = false;
        gedcomContainer.style.cursor = 'grab';
    });

    function applyTreeTransform() {
        treeViewport && treeViewport.setAttribute('transform',
            `translate(${treeTX},${treeTY}) scale(${treeScale})`);
    }

    let treeViewport = null;

    async function loadGedcomTree() {
        gedcomSvg.innerHTML = '<text x="50%" y="50%" text-anchor="middle" fill="#d4af37" font-family="JetBrains Mono" font-size="12">Loading genealogy…</text>';
        treeViewport = null;
        treeScale = 0.7; treeTX = 20; treeTY = 20;
        try {
            const res = await fetch('/REST/2.0/tickets/tree');
            if (!res.ok) throw new Error(res.statusText);
            const roots = await res.json();
            const total = countNodes(roots);
            renderGedcom(roots);
            addLog(`[REST 2.0] Genealogy: ${total} nodes, ${roots.length} root(s).`);
        } catch (err) {
            gedcomSvg.innerHTML = `<text x="20" y="40" fill="#e74c3c" font-family="JetBrains Mono" font-size="11">Tree error: ${err.message}</text>`;
            addLog(`[ERROR] Genealogy: ${err.message}`);
        }
    }

    function countNodes(nodes, n = 0) {
        return (nodes || []).reduce((acc, node) => countNodes(node.children, acc + 1), n);
    }

    /* ── Reingold-Tilford layout ─────────────────────────── */
    function rtLayout(nodes) {
        // Assign depth and preliminary x positions
        let leafX = 0;

        function firstWalk(node, depth) {
            node._depth = depth;
            node._visible = !collapsed.has(node.id);
            const kids = node._visible ? (node.children || []) : [];

            if (kids.length === 0) {
                node._prelim = leafX;
                leafX += NW + HGAP;
            } else {
                kids.forEach(c => firstWalk(c, depth + 1));
                const leftmost = kids[0]._prelim;
                const rightmost = kids[kids.length - 1]._prelim;
                node._prelim = (leftmost + rightmost) / 2;
            }
        }

        function secondWalk(node, offset) {
            node._x = node._prelim + offset;
            node._y = node._depth * (NH + VGAP) + 30;
            const kids = node._visible ? (node.children || []) : [];
            kids.forEach(c => secondWalk(c, offset));
        }

        nodes.forEach(r => firstWalk(r, 0));
        // Offset each root tree so they don't overlap
        let xOffset = 0;
        nodes.forEach(r => {
            secondWalk(r, xOffset);
            xOffset = r._x + NW + HGAP * 3;
        });
    }

    function flattenTree(nodes, out = []) {
        (nodes || []).forEach(n => {
            out.push(n);
            if (!collapsed.has(n.id)) flattenTree(n.children, out);
        });
        return out;
    }

    function renderGedcom(roots) {
        gedcomSvg.innerHTML = '';
        gedcomSvg.setAttribute('width', '100%');
        gedcomSvg.setAttribute('height', '100%');
        gedcomSvg.style.width = '100%';
        gedcomSvg.style.height = '100%';

        // Create a viewport group for pan/zoom
        treeViewport = document.createElementNS('http://www.w3.org/2000/svg', 'g');
        gedcomSvg.appendChild(treeViewport);
        applyTreeTransform();

        // Re-render function (called on collapse/expand)
        function rerender() {
            leafX_reset();
            rtLayout(roots);
            const all = flattenTree(roots);

            treeViewport.innerHTML = '';

            // Edges
            all.forEach(node => {
                if (collapsed.has(node.id)) return;
                (node.children || []).forEach(child => {
                    const x1 = node._x + NW / 2, y1 = node._y + NH;
                    const x2 = child._x + NW / 2, y2 = child._y;
                    const my = (y1 + y2) / 2;

                    const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
                    path.setAttribute('d', `M${x1},${y1} C${x1},${my} ${x2},${my} ${x2},${y2}`);
                    path.setAttribute('stroke', relColor(child.rel_type));
                    path.setAttribute('stroke-width', '1.5');
                    path.setAttribute('fill', 'none');
                    path.setAttribute('opacity', '0.55');
                    if (child.rel_type === 'CONTEXT') path.setAttribute('stroke-dasharray', '4,3');
                    treeViewport.appendChild(path);

                    // Inline relation badge on edge
                    if (child.rel_type) {
                        const badge = document.createElementNS('http://www.w3.org/2000/svg', 'text');
                        badge.setAttribute('x', (x1 + x2) / 2 + 4);
                        badge.setAttribute('y', my - 4);
                        badge.setAttribute('fill', relColor(child.rel_type));
                        badge.setAttribute('font-size', '7');
                        badge.setAttribute('font-family', 'JetBrains Mono');
                        badge.setAttribute('opacity', '0.6');
                        badge.textContent = child.rel_type;
                        treeViewport.appendChild(badge);
                    }
                });
            });

            // Nodes
            all.forEach(node => {
                const isGenesis = node.id === '00000000-0000-0000-0000-000000000000';
                const hasKids = (node.children || []).length > 0;
                const isCollapsed = collapsed.has(node.id);
                const sc = statusColor(node.status);

                const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
                g.setAttribute('transform', `translate(${node._x},${node._y})`);
                g.style.cursor = 'pointer';

                // Shadow glow
                const glow = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                glow.setAttribute('width', NW + 4); glow.setAttribute('height', NH + 4);
                glow.setAttribute('x', '-2'); glow.setAttribute('y', '-2');
                glow.setAttribute('rx', '8');
                glow.setAttribute('fill', 'none');
                glow.setAttribute('stroke', isGenesis ? '#d4af37' : sc);
                glow.setAttribute('stroke-width', '2');
                glow.setAttribute('opacity', '0.18');
                g.appendChild(glow);

                // Main card
                const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                rect.setAttribute('width', NW); rect.setAttribute('height', NH);
                rect.setAttribute('rx', '6');
                rect.setAttribute('fill', isGenesis ? 'rgba(212,175,55,0.12)' : 'rgba(14,14,26,0.9)');
                rect.setAttribute('stroke', isGenesis ? '#d4af37' : 'rgba(255,255,255,0.08)');
                rect.setAttribute('stroke-width', '1');
                g.appendChild(rect);

                // Status stripe (left edge)
                const stripe = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                stripe.setAttribute('width', '3'); stripe.setAttribute('height', NH);
                stripe.setAttribute('rx', '3'); stripe.setAttribute('fill', sc);
                g.appendChild(stripe);

                // Layer badge (top-right)
                const lbg = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                lbg.setAttribute('x', NW - 26); lbg.setAttribute('y', '1');
                lbg.setAttribute('width', '25'); lbg.setAttribute('height', '15');
                lbg.setAttribute('rx', '4');
                lbg.setAttribute('fill', 'rgba(0,243,255,0.15)');
                g.appendChild(lbg);
                const ltxt = svgText(treeViewport, `L${node.layer}`, NW - 13, 11, '#00f3ff', 8, 'JetBrains Mono', 'middle');
                g.appendChild(ltxt);

                // Status dot (top-right, above layer badge)
                const dot = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
                dot.setAttribute('cx', NW - 33); dot.setAttribute('cy', '9');
                dot.setAttribute('r', '4');
                dot.setAttribute('fill', sc);
                g.appendChild(dot);

                // Title
                g.appendChild(svgText(g, truncate(node.title, 26), 10, 20, isGenesis ? '#d4af37' : '#fff', 10.5, 'Outfit', 'start', '600'));

                // UUID row
                g.appendChild(svgText(g, node.id.substring(0, 20) + '…', 10, 34, 'rgba(255,255,255,0.3)', 8, 'JetBrains Mono'));

                // Creator row
                g.appendChild(svgText(g, `✦ ${node.creator}`, 10, 48, '#d4af37', 8, 'JetBrains Mono', 'start', '400', '0.65'));

                // Status text (bottom-right)
                g.appendChild(svgText(g, node.status, NW - 8, 56, sc, 7, 'JetBrains Mono', 'end', '400', '0.7'));

                // Collapse/expand toggle (▶ / ▼) if has children
                if (hasKids) {
                    const toggle = document.createElementNS('http://www.w3.org/2000/svg', 'text');
                    toggle.setAttribute('x', NW / 2);
                    toggle.setAttribute('y', NH + 14);
                    toggle.setAttribute('text-anchor', 'middle');
                    toggle.setAttribute('fill', '#00f3ff');
                    toggle.setAttribute('font-size', '10');
                    toggle.setAttribute('opacity', '0.7');
                    toggle.setAttribute('cursor', 'pointer');
                    toggle.textContent = isCollapsed
                        ? `▶ ${node.children.length} children`
                        : '▼';
                    toggle.addEventListener('click', ev => {
                        ev.stopPropagation();
                        isCollapsed ? collapsed.delete(node.id) : collapsed.add(node.id);
                        rerender();
                    });
                    g.appendChild(toggle);
                }

                // Click node → edit modal
                g.addEventListener('click', ev => {
                    if (ev.target.textContent && ev.target.textContent.includes('children')) return;
                    openEditModal({ id: node.id, title: node.title, status: node.status, creator: node.creator, intent: {} });
                });

                // Hover highlight
                g.addEventListener('mouseenter', () => rect.setAttribute('fill', isGenesis ? 'rgba(212,175,55,0.2)' : 'rgba(255,255,255,0.06)'));
                g.addEventListener('mouseleave', () => rect.setAttribute('fill', isGenesis ? 'rgba(212,175,55,0.12)' : 'rgba(14,14,26,0.9)'));

                treeViewport.appendChild(g);
            });
        }

        rerender();
    }

    // Helper: reset leafX counter before each layout pass
    let leafX_fn = null;
    function leafX_reset() {
        // leafX is scoped inside rtLayout, so we just call it fresh each time
    }

    function svgText(parent, content, x, y, fill, fontSize = 10, fontFamily = 'Outfit', anchor = 'start', fontWeight = '400', opacity = '1') {
        const t = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        t.setAttribute('x', x); t.setAttribute('y', y);
        t.setAttribute('fill', fill);
        t.setAttribute('font-size', fontSize);
        t.setAttribute('font-family', fontFamily);
        t.setAttribute('text-anchor', anchor);
        t.setAttribute('font-weight', fontWeight);
        t.setAttribute('opacity', opacity);
        t.textContent = content;
        return t;
    }

    function relColor(type) {
        switch ((type || '').toUpperCase()) {
            case 'EVOLUTION': return '#d4af37';
            case 'CONSEQUENCE': return '#00f3ff';
            case 'CONTEXT': return 'rgba(255,255,255,0.3)';
            case 'GENESIS': return '#9b59b6';
            default: return 'rgba(255,255,255,0.2)';
        }
    }

    function truncate(str, n) {
        return (str && str.length > n) ? str.substring(0, n) + '…' : (str || '');
    }

    refreshTreeBtn.addEventListener('click', loadGedcomTree);


    // ── Edit Modal ────────────────────────────────────────────
    const editModal = document.getElementById('editModal');
    const closeEdit = document.getElementById('closeEdit');
    const editForm = document.getElementById('editForm');

    function openEditModal(ticket) {
        document.getElementById('editTicketId').value = ticket.id;
        document.getElementById('editModalId').textContent = ticket.id.substring(0, 8) + '…';
        document.getElementById('editTitle').value = ticket.title || '';
        document.getElementById('editStatus').value = (ticket.status || 'PENDING').toUpperCase();
        document.getElementById('editPriority').value = ticket.priority || 20;
        document.getElementById('editAssignment').value = ticket.assigned_to || '';
        document.getElementById('editContent').value = JSON.stringify(ticket.intent || {}, null, 2);
        editModal.style.display = 'block';
    }

    closeEdit.onclick = () => editModal.style.display = 'none';
    window.addEventListener('click', e => { if (e.target === editModal) editModal.style.display = 'none'; });

    editForm.onsubmit = async e => {
        e.preventDefault();
        const id = document.getElementById('editTicketId').value;
        const payload = {
            Title: document.getElementById('editTitle').value,
            Status: document.getElementById('editStatus').value,
            Priority: parseInt(document.getElementById('editPriority').value, 10),
            AssignedTo: document.getElementById('editAssignment').value,
        };
        try {
            const res = await fetch(`/REST/2.0/ticket/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            if (res.ok) {
                addLog(`[REST 2.0] PUT /ticket/${id.substring(0, 8)}… → 200 OK`);
                editModal.style.display = 'none';
            } else {
                const txt = await res.text();
                addLog(`[REST 2.0] PUT failed: ${res.status} ${txt}`);
            }
        } catch (err) {
            addLog(`[ERROR] ${err.message}`);
        }
    };

    // ── Create Ticket Form ────────────────────────────────────
    const createForm = document.getElementById('createForm');
    const createResult = document.getElementById('createResult');

    createForm.onsubmit = async e => {
        e.preventDefault();
        const payload = {
            Subject: document.getElementById('createSubject').value,
            Queue: document.getElementById('createQueue').value,
            AgentID: document.getElementById('createAgent').value,
            Layer: parseInt(document.getElementById('createLayer').value, 10),
            Text: document.getElementById('createContent').value,
        };
        createResult.textContent = '▶ Submitting…';
        createResult.className = 'create-result';
        try {
            const res = await fetch('/REST/2.0/ticket', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const json = await res.json();
            if (res.ok) {
                createResult.textContent = `✓ Created: ${json.id}`;
                createResult.classList.add('success');
                addLog(`[REST 2.0] POST /ticket → ${json.id.substring(0, 8)}… created`);
                createForm.reset();
            } else {
                createResult.textContent = `✗ Error: ${json.error}`;
                createResult.classList.add('error');
            }
        } catch (err) {
            createResult.textContent = `✗ ${err.message}`;
            createResult.classList.add('error');
        }
    };

    // ── Log Helpers ───────────────────────────────────────────
    function addLog(msg) {
        const e = document.createElement('div');
        e.className = 'log-entry';
        e.textContent = `[${new Date().toLocaleTimeString()}] ${msg}`;
        logStream.prepend(e);
        if (logStream.children.length > 80) logStream.lastChild.remove();
    }

    clearLogsBtn.addEventListener('click', () => {
        logStream.innerHTML = '<div class="log-entry">[SYSTEM] Logs cleared.</div>';
    });
});
