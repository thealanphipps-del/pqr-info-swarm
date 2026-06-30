import { Router } from 'express';
import { ConversationOrchestrator } from '../core/orchestrator';
import { MemoryService } from '../core/memory/memoryService';
import { ModelRouter } from '../core/models/modelRouter';
import { ArbitrationEngine } from '../core/arbitration';
import { ThreadsRepo, MessagesRepo } from '../db/repositories/threadsRepo';
import { MemoryRepo } from '../db/repositories/memoryRepo';

export const chatRouter = Router();

// Singletons for orchestration
const threadsRepo = new ThreadsRepo();
const messagesRepo = new MessagesRepo();
const memoryRepo = new MemoryRepo();
const memoryService = new MemoryService(memoryRepo);
const modelRouter = new ModelRouter();
const arbitration = new ArbitrationEngine();
const orchestrator = new ConversationOrchestrator(
  memoryService,
  modelRouter,
  arbitration,
  threadsRepo,
  messagesRepo
);

chatRouter.post('/send', async (req, res) => {
  try {
    const { threadId, userId, message, participants, mode } = req.body;

    const result = await orchestrator.handleChat({
      threadId: threadId ?? null,
      userId: userId || 'anonymous',
      message: message || '',
      participants: participants ?? [],
      mode: mode ?? 'single',
    });

    res.json(result);
  } catch (err: any) {
    console.error(err);
    res.status(500).json({ error: err.message ?? 'Internal error' });
  }
});
