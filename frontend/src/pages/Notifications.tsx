import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useNotifications, useMarkRead, useMarkAllRead, useUnreadCount } from '@/api/hooks'
import { Bell, CheckCheck, Trash2 } from 'lucide-react'

const TYPE_STYLES: Record<string, { color: string; bg: string; emoji: string }> = {
  new_job: { color: 'rgb(129 140 248)', bg: 'rgba(99 102 241 / 0.12)', emoji: '💼' },
  keyword_match: { color: 'rgb(250 204 21)', bg: 'rgba(250 204 21 / 0.12)', emoji: '🔔' },
  company_opening: { color: 'rgb(34 197 94)', bg: 'rgba(34 197 94 / 0.12)', emoji: '🏢' },
  referral_update: { color: 'rgb(251 113 133)', bg: 'rgba(251 113 133 / 0.12)', emoji: '🤝' },
  interview_reminder: { color: 'rgb(96 165 250)', bg: 'rgba(96 165 250 / 0.12)', emoji: '📅' },
  daily_digest: { color: 'rgb(167 139 250)', bg: 'rgba(167 139 250 / 0.12)', emoji: '📊' },
}

export default function NotificationsPage() {
  const { data } = useNotifications({})
  const { data: unread } = useUnreadCount()
  const markRead = useMarkRead()
  const markAllRead = useMarkAllRead()

  const notifications = data?.data ?? []
  const unreadCount = unread?.count ?? 0

  return (
    <PageLayout title="Notifications">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            {unreadCount > 0 ? (
              <span style={{ color: 'rgb(129 140 248)' }}>{unreadCount} unread</span>
            ) : 'All caught up!'}
          </p>
        </div>
        {unreadCount > 0 && (
          <button
            id="mark-all-read-btn"
            className="btn-secondary flex items-center gap-2 text-sm"
            onClick={() => markAllRead.mutate()}
          >
            <CheckCheck size={14} /> Mark all read
          </button>
        )}
      </div>

      {notifications.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Bell size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No notifications yet</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Notifications will appear here when jobs match your alerts, interviews are scheduled, etc.
          </p>
        </div>
      ) : (
        <div className="space-y-2">
          {notifications.map((n: any, i: number) => {
            const style = TYPE_STYLES[n.type] ?? TYPE_STYLES.daily_digest
            return (
              <motion.div
                key={n.id}
                initial={{ opacity: 0, x: -8 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.03 }}
                className="glass-card p-4 flex items-start gap-4 cursor-pointer"
                style={{ opacity: n.is_read ? 0.65 : 1 }}
                onClick={() => !n.is_read && markRead.mutate(n.id)}
              >
                <div className="w-9 h-9 rounded-xl flex items-center justify-center text-lg shrink-0"
                  style={{ background: style.bg }}>
                  {style.emoji}
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-start justify-between gap-2">
                    <p className="font-medium text-sm" style={{ color: 'rgb(248 250 252)' }}>{n.title}</p>
                    {!n.is_read && (
                      <div className="w-2 h-2 rounded-full shrink-0 mt-1.5" style={{ background: 'rgb(99 102 241)' }} />
                    )}
                  </div>
                  <p className="text-xs mt-0.5 line-clamp-2" style={{ color: 'rgb(100 116 139)' }}>{n.body}</p>
                  <p className="text-xs mt-1" style={{ color: 'rgb(71 85 105)' }}>
                    {n.created_at ? new Date(n.created_at).toLocaleString() : ''}
                  </p>
                </div>
              </motion.div>
            )
          })}
        </div>
      )}
    </PageLayout>
  )
}
