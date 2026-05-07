import { useState, useMemo } from 'react';
import { Copy } from 'lucide-react';
import type { Prompt } from '../types';

function extractVariables(text: string): string[] {
  const matches = text.match(/\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g) || [];
  return [...new Set(matches.map((m) => m.slice(2, -2)))];
}

export function VariablesTab({ prompt }: { prompt: Prompt }) {
  const [values, setValues] = useState<Record<string, string>>({});
  const [copied, setCopied] = useState(false);

  const variables = useMemo(() => extractVariables(prompt.body || ''), [prompt.body]);

  const rendered = useMemo(() => {
    return (prompt.body || '').replace(
      /\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g,
      (_, name) => values[name] ?? `{{${name}}}`,
    );
  }, [prompt.body, values]);

  const unfilled = variables.filter((v) => !values[v]?.trim());
  const hasContent = rendered.trim().length > 0;

  const handleCopy = async () => {
    await navigator.clipboard.writeText(rendered);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  if (variables.length === 0) {
    return (
      <div className="py-12 text-center text-text-muted text-sm">
        没有检测到变量占位符。在正文中使用 {'{{variable_name}}'} 添加变量。
      </div>
    );
  }

  return (
    <div className="space-y-6">
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
          {hasContent ? rendered : <span className="text-code-comment">填写变量后查看渲染结果</span>}
        </div>

        {unfilled.length > 0 && hasContent && (
          <div className="mt-2 text-xs text-amber-600">
            仍有未填写的变量：{unfilled.map((v) => `{{${v}}}`).join(', ')}
          </div>
        )}
      </div>
    </div>
  );
}
