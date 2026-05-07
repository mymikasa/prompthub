import { useState, useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { z } from 'zod/v4';
import { zodResolver } from '@hookform/resolvers/zod';
import { Save, Loader2, Plus, X } from 'lucide-react';
import { updatePrompt } from '../api';
import type { Prompt } from '../types';

const editorSchema = z.object({
  title: z.string().min(1, '标题不能为空'),
  description: z.string().optional(),
  messageFormat: z.enum(['single_text', 'chat_messages']),
  body: z.string().optional(),
  visibility: z.enum(['private', 'workspace']),
  status: z.enum(['draft', 'active', 'deprecated']),
  targetProvider: z.string().optional(),
  targetModel: z.string().optional(),
  defaultTemperature: z.number().nullable().optional(),
  defaultMaxTokens: z.number().nullable().optional(),
  usageNotes: z.string().optional(),
});

type EditorForm = z.infer<typeof editorSchema>;

function extractVariables(text: string): string[] {
  const matches = text.match(/\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g) || [];
  return [...new Set(matches.map((m) => m.slice(2, -2)))];
}

export function EditorTab({ prompt }: { prompt: Prompt }) {
  const queryClient = useQueryClient();
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [tags, setTags] = useState<string[]>(prompt.tags || []);
  const [tagInput, setTagInput] = useState('');

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<EditorForm>({
    resolver: zodResolver(editorSchema),
    defaultValues: {
      title: prompt.title,
      description: prompt.description || '',
      messageFormat: prompt.messageFormat,
      body: prompt.body || '',
      visibility: prompt.visibility,
      status: prompt.status === 'archived' ? 'draft' : prompt.status,
      targetProvider: prompt.targetProvider || '',
      targetModel: prompt.targetModel || '',
      defaultTemperature: prompt.defaultTemperature,
      defaultMaxTokens: prompt.defaultMaxTokens,
      usageNotes: prompt.usageNotes || '',
    },
  });

  const body = watch('body') || '';
  const inferredVars = extractVariables(body);

  const onSubmit = async (data: EditorForm) => {
    setSaving(true);
    setSaveError('');
    setSaveSuccess(false);
    try {
      await updatePrompt(prompt.id, { ...data, tags });
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 2000);
      queryClient.invalidateQueries({ queryKey: ['prompt', prompt.id] });
      queryClient.invalidateQueries({ queryKey: ['prompts'] });
    } catch (err) {
      setSaveError((err as Error)?.message || '保存失败');
    } finally {
      setSaving(false);
    }
  };

  const addTag = () => {
    const t = tagInput.trim();
    if (t && !tags.includes(t)) setTags([...tags, t]);
    setTagInput('');
  };

  const removeTag = (tag: string) => setTags(tags.filter((t) => t !== tag));

  useEffect(() => {
    setTags(prompt.tags || []);
  }, [prompt]);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
        <div className="space-y-4">
          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">标题</label>
            <input
              className={`w-full h-9 px-3 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${errors.title ? 'border-danger' : 'border-border'}`}
              {...register('title')}
            />
            {errors.title && <span className="text-xs text-danger mt-1 block">{errors.title.message}</span>}
          </div>

          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">描述</label>
            <textarea
              className="w-full px-3 py-2 bg-bg-elevated border border-border rounded-lg text-sm text-text outline-none transition-all duration-150 hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] resize-y min-h-[60px]"
              rows={2}
              {...register('description')}
            />
          </div>

          <div>
            <div className="flex items-center gap-3 mb-2">
              <label className="text-xs font-medium text-text-muted">消息格式</label>
              <select
                className="h-8 px-2 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none cursor-pointer"
                {...register('messageFormat')}
              >
                <option value="single_text">Single Text</option>
                <option value="chat_messages">Chat Messages</option>
              </select>
            </div>

            <textarea
              className="w-full px-3 py-2 bg-code-bg text-code-text border border-[rgba(255,255,255,0.08)] rounded-lg text-sm font-mono outline-none resize-y min-h-[200px] focus:border-accent"
              rows={12}
              placeholder="输入提示词正文，使用 {{variable}} 作为变量占位符…"
              {...register('body')}
            />
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">标签</label>
            <div className="flex flex-wrap gap-1.5 mb-2">
              {tags.map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-bg-subtle border border-border text-xs text-text-muted font-mono"
                >
                  {tag}
                  <button
                    type="button"
                    onClick={() => removeTag(tag)}
                    className="text-text-subtle hover:text-danger cursor-pointer bg-transparent border-none p-0"
                  >
                    <X size={10} />
                  </button>
                </span>
              ))}
            </div>
            <div className="flex gap-1.5">
              <input
                className="flex-1 h-8 px-2.5 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none hover:border-border-strong focus:border-accent"
                placeholder="添加标签"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
              />
              <button
                type="button"
                onClick={addTag}
                className="h-8 w-8 rounded-md border border-border text-text-muted cursor-pointer bg-transparent hover:text-text hover:border-border-strong grid place-items-center"
              >
                <Plus size={12} />
              </button>
            </div>
          </div>

          <SelectField label="可见性" {...register('visibility')}>
            <option value="private">私有</option>
            <option value="workspace">工作空间</option>
          </SelectField>

          <SelectField label="状态" {...register('status')}>
            <option value="draft">草稿</option>
            <option value="active">可用</option>
            <option value="deprecated">已废弃</option>
          </SelectField>

          <InputField label="模型提供方" {...register('targetProvider')} placeholder="openai" />
          <InputField label="模型名称" {...register('targetModel')} placeholder="gpt-4o" />
          <InputField label="Temperature" {...register('defaultTemperature')} placeholder="0.7" type="number" step="0.1" />
          <InputField label="Max Tokens" {...register('defaultMaxTokens')} placeholder="2048" type="number" />

          <div>
            <label className="text-xs font-medium text-text-muted block mb-1.5">使用说明</label>
            <textarea
              className="w-full px-2.5 py-2 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none resize-y min-h-[60px] hover:border-border-strong focus:border-accent"
              rows={3}
              {...register('usageNotes')}
            />
          </div>

          {inferredVars.length > 0 && (
            <div>
              <div className="text-xs font-medium text-text-muted mb-1.5">推断变量</div>
              <div className="flex flex-wrap gap-1">
                {inferredVars.map((v) => (
                  <span
                    key={v}
                    className="px-1.5 py-0.5 rounded bg-[rgba(251,191,36,0.12)] border border-[rgba(251,191,36,0.25)] text-[11px] text-code-var font-mono"
                  >
                    {`{{${v}}}`}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3 pt-2">
        <button
          type="submit"
          disabled={saving}
          className="h-9 px-4 rounded-lg border border-transparent text-sm cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-2 font-medium"
        >
          {saving ? <Loader2 size={14} className="animate-spin" /> : <Save size={14} />}
          保存
        </button>
        {saveSuccess && <span className="text-sm text-success">已保存</span>}
        {saveError && <span className="text-sm text-danger">{saveError}</span>}
      </div>
    </form>
  );
}

function InputField({
  label,
  placeholder,
  type = 'text',
  step,
  ...registerProps
}: {
  label: string;
  placeholder?: string;
  type?: string;
  step?: string;
} & React.ComponentProps<'input'>) {
  return (
    <div>
      <label className="text-xs font-medium text-text-muted block mb-1.5">{label}</label>
      <input
        type={type}
        step={step}
        className="w-full h-8 px-2.5 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none hover:border-border-strong focus:border-accent"
        placeholder={placeholder}
        {...registerProps}
      />
    </div>
  );
}

function SelectField({
  label,
  children,
  ...registerProps
}: {
  label: string;
  children: React.ReactNode;
} & React.ComponentProps<'select'>) {
  return (
    <div>
      <label className="text-xs font-medium text-text-muted block mb-1.5">{label}</label>
      <select
        className="w-full h-8 px-2.5 bg-bg-elevated border border-border rounded-md text-xs text-text outline-none cursor-pointer hover:border-border-strong focus:border-accent"
        {...registerProps}
      >
        {children}
      </select>
    </div>
  );
}
