import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Settings, Loader2, AlertCircle, Save, Eye, EyeOff, Check } from 'lucide-react';
import { getProviderConfig, saveProviderConfig } from './api';
import type { SaveProviderConfigRequest } from './types';

export function SettingsPage() {
  return (
    <div>
      <div className="mb-8">
        <div className="text-[13px] text-text-muted mb-1 font-mono">$ settings</div>
        <h1 className="text-[28px] font-semibold tracking-tight">设置</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <SettingsCard />
      </div>
    </div>
  );
}

function SettingsCard() {
  const queryClient = useQueryClient();
  const {
    data: config,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['provider-config'],
    queryFn: getProviderConfig,
  });

  const [providerType, setProviderType] = useState('openai_compatible');
  const [baseUrl, setBaseUrl] = useState('');
  const [apiKey, setApiKey] = useState('');
  const [defaultModel, setDefaultModel] = useState('');
  const [showApiKey, setShowApiKey] = useState(false);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [initialized, setInitialized] = useState(false);

  if (config && !initialized) {
    setProviderType(config.providerType || 'openai_compatible');
    setBaseUrl(config.baseUrl || '');
    setDefaultModel(config.defaultModel || '');
    setInitialized(true);
  }

  const handleSave = async () => {
    setSaving(true);
    setSaveError('');
    setSaveSuccess(false);
    try {
      const req: SaveProviderConfigRequest = {
        providerType,
        baseUrl,
        defaultModel,
      };
      if (apiKey) req.apiKey = apiKey;
      await saveProviderConfig(req);
      setApiKey('');
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 2000);
      queryClient.invalidateQueries({ queryKey: ['provider-config'] });
    } catch (err) {
      setSaveError((err as any)?.message || '保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="border border-border rounded-xl p-6 bg-bg-elevated flex items-center justify-center py-12 text-text-muted text-sm">
        <Loader2 size={18} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (isError) {
    return (
      <div className="border border-border rounded-xl p-6 bg-bg-elevated flex items-center justify-center py-12 text-danger text-sm">
        <AlertCircle size={18} className="mr-2 shrink-0" />
        {(error as Error)?.message || '加载失败'}
      </div>
    );
  }

  return (
    <div className="border border-border rounded-xl bg-bg-elevated">
      <div className="px-5 py-4 border-b border-border flex items-center gap-2.5">
        <Settings size={16} className="text-text-muted" />
        <span className="text-sm font-medium">模型提供方配置</span>
        {config?.hasApiKey && (
          <span className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-medium bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950 dark:text-emerald-300 dark:border-emerald-800">
            <Check size={10} />
            已配置
          </span>
        )}
      </div>

      <div className="p-5 space-y-4">
        <div>
          <label className="text-xs font-medium text-text-muted block mb-1.5">提供方类型</label>
          <select
            value={providerType}
            onChange={(e) => setProviderType(e.target.value)}
            className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none cursor-pointer hover:border-border-strong focus:border-accent"
          >
            <option value="openai_compatible">OpenAI Compatible</option>
          </select>
        </div>

        <div>
          <label className="text-xs font-medium text-text-muted block mb-1.5">Base URL</label>
          <input
            type="text"
            value={baseUrl}
            onChange={(e) => setBaseUrl(e.target.value)}
            placeholder="https://api.openai.com/v1"
            className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] font-mono text-[13px]"
          />
        </div>

        <div>
          <label className="text-xs font-medium text-text-muted block mb-1.5">
            API Key
            {config?.hasApiKey && <span className="text-text-subtle font-normal ml-2">（已设置，留空则保持不变）</span>}
          </label>
          <div className="relative">
            <input
              type={showApiKey ? 'text' : 'password'}
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder={config?.hasApiKey ? 'sk-••••••••' : 'sk-...'}
              className="w-full h-9 px-3 pr-9 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] font-mono text-[13px]"
            />
            <button
              type="button"
              onClick={() => setShowApiKey(!showApiKey)}
              className="absolute right-2 top-1/2 -translate-y-1/2 w-6 h-6 grid place-items-center text-text-subtle hover:text-text cursor-pointer bg-transparent border-none rounded"
            >
              {showApiKey ? <EyeOff size={14} /> : <Eye size={14} />}
            </button>
          </div>
        </div>

        <div>
          <label className="text-xs font-medium text-text-muted block mb-1.5">默认模型</label>
          <input
            type="text"
            value={defaultModel}
            onChange={(e) => setDefaultModel(e.target.value)}
            placeholder="gpt-4o"
            className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] font-mono text-[13px]"
          />
        </div>

        <div className="flex items-center gap-3 pt-2">
          <button
            onClick={handleSave}
            disabled={saving}
            className="h-9 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-2 font-medium"
          >
            {saving ? <Loader2 size={14} className="animate-spin" /> : <Save size={14} />}
            保存
          </button>
          {saveSuccess && <span className="text-sm text-success">已保存</span>}
          {saveError && <span className="text-sm text-danger">{saveError}</span>}
        </div>
      </div>
    </div>
  );
}
