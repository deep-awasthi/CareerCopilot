import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/axios'

// -------- AUTH --------
export const useLogin = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { email: string; password: string }) =>
      api.post('/auth/login', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['me'] }),
  })
}

export const useRegister = () =>
  useMutation({
    mutationFn: (data: { email: string; password: string; name: string }) =>
      api.post('/auth/register', data).then(r => r.data),
  })

export const useMe = () =>
  useQuery({ queryKey: ['me'], queryFn: () => api.get('/auth/me').then(r => r.data.data) })

// -------- PROFILE --------
export const useProfile = () =>
  useQuery({ queryKey: ['profile'], queryFn: () => api.get('/profile').then(r => r.data.data) })

export const useUpdateProfile = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.put('/profile', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['profile'] }),
  })
}

// -------- RESUME --------
export const useResume = () =>
  useQuery({ queryKey: ['resume'], queryFn: () => api.get('/resume').then(r => r.data.data), retry: false })

export const useSubmitResume = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: { raw_text: string }) => api.post('/resume', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['resume'] }),
  })
}

// -------- JOBS --------
export const useJobs = (params: object) =>
  useQuery({
    queryKey: ['jobs', params],
    queryFn: () => api.get('/jobs', { params }).then(r => r.data),
  })

export const useJob = (id: number | string) =>
  useQuery({
    queryKey: ['job', id],
    queryFn: () => api.get(`/jobs/${id}`).then(r => r.data.data),
    enabled: !!id,
  })

export const useMatchedJobs = (params: object) =>
  useQuery({
    queryKey: ['jobs', 'matched', params],
    queryFn: () => api.get('/jobs/matched', { params }).then(r => r.data),
  })

// -------- APPLICATIONS --------
export const useApplications = (params: object = {}) =>
  useQuery({
    queryKey: ['applications', params],
    queryFn: () => api.get('/applications', { params }).then(r => r.data),
  })

export const useApplicationStats = () =>
  useQuery({
    queryKey: ['applications', 'stats'],
    queryFn: () => api.get('/applications/stats').then(r => r.data.data),
  })

export const useCreateApplication = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.post('/applications', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['applications'] }),
  })
}

export const useUpdateApplication = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: { id: number } & Record<string, any>) =>
      api.put(`/applications/${id}`, data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['applications'] }),
  })
}

export const useDeleteApplication = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/applications/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['applications'] }),
  })
}

// -------- INTERVIEWS --------
export const useInterviews = () =>
  useQuery({ queryKey: ['interviews'], queryFn: () => api.get('/interviews').then(r => r.data.data) })

export const useUpcomingInterviews = () =>
  useQuery({
    queryKey: ['interviews', 'upcoming'],
    queryFn: () => api.get('/interviews/upcoming').then(r => r.data.data),
  })

export const useCreateInterview = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.post('/interviews', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['interviews'] }),
  })
}

export const useAddRound = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: { id: number } & Record<string, any>) =>
      api.post(`/interviews/${id}/rounds`, data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['interviews'] }),
  })
}

// -------- REFERRALS --------
export const useReferrals = (params: object = {}) =>
  useQuery({
    queryKey: ['referrals', params],
    queryFn: () => api.get('/referrals', { params }).then(r => r.data),
  })

export const useCreateReferral = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.post('/referrals', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['referrals'] }),
  })
}

export const useUpdateReferral = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: { id: number } & Record<string, any>) =>
      api.put(`/referrals/${id}`, data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['referrals'] }),
  })
}

// -------- COMPANIES --------
export const useCompanies = (params: object = {}) =>
  useQuery({
    queryKey: ['companies', params],
    queryFn: () => api.get('/companies', { params }).then(r => r.data),
  })

export const useWatchlist = () =>
  useQuery({ queryKey: ['watchlist'], queryFn: () => api.get('/watchlists').then(r => r.data) })

export const useWatchCompany = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post(`/companies/${id}/watch`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['watchlist'] }),
  })
}

export const useUnwatchCompany = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/companies/${id}/watch`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['watchlist'] }),
  })
}

// -------- SEARCH PROFILES --------
export const useSearchProfiles = () =>
  useQuery({ queryKey: ['search-profiles'], queryFn: () => api.get('/search-profiles').then(r => r.data.data) })

export const useCreateSearchProfile = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.post('/search-profiles', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['search-profiles'] }),
  })
}

export const useUpdateSearchProfile = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: { id: number } & Record<string, any>) =>
      api.put(`/search-profiles/${id}`, data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['search-profiles'] }),
  })
}

export const useDeleteSearchProfile = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/search-profiles/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['search-profiles'] }),
  })
}

// -------- KEYWORD ALERTS --------
export const useKeywordAlerts = () =>
  useQuery({ queryKey: ['alerts'], queryFn: () => api.get('/alerts').then(r => r.data.data) })

export const useCreateAlert = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: object) => api.post('/alerts', data).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['alerts'] }),
  })
}

export const useDeleteAlert = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete(`/alerts/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['alerts'] }),
  })
}

// -------- NOTIFICATIONS --------
export const useNotifications = (params: object = {}) =>
  useQuery({
    queryKey: ['notifications', params],
    queryFn: () => api.get('/notifications', { params }).then(r => r.data),
  })

export const useUnreadCount = () =>
  useQuery({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => api.get('/notifications/unread-count').then(r => r.data.data),
    refetchInterval: 30000,
  })

export const useMarkRead = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.put(`/notifications/${id}/read`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notifications'] }),
  })
}

export const useMarkAllRead = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.put('/notifications/read-all').then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notifications'] }),
  })
}

// -------- ANALYTICS --------
export const useDashboardStats = () =>
  useQuery({
    queryKey: ['analytics', 'dashboard'],
    queryFn: () => api.get('/analytics/dashboard').then(r => r.data.data),
  })

export const useAnalytics = () =>
  useQuery({
    queryKey: ['analytics'],
    queryFn: () => api.get('/analytics').then(r => r.data.data),
  })

// -------- BOOKMARKS --------
export const useBookmarks = (type?: string) =>
  useQuery({
    queryKey: ['bookmarks', type],
    queryFn: () => api.get('/bookmarks', { params: type ? { type } : {} }).then(r => r.data),
  })

export const useBookmark = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ type, id }: { type: string; id: number }) =>
      api.post(`/bookmarks/${type}/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['bookmarks'] }),
  })
}

export const useUnbookmark = () => {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ type, id }: { type: string; id: number }) =>
      api.delete(`/bookmarks/${type}/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['bookmarks'] }),
  })
}

// -------- SEARCH --------
export const useSearch = (params: object) =>
  useQuery({
    queryKey: ['search', params],
    queryFn: () => api.get('/search', { params }).then(r => r.data.data),
    enabled: Object.values(params).some(v => v !== '' && v !== undefined),
  })
