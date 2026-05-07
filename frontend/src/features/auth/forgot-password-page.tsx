import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod/v4';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowRight, ArrowLeft, Loader2, CircleCheck } from 'lucide-react';
import { forgotPassword } from '@/features/auth/api';

const forgotSchema = z.object({
  email: z.email('请输入有效的邮箱地址'),
});

type ForgotForm = z.infer<typeof forgotSchema>;

export function ForgotPasswordPage() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);
  const [serverError, setServerError] = useState('');
  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm<ForgotForm>({
    resolver: zodResolver(forgotSchema),
  });

  const emailValue = getValues('email');

  const onSubmit = async (data: ForgotForm) => {
    setServerError('');
    setLoading(true);
    try {
      await forgotPassword(data.email);
      setSent(true);
    } catch {
      setServerError('发送失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  if (sent) {
    return (
      <>
        <div className="text-success mb-4">
          <CircleCheck size={20} />
        </div>
        <h1 className="text-[28px] font-semibold tracking-tight text-text leading-tight mb-2">
          邮件已发送
        </h1>
        <p className="text-sm text-text-muted leading-relaxed mb-2">
          我们已向{' '}
          <strong className="text-text">{emailValue}</strong>{' '}
          发送了重置链接。请检查你的收件箱（或垃圾邮件文件夹）。
        </p>
        <div className="flex flex-col gap-3 mt-2">
          <button
            onClick={() => navigate('/login')}
            className="h-10 px-4 rounded-lg border border-transparent font-medium text-sm cursor-pointer transition-all duration-150 inline-flex items-center justify-center gap-2 whitespace-nowrap bg-text text-bg border-text hover:bg-accent hover:border-accent hover:text-accent-on w-full"
          >
            返回登录
          </button>
          <button
            onClick={() => {
              setSent(false);
            }}
            className="h-10 px-4 rounded-lg border border-border font-medium text-sm cursor-pointer transition-all duration-150 inline-flex items-center justify-center gap-2 whitespace-nowrap w-full bg-transparent text-text hover:border-border-strong hover:bg-bg-subtle"
          >
            发送到其他邮箱
          </button>
        </div>
        <div className="mt-6 p-3 bg-bg-subtle border border-border rounded-lg text-xs text-text-muted leading-relaxed">
          <strong className="text-text font-mono text-[11px] uppercase tracking-[0.08em]">
            提示
          </strong>
          <br />
          没收到邮件？检查垃圾箱，或在 60 秒后重新发送。
        </div>
      </>
    );
  }

  return (
    <>
      <Link
        to="/login"
        className="inline-flex items-center gap-1.5 text-text-muted text-[13px] no-underline mb-6 transition-colors hover:text-text"
      >
        <ArrowLeft size={14} />
        返回登录
      </Link>
      <div className="font-mono text-[11px] font-medium uppercase tracking-[0.12em] text-text-subtle mb-4">
        / 找回密码
      </div>
      <h1 className="text-[28px] font-semibold tracking-tight text-text leading-tight mb-2">
        重置密码
      </h1>
      <p className="text-sm text-text-muted leading-relaxed mb-8">
        输入注册时使用的邮箱，我们会发送一条重置链接给你。
      </p>

      <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-col gap-1.5">
          <label className="text-[13px] font-medium text-text">
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
            <span className="text-xs text-danger">{errors.email.message}</span>
          )}
        </div>

        {serverError && (
          <div className="text-xs text-danger">{serverError}</div>
        )}

        <button
          type="submit"
          disabled={loading}
          className="h-10 px-4 rounded-lg border border-transparent font-medium text-sm cursor-pointer transition-all duration-150 inline-flex items-center justify-center gap-2 whitespace-nowrap bg-text text-bg border-text hover:bg-accent hover:border-accent hover:text-accent-on disabled:opacity-50 disabled:cursor-not-allowed w-full"
        >
          {loading ? (
            <>
              <Loader2 size={14} className="animate-spin" />
              发送中…
            </>
          ) : (
            <>
              发送重置链接
              <ArrowRight size={14} />
            </>
          )}
        </button>
      </form>
    </>
  );
}
