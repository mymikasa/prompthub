import { Copy, FileText } from 'lucide-react';
import { useState } from 'react';
import type { Prompt } from '../types';

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

export function OverviewTab({ prompt }: { prompt: Prompt }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(prompt.body || '');
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <MetaItem label="格式" value={prompt.messageFormat === 'chat_messages' ? 'Chat Messages' : 'Single Text'} />
        <MetaItem label="可见性" value={prompt.visibility === 'workspace' ? '工作空间' : '私有'} />
        <MetaItem label="目标模型" value={prompt.targetModel || '—'} mono />
        <MetaItem label="提供方" value={prompt.targetProvider || '—'} />
      </div>

      {prompt.tags?.length > 0 && (
        <div>
          <div className="text-xs font-medium text-text-muted mb-2">标签</div>
          <div className="flex flex-wrap gap-1.5">
            {prompt.tags.map((tag) => (
              <span
                key={tag}
                className="px-2 py-0.5 rounded bg-bg-subtle border border-border text-xs text-text-muted font-mono"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      )}

      {(prompt.defaultTemperature !== null || prompt.defaultMaxTokens !== null) && (
        <div className="grid grid-cols-2 gap-4">
          {prompt.defaultTemperature !== null && (
            <MetaItem label="Temperature" value={String(prompt.defaultTemperature)} mono />
          )}
          {prompt.defaultMaxTokens !== null && (
            <MetaItem label="Max Tokens" value={String(prompt.defaultMaxTokens)} mono />
          )}
        </div>
      )}

      {prompt.usageNotes && (
        <div>
          <div className="text-xs font-medium text-text-muted mb-2">使用说明</div>
          <div className="text-sm text-text leading-relaxed whitespace-pre-wrap bg-bg-subtle border border-border rounded-lg p-4">
            {prompt.usageNotes}
          </div>
        </div>
      )}

      <div>
        <div className="flex items-center justify-between mb-2">
          <div className="text-xs font-medium text-text-muted flex items-center gap-1.5">
            <FileText size={12} />
            提示词正文
          </div>
          <button
            onClick={handleCopy}
            className="h-7 px-2.5 rounded-md border border-border text-xs cursor-pointer bg-transparent text-text-muted hover:text-text hover:border-border-strong transition-all duration-150 inline-flex items-center gap-1.5"
          >
            <Copy size={11} />
            {copied ? '已复制' : '复制'}
          </button>
        </div>

        <div className="bg-code-bg rounded-lg p-4 px-5 font-mono text-[13px] leading-[1.7] text-code-text whitespace-pre-wrap max-h-[400px] overflow-auto">
          {highlightVariables(prompt.body || '')}
        </div>
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
