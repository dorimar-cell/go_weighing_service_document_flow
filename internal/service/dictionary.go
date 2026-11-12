package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/registry"
	"go_weighing_service/internal/repository"
)

// ValidationError — ошибка входных данных (400).
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type DictionaryService struct {
	repo *repository.DictionaryRepository
}

func NewDictionaryService(repo *repository.DictionaryRepository) *DictionaryService {
	return &DictionaryService{repo: repo}
}

// List возвращает страницу справочника.
func (s *DictionaryService) List(
	ctx context.Context,
	code string,
	page int,
	pageSize int,
	search string,
) (*dto.DictionaryListResponse, error) {
	def, err := resolve(code)
	if err != nil {
		return nil, err
	}

	page, pageSize = normalizePaging(page, pageSize)

	result, err := s.repo.List(ctx, def, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	return &dto.DictionaryListResponse{
		Items:    result.Items,
		Page:     page,
		PageSize: pageSize,
		Total:    result.Total,
	}, nil
}

// Create создаёт запись. Обязательные поля проверяются по реестру.
func (s *DictionaryService) Create(
	ctx context.Context,
	code string,
	raw map[string]json.RawMessage,
) (map[string]any, error) {
	def, err := resolve(code)
	if err != nil {
		return nil, err
	}

	values, err := parsePayload(def, raw)
	if err != nil {
		return nil, err
	}

	for _, field := range def.Fields {
		if field.Required && isEmpty(values[field.Key]) {
			return nil, &ValidationError{Message: fmt.Sprintf("Заполните поле «%s»", field.Key)}
		}
	}

	return s.repo.Create(ctx, def, values)
}

// Update меняет только переданные поля.
func (s *DictionaryService) Update(
	ctx context.Context,
	code string,
	id int64,
	raw map[string]json.RawMessage,
) (map[string]any, error) {
	def, err := resolve(code)
	if err != nil {
		return nil, err
	}

	values, err := parsePayload(def, raw)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, &ValidationError{Message: "Нет полей для обновления"}
	}

	for _, field := range def.Fields {
		if !field.Required {
			continue
		}
		if value, ok := values[field.Key]; ok && isEmpty(value) {
			return nil, &ValidationError{Message: fmt.Sprintf("Поле «%s» не может быть пустым", field.Key)}
		}
	}

	return s.repo.Update(ctx, def, id, values)
}

func (s *DictionaryService) Delete(ctx context.Context, code string, id int64) error {
	def, err := resolve(code)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, def, id)
}

func (s *DictionaryService) Restore(ctx context.Context, code string, id int64) error {
	def, err := resolve(code)
	if err != nil {
		return err
	}
	return s.repo.Restore(ctx, def, id)
}

func resolve(code string) (registry.Dictionary, error) {
	def, ok := registry.Get(code)
	if !ok {
		return registry.Dictionary{}, &ValidationError{Message: "Справочник не найден"}
	}
	return def, nil
}

func normalizePaging(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	if pageSize > 200 {
		pageSize = 200
	}
	return page, pageSize
}

func parsePayload(
	def registry.Dictionary,
	raw map[string]json.RawMessage,
) (map[string]any, error) {
	fields := make(map[string]registry.Field, len(def.Fields))
	for _, field := range def.Fields {
		fields[field.Key] = field
	}

	values := make(map[string]any, len(raw))

	for key, rawValue := range raw {
		field, ok := fields[key]
		if !ok {
			return nil, &ValidationError{Message: fmt.Sprintf("Неизвестное поле: %s", key)}
		}

		value, err := parseValue(field, rawValue)
		if err != nil {
			return nil, &ValidationError{
				Message: fmt.Sprintf("Поле «%s»: %s", key, err.Error()),
			}
		}
		values[key] = value
	}

	return values, nil
}

func parseValue(field registry.Field, raw json.RawMessage) (any, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "null" {
		return nil, nil
	}

	switch field.Kind {
	case registry.KindText:
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, fmt.Errorf("ожидается строка")
		}
		return strings.TrimSpace(value), nil

	case registry.KindBool:
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, fmt.Errorf("ожидается true/false")
		}
		return value, nil

	case registry.KindInt:
		var number int64
		if err := json.Unmarshal(raw, &number); err == nil {
			return number, nil
		}
		var text string
		if err := json.Unmarshal(raw, &text); err == nil {
			parsed, convErr := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
			if convErr == nil {
				return parsed, nil
			}
		}
		return nil, fmt.Errorf("ожидается целое число")

	case registry.KindNumeric:
		var number float64
		if err := json.Unmarshal(raw, &number); err == nil {
			return number, nil
		}
		var text string
		if err := json.Unmarshal(raw, &text); err == nil {
			parsed, convErr := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(text), ",", "."), 64)
			if convErr == nil {
				return parsed, nil
			}
		}
		return nil, fmt.Errorf("ожидается число")

	default:
		return nil, fmt.Errorf("неизвестный тип поля")
	}
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed) == ""
	default:
		return false
	}
}
