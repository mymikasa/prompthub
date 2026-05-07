import { useState, useMemo, useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import {
  Copy,
  Play,
  Loader2,
  Check,
} from 'lucide-react';
import { runPrompt } from '../api';
import { listProviders } from '../../settings/api';
import type { Prompt, RunRecord, ProviderConfig } from '../types';

function extractVariables(text: string): string[] {
  const matches = text.match(/\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g) || [];
  return [...new Set(matches.map((m) => m.slice(2, -2)))];
}

export function VariablesTab({ prompt }: { prompt: Prompt }) {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, string>>({});
  const [copied, setCopied] = useState(false);
  const [providers, setProviders] = useState<ProviderConfig[]>([]);
  const [selectedProviderId, setSelectedProviderId] = useState<number | null>(null);
  const [running, setRunning] = useState(false);
  const [runError, setRunError] = useState('');
  const [lastOutput, setLastOutput] = useState<RunRecord | null>(null);
  const [outputCopied, setOutputCopied] = useState(false);

  const variables = useMemo(() => extractVariables(prompt.body || ''), [prompt.body]);

  const rendered = useMemo(() => {
    return (prompt.body || '').replace(
      /\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g,
      (_, name) => values[name] ?? `{{${name}}}`,
    );
  }, [prompt.body, values]);

  const unfilled = variables.filter((v) => !values[v]?.trim());
  const hasContent = rendered.trim().length > 0;

  useEffect(() => {
    listProviders().then((list) => {
      setProviders(list);
      if (prompt.providerConfigId) {
        setSelectedProviderId(prompt.providerConfigId);
      } else if (list.length === 1) {
        setSelectedProviderId(list[0].id);
      }
    }).catch(() => {});
  }, [prompt.providerConfigId]);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(rendered);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const handleRun = async () => {
    setRunning(true);
    setRunError('');
    setLastOutput(null);
    try {
      const result = await runPrompt(prompt.id, {
        variables: values,
        providerId: selectedProviderId,
      });
      setLastOutput(result.run);
      queryClient.invalidateQueries({ queryKey: ['promptRuns', prompt.id] });
    } catch (err: any) {
      setRunError(err?.message || '执行失败');
    } finally {
      setRunning(false);
    }
  };

  const handleCopyOutput = () => {
    if (lastOutput?.outputText) {
      navigator.clipboard.writeText(lastOutput.outputText);
      setOutputCopied(true);
      setTimeout(() => setOutputCopied(false), 1500);
    }
  };

  return (
    <div className="space-y-6">
      {variables.length > 0 && (
        <div>
          <div className="text-xs font-medium text-text-muted mb-3">
            填写变量值
            {unfilled.length > 0 && (
              <span className="text-text-subtle ml-2">
                {unfilled.length} 个未填写
              </span>
            )}
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {variables.map((v) => (
              <div key={v} className="flex flex-col gap-1.5">
                <label className="font-mono text-xs text-text-subtle flex items-center gap-1">
                  <span className="opacity-60">{'{{'}</span>
                  {v}
                  <span className="opacity-60">{'}}'}</span>
                </label>
                <input
                  className="w-full h-9 px-3 bg-bg-elevated border border-border rounded-lg text-sm text-text outline-none transition-all duration-150 hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)]"
                  placeholder={`输入 ${v} 的值`}
                  value={values[v] || ''}
                  onChange={(e) => setValues({ ...values, [v]: e.target.value })}
                />
              </div>
            ))}
          </div>
        </div>
      )}

      <div>
        <div className="flex items-center justify-between mb-2">
          <div className="text-xs font-medium text-text-muted">渲染结果</div>
          <button
            onClick={handleCopy}
            disabled={!hasContent}
            className="h-7 px-2.5 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150 inline-flex items-center gap-1 disabled:opacity-40"
          >
            <Copy size={11} />
            {copied ? '已复制' : '复制文本'}
          </button>
        </div>

        <div className="bg-code-bg rounded-lg p-4 px-5 font-mono text-[13px] leading-[1.7] text-code-text whitespace-pre-wrap max-h-[400px] overflow-auto">
          {hasContent ? rendered : <span className="text-code-comment">暂无内容</span>}
        </div>

        {unfilled.length > 0 && hasContent && (
          <div className="mt-2 text-xs text-amber-600">
            仍有未填写的变量：{unfilled.map((v) => `{{${v}}}`).join(', ')}
          </div>
        )}
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <button
          onClick={handleRun}
          disabled={running || !selectedProviderId}
          className="h-9 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-2 font-medium"
        >
          {running ? <Loader2 size={14} className="animate-spin" /> : <Play size={14} />}
          {running ? '执行中…' : '执行'}
        </button>
        {providers.length > 1 && (
          <select
            value={selectedProviderId || ''}
            onChange={(e) => setSelectedProviderId(e.target.value ? Number(e.target.value) : null)}
            className="h-9 px-3 bg-bg-elevated border border-border rounded-lg text-xs text-text outline-none cursor-pointer hover:border-border-strong focus:border-accent"
          >
            {providers.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
        )}
        {!selectedProviderId && providers.length === 0 && (
          <span className="text-xs text-amber-600 dark:text-amber-400">请先在「设置 → API Keys」中配置大模型服务</span>
        )}
      </div>

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
              onClick={handleCopyOutput}
              className="ml-auto h-6 w-6 rounded border border-border text-text-subtle hover:text-text grid place-items-center cursor-pointer bg-transparent transition-colors"
            >
              {outputCopied ? <Check size={11} /> : <Copy size={11} />}
            </button>
          </div>
          <div className="p-4 text-sm text-text leading-relaxed whitespace-pre-wrap max-h-[400px] overflow-auto font-mono text-[13px]">
            {lastOutput.outputText || '（空输出）'}
          </div>
        </div>
      )}
    </div>
  );
}
