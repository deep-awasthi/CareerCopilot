-- 009_create_referrals.sql
CREATE TYPE referral_status AS ENUM (
    'not_contacted', 'contacted', 'follow_up',
    'referral_received', 'applied', 'rejected'
);

CREATE TABLE IF NOT EXISTS referrals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id BIGINT REFERENCES companies(id) ON DELETE SET NULL,
    referrer_name VARCHAR(255) NOT NULL,
    referrer_designation VARCHAR(255),
    referrer_department VARCHAR(255),
    referrer_office_location VARCHAR(255),
    referrer_profile_url VARCHAR(500),
    status referral_status DEFAULT 'not_contacted',
    notes TEXT,
    contacted_at TIMESTAMPTZ,
    follow_up_at TIMESTAMPTZ,
    referral_received_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_referrals_user_id ON referrals(user_id);
CREATE INDEX idx_referrals_company_id ON referrals(company_id);
CREATE INDEX idx_referrals_status ON referrals(status);
