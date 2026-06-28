-- 006_create_jobs.sql
CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    external_id VARCHAR(255),
    company_id BIGINT REFERENCES companies(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    short_description TEXT,
    location VARCHAR(255),
    locations TEXT[],
    is_remote BOOLEAN DEFAULT FALSE,
    is_hybrid BOOLEAN DEFAULT FALSE,
    employment_type VARCHAR(50),
    experience_min DECIMAL(4,1),
    experience_max DECIMAL(4,1),
    salary_min DECIMAL(15,2),
    salary_max DECIMAL(15,2),
    salary_currency VARCHAR(10) DEFAULT 'INR',
    skills TEXT[],
    application_url VARCHAR(1000),
    posted_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    dedup_hash VARCHAR(64) UNIQUE,
    view_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_jobs_company_id ON jobs(company_id);
CREATE INDEX idx_jobs_title ON jobs USING gin(to_tsvector('english', title));
CREATE INDEX idx_jobs_skills ON jobs USING gin(skills);
CREATE INDEX idx_jobs_locations ON jobs USING gin(locations);
CREATE INDEX idx_jobs_posted_at ON jobs(posted_at DESC);
CREATE INDEX idx_jobs_is_active ON jobs(is_active);
CREATE INDEX idx_jobs_dedup_hash ON jobs(dedup_hash);

-- Job sources (which portals list this job)
CREATE TABLE IF NOT EXISTS job_sources (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    provider VARCHAR(100) NOT NULL,
    external_id VARCHAR(255),
    source_url VARCHAR(1000),
    raw_data JSONB,
    scraped_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(job_id, provider)
);

CREATE INDEX idx_job_sources_job_id ON job_sources(job_id);
CREATE INDEX idx_job_sources_provider ON job_sources(provider);
