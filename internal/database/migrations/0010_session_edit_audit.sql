-- Этап 5/6: редактирование рейса и история изменений.

CREATE TABLE IF NOT EXISTS weighing_audit (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    action varchar(40) NOT NULL,
    user_id bigint NOT NULL,
    user_name varchar(150),
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_weighing_audit_weighing
    ON weighing_audit (id_weighing, created_at DESC);

-- Право на редактирование завершённых рейсов.
INSERT INTO system_dir_access_permissions (name, control_tag)
SELECT 'Редактирование взвешиваний', 'W008'
WHERE NOT EXISTS (
    SELECT 1 FROM system_dir_access_permissions WHERE control_tag = 'W008'
);

INSERT INTO com_dir_user_groups_access_permissions (id_dir_user_group, id_system_dir_access_permission)
SELECT 1, p.id
FROM system_dir_access_permissions p
WHERE p.control_tag = 'W008'
  AND NOT EXISTS (
      SELECT 1 FROM com_dir_user_groups_access_permissions ugap
      WHERE ugap.id_dir_user_group = 1 AND ugap.id_system_dir_access_permission = p.id
  );
