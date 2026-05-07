import { apiClient } from '@/lib/api';
import type {
  Prompt,
  PromptListParams,
  PromptListResponse,
  PromptVariable,
  TestCase,
  RunRecord,
  PromptVersion,
} from './types';

type Id = number | string;

export function getPrompts(params: PromptListParams) {
  const search = new URLSearchParams();
  if (params.keyword) search.set('keyword', params.keyword);
  if (params.status?.length) params.status.forEach((s) => search.append('status', s));
  if (params.tags?.length) params.tags.forEach((t) => search.append('tags', t));
  if (params.provider) search.set('provider', params.provider);
  if (params.model) search.set('model', params.model);
  if (params.page) search.set('page', String(params.page));
  if (params.pageSize) search.set('pageSize', String(params.pageSize));
  const qs = search.toString();
  return apiClient<PromptListResponse>(`/prompts${qs ? `?${qs}` : ''}`);
}

export function getPrompt(id: Id) {
  return apiClient<Prompt>(`/prompts/${id}`);
}

export function createPrompt(data: Partial<Prompt>) {
  return apiClient<Prompt>('/prompts', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function updatePrompt(id: Id, data: Partial<Prompt>) {
  return apiClient<Prompt>(`/prompts/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export function archivePrompt(id: Id) {
  return apiClient<void>(`/prompts/${id}/archive`, { method: 'POST' });
}

export function restorePrompt(id: Id) {
  return apiClient<void>(`/prompts/${id}/restore`, { method: 'POST' });
}

export function getTags() {
  return apiClient<string[]>('/tags');
}

export function getVariables(promptId: Id) {
  return apiClient<PromptVariable[]>(`/prompts/${promptId}/variables`);
}

export function updateVariable(promptId: Id, variableId: Id, data: Partial<PromptVariable>) {
  return apiClient<PromptVariable>(
    `/prompts/${promptId}/variables/${variableId}`,
    { method: 'PATCH', body: JSON.stringify(data) },
  );
}

export function getTestCases(promptId: Id) {
  return apiClient<TestCase[]>(`/prompts/${promptId}/test-cases`);
}

export function createTestCase(promptId: Id, data: Partial<TestCase>) {
  return apiClient<TestCase>(`/prompts/${promptId}/test-cases`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function updateTestCase(promptId: Id, testCaseId: Id, data: Partial<TestCase>) {
  return apiClient<TestCase>(
    `/prompts/${promptId}/test-cases/${testCaseId}`,
    { method: 'PATCH', body: JSON.stringify(data) },
  );
}

export function deleteTestCase(promptId: Id, testCaseId: Id) {
  return apiClient<void>(
    `/prompts/${promptId}/test-cases/${testCaseId}`,
    { method: 'DELETE' },
  );
}

export function getRuns(promptId: Id) {
  return apiClient<RunRecord[]>(`/prompts/${promptId}/runs`);
}

export function getVersions(promptId: Id) {
  return apiClient<PromptVersion[]>(`/prompts/${promptId}/versions`);
}

export function getVersion(promptId: Id, versionId: Id) {
  return apiClient<PromptVersion>(`/prompts/${promptId}/versions/${versionId}`);
}

export function restoreVersion(promptId: Id, versionId: Id) {
  return apiClient<PromptVersion>(
    `/prompts/${promptId}/versions/${versionId}/restore`,
    { method: 'POST' },
  );
}
