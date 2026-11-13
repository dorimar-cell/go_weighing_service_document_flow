-- Этап 7: терминал весовой по умолчанию и экран станций.

INSERT INTO settings_terminal (
    id, terminal_model_value, interface_type, network_address, multiplier_weight,
    com_name, baudrate, parity, databits, stopbits, handshake, rts, dtr,
    platform_number, id_dir_weight_room, updated_date, is_modbus, ip, port,
    use_socket, socket_ip, socket_port, enable, terminal_interval,
    updated_by_user_id, weight_conditional_zero
) VALUES (
    1, 'Emulator', 'Emulator', '1', 1,
    'COM1', '9600', 'None', '8', 'One', 'None', false, false,
    1, 1, now(), false, '127.0.0.1', '502',
    false, '127.0.0.1', '4001', true, '500',
    1, 0
)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('settings_terminal', 'id'),
    GREATEST(1, (SELECT COALESCE(MAX(id), 1) FROM settings_terminal))
);
