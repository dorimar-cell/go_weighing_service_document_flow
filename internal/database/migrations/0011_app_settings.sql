-- Этап 7 (частично): настройки полей рейса — обязательность и видимость.

ALTER TABLE settings_app
    ADD COLUMN IF NOT EXISTS require_driver boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS require_trailer boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS require_cargo boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS require_counterparty boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS require_carrier boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS require_storage boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS visible_driver boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_trailer boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_cargo boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_counterparty boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_carrier boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_storage boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_operation_type boolean NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS visible_pass boolean NOT NULL DEFAULT true;

INSERT INTO settings_app (
    id, ts_number_mask_use, ts_number_mask, need_operation_type, need_fund,
    need_required_fields, placement_visible, ves_doc_visible, axle_load_visible,
    need_cargo_type, visible_cargo_type, need_placement, visible_fund,
    visible_placement, conditional_zero, use_com_object, use_region,
    union_recognize_skud
) VALUES (
    1, false, '', false, false,
    true, true, true, false,
    false, true, false, true,
    true, 0, false, false,
    false
)
ON CONFLICT (id) DO NOTHING;

-- Право на изменение настроек весовой.
INSERT INTO system_dir_access_permissions (name, control_tag)
SELECT 'Настройки весовой', 'W009'
WHERE NOT EXISTS (
    SELECT 1 FROM system_dir_access_permissions WHERE control_tag = 'W009'
);

INSERT INTO com_dir_user_groups_access_permissions (id_dir_user_group, id_system_dir_access_permission)
SELECT 1, p.id
FROM system_dir_access_permissions p
WHERE p.control_tag = 'W009'
  AND NOT EXISTS (
      SELECT 1 FROM com_dir_user_groups_access_permissions ugap
      WHERE ugap.id_dir_user_group = 1 AND ugap.id_system_dir_access_permission = p.id
  );
