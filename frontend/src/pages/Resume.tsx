import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useResume, useSubmitResume } from '@/api/hooks'
import { useForm } from 'react-hook-form'
import { FileText, Upload, CheckCircle2, ChevronRight, Code2, GraduationCap, Building2, Award } from 'lucide-react'

export default function ResumePage() {
  const { data: resume, isLoading } = useResume()
  const submitResume = useSubmitResume()
  const [showForm, setShowForm] = useState(false)
  const { register, handleSubmit, formState: { errors } } = useForm<{ raw_text: string }>()

  const onSubmit = (data: { raw_text: string }) => {
    submitResume.mutate(data, { onSuccess: () => setShowForm(false) })
  }

  return (
    <PageLayout title="Resume Parser">
      <div className="max-w-4xl">
        <div className="flex items-center justify-between mb-6">
          <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
            Paste your resume as plain text. Skills are extracted using deterministic keyword matching — no AI.
          </p>
          <button
            id="update-resume-btn"
            className="btn-primary flex items-center gap-2"
            onClick={() => setShowForm(v => !v)}
          >
            <Upload size={16} /> {resume ? 'Update Resume' : 'Upload Resume'}
          </button>
        </div>

        {/* Resume input form */}
        <AnimatePresence>
          {(showForm || !resume) && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: 'auto', opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden mb-6"
            >
              <div className="glass-card p-6">
                <h3 className="font-semibold mb-4" style={{ color: 'rgb(248 250 252)' }}>
                  Paste Your Resume
                </h3>
                <form onSubmit={handleSubmit(onSubmit)}>
                  <textarea
                    id="resume-text"
                    className="input-field mb-4 font-mono text-sm"
                    rows={20}
                    placeholder="Paste your full resume here as plain text...&#10;&#10;Include your work experience, skills, education, projects, and certifications.&#10;&#10;Example:&#10;Skills: Python, React, TypeScript, AWS, Docker, Kubernetes&#10;&#10;Experience:&#10;Senior Software Engineer at Stripe (2022–Present)&#10;  - Built payment reconciliation system processing $2B/month..."
                    {...register('raw_text', {
                      required: 'Resume text is required',
                      minLength: { value: 100, message: 'Resume must be at least 100 characters' }
                    })}
                  />
                  {errors.raw_text && (
                    <p className="text-sm mb-3" style={{ color: 'rgb(239 68 68)' }}>{errors.raw_text.message}</p>
                  )}
                  <div className="flex gap-3">
                    <button
                      type="submit"
                      className="btn-primary flex items-center gap-2"
                      disabled={submitResume.isPending}
                    >
                      {submitResume.isPending ? (
                        <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      ) : <CheckCircle2 size={16} />}
                      {submitResume.isPending ? 'Parsing...' : 'Parse & Save'}
                    </button>
                    {resume && (
                      <button type="button" className="btn-secondary" onClick={() => setShowForm(false)}>
                        Cancel
                      </button>
                    )}
                  </div>
                </form>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Parsed resume display */}
        {isLoading ? (
          <div className="glass-card p-8 text-center">
            <div className="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto" />
          </div>
        ) : resume && !showForm ? (
          <div className="space-y-5">
            {/* Skills */}
            <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} className="glass-card p-6">
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 rounded-lg flex items-center justify-center"
                  style={{ background: 'rgba(99 102 241 / 0.15)' }}>
                  <Code2 size={16} style={{ color: 'rgb(129 140 248)' }} />
                </div>
                <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>
                  Extracted Skills
                  <span className="ml-2 badge badge-purple text-xs">{resume.parsed_skills?.length ?? 0} found</span>
                </h3>
              </div>
              <div className="flex flex-wrap gap-2">
                {resume.parsed_skills?.length > 0 ? (
                  resume.parsed_skills.map((skill: string) => (
                    <span key={skill} className="badge badge-purple text-sm">{skill}</span>
                  ))
                ) : (
                  <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>No skills detected. Make sure to include a skills section.</p>
                )}
              </div>
            </motion.div>

            {/* Companies */}
            {resume.companies?.length > 0 && (
              <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }} className="glass-card p-6">
                <div className="flex items-center gap-2 mb-4">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: 'rgba(250 204 21 / 0.12)' }}>
                    <Building2 size={16} style={{ color: 'rgb(250 204 21)' }} />
                  </div>
                  <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>
                    Companies Detected <span className="ml-2 badge badge-yellow text-xs">{resume.companies.length}</span>
                  </h3>
                </div>
                <div className="flex flex-wrap gap-2">
                  {resume.companies.map((co: string) => (
                    <span key={co} className="badge badge-yellow">{co}</span>
                  ))}
                </div>
              </motion.div>
            )}

            {/* Education */}
            {resume.education?.length > 0 && (
              <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.15 }} className="glass-card p-6">
                <div className="flex items-center gap-2 mb-4">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: 'rgba(34 197 94 / 0.12)' }}>
                    <GraduationCap size={16} style={{ color: 'rgb(34 197 94)' }} />
                  </div>
                  <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>Education</h3>
                </div>
                <div className="space-y-2">
                  {resume.education.map((edu: any, i: number) => (
                    <div key={i} className="p-3 rounded-lg" style={{ background: 'rgba(28 28 42 / 0.6)' }}>
                      <p className="font-medium text-sm" style={{ color: 'rgb(248 250 252)' }}>{edu.institution}</p>
                      <p className="text-xs mt-0.5" style={{ color: 'rgb(100 116 139)' }}>{edu.degree} {edu.year && `· ${edu.year}`}</p>
                    </div>
                  ))}
                </div>
              </motion.div>
            )}

            {/* Certifications */}
            {resume.certifications?.length > 0 && (
              <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }} className="glass-card p-6">
                <div className="flex items-center gap-2 mb-4">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: 'rgba(251 113 133 / 0.12)' }}>
                    <Award size={16} style={{ color: 'rgb(251 113 133)' }} />
                  </div>
                  <h3 className="font-semibold" style={{ color: 'rgb(248 250 252)' }}>Certifications</h3>
                </div>
                <div className="flex flex-wrap gap-2">
                  {resume.certifications.map((cert: string) => (
                    <span key={cert} className="badge badge-red">{cert}</span>
                  ))}
                </div>
              </motion.div>
            )}

            <p className="text-xs text-center" style={{ color: 'rgb(45 45 65)' }}>
              Parsed using deterministic keyword matching · No AI involved
            </p>
          </div>
        ) : !resume && !showForm ? (
          <div className="glass-card p-16 text-center">
            <FileText size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
            <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No resume uploaded yet</p>
            <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
              Upload your resume to automatically extract skills, companies, and education.
            </p>
          </div>
        ) : null}
      </div>
    </PageLayout>
  )
}
