import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useWatchlist, useUnwatchCompany } from '@/api/hooks'
import { Target, X, Globe, Eye, Building2 } from 'lucide-react'

export default function WatchlistsPage() {
  const { data, isLoading } = useWatchlist()
  const unwatch = useUnwatchCompany()

  const watchlists = data?.data ?? []

  return (
    <PageLayout title="Watchlists">
      <p className="text-sm mb-6" style={{ color: 'rgb(71 85 105)' }}>
        Get notified when these companies post new jobs.
      </p>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="glass-card p-5 animate-pulse h-28" />
          ))}
        </div>
      ) : watchlists.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Target size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No companies in watchlist</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Go to <strong>Companies</strong> and click "Watch" to add companies here.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {watchlists.map((wl: any, i: number) => {
            const co = wl.company
            return (
              <motion.div
                key={wl.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.04 }}
                className="glass-card p-5 group relative"
              >
                <button
                  id={`unwatch-${wl.id}`}
                  onClick={() => unwatch.mutate(wl.company_id)}
                  className="absolute top-3 right-3 p-1 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity"
                  style={{ color: 'rgb(239 68 68)', background: 'rgba(239 68 68 / 0.1)' }}
                >
                  <X size={14} />
                </button>

                <div className="flex items-center gap-3 mb-3">
                  {co?.logo_url ? (
                    <img src={co.logo_url} alt={co?.name} className="w-9 h-9 rounded-lg" />
                  ) : (
                    <div className="w-9 h-9 rounded-lg flex items-center justify-center font-bold text-sm"
                      style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)', color: 'white' }}>
                      {co?.name?.[0] ?? <Building2 size={14} />}
                    </div>
                  )}
                  <div>
                    <p className="font-semibold text-sm" style={{ color: 'rgb(248 250 252)' }}>{co?.name ?? 'Unknown'}</p>
                    <p className="text-xs" style={{ color: 'rgb(100 116 139)' }}>{co?.industry}</p>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {wl.notify_new_jobs && (
                    <span className="badge badge-green text-xs">🔔 Alerts On</span>
                  )}
                  {wl.last_notified_at && (
                    <span className="text-xs ml-auto" style={{ color: 'rgb(71 85 105)' }}>
                      Last notified: {new Date(wl.last_notified_at).toLocaleDateString()}
                    </span>
                  )}
                </div>
              </motion.div>
            )
          })}
        </div>
      )}
    </PageLayout>
  )
}
