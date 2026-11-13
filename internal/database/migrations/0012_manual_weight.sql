-- Этап 4: ручной и нулевой вес.
-- W010 — ручной вес, W011 — нулевой вес (значение 0).

INSERT INTO system_dir_access_permissions (name, control_tag)
SELECT seed.name, seed.tag
FROM (VALUES
    ('Ручной ввод веса', 'W010'),
    ('Нулевой вес', 'W011')
) AS seed(name, tag)
WHERE NOT EXISTS (
    SELECT 1 FROM system_dir_access_permissions p WHERE p.control_tag = seed.tag
);

INSERT INTO com_dir_user_groups_access_permissions (id_dir_user_group, id_system_dir_access_permission)
SELECT 1, p.id
FROM system_dir_access_permissions p
WHERE p.control_tag IN ('W010', 'W011')
  AND NOT EXISTS (
      SELECT 1 FROM com_dir_user_groups_access_permissions ugap
      WHERE ugap.id_dir_user_group = 1 AND ugap.id_system_dir_access_permission = p.id
  );
