import { NormalizedLLMResponse } from '../types';
import { loadConfig, loadJsonConfig } from '../../config';
import { VertexGeminiAdapter } from './adapters/vertexGemini';
import { LmStudioAdapter } from './adapters/lmStudio';

type ModelConfig = {
  id: string;
  provider: 'vertex' | 'lmstudio' | 'generic';
  endpoint?: string;
  weight?: number;
  tags?: string[];
};

export type DispatchRequest = {
  mode: 'single' | 'parallel' | 'arbitrated' | 'swarm';
  participants: string[];
  prompt: string;
};

export class ModelRouter {
  private models: ModelConfig[];
  private adapters: Map<string, any>;

  constructor() {
    const cfg = loadConfig();
    try {
      this.models = loadJsonConfig<ModelConfig[]>(cfg.MODELS_CONFIG_PATH);
    } catch {
      this.models = [
        {
          id: "gemini-1.5-pro",
          provider: "vertex",
          weight: 1.0,
          tags: ["primary"]
        },
        {
          id: "lmstudio-local",
          provider: "lmstudio",
          endpoint: "http://localhost:1234/v1/chat/completions",
          weight: 0.8,
          tags: ["local"]
        }
      ];
    }
    this.adapters = new Map();

    for (const m of this.models) {
      if (m.provider === 'vertex') {
        this.adapters.set(m.id, new VertexGeminiAdapter(m));
      } else if (m.provider === 'lmstudio') {
        this.adapters.set(m.id, new LmStudioAdapter(m));
      }
    }
  }

  async dispatch(req: DispatchRequest): Promise<NormalizedLLMResponse[]> {
    const targets =
      req.participants.length > 0
        ? this.models.filter((m) => req.participants.includes(m.id))
        : this.models.filter((m) => m.tags?.includes('primary'));

    if (targets.length === 0 && this.models.length > 0) {
      targets.push(this.models[0]);
    }

    const calls = targets.map(async (m) => {
      const adapter = this.adapters.get(m.id);
      if (!adapter) throw new Error(`No adapter for model ${m.id}`);
      const res = await adapter.complete(req.prompt);
      return { ...res, modelId: m.id } as NormalizedLLMResponse;
    });

    return Promise.all(calls);
  }
}
