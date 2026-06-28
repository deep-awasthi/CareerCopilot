import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { PageLayout } from '@/components/layout/Layout'
import { useInterviews, useCreateInterview, useAddRound } from '@/api/hooks'
import { useForm } from 'react-hook-form'
import {
  Calendar, Plus, ChevronRight, Clock,
  CheckCircle2, XCircle, Circle, Award
} from 'lucide-react'

const STAGES = [
  'recruiter_call', 'online_assessment', 'technical_round_1',
  'technical_round_2', 'system_design', 'manager_round', 'hr_round', 'offer'
]

const STAGE_LABELS: Record<string, string> = {
  applied: 'Applied', recruiter_call: 'Recruiter Call', online_assessment: 'OA',
  technical_round_1: 'Tech Round 1', technical_round_2: 'Tech Round 2',
  system_design: 'System Design', manager_round: 'Manager', hr_round: 'HR', offer: 'Offer'
}

const RESULT_COLOR: Record<string, string> = {
  passed: 'rgb(34 197 94)', failed: 'rgb(239 68 68)', pending: 'rgb(250 204 21)', cancelled: 'rgb(100 116 139)'
}

const ResultIcon = ({ result }: { result: string }) => {
  if (result === 'passed') return <CheckCircle2 size={16} style={{ color: RESULT_COLOR.passed }} />
  if (result === 'failed') return <XCircle size={16} style={{ color: RESULT_COLOR.failed }} />
  return <Circle size={16} style={{ color: RESULT_COLOR.pending }} />
}

export default function InterviewsPage() {
  const { data: interviews, isLoading } = useInterviews()
  const createInterview = useCreateInterview()
  const addRound = useAddRound()
  const [showCreate, setShowCreate] = useState(false)
  const [selectedIv, setSelectedIv] = useState<any>(null)
  const [showAddRound, setShowAddRound] = useState(false)

  const { register, handleSubmit, reset } = useForm()
  const { register: regRound, handleSubmit: submitRound, reset: resetRound } = useForm()

  const onCreateInterview = (data: any) => {
    createInterview.mutate(data, {
      onSuccess: () => { reset(); setShowCreate(false) }
    })
  }

  const onAddRound = (data: any) => {
    if (!selectedIv) return
    addRound.mutate({ id: selectedIv.id, ...data }, {
      onSuccess: () => { resetRound(); setShowAddRound(false) }
    })
  }

  return (
    <PageLayout title="Interview Tracker">
      <div className="flex items-center justify-between mb-6">
        <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
          Track every stage of your interview process
        </p>
        <button
          id="add-interview-btn"
          className="btn-primary flex items-center gap-2"
          onClick={() => setShowCreate(true)}
        >
          <Plus size={16} /> Add Interview
        </button>
      </div>

      {/* Create interview modal */}
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
              initial={{ scale: 0.95 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0.95 }}
              className="glass-card p-8 w-full max-w-md mx-4"
            >
              <h3 className="text-lg font-bold mb-5" style={{ color: 'rgb(248 250 252)' }}>Start Interview Tracking</h3>
              <form onSubmit={handleSubmit(onCreateInterview)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Application ID</label>
                  <input className="input-field" placeholder="Application ID" {...register('application_id', { required: true, valueAsNumber: true })} />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Notes</label>
                  <textarea className="input-field" rows={3} placeholder="Any initial notes..." {...register('notes')} />
                </div>
                <div className="flex gap-3">
                  <button type="submit" className="btn-primary flex-1">Create</button>
                  <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                </div>
              </form>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {isLoading ? (
        <div className="glass-card p-8 text-center">
          <div className="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto" />
        </div>
      ) : !interviews || interviews.length === 0 ? (
        <div className="glass-card p-16 text-center">
          <Calendar size={48} className="mx-auto mb-4" style={{ color: 'rgba(45 45 65 / 0.8)' }} />
          <p className="text-lg font-semibold mb-2" style={{ color: 'rgb(100 116 139)' }}>No interviews tracked yet</p>
          <p className="text-sm mb-4" style={{ color: 'rgb(71 85 105)' }}>Start tracking your interview progress by adding an interview.</p>
          <button className="btn-primary" onClick={() => setShowCreate(true)}>
            <Plus size={16} className="inline mr-2" /> Add Interview
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Interview list */}
          <div className="lg:col-span-1 space-y-3">
            {interviews.map((iv: any) => (
              <motion.div
                key={iv.id}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                onClick={() => setSelectedIv(iv)}
                className="glass-card p-4 cursor-pointer"
                style={{
                  borderColor: selectedIv?.id === iv.id ? 'rgba(99 102 241 / 0.5)' : undefined,
                }}
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-sm" style={{ color: 'rgb(248 250 252)' }}>
                      Application #{iv.application_id}
                    </p>
                    <p className="text-xs mt-1 capitalize" style={{ color: 'rgb(100 116 139)' }}>
                      Current: {STAGE_LABELS[iv.current_stage] || iv.current_stage}
                    </p>
                  </div>
                  <div className="flex items-center gap-1">
                    <span className="badge badge-purple text-xs">{iv.rounds?.length ?? 0} rounds</span>
                    <ChevronRight size={14} style={{ color: 'rgb(71 85 105)' }} />
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Interview detail / timeline */}
          {selectedIv && (
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              className="lg:col-span-2 glass-card p-6"
            >
              <div className="flex items-center justify-between mb-6">
                <h3 className="font-bold text-lg" style={{ color: 'rgb(248 250 252)' }}>
                  Application #{selectedIv.application_id}
                </h3>
                <button
                  id="add-round-btn"
                  className="btn-primary flex items-center gap-2 text-sm"
                  onClick={() => setShowAddRound(v => !v)}
                >
                  <Plus size={14} /> Add Round
                </button>
              </div>

              {/* Add round form */}
              <AnimatePresence>
                {showAddRound && (
                  <motion.form
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: 'auto', opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    onSubmit={submitRound(onAddRound)}
                    className="overflow-hidden mb-6"
                  >
                    <div className="p-4 rounded-xl mb-4" style={{ background: 'rgba(28 28 42 / 0.7)', border: '1px solid rgba(45 45 65 / 0.5)' }}>
                      <div className="grid grid-cols-2 gap-3">
                        <div>
                          <label className="block text-xs mb-1" style={{ color: 'rgb(71 85 105)' }}>Stage</label>
                          <select className="input-field" {...regRound('stage', { required: true })}>
                            {STAGES.map(s => <option key={s} value={s}>{STAGE_LABELS[s]}</option>)}
                          </select>
                        </div>
                        <div>
                          <label className="block text-xs mb-1" style={{ color: 'rgb(71 85 105)' }}>Scheduled At</label>
                          <input type="datetime-local" className="input-field" {...regRound('scheduled_at')} />
                        </div>
                        <div>
                          <label className="block text-xs mb-1" style={{ color: 'rgb(71 85 105)' }}>Interviewer</label>
                          <input className="input-field" placeholder="Name" {...regRound('interviewer_name')} />
                        </div>
                        <div>
                          <label className="block text-xs mb-1" style={{ color: 'rgb(71 85 105)' }}>Result</label>
                          <select className="input-field" {...regRound('result')}>
                            <option value="pending">Pending</option>
                            <option value="passed">Passed</option>
                            <option value="failed">Failed</option>
                            <option value="cancelled">Cancelled</option>
                          </select>
                        </div>
                      </div>
                      <div className="mt-3">
                        <label className="block text-xs mb-1" style={{ color: 'rgb(71 85 105)' }}>Feedback / Notes</label>
                        <textarea className="input-field" rows={2} {...regRound('feedback')} />
                      </div>
                      <div className="flex gap-2 mt-3">
                        <button type="submit" className="btn-primary text-sm">Add Round</button>
                        <button type="button" className="btn-secondary text-sm" onClick={() => setShowAddRound(false)}>Cancel</button>
                      </div>
                    </div>
                  </motion.form>
                )}
              </AnimatePresence>

              {/* Timeline */}
              <div className="relative">
                {(!selectedIv.rounds || selectedIv.rounds.length === 0) ? (
                  <p className="text-sm text-center py-8" style={{ color: 'rgb(71 85 105)' }}>
                    No rounds added yet. Click "Add Round" to log an interview round.
                  </p>
                ) : (
                  <div className="space-y-4">
                    {selectedIv.rounds.map((round: any, idx: number) => (
                      <motion.div
                        key={round.id}
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: idx * 0.05 }}
                        className="flex gap-4"
                      >
                        <div className="flex flex-col items-center">
                          <ResultIcon result={round.result} />
                          {idx < selectedIv.rounds.length - 1 && (
                            <div className="w-0.5 flex-1 mt-2" style={{ background: 'rgba(45 45 65 / 0.6)' }} />
                          )}
                        </div>
                        <div className="flex-1 pb-4">
                          <div className="flex items-start justify-between">
                            <div>
                              <p className="font-semibold text-sm" style={{ color: 'rgb(248 250 252)' }}>
                                {STAGE_LABELS[round.stage] || round.stage}
                              </p>
                              {round.interviewer_name && (
                                <p className="text-xs mt-0.5" style={{ color: 'rgb(100 116 139)' }}>
                                  with {round.interviewer_name}
                                  {round.interviewer_role && ` (${round.interviewer_role})`}
                                </p>
                              )}
                            </div>
                            <div className="text-right">
                              <span className="badge text-xs capitalize" style={{
                                background: `${RESULT_COLOR[round.result]}20`,
                                color: RESULT_COLOR[round.result],
                                border: `1px solid ${RESULT_COLOR[round.result]}40`,
                              }}>
                                {round.result}
                              </span>
                              {round.scheduled_at && (
                                <p className="text-xs mt-1 flex items-center gap-1 justify-end" style={{ color: 'rgb(71 85 105)' }}>
                                  <Clock size={11} />
                                  {new Date(round.scheduled_at).toLocaleDateString()}
                                </p>
                              )}
                            </div>
                          </div>
                          {round.feedback && (
                            <p className="text-xs mt-2 p-2 rounded-lg" style={{ background: 'rgba(28 28 42 / 0.5)', color: 'rgb(148 163 184)' }}>
                              {round.feedback}
                            </p>
                          )}
                        </div>
                      </motion.div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </div>
      )}
    </PageLayout>
  )
}
