-- 008_create_interviews.sql
CREATE TYPE interview_stage AS ENUM (
    'applied', 'recruiter_call', 'online_assessment',
    'technical_round_1', 'technical_round_2', 'system_design',
    'manager_round', 'hr_round', 'offer', 'rejected'
);

CREATE TYPE interview_result AS ENUM ('pending', 'passed', 'failed', 'cancelled');

CREATE TABLE IF NOT EXISTS interviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_id BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    current_stage interview_stage DEFAULT 'applied',
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(application_id)
);

CREATE TABLE IF NOT EXISTS interview_rounds (
    id BIGSERIAL PRIMARY KEY,
    interview_id BIGINT NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    stage interview_stage NOT NULL,
    scheduled_at TIMESTAMPTZ,
    duration_minutes INT,
    interviewer_name VARCHAR(255),
    interviewer_role VARCHAR(255),
    interviewer_linkedin VARCHAR(500),
    feedback TEXT,
    notes TEXT,
    result interview_result DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_interviews_user_id ON interviews(user_id);
CREATE INDEX idx_interviews_application_id ON interviews(application_id);
CREATE INDEX idx_interview_rounds_interview_id ON interview_rounds(interview_id);
CREATE INDEX idx_interview_rounds_scheduled_at ON interview_rounds(scheduled_at);
