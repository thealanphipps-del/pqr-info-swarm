export type ChatMode = 'single' | 'parallel' | 'arbitrated' | 'swarm';

export type ToolCall = {
  name: string;
  arguments: any;
};

export type NormalizedLLMResponse = {
  modelId: string;
  content: string;
  toolCalls?: ToolCall[];
  reasoning?: string | null;
  usage?: { promptTokens: number; completionTokens: number };
  raw?: any;
};

export type ChatRequest = {
  threadId: string | null;
  userId: string;
  message: string;
  participants: string[]; // model IDs or agent IDs
  mode: ChatMode;
};

export type ChatTurnContext = {
  threadId: string;
  userId: string;
  history: any[]; // normalized messages
  memoryContext: {
    summary: string | null;
    facts: string[];
    snippets: string[];
  };
};
