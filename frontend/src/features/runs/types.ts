export type RunRecord = {
  id: number;
  promptId: number;
  promptVersionId: number;
  testCaseId: number | null;
  provider: string;
  model: string;
  inputVariables: Record<string, string>;
  renderedPromptSnapshot: any;
  outputText: string;
  latency: number;
  tokenUsage: { promptTokens: number; completionTokens: number; totalTokens: number } | null;
  errorMessage: string;
  createdBy: number;
  createdAt: string;
};

export type RunListResponse = {
  items: RunRecord[];
  total: number;
  page: number;
  pageSize: number;
};
