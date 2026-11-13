package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/model"
	"go_weighing_service/internal/registry"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrNoFields = errors.New("no fields to update")
)

type DictionaryRepository struct {
	pool *pgxpool.Pool
}

func NewDictionaryRepository(pool *pgxpool.Pool) *DictionaryRepository {
	return &DictionaryRepository{pool: pool}
}

// List возвращает страницу записей с поиском по полям Search.
func (r *DictionaryRepository) List(
	ctx context.Context,
	def registry.Dictionary,
	page int,
	pageSize int,
	search string,
) (*model.DictionaryList, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0, 3)

	if term := strings.TrimSpace(search); term != "" {
		ors := make([]string, 0, len(def.Fields))
		for _, field := range def.Fields {
			if field.Search {
				ors = append(ors, fmt.Sprintf("%s ILIKE $1", field.Column))
			}
		}
		if len(ors) > 0 {
			where = append(where, "("+strings.Join(ors, " OR ")+")")
			args = append(args, "%"+term+"%")
		}
	}

	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.pool.QueryRow(ctx,
		"SELECT count(*) FROM "+def.Table+whereSQL, args...,
	).Scan(&total); err != nil {
		return nil, fmt.Errorf("count %s: %w", def.Table, err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	query := "SELECT " + strings.Join(selectExpressions(def), ", ") +
		" FROM " + def.Table + whereSQL +
		" ORDER BY " + def.OrderBy +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitArg, offsetArg)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", def.Table, err)
	}
	defer rows.Close()

	items := make([]model.DictionaryItem, 0, pageSize)
	for rows.Next() {
		item, err := scanItem(rows, def)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", def.Table, err)
	}

	return &model.DictionaryList{Items: items, Total: total}, nil
}

// Audit пишет действие над справочником.
func (r *DictionaryRepository) Audit(
	ctx context.Context,
	code string,
	recordID int64,
	action string,
	userID int64,
	payload any,
) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal dictionary audit: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO dictionary_audit (code, record_id, action, user_id, user_name, payload)
		VALUES ($1, $2, $3, $4, (SELECT name FROM dir_users WHERE id = $4), $5::jsonb)`,
		code, recordID, action, userID, string(data),
	)
	if err != nil {
		return fmt.Errorf("insert dictionary audit: %w", err)
	}
	return nil
}

// AuditList возвращает журнал аудита справочников.
func (r *DictionaryRepository) AuditList(
	ctx context.Context,
	code string,
	action string,
	page int,
	pageSize int,
) (*dto.DictionaryAuditResponse, error) {
	where := []string{"1 = 1"}
	args := make([]any, 0, 2)

	if strings.TrimSpace(code) != "" {
		args = append(args, strings.TrimSpace(code))
		where = append(where, fmt.Sprintf("code = $%d", len(args)))
	}
	if strings.TrimSpace(action) != "" {
		args = append(args, strings.TrimSpace(action))
		where = append(where, fmt.Sprintf("action = $%d", len(args)))
	}
	whereSQL := " WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM dictionary_audit"+whereSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count dictionary audit: %w", err)
	}

	args = append(args, pageSize)
	limitArg := len(args)
	args = append(args, (page-1)*pageSize)
	offsetArg := len(args)

	rows, err := r.pool.Query(ctx, `
		SELECT id, code, COALESCE(record_id, 0), action, user_id, COALESCE(user_name, ''), payload, created_at
		FROM dictionary_audit`+whereSQL+fmt.Sprintf(
		" ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d", limitArg, offsetArg,
	), args...)
	if err != nil {
		return nil, fmt.Errorf("list dictionary audit: %w", err)
	}
	defer rows.Close()

	items := make([]dto.DictionaryAuditItem, 0, pageSize)
	for rows.Next() {
		var item dto.DictionaryAuditItem
		var payload []byte
		if err := rows.Scan(
			&item.ID, &item.Code, &item.RecordID, &item.Action, &item.UserID,
			&item.UserName, &payload, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dictionary audit: %w", err)
		}
		item.Payload = json.RawMessage(payload)
		items = append(items, item)
	}

	return &dto.DictionaryAuditResponse{
		Items: items, Page: page, PageSize: pageSize, Total: total,
	}, rows.Err()
}

// Export возвращает все записи справочника (до 10000) для выгрузки.
func (r *DictionaryRepository) Export(
	ctx context.Context,
	def registry.Dictionary,
	search string,
) ([]model.DictionaryItem, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0, 1)

	if term := strings.TrimSpace(search); term != "" {
		ors := make([]string, 0, len(def.Fields))
		for _, field := range def.Fields {
			if field.Search {
				ors = append(ors, fmt.Sprintf("%s ILIKE $1", field.Column))
			}
		}
		if len(ors) > 0 {
			where = append(where, "("+strings.Join(ors, " OR ")+")")
			args = append(args, "%"+term+"%")
		}
	}

	query := "SELECT " + strings.Join(selectExpressions(def), ", ") +
		" FROM " + def.Table +
		" WHERE " + strings.Join(where, " AND ") +
		" ORDER BY " + def.OrderBy +
		" LIMIT 10000"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("export %s: %w", def.Table, err)
	}
	defer rows.Close()

	items := make([]model.DictionaryItem, 0, 256)
	for rows.Next() {
		item, err := scanItem(rows, def)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// Create вставляет запись и возвращает её в API-представлении.
func (r *DictionaryRepository) Create(
	ctx context.Context,
	def registry.Dictionary,
	values map[string]any,
) (model.DictionaryItem, error) {
	columns := make([]string, 0, len(def.Fields))
	placeholders := make([]string, 0, len(def.Fields))
	args := make([]any, 0, len(def.Fields))

	for _, field := range def.Fields {
		columns = append(columns, field.Column)
		args = append(args, normalizedValue(field, values[field.Key]))
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		def.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(returningExpressions(def), ", "),
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", def.Table, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotFound
	}
	return scanItem(rows, def)
}

// Update обновляет только переданные поля (PATCH-семантика).
func (r *DictionaryRepository) Update(
	ctx context.Context,
	def registry.Dictionary,
	id int64,
	values map[string]any,
) (model.DictionaryItem, error) {
	sets := make([]string, 0, len(values)+1)
	args := make([]any, 0, len(values)+1)

	for _, field := range def.Fields {
		value, ok := values[field.Key]
		if !ok {
			continue
		}
		args = append(args, normalizedValue(field, value))
		sets = append(sets, fmt.Sprintf("%s = $%d", field.Column, len(args)))
	}

	if len(sets) == 0 {
		return nil, ErrNoFields
	}

	sets = append(sets, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING %s",
		def.Table,
		strings.Join(sets, ", "),
		len(args),
		strings.Join(returningExpressions(def), ", "),
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update %s: %w", def.Table, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotFound
	}
	return scanItem(rows, def)
}

// Delete выполняет мягкое удаление.
func (r *DictionaryRepository) Delete(ctx context.Context, def registry.Dictionary, id int64) error {
	tag, err := r.pool.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL", def.Table),
		id,
	)
	if err != nil {
		return fmt.Errorf("delete %s: %w", def.Table, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Restore снимает мягкое удаление.
func (r *DictionaryRepository) Restore(ctx context.Context, def registry.Dictionary, id int64) error {
	tag, err := r.pool.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET deleted_at = NULL, updated_at = now() WHERE id = $1 AND deleted_at IS NOT NULL", def.Table),
		id,
	)
	if err != nil {
		return fmt.Errorf("restore %s: %w", def.Table, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func selectExpressions(def registry.Dictionary) []string {
	return returningExpressions(def)
}

func returningExpressions(def registry.Dictionary) []string {
	expressions := make([]string, 0, len(def.Fields)+2)
	expressions = append(expressions, "id")

	for _, field := range def.Fields {
		expr := field.Column
		if field.Kind == registry.KindNumeric {
			expr = field.Column + "::float8"
		}
		expressions = append(expressions, expr)
	}

	expressions = append(expressions, "deleted_at IS NOT NULL AS deleted")
	return expressions
}

func scanItem(rows pgx.Rows, def registry.Dictionary) (model.DictionaryItem, error) {
	values, err := rows.Values()
	if err != nil {
		return nil, fmt.Errorf("scan %s: %w", def.Table, err)
	}

	item := make(model.DictionaryItem, len(def.Fields)+2)
	if len(values) > 0 {
		item["id"] = values[0]
	}
	for index, field := range def.Fields {
		item[field.Key] = values[index+1]
	}
	if len(values) > len(def.Fields)+1 {
		item["deleted"] = values[len(def.Fields)+1]
	}

	return item, nil
}

// normalizedValue приводит отсутствующие значения к безопасным дефолтам.
func normalizedValue(field registry.Field, value any) any {
	if value != nil {
		return value
	}
	if field.Kind == registry.KindBool {
		return false
	}
	return nil
}
