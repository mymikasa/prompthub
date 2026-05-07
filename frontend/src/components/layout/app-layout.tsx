import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import {
  LayoutGrid,
  FileText,
  Activity,
  User,
  KeyRound,
  LogOut,
  Sun,
  Moon,
} from 'lucide-react';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/features/auth/auth-context';
import { useTheme } from '@/components/ui/theme-provider';

const NAV_ITEMS = [
  { icon: LayoutGrid, label: '概览', path: '/prompts' },
  { icon: FileText, label: '提示词', path: '/prompts', exact: true },
  { icon: Activity, label: '运行记录', path: '/runs' },
];

const SETTING_ITEMS = [
  { icon: User, label: '工作空间', path: '/settings' },
  { icon: KeyRound, label: 'API Keys', path: '/settings/provider' },
];

export function AppLayout() {
  const { user, logout } = useAuth();
  const { theme, toggle } = useTheme();
  const location = useLocation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const isActive = (path: string, exact?: boolean) => {
    if (exact) return location.pathname === path;
    return location.pathname.startsWith(path);
  };

  const handleLogout = async () => {
    await logout();
    queryClient.clear();
    navigate('/login', { replace: true });
  };

  return (
    <div className="min-h-screen bg-bg flex">
      <aside className="w-60 border-r border-border p-6 px-4 flex flex-col gap-1 bg-bg-subtle shrink-0">
        <Link to="/prompts" className="flex items-center gap-2.5 no-underline text-text mb-6 px-2">
          <div className="w-7 h-7 rounded-md bg-text text-bg grid place-items-center font-mono font-bold text-sm tracking-tight">
            P
          </div>
          <span className="font-semibold text-[15px] tracking-tight">PromptHub</span>
        </Link>

        {NAV_ITEMS.map((item) => (
          <div
            key={item.label}
            onClick={() => navigate(item.path)}
            className={`flex items-center gap-2.5 px-2.5 py-2 rounded-md text-[13px] cursor-pointer transition-all duration-150 ${
              isActive(item.path, item.exact)
                ? 'bg-bg-elevated text-text font-medium'
                : 'text-text-muted hover:bg-bg-elevated hover:text-text'
            }`}
          >
            <item.icon size={14} />
            {item.label}
          </div>
        ))}

        <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-text-subtle px-2.5 pt-4 pb-1.5">
          设置
        </div>

        {SETTING_ITEMS.map((item) => (
          <div
            key={item.label}
            onClick={() => navigate(item.path)}
            className={`flex items-center gap-2.5 px-2.5 py-2 rounded-md text-[13px] cursor-pointer transition-all duration-150 ${
              isActive(item.path)
                ? 'bg-bg-elevated text-text font-medium'
                : 'text-text-muted hover:bg-bg-elevated hover:text-text'
            }`}
          >
            <item.icon size={14} />
            {item.label}
          </div>
        ))}

        <div className="mt-auto pt-3 px-2 border-t border-border flex items-center gap-2.5">
          <div className="w-7 h-7 rounded-full bg-accent text-accent-on grid place-items-center text-xs font-semibold shrink-0">
            {(user?.name || user?.email || 'A')[0].toUpperCase()}
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-[13px] font-medium truncate">{user?.name || '我的账号'}</div>
            <div className="text-[11px] text-text-muted truncate">{user?.email}</div>
          </div>
          <button
            onClick={toggle}
            className="w-7 h-7 rounded-md border border-border bg-bg-elevated text-text-muted grid place-items-center cursor-pointer transition-all duration-150 hover:text-text hover:border-border-strong shrink-0"
            title="切换主题"
          >
            {theme === 'light' ? <Moon size={12} /> : <Sun size={12} />}
          </button>
          <button
            onClick={handleLogout}
            className="w-7 h-7 rounded-md border border-border bg-bg-elevated text-text-muted grid place-items-center cursor-pointer transition-all duration-150 hover:text-text hover:border-border-strong shrink-0"
            title="退出登录"
          >
            <LogOut size={12} />
          </button>
        </div>
      </aside>

      <main className="flex-1 p-8 px-12 overflow-auto min-w-0">
        <Outlet />
      </main>
    </div>
  );
}
