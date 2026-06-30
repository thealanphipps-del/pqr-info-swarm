import { Pool } from 'pg';
import { loadConfig } from '../../config';

export class MemoryRepo {
  private pool: Pool | null = null;
  private summaries: Map<string, string> = new Map();
  private facts: Map<string, string[]> = new Map();
  private snippets: Map<string, string[]> = new Map();

  constructor() {
    const config = loadConfig();
    try {
      this.pool = new Pool({
        connectionString: config.DATABASE_URL,
      });
    } catch {
      this.pool = null;
    }
  }

  async getSummary(threadId: string): Promise<string | null> {
    if (this.pool) {
      try {
        const res = await this.pool.query(
          'SELECT summary FROM memory_summaries WHERE thread_id = $1',
          [threadId]
        );
        return res.rows[0]?.summary || null;
      } catch (err) {
        console.error('[DB] Failed to get summary, fallback to memory:', err);
      }
    }
    return this.summaries.get(threadId) || null;
  }

  async getFacts(threadId: string): Promise<string[]> {
    return this.facts.get(threadId) || [];
  }

  async searchSnippets(threadId: string, query: string): Promise<string[]> {
    return this.snippets.get(threadId) || [];
  }

  async updateSummary(threadId: string, turn: any) {
    const nextSummary = `This conversation contains historical feedback loops and execution states including model coordination. User request: "${turn.userMessage}". Chosen response: "${turn.finalResponse.chosen.content.slice(0, 100)}..."`;
    
    if (this.pool) {
      try {
        await this.pool.query(
          'INSERT INTO memory_summaries (thread_id, summary, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (thread_id) DO UPDATE SET summary = $2, updated_at = NOW()',
          [threadId, nextSummary]
        );
        return;
      } catch (err) {
        console.error('[DB] Failed to update summary, fallback to memory:', err);
      }
    }
    this.summaries.set(threadId, nextSummary);
  }

  async updateFacts(threadId: string, turn: any) {
    const current = this.facts.get(threadId) || [];
    if (turn.userMessage.toLowerCase().includes('fact:')) {
      current.push(turn.userMessage);
    }
    this.facts.set(threadId, current);
  }

  async indexSnippets(threadId: string, turn: any) {
    const current = this.snippets.get(threadId) || [];
    current.push(turn.userMessage);
    this.snippets.set(threadId, current);
  }
}
