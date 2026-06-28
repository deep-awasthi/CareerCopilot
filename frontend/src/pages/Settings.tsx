import { useState } from 'react'
import { motion } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useProfile, useUpdateProfile } from '@/api/hooks'
import { useAuthStore } from '@/stores/authStore'
import { useForm } from 'react-hook-form'
import { User, Phone, MapPin, Link2, Code2, Globe, Save, Briefcase, Calendar } from 'lucide-react'

export default function SettingsPage() {
  const { data: profile, isLoading } = useProfile()
  const updateProfile = useUpdateProfile()
  const { user } = useAuthStore()
  const [activeTab, setActiveTab] = useState<'profile' | 'preferences' | 'notifications'>('profile')
  const { register, handleSubmit, formState: { isDirty } } = useForm({
    values: profile ?? {},
  })

  const onSubmit = (data: any) => {
    updateProfile.mutate(data)
  }

  const tabs = [
    { id: 'profile', label: 'Profile' },
    { id: 'preferences', label: 'Preferences' },
    { id: 'notifications', label: 'Notifications' },
  ] as const

  return (
    <PageLayout title="Settings">
      {/* Tabs */}
      <div className="flex gap-1 mb-8 p-1 rounded-xl w-fit" style={{ background: 'rgba(20 20 30 / 0.8)', border: '1px solid rgba(45 45 65 / 0.5)' }}>
        {tabs.map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className="px-5 py-2 rounded-lg text-sm font-medium transition-all"
            style={{
              background: activeTab === tab.id ? 'rgba(99 102 241 / 0.2)' : 'transparent',
              color: activeTab === tab.id ? 'rgb(129 140 248)' : 'rgb(71 85 105)',
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'profile' && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="max-w-2xl"
        >
          {/* Account info */}
          <div className="glass-card p-6 mb-6">
            <h3 className="font-semibold mb-4" style={{ color: 'rgb(248 250 252)' }}>Account</h3>
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full flex items-center justify-center text-2xl font-bold"
                style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)', color: 'white' }}>
                {user?.email?.[0]?.toUpperCase()}
              </div>
              <div>
                <p className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>{user?.email}</p>
                <p className="text-sm mt-0.5" style={{ color: 'rgb(71 85 105)' }}>
                  {user?.is_email_verified ? '✅ Email verified' : '⚠️ Email not verified'}
                </p>
              </div>
            </div>
          </div>

          {/* Profile form */}
          {isLoading ? (
            <div className="glass-card p-8 text-center">
              <div className="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto" />
            </div>
          ) : (
            <div className="glass-card p-6">
              <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Profile Information</h3>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      <User size={12} className="inline mr-1" /> Full Name
                    </label>
                    <input className="input-field" placeholder="Jane Doe" {...register('name')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      <Phone size={12} className="inline mr-1" /> Phone
                    </label>
                    <input className="input-field" placeholder="+91 9999999999" {...register('phone')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      <Briefcase size={12} className="inline mr-1" /> Current Role
                    </label>
                    <input className="input-field" placeholder="Senior Engineer" {...register('current_role')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      <Calendar size={12} className="inline mr-1" /> Experience (years)
                    </label>
                    <input type="number" className="input-field" placeholder="5" {...register('experience_years', { valueAsNumber: true })} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      <MapPin size={12} className="inline mr-1" /> Location
                    </label>
                    <input className="input-field" placeholder="Bangalore, India" {...register('location')} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                      Preferred Work Type
                    </label>
                    <select className="input-field" {...register('preferred_work_type')}>
                      <option value="any">Any</option>
                      <option value="remote">Remote</option>
                      <option value="hybrid">Hybrid</option>
                      <option value="onsite">Onsite</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                    <Link2 size={12} className="inline mr-1" /> LinkedIn URL
                  </label>
                  <input className="input-field" placeholder="https://linkedin.com/in/yourname" {...register('linkedin_url')} />
                </div>
                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                    <Code2 size={12} className="inline mr-1" /> GitHub URL
                  </label>
                  <input className="input-field" placeholder="https://github.com/yourname" {...register('github_url')} />
                </div>
                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                    <Globe size={12} className="inline mr-1" /> Portfolio URL
                  </label>
                  <input className="input-field" placeholder="https://yoursite.dev" {...register('portfolio_url')} />
                </div>
                <div>
                  <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Bio</label>
                  <textarea className="input-field" rows={3} placeholder="Brief professional bio..." {...register('bio')} />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Current CTC (₹)</label>
                    <input type="number" className="input-field" placeholder="1200000" {...register('current_ctc', { valueAsNumber: true })} />
                  </div>
                  <div>
                    <label className="block text-xs font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Expected CTC (₹)</label>
                    <input type="number" className="input-field" placeholder="1800000" {...register('expected_ctc', { valueAsNumber: true })} />
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" {...register('is_open_to_work')} className="w-4 h-4" style={{ accentColor: 'rgb(99 102 241)' }} />
                    <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>Open to work</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" {...register('is_open_to_relocate')} className="w-4 h-4" style={{ accentColor: 'rgb(99 102 241)' }} />
                    <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>Open to relocate</span>
                  </label>
                </div>

                <button
                  type="submit"
                  className="btn-primary flex items-center gap-2"
                  disabled={updateProfile.isPending}
                >
                  {updateProfile.isPending ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : <Save size={16} />}
                  Save Profile
                </button>
                {updateProfile.isSuccess && (
                  <p className="text-sm" style={{ color: 'rgb(34 197 94)' }}>✅ Profile saved successfully</p>
                )}
              </form>
            </div>
          )}
        </motion.div>
      )}

      {activeTab === 'preferences' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="glass-card p-6 max-w-2xl">
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Job Preferences</h3>
          <div className="space-y-4">
            <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
              Manage your search profiles and keyword alerts from their dedicated pages.
            </p>
            <div className="grid grid-cols-2 gap-3">
              <a href="/search-profiles" className="btn-secondary text-center text-sm">→ Search Profiles</a>
              <a href="/alerts" className="btn-secondary text-center text-sm">→ Keyword Alerts</a>
            </div>
          </div>
        </motion.div>
      )}

      {activeTab === 'notifications' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="glass-card p-6 max-w-2xl">
          <h3 className="font-semibold mb-5" style={{ color: 'rgb(248 250 252)' }}>Notification Preferences</h3>
          <div className="space-y-4">
            {[
              { label: 'Daily job digest email', desc: 'Sent every morning at 7 AM', defaultChecked: true },
              { label: 'Keyword match alerts', desc: 'When new jobs match your keywords', defaultChecked: true },
              { label: 'Company opening alerts', desc: 'When watched companies post new jobs', defaultChecked: true },
              { label: 'Interview reminders', desc: '24 hours before scheduled interviews', defaultChecked: true },
            ].map(item => (
              <div key={item.label} className="flex items-center justify-between p-4 rounded-xl"
                style={{ background: 'rgba(28 28 42 / 0.5)', border: '1px solid rgba(45 45 65 / 0.5)' }}>
                <div>
                  <p className="text-sm font-medium" style={{ color: 'rgb(248 250 252)' }}>{item.label}</p>
                  <p className="text-xs mt-0.5" style={{ color: 'rgb(71 85 105)' }}>{item.desc}</p>
                </div>
                <input type="checkbox" defaultChecked={item.defaultChecked} className="w-5 h-5" style={{ accentColor: 'rgb(99 102 241)' }} />
              </div>
            ))}
          </div>
        </motion.div>
      )}
    </PageLayout>
  )
}
