export type PromptMessageFormat = 'single_text' | 'chat_messages';
export type PromptVisibility = 'private' | 'workspace';
export type PromptStatus = 'draft' | 'active' | 'deprecated' | 'archived';
export type MessageRole = 'system' | 'user' | 'assistant';

export type Prompt = {
  id: number;
  workspaceId: number;
  createdBy: number;
  title: string;
  slug: string;
  description: string;
  body: string;
  messageFormat: PromptMessageFormat;
  visibility: PromptVisibility;
  status: PromptStatus;
  targetProvider: string;
  targetModel: string;
  defaultTemperature: number | null;
  defaultMaxTokens: number | null;
  usageNotes: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
};

export type ChatMessage = {
  role: MessageRole;
  content: string;
};

export type PromptListParams = {
  keyword?: string;
  status?: PromptStatus[];
  tags?: string[];
  provider?: string;
  model?: string;
  page?: number;
  pageSize?: number;
};

export type PromptListResponse = {
  items: Prompt[];
  total: number;
  page: number;
  pageSize: number;
};

export type PromptVariable = {
  id: number;
  name: string;
  label: string;
  description: string;
  required: boolean;
  defaultValue: string;
  exampleValue: string;
};

export type TestCase = {
  id: string;
  name: string;
  variableValues: Record<string, string>;
  expectedBehavior: string;
  expectedOutput: string;
  createdAt: string;
  updatedAt: string;
};

export type RunRecord = {
  id: string;
  promptVersionId: string;
  testCaseId: string | null;
  provider: string;
  model: string;
  inputVariables: Record<string, string>;
  renderedPrompt: string;
  output: string;
  latencyMs: number;
  tokenUsage: { promptTokens: number; completionTokens: number; totalTokens: number } | null;
  error: string | null;
  createdBy: string;
  createdAt: string;
};

export type VersionSnapshot = {
  content: string;
  messages: ChatMessage[];
  variables: PromptVariable[];
  targetProvider: string;
  targetModel: string;
  status: PromptStatus;
  tags: string[];
};

export type PromptVersion = {
  id: number;
  promptId: number;
  versionNumber: number;
  changeNote: string;
  author: string;
  snapshot: VersionSnapshot;
  createdAt: string;
};
