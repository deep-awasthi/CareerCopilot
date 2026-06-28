import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ProtectedRoute, PublicRoute } from '@/components/layout/ProtectedRoute'
// @ts-ignore
import './index.css'

// Pages (lazy-loaded)
const LoginPage = React.lazy(() => import('@/pages/Login'))
const RegisterPage = React.lazy(() => import('@/pages/Register'))
const DashboardPage = React.lazy(() => import('@/pages/Dashboard'))
const JobsPage = React.lazy(() => import('@/pages/Jobs'))
const ApplicationsPage = React.lazy(() => import('@/pages/Applications'))
const InterviewsPage = React.lazy(() => import('@/pages/Interviews'))
const ReferralsPage = React.lazy(() => import('@/pages/Referrals'))
const CompaniesPage = React.lazy(() => import('@/pages/Companies'))
const WatchlistsPage = React.lazy(() => import('@/pages/Watchlists'))
const ResumePage = React.lazy(() => import('@/pages/Resume'))
const SearchProfilesPage = React.lazy(() => import('@/pages/SearchProfiles'))
const AlertsPage = React.lazy(() => import('@/pages/Alerts'))
const BookmarksPage = React.lazy(() => import('@/pages/Bookmarks'))
const AnalyticsPage = React.lazy(() => import('@/pages/Analytics'))
const NotificationsPage = React.lazy(() => import('@/pages/Notifications'))
const SettingsPage = React.lazy(() => import('@/pages/Settings'))

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
      refetchOnWindowFocus: false,
    },
  },
})

const Fallback = () => (
  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: 'rgb(10 10 15)' }}>
    <div style={{ width: 32, height: 32, border: '3px solid rgba(99 102 241 / 0.3)', borderTopColor: 'rgb(99 102 241)', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
    <style>{`@keyframes spin { to { transform: rotate(360deg) } }`}</style>
  </div>
)

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <React.Suspense fallback={<Fallback />}>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<PublicRoute><LoginPage /></PublicRoute>} />
            <Route path="/register" element={<PublicRoute><RegisterPage /></PublicRoute>} />

            {/* Protected routes */}
            <Route path="/" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
            <Route path="/jobs" element={<ProtectedRoute><JobsPage /></ProtectedRoute>} />
            <Route path="/applications" element={<ProtectedRoute><ApplicationsPage /></ProtectedRoute>} />
            <Route path="/interviews" element={<ProtectedRoute><InterviewsPage /></ProtectedRoute>} />
            <Route path="/referrals" element={<ProtectedRoute><ReferralsPage /></ProtectedRoute>} />
            <Route path="/companies" element={<ProtectedRoute><CompaniesPage /></ProtectedRoute>} />
            <Route path="/watchlists" element={<ProtectedRoute><WatchlistsPage /></ProtectedRoute>} />
            <Route path="/resume" element={<ProtectedRoute><ResumePage /></ProtectedRoute>} />
            <Route path="/search-profiles" element={<ProtectedRoute><SearchProfilesPage /></ProtectedRoute>} />
            <Route path="/alerts" element={<ProtectedRoute><AlertsPage /></ProtectedRoute>} />
            <Route path="/bookmarks" element={<ProtectedRoute><BookmarksPage /></ProtectedRoute>} />
            <Route path="/analytics" element={<ProtectedRoute><AnalyticsPage /></ProtectedRoute>} />
            <Route path="/notifications" element={<ProtectedRoute><NotificationsPage /></ProtectedRoute>} />
            <Route path="/settings" element={<ProtectedRoute><SettingsPage /></ProtectedRoute>} />

            {/* Fallback */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </React.Suspense>
      </BrowserRouter>
    </QueryClientProvider>
  </React.StrictMode>
)
