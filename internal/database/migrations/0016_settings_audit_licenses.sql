-- Этап 7/8: аудит настроек и реестр лицензий.

CREATE TABLE IF NOT EXISTS settings_audit (
    id bigserial PRIMARY KEY,
    scope varchar(40) NOT NULL,
    entity_id bigint,
    action varchar(20) NOT NULL,
    user_id bigint NOT NULL,
    user_name varchar(150),
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_settings_audit_created
    ON settings_audit (created_at DESC);

CREATE TABLE IF NOT EXISTS weighing_license (
    id bigserial PRIMARY KEY,
    name varchar(150) NOT NULL,
    module_code varchar(40) NOT NULL,
    valid_until date,
    seats bigint NOT NULL DEFAULT 1,
    active boolean NOT NULL DEFAULT true,
    notes text,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    modified_by_user_id bigint,
    modified_date timestamptz
);
