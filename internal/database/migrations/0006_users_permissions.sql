-- 0006: привязка пользователей портала и права на действия весовой.

ALTER TABLE dir_users
    ADD COLUMN IF NOT EXISTS portal_user_id bigint,
    ADD COLUMN IF NOT EXISTS portal_username text;

CREATE UNIQUE INDEX IF NOT EXISTS ux_dir_users_portal_user
    ON dir_users (portal_user_id)
    WHERE portal_user_id IS NOT NULL;

-- Права на действия (теги W001..W007).
INSERT INTO system_dir_access_permissions (name, control_tag)
SELECT seed.name, seed.tag
FROM (VALUES
    ('Первое взвешивание', 'W001'),
    ('Второе взвешивание', 'W002'),
    ('Поточное взвешивание', 'W003'),
    ('Отмена взвешивания', 'W004'),
    ('Журнал взвешиваний', 'W005'),
    ('Экспорт журнала', 'W006'),
    ('Управление пропусками', 'W007')
) AS seed(name, tag)
WHERE NOT EXISTS (
    SELECT 1 FROM system_dir_access_permissions p WHERE p.control_tag = seed.tag
);

-- Группа операторов.
INSERT INTO dir_user_group (id, name, system_record)
VALUES (2, 'Операторы', false)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('dir_user_group', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM dir_user_group))
);

-- Администраторы: все права весовой.
INSERT INTO com_dir_user_groups_access_permissions (id_dir_user_group, id_system_dir_access_permission)
SELECT 1, p.id
FROM system_dir_access_permissions p
WHERE p.control_tag LIKE 'W%'
  AND NOT EXISTS (
      SELECT 1 FROM com_dir_user_groups_access_permissions ugap
      WHERE ugap.id_dir_user_group = 1 AND ugap.id_system_dir_access_permission = p.id
  );

-- Операторы: взвешивание, поточный режим, журнал, пропуска (без экспорта и отмены).
INSERT INTO com_dir_user_groups_access_permissions (id_dir_user_group, id_system_dir_access_permission)
SELECT 2, p.id
FROM system_dir_access_permissions p
WHERE p.control_tag IN ('W001', 'W002', 'W003', 'W005', 'W007')
  AND NOT EXISTS (
      SELECT 1 FROM com_dir_user_groups_access_permissions ugap
      WHERE ugap.id_dir_user_group = 2 AND ugap.id_system_dir_access_permission = p.id
  );

UPDATE dir_users
SET portal_username = 'system'
WHERE id = 1 AND portal_username IS NULL;
