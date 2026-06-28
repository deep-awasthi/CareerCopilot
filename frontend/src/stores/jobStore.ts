import { create } from 'zustand'

interface NotificationState {
  unreadCount: number
  setUnreadCount: (count: number) => void
  incrementUnread: () => void
  resetUnread: () => void
}

export const useNotificationStore = create<NotificationState>((set) => ({
  unreadCount: 0,
  setUnreadCount: (count) => set({ unreadCount: count }),
  incrementUnread: () => set((s) => ({ unreadCount: s.unreadCount + 1 })),
  resetUnread: () => set({ unreadCount: 0 }),
}))

interface JobFilter {
  query: string
  location: string
  remote: boolean
  experienceMin: number
  experienceMax: number
  salaryMin: number
  provider: string
}

interface JobState {
  filters: JobFilter
  selectedJobId: number | null
  setFilters: (filters: Partial<JobFilter>) => void
  resetFilters: () => void
  setSelectedJob: (id: number | null) => void
}

const defaultFilters: JobFilter = {
  query: '',
  location: '',
  remote: false,
  experienceMin: 0,
  experienceMax: 30,
  salaryMin: 0,
  provider: '',
}

export const useJobStore = create<JobState>((set) => ({
  filters: defaultFilters,
  selectedJobId: null,
  setFilters: (filters) => set((s) => ({ filters: { ...s.filters, ...filters } })),
  resetFilters: () => set({ filters: defaultFilters }),
  setSelectedJob: (id) => set({ selectedJobId: id }),
}))
