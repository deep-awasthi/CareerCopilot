import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useBookmarks, useUnbookmark } from '@/api/hooks'
import { Bookmark, X, Briefcase, Building2, Users } from 'lucide-react'
import { useState } from 'react'

const TYPE_ICONS: Record<string, any> = { job: Briefcase, company: Building2, referral: Users }

export default function BookmarksPage() {
  const [activeType, setActiveType] = useState('')
  const { data, isLoading } = useBookmarks(activeType || undefined)
  const unbookmark = useUnbookmark()

  const bookmarks = data?.data ?? []

  return (
    <PageLayout title="Bookmarks">
      <div className="flex gap-3 mb-6">
        {['', 'job', 'company', 'referral'].map(type => (
          <button
            key={type}
            onClick={() => setActiveType(type)}
            className={`px-4 py-2 rounded-full text-sm font-medium transition-all capitalize`}
            style={{
              background: activeType === type ? 'rgba(99 102 241 / 0.15)' : 'rgba(20 20 30 / 0.8)',
              color: activeType === type ? 'rgb(129 140 248)' : 'rgb(71 85 105)',
              border: `1px solid ${activeType === type ? 'rgba(99 102 241 / 0.3)' : 'rgba(45 45 65 / 0.5)'}`,
            }}
          >
            {type || 'All'}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => <div key={i} className="glass-card p-5 animate-pulse h-24" />)}
        </div>
      ) : bookmarks.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Bookmark size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No bookmarks</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Bookmark jobs, companies, or referrals from their respective pages.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {bookmarks.map((bm: any, i: number) => {
            const Icon = TYPE_ICONS[bm.type] ?? Briefcase
            return (
              <motion.div
                key={bm.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.04 }}
                className="glass-card p-5 group relative"
              >
                <button
                  onClick={() => unbookmark.mutate({ type: bm.type, id: bm.target_id })}
                  className="absolute top-3 right-3 p-1 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity"
                  style={{ color: 'rgb(239 68 68)', background: 'rgba(239 68 68 / 0.1)' }}
                >
                  <X size={14} />
                </button>
                <div className="flex items-center gap-3">
                  <div className="w-9 h-9 rounded-xl flex items-center justify-center"
                    style={{ background: 'rgba(99 102 241 / 0.12)' }}>
                    <Icon size={16} style={{ color: 'rgb(129 140 248)' }} />
                  </div>
                  <div>
                    <p className="font-medium text-sm capitalize" style={{ color: 'rgb(248 250 252)' }}>
                      {bm.type} #{bm.target_id}
                    </p>
                    <p className="text-xs" style={{ color: 'rgb(71 85 105)' }}>
                      {new Date(bm.created_at).toLocaleDateString()}
                    </p>
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
