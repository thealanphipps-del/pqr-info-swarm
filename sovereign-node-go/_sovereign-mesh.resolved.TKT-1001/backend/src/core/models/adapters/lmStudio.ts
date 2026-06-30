import { NormalizedLLMResponse } from '../../types';

type ModelConfig = {
  id: string;
  endpoint?: string;
};

export class LmStudioAdapter {
  private endpoint: string;
  private modelId: string;

  constructor(cfg: ModelConfig) {
    this.endpoint = cfg.endpoint || 'http://localhost:1234/v1/chat/completions';
    this.modelId = cfg.id;
  }

  async complete(prompt: string): Promise<NormalizedLLMResponse> {
    try {
      const res = await fetch(this.endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: this.modelId,
          messages: [{ role: 'user', content: prompt }],
          stream: false,
        }),
      });

      if (!res.ok) {
        throw new Error(`LM Studio / Ollama error: ${res.status} ${res.statusText}`);
      }

      const data: any = await res.json();
      const content = data.choices?.[0]?.message?.content ?? '';

      return {
        modelId: this.modelId,
        content,
        raw: data,
      };
    } catch (err: any) {
      console.warn(`[Local LLM Adapter] Connection failed to ${this.endpoint}: ${err.message}. Falling back to Sandbox Mode.`);
      return {
        modelId: this.modelId,
        content: `[LM Studio Local Sandbox] Connection to ${this.endpoint} failed. Prompt received: "${prompt.slice(0, 100)}..."`,
        raw: { sandbox: true }
      };
    }
  }
}
