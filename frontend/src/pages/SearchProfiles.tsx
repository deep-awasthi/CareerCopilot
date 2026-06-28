import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useSearchProfiles, useCreateSearchProfile, useDeleteSearchProfile, useUpdateSearchProfile } from '@/api/hooks'
import { useForm } from 'react-hook-form'
import { Search, Plus, Trash2, Edit2, MapPin, Briefcase, DollarSign, Wifi, ToggleLeft, ToggleRight } from 'lucide-react'

export default function SearchProfilesPage() {
  const { data: profiles, isLoading } = useSearchProfiles()
  const createProfile = useCreateSearchProfile()
  const updateProfile = useUpdateSearchProfile()
  const deleteProfile = useDeleteSearchProfile()
  const [showCreate, setShowCreate] = useState(false)
  const { register, handleSubmit, reset } = useForm()

  const onSubmit = (data: any) => {
    const processed = {
      ...data,
      keywords: data.keywords ? data.keywords.split(',').map((k: string) => k.trim()) : [],
      locations: data.locations ? data.locations.split(',').map((l: string) => l.trim()) : [],
      experience_min: parseFloat(data.experience_min) || 0,
      experience_max: parseFloat(data.experience_max) || 0,
      salary_min: parseFloat(data.salary_min) || 0,
    }
    createProfile.mutate(processed, { onSuccess: () => { reset(); setShowCreate(false) } })
  }

  const toggleActive = (profile: any) => {
    updateProfile.mutate({ id: profile.id, is_active: !profile.is_active })
  }

  return (
    <PageLayout title="Search Profiles">
      <div className="flex items-center justify-between mb-6">
        <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
          Save search criteria — the scheduler runs them every 6 hours to fetch new jobs
        </p>
        <button id="create-profile-btn" className="btn-primary flex items-center gap-2" onClick={() => setShowCreate(v => !v)}>
          <Plus size={16} /> New Profile
        </button>
      </div>

      {/* Create form */}
      <AnimatePresence>
        {showCreate && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="overflow-hidden mb-6"
          >
            <div className="glass-card p-6">
              <h3 className="font-semibold mb-4" style={{ color: 'rgb(248 250 252)' }}>Create Search Profile</h3>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Profile Name *</label>
                    <input className="input-field" placeholder="e.g. Senior Backend Engineer" {...register('name', { required: true })} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Keywords (comma-separated)</label>
                    <input className="input-field" placeholder="golang, microservices, kubernetes" {...register('keywords')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Locations (comma-separated)</label>
                    <input className="input-field" placeholder="Bangalore, Remote, Mumbai" {...register('locations')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Job Type</label>
                    <select className="input-field" {...register('job_type')}>
                      <option value="">Any</option>
                      <option value="full-time">Full-time</option>
                      <option value="contract">Contract</option>
                      <option value="internship">Internship</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Min Experience (yrs)</label>
                    <input type="number" className="input-field" placeholder="0" {...register('experience_min')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Max Experience (yrs)</label>
                    <input type="number" className="input-field" placeholder="10" {...register('experience_max')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Min Salary</label>
                    <input type="number" className="input-field" placeholder="1000000" {...register('salary_min')} />
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" {...register('is_remote')} className="w-4 h-4" style={{ accentColor: 'rgb(99 102 241)' }} />
                    <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>Remote only</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" {...register('is_hybrid')} className="w-4 h-4" style={{ accentColor: 'rgb(99 102 241)' }} />
                    <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>Include hybrid</span>
                  </label>
                </div>
                <div className="flex gap-3">
                  <button type="submit" className="btn-primary" disabled={createProfile.isPending}>
                    {createProfile.isPending ? 'Creating...' : 'Create Profile'}
                  </button>
                  <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                </div>
              </form>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Profile cards */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="glass-card p-5 animate-pulse h-40" />
          ))}
        </div>
      ) : !profiles || profiles.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Search size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No search profiles</p>
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Create a profile to automatically scrape matching jobs every 6 hours.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {profiles.map((profile: any, i: number) => (
            <motion.div
              key={profile.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              className="glass-card p-5"
              style={{ opacity: profile.is_active ? 1 : 0.6 }}
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>{profile.name}</h3>
                  {profile.last_run_at && (
                    <p className="text-xs mt-0.5" style={{ color: 'rgb(71 85 105)' }}>
                      Last run: {new Date(profile.last_run_at).toLocaleDateString()}
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => toggleActive(profile)}
                    style={{ color: profile.is_active ? 'rgb(34 197 94)' : 'rgb(71 85 105)' }}
                    title={profile.is_active ? 'Deactivate' : 'Activate'}
                  >
                    {profile.is_active ? <ToggleRight size={22} /> : <ToggleLeft size={22} />}
                  </button>
                  <button
                    onClick={() => deleteProfile.mutate(profile.id)}
                    style={{ color: 'rgb(71 85 105)' }}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>

              <div className="space-y-2">
                {profile.keywords?.length > 0 && (
                  <div className="flex items-center gap-2 flex-wrap">
                    <Search size={12} style={{ color: 'rgb(71 85 105)', flexShrink: 0 }} />
                    {profile.keywords.slice(0, 4).map((k: string) => (
                      <span key={k} className="badge badge-purple text-xs">{k}</span>
                    ))}
                    {profile.keywords.length > 4 && <span className="text-xs" style={{ color: 'rgb(71 85 105)' }}>+{profile.keywords.length - 4}</span>}
                  </div>
                )}
                {profile.locations?.length > 0 && (
                  <div className="flex items-center gap-2">
                    <MapPin size={12} style={{ color: 'rgb(71 85 105)' }} />
                    <span className="text-xs" style={{ color: 'rgb(148 163 184)' }}>{profile.locations.join(', ')}</span>
                  </div>
                )}
                {(profile.experience_min > 0 || profile.experience_max > 0) && (
                  <div className="flex items-center gap-2">
                    <Briefcase size={12} style={{ color: 'rgb(71 85 105)' }} />
                    <span className="text-xs" style={{ color: 'rgb(148 163 184)' }}>
                      {profile.experience_min}–{profile.experience_max || '∞'} years exp
                    </span>
                  </div>
                )}
                {profile.salary_min > 0 && (
                  <div className="flex items-center gap-2">
                    <DollarSign size={12} style={{ color: 'rgb(71 85 105)' }} />
                    <span className="text-xs" style={{ color: 'rgb(148 163 184)' }}>
                      Min ₹{(profile.salary_min / 100000).toFixed(0)}L
                    </span>
                  </div>
                )}
                {profile.is_remote && (
                  <div className="flex items-center gap-2">
                    <Wifi size={12} style={{ color: 'rgb(34 197 94)' }} />
                    <span className="text-xs" style={{ color: 'rgb(34 197 94)' }}>Remote only</span>
                  </div>
                )}
              </div>
            </motion.div>
          ))}
        </div>
      )}
    </PageLayout>
  )
}
