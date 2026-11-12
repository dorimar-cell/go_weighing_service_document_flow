-- 0004: измерения с подписью, сервисный пользователь и цикл взвешивания.
--
-- Измерение хранится отдельно от weighing: сервер проверяет HMAC и
-- идемпотентность по operation_id, а не доверяет весу из формы.

CREATE TABLE IF NOT EXISTS weighing_measurement (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    kind smallint NOT NULL CHECK (kind IN (1, 2)),
    operation_id uuid NOT NULL UNIQUE,
    weight_kg numeric(12, 3) NOT NULL,
    stable boolean NOT NULL,
    terminal text,
    platform integer NOT NULL,
    measured_at timestamptz NOT NULL,
    checksum text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE weighing_measurement
    ADD CONSTRAINT fk_measurement_weighing
    FOREIGN KEY (id_weighing) REFERENCES weighing(id)
    DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX IF NOT EXISTS idx_measurement_weighing ON weighing_measurement (id_weighing, kind);

-- Сервисный пользователь и весовая для MVP: позже заменяются реальными
-- пользователями портала и привязкой станции.
INSERT INTO dir_weight_room (id, name, description)
VALUES (1, 'Весовая №1', 'Пост по умолчанию')
ON CONFLICT (id) DO NOTHING;

INSERT INTO dir_user_group (id, name, system_record)
VALUES (1, 'Администраторы', true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO dir_users (
    id, name, password_hash, legacy_password_md5, deleted,
    id_dir_weight_room, id_dir_user_groups,
    created_by_user_id, created_date, ref
)
VALUES (
    1, 'system', NULL, NULL, false,
    1, 1,
    1, now(), gen_random_uuid()
)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('dir_weight_room', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM dir_weight_room))
);

SELECT setval(
    pg_get_serial_sequence('dir_user_group', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM dir_user_group))
);

SELECT setval(
    pg_get_serial_sequence('dir_users', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM dir_users))
);
