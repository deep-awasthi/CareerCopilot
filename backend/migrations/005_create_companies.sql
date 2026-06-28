-- 005_create_companies.sql
CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    domain VARCHAR(255),
    career_page_url VARCHAR(500),
    logo_url VARCHAR(500),
    industry VARCHAR(100),
    size VARCHAR(50),
    headquarters VARCHAR(255),
    description TEXT,
    linkedin_url VARCHAR(500),
    glassdoor_url VARCHAR(500),
    founded_year INT,
    is_active BOOLEAN DEFAULT TRUE,
    last_scraped_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_slug ON companies(slug);
CREATE INDEX idx_companies_domain ON companies(domain);
