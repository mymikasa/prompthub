import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { Loader2, AlertCircle, RotateCcw, Eye } from 'lucide-react';
import { getVersions, restoreVersion } from '../api';

function formatDateTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

const STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  active: '可用',
  deprecated: '已废弃',
  archived: '已归档',
};

const STATUS_BADGE: Record<string, string> = {
  draft: 'bg-amber-50 text-amber-700 border-amber-200',
  active: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  deprecated: 'bg-gray-50 text-gray-600 border-gray-200',
  archived: 'bg-stone-50 text-stone-500 border-stone-200',
};

export function VersionsTab({ promptId }: { promptId: number | string }) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: versions, isLoading, isError } = useQuery({
    queryKey: ['promptVersions', promptId],
    queryFn: () => getVersions(promptId),
  });

  const handleRestore = async (versionId: number | string, versionNumber: number) => {
    if (!window.confirm(`确定要恢复 v${versionNumber} 吗？这将基于此版本创建一个新版本。`)) return;
    await restoreVersion(promptId, versionId);
    queryClient.invalidateQueries({ queryKey: ['promptVersions', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompt', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompts'] });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-text-muted text-sm">
        <Loader2 size={16} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (isError || !versions?.length) {
    return (
      <div className="py-12 text-center text-text-muted text-sm">
        <AlertCircle size={20} className="mx-auto mb-2 text-text-subtle" />
        暂无版本记录
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {versions.map((v) => (
        <div
          key={v.id}
          className="flex items-center gap-4 px-4 py-3 border border-border rounded-lg bg-bg-elevated hover:bg-bg-subtle transition-colors duration-100"
        >
          <div className="font-mono text-sm font-semibold text-accent shrink-0 w-10">
            v{v.versionNumber}
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 text-sm">
              <span className="text-text-muted">{v.author}</span>
              <span className="text-text-subtle text-xs">{formatDateTime(v.createdAt)}</span>
              {v.snapshot?.status && (
                <span className={`inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium border ${STATUS_BADGE[v.snapshot.status] || ''}`}>
                  {STATUS_LABEL[v.snapshot.status] || v.snapshot.status}
                </span>
              )}
            </div>
            {v.changeNote && (
              <div className="text-xs text-text-subtle mt-0.5 truncate">{v.changeNote}</div>
            )}
            {v.snapshot?.tags?.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-1">
                {v.snapshot.tags.map((tag) => (
                  <span key={tag} className="px-1.5 py-0.5 rounded bg-bg-subtle border border-border text-[10px] text-text-muted font-mono">
                    {tag}
                  </span>
                ))}
              </div>
            )}
          </div>

          <div className="flex items-center gap-2 shrink-0">
            <button
              onClick={() => navigate(`/prompts/${promptId}/versions/${v.id}`)}
              className="h-7 px-2.5 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150 inline-flex items-center gap-1"
            >
              <Eye size={11} />
              快照
            </button>
            <button
              onClick={() => handleRestore(v.id, v.versionNumber)}
              className="h-7 px-2.5 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-accent hover:border-accent transition-all duration-150 inline-flex items-center gap-1"
            >
              <RotateCcw size={11} />
              恢复
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
