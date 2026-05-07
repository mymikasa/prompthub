import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AuthLayout } from '@/components/layout/auth-layout';
import { AppLayout } from '@/components/layout/app-layout';
import { LoginPage } from '@/features/auth/login-page';
import { SignupPage } from '@/features/auth/signup-page';
import { ForgotPasswordPage } from '@/features/auth/forgot-password-page';
import { PromptListPage } from '@/features/prompts/prompt-list-page';
import { PromptDetailPage } from '@/features/prompts/prompt-detail-page';
import { VersionDetailPage } from '@/features/prompts/version-detail-page';

export const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/signup', element: <SignupPage /> },
      { path: '/forgot-password', element: <ForgotPasswordPage /> },
    ],
  },
  {
    element: <AppLayout />,
    children: [
      { path: '/prompts', element: <PromptListPage /> },
      {
        path: '/prompts/new',
        element: (
          <div className="flex items-center justify-center py-20 text-text-muted text-sm">
            创建提示词（开发中）
          </div>
        ),
      },
      {
        path: '/prompts/:promptId',
        element: <PromptDetailPage />,
      },
      {
        path: '/prompts/:promptId/edit',
        element: (
          <div className="flex items-center justify-center py-20 text-text-muted text-sm">
            编辑提示词（开发中）
          </div>
        ),
      },
      {
        path: '/prompts/:promptId/versions/:versionId',
        element: <VersionDetailPage />,
      },
      {
        path: '/runs',
        element: (
          <div className="flex items-center justify-center py-20 text-text-muted text-sm">
            运行记录页（开发中）
          </div>
        ),
      },
      {
        path: '/settings',
        element: (
          <div className="flex items-center justify-center py-20 text-text-muted text-sm">
            设置页（开发中）
          </div>
        ),
      },
      {
        path: '/settings/provider',
        element: (
          <div className="flex items-center justify-center py-20 text-text-muted text-sm">
            模型提供方配置（开发中）
          </div>
        ),
      },
    ],
  },
  { path: '/', element: <Navigate to="/login" replace /> },
  { path: '*', element: <Navigate to="/login" replace /> },
]);
