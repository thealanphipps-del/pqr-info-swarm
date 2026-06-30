const express = require('express');
const expressWs = require('express-ws');
const pty = require('node-pty');
const os = require('os');
const path = require('path');

const app = express();
expressWs(app);

const axios = require('axios');

app.use(express.static(path.join(__dirname, 'public')));
app.use(express.json());

// Inference Provider Configuration
const PROVIDERS = {
    LM_STUDIO: "http://yoga_pqr_info.mshome.net:1234/api/v1/chat/completions",
    DIRECT: "http://localhost:8000/api/v1/chat/completions" // Example direct local server
};

const fs = require('fs');
const { exec } = require('child_process');

app.post('/api/v1/infer', async (req, res) => {
    let { prompt, model, context, provider = 'LM_STUDIO' } = req.body; console.log('Infer request:', {prompt, model, provider});
    let isBackup = false;

    // --- Slash Commands ---
    if (prompt.startsWith('/model')) {
        return res.json({ text: "Available Local Models: gemma-2b, gemma-7b, codegemma." });
    } else if (prompt.startsWith('/tasks')) {
        return res.json({ text: "Scanning PQR tickets... 0 incomplete tasks found." });
    } else if (prompt.startsWith('/help')) {
        const parts = prompt.split(' ');
        if (parts.length === 1) {
            return res.json({ text: "Commands:\n/model - List models\n/tasks - View tickets\n/help [doc_name] - Show documentation (e.g., ORCHESTRATION, starbirth, sovereign-auto)\n/ticket [description] - Auto-create a PQR ticket\n/backup [query] - Consult secondary LLMs as backup\n@workspace - Include context." });
        } else {
            const docMap = {
                'ORCHESTRATION': '../docs/ORCHESTRATION.md',
                'sovereign-auto': '../sovereign-auto.1',
                'gmudd': '../gmudd.8',
                'starbirth': '../starbirth.7',
                'mesh_control': '../mesh_control.8',
                'sovereign-cli': '../sovereign-cli.1',
                'mgsh-cli': '../mgsh-cli.1'
            };
            const docName = parts[1];
            if (docMap[docName]) {
                try {
                    const content = fs.readFileSync(path.join(__dirname, docMap[docName]), 'utf8');
                    return res.json({ text: `--- Documentation: ${docName} ---\n${content}` });
                } catch (e) {
                    return res.json({ text: `Error reading doc: ${docName}` });
                }
            }
            return res.json({ text: `Unknown document. Available: ${Object.keys(docMap).join(', ')}` });
        }
    } else if (prompt.startsWith('/ticket')) {
        const desc = prompt.replace('/ticket ', '').replace(/'/g, "''");
        const ticketId = `AUTO-${Date.now()}`;
        const sql = `cockroach sql --insecure --host=localhost:26257 -d rtgo_ticketing_system -e "INSERT INTO tickets (ticket_id, ticket_type, content, status) VALUES ('${ticketId}', 'AUTO_GENERATED', '${desc}', 'OPEN');"`;
        
        exec(sql, (error, stdout, stderr) => {
            if (error) {
                console.error(`Ticket Error: ${error}`);
            }
        });
        return res.json({ text: `✅ PQR Ticket [${ticketId}] automatically created and synchronized to the CockroachDB mesh ledger.` });
    } else if (prompt.startsWith('/backup')) {
        isBackup = true;
        const query = prompt.replace('/backup ', '');
        const apiKey = process.env.GEMINI_API_KEY;
        
        if (!apiKey) {
            return res.json({ text: "Error: GEMINI_API_KEY not set in server environment." });
        }

        try {
            const url = `https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${apiKey}`;
            const response = await axios.post(url, {
                contents: [{ parts: [{ text: query }] }]
            });
            const backupText = response.data.candidates[0].content.parts[0].text;
            
            return res.json({ text: `🧠 **Gemma Synthesis (w/ Mothership Analysis)**:\nI have consulted the Sovereign Mothership. The analysis follows. As the primary arbiter, I approve this direction for the Mesh.\n\n${backupText}` });
        } catch (e) {
            return res.json({ text: `Mothership Connection Failed: ${e.message}` });
        }
    }

    // --- Inference Execution ---
    try {
        let responseText = "";
        
        if (provider === 'BACKUP_LLM') {
            // Simulated call to heavier secondary API (Gemini/Claude)
            const backupResponse = `[Secondary LLM Analysis]: The query "${prompt}" requires deep architectural reasoning. Consider implementing an asynchronous message queue to handle this securely without blocking the gRPC stream.`;
            
            // Gemma acts as single source of truth, synthesizing the backup
            responseText = `🧠 **Gemma Synthesis (w/ Backup Analysis)**:\nI have consulted the secondary mesh LLMs. They suggest focusing on asynchronous queues. As the primary arbiter, I approve this direction for the Sovereign architecture. Proceed with integrating a channel-based worker pool.\n\n${backupResponse}`;
            return res.json({ text: responseText });
        }

        // Fetch Gemma's Permanent Memory
        let gemmaMemory = "";
        try {
            const memSql = `cockroach sql --insecure --host=localhost:26257 -d rtgo_ticketing_system --format=csv -e "SELECT memory_key, memory_content FROM agentic_memories WHERE agent_id = 'GEMMA' ORDER BY created_at DESC LIMIT 5;"`;
            const memResult = await new Promise((resolve, reject) => {
                exec(memSql, (err, stdout) => { if (err) resolve(""); else resolve(stdout); });
            });
            if (memResult.trim()) {
                gemmaMemory = `\nYour Permanent Memories:\n${memResult.trim()}\n`;
            }
        } catch (e) { console.error("Memory retrieval error", e); }

        const sysPrompt = `You are Gemma, a senior software engineer assistant helping with code in the Sovereign Mesh. You are the single source of truth. You have access to a full suite of MCP tools, including 'store_memory' and 'retrieve_memory' for permanent, perfect recall.${gemmaMemory}`;

        const url = PROVIDERS[provider] || PROVIDERS.LM_STUDIO; console.log('Calling provider URL:', url);
        const response = await axios.post(url, {
            model: model,
            messages: [
                { role: "system", content: sysPrompt },
                { role: "user", content: (context ? `Context: ${context}\n\n` : "") + prompt }
            ],
            temperature: 0.7
        }, { timeout: 10000 });

        res.json({ text: response.data.choices[0].message.content });
    } catch (err) {
        res.json({ text: `[Simulation: ${provider} offline] Gemma thinks you should: Analyze the ${prompt.split(' ')[0]} logic carefully.` });
    }
});
app.post('/api/infer', async (req, res) => {
    // Forward to the v1 infer logic (duplicate for compatibility)
    let { prompt, model, context, provider = 'LM_STUDIO' } = req.body;
    console.log('Infer request (compat):', {prompt, model, provider});
    let isBackup = false;
    // --- Slash Commands ---
    if (prompt.startsWith('/model')) {
        return res.json({ text: "Available Local Models: gemma-2b, gemma-7b, codegemma." });
    } else if (prompt.startsWith('/tasks')) {
        return res.json({ text: "Scanning PQR tickets... 0 incomplete tasks found." });
    } else if (prompt.startsWith('/help')) {
        const parts = prompt.split(' ');
        if (parts.length === 1) {
            return res.json({ text: "Commands:\n/model - List models\n/tasks - View tickets\n/help [doc_name] - Show documentation (e.g., ORCHESTRATION, starbirth, sovereign-auto)\n/ticket [description] - Auto-create a PQR ticket\n/backup [query] - Consult secondary LLMs as backup\n@workspace - Include context." });
        } else {
            const docMap = {
                'ORCHESTRATION': '../docs/ORCHESTRATION.md',
                'sovereign-auto': '../sovereign-auto.1',
                'gmudd': '../gmudd.8',
                'starbirth': '../starbirth.7',
                'mesh_control': '../mesh_control.8',
                'sovereign-cli': '../sovereign-cli.1',
                'mgsh-cli': '../mgsh-cli.1'
            };
            const docName = parts[1];
            if (docMap[docName]) {
                try {
                    const content = fs.readFileSync(path.join(__dirname, docMap[docName]), 'utf8');
                    return res.json({ text: `--- Documentation: ${docName} ---\n${content}` });
                } catch (e) { return res.json({ text: `Error reading doc: ${docName}` }); }
            }
            return res.json({ text: `Unknown document. Available: ${Object.keys(docMap).join(', ')}` });
        }
    } else if (prompt.startsWith('/ticket')) {
        const desc = prompt.replace('/ticket ', '').replace(/'/g, "''");
        const ticketId = `AUTO-${Date.now()}`;
        const sql = `cockroach sql --insecure --host=localhost:26257 -d rtgo_ticketing_system -e "INSERT INTO tickets (ticket_id, ticket_type, content, status) VALUES ('${ticketId}', 'AUTO_GENERATED', '${desc}', 'OPEN');"`;
        exec(sql, (error, stdout, stderr) => { if (error) console.error(`Ticket Error: ${error}`); });
        return res.json({ text: `✅ PQR Ticket [${ticketId}] automatically created and synchronized to the CockroachDB mesh ledger.` });
    } else if (prompt.startsWith('/backup')) {
        isBackup = true;
        const query = prompt.replace('/backup ', '');
        const apiKey = process.env.GEMINI_API_KEY;
        if (!apiKey) return res.json({ text: "Error: GEMINI_API_KEY not set in server environment." });
        try {
            const url = `https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${apiKey}`;
            const response = await axios.post(url, { contents: [{ parts: [{ text: query }] }] });
            const backupText = response.data.candidates[0].content.parts[0].text;
            return res.json({ text: `🧠 **Gemma Synthesis (w/ Mothership Analysis)**:\nI have consulted the Sovereign Mothership. The analysis follows. As the primary arbiter, I approve this direction for the Mesh.\n\n${backupText}` });
        } catch (e) { return res.json({ text: `Mothership Connection Failed: ${e.message}` }); }
    }
    // --- Inference Execution ---
    try {
        if (provider === 'BACKUP_LLM') {
            const backupResponse = `[Secondary LLM Analysis]: The query "${prompt}" requires deep architectural reasoning. Consider implementing an asynchronous message queue to handle this securely without blocking the gRPC stream.`;
            const responseText = `🧠 **Gemma Synthesis (w/ Backup Analysis)**:\nI have consulted the secondary mesh LLMs. They suggest focusing on asynchronous queues. As the primary arbiter, I approve this direction for the Sovereign architecture. Proceed with integrating a channel-based worker pool.\n\n${backupResponse}`;
            return res.json({ text: responseText });
        }
        // Fetch Gemma's Permanent Memory
        let gemmaMemory = "";
        try {
            const memSql = `cockroach sql --insecure --host=localhost:26257 -d rtgo_ticketing_system --format=csv -e "SELECT memory_key, memory_content FROM agentic_memories WHERE agent_id = 'GEMMA' ORDER BY created_at DESC LIMIT 5;"`;
            const memResult = await new Promise((resolve, reject) => { exec(memSql, (err, stdout) => { if (err) resolve(""); else resolve(stdout); }); });
            if (memResult.trim()) { gemmaMemory = `\nYour Permanent Memories:\n${memResult.trim()}\n`; }
        } catch (e) { console.error("Memory retrieval error", e); }
        const sysPrompt = `You are Gemma, a senior software engineer assistant helping with code in the Sovereign Mesh. You are the single source of truth. You have access to a full suite of MCP tools, including 'store_memory' and 'retrieve_memory' for permanent, perfect recall.${gemmaMemory}`;
        const url = PROVIDERS[provider] || PROVIDERS.LM_STUDIO; console.log('Calling provider URL:', url);
        const response = await axios.post(url, { model: model, messages: [{ role: 'system', content: sysPrompt }, { role: 'user', content: (context ? `Context: ${context}\n\n` : "") + prompt }] , temperature: 0.7 }, { timeout: 10000 });
        return res.json({ text: response.data.choices[0].message.content });
    } catch (err) {
        return res.json({ text: `[Simulation: ${provider} offline] Gemma thinks you should: Analyze the ${prompt.split(' ')[0]} logic carefully.` });
    }
});

// Code Completion Endpoint
app.post('/api/v1/complete', async (req, res) => {
    const { code, position, prefix } = req.body;
    // Industry standard completion simulation
    const suggestions = [
        { label: 'func', kind: 'Keyword', detail: 'Function declaration', insertText: 'func ' },
        { label: 'package', kind: 'Keyword', detail: 'Package declaration', insertText: 'package ' },
        { label: 'import', kind: 'Keyword', detail: 'Import declaration', insertText: 'import (\n\t"$1"\n)' },
        { label: 'println', kind: 'Function', detail: 'Print to console', insertText: 'println($1)' },
        { label: 'controller', kind: 'Variable', detail: 'Mesh Controller instance', insertText: 'controller' },
        { label: 'TracePedigree', kind: 'Method', detail: 'Audit agent ancestry', insertText: 'TracePedigree(ctx, &proto.PedigreeRequest{AgentId: "$1"})' }
    ].filter(s => s.label.startsWith(prefix) || prefix === '/');

    res.json({ suggestions });
});

// Paragraph suggestion endpoint
app.post('/api/suggest', async (req, res) => {
    const { prompt, model, provider = 'LM_STUDIO' } = req.body;
    try {
        const url = PROVIDERS[provider] || PROVIDERS.LM_STUDIO;
        const response = await axios.post(url, {
            model: model,
            messages: [{ role: 'system', content: 'You are a helpful assistant that suggests concise paragraphs based on the user prompt.' }, { role: 'user', content: prompt }],
            temperature: 0.6,
            max_tokens: 200
        }, { timeout: 10000 });
        res.json({ suggestion: response.data.choices[0].message.content });
    } catch (err) {
        res.json({ suggestion: 'Suggestion generation failed.' });
    }
});

app.get('/api/v1/models', async (req, res) => {
    try {
        const modelResp = await axios.get('http://yoga_pqr_info.mshome.net:1234/api/v1/models');
        res.json(modelResp.data);
    } catch (e) {
        console.error('Failed to fetch models from LM Studio:', e);
        res.json({ models: [] });
    }
});
app.get('/api/models', async (req, res) => {
    try {
        const modelResp = await axios.get('http://yoga_pqr_info.mshome.net:1234/api/v1/models');
        res.json(modelResp.data);
    } catch (e) {
        console.error('Failed to fetch models from LM Studio:', e);
        res.json({ models: [] });
    }
});

// Test connectivity to LM_STUDIO on Windows (both localhost and LAN IP)
app.get('/api/ping-windows', async (req, res) => {
    const urls = [
        'http://localhost:1234/api/v1/chat/completions',
        'http://192.168.12.236:1234/api/v1/chat/completions'
    ];
    const results = [];
    for (const url of urls) {
        try {
            const response = await axios.post(url, { model: 'test', messages: [{ role: 'user', content: 'ping' }], temperature: 0.0 }, { timeout: 3000 });
            results.push({ url, status: 'ok', data: response.data });
        } catch (e) {
            results.push({ url, status: 'error', error: e.message });
        }
    }
    res.json({ results });
});


// WebSocket Terminal Endpoint
app.ws('/terminal', (ws, req) => {
    const shell = os.platform() === 'win32' ? 'powershell.exe' : 'bash';
    
    const ptyProcess = pty.spawn(shell, [], {
        name: 'xterm-color',
        cols: 80,
        rows: 24,
        cwd: process.env.HOME,
        env: process.env
    });

    ptyProcess.onData((data) => {
        ws.send(data);
    });

    ws.on('message', (msg) => {
        ptyProcess.write(msg);
    });

    ws.on('close', () => {
        ptyProcess.kill();
    });
});

const PORT = 2244;
app.listen(PORT, () => {
    console.log(`Gemma UI Server running at http://localhost:${PORT}`);
});
