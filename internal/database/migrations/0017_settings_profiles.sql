-- Этап 7: версионируемые профили настроек (draft -> apply).
-- Легаси-таблицы Settings_* остаются для будущего ETL; рабочее состояние
-- портала и агента хранится в одном JSON-профиле на раздел.

CREATE TABLE IF NOT EXISTS weighing_settings_profile (
    scope varchar(40) PRIMARY KEY,
    draft jsonb NOT NULL DEFAULT '{}'::jsonb,
    applied jsonb NOT NULL DEFAULT '{}'::jsonb,
    version integer NOT NULL DEFAULT 1,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by_user_id bigint,
    updated_by_name varchar(150),
    applied_at timestamptz,
    applied_by_user_id bigint,
    applied_by_name varchar(150)
);

ALTER TABLE settings_audit ADD COLUMN IF NOT EXISTS version integer NOT NULL DEFAULT 0;
