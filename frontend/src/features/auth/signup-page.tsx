import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod/v4';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowRight, Eye, EyeOff, Loader2 } from 'lucide-react';
import { useAuth } from '@/features/auth/auth-context';
import type { ApiError } from '@/lib/api';

const signupSchema = z
  .object({
    name: z.string().min(1, '请填写姓名'),
    email: z.email('请输入有效的邮箱地址'),
    password: z.string().min(8, '密码至少 8 位'),
    agree: z.literal(true, { message: '请先同意服务条款' }),
  })
  .transform((data) => ({ ...data, agree: data.agree as true }));

type SignupForm = z.input<typeof signupSchema>;

function getPasswordStrength(password: string) {
  let score = 0;
  if (password.length >= 8) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^A-Za-z0-9]/.test(password)) score++;
  return score;
}

const strengthLabels = ['—', '弱', '一般', '良好', '很强'];

export function SignupPage() {
  const navigate = useNavigate();
  const { signup } = useAuth();
  const [showPw, setShowPw] = useState(false);
  const [serverError, setServerError] = useState('');
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<SignupForm>({
    resolver: zodResolver(signupSchema),
  });

  const password = watch('password', '');
  const pwScore = getPasswordStrength(password);

  const onSubmit = async (data: SignupForm) => {
    setServerError('');
    try {
      await signup(data.name, data.email, data.password);
      navigate('/prompts', { replace: true });
    } catch (err) {
      const apiErr = err as ApiError;
      setServerError(apiErr.message || '注册失败，请重试');
    }
  };

  return (
    <>
      <div className="font-mono text-[11px] font-medium uppercase tracking-[0.12em] text-text-subtle mb-4">
        / 注册
      </div>
      <h1 className="text-[28px] font-semibold tracking-tight text-text leading-tight mb-2">
        创建你的工作空间
      </h1>
      <p className="text-sm text-text-muted leading-relaxed mb-8">
        几秒钟开始整理你的提示词。注册后会自动创建一个个人工作空间。
      </p>

      <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-medium text-text">
            <span>姓名</span>
          </label>
          <input
            className={`w-full h-10 px-3 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${
              errors.name ? 'border-danger' : 'border-border'
            }`}
            placeholder="张三"
            autoFocus
            {...register('name')}
          />
          {errors.name && (
            <span className="text-xs text-danger">{errors.name.message}</span>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-medium text-text">
            <span>工作邮箱</span>
          </label>
          <input
            type="email"
            className={`w-full h-10 px-3 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${
              errors.email ? 'border-danger' : 'border-border'
            }`}
            placeholder="you@company.com"
            {...register('email')}
          />
          {errors.email && (
            <span className="text-xs text-danger">{errors.email.message}</span>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-medium text-text">
            <span>设置密码</span>
          </label>
          <div className="relative">
            <input
              type={showPw ? 'text' : 'password'}
              className={`w-full h-10 px-3 pr-10 bg-bg-elevated border rounded-lg text-sm text-text outline-none transition-all duration-150 placeholder:text-text-subtle hover:border-border-strong focus:border-accent focus:shadow-[0_0_0_3px_var(--color-accent-soft)] ${
                errors.password ? 'border-danger' : 'border-border'
              }`}
              placeholder="至少 8 位，建议混合大小写和数字"
              {...register('password')}
            />
            <button
              type="button"
              onClick={() => setShowPw((s) => !s)}
              className="absolute right-1 top-1/2 -translate-y-1/2 w-8 h-8 grid place-items-center text-text-subtle bg-transparent border-none cursor-pointer rounded-md transition-colors duration-150 hover:text-text"
            >
              {showPw ? <EyeOff size={16} /> : <Eye size={16} />}
            </button>
          </div>
          {password.length > 0 && (
            <div className="flex items-center gap-2 mt-1.5">
              <div className="flex-1 h-[3px] bg-bg-subtle rounded overflow-hidden">
                <div
                  className="h-full transition-all duration-200"
                  style={{
                    width: `${(pwScore / 4) * 100}%`,
                    background:
                      pwScore < 2
                        ? 'var(--color-danger)'
                        : pwScore < 3
                          ? '#f59e0b'
                          : 'var(--color-success)',
                  }}
                />
              </div>
              <span className="text-[11px] text-text-muted font-mono min-w-[28px]">
                {strengthLabels[pwScore]}
              </span>
            </div>
          )}
          {errors.password && (
            <span className="text-xs text-danger">{errors.password.message}</span>
          )}
        </div>

        <label className="flex items-start gap-2 text-[13px] text-text-muted cursor-pointer select-none leading-relaxed">
          <input
            type="checkbox"
            className="appearance-none w-4 h-4 shrink-0 mt-0.5 border border-border-strong rounded bg-bg-elevated cursor-pointer grid place-items-center transition-all duration-150 checked:bg-accent checked:border-accent after:content-[''] after:w-2 after:h-[5px] after:border-l-[1.5px] after:border-b-[1.5px] after:border-accent-on after:rotate-[-45deg] after:translate-x-[1px] after:-translate-y-[1px] checked:after:grid"
            {...register('agree')}
          />
          <span>
            我同意{' '}
            <a
              href="#"
              onClick={(e) => e.preventDefault()}
              className="text-text underline decoration-border-strong underline-offset-[2px] hover:decoration-accent transition-colors"
            >
              服务条款
            </a>{' '}
            和{' '}
            <a
              href="#"
              onClick={(e) => e.preventDefault()}
              className="text-text underline decoration-border-strong underline-offset-[2px] hover:decoration-accent transition-colors"
            >
              隐私政策
            </a>
          </span>
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
              创建中…
            </>
          ) : (
            <>
              创建账号
              <ArrowRight size={14} />
            </>
          )}
        </button>
      </form>

      <div className="text-center text-[13px] text-text-muted mt-6">
        已经有账号了？{' '}
        <Link
          to="/login"
          className="text-text font-medium no-underline hover:text-accent transition-colors"
        >
          直接登录
        </Link>
      </div>
    </>
  );
}
