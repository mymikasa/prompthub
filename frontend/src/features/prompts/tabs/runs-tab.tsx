import { useState, useEffect } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Play,
  Loader2,
  AlertCircle,
  Clock,
  Cpu,
  ChevronDown,
  ChevronRight,
  Copy,
  Check,
} from 'lucide-react';
import { runPrompt, getRuns, getVariables } from '../api';
import { listProviders } from '../../settings/api';
import type { Prompt, PromptVariable, RunRecord, ProviderConfig } from '../types';

export function RunsTab({ prompt }: { prompt: Prompt }) {
  const queryClient = useQueryClient();
  const [providers, setProviders] = useState<ProviderConfig[]>([]);
  const [variables, setVariables] = useState<PromptVariable[]>([]);
  const [varValues, setVarValues] = useState<Record<string, string>>({});
  const [selectedProviderId, setSelectedProviderId] = useState<number | null>(null);
  const [modelOverride, setModelOverride] = useState('');
  const [running, setRunning] = useState(false);
  const [runError, setRunError] = useState('');
  const [lastOutput, setLastOutput] = useState<RunRecord | null>(null);
  const [copied, setCopied] = useState(false);
  const [expandedRun, setExpandedRun] = useState<number | null>(null);

  const selectedProvider = providers.find((p) => p.id === selectedProviderId);

  const { data: runsData, isLoading: runsLoading } = useQuery({
    queryKey: ['promptRuns', prompt.id],
    queryFn: () => getRuns(prompt.id),
  });

  useEffect(() => {
    listProviders().then((list) => {
      setProviders(list);
      if (prompt.providerConfigId) {
        setSelectedProviderId(prompt.providerConfigId);
      } else if (list.length === 1) {
        setSelectedProviderId(list[0].id);
      }
    }).catch(() => {});

    getVariables(prompt.id).then((vars) => {
      setVariables(vars);
      const defaults: Record<string, string> = {};
      vars.forEach((v) => {
        if (v.defaultValue) defaults[v.name] = v.defaultValue;
        else if (v.exampleValue) defaults[v.name] = v.exampleValue;
      });
      setVarValues(defaults);
    }).catch(() => {});
  }, [prompt.id, prompt.providerConfigId]);

  const handleRun = async () => {
    setRunning(true);
    setRunError('');
    setLastOutput(null);
    try {
      const result = await runPrompt(prompt.id, {
        variables: varValues,
        providerId: selectedProviderId,
        model: modelOverride || undefined,
      });
      setLastOutput(result.run);
      queryClient.invalidateQueries({ queryKey: ['promptRuns', prompt.id] });
    } catch (err: any) {
      setRunError(err?.message || '执行失败');
    } finally {
      setRunning(false);
    }
  };

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const runs = runsData?.items ?? [];

  return (
    <div className="space-y-6">
      <div className="border border-border rounded-xl bg-bg-elevated overflow-hidden">
        <div className="px-5 py-3.5 border-b border-border flex items-center gap-2">
          <Play size={15} className="text-accent" />
          <span className="text-sm font-semibold">执行</span>
        </div>

        <div className="p-5 space-y-4">
          {variables.length > 0 && (
            <div>
              <div className="text-xs font-medium text-text-muted mb-2">变量</div>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
                {variables.map((v) => (
                  <div key={v.name}>
                    <label className="text-[11px] font-mono text-text-subtle block mb-1">
                      {`{{${v.name}}}`}
                      {v.required && <span className="text-danger ml-0.5">*</span>}
                    </label>
                    <input
                      className="w-full h-8 px-2.5 bg-bg border border-border rounded-md text-sm text-text outline-none hover:border-border-strong focus:border-accent font-mono text-[13px]"
                      placeholder={v.exampleValue || v.defaultValue || v.name}
                      value={varValues[v.name] || ''}
                      onChange={(e) => setVarValues({ ...varValues, [v.name]: e.target.value })}
                    />
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="flex flex-wrap items-end gap-3">
            <div>
              <label className="text-[11px] font-medium text-text-muted block mb-1">服务</label>
              <select
                value={selectedProviderId || ''}
                onChange={(e) => setSelectedProviderId(e.target.value ? Number(e.target.value) : null)}
                className="h-8 px-2.5 bg-bg border border-border rounded-md text-xs text-text outline-none cursor-pointer hover:border-border-strong focus:border-accent min-w-[160px]"
              >
                <option value="">未选择</option>
                {providers.map((p) => (
                  <option key={p.id} value={p.id}>{p.name}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="text-[11px] font-medium text-text-muted block mb-1">模型</label>
              <input
                className="h-8 px-2.5 bg-bg border border-border rounded-md text-xs text-text outline-none hover:border-border-strong focus:border-accent font-mono min-w-[140px]"
                placeholder={selectedProvider?.defaultModel || 'gpt-4o'}
                value={modelOverride}
                onChange={(e) => setModelOverride(e.target.value)}
              />
            </div>

            <button
              onClick={handleRun}
              disabled={running || !selectedProviderId}
              className="h-8 px-4 rounded-lg border border-transparent text-xs cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-1.5 font-medium"
            >
              {running ? <Loader2 size={12} className="animate-spin" /> : <Play size={12} />}
              {running ? '执行中…' : '执行'}
            </button>
          </div>

          {!selectedProviderId && providers.length === 0 && (
            <div className="text-xs text-amber-600 dark:text-amber-400">
              请先在「设置 → API Keys」中配置一个大模型服务。
            </div>
          )}

          {runError && (
            <div className="p-3 rounded-lg bg-red-50/50 border border-red-200 dark:bg-red-950/30 dark:border-red-900 text-xs text-danger">
              {runError}
            </div>
          )}

          {lastOutput && (
            <div className="border border-border rounded-lg overflow-hidden">
              <div className="px-4 py-2.5 bg-bg-subtle border-b border-border flex items-center gap-2">
                <span className="text-xs font-medium">执行结果</span>
                <span className="text-[11px] font-mono text-text-muted">{lastOutput.provider}/{lastOutput.model}</span>
                {lastOutput.latency > 0 && (
                  <span className="text-[11px] font-mono text-text-subtle">{lastOutput.latency}ms</span>
                )}
                {lastOutput.tokenUsage && (
                  <span className="text-[11px] font-mono text-text-subtle">{lastOutput.tokenUsage.totalTokens} tokens</span>
                )}
                <button
                  onClick={() => handleCopy(lastOutput.outputText)}
                  className="ml-auto h-6 w-6 rounded border border-border text-text-subtle hover:text-text grid place-items-center cursor-pointer bg-transparent transition-colors"
                >
                  {copied ? <Check size={11} /> : <Copy size={11} />}
                </button>
              </div>
              <div className="p-4 text-sm text-text leading-relaxed whitespace-pre-wrap max-h-[400px] overflow-auto font-mono text-[13px]">
                {lastOutput.outputText || '（空输出）'}
              </div>
            </div>
          )}
        </div>
      </div>

      <div>
        <div className="text-sm font-semibold mb-3">运行记录</div>

        {runsLoading && (
          <div className="flex items-center justify-center py-12 text-text-muted text-sm">
            <Loader2 size={16} className="animate-spin mr-2" />
            加载中…
          </div>
        )}

        {!runsLoading && runs.length === 0 && (
          <div className="py-12 text-center text-text-muted text-sm">
            <AlertCircle size={20} className="mx-auto mb-2 text-text-subtle" />
            暂无运行记录
          </div>
        )}

        {!runsLoading && runs.length > 0 && (
          <div className="space-y-2">
            {runs.map((run) => (
              <div
                key={run.id}
                className="border border-border rounded-lg bg-bg-elevated overflow-hidden"
              >
                <div
                  className="flex items-center gap-3 px-4 py-2.5 cursor-pointer hover:bg-bg-subtle transition-colors"
                  onClick={() => setExpandedRun(expandedRun === run.id ? null : run.id)}
                >
                  {expandedRun === run.id ? (
                    <ChevronDown size={12} className="text-text-subtle shrink-0" />
                  ) : (
                    <ChevronRight size={12} className="text-text-subtle shrink-0" />
                  )}
                  <div className="flex items-center gap-1.5 text-xs text-text-muted">
                    <Cpu size={12} />
                    <span className="font-mono">{run.provider}/{run.model}</span>
                  </div>
                  <div className="flex items-center gap-1.5 text-xs text-text-subtle">
                    <Clock size={12} />
                    {formatDateTime(run.createdAt)}
                  </div>
                  {run.latency > 0 && (
                    <span className="text-xs text-text-subtle font-mono">{formatMs(run.latency)}</span>
                  )}
                  {run.tokenUsage && (
                    <span className="text-xs text-text-subtle font-mono">{run.tokenUsage.totalTokens} tokens</span>
                  )}
                  <div className="ml-auto">
                    {run.errorMessage ? (
                      <span className="text-xs text-danger">失败</span>
                    ) : (
                      <span className="text-xs text-success">成功</span>
                    )}
                  </div>
                </div>

                {expandedRun === run.id && (
                  <div className="border-t border-border">
                    {run.errorMessage && (
                      <div className="px-4 py-2 text-xs text-danger bg-red-50/50 border-b border-border">
                        {run.errorMessage}
                      </div>
                    )}
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-0 divide-x divide-border">
                      <div className="p-4">
                        <div className="text-[11px] font-medium text-text-subtle mb-1.5">输入变量</div>
                        {Object.keys(run.inputVariables || {}).length > 0 ? (
                          <div className="flex flex-wrap gap-1">
                            {Object.entries(run.inputVariables).map(([k, v]) => (
                              <span key={k} className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-bg-subtle border border-border text-text-muted">
                                {k}={v}
                              </span>
                            ))}
                          </div>
                        ) : (
                          <span className="text-[11px] text-text-subtle">无变量</span>
                        )}
                      </div>
                      <div className="p-4">
                        <div className="flex items-center gap-2 mb-1.5">
                          <span className="text-[11px] font-medium text-text-subtle">输出</span>
                          {run.outputText && (
                            <button
                              onClick={() => handleCopy(run.outputText)}
                              className="h-4 w-4 rounded text-text-subtle hover:text-text grid place-items-center cursor-pointer bg-transparent border-none"
                            >
                              <Copy size={10} />
                            </button>
                          )}
                        </div>
                        <div className="text-xs text-text leading-relaxed whitespace-pre-wrap max-h-[160px] overflow-auto font-mono">
                          {run.outputText || '—'}
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function formatDateTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

function formatMs(ms: number) {
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(1)}s`;
}
