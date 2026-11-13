-- Этап 5/6: пресеты фильтров журнала и избранные отчёты.

ALTER TABLE user_journal_preferences
    ADD COLUMN IF NOT EXISTS presets jsonb NOT NULL DEFAULT '[]'::jsonb;

CREATE TABLE IF NOT EXISTS user_report_favorites (
    user_id bigint NOT NULL,
    code varchar(60) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, code)
);
