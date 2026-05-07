import { useQuery } from '@tanstack/react-query';
import { Loader2, AlertCircle, Clock, Cpu } from 'lucide-react';
import { getRuns } from '../api';

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

export function RunsTab({ promptId }: { promptId: number | string }) {
  const { data: runs, isLoading, isError } = useQuery({
    queryKey: ['promptRuns', promptId],
    queryFn: () => getRuns(promptId),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-text-muted text-sm">
        <Loader2 size={16} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (isError || !runs?.length) {
    return (
      <div className="py-12 text-center text-text-muted text-sm">
        <AlertCircle size={20} className="mx-auto mb-2 text-text-subtle" />
        暂无运行记录
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
          <div className="flex items-center gap-3 px-4 py-2.5 border-b border-border">
            <div className="flex items-center gap-1.5 text-xs text-text-muted">
              <Cpu size={12} />
              <span className="font-mono">{run.provider}/{run.model}</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-text-subtle">
              <Clock size={12} />
              {formatDateTime(run.createdAt)}
            </div>
            <span className="text-xs text-text-subtle font-mono">{formatMs(run.latencyMs)}</span>
            {run.tokenUsage && (
              <span className="text-xs text-text-subtle font-mono">
                {run.tokenUsage.totalTokens} tokens
              </span>
            )}
            <div className="ml-auto">
              {run.error ? (
                <span className="text-xs text-danger">失败</span>
              ) : (
                <span className="text-xs text-success">成功</span>
              )}
            </div>
          </div>

          {run.error && (
            <div className="px-4 py-2 text-xs text-danger bg-red-50/50 border-b border-border">
              {run.error}
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-0 divide-x divide-border">
            <div className="p-4">
              <div className="text-[11px] font-medium text-text-subtle mb-1.5">输入变量</div>
              <div className="flex flex-wrap gap-1">
                {Object.entries(run.inputVariables || {}).map(([k, v]) => (
                  <span key={k} className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-bg-subtle border border-border text-text-muted">
                    {k}={v}
                  </span>
                ))}
              </div>
            </div>
            <div className="p-4">
              <div className="text-[11px] font-medium text-text-subtle mb-1.5">输出</div>
              <div className="text-xs text-text leading-relaxed whitespace-pre-wrap max-h-[120px] overflow-auto font-mono">
                {run.output || '—'}
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
