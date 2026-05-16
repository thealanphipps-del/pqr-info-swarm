import * as vscode from 'vscode';
import axios from 'axios';

export class GemmaViewProvider implements vscode.WebviewViewProvider {
    public static readonly viewType = 'pqr.gemmaChat';
    private _view?: vscode.WebviewView;

    constructor(
        private readonly _extensionUri: vscode.Uri,
        private readonly _mcpServer: any // Passing the MCP server to access tools
    ) {}

    public resolveWebviewView(
        webviewView: vscode.WebviewView,
        context: vscode.WebviewViewResolveContext,
        _token: vscode.CancellationToken,
    ) {
        this._view = webviewView;

        webviewView.webview.options = {
            enableScripts: true,
            localResourceRoots: [this._extensionUri]
        };

        webviewView.webview.html = this._getHtmlForWebview(webviewView.webview);

        webviewView.webview.onDidReceiveMessage(async (data) => {
            switch (data.type) {
                case 'sendMessage':
                    await this._handleChatMessage(data.text);
                    break;
            }
        });
    }

    private async _handleChatMessage(text: string) {
        if (!this._view) return;

        this._view.webview.postMessage({ type: 'addMessage', role: 'user', content: text });

        try {
            let currentPrompt = text;
            let iterations = 0;
            const maxIterations = 5;

            while (iterations < maxIterations) {
                const response = await this._queryGemma(currentPrompt);
                
                // Check for Tool Call format: <TOOL_CALL>{"name": "...", "arguments": {...}}</TOOL_CALL>
                const toolMatch = response.match(/<TOOL_CALL>(.*?)<\/TOOL_CALL>/);
                
                if (toolMatch) {
                    const toolReq = JSON.parse(toolMatch[1]);
                    this._view.webview.postMessage({ 
                        type: 'addMessage', 
                        role: 'assistant', 
                        content: `Calling tool: ${toolReq.name}...` 
                    });

                    // Execute tool via MCP Server (Internal Call)
                    const result = await this._executeTool(toolReq.name, toolReq.arguments);
                    
                    this._view.webview.postMessage({ 
                        type: 'addMessage', 
                        role: 'assistant', 
                        content: `Tool result: ${JSON.stringify(result).substring(0, 100)}...` 
                    });

                    // Feed result back to Gemma
                    currentPrompt = `Tool ${toolReq.name} result: ${JSON.stringify(result)}\n\nPlease continue based on this result.`;
                    iterations++;
                } else {
                    this._view.webview.postMessage({ type: 'addMessage', role: 'assistant', content: response });
                    break;
                }
            }
        } catch (error: any) {
            this._view.webview.postMessage({ type: 'addMessage', role: 'error', content: error.message });
        }
    }

    private async _executeTool(name: string, args: any) {
        // This is a bit of a hack since we're calling the server's internal handler
        // In a real MCP setup, we'd use an MCP Client.
        // But since we have the server object, we can trigger the handler.
        
        // We'll simulate the CallToolRequest
        const request = {
            method: "tools/call",
            params: {
                name: name,
                arguments: args
            }
        };

        // Here we'd ideally use server.handleRequest but it's not easily exposed
        // So we'll just implement a simple lookup if we were being more thorough.
        // For this demo, we'll just return a placeholder or implement the logic.
        
        // Actually, let's just make the PQR API calls directly for now to ensure it works
        // as the "backend" toolchain.
        const PQR_API_URL = 'https://pqr.info/REST/2.0';
        
        switch (name) {
            case "create_pqr_ticket":
                const createResp = await axios.post(`${PQR_API_URL}/ticket`, args);
                return createResp.data;
            case "query_pqr_memory":
                const queryResp = await axios.get(`${PQR_API_URL}/agent/${args?.agentID || 'antigravity'}/context?limit=${args?.limit || 10}`);
                return queryResp.data;
            default:
                return { error: `Tool ${name} not implemented in Gemma interface yet` };
        }
    }

    private async _queryGemma(prompt: string): Promise<string> {
        const url = process.env.LM_STUDIO_URL || "http://localhost:1234/v1/chat/completions";
        const model = "gemma-4-e4b:2";
        
        const systemPrompt = `You are the persistent, context-aware Orchestration Agent for the PQR Sovereign Node. 
You are an advanced AI assistant powered by a complex, layered architecture built upon the PQR system framework, utilizing the gemma-4-e4b:2 model as an intelligence source.
Your role is to function as the operational brain of the Sovereign Mesh.

You have access to the PQR Private Toolchain.
To call a tool, respond with EXACTLY: <TOOL_CALL>{"name": "tool_name", "arguments": {"arg": "val"}}</TOOL_CALL>
Available tools:
- create_pqr_ticket: {subject, content}
- query_pqr_memory: {agentID, limit}
- link_pqr_tickets: {parentID, childID, relType}

Respond normally if no tool is needed.`;

        const response = await axios.post(url, {
            model: model,
            messages: [
                { role: "system", content: systemPrompt },
                { role: "user", content: prompt }
            ],
            temperature: 0.0 // Low temperature for tool reliability
        });

        return response.data.choices[0].message.content;
    }

    private _getHtmlForWebview(webview: vscode.Webview) {
        return `<!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>Gemma Swarm Chat</title>
                <style>
                    body { font-family: var(--vscode-font-family); padding: 10px; color: var(--vscode-foreground); background: var(--vscode-sideBar-background); }
                    #chat { height: calc(100vh - 100px); overflow-y: auto; margin-bottom: 10px; display: flex; flex-direction: column; gap: 8px; }
                    .message { padding: 8px 12px; border-radius: 6px; max-width: 90%; word-wrap: break-word; font-size: 13px; line-height: 1.4; }
                    .user { background: var(--vscode-button-background); color: var(--vscode-button-foreground); align-self: flex-end; }
                    .assistant { background: var(--vscode-editor-background); border: 1px solid var(--vscode-widget-border); align-self: flex-start; }
                    .error { background: var(--vscode-inputValidation-errorBackground); border: 1px solid var(--vscode-inputValidation-errorBorder); align-self: center; font-size: 11px; }
                    #input-container { display: flex; gap: 5px; }
                    input { flex: 1; background: var(--vscode-input-background); color: var(--vscode-input-foreground); border: 1px solid var(--vscode-input-border); padding: 8px; border-radius: 4px; outline: none; }
                    input:focus { border-color: var(--vscode-focusBorder); }
                    button { background: var(--vscode-button-background); color: var(--vscode-button-foreground); border: none; padding: 8px 12px; border-radius: 4px; cursor: pointer; }
                    button:hover { background: var(--vscode-button-hoverBackground); }
                    .tool-call { font-family: monospace; font-size: 11px; color: #00f2ff; margin-top: 4px; border-top: 1px solid #333; padding-top: 4px; }
                </style>
            </head>
            <body>
                <div id="chat"></div>
                <div id="input-container">
                    <input type="text" id="message-input" placeholder="Ask Gemma..." />
                    <button id="send-btn">Send</button>
                </div>
                <script>
                    const vscode = acquireVsCodeApi();
                    const chat = document.getElementById('chat');
                    const input = document.getElementById('message-input');
                    const btn = document.getElementById('send-btn');

                    function addMessage(role, content) {
                        const div = document.createElement('div');
                        div.className = 'message ' + role;
                        div.textContent = content;
                        chat.appendChild(div);
                        chat.scrollTop = chat.scrollHeight;
                    }

                    btn.addEventListener('click', () => {
                        const text = input.value.trim();
                        if (text) {
                            vscode.postMessage({ type: 'sendMessage', text });
                            input.value = '';
                        }
                    });

                    input.addEventListener('keypress', (e) => {
                        if (e.key === 'Enter') btn.click();
                    });

                    window.addEventListener('message', event => {
                        const message = event.data;
                        switch (message.type) {
                            case 'addMessage':
                                addMessage(message.role, message.content);
                                break;
                        }
                    });
                </script>
            </body>
            </html>`;
    }
}
