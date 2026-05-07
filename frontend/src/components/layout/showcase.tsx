import { useState, useEffect, useCallback } from 'react';

const TEMPLATE = [
  {
    role: 'system',
    text: 'You are a senior {{role}} reviewing code for {{company}}. Be concise, direct, and prioritize {{priority}}.',
  },
  { role: 'user', text: 'Review this pull request:\n\n{{diff}}' },
];

const VARIABLES = [
  { key: 'role', value: 'engineer' },
  { key: 'company', value: 'Acme Corp' },
  { key: 'priority', value: 'security' },
  { key: 'diff', value: '<patch.diff>' },
];

const FULL_SOURCE = TEMPLATE.map((m) => `${m.role.toUpperCase()}\n${m.text}`).join('\n\n');
const TOTAL_CHARS = FULL_SOURCE.length;

type Phase = 'typing' | 'filling' | 'rendering' | 'reset';

function useAnimationLoop() {
  const [phase, setPhase] = useState<Phase>('typing');
  const [typedChars, setTypedChars] = useState(0);
  const [filledVars, setFilledVars] = useState<Record<string, number>>({});
  const [activeVar, setActiveVar] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'template' | 'rendered'>('template');

  const reset = useCallback(() => {
    setTypedChars(0);
    setFilledVars({});
    setActiveVar(null);
    setActiveTab('template');
    setPhase('typing');
  }, []);

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout>;

    if (phase === 'typing') {
      if (typedChars < TOTAL_CHARS) {
        timer = setTimeout(() => setTypedChars((c) => c + 2), 18);
      } else {
        timer = setTimeout(() => {
          setPhase('filling');
          setActiveTab('rendered');
        }, 600);
      }
    } else if (phase === 'filling') {
      const currentVar = VARIABLES.find((v) => (filledVars[v.key] || 0) < v.value.length);
      if (!currentVar) {
        timer = setTimeout(() => setPhase('rendering'), 500);
      } else {
        setActiveVar(currentVar.key);
        const filled = filledVars[currentVar.key] || 0;
        if (filled < currentVar.value.length) {
          timer = setTimeout(() => {
            setFilledVars((prev) => ({ ...prev, [currentVar.key]: filled + 1 }));
          }, 55);
        }
      }
    } else if (phase === 'rendering') {
      timer = setTimeout(() => setPhase('reset'), 2400);
    } else {
      timer = setTimeout(reset, 400);
    }

    return () => clearTimeout(timer);
  }, [phase, typedChars, filledVars, reset]);

  return { phase, typedChars, filledVars, activeVar, activeTab };
}

function BlinkingCursor({ color }: { color?: string }) {
  return (
    <span
      className="inline-block w-[7px] h-[1em] align-text-bottom ml-px"
      style={{
        background: color ?? 'var(--color-code-text)',
        animation: 'blink 1s steps(2) infinite',
      }}
    />
  );
}

function renderTypingBlock(typedChars: number) {
  let remaining = typedChars;
  const blocks: React.ReactNode[] = [];

  for (let i = 0; i < TEMPLATE.length; i++) {
    const msg = TEMPLATE[i];
    const roleStr = msg.role.toUpperCase();
    const fullBlock = `${roleStr}\n${msg.text}`;
    if (remaining <= 0) break;

    const slice = fullBlock.slice(0, remaining);
    const rolePart = slice.slice(0, Math.min(roleStr.length, slice.length));
    const bodyPart = slice.length > roleStr.length ? slice.slice(roleStr.length + 1) : '';
    const parts = bodyPart.split(/(\{\{[a-z]*\}?\}?)/g);

    blocks.push(
      <div key={i} className="min-h-[1.7em] block">
        <span className="text-code-role font-semibold">{rolePart}</span>
        {bodyPart && '\n'}
        {parts.map((p, j) => {
          if (/^\{\{[a-z]+\}\}$/.test(p)) {
            return (
              <span
                key={j}
                className="text-code-var bg-[rgba(251,191,36,0.12)] px-1 rounded border border-[rgba(251,191,36,0.25)]"
              >
                {p}
              </span>
            );
          }
          if (/^\{\{/.test(p)) {
            return (
              <span
                key={j}
                className="text-code-var bg-[rgba(251,191,36,0.12)] px-1 rounded border border-[rgba(251,191,36,0.25)] opacity-70"
              >
                {p}
              </span>
            );
          }
          return <span key={j}>{p}</span>;
        })}
        {(i === TEMPLATE.length - 1 || remaining < fullBlock.length) &&
          remaining < fullBlock.length && <BlinkingCursor />}
      </div>,
    );

    remaining -= fullBlock.length + 2;
  }

  return blocks;
}

function renderFilledBlock(
  msgIdx: number,
  filledVars: Record<string, number>,
  activeVar: string | null,
) {
  const msg = TEMPLATE[msgIdx];
  const parts = msg.text.split(/(\{\{[a-z]+\}\})/g);

  return (
    <div key={msgIdx} className="min-h-[1.7em] block">
      <span className="text-code-role font-semibold">{msg.role.toUpperCase()}</span>
      {'\n'}
      {parts.map((p, i) => {
        const m = p.match(/^\{\{([a-z]+)\}\}$/);
        if (!m) return <span key={i}>{p}</span>;

        const varName = m[1];
        const variable = VARIABLES.find((v) => v.key === varName);
        const filled = filledVars[varName] || 0;
        const fullVal = variable?.value ?? '';
        const partialVal = fullVal.slice(0, filled);
        const isActive = activeVar === varName && filled < fullVal.length;

        if (filled === 0) {
          return (
            <span
              key={i}
              className="text-code-var bg-[rgba(251,191,36,0.12)] px-1 rounded border border-[rgba(251,191,36,0.25)]"
            >
              {p}
            </span>
          );
        }

        return (
          <span
            key={i}
            className="text-[#4ade80] bg-[rgba(74,222,128,0.1)] px-1 rounded border border-[rgba(74,222,128,0.25)]"
          >
            {partialVal}
            {isActive && <BlinkingCursor color="#4ade80" />}
          </span>
        );
      })}
    </div>
  );
}

export function Showcase() {
  const { phase, typedChars, filledVars, activeVar, activeTab } = useAnimationLoop();

  return (
    <div className="hidden lg:flex relative bg-bg-showcase border-l border-border p-12 overflow-hidden flex-col justify-center">
      <style>{`
        @keyframes blink { 50% { opacity: 0; } }
        @keyframes showcasePulse { 0%,100%{opacity:1} 50%{opacity:0.4} }
      `}</style>

      <div
        className="absolute inset-0 opacity-40"
        style={{
          backgroundImage:
            'linear-gradient(var(--color-border) 1px, transparent 1px), linear-gradient(90deg, var(--color-border) 1px, transparent 1px)',
          backgroundSize: '48px 48px',
          maskImage: 'radial-gradient(circle at 50% 50%, #000 30%, transparent 75%)',
          WebkitMaskImage: 'radial-gradient(circle at 50% 50%, #000 30%, transparent 75%)',
        }}
      />
      <div className="relative max-w-[560px] mx-auto w-full">
        <div className="flex items-center gap-2 font-mono text-[11px] uppercase tracking-[0.14em] text-text-subtle mb-3">
          <span className="w-1.5 h-1.5 rounded-full bg-success" style={{ animation: 'showcasePulse 2s ease-in-out infinite' }} />
          <span>实时预览 &middot; LIVE PREVIEW</span>
        </div>
        <h2 className="text-[28px] font-semibold tracking-tight leading-snug text-text mb-2">
          为团队管理可复用的提示词模板
        </h2>
        <p className="text-sm text-text-muted leading-relaxed mb-7 max-w-[460px]">
          用变量、版本和测试用例把零散的 prompt 整合成可维护的资产。下面是一个真实的提示词在被渲染。
        </p>

        <div className="bg-bg-elevated border border-border rounded-xl overflow-hidden shadow-lg">
          <div className="flex border-b border-border px-1 bg-bg-subtle">
            <div
              className={`px-3.5 py-2.5 text-xs font-mono cursor-pointer border-b-[1.5px] mb-[-1px] flex items-center gap-1.5 transition-colors duration-150 ${
                activeTab === 'template'
                  ? 'text-text border-b-accent'
                  : 'text-text-subtle border-b-transparent'
              }`}
            >
              <span className={`w-1.5 h-1.5 rounded-full ${activeTab === 'template' ? 'bg-accent' : 'bg-text-subtle'}`} />
              <span>template.prompt</span>
            </div>
            <div
              className={`px-3.5 py-2.5 text-xs font-mono cursor-pointer border-b-[1.5px] mb-[-1px] flex items-center gap-1.5 transition-colors duration-150 ${
                activeTab === 'rendered'
                  ? 'text-text border-b-accent'
                  : 'text-text-subtle border-b-transparent'
              }`}
            >
              <span className={`w-1.5 h-1.5 rounded-full ${activeTab === 'rendered' ? 'bg-accent' : 'bg-text-subtle'}`} />
              <span>rendered.txt</span>
            </div>
          </div>

          <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border text-[11px] font-mono text-text-muted">
            <span className="px-1.5 py-0.5 rounded bg-bg-subtle border border-border">code-review</span>
            <span className="px-1.5 py-0.5 rounded bg-bg-subtle border border-border">chat_messages</span>
            <span className="px-1.5 py-0.5 rounded bg-bg-subtle border border-border">gpt-4o</span>
            <span className="ml-auto text-text-subtle">v3 &middot; active</span>
          </div>

          <div className="p-4 px-5 font-mono text-[13px] leading-[1.7] text-code-text bg-code-bg whitespace-pre-wrap min-h-[220px]">
            {phase === 'typing' && renderTypingBlock(typedChars)}
            {phase !== 'typing' &&
              TEMPLATE.map((_, i) => renderFilledBlock(i, filledVars, activeVar))}
          </div>

          <div className="border-t border-border p-3.5 px-5 bg-bg-elevated grid grid-cols-2 gap-3">
            {VARIABLES.map((v) => {
              const filled = filledVars[v.key] || 0;
              const partial = v.value.slice(0, filled);
              const isActive = activeVar === v.key && filled < v.value.length;
              return (
                <div key={v.key} className="flex flex-col gap-1 min-w-0">
                  <div className="font-mono text-[11px] text-text-subtle flex items-center gap-1">
                    <span className="opacity-60">{'{{'}</span>
                    {v.key}
                    <span className="opacity-60">{'}}'}</span>
                  </div>
                  <div
                    className={`text-[13px] text-text font-mono bg-bg-subtle border rounded-md px-2.5 py-1.5 min-h-[30px] whitespace-nowrap overflow-hidden text-ellipsis transition-all duration-300 ${
                      isActive ? 'border-accent bg-accent-soft' : 'border-border'
                    }`}
                  >
                    {partial || <span className="text-text-subtle">&mdash;</span>}
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <div className="mt-6 flex flex-wrap gap-2">
          {['版本历史', '变量提取', '测试用例', '团队工作空间'].map((label) => (
            <div
              key={label}
              className="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-full bg-bg-elevated border border-border text-xs text-text-muted"
            >
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" className="text-accent">
                <path d="M2 6L5 9L10 3" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              {label}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
