import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useJobs, useCreateApplication } from '@/api/hooks'
import {
  Search, MapPin, Clock, Briefcase, ExternalLink,
  ChevronLeft, ChevronRight, SlidersHorizontal, Bookmark,
  Building2, Wifi, Star, Filter, X
} from 'lucide-react'

const PROVIDERS = ['all', 'greenhouse', 'lever', 'linkedin', 'indeed', 'naukri', 'wellfound']
const TYPES = ['', 'full-time', 'part-time', 'contract', 'internship']

export default function JobsPage() {
  const [filters, setFilters] = useState({
    q: '', location: '', remote: undefined as boolean | undefined,
    per_page: 20, page: 1, provider: '', employment_type: '',
    exp_min: '', salary_min: '',
  })
  const [showFilters, setShowFilters] = useState(false)
  const [selectedJob, setSelectedJob] = useState<any>(null)

  const { data, isLoading, isFetching } = useJobs(
    Object.fromEntries(Object.entries(filters).filter(([, v]) => v !== '' && v !== undefined))
  )
  const createApplication = useCreateApplication()

  const jobs = data?.data ?? []
  const total = data?.total ?? 0
  const totalPages = Math.ceil(total / filters.per_page)

  const save = (job: any) => {
    createApplication.mutate({ job_id: job.id, status: 'saved' })
  }

  return (
    <PageLayout title="Job Discovery">
      {/* Search bar */}
      <div className="flex gap-3 mb-5">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2" style={{ color: 'rgb(71 85 105)' }} />
          <input
            id="job-search-q"
            className="input-field pl-9"
            placeholder="Search by title, skill, or keyword..."
            value={filters.q}
            onChange={e => setFilters(f => ({ ...f, q: e.target.value, page: 1 }))}
          />
        </div>
        <div className="relative">
          <MapPin size={16} className="absolute left-3 top-1/2 -translate-y-1/2" style={{ color: 'rgb(71 85 105)' }} />
          <input
            id="job-search-location"
            className="input-field pl-9 w-48"
            placeholder="Location"
            value={filters.location}
            onChange={e => setFilters(f => ({ ...f, location: e.target.value, page: 1 }))}
          />
        </div>
        <button
          id="job-filter-toggle"
          className="btn-secondary flex items-center gap-2"
          onClick={() => setShowFilters(v => !v)}
        >
          <Filter size={16} />
          Filters
          {showFilters && <X size={14} />}
        </button>
      </div>

      {/* Advanced filters */}
      <AnimatePresence>
        {showFilters && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="glass-card p-5 mb-5 overflow-hidden"
          >
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div>
                <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(71 85 105)' }}>Provider</label>
                <select
                  id="job-filter-provider"
                  className="input-field"
                  value={filters.provider}
                  onChange={e => setFilters(f => ({ ...f, provider: e.target.value, page: 1 }))}
                >
                  {PROVIDERS.map(p => <option key={p} value={p === 'all' ? '' : p}>{p}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(71 85 105)' }}>Job Type</label>
                <select
                  id="job-filter-type"
                  className="input-field"
                  value={filters.employment_type}
                  onChange={e => setFilters(f => ({ ...f, employment_type: e.target.value, page: 1 }))}
                >
                  {TYPES.map(t => <option key={t} value={t}>{t || 'Any type'}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(71 85 105)' }}>Min Experience (yrs)</label>
                <input
                  id="job-filter-exp"
                  type="number" min="0" max="30"
                  className="input-field"
                  placeholder="e.g. 2"
                  value={filters.exp_min}
                  onChange={e => setFilters(f => ({ ...f, exp_min: e.target.value, page: 1 }))}
                />
              </div>
              <div>
                <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(71 85 105)' }}>Min Salary</label>
                <input
                  id="job-filter-salary"
                  type="number" min="0"
                  className="input-field"
                  placeholder="e.g. 800000"
                  value={filters.salary_min}
                  onChange={e => setFilters(f => ({ ...f, salary_min: e.target.value, page: 1 }))}
                />
              </div>
              <label className="flex items-center gap-2 cursor-pointer col-span-2">
                <input
                  id="job-filter-remote"
                  type="checkbox"
                  checked={filters.remote === true}
                  onChange={e => setFilters(f => ({ ...f, remote: e.target.checked ? true : undefined, page: 1 }))}
                  className="w-4 h-4 rounded"
                  style={{ accentColor: 'rgb(99 102 241)' }}
                />
                <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>Remote only</span>
              </label>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Results bar */}
      <div className="flex items-center justify-between mb-4">
        <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
          {isLoading ? 'Loading...' : `${total.toLocaleString()} jobs found`}
        </p>
        <div className="flex gap-2">
          {PROVIDERS.slice(1).map(p => (
            <button
              key={p}
              onClick={() => setFilters(f => ({ ...f, provider: f.provider === p ? '' : p, page: 1 }))}
              className={`badge text-xs cursor-pointer ${filters.provider === p ? 'badge-purple' : ''}`}
              style={filters.provider !== p ? { background: 'rgba(30 30 45 / 0.8)', color: 'rgb(100 116 139)', border: '1px solid rgba(45 45 65 / 0.5)' } : {}}
            >
              {p}
            </button>
          ))}
        </div>
      </div>

      <div className="flex gap-6" style={{ minHeight: '60vh' }}>
        {/* Job list */}
        <div className="flex-1 space-y-3">
          {isLoading ? (
            Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="glass-card p-5 animate-pulse">
                <div className="h-5 rounded mb-3" style={{ background: 'rgba(45 45 65 / 0.5)', width: '60%' }} />
                <div className="h-3 rounded mb-2" style={{ background: 'rgba(45 45 65 / 0.3)', width: '40%' }} />
                <div className="h-3 rounded" style={{ background: 'rgba(45 45 65 / 0.2)', width: '80%' }} />
              </div>
            ))
          ) : jobs.length === 0 ? (
            <div className="glass-card p-16 text-center">
              <Briefcase size={48} className="mx-auto mb-4" style={{ color: 'rgb(45 45 65)' }} />
              <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No jobs found</p>
              <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>Try adjusting your search filters or check back later as new jobs are scraped every 6 hours.</p>
            </div>
          ) : (
            jobs.map((job: any, i: number) => (
              <motion.div
                key={job.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.02 }}
                onClick={() => setSelectedJob(job)}
                className="glass-card p-5 cursor-pointer transition-all"
                style={{
                  borderColor: selectedJob?.id === job.id ? 'rgba(99 102 241 / 0.5)' : undefined,
                  boxShadow: selectedJob?.id === job.id ? '0 0 0 1px rgba(99 102 241 / 0.3)' : undefined,
                }}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-base mb-1 truncate" style={{ color: 'rgb(248 250 252)' }}>
                      {job.title}
                    </h3>
                    <div className="flex items-center gap-3 flex-wrap">
                      {job.company && (
                        <span className="flex items-center gap-1 text-sm" style={{ color: 'rgb(148 163 184)' }}>
                          <Building2 size={13} /> {job.company}
                        </span>
                      )}
                      <span className="flex items-center gap-1 text-sm" style={{ color: 'rgb(100 116 139)' }}>
                        <MapPin size={13} /> {job.location || 'Location not specified'}
                      </span>
                      {job.is_remote && (
                        <span className="flex items-center gap-1 text-xs" style={{ color: 'rgb(34 197 94)' }}>
                          <Wifi size={12} /> Remote
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex flex-col items-end gap-2 shrink-0">
                    {job.match_score > 0 && (
                      <span className="badge text-xs font-semibold" style={{ background: 'rgba(99 102 241 / 0.15)', color: 'rgb(129 140 248)', border: '1px solid rgba(99 102 241 / 0.3)' }}>
                        {job.match_score}% Match
                      </span>
                    )}
                    {job.salary_max > 0 && (
                      <span className="text-xs font-medium" style={{ color: 'rgb(250 204 21)' }}>
                        ₹{(job.salary_min / 100000).toFixed(0)}–{(job.salary_max / 100000).toFixed(0)}L
                      </span>
                    )}
                    <button
                      id={`save-job-${job.id}`}
                      onClick={e => { e.stopPropagation(); save(job) }}
                      className="badge badge-purple text-xs"
                    >
                      <Bookmark size={11} className="mr-1" /> Save
                    </button>
                  </div>
                </div>
                <div className="flex items-center gap-2 mt-3 flex-wrap">
                  {job.skills?.slice(0, 5).map((skill: string) => (
                    <span key={skill} className="badge" style={{ background: 'rgba(30 30 45 / 0.8)', color: 'rgb(100 116 139)', border: '1px solid rgba(45 45 65 / 0.5)', fontSize: 11 }}>
                      {skill}
                    </span>
                  ))}
                  {job.skills?.length > 5 && (
                    <span className="text-xs" style={{ color: 'rgb(71 85 105)' }}>+{job.skills.length - 5}</span>
                  )}
                  <div className="ml-auto flex items-center gap-1 text-xs" style={{ color: 'rgb(71 85 105)' }}>
                    <Clock size={11} />
                    {job.posted_at ? new Date(job.posted_at).toLocaleDateString() : 'Recently'}
                  </div>
                </div>
              </motion.div>
            ))
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-3 mt-6">
              <button
                className="btn-secondary flex items-center gap-1"
                disabled={filters.page <= 1}
                onClick={() => setFilters(f => ({ ...f, page: f.page - 1 }))}
              >
                <ChevronLeft size={16} /> Prev
              </button>
              <span className="text-sm" style={{ color: 'rgb(100 116 139)' }}>
                Page {filters.page} of {totalPages}
              </span>
              <button
                className="btn-secondary flex items-center gap-1"
                disabled={filters.page >= totalPages}
                onClick={() => setFilters(f => ({ ...f, page: f.page + 1 }))}
              >
                Next <ChevronRight size={16} />
              </button>
            </div>
          )}
        </div>

        {/* Job detail panel */}
        <AnimatePresence>
          {selectedJob && (
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
              className="w-96 glass-card p-6 shrink-0 sticky top-0 self-start overflow-y-auto"
              style={{ maxHeight: '80vh' }}
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h2 className="font-bold text-lg leading-tight mb-1" style={{ color: 'rgb(248 250 252)' }}>
                    {selectedJob.title}
                  </h2>
                  <p className="text-sm" style={{ color: 'rgb(100 116 139)' }}>{selectedJob.location}</p>
                </div>
                <button onClick={() => setSelectedJob(null)} style={{ color: 'rgb(71 85 105)' }}>
                  <X size={18} />
                </button>
              </div>

              <div className="flex flex-wrap gap-2 mb-5">
                {selectedJob.is_remote && <span className="badge badge-green">Remote</span>}
                {selectedJob.is_hybrid && <span className="badge badge-blue">Hybrid</span>}
                {selectedJob.employment_type && <span className="badge badge-purple capitalize">{selectedJob.employment_type}</span>}
                {selectedJob.experience_min > 0 && (
                  <span className="badge badge-yellow">{selectedJob.experience_min}–{selectedJob.experience_max || '∞'} yrs</span>
                )}
              </div>

              {selectedJob.match_score > 0 && (
                <div className="p-4 rounded-xl mb-4" style={{ background: 'rgba(99 102 241 / 0.08)', border: '1px solid rgba(99 102 241 / 0.2)' }}>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs font-semibold uppercase tracking-wider" style={{ color: 'rgb(129 140 248)' }}>Match Score</span>
                    <span className="text-sm font-bold" style={{ color: 'rgb(129 140 248)' }}>{selectedJob.match_score}%</span>
                  </div>
                  <div className="space-y-1.5 text-xs animate-fadeIn" style={{ color: 'rgb(148 163 184)' }}>
                    {selectedJob.match_breakdown && (
                      <>
                        <div className="flex justify-between">
                          <span>Skills Match</span>
                          <span style={{ color: selectedJob.match_breakdown.skill_match > 0 ? 'rgb(34 197 94)' : 'rgb(100 116 139)' }}>
                            {selectedJob.match_breakdown.skill_match > 0 ? '✓ Match' : '—'}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span>Location Match</span>
                          <span style={{ color: selectedJob.match_breakdown.location_match > 0 ? 'rgb(34 197 94)' : 'rgb(100 116 139)' }}>
                            {selectedJob.match_breakdown.location_match > 0 ? '✓ Match' : '—'}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span>Salary Match</span>
                          <span style={{ color: selectedJob.match_breakdown.salary_match > 0 ? 'rgb(34 197 94)' : 'rgb(100 116 139)' }}>
                            {selectedJob.match_breakdown.salary_match > 0 ? '✓ Match' : '—'}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span>Experience Match</span>
                          <span style={{ color: selectedJob.match_breakdown.experience_match > 0 ? 'rgb(34 197 94)' : 'rgb(100 116 139)' }}>
                            {selectedJob.match_breakdown.experience_match > 0 ? '✓ Match' : '—'}
                          </span>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              )}

              {selectedJob.salary_min > 0 && (
                <div className="p-3 rounded-lg mb-4" style={{ background: 'rgba(250 204 21 / 0.08)', border: '1px solid rgba(250 204 21 / 0.2)' }}>
                  <p className="text-xs" style={{ color: 'rgb(100 116 139)' }}>Salary Range</p>
                  <p className="font-bold" style={{ color: 'rgb(250 204 21)' }}>
                    ₹{(selectedJob.salary_min / 100000).toFixed(0)}L – ₹{(selectedJob.salary_max / 100000).toFixed(0)}L / yr
                  </p>
                </div>
              )}

              <div className="mb-4">
                <p className="text-xs font-semibold uppercase tracking-wider mb-2" style={{ color: 'rgb(71 85 105)' }}>Required Skills</p>
                <div className="flex flex-wrap gap-1.5">
                  {selectedJob.skills?.map((s: string) => (
                    <span key={s} className="badge badge-purple text-xs">{s}</span>
                  ))}
                </div>
              </div>

              {selectedJob.sources?.length > 0 && (
                <div className="mb-4">
                  <p className="text-xs font-semibold uppercase tracking-wider mb-2" style={{ color: 'rgb(71 85 105)' }}>Sources</p>
                  {selectedJob.sources.map((src: any) => (
                    <div key={src.provider} className="flex items-center gap-2 text-sm mb-1" style={{ color: 'rgb(148 163 184)' }}>
                      <span className="capitalize">{src.provider}</span>
                    </div>
                  ))}
                </div>
              )}

              <div className="space-y-2">
                <a
                  href={selectedJob.application_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="btn-primary w-full flex items-center justify-center gap-2"
                >
                  Apply Now <ExternalLink size={14} />
                </a>
                <button
                  onClick={() => save(selectedJob)}
                  className="btn-secondary w-full flex items-center justify-center gap-2"
                >
                  <Bookmark size={14} /> Save Job
                </button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </PageLayout>
  )
}
