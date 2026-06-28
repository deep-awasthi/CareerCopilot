# CareerCopilot 🚀

> **Track jobs. Find referrals. Grow your career.**

A production-grade, full-stack career management platform for software engineers. Discover jobs from 8+ portals, monitor company career pages, find referral opportunities, track applications, and manage interview progress — all from one dashboard.

---

## 🏗 Architecture

```
CareerCopilot/
├── backend/                 # Go + Gin + GORM + PostgreSQL
│   ├── cmd/server/main.go   # Server entrypoint
│   ├── internal/
│   │   ├── auth/            # JWT auth (register, login, refresh)
│   │   ├── user/            # Profile management
│   │   ├── resume/          # Deterministic skill parsing (no AI)
│   │   ├── job/             # Job entity, deduplication, matching
│   │   ├── provider/        # 8 job provider adapters
│   │   ├── application/     # Application tracking
│   │   ├── interview/       # Interview rounds & scheduling
│   │   ├── referral/        # Referral lifecycle
│   │   ├── company/         # Company management
│   │   ├── keyword_alert/   # Keyword-based job alerts
│   │   ├── notification/    # In-app + email notifications
│   │   ├── bookmark/        # Polymorphic bookmarks
│   │   ├── search_profile/  # Saved search criteria
│   │   ├── analytics/       # Dashboard stats + charts
│   │   ├── search/          # Elasticsearch full-text search
│   │   └── scheduler/       # Background cron jobs
│   └── pkg/
│       ├── config/          # Env-based configuration
│       ├── database/        # PG, Redis, ES connections
│       ├── middleware/       # JWT auth, CORS, rate limiting, logging
│       ├── email/           # SMTP + HTML email templates
│       └── logger/          # Zap structured logging
│
└── frontend/                # React + Vite + TypeScript
    └── src/
        ├── pages/           # 16 pages (Dashboard, Jobs, Applications, ...)
        ├── components/      # Layout (Sidebar, Topbar), UI
        ├── api/hooks.ts     # TanStack Query hooks for all 15 modules
        ├── stores/          # Zustand (auth + job filters)
        └── lib/axios.ts     # Axios with JWT interceptor + auto-refresh
```

---

## 🛠 Tech Stack

| Layer | Technologies |
|-------|-------------|
| **Frontend** | React 19, TypeScript, Vite, React Router, TanStack Query, Zustand, Tailwind CSS v4, Framer Motion, Recharts |
| **Backend** | Go 1.23, Gin, GORM, PostgreSQL 16 |
| **Cache** | Redis 7 |
| **Search** | Elasticsearch 8.13 |
| **Jobs** | robfig/cron (every 6h scraping, daily digests, interview reminders) |
| **Email** | SMTP + HTML templates |
| **Auth** | JWT (access 15min + refresh 7d tokens) |

---

## 🚀 Quick Start

### Option 1: Docker Compose (recommended)

```bash
# Clone the repo
git clone https://github.com/yourname/careercopilot.git
cd careercopilot

# Create .env
cp .env.example .env
# Edit .env with your SMTP credentials and JWT secrets

# Start everything
docker compose up -d

# Open the app
open http://localhost:3000
```

### Option 2: Local Development

**Backend:**
```bash
cd backend
cp .env.example .env
# Edit .env

# Start deps (PostgreSQL, Redis, Elasticsearch)
docker compose up -d postgres redis elasticsearch

go mod download
go run ./cmd/server
# API running at http://localhost:8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# UI running at http://localhost:3000
```

---

## 🌐 API Reference

Base URL: `http://localhost:8080/api/v1`

All protected routes require: `Authorization: Bearer <token>`

| Module | Endpoints |
|--------|-----------|
| **Auth** | `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`, `GET /auth/me` |
| **Profile** | `GET /profile`, `PUT /profile` |
| **Resume** | `GET /resume`, `POST /resume` |
| **Jobs** | `GET /jobs`, `GET /jobs/:id`, `GET /jobs/matched` |
| **Applications** | `GET/POST /applications`, `PUT /applications/:id`, `DELETE /applications/:id` |
| **Interviews** | `GET/POST /interviews`, `POST /interviews/:id/rounds` |
| **Referrals** | `GET/POST /referrals`, `PUT /referrals/:id` |
| **Companies** | `GET /companies`, `POST/DELETE /companies/:id/watch` |
| **Watchlists** | `GET /watchlists` |
| **Search Profiles** | `GET/POST/PUT/DELETE /search-profiles` |
| **Alerts** | `GET/POST/DELETE /alerts` |
| **Notifications** | `GET /notifications`, `PUT /notifications/:id/read`, `PUT /notifications/read-all` |
| **Analytics** | `GET /analytics/dashboard`, `GET /analytics` |
| **Search** | `GET /search?q=...` |
| **Bookmarks** | `GET/POST/DELETE /bookmarks/:type/:id` |

---

## 📱 Pages

| Page | Route | Description |
|------|-------|-------------|
| Dashboard | `/` | Stats cards + pipeline + upcoming interviews |
| Jobs | `/jobs` | Search + filter + apply from 8 providers |
| Applications | `/applications` | Kanban-style status tracker |
| Interviews | `/interviews` | Timeline with rounds and feedback |
| Referrals | `/referrals` | LinkedIn contact lifecycle tracker |
| Companies | `/companies` | Company search + career page links |
| Watchlists | `/watchlists` | Monitored companies |
| Search Profiles | `/search-profiles` | Saved search criteria |
| Resume | `/resume` | Paste & parse with skill extraction |
| Alerts | `/alerts` | Keyword-based job alerts |
| Bookmarks | `/bookmarks` | Saved jobs/companies/referrals |
| Analytics | `/analytics` | Charts: applications, skills, companies |
| Notifications | `/notifications` | Real-time notification center |
| Settings | `/settings` | Profile, preferences, notification toggles |

---

## ⚙️ Configuration

All configuration is via environment variables:

```env
# App
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=careercopilot
DB_PASSWORD=careercopilot_secret
DB_NAME=careercopilot

# JWT
JWT_ACCESS_SECRET=change_me
JWT_REFRESH_SECRET=change_me_too

# SMTP (optional — for digest emails)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@gmail.com
SMTP_PASSWORD=app_password

# Scheduler
SCHEDULER_JOB_INTERVAL_HOURS=6
SCHEDULER_DIGEST_CRON=0 7 * * *

# Frontend URL (CORS)
FRONTEND_URL=http://localhost:3000
```

---

## 🧠 Key Design Decisions

- **No AI**: Skill extraction uses deterministic regex + keyword matching from a curated skill dictionary. No OpenAI, Gemini, or ML models.
- **Deduplication**: Jobs are deduplicated using SHA-256 hash of `normalize(company + title + location)`, preventing the same job from appearing twice across providers.
- **Job Matching**: Rule-based scoring (max 100 points) based on skill overlap, location match, salary range, experience, and company preferences.
- **Provider Architecture**: Registry pattern — all 8 providers implement a common `Provider` interface. Easy to add new providers.
- **Clean Architecture**: Each module follows `entity → repository → service → controller → routes`.

---

## 📦 Job Providers

| Provider | Type |
|----------|------|
| Greenhouse | API |
| Lever | API |
| Workday | API |
| Indeed | Scraper |
| LinkedIn | Scraper |
| Wellfound | Scraper |
| Naukri | Scraper |
| CareerPage | Generic |

---

## 📄 License

MIT © CareerCopilot
