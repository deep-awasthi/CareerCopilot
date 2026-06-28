import { NavLink, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useAuthStore } from '@/stores/authStore'
import { useUnreadCount } from '@/api/hooks'
import {
  LayoutDashboard, FileText, Search, Briefcase, ClipboardList,
  Calendar, Users, Building2, Bell, BarChart3, Bookmark,
  Target, Settings, LogOut, Zap, ChevronRight
} from 'lucide-react'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/jobs', icon: Briefcase, label: 'Jobs' },
  { to: '/applications', icon: ClipboardList, label: 'Applications' },
  { to: '/interviews', icon: Calendar, label: 'Interviews' },
  { to: '/referrals', icon: Users, label: 'Referrals' },
  { to: '/companies', icon: Building2, label: 'Companies' },
  { to: '/watchlists', icon: Target, label: 'Watchlists' },
  { to: '/search-profiles', icon: Search, label: 'Search Profiles' },
  { to: '/resume', icon: FileText, label: 'Resume' },
  { to: '/alerts', icon: Zap, label: 'Keyword Alerts' },
  { to: '/bookmarks', icon: Bookmark, label: 'Bookmarks' },
  { to: '/analytics', icon: BarChart3, label: 'Analytics' },
  { to: '/notifications', icon: Bell, label: 'Notifications' },
  { to: '/settings', icon: Settings, label: 'Settings' },
]

export function Sidebar() {
  const { logout, user } = useAuthStore()
  const navigate = useNavigate()
  const { data: unreadData } = useUnreadCount()
  const unreadCount = unreadData?.count ?? 0

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <motion.aside
      initial={{ x: -20, opacity: 0 }}
      animate={{ x: 0, opacity: 1 }}
      className="fixed left-0 top-0 h-screen w-64 flex flex-col z-50"
      style={{
        background: 'rgba(15 15 22 / 0.95)',
        borderRight: '1px solid rgba(45 45 65 / 0.8)',
        backdropFilter: 'blur(20px)',
      }}
    >
      {/* Logo */}
      <div className="px-6 py-5 flex items-center gap-3" style={{ borderBottom: '1px solid rgba(45 45 65 / 0.5)' }}>
        <div className="w-8 h-8 rounded-lg flex items-center justify-center"
          style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
          <Briefcase size={16} color="white" />
        </div>
        <div>
          <span className="font-bold text-sm" style={{ color: 'rgb(248 250 252)' }}>CareerCopilot</span>
          <p className="text-xs" style={{ color: 'rgb(71 85 105)' }}>Track. Refer. Grow.</p>
        </div>
      </div>

      {/* User info */}
      <div className="px-6 py-4" style={{ borderBottom: '1px solid rgba(45 45 65 / 0.5)' }}>
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold"
            style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
            {user?.email?.[0]?.toUpperCase() || 'U'}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium truncate" style={{ color: 'rgb(248 250 252)' }}>
              {user?.email?.split('@')[0] || 'User'}
            </p>
            <p className="text-xs truncate" style={{ color: 'rgb(71 85 105)' }}>{user?.email}</p>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto px-3 py-3 space-y-1">
        {navItems.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
          >
            <Icon size={16} />
            <span className="flex-1">{label}</span>
            {label === 'Notifications' && unreadCount > 0 && (
              <span className="badge badge-purple text-xs px-2 py-0.5">{unreadCount}</span>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Logout */}
      <div className="px-3 py-3" style={{ borderTop: '1px solid rgba(45 45 65 / 0.5)' }}>
        <button onClick={handleLogout} className="nav-item w-full text-left"
          style={{ color: 'rgb(239 68 68)' }}>
          <LogOut size={16} />
          <span>Sign Out</span>
        </button>
      </div>
    </motion.aside>
  )
}

export function Topbar({ title }: { title: string }) {
  return (
    <div className="flex items-center justify-between px-8 py-5"
      style={{ borderBottom: '1px solid rgba(45 45 65 / 0.5)' }}>
      <h1 className="text-xl font-bold" style={{ color: 'rgb(248 250 252)' }}>{title}</h1>
    </div>
  )
}

export function PageLayout({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <div className="flex-1 ml-64 flex flex-col min-h-screen overflow-auto"
        style={{ background: 'rgb(10 10 15)' }}>
        <Topbar title={title} />
        <main className="flex-1 p-8">
          {children}
        </main>
      </div>
    </div>
  )
}
