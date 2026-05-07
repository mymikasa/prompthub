export type ProviderConfig = {
  id: number;
  providerType: string;
  baseUrl: string;
  hasApiKey: boolean;
  defaultModel: string;
};

export type SaveProviderConfigRequest = {
  providerType: string;
  baseUrl: string;
  apiKey?: string;
  defaultModel?: string;
};
