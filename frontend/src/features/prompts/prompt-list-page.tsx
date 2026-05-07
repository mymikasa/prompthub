import { useState, useEffect, useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Search,
  Plus,
  SlidersHorizontal,
  X,
  FileText,
  Loader2,
  AlertCircle,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { getPrompts } from '@/features/prompts/api';
import type { Prompt, PromptStatus, PromptListParams } from '@/features/prompts/types';

const PAGE_SIZE = 20;

const STATUS_OPTIONS: { value: PromptStatus; label: string }[] = [
  { value: 'draft', label: '草稿' },
  { value: 'active', label: '可用' },
  { value: 'deprecated', label: '已废弃' },
  { value: 'archived', label: '已归档' },
];

const STATUS_BADGE: Record<PromptStatus, string> = {
  draft: 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950 dark:text-amber-300 dark:border-amber-800',
  active: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950 dark:text-emerald-300 dark:border-emerald-800',
  deprecated: 'bg-gray-50 text-gray-600 border-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-700',
  archived: 'bg-stone-50 text-stone-500 border-stone-200 dark:bg-stone-800 dark:text-stone-400 dark:border-stone-700',
};

const STATUS_LABEL: Record<PromptStatus, string> = {
  draft: '草稿',
  active: '可用',
  deprecated: '已废弃',
  archived: '已归档',
};

function formatRelativeTime(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days} 天前`;
  const months = Math.floor(days / 30);
  if (months < 12) return `${months} 个月前`;
  return `${Math.floor(months / 12)} 年前`;
}

function useDebouncedValue<T>(value: T, delay: number): T {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  return debounced;
}

export function PromptListPage() {
  const navigate = useNavigate();
  const [keyword, setKeyword] = useState('');
  const [selectedStatuses, setSelectedStatuses] = useState<PromptStatus[]>([]);
  const [showArchived, setShowArchived] = useState(false);
  const [showFilters, setShowFilters] = useState(false);
  const [page, setPage] = useState(1);

  const debouncedKeyword = useDebouncedValue(keyword, 300);

  const params: PromptListParams = useMemo(() => {
    const statuses = showArchived
      ? selectedStatuses
      : selectedStatuses.filter((s) => s !== 'archived');
    return {
      keyword: debouncedKeyword || undefined,
      status: statuses.length > 0 ? statuses : undefined,
      page,
      pageSize: PAGE_SIZE,
    };
  }, [debouncedKeyword, selectedStatuses, showArchived, page]);

  const {
    data,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['prompts', params],
    queryFn: () => getPrompts(params),
    placeholderData: (prev) => prev,
  });

  const prompts = data?.items ?? [];
  const hasFilters = keyword || selectedStatuses.length > 0;
  const totalPages = data ? Math.ceil(data.total / data.pageSize) : 1;

  const toggleStatus = useCallback((status: PromptStatus) => {
    setSelectedStatuses((prev) =>
      prev.includes(status) ? prev.filter((s) => s !== status) : [...prev, status],
    );
    setPage(1);
  }, []);

  const clearFilters = useCallback(() => {
    setKeyword('');
    setSelectedStatuses([]);
    setPage(1);
  }, []);

  useEffect(() => {
    setPage(1);
  }, [debouncedKeyword]);

  const renderEmptyState = () => {
    if (hasFilters) {
      return (
        <div className="border border-dashed border-border-strong rounded-xl py-16 px-8 text-center bg-bg-subtle">
          <Search size={24} className="mx-auto text-text-subtle mb-3" />
          <div className="text-base font-semibold mb-1">没有匹配的提示词</div>
          <p className="text-sm text-text-muted mb-5">试试其他关键词或筛选条件</p>
          <button
            onClick={clearFilters}
            className="h-9 px-4 rounded-lg border border-border text-sm cursor-pointer bg-transparent text-text hover:border-border-strong hover:bg-bg-subtle transition-all duration-150"
          >
            清除筛选
          </button>
        </div>
      );
    }

    return (
      <div className="border border-dashed border-border-strong rounded-xl py-16 px-8 text-center bg-bg-subtle">
        <div className="inline-flex p-4 bg-bg-elevated border border-border rounded-xl text-accent mb-4">
          <FileText size={24} />
        </div>
        <div className="text-base font-semibold mb-1">还没有提示词</div>
        <p className="text-sm text-text-muted mb-5">创建你的第一个提示词，开始构建可复用的提示词资产。</p>
        <button
          onClick={() => navigate('/prompts/new')}
          className="h-10 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on transition-all duration-150 inline-flex items-center gap-2"
        >
          <Plus size={14} />
          创建提示词
        </button>
      </div>
    );
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="text-[13px] text-text-muted mb-1 font-mono">$ prompts list</div>
          <h1 className="text-[28px] font-semibold tracking-tight">提示词</h1>
        </div>
        <button
          onClick={() => navigate('/prompts/new')}
          className="h-10 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on transition-all duration-150 inline-flex items-center gap-2 font-medium"
        >
          <Plus size={14} />
          创建提示词
        </button>
      </div>

      <div className="flex items-center gap-3 mb-5">
        <div className="flex-1 relative">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-subtle" />
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="搜索标题、描述或正文…"
            className="w-full h-10 pl-9 pr-3 bg-bg-elevated border border-border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)]"
          />
          {keyword && (
            <button
              onClick={() => setKeyword('')}
              className="absolute right-2 top-1/2 -translate-y-1/2 w-6 h-6 grid place-items-center text-text-subtle hover:text-text cursor-pointer bg-transparent border-none rounded transition-colors"
            >
              <X size={12} />
            </button>
          )}
        </div>
        <button
          onClick={() => setShowFilters(!showFilters)}
          className={`h-10 px-3 rounded-lg border text-sm cursor-pointer transition-all duration-150 inline-flex items-center gap-2 ${
            showFilters
              ? 'border-accent bg-accent-soft text-accent'
              : 'border-border bg-bg-elevated text-text-muted hover:border-border-strong hover:text-text'
          }`}
        >
          <SlidersHorizontal size={14} />
          筛选
        </button>
      </div>

      {showFilters && (
        <div className="flex flex-wrap items-center gap-3 mb-5 p-4 bg-bg-elevated border border-border rounded-lg">
          <div className="text-xs font-medium text-text-muted mr-1">状态</div>
          {STATUS_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              onClick={() => {
                if (opt.value === 'archived') {
                  setShowArchived(!showArchived);
                  if (!showArchived) toggleStatus('archived');
                } else {
                  toggleStatus(opt.value);
                }
              }}
              className={`h-7 px-2.5 rounded-md border text-xs cursor-pointer transition-all duration-150 ${
                opt.value === 'archived'
                  ? showArchived
                    ? 'border-accent bg-accent-soft text-accent'
                    : 'border-border bg-bg text-text-muted hover:border-border-strong'
                  : selectedStatuses.includes(opt.value)
                    ? 'border-accent bg-accent-soft text-accent'
                    : 'border-border bg-bg text-text-muted hover:border-border-strong'
              }`}
            >
              {opt.label}
            </button>
          ))}
          {hasFilters && (
            <button
              onClick={clearFilters}
              className="h-7 px-2.5 rounded-md border border-border text-xs cursor-pointer text-text-muted hover:text-text transition-colors bg-transparent"
            >
              清除全部
            </button>
          )}
        </div>
      )}

      {isLoading && (
        <div className="flex items-center justify-center py-20 text-text-muted text-sm">
          <Loader2 size={18} className="animate-spin mr-2" />
          加载中…
        </div>
      )}

      {isError && (
        <div className="flex items-center justify-center py-20 text-danger text-sm">
          <AlertCircle size={18} className="mr-2 shrink-0" />
          {(error as Error)?.message || '加载失败，请重试'}
        </div>
      )}

      {!isLoading && !isError && prompts.length === 0 && renderEmptyState()}

      {!isLoading && !isError && prompts.length > 0 && (
        <div className="border border-border rounded-xl overflow-hidden bg-bg-elevated">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-bg-subtle text-left">
                <th className="px-4 py-3 font-medium text-text-muted text-xs">标题</th>
                <th className="px-4 py-3 font-medium text-text-muted text-xs hidden md:table-cell">状态</th>
                <th className="px-4 py-3 font-medium text-text-muted text-xs hidden lg:table-cell">标签</th>
                <th className="px-4 py-3 font-medium text-text-muted text-xs hidden xl:table-cell">模型</th>
                <th className="px-4 py-3 font-medium text-text-muted text-xs text-right">更新时间</th>
              </tr>
            </thead>
            <tbody>
              {prompts.map((prompt: Prompt) => (
                <tr
                  key={prompt.id}
                  onClick={() => navigate(`/prompts/${prompt.id}`)}
                  className="border-b border-border last:border-b-0 cursor-pointer transition-colors duration-100 hover:bg-bg-subtle"
                >
                  <td className="px-4 py-3">
                    <div className="font-medium text-text truncate max-w-[320px]">{prompt.title}</div>
                    {prompt.description && (
                      <div className="text-text-muted text-xs mt-0.5 truncate max-w-[320px]">
                        {prompt.description}
                      </div>
                    )}
                  </td>
                  <td className="px-4 py-3 hidden md:table-cell">
                    <span
                      className={`inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium border ${STATUS_BADGE[prompt.status]}`}
                    >
                      {STATUS_LABEL[prompt.status]}
                    </span>
                  </td>
                  <td className="px-4 py-3 hidden lg:table-cell">
                    <div className="flex flex-wrap gap-1">
                      {prompt.tags?.map((tag) => (
                        <span
                          key={tag}
                          className="inline-flex px-1.5 py-0.5 rounded bg-bg-subtle border border-border text-[11px] text-text-muted font-mono"
                        >
                          {tag}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td className="px-4 py-3 hidden xl:table-cell">
                    {prompt.targetModel ? (
                      <span className="text-xs text-text-muted font-mono">{prompt.targetModel}</span>
                    ) : (
                      <span className="text-xs text-text-subtle">—</span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <span className="text-xs text-text-subtle whitespace-nowrap">
                      {formatRelativeTime(prompt.updatedAt)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {data && data.total > 0 && (
            <div className="px-4 py-3 border-t border-border text-xs text-text-muted flex items-center justify-between">
              <span>
                共 {data.total} 个提示词
              </span>
              <div className="flex items-center gap-1">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page <= 1}
                  className="w-7 h-7 rounded-md border border-border grid place-items-center cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong disabled:opacity-30 disabled:cursor-default transition-all duration-150"
                >
                  <ChevronLeft size={14} />
                </button>
                {totalPages > 7 ? (
                  <PageNumbers page={page} total={totalPages} onChange={setPage} />
                ) : (
                  Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
                    <button
                      key={p}
                      onClick={() => setPage(p)}
                      className={`w-7 h-7 rounded-md text-xs grid place-items-center cursor-pointer transition-all duration-150 ${
                        p === page
                          ? 'bg-accent text-accent-on font-medium'
                          : 'bg-transparent text-text-muted hover:text-text'
                      }`}
                    >
                      {p}
                    </button>
                  ))
                )}
                <button
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="w-7 h-7 rounded-md border border-border grid place-items-center cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong disabled:opacity-30 disabled:cursor-default transition-all duration-150"
                >
                  <ChevronRight size={14} />
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function PageNumbers({ page, total, onChange }: { page: number; total: number; onChange: (p: number) => void }) {
  const pages: (number | '...')[] = [];
  if (page > 3) pages.push(1, '...');
  for (let i = Math.max(1, page - 1); i <= Math.min(total, page + 1); i++) {
    pages.push(i);
  }
  if (page < total - 2) pages.push('...', total);

  return (
    <>
      {pages.map((p, i) =>
        p === '...' ? (
          <span key={`ellipsis-${i}`} className="w-7 h-7 grid place-items-center text-text-subtle text-xs">
            …
          </span>
        ) : (
          <button
            key={p}
            onClick={() => onChange(p)}
            className={`w-7 h-7 rounded-md text-xs grid place-items-center cursor-pointer transition-all duration-150 ${
              p === page
                ? 'bg-accent text-accent-on font-medium'
                : 'bg-transparent text-text-muted hover:text-text'
            }`}
          >
            {p}
          </button>
        ),
      )}
    </>
  );
}
