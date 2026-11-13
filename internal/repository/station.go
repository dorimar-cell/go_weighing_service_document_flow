package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// StationRepository — весовые и настройки терминалов.
type StationRepository struct {
	pool *pgxpool.Pool
}

func NewStationRepository(pool *pgxpool.Pool) *StationRepository {
	return &StationRepository{pool: pool}
}

func (r *StationRepository) List(ctx context.Context) ([]dto.StationDTO, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT st.id, COALESCE(wr.id, 0), COALESCE(wr.name, ''), st.platform_number,
		       COALESCE(st.terminal_model_value, ''), COALESCE(st.interface_type, ''),
		       COALESCE(st.com_name, ''), COALESCE(st.baudrate, ''), COALESCE(st.parity, ''),
		       COALESCE(st.databits, ''), COALESCE(st.stopbits, ''), st.multiplier_weight::float8,
		       COALESCE(st.is_modbus, false), COALESCE(st.ip, ''), COALESCE(st.port, ''),
		       COALESCE(st.use_socket, false), COALESCE(st.socket_ip, ''), COALESCE(st.socket_port, ''),
		       st.enable, COALESCE(st.terminal_interval, ''), st.updated_date
		FROM settings_terminal st
		LEFT JOIN dir_weight_room wr ON wr.id = st.id_dir_weight_room
		ORDER BY wr.name, st.platform_number, st.id`)
	if err != nil {
		return nil, fmt.Errorf("list stations: %w", err)
	}
	defer rows.Close()

	items := make([]dto.StationDTO, 0, 4)
	for rows.Next() {
		var item dto.StationDTO
		if err := rows.Scan(
			&item.TerminalID, &item.WeightRoomID, &item.WeightRoomName, &item.PlatformNumber,
			&item.TerminalModel, &item.InterfaceType, &item.ComName, &item.Baudrate,
			&item.Parity, &item.Databits, &item.Stopbits, &item.MultiplierWeight,
			&item.IsModbus, &item.IP, &item.Port, &item.UseSocket, &item.SocketIP,
			&item.SocketPort, &item.Enable, &item.TerminalInterval, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan station: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *StationRepository) Save(
	ctx context.Context,
	id int64,
	req dto.SaveStationRequest,
	userID int64,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE settings_terminal
		SET platform_number = $2, terminal_model_value = $3, interface_type = $4,
		    com_name = $5, baudrate = $6, parity = $7, databits = $8, stopbits = $9,
		    multiplier_weight = $10, is_modbus = $11, ip = $12, port = $13,
		    use_socket = $14, socket_ip = $15, socket_port = $16, enable = $17,
		    terminal_interval = $18, updated_date = now(), updated_by_user_id = $19
		WHERE id = $1`,
		id, req.PlatformNumber, req.TerminalModel, req.InterfaceType, req.ComName,
		req.Baudrate, req.Parity, req.Databits, req.Stopbits, req.MultiplierWeight,
		req.IsModbus, req.IP, req.Port, req.UseSocket, req.SocketIP, req.SocketPort,
		req.Enable, req.TerminalInterval, userID,
	)
	if err != nil {
		return fmt.Errorf("save station: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
