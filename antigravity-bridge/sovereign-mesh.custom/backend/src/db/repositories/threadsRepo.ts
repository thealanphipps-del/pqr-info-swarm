import { Pool } from 'pg';
import { loadConfig } from '../../config';
import { v4 as uuidv4 } from 'uuid';

export class ThreadsRepo {
  private pool: Pool | null = null;
  private memoryDb: Map<string, any> = new Map();

  constructor() {
    const config = loadConfig();
    try {
      this.pool = new Pool({
        connectionString: config.DATABASE_URL,
        connectionTimeoutMillis: 2000,
      });
      // Test connection asynchronously
      this.pool.query('SELECT 1').catch(() => {
        console.warn('[DB] PostgreSQL connection failed. Falling back to In-Memory DB.');
        this.pool = null;
      });
    } catch {
      this.pool = null;
    }
  }

  async createThread(userId: string): Promise<string> {
    const id = uuidv4();
    if (this.pool) {
      try {
        await this.pool.query(
          'INSERT INTO threads (id, user_id, created_at) VALUES ($1, $2, NOW())',
          [id, userId]
        );
        return id;
      } catch (err) {
        console.error('[DB] Failed to insert thread, fallback to memory:', err);
      }
    }
    this.memoryDb.set(id, { userId, createdAt: new Date() });
    return id;
  }
}

export class MessagesRepo {
  private pool: Pool | null = null;
  private memoryDb: Map<string, any[]> = new Map();

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

  async getThreadMessages(threadId: string): Promise<any[]> {
    if (this.pool) {
      try {
        const res = await this.pool.query(
          'SELECT role, content, metadata, created_at FROM messages WHERE thread_id = $1 ORDER BY created_at ASC',
          [threadId]
        );
        return res.rows;
      } catch (err) {
        console.error('[DB] Failed to get messages, fallback to memory:', err);
      }
    }
    return this.memoryDb.get(threadId) || [];
  }

  async appendTurn(threadId: string, turn: { userId: string; userMessage: string; modelResponses: any[]; final: any }) {
    const userMsg = {
      role: 'user',
      content: turn.userMessage,
      metadata: JSON.stringify({ userId: turn.userId }),
      created_at: new Date()
    };
    const assistantMsg = {
      role: 'assistant',
      content: turn.final.chosen.content,
      metadata: JSON.stringify({
        chosenModel: turn.final.chosen.modelId,
        rationale: turn.final.rationale,
        allResponses: turn.modelResponses
      }),
      created_at: new Date()
    };

    if (this.pool) {
      try {
        await this.pool.query(
          'INSERT INTO messages (thread_id, role, content, metadata, created_at) VALUES ($1, $2, $3, $4, NOW())',
          [threadId, userMsg.role, userMsg.content, userMsg.metadata]
        );
        await this.pool.query(
          'INSERT INTO messages (thread_id, role, content, metadata, created_at) VALUES ($1, $2, $3, $4, NOW())',
          [threadId, assistantMsg.role, assistantMsg.content, assistantMsg.metadata]
        );
        return;
      } catch (err) {
        console.error('[DB] Failed to append turn, fallback to memory:', err);
      }
    }

    const current = this.memoryDb.get(threadId) || [];
    current.push(userMsg);
    current.push(assistantMsg);
    this.memoryDb.set(threadId, current);
  }
}
