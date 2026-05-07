import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Loader2,
  AlertCircle,
  Clock,
  Cpu,
  ChevronDown,
  ChevronRight,
  Copy,
  Check,
} from 'lucide-react';
import { getRuns } from '../api';

export function RunsTab({ promptId }: { promptId: number | string }) {
  const [expandedRun, setExpandedRun] = useState<number | null>(null);
  const [copied, setCopied] = useState<number | null>(null);

  const { data: runsData, isLoading } = useQuery({
    queryKey: ['promptRuns', promptId],
    queryFn: () => getRuns(promptId),
  });

  const runs = runsData?.items ?? [];

  const handleCopy = (id: number, text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 1500);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-text-muted text-sm">
        <Loader2 size={16} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (runs.length === 0) {
    return (
      <div className="py-12 text-center text-text-muted text-sm">
        <AlertCircle size={20} className="mx-auto mb-2 text-text-subtle" />
        暂无运行记录。前往「变量」tab 填写变量并执行。
      </div>
    );
  }

  return (
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
                        onClick={() => handleCopy(run.id, run.outputText)}
                        className="h-4 w-4 rounded text-text-subtle hover:text-text grid place-items-center cursor-pointer bg-transparent border-none"
                      >
                        {copied === run.id ? <Check size={10} /> : <Copy size={10} />}
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
