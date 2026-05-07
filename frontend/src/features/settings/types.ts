export type ProviderConfig = {
  id: number;
  providerType: string;
  name: string;
  baseUrl: string;
  hasApiKey: boolean;
  defaultModel: string;
};

export type SaveProviderConfigRequest = {
  name: string;
  providerType: string;
  baseUrl: string;
  apiKey?: string;
  defaultModel?: string;
};

export type ProviderTypeInfo = {
  type: string;
  label: string;
  defaultBaseUrl: string;
};
