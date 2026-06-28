-- 004_create_search_profiles.sql
CREATE TABLE IF NOT EXISTS search_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    keywords TEXT[],
    experience_min DECIMAL(4,1),
    experience_max DECIMAL(4,1),
    locations TEXT[],
    salary_min DECIMAL(15,2),
    salary_max DECIMAL(15,2),
    is_remote BOOLEAN DEFAULT FALSE,
    is_hybrid BOOLEAN DEFAULT FALSE,
    job_type VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    last_run_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_search_profiles_user_id ON search_profiles(user_id);
CREATE INDEX idx_search_profiles_active ON search_profiles(is_active);
