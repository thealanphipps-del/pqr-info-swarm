import { NormalizedLLMResponse } from '../../types';
import { AppConfig, loadConfig } from '../../../config';

type ModelConfig = {
  id: string;
};

export class VertexGeminiAdapter {
  private cfg: AppConfig;
  private modelId: string;

  constructor(modelConfig: ModelConfig) {
    this.cfg = loadConfig();
    this.modelId = modelConfig.id;
  }

  async complete(prompt: string): Promise<NormalizedLLMResponse> {
    const vertexKey = process.env.GEMINI_API_KEY || process.env.VERTEX_API_KEY;
    if (!vertexKey) {
      // Graceful local test fallback
      return {
        modelId: this.modelId,
        content: `[Gemini Sandbox Local Fallback] Received prompt: "${prompt.slice(0, 100)}...". Set GEMINI_API_KEY environment variable to activate real calls.`,
        usage: { promptTokens: prompt.length / 4, completionTokens: 25 },
        raw: { sandbox: true }
      };
    }

    try {
      const response = await fetch(`https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent?key=${vertexKey}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          contents: [{ parts: [{ text: prompt }] }]
        })
      });

      if (!response.ok) {
        throw new Error(`Gemini API returned status ${response.status}`);
      }

      const json: any = await response.json();
      const text = json.candidates?.[0]?.content?.parts?.[0]?.text || 'No response content';
      return {
        modelId: this.modelId,
        content: text,
        raw: json
      };
    } catch (err: any) {
      console.error('[Gemini Adapter] API call failed:', err);
      return {
        modelId: this.modelId,
        content: `[Gemini Error Fallback] Failed to call real API: ${err.message}`,
        raw: { error: err.message }
      };
    }
  }
}
