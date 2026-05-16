import * as vscode from 'vscode';
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import axios from 'axios';
import { Octokit } from "@octokit/rest";
import * as fs from 'fs';
import * as path from 'path';

const PQR_API_URL = process.env.PQR_API_URL || 'https://pqr.info/REST/2.0';
const VAULT_ADDR = process.env.VAULT_ADDR || 'http://localhost:8200';
const VAULT_TOKEN = process.env.VAULT_TOKEN || 'pqr-vault-token';
const octokit = new Octokit({ auth: process.env.GITHUB_TOKEN });

let tokenStatusBarItem: vscode.StatusBarItem;

export function activate(context: vscode.ExtensionContext) {
    console.log('PQR Info Swarm is now active');

    // 1. Initialize MCP Server
    const server = new Server(
        {
            name: "pqr-info-swarm",
            version: "0.1.0",
        },
        {
            capabilities: {
                tools: {},
            },
        }
    );

    // 2. Define MCP Tools
    server.setRequestHandler(ListToolsRequestSchema, async () => ({
        tools: [
            {
                name: "create_pqr_ticket",
                description: "Create a new forensic ticket in the PQR Info Swarm fabric",
                inputSchema: {
                    type: "object",
                    properties: {
                        subject: { type: "string" },
                        content: { type: "string" },
                        agentID: { type: "string", default: "antigravity" },
                        intent: { type: "object" }
                    },
                    required: ["subject", "content"]
                },
            },
            {
                name: "query_pqr_memory",
                description: "Search the agent memory and fabric for relevant context",
                inputSchema: {
                    type: "object",
                    properties: {
                        agentID: { type: "string", default: "antigravity" },
                        limit: { type: "number", default: 10 }
                    }
                },
            },
            {
                name: "link_pqr_tickets",
                description: "Establish a forensic lineage link between two tickets",
                inputSchema: {
                    type: "object",
                    properties: {
                        parentID: { type: "string" },
                        childID: { type: "string" },
                        relType: { type: "string", enum: ["EVOLUTION", "CONSEQUENCE", "CONTEXT", "GENESIS"] }
                    },
                    required: ["parentID", "childID", "relType"]
                },
            },
            {
                name: "sync_github_issue",
                description: "Synchronize a PQR ticket with a GitHub Issue's state",
                inputSchema: {
                    type: "object",
                    properties: {
                        ticketID: { type: "string" },
                        owner: { type: "string" },
                        repo: { type: "string" },
                        issueNumber: { type: "number" }
                    },
                    required: ["ticketID", "owner", "repo", "issueNumber"]
                },
            },
            {
                name: "create_healing_ticket",
                description: "Create a new self-healing ticket from a log issue",
                inputSchema: {
                    type: "object",
                    properties: {
                        issue: { type: "string" },
                        logSnippet: { type: "string" }
                    },
                    required: ["issue", "logSnippet"]
                },
            },
            {
                name: "process_healing_iteration",
                description: "Advance a self-healing ticket to the next escalation iteration",
                inputSchema: {
                    type: "object",
                    properties: {
                        ticketID: { type: "string" }
                    },
                    required: ["ticketID"]
                },
            },
            {
                name: "record_healing_failure",
                description: "Log a failed attempt for a self-healing ticket to avoid repetition in higher tiers",
                inputSchema: {
                    type: "object",
                    properties: {
                        ticketID: { type: "string" },
                        failure: { type: "string" }
                    },
                    required: ["ticketID", "failure"]
                },
            },
            {
                name: "resolve_healing_ticket",
                description: "Finalize a self-healing ticket and add the resolution to the evolutionary knowledge base",
                inputSchema: {
                    type: "object",
                    properties: {
                        ticketID: { type: "string" },
                        resolution: { type: "string" },
                        agentID: { type: "string", default: "antigravity" }
                    },
                    required: ["ticketID", "resolution"]
                },
            },
            {
                name: "sovereign_override",
                description: "Execute a high-level emergency command via the Gemini Bridge",
                inputSchema: {
                    type: "object",
                    properties: {
                        command: { type: "string" },
                        params: { type: "object" }
                    },
                    required: ["command"]
                },
            },
            {
                name: "get_forensic_audit",
                description: "Retrieve a deep forensic audit of the last active swarm tickets",
                inputSchema: {
                    type: "object",
                    properties: {
                        limit: { type: "number", default: 10 }
                    }
                },
            }
        ],
    }));

    // 3. Handle Tool Calls
    server.setRequestHandler(CallToolRequestSchema, async (request) => {
        const { name, arguments: args } = request.params;

        try {
            switch (name) {
                case "create_pqr_ticket":
                    const createResp = await axios.post(`${PQR_API_URL}/ticket`, args);
                    return { content: [{ type: "text", text: JSON.stringify(createResp.data, null, 2) }] };

                case "query_pqr_memory":
                    const queryResp = await axios.get(`${PQR_API_URL}/agent/${args?.agentID || 'antigravity'}/context?limit=${args?.limit || 10}`);
                    return { content: [{ type: "text", text: JSON.stringify(queryResp.data, null, 2) }] };

                case "link_pqr_tickets":
                    const linkResp = await axios.post(`${PQR_API_URL}/ticket/${args?.parentID}/link/${args?.childID}`, {
                        relationship_type: args?.relType,
                        agent_id: "antigravity"
                    });
                    return { content: [{ type: "text", text: JSON.stringify(linkResp.data, null, 2) }] };

                case "sync_github_issue":
                    const { ticketID, owner, repo, issueNumber } = args as any;
                    const issue = await octokit.issues.get({ owner, repo, issue_number: issueNumber });
                    
                    // Map GitHub state to PQR status
                    let pqrStatus = "ACTIVE";
                    if (issue.data.state === 'closed') {
                        pqrStatus = "COMPLETED";
                    }

                    // Update PQR ticket
                    const syncResp = await axios.put(`${PQR_API_URL}/ticket/${ticketID}`, {
                        Status: pqrStatus,
                        Title: issue.data.title
                    });

                    return { content: [{ type: "text", text: `Synced with GitHub #${issueNumber}: Status set to ${pqrStatus}` }] };

                case "create_healing_ticket":
                    const healingResp = await axios.post(`${PQR_API_URL}/healing/ticket`, request.params.arguments);
                    return { content: [{ type: "text", text: JSON.stringify(healingResp.data, null, 2) }] };

                case "process_healing_iteration":
                    const iterResp = await axios.post(`${PQR_API_URL}/healing/iterate/${request.params.arguments?.ticketID}`);
                    return { content: [{ type: "text", text: JSON.stringify(iterResp.data, null, 2) }] };

                case "record_healing_failure":
                    const failResp = await axios.post(`${PQR_API_URL}/healing/failure`, request.params.arguments);
                    return { content: [{ type: "text", text: JSON.stringify(failResp.data, null, 2) }] };

                case "resolve_healing_ticket":
                    const resResp = await axios.post(`${PQR_API_URL}/healing/resolve`, args);
                    return { content: [{ type: "text", text: JSON.stringify(resResp.data, null, 2) }] };
                
                case "sovereign_override":
                    // Pull Gemini Key from environment or vault
                    const geminiKey = process.env.GEMINI_API_KEY || "AIzaSyCqMMdPm1s6MuXy06yiWUlIQ0CJ1C-rPWk";
                    const overrideResp = await axios.post(`${PQR_API_URL}/emergency/bridge`, args, {
                        headers: { "X-Gemini-Key": geminiKey }
                    });
                    return { content: [{ type: "text", text: JSON.stringify(overrideResp.data, null, 2) }] };

                case "get_forensic_audit":
                    const auditResp = await axios.get(`${PQR_API_URL}/tickets?limit=${args?.limit || 10}`);
                    return { content: [{ type: "text", text: JSON.stringify(auditResp.data, null, 2) }] };

                default:
                    throw new Error(`Unknown tool: ${name}`);
            }
        } catch (error: any) {
            return {
                content: [{ type: "text", text: `Error: ${error.message}` }],
                isError: true,
            };
        }
    });

    // 4. Start Transport
    const transport = new StdioServerTransport();
    server.connect(transport).catch(console.error);

    context.subscriptions.push(disposable);
    
    // 5b. Emergency Repair Override
    let emergencyRepair = vscode.commands.registerCommand('pqr.emergencyRepair', async () => {
        const directive = await vscode.window.showInputBox({
            prompt: 'Enter High Justiciar Repair Directive',
            placeHolder: 'e.g. FORCE_HEALING_ITERATION, STOP_TOKEN_BURN'
        });

        if (!directive) return;

        try {
            const geminiKey = process.env.GEMINI_API_KEY || "AIzaSyCqMMdPm1s6MuXy06yiWUlIQ0CJ1C-rPWk";
            const response = await axios.post(`${PQR_API_URL}/emergency/bridge`, {
                command: "TRIGGER_HEALING",
                params: { issue: directive }
            }, {
                headers: { "X-Gemini-Key": geminiKey }
            });
            vscode.window.showInformationMessage(`✓ Repair Directive Executed: ${JSON.stringify(response.data)}`);
        } catch (error: any) {
            vscode.window.showErrorMessage(`Repair Failed: ${error.message}`);
        }
    });

    context.subscriptions.push(emergencyRepair);

    // 6. Setup Vault (Sweep .env)
    let vaultSetup = vscode.commands.registerCommand('pqr.setupVault', async () => {
        const fileUri = await vscode.window.showOpenDialog({
            canSelectFiles: true,
            canSelectFolders: false,
            canSelectMany: false,
            filters: { 'Env Files': ['env'] },
            title: 'Select .env file to sweep into PQR Vault'
        });

        if (!fileUri || fileUri.length === 0) return;

        const envPath = fileUri[0].fsPath;
        const auth = await vscode.window.showInformationMessage(
            `Authorize sweeping secrets from ${path.basename(envPath)} into PQR Vault? This will securely migrate your credentials.`,
            { modal: true },
            'Authorize & Sweep'
        );

        if (auth === 'Authorize & Sweep') {
            try {
                const content = fs.readFileSync(envPath, 'utf8');
                const secrets: Record<string, string> = {};
                
                content.split('\n').forEach(line => {
                    const match = line.match(/^([^#\s][^=]*)=(.*)$/);
                    if (match) {
                        secrets[match[1].trim()] = match[2].trim();
                    }
                });

                if (Object.keys(secrets).length === 0) {
                    vscode.window.showWarningMessage('No secrets found in the selected file.');
                    return;
                }

                await axios.post(`${VAULT_ADDR}/v1/secret/data/pqr`, { data: secrets }, {
                    headers: { 'X-Vault-Token': VAULT_TOKEN }
                });

                vscode.window.showInformationMessage('✓ Secrets successfully swept into PQR Vault. You can now safely remove the .env file.');
            } catch (error: any) {
                vscode.window.showErrorMessage(`Vault Sweep Failed: ${error.message}`);
            }
        }
    });

    context.subscriptions.push(vaultSetup);

    // 7. Initialize Token Sentinel (Status Bar)
    tokenStatusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    tokenStatusBarItem.command = 'pqr.openHUD';
    context.subscriptions.push(tokenStatusBarItem);
    tokenStatusBarItem.show();

    // Start polling metrics
    updateTokenSentinel();
    setInterval(updateTokenSentinel, 30000); // Poll every 30 seconds
}

async function updateTokenSentinel() {
    try {
        const response = await axios.get(`${PQR_API_URL}/metrics/tokens`);
        const { usage_percentage } = response.data;
        const percent = usage_percentage.toFixed(1);

        tokenStatusBarItem.text = `$(circuit-board) PQR Tokens: ${percent}%`;
        tokenStatusBarItem.tooltip = `Swarm Token Quota Usage: ${percent}%`;
        
        // Color coding for sentinel
        if (usage_percentage > 90) {
            tokenStatusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
        } else if (usage_percentage > 75) {
            tokenStatusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
        } else {
            tokenStatusBarItem.backgroundColor = undefined;
            tokenStatusBarItem.color = '#00f2ff'; // Sovereign Blue
        }
    } catch (error) {
        tokenStatusBarItem.text = `$(warning) PQR Tokens: Offline`;
        tokenStatusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
    }
}

export function deactivate() {}
