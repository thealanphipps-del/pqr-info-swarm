import * as fs from 'fs';
import * as path from 'path';
import * as dotenv from 'dotenv';

dotenv.config();

export type AppConfig = {
  PORT: number;
  DATABASE_URL: string;
  VERTEX_PROJECT_ID: string;
  VERTEX_LOCATION: string;
  MODELS_CONFIG_PATH: string;
  MCP_CONFIG_PATH: string;
};

export function loadConfig(): AppConfig {
  return {
    PORT: Number(process.env.PORT ?? 4000),
    DATABASE_URL: process.env.DATABASE_URL ?? 'postgres://localhost/pqr',
    VERTEX_PROJECT_ID: process.env.VERTEX_PROJECT_ID ?? '',
    VERTEX_LOCATION: process.env.VERTEX_LOCATION ?? 'us-central1',
    MODELS_CONFIG_PATH:
      process.env.MODELS_CONFIG_PATH ??
      path.join(__dirname, '../config/models.json'),
    MCP_CONFIG_PATH:
      process.env.MCP_CONFIG_PATH ??
      path.join(__dirname, '../config/mcp.json'),
  };
}

export function loadJsonConfig<T = any>(filePath: string): T {
  const raw = fs.readFileSync(filePath, 'utf8');
  return JSON.parse(raw) as T;
}
