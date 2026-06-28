-- 007_create_applications.sql
CREATE TYPE application_status AS ENUM (
    'saved', 'applied', 'interview', 'offer', 'rejected', 'archived'
);

CREATE TABLE IF NOT EXISTS applications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    status application_status DEFAULT 'saved',
    notes TEXT,
    follow_up_date DATE,
    salary_offered DECIMAL(15,2),
    referral_used BOOLEAN DEFAULT FALSE,
    referral_id BIGINT,
    applied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, job_id)
);

CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_job_id ON applications(job_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_follow_up ON applications(follow_up_date);
