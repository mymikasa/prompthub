import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import {
  Activity,
  Loader2,
  AlertCircle,
  Clock,
  Cpu,
  ChevronDown,
  ChevronRight,
  Copy,
  Check,
  Filter,
  FileText,
} from 'lucide-react';
import { getAllRuns } from './api';

export function RunsPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState('');
  const [expandedRun, setExpandedRun] = useState<number | null>(null);
  const [copied, setCopied] = useState<number | null>(null);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['allRuns', page, statusFilter],
    queryFn: () => getAllRuns(page, 20, statusFilter || undefined),
  });

  const runs = data?.items ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / 20);

  const handleCopy = (id: number, text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 1500);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="text-[13px] text-text-muted mb-1 font-mono">$ runs list</div>
          <h1 className="text-[28px] font-semibold tracking-tight">运行记录</h1>
        </div>
        <div className="flex items-center gap-2">
          <Filter size={14} className="text-text-muted" />
          <select
            value={statusFilter}
            onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }}
            className="h-8 px-2.5 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none cursor-pointer hover:border-border-strong focus:border-accent"
          >
            <option value="">全部状态</option>
            <option value="success">成功</option>
            <option value="failed">失败</option>
          </select>
        </div>
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-20 text-text-muted text-sm">
          <Loader2 size={18} className="animate-spin mr-2" />
          加载中…
        </div>
      )}

      {isError && (
        <div className="flex items-center justify-center py-20 text-danger text-sm">
          <AlertCircle size={18} className="mr-2" />
          加载失败
        </div>
      )}

      {!isLoading && !isError && runs.length === 0 && (
        <div className="py-20 text-center text-text-muted text-sm">
          <Activity size={28} className="mx-auto mb-3 text-text-subtle" />
          暂无运行记录。前往提示词的「变量」tab 执行一次。
        </div>
      )}

      {!isLoading && !isError && runs.length > 0 && (
        <div className="space-y-2">
          {runs.map((run) => (
            <div
              key={run.id}
              className="border border-border rounded-lg bg-bg-elevated overflow-hidden"
            >
              <div
                className="flex items-center gap-3 px-4 py-3 cursor-pointer hover:bg-bg-subtle transition-colors"
                onClick={() => setExpandedRun(expandedRun === run.id ? null : run.id)}
              >
                {expandedRun === run.id ? (
                  <ChevronDown size={12} className="text-text-subtle shrink-0" />
                ) : (
                  <ChevronRight size={12} className="text-text-subtle shrink-0" />
                )}
                <button
                  onClick={(e) => { e.stopPropagation(); navigate(`/prompts/${run.promptId}`); }}
                  className="inline-flex items-center gap-1.5 text-xs text-accent hover:underline bg-transparent border-none cursor-pointer p-0 font-medium"
                >
                  <FileText size={11} />
                  Prompt #{run.promptId}
                </button>
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
                      <div className="text-xs text-text leading-relaxed whitespace-pre-wrap max-h-[200px] overflow-auto font-mono">
                        {run.outputText || '—'}
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}

          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 pt-4">
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page <= 1}
                className="h-8 px-3 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text disabled:opacity-40 disabled:cursor-default"
              >
                上一页
              </button>
              <span className="text-xs text-text-muted">
                {page} / {totalPages}
              </span>
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page >= totalPages}
                className="h-8 px-3 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text disabled:opacity-40 disabled:cursor-default"
              >
                下一页
              </button>
            </div>
          )}
        </div>
      )}
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
