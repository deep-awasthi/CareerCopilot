import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useDashboardStats, useUpcomingInterviews, useApplicationStats } from '@/api/hooks'
import {
  Briefcase, ClipboardList, Calendar, Users,
  Building2, Bell, Zap, TrendingUp, ChevronRight, Clock
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

interface StatCard {
  label: string
  value: number | string
  icon: React.ElementType
  color: string
  bg: string
  to: string
  change?: string
}

export default function DashboardPage() {
  const { data: stats, isLoading } = useDashboardStats()
  const { data: upcoming } = useUpcomingInterviews()
  const { data: appStats } = useApplicationStats()
  const { user } = useAuthStore()

  const firstName = user?.email?.split('@')[0] || 'there'

  const statCards: StatCard[] = [
    {
      label: 'New Jobs Today',
      value: stats?.jobs_found_today ?? '—',
      icon: Briefcase,
      color: 'rgb(129 140 248)',
      bg: 'rgba(99 102 241 / 0.12)',
      to: '/jobs',
      change: '+12%',
    },
    {
      label: 'Applications',
      value: (stats?.applied_jobs ?? 0) + (stats?.saved_jobs ?? 0),
      icon: ClipboardList,
      color: 'rgb(34 197 94)',
      bg: 'rgba(34 197 94 / 0.12)',
      to: '/applications',
    },
    {
      label: 'Interviews',
      value: stats?.interviews ?? '—',
      icon: Calendar,
      color: 'rgb(250 204 21)',
      bg: 'rgba(250 204 21 / 0.12)',
      to: '/interviews',
    },
    {
      label: 'Referrals',
      value: stats?.referral_opportunities ?? '—',
      icon: Users,
      color: 'rgb(251 113 133)',
      bg: 'rgba(251 113 133 / 0.12)',
      to: '/referrals',
    },
    {
      label: 'Companies Watching',
      value: stats?.companies_following ?? '—',
      icon: Building2,
      color: 'rgb(96 165 250)',
      bg: 'rgba(96 165 250 / 0.12)',
      to: '/watchlists',
    },
    {
      label: 'Active Alerts',
      value: stats?.keyword_alerts ?? '—',
      icon: Zap,
      color: 'rgb(167 139 250)',
      bg: 'rgba(167 139 250 / 0.12)',
      to: '/alerts',
    },
  ]

  const applicationStatuses = [
    { label: 'Saved', status: 'saved', color: 'rgb(100 116 139)', value: stats?.saved_jobs ?? 0 },
    { label: 'Applied', status: 'applied', color: 'rgb(59 130 246)', value: stats?.applied_jobs ?? 0 },
    { label: 'Interview', status: 'interview', color: 'rgb(250 204 21)', value: stats?.interviews ?? 0 },
    { label: 'Offer', status: 'offer', color: 'rgb(34 197 94)', value: stats?.offers ?? 0 },
  ]

  return (
    <PageLayout title="Dashboard">
      {/* Greeting */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-8"
      >
        <h2 className="text-2xl font-bold mb-1" style={{ color: 'rgb(248 250 252)' }}>
          Good morning, {firstName}! 👋
        </h2>
        <p style={{ color: 'rgb(71 85 105)' }}>
          Here's a snapshot of your career progress.
        </p>
      </motion.div>

      {/* Stat Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
        {statCards.map((card, i) => (
          <motion.div
            key={card.label}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.05 }}
          >
            <Link to={card.to} className="stat-card flex items-start justify-between group cursor-pointer block">
              <div>
                <p className="text-xs font-semibold uppercase tracking-wider mb-2" style={{ color: 'rgb(71 85 105)' }}>
                  {card.label}
                </p>
                <p className="text-3xl font-bold" style={{ color: 'rgb(248 250 252)' }}>
                  {isLoading ? '—' : card.value}
                </p>
                {card.change && (
                  <span className="text-xs mt-1 flex items-center gap-1" style={{ color: 'rgb(34 197 94)' }}>
                    <TrendingUp size={12} /> {card.change}
                  </span>
                )}
              </div>
              <div className="w-10 h-10 rounded-xl flex items-center justify-center"
                style={{ background: card.bg }}>
                <card.icon size={18} style={{ color: card.color }} />
              </div>
            </Link>
          </motion.div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Application Pipeline */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="glass-card p-6"
        >
          <div className="flex items-center justify-between mb-5">
            <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>Application Pipeline</h3>
            <Link to="/applications" className="text-xs flex items-center gap-1" style={{ color: 'rgb(129 140 248)' }}>
              View all <ChevronRight size={14} />
            </Link>
          </div>
          <div className="space-y-3">
            {applicationStatuses.map(s => (
              <div key={s.status} className="flex items-center gap-3">
                <div className="w-24 text-xs capitalize" style={{ color: 'rgb(148 163 184)' }}>{s.label}</div>
                <div className="flex-1 h-2 rounded-full" style={{ background: 'rgba(45 45 65 / 0.8)' }}>
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${Math.min((s.value / Math.max(...applicationStatuses.map(x => x.value), 1)) * 100, 100)}%` }}
                    transition={{ duration: 0.8, delay: 0.5 }}
                    className="h-full rounded-full"
                    style={{ background: s.color }}
                  />
                </div>
                <div className="w-6 text-right text-sm font-bold" style={{ color: 'rgb(248 250 252)' }}>{s.value}</div>
              </div>
            ))}
          </div>
        </motion.div>

        {/* Upcoming Interviews */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="glass-card p-6"
        >
          <div className="flex items-center justify-between mb-5">
            <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>Upcoming Interviews</h3>
            <Link to="/interviews" className="text-xs flex items-center gap-1" style={{ color: 'rgb(129 140 248)' }}>
              View all <ChevronRight size={14} />
            </Link>
          </div>
          {!upcoming || upcoming.length === 0 ? (
            <div className="text-center py-8">
              <Calendar size={32} className="mx-auto mb-3" style={{ color: 'rgb(45 45 65)' }} />
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>No upcoming interviews</p>
              <Link to="/interviews" className="text-xs mt-2 inline-block" style={{ color: 'rgb(129 140 248)' }}>
                Add interview →
              </Link>
            </div>
          ) : (
            <div className="space-y-3">
              {upcoming.slice(0, 4).map((round: any, i: number) => (
                <div key={round.id} className="flex items-start gap-3 p-3 rounded-lg"
                  style={{ background: 'rgba(28 28 42 / 0.6)' }}>
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center text-xs font-bold"
                    style={{ background: 'rgba(99 102 241 / 0.15)', color: 'rgb(129 140 248)' }}>
                    {i + 1}
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium" style={{ color: 'rgb(248 250 252)' }}>{round.stage?.replace(/_/g, ' ')}</p>
                    <div className="flex items-center gap-1 mt-1">
                      <Clock size={12} style={{ color: 'rgb(71 85 105)' }} />
                      <span className="text-xs" style={{ color: 'rgb(71 85 105)' }}>
                        {round.scheduled_at ? new Date(round.scheduled_at).toLocaleDateString() : 'Not scheduled'}
                      </span>
                    </div>
                  </div>
                  <span className="badge badge-yellow text-xs capitalize">{round.result}</span>
                </div>
              ))}
            </div>
          )}
        </motion.div>
      </div>
    </PageLayout>
  )
}
