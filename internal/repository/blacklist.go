package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// BlacklistMatch ищет совпадение по чёрным спискам транспорта, прицепа и водителя.
func (r *SessionRepository) BlacklistMatch(
	ctx context.Context,
	transport string,
	trailer string,
	driver string,
) (string, string, bool, error) {
	checks := []struct {
		kind  string
		table string
		value string
	}{
		{kind: "Транспортное средство", table: "dir_transport", value: transport},
		{kind: "Прицеп", table: "dir_trailer", value: trailer},
		{kind: "Водитель", table: "dir_driver", value: driver},
	}

	for _, check := range checks {
		value := strings.TrimSpace(check.value)
		if value == "" {
			continue
		}

		var name string
		err := r.pool.QueryRow(ctx, fmt.Sprintf(`
			SELECT name FROM %s
			WHERE lower(name) = lower($1) AND black_list = true AND deleted_at IS NULL
			LIMIT 1`, check.table), value,
		).Scan(&name)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return "", "", false, fmt.Errorf("blacklist %s: %w", check.table, err)
		}

		return check.kind, name, true, nil
	}

	return "", "", false, nil
}
