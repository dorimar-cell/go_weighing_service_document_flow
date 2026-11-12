-- 0001: справочники весовой.
--
-- Названия таблиц повторяют легаси-имена (dir_*), но колонки — snake_case.
-- Каждая таблица: id, предметные поля, deleted_at (мягкое удаление), created_at/updated_at.

CREATE TABLE IF NOT EXISTS dir_transport (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    black_list boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_trailer (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    black_list boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_type_of_transport (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_type_of_trailer (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_driver (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    phone text,
    inn text,
    passport_series text,
    passport_number text,
    license_series text,
    license_number text,
    black_list boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_cargo (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    price_per_tonne numeric(14, 3),
    selling_price numeric(14, 3),
    is_danger boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_type_of_cargo (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    cargo_id bigint,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_storage (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    company text,
    phone text,
    address text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_placement (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_fund (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_carrier (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    inn text,
    kpp text,
    phone text,
    address text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_counterparty (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    inn text,
    kpp text,
    phone text,
    director text,
    address text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_provider (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    inn text,
    phone text,
    address text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_type_of_operation (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    direction bigint,
    main boolean NOT NULL DEFAULT false,
    system_record boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_kfkh_plantation (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_kfkh_combine (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    make text,
    scope numeric(14, 3),
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_kfkh_combiner (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    talon text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_weight_room (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dir_user_group (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    system_record boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS system_dir_type_of_pass (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS system_dir_identifier_state (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dir_transport_deleted ON dir_transport (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_trailer_deleted ON dir_trailer (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_type_of_transport_deleted ON dir_type_of_transport (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_type_of_trailer_deleted ON dir_type_of_trailer (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_driver_deleted ON dir_driver (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_cargo_deleted ON dir_cargo (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_type_of_cargo_deleted ON dir_type_of_cargo (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_storage_deleted ON dir_storage (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_placement_deleted ON dir_placement (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_fund_deleted ON dir_fund (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_carrier_deleted ON dir_carrier (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_counterparty_deleted ON dir_counterparty (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_provider_deleted ON dir_provider (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_type_of_operation_deleted ON dir_type_of_operation (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_kfkh_plantation_deleted ON dir_kfkh_plantation (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_kfkh_combine_deleted ON dir_kfkh_combine (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_kfkh_combiner_deleted ON dir_kfkh_combiner (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_weight_room_deleted ON dir_weight_room (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dir_user_group_deleted ON dir_user_group (deleted_at);
CREATE INDEX IF NOT EXISTS idx_system_dir_type_of_pass_deleted ON system_dir_type_of_pass (deleted_at);
CREATE INDEX IF NOT EXISTS idx_system_dir_identifier_state_deleted ON system_dir_identifier_state (deleted_at);

-- Системные справочники: значения повторяют легаси-семантику.
INSERT INTO system_dir_type_of_pass (name)
SELECT value FROM (VALUES ('Постоянный'), ('Временный'), ('Разовый')) AS seed(value)
WHERE NOT EXISTS (SELECT 1 FROM system_dir_type_of_pass);

INSERT INTO system_dir_identifier_state (name)
SELECT value FROM (VALUES ('Действителен'), ('Просрочен'), ('Выдан'), ('Не выдан'), ('Деактивирован')) AS seed(value)
WHERE NOT EXISTS (SELECT 1 FROM system_dir_identifier_state);

INSERT INTO dir_type_of_operation (name, direction, main, system_record)
SELECT name, direction, is_main, false
FROM (VALUES
    ('Приход', 1, true),
    ('Расход', 2, true),
    ('Перемещение', 3, false)
) AS seed(name, direction, is_main)
WHERE NOT EXISTS (SELECT 1 FROM dir_type_of_operation);
