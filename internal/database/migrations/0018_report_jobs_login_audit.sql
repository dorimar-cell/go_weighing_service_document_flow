-- Этап 6/7: задания отчётов, кэш результатов и аудит входов.

CREATE TABLE IF NOT EXISTS report_jobs (
    id bigserial PRIMARY KEY,
    code varchar(40) NOT NULL,
    format varchar(10) NOT NULL,
    params jsonb NOT NULL DEFAULT '{}'::jsonb,
    status varchar(20) NOT NULL DEFAULT 'pending',
    error text,
    file_name varchar(200),
    content_type varchar(120),
    data bytea,
    created_by_user_id bigint,
    created_at timestamptz NOT NULL DEFAULT now(),
    started_at timestamptz,
    finished_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_report_jobs_created
    ON report_jobs (created_at DESC);

CREATE TABLE IF NOT EXISTS report_cache (
    cache_key varchar(64) PRIMARY KEY,
    code varchar(40) NOT NULL,
    format varchar(10) NOT NULL,
    params jsonb NOT NULL DEFAULT '{}'::jsonb,
    file_name varchar(200) NOT NULL,
    content_type varchar(120) NOT NULL,
    data bytea NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_report_cache_expires
    ON report_cache (expires_at);

CREATE TABLE IF NOT EXISTS login_audit (
    id bigserial PRIMARY KEY,
    user_id bigint,
    username varchar(150),
    ip varchar(64),
    user_agent varchar(300),
    result varchar(20) NOT NULL,
    reason varchar(80),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_login_audit_created
    ON login_audit (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_login_audit_user
    ON login_audit (user_id, created_at DESC);
