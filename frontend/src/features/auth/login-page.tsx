import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod/v4';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowRight, Eye, EyeOff, Loader2 } from 'lucide-react';
import { useAuth } from '@/features/auth/auth-context';
import type { ApiError } from '@/lib/api';

const loginSchema = z.object({
  email: z.email('请输入有效的邮箱地址'),
  password: z.string().min(6, '密码至少 6 位'),
  remember: z.boolean().optional(),
});

type LoginForm = z.infer<typeof loginSchema>;

export function LoginPage() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [showPw, setShowPw] = useState(false);
  const [serverError, setServerError] = useState('');
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: { remember: true },
  });

  const onSubmit = async (data: LoginForm) => {
    setServerError('');
    try {
      await login(data.email, data.password, data.remember);
      navigate('/prompts', { replace: true });
    } catch (err) {
      const apiErr = err as ApiError;
      setServerError(apiErr.message || '登录失败，请重试');
    }
  };

  return (
    <>
      <div className="font-mono text-[11px] font-medium uppercase tracking-[0.12em] text-text-subtle mb-4" />
      <h1 className="text-[28px] font-semibold tracking-tight text-text leading-tight mb-2">
        欢迎回到 PromptHub
      </h1>
      <p className="text-sm text-text-muted leading-relaxed mb-8">
        登录你的工作空间，继续管理团队的提示词资产。
      </p>

      <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-col gap-1.5">
          <label className="flex justify-between items-baseline text-[13px] font-medium text-text">
            <span>邮箱</span>
          </label>
          <input
            type="email"
            className={`w-full h-10 px-3 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${
              errors.email ? 'border-danger' : 'border-border'
            }`}
            placeholder="you@company.com"
            autoFocus
            {...register('email')}
          />
          {errors.email && (
            <span className="text-xs text-danger flex items-center gap-1">
              {errors.email.message}
            </span>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="flex justify-between items-baseline text-[13px] font-medium text-text">
            <span>密码</span>
            <Link
              to="/forgot-password"
              className="text-xs text-text-muted no-underline transition-colors duration-150 hover:text-accent"
            >
              忘记密码？
            </Link>
          </label>
          <div className="relative">
            <input
              type={showPw ? 'text' : 'password'}
              className={`w-full h-10 px-3 pr-10 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${
                errors.password ? 'border-danger' : 'border-border'
              }`}
              placeholder="至少 6 位"
              {...register('password')}
            />
            <button
              type="button"
              onClick={() => setShowPw((s) => !s)}
              className="absolute right-1 top-1/2 -translate-y-1/2 w-8 h-8 grid place-items-center text-text-subtle bg-transparent border-none cursor-pointer rounded-md transition-colors duration-150 hover:text-text"
              aria-label="切换密码可见"
            >
              {showPw ? <EyeOff size={16} /> : <Eye size={16} />}
            </button>
          </div>
          {errors.password && (
            <span className="text-xs text-danger flex items-center gap-1">
              {errors.password.message}
            </span>
          )}
        </div>

        <label className="flex items-start gap-2 text-[13px] text-text-muted cursor-pointer select-none leading-relaxed">
          <input
            type="checkbox"
            className="appearance-none w-4 h-4 shrink-0 mt-0.5 border border-border-strong rounded bg-bg-elevated cursor-pointer grid place-items-center transition-all duration-150 checked:bg-accent checked:border-accent after:content-[''] after:w-2 after:h-[5px] after:border-l-[1.5px] after:border-b-[1.5px] after:border-accent-on after:rotate-[-45deg] after:translate-x-[1px] after:-translate-y-[1px] checked:after:grid"
            {...register('remember')}
          />
          <span>30 天内保持登录状态</span>
        </label>

        {serverError && (
          <div className="text-xs text-danger">{serverError}</div>
        )}

        <button
          type="submit"
          disabled={isSubmitting}
          className="h-10 px-4 rounded-lg border border-transparent font-medium text-sm cursor-pointer transition-all duration-150 inline-flex items-center justify-center gap-2 whitespace-nowrap bg-text text-bg border-text hover:bg-accent hover:border-accent hover:text-accent-on disabled:opacity-50 disabled:cursor-not-allowed w-full"
        >
          {isSubmitting ? (
            <>
              <Loader2 size={14} className="animate-spin" />
              登录中…
            </>
          ) : (
            <>
              登录
              <ArrowRight size={14} />
            </>
          )}
        </button>

        <div className="flex items-center gap-3 my-6 text-[11px] font-mono uppercase tracking-[0.1em] text-text-subtle before:content-[''] before:flex-1 before:h-px before:bg-border after:content-[''] after:flex-1 after:h-px after:bg-border">
          OR
        </div>

        <button
          type="button"
          disabled
          className="h-10 px-4 rounded-lg border border-border font-medium text-sm cursor-default transition-all duration-150 inline-flex items-center justify-center gap-2 whitespace-nowrap w-full bg-transparent text-text hover:border-border-strong hover:bg-bg-subtle disabled:opacity-50"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0016 8c0-4.42-3.58-8-8-8z" />
          </svg>
          使用 GitHub 登录
          <span className="inline-flex items-center gap-1.5 text-[10px] font-mono uppercase tracking-[0.08em] px-1.5 py-0.5 rounded bg-bg-subtle text-text-subtle border border-border">
            即将上线
          </span>
        </button>
      </form>

      <div className="text-center text-[13px] text-text-muted mt-6">
        还没有账号？{' '}
        <Link
          to="/signup"
          className="text-text font-medium no-underline hover:text-accent transition-colors"
        >
          创建工作空间
        </Link>
      </div>
    </>
  );
}
