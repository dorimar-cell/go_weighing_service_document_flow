-- 0007: камера по умолчанию для MVP-фото.

INSERT INTO settings_video_surveillance (
    id, enable, name, ip, port, wait_time, try_times,
    stream_type_mini_window, stream_type_full_window,
    platform_number, id_dir_weight_room, updated_date, updated_by_user_id
)
VALUES (
    1, true, 'Emulator', '127.0.0.1', 0, 3000, 3,
    0, 0,
    1, 1, now(), 1
)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('settings_video_surveillance', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM settings_video_surveillance))
);
