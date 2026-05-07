import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Loader2, AlertCircle, RotateCcw } from 'lucide-react';
import { getVersion, restoreVersion } from '@/features/prompts/api';

function formatDateTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function highlightVariables(text: string) {
  const parts = text.split(/(\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\})/g);
  return parts.map((part, i) => {
    if (/^\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}$/.test(part)) {
      return (
        <span
          key={i}
          className="text-code-var bg-[rgba(251,191,36,0.12)] px-1 rounded border border-[rgba(251,191,36,0.25)]"
        >
          {part}
        </span>
      );
    }
    return <span key={i}>{part}</span>;
  });
}

export function VersionDetailPage() {
  const { promptId, versionId } = useParams<{ promptId: string; versionId: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: version, isLoading, isError, error } = useQuery({
    queryKey: ['promptVersion', promptId, versionId],
    queryFn: () => getVersion(promptId!, versionId!),
    enabled: !!promptId && !!versionId,
  });

  const handleRestore = async () => {
    if (!version || !window.confirm(`确定要恢复 v${version.versionNumber} 吗？这将基于此版本创建一个新版本。`)) return;
    await restoreVersion(promptId!, versionId!);
    queryClient.invalidateQueries({ queryKey: ['promptVersions', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompt', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompts'] });
    navigate(`/prompts/${promptId}`, { replace: true });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20 text-text-muted text-sm">
        <Loader2 size={18} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (isError || !version) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-text-muted text-sm">
        <AlertCircle size={24} className="text-danger mb-3" />
        <p className="text-text mb-4">{(error as Error)?.message || '加载版本失败'}</p>
        <button
          onClick={() => navigate(-1)}
          className="h-9 px-4 rounded-lg border border-border text-sm cursor-pointer bg-transparent text-text hover:border-border-strong transition-all duration-150"
        >
          返回
        </button>
      </div>
    );
  }

  const snap = version.snapshot;

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <button
          onClick={() => navigate(`/prompts/${promptId}`)}
          className="w-8 h-8 rounded-lg border border-border bg-bg-elevated text-text-muted grid place-items-center cursor-pointer transition-all duration-150 hover:text-text hover:border-border-strong"
        >
          <ArrowLeft size={14} />
        </button>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2.5">
            <h1 className="text-xl font-semibold tracking-tight">
              版本快照 v{version.versionNumber}
            </h1>
          </div>
          <div className="text-sm text-text-muted mt-0.5">
            {version.author} · {formatDateTime(version.createdAt)}
          </div>
        </div>
        <button
          onClick={handleRestore}
          className="h-8 px-3 rounded-lg border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-accent hover:border-accent transition-all duration-150 inline-flex items-center gap-1.5"
        >
          <RotateCcw size={12} />
          恢复此版本
        </button>
      </div>

      <div className="space-y-5">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <MetaItem label="版本号" value={`v${version.versionNumber}`} mono />
          <MetaItem label="状态" value={snap?.status || '—'} />
          <MetaItem label="目标模型" value={snap?.targetModel || '—'} mono />
          <MetaItem label="提供方" value={snap?.targetProvider || '—'} />
        </div>

        {snap?.tags?.length > 0 && (
          <div>
            <div className="text-xs font-medium text-text-muted mb-2">标签</div>
            <div className="flex flex-wrap gap-1.5">
              {snap.tags.map((tag) => (
                <span key={tag} className="px-2 py-0.5 rounded bg-bg-subtle border border-border text-xs text-text-muted font-mono">
                  {tag}
                </span>
              ))}
            </div>
          </div>
        )}

        {version.changeNote && (
          <div>
            <div className="text-xs font-medium text-text-muted mb-2">变更说明</div>
            <div className="text-sm text-text leading-relaxed whitespace-pre-wrap bg-bg-subtle border border-border rounded-lg p-4">
              {version.changeNote}
            </div>
          </div>
        )}

        <div>
          <div className="text-xs font-medium text-text-muted mb-2">提示词正文快照</div>
          <div className="bg-code-bg rounded-lg p-4 px-5 font-mono text-[13px] leading-[1.7] text-code-text whitespace-pre-wrap max-h-[500px] overflow-auto">
            {snap?.content
              ? highlightVariables(snap.content)
              : <span className="text-code-comment">无内容</span>
            }
          </div>
        </div>

        {snap?.variables?.length > 0 && (
          <div>
            <div className="text-xs font-medium text-text-muted mb-2">变量快照</div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
              {snap.variables.map((v) => (
                <div key={v.name} className="flex items-center gap-2 px-3 py-2 bg-bg-subtle border border-border rounded-lg">
                  <span className="font-mono text-xs text-text-subtle">
                    {'{{'}{v.name}{'}}'}
                  </span>
                  <span className="text-xs text-text-muted flex-1 truncate">
                    {v.label || v.name}
                  </span>
                  {v.required && (
                    <span className="text-[10px] text-danger">必填</span>
                  )}
                  {v.defaultValue && (
                    <span className="text-[10px] text-text-subtle truncate max-w-[120px]">
                      默认: {v.defaultValue}
                    </span>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function MetaItem({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div>
      <div className="text-[11px] text-text-subtle mb-0.5">{label}</div>
      <div className={`text-sm text-text ${mono ? 'font-mono' : ''}`}>{value}</div>
    </div>
  );
}
