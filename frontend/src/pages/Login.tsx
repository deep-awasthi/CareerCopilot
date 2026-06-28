import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useForm } from 'react-hook-form'
import { useLogin } from '@/api/hooks'
import { useAuthStore } from '@/stores/authStore'
import { Eye, EyeOff, Briefcase, ArrowRight } from 'lucide-react'

interface LoginForm {
  email: string
  password: string
}

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false)
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginForm>()
  const login = useLogin()
  const setAuth = useAuthStore(s => s.setAuth)
  const navigate = useNavigate()

  const onSubmit = async (data: LoginForm) => {
    try {
      const res = await login.mutateAsync(data)
      setAuth(res.data.user, res.data.access_token, res.data.refresh_token)
      navigate('/')
    } catch {
      // error handled below
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center relative overflow-hidden"
      style={{ background: 'rgb(10 10 15)' }}>

      {/* Background orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 rounded-full blur-3xl pointer-events-none"
        style={{ background: 'rgba(99 102 241 / 0.08)' }} />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 rounded-full blur-3xl pointer-events-none"
        style={{ background: 'rgba(139 92 246 / 0.06)' }} />

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md mx-4"
      >
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-xl mb-4"
            style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}>
            <Briefcase size={20} color="white" />
          </div>
          <h1 className="text-3xl font-bold mb-2" style={{ color: 'rgb(248 250 252)' }}>
            Welcome back
          </h1>
          <p style={{ color: 'rgb(71 85 105)' }}>
            Track jobs. Find referrals. Grow your career.
          </p>
        </div>

        {/* Card */}
        <div className="glass-card p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                Email address
              </label>
              <input
                id="login-email"
                type="email"
                className="input-field"
                placeholder="you@example.com"
                {...register('email', { required: 'Email is required' })}
              />
              {errors.email && <p className="text-xs mt-1.5" style={{ color: 'rgb(239 68 68)' }}>{errors.email.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'rgb(148 163 184)' }}>
                Password
              </label>
              <div className="relative">
                <input
                  id="login-password"
                  type={showPassword ? 'text' : 'password'}
                  className="input-field pr-10"
                  placeholder="••••••••"
                  {...register('password', { required: 'Password is required' })}
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 -translate-y-1/2"
                  style={{ color: 'rgb(71 85 105)' }}
                  onClick={() => setShowPassword(v => !v)}
                >
                  {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
              {errors.password && <p className="text-xs mt-1.5" style={{ color: 'rgb(239 68 68)' }}>{errors.password.message}</p>}
            </div>

            {login.isError && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="rounded-lg p-3 text-sm"
                style={{ background: 'rgba(239 68 68 / 0.1)', border: '1px solid rgba(239 68 68 / 0.2)', color: 'rgb(239 68 68)' }}
              >
                Invalid email or password. Please try again.
              </motion.div>
            )}

            <button
              id="login-submit"
              type="submit"
              disabled={isSubmitting || login.isPending}
              className="btn-primary w-full flex items-center justify-center gap-2"
            >
              {isSubmitting || login.isPending ? (
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              ) : (
                <>Sign in <ArrowRight size={16} /></>
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm" style={{ color: 'rgb(71 85 105)' }}>
              Don't have an account?{' '}
              <Link to="/register" className="font-medium" style={{ color: 'rgb(129 140 248)' }}>
                Create account
              </Link>
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  )
}
