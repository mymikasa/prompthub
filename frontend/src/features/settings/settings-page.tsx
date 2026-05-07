import { useState } from 'react';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import {
  KeyRound,
  Loader2,
  AlertCircle,
  Plus,
  Trash2,
  Pencil,
  Eye,
  EyeOff,
  Check,
  X,
  Server,
} from 'lucide-react';
import { listProviders, saveProviderConfig, deleteProviderConfig } from './api';
import type { ProviderConfig, SaveProviderConfigRequest, ProviderTypeInfo } from './types';

const PROVIDER_TYPES: ProviderTypeInfo[] = [
  { type: 'openai_compatible', label: 'OpenAI Compatible', defaultBaseUrl: 'https://api.openai.com/v1' },
  { type: 'deepseek', label: 'DeepSeek', defaultBaseUrl: 'https://api.deepseek.com/v1' },
  { type: 'zhipu', label: '智谱 AI', defaultBaseUrl: 'https://open.bigmodel.cn/api/paas/v4' },
  { type: 'minimax', label: 'MiniMax', defaultBaseUrl: 'https://api.minimax.chat/v1' },
  { type: 'moonshot', label: 'Moonshot', defaultBaseUrl: 'https://api.moonshot.cn/v1' },
  { type: 'qwen', label: '通义千问', defaultBaseUrl: 'https://dashscope.aliyuncs.com/compatible-mode/v1' },
];

function getProviderLabel(type: string) {
  return PROVIDER_TYPES.find((p) => p.type === type)?.label ?? type;
}

export function SettingsPage() {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingConfig, setEditingConfig] = useState<ProviderConfig | null>(null);

  const {
    data: providers = [],
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['providers'],
    queryFn: listProviders,
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProviderConfig,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['providers'] }),
  });

  const handleAdd = () => {
    setEditingConfig(null);
    setModalOpen(true);
  };

  const handleEdit = (pc: ProviderConfig) => {
    setEditingConfig(pc);
    setModalOpen(true);
  };

  const handleDelete = (pc: ProviderConfig) => {
    if (!window.confirm(`确定要删除「${pc.name}」吗？`)) return;
    deleteMutation.mutate(pc.id);
  };

  const handleModalClose = () => {
    setModalOpen(false);
    setEditingConfig(null);
  };

  const handleSaved = () => {
    setModalOpen(false);
    setEditingConfig(null);
    queryClient.invalidateQueries({ queryKey: ['providers'] });
  };

  const configuredTypes = new Set(providers.map((p) => p.providerType));
  const availableTypes = PROVIDER_TYPES.filter((p) => !configuredTypes.has(p.type));
  const canAdd = availableTypes.length > 0;

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="text-[13px] text-text-muted mb-1 font-mono">$ settings providers</div>
          <h1 className="text-[28px] font-semibold tracking-tight">API Keys</h1>
          <p className="text-sm text-text-muted mt-1">配置大模型服务的 API Key，用于运行提示词。</p>
        </div>
        {canAdd && (
          <button
            onClick={handleAdd}
            className="h-10 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on transition-all duration-150 inline-flex items-center gap-2 font-medium"
          >
            <Plus size={14} />
            添加服务
          </button>
        )}
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-20 text-text-muted text-sm">
          <Loader2 size={18} className="animate-spin mr-2" />
          加载中…
        </div>
      )}

      {isError && (
        <div className="flex items-center justify-center py-20 text-danger text-sm">
          <AlertCircle size={18} className="mr-2 shrink-0" />
          {(error as Error)?.message || '加载失败'}
        </div>
      )}

      {!isLoading && !isError && providers.length === 0 && (
        <div className="border border-dashed border-border-strong rounded-xl py-16 px-8 text-center bg-bg-subtle">
          <div className="inline-flex p-4 bg-bg-elevated border border-border rounded-xl text-accent mb-4">
            <KeyRound size={24} />
          </div>
          <div className="text-base font-semibold mb-1">还没有配置 API Key</div>
          <p className="text-sm text-text-muted mb-5">添加一个大模型服务，即可开始运行提示词。</p>
          {canAdd && (
            <button
              onClick={handleAdd}
              className="h-10 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on transition-all duration-150 inline-flex items-center gap-2"
            >
              <Plus size={14} />
              添加服务
            </button>
          )}
        </div>
      )}

      {!isLoading && !isError && providers.length > 0 && (
        <div className="space-y-3">
          {providers.map((pc) => (
            <ProviderCard
              key={pc.id}
              config={pc}
              onEdit={() => handleEdit(pc)}
              onDelete={() => handleDelete(pc)}
              deleting={deleteMutation.isPending}
            />
          ))}
        </div>
      )}

      {modalOpen && (
        <ProviderModal
          config={editingConfig}
          availableTypes={editingConfig ? [] : availableTypes}
          onSaved={handleSaved}
          onClose={handleModalClose}
        />
      )}
    </div>
  );
}

function ProviderCard({
  config,
  onEdit,
  onDelete,
  deleting,
}: {
  config: ProviderConfig;
  onEdit: () => void;
  onDelete: () => void;
  deleting: boolean;
}) {
  return (
    <div className="border border-border rounded-xl bg-bg-elevated hover:border-border-strong transition-all duration-150">
      <div className="px-5 py-4 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-bg-subtle border border-border grid place-items-center shrink-0">
          <Server size={18} className="text-accent" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium">{config.name}</span>
            <span className="px-1.5 py-0.5 rounded text-[11px] font-mono text-text-muted bg-bg-subtle border border-border">
              {getProviderLabel(config.providerType)}
            </span>
            {config.hasApiKey ? (
              <span className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-medium bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950 dark:text-emerald-300 dark:border-emerald-800">
                <Check size={10} />
                已配置
              </span>
            ) : (
              <span className="px-1.5 py-0.5 rounded text-[11px] font-medium bg-amber-50 text-amber-700 border border-amber-200 dark:bg-amber-950 dark:text-amber-300 dark:border-amber-800">
                未配置 Key
              </span>
            )}
          </div>
          <div className="text-xs text-text-muted mt-1 truncate font-mono">
            {config.baseUrl}
            {config.defaultModel && <span className="ml-2 text-text-subtle">· {config.defaultModel}</span>}
          </div>
        </div>

        <div className="flex items-center gap-1.5 shrink-0">
          <button
            onClick={onEdit}
            className="h-8 w-8 rounded-lg border border-border text-text-muted grid place-items-center cursor-pointer bg-transparent hover:text-text hover:border-border-strong transition-all duration-150"
            title="编辑"
          >
            <Pencil size={13} />
          </button>
          <button
            onClick={onDelete}
            disabled={deleting}
            className="h-8 w-8 rounded-lg border border-border text-text-muted grid place-items-center cursor-pointer bg-transparent hover:text-danger hover:border-danger/30 disabled:opacity-50 transition-all duration-150"
            title="删除"
          >
            <Trash2 size={13} />
          </button>
        </div>
      </div>
    </div>
  );
}

function ProviderModal({
  config,
  availableTypes,
  onSaved,
  onClose,
}: {
  config: ProviderConfig | null;
  availableTypes: ProviderTypeInfo[];
  onSaved: () => void;
  onClose: () => void;
}) {
  const isEdit = !!config;
  const [providerType, setProviderType] = useState(config?.providerType ?? availableTypes[0]?.type ?? '');
  const [name, setName] = useState(config?.name ?? '');
  const [baseUrl, setBaseUrl] = useState(config?.baseUrl ?? '');
  const [apiKey, setApiKey] = useState('');
  const [defaultModel, setDefaultModel] = useState(config?.defaultModel ?? '');
  const [showApiKey, setShowApiKey] = useState(false);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  const handleTypeChange = (type: string) => {
    setProviderType(type);
    const info = PROVIDER_TYPES.find((p) => p.type === type);
    if (info && !baseUrl) {
      setBaseUrl(info.defaultBaseUrl);
    }
    if (!name) {
      setName(info?.label ?? type);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setSaveError('');
    try {
      const req: SaveProviderConfigRequest = {
        name,
        providerType,
        baseUrl,
        defaultModel: defaultModel || undefined,
      };
      if (apiKey) req.apiKey = apiKey;
      await saveProviderConfig(req);
      onSaved();
    } catch (err) {
      setSaveError((err as any)?.message || '保存失败');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm" onClick={onClose}>
      <div
        className="w-full max-w-lg bg-bg-elevated border border-border rounded-xl shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <span className="text-sm font-semibold">{isEdit ? '编辑服务' : '添加服务'}</span>
          <button
            onClick={onClose}
            className="w-7 h-7 rounded-md border border-border text-text-muted grid place-items-center cursor-pointer bg-transparent hover:text-text hover:border-border-strong"
          >
            <X size={14} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-5 space-y-4">
          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">服务类型</label>
            {isEdit ? (
              <div className="h-9 px-3 bg-bg-subtle border border-border rounded-lg text-sm text-text-muted flex items-center">
                {getProviderLabel(providerType)}
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-2">
                {availableTypes.map((pt) => (
                  <button
                    key={pt.type}
                    type="button"
                    onClick={() => handleTypeChange(pt.type)}
                    className={`h-10 px-3 rounded-lg border text-sm text-left cursor-pointer transition-all duration-150 ${
                      providerType === pt.type
                        ? 'border-accent bg-accent-soft text-accent font-medium'
                        : 'border-border bg-bg text-text-muted hover:border-border-strong hover:text-text'
                    }`}
                  >
                    {pt.label}
                  </button>
                ))}
              </div>
            )}
          </div>

          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">名称</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例如：我的 DeepSeek"
              className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)]"
              required
            />
          </div>

          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">Base URL</label>
            <input
              type="text"
              value={baseUrl}
              onChange={(e) => setBaseUrl(e.target.value)}
              placeholder="https://api.openai.com/v1"
              className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] font-mono text-[13px]"
              required
            />
          </div>

          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">
              API Key
              {isEdit && config?.hasApiKey && (
                <span className="text-text-subtle font-normal ml-2">（已设置，留空则保持不变）</span>
              )}
            </label>
            <div className="relative">
              <input
                type={showApiKey ? 'text' : 'password'}
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                placeholder={isEdit && config?.hasApiKey ? '••••••••' : 'sk-...'}
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
              placeholder="gpt-4o / deepseek-chat / glm-4-flash …"
              className="w-full h-9 px-3 bg-bg border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] font-mono text-[13px]"
            />
          </div>

          {saveError && (
            <div className="text-sm text-danger">{saveError}</div>
          )}

          <div className="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="h-9 px-4 rounded-lg border border-border text-sm cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={saving}
              className="h-9 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-2 font-medium"
            >
              {saving && <Loader2 size={14} className="animate-spin" />}
              {isEdit ? '保存' : '添加'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
