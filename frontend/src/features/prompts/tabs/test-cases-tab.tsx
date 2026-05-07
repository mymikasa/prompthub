import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Plus, Pencil, Trash2, Loader2, X, Save } from 'lucide-react';
import { getTestCases, createTestCase, updateTestCase, deleteTestCase } from '../api';
import type { Prompt, TestCase } from '../types';

export function TestCasesTab({ prompt }: { prompt: Prompt }) {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formName, setFormName] = useState('');
  const [formExpected, setFormExpected] = useState('');
  const [formVariables, setFormVariables] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);

  const { data: testCases, isLoading } = useQuery({
    queryKey: ['testCases', prompt.id],
    queryFn: () => getTestCases(prompt.id),
  });

  const handleEdit = (tc: TestCase) => {
    setEditingId(tc.id);
    setFormName(tc.name);
    setFormExpected(tc.expectedBehavior || '');
    setFormVariables(tc.variableValues || {});
    setShowForm(true);
  };

  const handleNew = () => {
    setEditingId(null);
    setFormName('');
    setFormExpected('');
    setFormVariables({});
    setShowForm(true);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (editingId) {
        await updateTestCase(prompt.id, editingId, {
          name: formName,
          expectedBehavior: formExpected,
          variableValues: formVariables,
        });
      } else {
        await createTestCase(prompt.id, {
          name: formName,
          expectedBehavior: formExpected,
          variableValues: formVariables,
        });
      }
      queryClient.invalidateQueries({ queryKey: ['testCases', prompt.id] });
      setShowForm(false);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('确定删除此测试用例？')) return;
    await deleteTestCase(prompt.id, id);
    queryClient.invalidateQueries({ queryKey: ['testCases', prompt.id] });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-text-muted text-sm">
        <Loader2 size={16} className="animate-spin mr-2" />
        加载中…
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="text-xs font-medium text-text-muted">
          {testCases?.length || 0} 个测试用例
        </div>
        <button
          onClick={handleNew}
          className="h-8 px-3 rounded-lg border border-transparent text-xs cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on transition-all duration-150 inline-flex items-center gap-1.5 font-medium"
        >
          <Plus size={12} />
          新增
        </button>
      </div>

      {showForm && (
        <div className="border border-accent/30 bg-accent-soft/30 rounded-lg p-4 space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">{editingId ? '编辑测试用例' : '新建测试用例'}</span>
            <button onClick={() => setShowForm(false)} className="w-6 h-6 grid place-items-center text-text-muted hover:text-text cursor-pointer bg-transparent border-none">
              <X size={14} />
            </button>
          </div>
          <input
            className="w-full h-9 px-3 bg-bg-elevated border border-border rounded-lg text-sm text-text outline-none hover:border-border-strong focus:border-accent"
            placeholder="用例名称"
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
          />
          <textarea
            className="w-full px-3 py-2 bg-bg-elevated border border-border rounded-lg text-sm text-text outline-none resize-y min-h-[60px] hover:border-border-strong focus:border-accent"
            placeholder="预期行为说明"
            value={formExpected}
            onChange={(e) => setFormExpected(e.target.value)}
          />
          <div className="grid grid-cols-2 gap-2">
            {Object.entries(formVariables).map(([key, val]) => (
              <div key={key} className="flex items-center gap-1.5">
                <span className="text-xs font-mono text-text-subtle shrink-0">{key}:</span>
                <input
                  className="flex-1 h-7 px-2 bg-bg-elevated border border-border rounded text-xs text-text outline-none focus:border-accent"
                  value={val}
                  onChange={(e) => setFormVariables({ ...formVariables, [key]: e.target.value })}
                />
                <button onClick={() => { const nv = { ...formVariables }; delete nv[key]; setFormVariables(nv); }} className="text-text-subtle hover:text-danger cursor-pointer bg-transparent border-none">
                  <X size={10} />
                </button>
              </div>
            ))}
          </div>
          <button
            onClick={handleSave}
            disabled={saving || !formName.trim()}
            className="h-8 px-3 rounded-lg border border-transparent text-xs cursor-pointer bg-text text-bg hover:bg-accent hover:text-accent-on disabled:opacity-50 transition-all duration-150 inline-flex items-center gap-1.5"
          >
            {saving ? <Loader2 size={12} className="animate-spin" /> : <Save size={12} />}
            保存
          </button>
        </div>
      )}

      {!testCases?.length && !showForm ? (
        <div className="py-12 text-center text-text-muted text-sm">
          暂无测试用例。点击「新增」创建第一个测试用例。
        </div>
      ) : (
        <div className="space-y-2">
          {testCases?.map((tc) => (
            <div key={tc.id} className="flex items-start gap-3 px-4 py-3 border border-border rounded-lg bg-bg-elevated">
              <div className="flex-1 min-w-0">
                <div className="text-sm font-medium text-text">{tc.name}</div>
                {tc.expectedBehavior && (
                  <div className="text-xs text-text-muted mt-0.5">{tc.expectedBehavior}</div>
                )}
                {Object.keys(tc.variableValues || {}).length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1.5">
                    {Object.entries(tc.variableValues).map(([k, v]) => (
                      <span key={k} className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-bg-subtle border border-border text-text-muted">
                        {k}={v}
                      </span>
                    ))}
                  </div>
                )}
              </div>
              <div className="flex items-center gap-1 shrink-0">
                <button
                  onClick={() => handleEdit(tc)}
                  className="w-7 h-7 rounded-md border border-border text-text-subtle hover:text-text cursor-pointer bg-transparent grid place-items-center transition-colors"
                >
                  <Pencil size={11} />
                </button>
                <button
                  onClick={() => handleDelete(tc.id)}
                  className="w-7 h-7 rounded-md border border-border text-text-subtle hover:text-danger cursor-pointer bg-transparent grid place-items-center transition-colors"
                >
                  <Trash2 size={11} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
