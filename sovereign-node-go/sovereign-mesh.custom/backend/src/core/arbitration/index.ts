import { NormalizedLLMResponse } from '../types';

export type ArbitrationContext = {
  threadId: string;
  userId: string;
  prompt: string;
  memoryContext: any;
};

export type ArbitrationResult = {
  chosen: NormalizedLLMResponse;
  rationale?: string;
};

export class ArbitrationEngine {
  async decide(
    mode: 'single' | 'parallel' | 'arbitrated' | 'swarm',
    inputs: NormalizedLLMResponse[],
    ctx: ArbitrationContext
  ): Promise<ArbitrationResult> {
    if (mode === 'single' || inputs.length === 0) {
      return { chosen: inputs[0] || { modelId: 'fallback', content: 'No responses received' }, rationale: 'single-model' };
    }

    if (mode === 'parallel') {
      return { chosen: inputs[0], rationale: 'parallel-first' };
    }

    // Heuristics: pick longest answer
    const sorted = [...inputs].sort(
      (a, b) => (b.content?.length ?? 0) - (a.content?.length ?? 0)
    );
    return {
      chosen: sorted[0],
      rationale: 'heuristic-longest',
    };
  }
}
