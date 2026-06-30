import { ChatRequest, ChatTurnContext, NormalizedLLMResponse } from './types';
import { MemoryService } from './memory/memoryService';
import { ModelRouter } from './models/modelRouter';
import { ArbitrationEngine } from './arbitration';
import { ThreadsRepo, MessagesRepo } from '../db/repositories/threadsRepo';

export class ConversationOrchestrator {
  constructor(
    private memory: MemoryService,
    private models: ModelRouter,
    private arbitration: ArbitrationEngine,
    private threadsRepo: ThreadsRepo,
    private messagesRepo: MessagesRepo
  ) {}

  async handleChat(req: ChatRequest) {
    const threadId =
      req.threadId ?? (await this.threadsRepo.createThread(req.userId));

    const history = await this.messagesRepo.getThreadMessages(threadId);

    const memoryContext = await this.memory.getContext(threadId, req.message);

    const ctx: ChatTurnContext = {
      threadId,
      userId: req.userId,
      history,
      memoryContext,
    };

    const prompt = this.buildPrompt(ctx, req.message);

    const modelResponses = await this.models.dispatch({
      mode: req.mode,
      participants: req.participants,
      prompt,
    });

    const final = await this.arbitration.decide(req.mode, modelResponses, {
      threadId,
      userId: req.userId,
      prompt,
      memoryContext,
    });

    await this.messagesRepo.appendTurn(threadId, {
      userId: req.userId,
      userMessage: req.message,
      modelResponses,
      final,
    });

    await this.memory.updateMemory(threadId, {
      userMessage: req.message,
      finalResponse: final,
      modelResponses,
    });

    return { threadId, final, modelResponses };
  }

  private buildPrompt(ctx: ChatTurnContext, userMessage: string): string {
    const { summary, facts, snippets } = ctx.memoryContext;

    const memoryBlock = [
      summary ? `Summary:\n${summary}` : null,
      facts.length ? `Facts:\n- ${facts.join('\n- ')}` : null,
      snippets.length ? `Relevant snippets:\n${snippets.join('\n')}` : null,
    ]
      .filter(Boolean)
      .join('\n\n');

    const recent = ctx.history.slice(-10); // rolling context trimming

    const historyText = recent
      .map((m: any) => `${m.role.toUpperCase()}: ${m.content}`)
      .join('\n');

    return [
      'You are part of a multi-model system. Respond clearly and concisely.',
      memoryBlock ? `\n[MEMORY]\n${memoryBlock}` : '',
      '\n[HISTORY]\n',
      historyText,
      '\n[USER]\n',
      userMessage,
    ].join('');
  }
}
