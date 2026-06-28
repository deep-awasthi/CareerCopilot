import { useState } from 'react'
import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useApplications, useUpdateApplication, useDeleteApplication } from '@/api/hooks'
import {
  ClipboardList, Briefcase, ChevronDown, Trash2,
  CheckCircle2, Clock, Star, XCircle, Archive, BookmarkIcon
} from 'lucide-react'

const STATUSES = [
  { value: '', label: 'All', color: 'rgb(100 116 139)', bg: 'rgba(100 116 139 / 0.1)' },
  { value: 'saved', label: 'Saved', color: 'rgb(100 116 139)', bg: 'rgba(100 116 139 / 0.1)', icon: BookmarkIcon },
  { value: 'applied', label: 'Applied', color: 'rgb(59 130 246)', bg: 'rgba(59 130 246 / 0.1)', icon: CheckCircle2 },
  { value: 'interview', label: 'Interview', color: 'rgb(250 204 21)', bg: 'rgba(250 204 21 / 0.1)', icon: Clock },
  { value: 'offer', label: 'Offer', color: 'rgb(34 197 94)', bg: 'rgba(34 197 94 / 0.1)', icon: Star },
  { value: 'rejected', label: 'Rejected', color: 'rgb(239 68 68)', bg: 'rgba(239 68 68 / 0.1)', icon: XCircle },
  { value: 'archived', label: 'Archived', color: 'rgb(71 85 105)', bg: 'rgba(71 85 105 / 0.1)', icon: Archive },
]

export default function ApplicationsPage() {
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const { data, isLoading } = useApplications({ status: statusFilter, page, per_page: 20 })
  const updateApp = useUpdateApplication()
  const deleteApp = useDeleteApplication()

  const applications = data?.data ?? []
  const total = data?.total ?? 0

  const changeStatus = (id: number, status: string) => {
    updateApp.mutate({ id, status })
  }

  return (
    <PageLayout title="Applications">
      {/* Status filter bar */}
      <div className="flex gap-2 mb-6 flex-wrap">
        {STATUSES.map(s => (
          <button
            key={s.value}
            id={`filter-${s.value || 'all'}`}
            onClick={() => { setStatusFilter(s.value); setPage(1) }}
            className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all"
            style={{
              background: statusFilter === s.value ? s.bg : 'rgba(20 20 30 / 0.8)',
              color: statusFilter === s.value ? s.color : 'rgb(71 85 105)',
              border: `1px solid ${statusFilter === s.value ? s.color + '40' : 'rgba(45 45 65 / 0.5)'}`,
            }}
          >
            {s.icon && <s.icon size={13} />}
            {s.label}
          </button>
        ))}
        <div className="ml-auto text-sm flex items-center" style={{ color: 'rgb(71 85 105)' }}>
          {total} total
        </div>
      </div>

      {/* Applications table */}
      {isLoading ? (
        <div className="glass-card p-8 text-center">
          <div className="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto" />
        </div>
      ) : applications.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <ClipboardList size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No applications yet</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Save or apply to jobs to start tracking them here.
          </p>
        </div>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="data-table">
            <thead>
              <tr>
                <th>Job</th>
                <th>Status</th>
                <th>Applied</th>
                <th>Follow-Up</th>
                <th>Salary Offered</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {applications.map((app: any, i: number) => {
                const status = STATUSES.find(s => s.value === app.status) ?? STATUSES[0]
                return (
                  <motion.tr
                    key={app.id}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: i * 0.03 }}
                  >
                    <td>
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-lg flex items-center justify-center text-xs font-bold"
                          style={{ background: 'rgba(99 102 241 / 0.15)', color: 'rgb(129 140 248)' }}>
                          <Briefcase size={14} />
                        </div>
                        <div>
                          <p className="font-medium text-sm" style={{ color: 'rgb(248 250 252)' }}>
                            Job #{app.job_id}
                          </p>
                          {app.referral_used && (
                            <span className="text-xs" style={{ color: 'rgb(250 204 21)' }}>via referral</span>
                          )}
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="relative">
                        <select
                          value={app.status}
                          onChange={e => changeStatus(app.id, e.target.value)}
                          className="appearance-none text-xs font-semibold px-3 py-1.5 rounded-full pr-7 cursor-pointer border-0 outline-none"
                          style={{ background: status.bg, color: status.color }}
                        >
                          {STATUSES.slice(1).map(s => (
                            <option key={s.value} value={s.value}>{s.label}</option>
                          ))}
                        </select>
                        <ChevronDown size={12} className="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none" style={{ color: status.color }} />
                      </div>
                    </td>
                    <td>
                      <span className="text-sm">
                        {app.applied_at ? new Date(app.applied_at).toLocaleDateString() : '—'}
                      </span>
                    </td>
                    <td>
                      <span className="text-sm">
                        {app.follow_up_date ? new Date(app.follow_up_date).toLocaleDateString() : '—'}
                      </span>
                    </td>
                    <td>
                      {app.salary_offered > 0 ? (
                        <span className="font-medium" style={{ color: 'rgb(34 197 94)' }}>
                          ₹{(app.salary_offered / 100000).toFixed(1)}L
                        </span>
                      ) : '—'}
                    </td>
                    <td>
                      <button
                        onClick={() => deleteApp.mutate(app.id)}
                        className="p-1.5 rounded-lg transition-colors"
                        style={{ color: 'rgb(71 85 105)' }}
                      >
                        <Trash2 size={14} />
                      </button>
                    </td>
                  </motion.tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}
    </PageLayout>
  )
}
