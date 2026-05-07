import { Link, Outlet } from 'react-router-dom';
import { Sun, Moon } from 'lucide-react';
import { useTheme } from '@/components/ui/theme-provider';
import { Showcase } from './showcase';

export function AuthLayout() {
  const { theme, toggle } = useTheme();

  return (
    <div className="min-h-screen grid grid-cols-1 lg:grid-cols-[minmax(440px,1fr)_1.15fr] bg-bg">
      <div className="relative flex flex-col px-8 py-8 lg:px-12 lg:py-8 min-h-screen bg-bg">
        <div className="flex items-center justify-between">
          <Link to="/login" className="flex items-center gap-2.5 no-underline text-text">
            <div className="w-7 h-7 rounded-md bg-text text-bg grid place-items-center font-mono font-bold text-sm tracking-tight">
              P
            </div>
            <span className="font-semibold text-[15px] tracking-tight">
              PromptHub
            </span>
          </Link>
          <button
            onClick={toggle}
            className="w-8 h-8 rounded-lg border border-border bg-bg-elevated text-text-muted grid place-items-center cursor-pointer transition-all duration-150 hover:text-text hover:border-border-strong"
            title="切换主题"
          >
            {theme === 'light' ? <Moon size={14} /> : <Sun size={14} />}
          </button>
        </div>

        <div className="flex-1 flex items-center justify-center py-10">
          <div className="w-full max-w-[380px]">
            <Outlet />
          </div>
        </div>

        <div className="flex justify-between items-center text-xs text-text-subtle font-mono">
          <span>&copy; 2026 PromptHub</span>
          <span className="flex gap-4">
            <a href="#" onClick={(e) => e.preventDefault()} className="text-text-subtle no-underline hover:text-text-muted transition-colors">
              服务条款
            </a>
            <a href="#" onClick={(e) => e.preventDefault()} className="text-text-subtle no-underline hover:text-text-muted transition-colors">
              隐私政策
            </a>
            <a href="#" onClick={(e) => e.preventDefault()} className="text-text-subtle no-underline hover:text-text-muted transition-colors">
              文档
            </a>
          </span>
        </div>
      </div>

      <Showcase />
    </div>
  );
}
