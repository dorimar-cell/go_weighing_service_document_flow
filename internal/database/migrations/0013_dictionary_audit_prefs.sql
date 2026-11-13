-- Этап 2/5: аудит справочников и серверные настройки журнала.

CREATE TABLE IF NOT EXISTS dictionary_audit (
    id bigserial PRIMARY KEY,
    code varchar(60) NOT NULL,
    record_id bigint,
    action varchar(20) NOT NULL,
    user_id bigint NOT NULL,
    user_name varchar(150),
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dictionary_audit_created
    ON dictionary_audit (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_dictionary_audit_code
    ON dictionary_audit (code, created_at DESC);

CREATE TABLE IF NOT EXISTS user_journal_preferences (
    user_id bigint PRIMARY KEY,
    filters jsonb NOT NULL DEFAULT '{}'::jsonb,
    columns jsonb NOT NULL DEFAULT '[]'::jsonb,
    updated_at timestamptz NOT NULL DEFAULT now()
);
