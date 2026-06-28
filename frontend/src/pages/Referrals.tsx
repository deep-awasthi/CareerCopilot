import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useReferrals, useCreateReferral, useUpdateReferral } from '@/api/hooks'
import { useForm } from 'react-hook-form'
import { Users, Plus, ChevronDown, Building2, ExternalLink } from 'lucide-react'

const STATUSES = [
  { value: 'not_contacted', label: 'Not Contacted', color: 'rgb(100 116 139)', bg: 'rgba(100 116 139 / 0.1)' },
  { value: 'contacted', label: 'Contacted', color: 'rgb(59 130 246)', bg: 'rgba(59 130 246 / 0.1)' },
  { value: 'follow_up', label: 'Follow Up', color: 'rgb(250 204 21)', bg: 'rgba(250 204 21 / 0.1)' },
  { value: 'referral_received', label: 'Referral Received', color: 'rgb(34 197 94)', bg: 'rgba(34 197 94 / 0.1)' },
  { value: 'applied', label: 'Applied', color: 'rgb(129 140 248)', bg: 'rgba(129 140 248 / 0.1)' },
  { value: 'rejected', label: 'Rejected', color: 'rgb(239 68 68)', bg: 'rgba(239 68 68 / 0.1)' },
]

export default function ReferralsPage() {
  const [showCreate, setShowCreate] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const { data } = useReferrals({ status: statusFilter })
  const createReferral = useCreateReferral()
  const updateReferral = useUpdateReferral()
  const { register, handleSubmit, reset } = useForm()

  const referrals = data?.data ?? []

  const onSubmit = (d: any) => {
    createReferral.mutate(d, { onSuccess: () => { reset(); setShowCreate(false) } })
  }

  const getStatus = (val: string) => STATUSES.find(s => s.value === val) ?? STATUSES[0]

  return (
    <PageLayout title="Referral Finder">
      <div className="flex items-center justify-between mb-6">
        <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
          Track referral opportunities from LinkedIn connections
        </p>
        <button id="add-referral-btn" className="btn-primary flex items-center gap-2" onClick={() => setShowCreate(true)}>
          <Plus size={16} /> Add Referral
        </button>
      </div>

      {/* Status filters */}
      <div className="flex gap-2 mb-6 flex-wrap">
        <button
          onClick={() => setStatusFilter('')}
          className="px-3 py-1.5 rounded-full text-xs font-medium transition-all"
          style={{
            background: statusFilter === '' ? 'rgba(99 102 241 / 0.15)' : 'rgba(20 20 30 / 0.8)',
            color: statusFilter === '' ? 'rgb(129 140 248)' : 'rgb(71 85 105)',
            border: `1px solid ${statusFilter === '' ? 'rgba(99 102 241 / 0.3)' : 'rgba(45 45 65 / 0.5)'}`,
          }}
        >All</button>
        {STATUSES.map(s => (
          <button
            key={s.value}
            onClick={() => setStatusFilter(s.value)}
            className="px-3 py-1.5 rounded-full text-xs font-medium transition-all"
            style={{
              background: statusFilter === s.value ? s.bg : 'rgba(20 20 30 / 0.8)',
              color: statusFilter === s.value ? s.color : 'rgb(71 85 105)',
              border: `1px solid ${statusFilter === s.value ? s.color + '40' : 'rgba(45 45 65 / 0.5)'}`,
            }}
          >
            {s.label}
          </button>
        ))}
      </div>

      {/* Create modal */}
      <AnimatePresence>
        {showCreate && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 flex items-center justify-center"
            style={{ background: 'rgba(0 0 0 / 0.7)', backdropFilter: 'blur(4px)' }}
          >
            <motion.div
              initial={{ scale: 0.95, y: 10 }}
              animate={{ scale: 1, y: 0 }}
              exit={{ scale: 0.95 }}
              className="glass-card p-8 w-full max-w-lg mx-4"
            >
              <h3 className="text-lg font-bold mb-6" style={{ color: 'rgb(248 250 252)' }}>Add Referral Contact</h3>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Name *</label>
                    <input className="input-field" placeholder="John Doe" {...register('referrer_name', { required: true })} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Designation</label>
                    <input className="input-field" placeholder="Senior Engineer" {...register('referrer_designation')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Department</label>
                    <input className="input-field" placeholder="Engineering" {...register('referrer_department')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Office Location</label>
                    <input className="input-field" placeholder="Bangalore" {...register('referrer_office_location')} />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>LinkedIn Profile URL</label>
                  <div className="relative">
                    <ExternalLink size={14} className="absolute left-3 top-1/2 -translate-y-1/2" style={{ color: 'rgb(59 130 246)' }} />
                    <input className="input-field pl-9" placeholder="https://linkedin.com/in/johndoe" {...register('referrer_profile_url')} />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Notes</label>
                  <textarea className="input-field" rows={2} placeholder="Connection notes..." {...register('notes')} />
                </div>
                <div className="flex gap-3">
                  <button type="submit" className="btn-primary flex-1" disabled={createReferral.isPending}>
                    {createReferral.isPending ? 'Saving...' : 'Save Contact'}
                  </button>
                  <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                </div>
              </form>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Referral cards */}
      {referrals.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Users size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No referrals tracked</p>
          <p className="text-sm mb-4" style={{ color: 'rgb(71 85 105)' }}>
            Add LinkedIn connections who work at companies you're interested in.
          </p>
          <button className="btn-primary" onClick={() => setShowCreate(true)}>
            <Plus size={16} className="inline mr-2" /> Add First Contact
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {referrals.map((ref: any, i: number) => {
            const status = getStatus(ref.status)
            return (
              <motion.div
                key={ref.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.04 }}
                className="glass-card p-5"
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold"
                      style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)', color: 'white' }}>
                      {ref.referrer_name?.[0]?.toUpperCase()}
                    </div>
                    <div>
                      <p className="font-semibold text-sm" style={{ color: 'rgb(248 250 252)' }}>{ref.referrer_name}</p>
                      <p className="text-xs" style={{ color: 'rgb(100 116 139)' }}>{ref.referrer_designation}</p>
                    </div>
                  </div>
                  <span className="badge text-xs" style={{ background: status.bg, color: status.color, border: `1px solid ${status.color}40` }}>
                    {status.label}
                  </span>
                </div>

                {(ref.referrer_department || ref.referrer_office_location) && (
                  <div className="flex items-center gap-1 text-xs mb-3" style={{ color: 'rgb(71 85 105)' }}>
                    <Building2 size={11} />
                    {[ref.referrer_department, ref.referrer_office_location].filter(Boolean).join(' · ')}
                  </div>
                )}

                {ref.notes && (
                  <p className="text-xs mb-3 line-clamp-2" style={{ color: 'rgb(100 116 139)' }}>{ref.notes}</p>
                )}

                <div className="flex items-center gap-2 mt-3">
                  {ref.referrer_profile_url && (
                    <a
                      href={ref.referrer_profile_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs px-3 py-1.5 rounded-lg"
                      style={{ background: 'rgba(59 130 246 / 0.1)', color: 'rgb(59 130 246)', border: '1px solid rgba(59 130 246 / 0.2)' }}
                    >
                      <ExternalLink size={12} /> LinkedIn
                    </a>
                  )}
                  <div className="relative ml-auto">
                    <select
                      value={ref.status}
                      onChange={e => updateReferral.mutate({ id: ref.id, ...({ status: e.target.value } as any) })}
                      className="appearance-none text-xs px-3 py-1.5 rounded-lg cursor-pointer border-0 outline-none pr-6"
                      style={{ background: 'rgba(28 28 42 / 0.8)', color: 'rgb(148 163 184)', border: '1px solid rgba(45 45 65 / 0.5)' }}
                    >
                      {STATUSES.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}
                    </select>
                    <ChevronDown size={11} className="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none" style={{ color: 'rgb(71 85 105)' }} />
                  </div>
                </div>
              </motion.div>
            )
          })}
        </div>
      )}
    </PageLayout>
  )
}
