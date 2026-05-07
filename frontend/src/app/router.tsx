import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AuthLayout } from '@/components/layout/auth-layout';
import { AppLayout } from '@/components/layout/app-layout';
import { LoginPage } from '@/features/auth/login-page';
import { SignupPage } from '@/features/auth/signup-page';
import { ForgotPasswordPage } from '@/features/auth/forgot-password-page';
import { PromptListPage } from '@/features/prompts/prompt-list-page';
import { PromptDetailPage } from '@/features/prompts/prompt-detail-page';
import { VersionDetailPage } from '@/features/prompts/version-detail-page';
import { CreatePromptPage } from '@/features/prompts/create-prompt-page';
import { SettingsPage } from '@/features/settings/settings-page';
import { RunsPage } from '@/features/runs/runs-page';

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
        element: <CreatePromptPage />,
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
        element: <RunsPage />,
      },
      {
        path: '/settings',
        element: <SettingsPage />,
      },
      {
        path: '/settings/provider',
        element: <SettingsPage />,
      },
    ],
  },
  { path: '/', element: <Navigate to="/login" replace /> },
  { path: '*', element: <Navigate to="/login" replace /> },
]);
