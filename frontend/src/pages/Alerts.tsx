import { useState } from 'react'
import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useKeywordAlerts, useCreateAlert, useDeleteAlert } from '@/api/hooks'
import { useForm } from 'react-hook-form'
import { Zap, Plus, Trash2, Bell, BellOff } from 'lucide-react'

export default function AlertsPage() {
  const { data: alerts, isLoading } = useKeywordAlerts()
  const createAlert = useCreateAlert()
  const deleteAlert = useDeleteAlert()
  const [showForm, setShowForm] = useState(false)
  const { register, handleSubmit, reset } = useForm()

  const onSubmit = (data: any) => {
    createAlert.mutate({ ...data, email_notify: true, browser_notify: true }, {
      onSuccess: () => { reset(); setShowForm(false) }
    })
  }

  return (
    <PageLayout title="Keyword Alerts">
      <div className="flex items-center justify-between mb-6">
        <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
          Get notified when new jobs match your keywords. Checked every 6 hours.
        </p>
        <button id="add-alert-btn" className="btn-primary flex items-center gap-2" onClick={() => setShowForm(v => !v)}>
          <Plus size={16} /> Add Alert
        </button>
      </div>

      {showForm && (
        <motion.div
          initial={{ opacity: 0, y: -8 }}
          animate={{ opacity: 1, y: 0 }}
          className="glass-card p-5 mb-6"
        >
          <form onSubmit={handleSubmit(onSubmit)} className="flex items-end gap-3">
            <div className="flex-1">
              <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Keyword or Phrase</label>
              <div className="relative">
                <Zap size={14} className="absolute left-3 top-1/2 -translate-y-1/2" style={{ color: 'rgb(250 204 21)' }} />
                <input
                  id="alert-keyword"
                  className="input-field pl-9"
                  placeholder='e.g. "Senior Go Engineer" or "React Native"'
                  {...register('keyword', { required: true, minLength: 2 })}
                />
              </div>
            </div>
            <button type="submit" className="btn-primary" disabled={createAlert.isPending}>
              {createAlert.isPending ? 'Adding...' : 'Create Alert'}
            </button>
            <button type="button" className="btn-secondary" onClick={() => setShowForm(false)}>Cancel</button>
          </form>
        </motion.div>
      )}

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="glass-card p-4 animate-pulse h-16" />
          ))}
        </div>
      ) : !alerts || alerts.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Zap size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No alerts set up</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Add keywords to get notified when matching jobs appear.
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {alerts.map((alert: any, i: number) => (
            <motion.div
              key={alert.id}
              initial={{ opacity: 0, x: -8 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.04 }}
              className="glass-card p-4 flex items-center gap-4"
              style={{ opacity: alert.is_active ? 1 : 0.5 }}
            >
              <div className="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
                style={{ background: 'rgba(250 204 21 / 0.12)' }}>
                <Zap size={16} style={{ color: 'rgb(250 204 21)' }} />
              </div>

              <div className="flex-1">
                <p className="font-semibold text-sm" style={{ color: 'rgb(248 250 252)' }}>
                  "{alert.keyword}"
                </p>
                <div className="flex items-center gap-3 mt-0.5">
                  {alert.match_count > 0 && (
                    <span className="text-xs" style={{ color: 'rgb(34 197 94)' }}>
                      {alert.match_count} matches
                    </span>
                  )}
                  {alert.last_matched_at && (
                    <span className="text-xs" style={{ color: 'rgb(71 85 105)' }}>
                      Last: {new Date(alert.last_matched_at).toLocaleDateString()}
                    </span>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-3 shrink-0">
                {alert.email_notify ? (
                  <Bell size={15} style={{ color: 'rgb(129 140 248)' }} aria-label="Email on" />
                ) : (
                  <BellOff size={15} style={{ color: 'rgb(71 85 105)' }} aria-label="Email off" />
                )}
                <span className={`badge text-xs ${alert.is_active ? 'badge-green' : ''}`}
                  style={!alert.is_active ? { background: 'rgba(71 85 105 / 0.1)', color: 'rgb(71 85 105)', border: '1px solid rgba(71 85 105 / 0.2)' } : {}}>
                  {alert.is_active ? 'Active' : 'Paused'}
                </span>
                <button
                  onClick={() => deleteAlert.mutate(alert.id)}
                  className="p-1.5 rounded-lg"
                  style={{ color: 'rgb(71 85 105)' }}
                >
                  <Trash2 size={14} />
                </button>
              </div>
            </motion.div>
          ))}
        </div>
      )}
    </PageLayout>
  )
}
