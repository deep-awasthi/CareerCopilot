import { useState } from 'react'
import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useCompanies, useWatchCompany } from '@/api/hooks'
import { Search, Building2, Globe, Eye, Star } from 'lucide-react'

export default function CompaniesPage() {
  const [q, setQ] = useState('')
  const { data, isLoading } = useCompanies({ q })
  const watchCompany = useWatchCompany()

  const companies = data?.data ?? []

  return (
    <PageLayout title="Companies">
      <div className="flex gap-3 mb-6">
        <div className="relative flex-1 max-w-md">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2" style={{ color: 'rgb(71 85 105)' }} />
          <input
            id="company-search"
            className="input-field pl-9"
            placeholder="Search companies..."
            value={q}
            onChange={e => setQ(e.target.value)}
          />
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 9 }).map((_, i) => (
            <div key={i} className="glass-card p-5 animate-pulse">
              <div className="h-5 rounded mb-3" style={{ background: 'rgba(45 45 65 / 0.5)', width: '60%' }} />
              <div className="h-3 rounded" style={{ background: 'rgba(45 45 65 / 0.3)', width: '40%' }} />
            </div>
          ))}
        </div>
      ) : companies.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Building2 size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No companies found</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>Companies are populated automatically when jobs are scraped.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {companies.map((co: any, i: number) => (
            <motion.div
              key={co.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.04 }}
              className="glass-card p-5 flex flex-col"
            >
              <div className="flex items-start gap-3 mb-3">
                {co.logo_url ? (
                  <img src={co.logo_url} alt={co.name} className="w-10 h-10 rounded-lg object-contain"
                    style={{ background: 'rgba(255 255 255 / 0.05)' }} />
                ) : (
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center font-bold text-sm"
                    style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)', color: 'white' }}>
                    {co.name?.[0]}
                  </div>
                )}
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold text-sm truncate" style={{ color: 'rgb(248 250 252)' }}>{co.name}</h3>
                  <p className="text-xs" style={{ color: 'rgb(100 116 139)' }}>{co.industry} · {co.size}</p>
                </div>
              </div>

              {co.headquarters && (
                <p className="text-xs mb-3" style={{ color: 'rgb(71 85 105)' }}>📍 {co.headquarters}</p>
              )}

              {co.description && (
                <p className="text-xs mb-3 line-clamp-2 flex-1" style={{ color: 'rgb(100 116 139)' }}>{co.description}</p>
              )}

              <div className="flex items-center gap-2 mt-auto">
                {co.domain && (
                  <a
                    href={`https://${co.domain}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-xs px-3 py-1.5 rounded-lg"
                    style={{ background: 'rgba(45 45 65 / 0.5)', color: 'rgb(148 163 184)' }}
                  >
                    <Globe size={12} /> Website
                  </a>
                )}
                {co.career_page_url && (
                  <a
                    href={co.career_page_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-xs px-3 py-1.5 rounded-lg"
                    style={{ background: 'rgba(45 45 65 / 0.5)', color: 'rgb(148 163 184)' }}
                  >
                    <Eye size={12} /> Careers
                  </a>
                )}
                <button
                  id={`watch-company-${co.id}`}
                  className="ml-auto flex items-center gap-1 text-xs px-3 py-1.5 rounded-lg"
                  style={{ background: 'rgba(99 102 241 / 0.12)', color: 'rgb(129 140 248)', border: '1px solid rgba(99 102 241 / 0.2)' }}
                  onClick={() => watchCompany.mutate(co.id)}
                >
                  <Star size={12} /> Watch
                </button>
              </div>
            </motion.div>
          ))}
        </div>
      )}
    </PageLayout>
  )
}
