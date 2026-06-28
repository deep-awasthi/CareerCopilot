-- 010_create_bookmarks.sql
CREATE TYPE bookmark_type AS ENUM ('job', 'company', 'referral');

CREATE TABLE IF NOT EXISTS bookmarks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type bookmark_type NOT NULL,
    target_id BIGINT NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, type, target_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_type ON bookmarks(type);

-- 010b: Company watchlists
CREATE TABLE IF NOT EXISTS company_watchlists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    notify_new_jobs BOOLEAN DEFAULT TRUE,
    last_notified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, company_id)
);

CREATE INDEX idx_company_watchlists_user_id ON company_watchlists(user_id);
CREATE INDEX idx_company_watchlists_company_id ON company_watchlists(company_id);
