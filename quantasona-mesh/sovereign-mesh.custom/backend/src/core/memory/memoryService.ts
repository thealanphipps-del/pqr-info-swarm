import { MemoryRepo } from '../../db/repositories/memoryRepo';

export type MemoryContext = {
  summary: string | null;
  facts: string[];
  snippets: string[];
};

export class MemoryService {
  constructor(private memoryRepo: MemoryRepo) {}

  async getContext(
    threadId: string,
    userMessage: string
  ): Promise<MemoryContext> {
    const summary = await this.memoryRepo.getSummary(threadId);
    const facts = await this.memoryRepo.getFacts(threadId);
    const snippets = await this.memoryRepo.searchSnippets(threadId, userMessage);

    return { summary, facts, snippets };
  }

  async updateMemory(threadId: string, turn: any) {
    await this.memoryRepo.updateSummary(threadId, turn);
    await this.memoryRepo.updateFacts(threadId, turn);
    await this.memoryRepo.indexSnippets(threadId, turn);
  }
}
