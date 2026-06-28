-- 002_create_profiles.sql
CREATE TABLE IF NOT EXISTS profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(20),
    experience_years DECIMAL(4,1) DEFAULT 0,
    current_company VARCHAR(255),
    current_ctc DECIMAL(15,2),
    expected_ctc DECIMAL(15,2),
    notice_period_days INT DEFAULT 0,
    preferred_locations TEXT[],
    preferred_roles TEXT[],
    preferred_skills TEXT[],
    bio TEXT,
    linkedin_url VARCHAR(500),
    github_url VARCHAR(500),
    portfolio_url VARCHAR(500),
    avatar_url VARCHAR(500),
    is_open_to_work BOOLEAN DEFAULT TRUE,
    preferred_work_type VARCHAR(20) DEFAULT 'any',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);
