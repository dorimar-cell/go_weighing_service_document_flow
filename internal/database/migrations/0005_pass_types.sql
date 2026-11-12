-- 0005: типы регистрации пропусков для MVP.
-- Легаси-семантика: 1 — ключ, 2 — распознавание номера, 3 — штрихкод, 4 — въезд/выезд.

INSERT INTO system_dir_type_of_regweight (id, name)
SELECT value, label
FROM (VALUES
    (1, 'Ключ'),
    (2, 'Распознавание номера'),
    (3, 'Штрихкод'),
    (4, 'Въезд/выезд')
) AS seed(value, label)
WHERE NOT EXISTS (SELECT 1 FROM system_dir_type_of_regweight);

SELECT setval(
    pg_get_serial_sequence('system_dir_type_of_regweight', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM system_dir_type_of_regweight))
);
