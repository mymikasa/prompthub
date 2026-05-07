import { apiClient } from '@/lib/api';
import type { ProviderConfig, SaveProviderConfigRequest } from './types';

export function getProviderConfig() {
  return apiClient<ProviderConfig>('/settings/provider');
}

export function saveProviderConfig(data: SaveProviderConfigRequest) {
  return apiClient<ProviderConfig>('/settings/provider', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}
