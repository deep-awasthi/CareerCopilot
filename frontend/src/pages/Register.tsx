import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useForm } from 'react-hook-form'
import { useRegister } from '@/api/hooks'
import { Eye, EyeOff, Briefcase, ArrowRight, CheckCircle2 } from 'lucide-react'

interface RegisterForm {
  name: string
  email: string
  password: string
  confirmPassword: string
}

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false)
  const { register, handleSubmit, watch, formState: { errors, isSubmitting } } = useForm<RegisterForm>()
  const registerMutation = useRegister()
  const navigate = useNavigate()

  const onSubmit = async (data: RegisterForm) => {
    try {
      await registerMutation.mutateAsync({ email: data.email, password: data.password, name: data.name })
      navigate('/login?registered=true')
    } catch {}
  }

  const features = [
    'Discover jobs from 8+ portals',
    'Deduplicated job listings',
    'Referral opportunity finder',
    'Interview stage tracker',
    'Daily email digest',
  ]

  return (
    <div className="min-h-screen flex items-center justify-center relative overflow-hidden"
      style={{ background: 'rgb(10 10 15)' }}>

      <div className="absolute top-1/3 left-1/3 w-96 h-96 rounded-full blur-3xl pointer-events-none"
        style={{ background: 'rgba(99 102 241 / 0.07)' }} />
      <div className="absolute bottom-1/3 right-1/3 w-96 h-96 rounded-full blur-3xl pointer-events-none"
        style={{ background: 'rgba(236 72 153 / 0.05)' }} />

      <div className="w-full max-w-4xl mx-4 grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
        {/* Left: Benefits */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5 }}
          className="hidden md:block"
        >
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-xl flex items-center justify-center"
              style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
              <Briefcase size={18} color="white" />
            </div>
            <span className="text-xl font-bold" style={{ color: 'rgb(248 250 252)' }}>CareerCopilot</span>
          </div>
          <h2 className="text-4xl font-bold mb-4 leading-tight" style={{ color: 'rgb(248 250 252)' }}>
            Your career,<br />
            <span className="gradient-text">centralized.</span>
          </h2>
          <p className="text-lg mb-8" style={{ color: 'rgb(100 116 139)' }}>
            One platform to discover jobs, track applications, find referrals, and land your next role.
          </p>
          <div className="space-y-3">
            {features.map((f, i) => (
              <motion.div
                key={f}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.1 + 0.3 }}
                className="flex items-center gap-3"
              >
                <CheckCircle2 size={18} style={{ color: 'rgb(129 140 248)' }} />
                <span className="text-sm" style={{ color: 'rgb(148 163 184)' }}>{f}</span>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Right: Form */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <div className="text-center mb-6">
            <h1 className="text-2xl font-bold" style={{ color: 'rgb(248 250 252)' }}>Create your account</h1>
            <p className="text-sm mt-1" style={{ color: 'rgb(71 85 105)' }}>Free forever. No credit card needed.</p>
          </div>

          <div className="glass-card p-8">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Full name</label>
                <input
                  id="register-name"
                  className="input-field"
                  placeholder="Jane Doe"
                  {...register('name', { required: 'Name is required' })}
                />
                {errors.name && <p className="text-xs mt-1" style={{ color: 'rgb(239 68 68)' }}>{errors.name.message}</p>}
              </div>

              <div>
                <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Email address</label>
                <input
                  id="register-email"
                  type="email"
                  className="input-field"
                  placeholder="jane@example.com"
                  {...register('email', { required: 'Email is required', pattern: { value: /\S+@\S+\.\S+/, message: 'Invalid email' } })}
                />
                {errors.email && <p className="text-xs mt-1" style={{ color: 'rgb(239 68 68)' }}>{errors.email.message}</p>}
              </div>

              <div>
                <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Password</label>
                <div className="relative">
                  <input
                    id="register-password"
                    type={showPassword ? 'text' : 'password'}
                    className="input-field pr-10"
                    placeholder="Min. 8 characters"
                    {...register('password', { required: 'Password is required', minLength: { value: 8, message: 'Min 8 characters' } })}
                  />
                  <button type="button" className="absolute right-3 top-1/2 -translate-y-1/2"
                    style={{ color: 'rgb(71 85 105)' }} onClick={() => setShowPassword(v => !v)}>
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
                {errors.password && <p className="text-xs mt-1" style={{ color: 'rgb(239 68 68)' }}>{errors.password.message}</p>}
              </div>

              <div>
                <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>Confirm password</label>
                <input
                  id="register-confirm-password"
                  type="password"
                  className="input-field"
                  placeholder="••••••••"
                  {...register('confirmPassword', {
                    validate: v => v === watch('password') || 'Passwords do not match'
                  })}
                />
                {errors.confirmPassword && <p className="text-xs mt-1" style={{ color: 'rgb(239 68 68)' }}>{errors.confirmPassword.message}</p>}
              </div>

              {registerMutation.isError && (
                <div className="rounded-lg p-3 text-sm"
                  style={{ background: 'rgba(239 68 68 / 0.1)', border: '1px solid rgba(239 68 68 / 0.2)', color: 'rgb(239 68 68)' }}>
                  Registration failed. Email may already be in use.
                </div>
              )}

              <button
                id="register-submit"
                type="submit"
                disabled={isSubmitting || registerMutation.isPending}
                className="btn-primary w-full flex items-center justify-center gap-2"
              >
                {registerMutation.isPending ? (
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                ) : (
                  <>Create account <ArrowRight size={16} /></>
                )}
              </button>
            </form>

            <p className="text-center mt-5 text-sm" style={{ color: 'rgb(71 85 105)' }}>
              Already have an account?{' '}
              <Link to="/login" className="font-medium" style={{ color: 'rgb(129 140 248)' }}>Sign in</Link>
            </p>
          </div>
        </motion.div>
      </div>
    </div>
  )
}
