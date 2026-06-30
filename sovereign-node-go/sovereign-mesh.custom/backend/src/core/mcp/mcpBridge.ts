import { loadConfig, loadJsonConfig } from '../../config';

type MCPConfig = {
  endpoint: string;
  tools: { name: string }[];
};

export class MCPBridge {
  private cfg: MCPConfig | null = null;

  constructor() {
    const appCfg = loadConfig();
    try {
      this.cfg = loadJsonConfig<MCPConfig>(appCfg.MCP_CONFIG_PATH);
    } catch {
      this.cfg = {
        endpoint: 'ssh://localhost',
        tools: [
          { name: 'fs.readFile' },
          { name: 'fs.writeFile' },
          { name: 'sshfs.mount' },
          { name: 'sshfs.unmount' }
        ]
      };
    }
  }

  async callTool(name: string, args: any): Promise<any> {
    console.log(`[MCP Bridge] Calling tool ${name} with args:`, args);
    return { ok: true, name, args, fake: true };
  }
}
