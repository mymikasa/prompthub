import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import {
  ArrowLeft,
  Eye,
  Pencil,
  Variable,
  FlaskConical,
  Activity,
  History,
  Archive,
  RotateCcw,
  Loader2,
  AlertCircle,
} from 'lucide-react';
import { getPrompt, archivePrompt, restorePrompt } from '@/features/prompts/api';
import type { PromptStatus } from '@/features/prompts/types';
import { OverviewTab } from './tabs/overview-tab';
import { EditorTab } from './tabs/editor-tab';
import { VariablesTab } from './tabs/variables-tab';
import { VersionsTab } from './tabs/versions-tab';
import { TestCasesTab } from './tabs/test-cases-tab';
import { RunsTab } from './tabs/runs-tab';

const STATUS_BADGE: Record<PromptStatus, string> = {
  draft: 'bg-amber-50 text-amber-700 border-amber-200',
  active: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  deprecated: 'bg-gray-50 text-gray-600 border-gray-200',
  archived: 'bg-stone-50 text-stone-500 border-stone-200',
};

const STATUS_LABEL: Record<PromptStatus, string> = {
  draft: '草稿',
  active: '可用',
  deprecated: '已废弃',
  archived: '已归档',
};

type TabKey = 'overview' | 'editor' | 'variables' | 'test-cases' | 'runs' | 'versions';

const TABS: { key: TabKey; label: string; icon: typeof Eye }[] = [
  { key: 'overview', label: '概览', icon: Eye },
  { key: 'editor', label: '编辑器', icon: Pencil },
  { key: 'variables', label: '变量', icon: Variable },
  { key: 'test-cases', label: '测试用例', icon: FlaskConical },
  { key: 'runs', label: '运行记录', icon: Activity },
  { key: 'versions', label: '版本', icon: History },
];

export function PromptDetailPage() {
  const { promptId } = useParams<{ promptId: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<TabKey>('overview');

  const {
    data: prompt,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['prompt', promptId],
    queryFn: () => getPrompt(promptId!),
    enabled: !!promptId,
  });

  const handleArchive = async () => {
    if (!promptId || !prompt) return;
    if (!window.confirm(`确定要归档「${prompt.title}」吗？`)) return;
    await archivePrompt(promptId);
    queryClient.invalidateQueries({ queryKey: ['prompt', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompts'] });
  };

  const handleRestore = async () => {
    if (!promptId) return;
    await restorePrompt(promptId);
    queryClient.invalidateQueries({ queryKey: ['prompt', promptId] });
    queryClient.invalidateQueries({ queryKey: ['prompts'] });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20 text-text-muted text-sm">
        <Loader2 size={18} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  if (isError || !prompt) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-text-muted text-sm">
        <AlertCircle size={24} className="text-danger mb-3" />
        <p className="text-text mb-4">{(error as Error)?.message || '加载提示词失败'}</p>
        <button
          onClick={() => navigate('/prompts')}
          className="h-9 px-4 rounded-lg border border-border text-sm cursor-pointer bg-transparent text-text hover:border-border-strong transition-all duration-150"
        >
          返回列表
        </button>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <button
          onClick={() => navigate('/prompts')}
          className="w-8 h-8 rounded-lg border border-border bg-bg-elevated text-text-muted grid place-items-center cursor-pointer transition-all duration-150 hover:text-text hover:border-border-strong"
        >
          <ArrowLeft size={14} />
        </button>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2.5">
            <h1 className="text-xl font-semibold tracking-tight truncate">{prompt.title}</h1>
            <span
              className={`inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium border shrink-0 ${STATUS_BADGE[prompt.status]}`}
            >
              {STATUS_LABEL[prompt.status]}
            </span>
          </div>
          {prompt.description && (
            <p className="text-sm text-text-muted mt-0.5 truncate">{prompt.description}</p>
          )}
        </div>
        <div className="flex items-center gap-2 shrink-0">
          {prompt.status !== 'archived' ? (
            <button
              onClick={handleArchive}
              className="h-8 px-3 rounded-lg border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150 inline-flex items-center gap-1.5"
            >
              <Archive size={12} />
              归档
            </button>
          ) : (
            <button
              onClick={handleRestore}
              className="h-8 px-3 rounded-lg border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150 inline-flex items-center gap-1.5"
            >
              <RotateCcw size={12} />
              恢复
            </button>
          )}
        </div>
      </div>

      <div className="flex gap-0 border-b border-border mb-6">
        {TABS.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2.5 text-sm cursor-pointer border-b-[1.5px] mb-[-1px] flex items-center gap-1.5 transition-colors duration-150 bg-transparent ${
              activeTab === tab.key
                ? 'text-text border-b-accent font-medium'
                : 'text-text-muted border-b-transparent hover:text-text'
            }`}
          >
            <tab.icon size={14} />
            {tab.label}
          </button>
        ))}
      </div>

      <div>
        {activeTab === 'overview' && <OverviewTab prompt={prompt} />}
        {activeTab === 'editor' && <EditorTab prompt={prompt} />}
        {activeTab === 'variables' && <VariablesTab prompt={prompt} />}
        {activeTab === 'test-cases' && <TestCasesTab prompt={prompt} />}
        {activeTab === 'runs' && <RunsTab promptId={prompt.id} />}
        {activeTab === 'versions' && <VersionsTab promptId={prompt.id} />}
      </div>
    </div>
  );
}
