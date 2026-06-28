import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useAnalytics, useDashboardStats } from '@/api/hooks'
import {
  BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
  LineChart, Line, PieChart, Pie, Cell, Legend, CartesianGrid
} from 'recharts'
import { TrendingUp, Target, Star, Users } from 'lucide-react'

const COLORS = ['#6366f1', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#06b6d4']

const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null
  return (
    <div className="glass-card px-3 py-2 text-xs" style={{ color: 'rgb(248 250 252)' }}>
      <p className="font-semibold mb-1">{label}</p>
      {payload.map((p: any) => (
        <p key={p.name} style={{ color: p.color }}>{p.name}: {p.value}</p>
      ))}
    </div>
  )
}

export default function AnalyticsPage() {
  const { data: stats } = useDashboardStats()
  const { data: analytics } = useAnalytics()

  const kpiCards = [
    {
      label: 'Response Rate',
      value: `${(analytics?.response_rate ?? 0).toFixed(1)}%`,
      icon: TrendingUp,
      color: 'rgb(129 140 248)',
      bg: 'rgba(99 102 241 / 0.12)',
      desc: 'Applied → Interview',
    },
    {
      label: 'Interview Success',
      value: `${(analytics?.interview_success_rate ?? 0).toFixed(1)}%`,
      icon: Target,
      color: 'rgb(34 197 94)',
      bg: 'rgba(34 197 94 / 0.12)',
      desc: 'Rounds passed',
    },
    {
      label: 'Avg Offer Salary',
      value: analytics?.average_salary > 0 ? `₹${(analytics.average_salary / 100000).toFixed(1)}L` : '—',
      icon: Star,
      color: 'rgb(250 204 21)',
      bg: 'rgba(250 204 21 / 0.12)',
      desc: 'From accepted offers',
    },
    {
      label: 'Referral Success',
      value: `${(analytics?.referral_success_rate ?? 0).toFixed(1)}%`,
      icon: Users,
      color: 'rgb(251 113 133)',
      bg: 'rgba(251 113 133 / 0.12)',
      desc: 'Got referral received',
    },
  ]

  return (
    <PageLayout title="Analytics">
      {/* KPI cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {kpiCards.map((card, i) => (
          <motion.div
            key={card.label}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.06 }}
            className="stat-card flex items-start justify-between"
          >
            <div>
              <p className="text-xs font-semibold uppercase tracking-wider mb-1" style={{ color: 'rgb(71 85 105)' }}>
                {card.label}
              </p>
              <p className="text-3xl font-bold" style={{ color: 'rgb(248 250 252)' }}>{card.value}</p>
              <p className="text-xs mt-1" style={{ color: 'rgb(71 85 105)' }}>{card.desc}</p>
            </div>
            <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ background: card.bg }}>
              <card.icon size={18} style={{ color: card.color }} />
            </div>
          </motion.div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Applications per month */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="glass-card p-6"
        >
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Applications per Month</h3>
          {analytics?.applications_per_month?.length > 0 ? (
            <ResponsiveContainer width="100%" height={220}>
              <BarChart data={analytics.applications_per_month}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(45 45 65 / 0.5)" />
                <XAxis dataKey="month" tick={{ fontSize: 11, fill: 'rgb(71 85 105)' }} />
                <YAxis tick={{ fontSize: 11, fill: 'rgb(71 85 105)' }} />
                <Tooltip content={<CustomTooltip />} />
                <Bar dataKey="count" fill="#6366f1" radius={[4, 4, 0, 0]} name="Applications" />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-40 flex items-center justify-center">
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>No application data yet</p>
            </div>
          )}
        </motion.div>

        {/* Top Skills Demand */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="glass-card p-6"
        >
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Top In-Demand Skills</h3>
          {analytics?.top_skills?.length > 0 ? (
            <ResponsiveContainer width="100%" height={220}>
              <BarChart data={analytics.top_skills.slice(0, 10)} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(45 45 65 / 0.5)" />
                <XAxis type="number" tick={{ fontSize: 11, fill: 'rgb(71 85 105)' }} />
                <YAxis dataKey="skill" type="category" tick={{ fontSize: 11, fill: 'rgb(148 163 184)' }} width={100} />
                <Tooltip content={<CustomTooltip />} />
                <Bar dataKey="count" fill="#8b5cf6" radius={[0, 4, 4, 0]} name="Jobs" />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-40 flex items-center justify-center">
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>No skill data yet</p>
            </div>
          )}
        </motion.div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top companies */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="glass-card p-6"
        >
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Companies Applied To</h3>
          {analytics?.top_companies?.length > 0 ? (
            <div className="space-y-3">
              {analytics.top_companies.slice(0, 8).map((co: any, i: number) => (
                <div key={co.company} className="flex items-center gap-3">
                  <span className="text-xs w-5 text-right font-bold" style={{ color: 'rgb(71 85 105)' }}>{i + 1}</span>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-medium" style={{ color: 'rgb(248 250 252)' }}>{co.company}</span>
                      <span className="text-xs" style={{ color: 'rgb(100 116 139)' }}>{co.count}</span>
                    </div>
                    <div className="h-1.5 rounded-full" style={{ background: 'rgba(45 45 65 / 0.6)' }}>
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${(co.count / analytics.top_companies[0].count) * 100}%` }}
                        transition={{ duration: 0.8, delay: 0.6 + i * 0.05 }}
                        className="h-full rounded-full"
                        style={{ background: COLORS[i % COLORS.length] }}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="h-40 flex items-center justify-center">
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>Apply to jobs to see company data</p>
            </div>
          )}
        </motion.div>

        {/* Application status pie */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="glass-card p-6"
        >
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Pipeline Breakdown</h3>
          {stats ? (
            <ResponsiveContainer width="100%" height={220}>
              <PieChart>
                <Pie
                  data={[
                    { name: 'Saved', value: stats.saved_jobs ?? 0 },
                    { name: 'Applied', value: stats.applied_jobs ?? 0 },
                    { name: 'Interview', value: stats.interviews ?? 0 },
                    { name: 'Offer', value: stats.offers ?? 0 },
                  ].filter(d => d.value > 0)}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={90}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {COLORS.map((c, i) => <Cell key={i} fill={c} />)}
                </Pie>
                <Tooltip content={<CustomTooltip />} />
                <Legend iconType="circle" iconSize={8} wrapperStyle={{ fontSize: 12, color: 'rgb(148 163 184)' }} />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-40 flex items-center justify-center">
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>No pipeline data yet</p>
            </div>
          )}
        </motion.div>
      </div>
    </PageLayout>
  )
}
