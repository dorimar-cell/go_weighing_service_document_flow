package service

import (
	"context"
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// secretMask — значение, которое UI получает вместо пароля.
const secretMask = "***"

// settingsSecretFields — поля, которые не отдаются в портал и не пишутся в аудит.
var settingsSecretFields = map[string]map[string]bool{
	"email": {"userPassword": true},
	"c1":    {"password": true},
}

// settingsAgentScopes — разделы, которые применяет агент станции.
var settingsAgentScopes = []string{"scoreboard", "speech", "skud", "axle", "kfkh"}

func settingsDefaults(scope string) (map[string]any, bool) {
	switch scope {
	case "scoreboard":
		return map[string]any{
			"enable": false, "model": "CD-3000", "interfaceType": "socket",
			"socketIp": "", "socketPort": 5000, "comName": "COM1", "baudrate": "9600",
			"intervalSec": 5, "platformNumber": 1, "units": "kg", "rts": false, "dtr": false,
		}, true
	case "speech":
		return map[string]any{
			"enable": false, "speechType": 1, "volume": 80, "rate": 0,
			"platformNumber": 1, "messageText": "Взвешивание завершено",
			"messageTextTwo": "", "messageTextThree": "",
		}, true
	case "skud":
		return map[string]any{
			"enableSkud": false, "barrierWorkMode": "1",
			"barrierInOpen": "1", "barrierInClose": "2",
			"barrierOutOpen": "1", "barrierOutClose": "2",
			"enableTrafficLight": false, "trafficLightWorkMode": "1",
			"useTrafficLightDelay": false, "trafficLightDelayValue": 10,
			"enableIrSensors": false, "invertS1": false, "invertS2": false,
			"enableReader": false, "readerWorkMode": "1",
			"interfaceType": "com", "comName": "COM2", "socketIp": "", "socketPort": 0,
		}, true
	case "email":
		return map[string]any{
			"enable": false, "smtpHost": "", "smtpPort": 465, "useSsl": true,
			"userMail": "", "userPassword": "", "adresantName": "", "adresantMail": "",
			"message": "", "messageBody": "", "periodFlag": false, "timePeriod": 60,
			"sendFlag": true, "excelFlag": false,
		}, true
	case "integration":
		return map[string]any{
			"enable": false, "integrationMode": 1,
			"uploadFilePath": "", "downloadFilePath": "",
			"timePeriod": 60, "downloadPeriodTime": 60,
			"withCargo": true, "withCounterparty": true, "withStorage": true,
			"withUsers": false, "withWeighing": true,
		}, true
	case "c1":
		return map[string]any{
			"enable": false, "ip": "", "port": 0, "server": "", "refr": "",
			"userName": "", "password": "", "pathToDb": "", "pathToImage": "",
		}, true
	case "axle":
		return map[string]any{
			"enable": false, "platformNumber": 1, "interval": 500,
			"axleMode": 1, "axleNotification": 0, "limitOfStart": 0.0,
			"interAxle": 0.0, "minAxle": 0.0, "windowSize": 10,
			"periodDelay": 0, "storeData": true, "type": 1,
		}, true
	case "kfkh":
		return map[string]any{"enable": false}, true
	}
	return nil, false
}

// SettingsProfileScopes возвращает поддерживаемые разделы настроек.
func SettingsProfileScopes() []string {
	return []string{"scoreboard", "speech", "skud", "email", "integration", "c1", "axle", "kfkh"}
}

// Profile возвращает профиль раздела: applied + draft, секреты замаскированы.
func (s *SettingsExtraService) Profile(ctx context.Context, scope string) (*dto.SettingsProfile, error) {
	scope = strings.TrimSpace(scope)
	defaults, ok := settingsDefaults(scope)
	if !ok {
		return nil, &ValidationError{Message: "Неизвестный раздел настроек: " + scope}
	}

	stored, err := s.repo.Profile(ctx, scope)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return &dto.SettingsProfile{
			Scope: scope, Draft: cloneValues(defaults), Applied: cloneValues(defaults),
			Version: 0,
		}, nil
	}

	stored.Draft = maskSecrets(scope, stored.Draft)
	stored.Applied = maskSecrets(scope, stored.Applied)
	return stored, nil
}

// Profiles возвращает все разделы, включая ещё не сохранённые значения по умолчанию.
func (s *SettingsExtraService) Profiles(ctx context.Context) ([]dto.SettingsProfile, error) {
	stored, err := s.repo.Profiles(ctx)
	if err != nil {
		return nil, err
	}

	byScope := make(map[string]dto.SettingsProfile, len(stored))
	for _, item := range stored {
		byScope[item.Scope] = item
	}

	scopes := SettingsProfileScopes()
	items := make([]dto.SettingsProfile, 0, len(scopes))
	for _, scope := range scopes {
		defaults, _ := settingsDefaults(scope)
		item, ok := byScope[scope]
		if !ok {
			item = dto.SettingsProfile{
				Scope: scope, Draft: cloneValues(defaults), Applied: cloneValues(defaults),
			}
		} else {
			item.Draft = maskSecrets(scope, item.Draft)
			item.Applied = maskSecrets(scope, item.Applied)
		}
		items = append(items, item)
	}
	return items, nil
}

// SaveDraft сохраняет черновик и пишет аудит с номером версии.
func (s *SettingsExtraService) SaveDraft(
	ctx context.Context,
	scope string,
	values map[string]any,
	userID int64,
) (*dto.SettingsProfile, error) {
	scope = strings.TrimSpace(scope)
	defaults, ok := settingsDefaults(scope)
	if !ok {
		return nil, &ValidationError{Message: "Неизвестный раздел настроек: " + scope}
	}
	if values == nil {
		return nil, &ValidationError{Message: "Передайте значения настроек"}
	}

	current, err := s.repo.Profile(ctx, scope)
	if err != nil {
		return nil, err
	}
	base := cloneValues(defaults)
	if current != nil {
		base = current.Draft
	}

	merged := mergeValues(base, values)
	merged = restoreSecrets(scope, merged, base)
	if err := validateProfile(scope, merged); err != nil {
		return nil, err
	}

	if err := s.repo.SaveDraft(ctx, scope, merged, userID); err != nil {
		return nil, err
	}

	saved, err := s.repo.Profile(ctx, scope)
	if err != nil {
		return nil, err
	}
	if saved != nil {
		_ = s.repo.LogSettingsVersion(ctx, "profile:"+scope, 0, "draft_saved", userID,
			saved.Version, map[string]any{"scope": scope, "values": maskSecrets(scope, merged)})
		saved.Draft = maskSecrets(scope, saved.Draft)
		saved.Applied = maskSecrets(scope, saved.Applied)
	}
	return saved, nil
}

// Apply применяет черновик раздела и поднимает версию.
func (s *SettingsExtraService) Apply(
	ctx context.Context,
	scope string,
	userID int64,
) (*dto.SettingsProfile, error) {
	scope = strings.TrimSpace(scope)
	if _, ok := settingsDefaults(scope); !ok {
		return nil, &ValidationError{Message: "Неизвестный раздел настроек: " + scope}
	}

	current, err := s.repo.Profile(ctx, scope)
	if err != nil {
		return nil, err
	}
	if current == nil {
		defaults, _ := settingsDefaults(scope)
		if err := s.repo.SaveDraft(ctx, scope, defaults, userID); err != nil {
			return nil, err
		}
	}

	applied, err := s.repo.Apply(ctx, scope, userID)
	if err != nil {
		return nil, err
	}

	_ = s.repo.LogSettingsVersion(ctx, "profile:"+scope, 0, "applied", userID,
		applied.Version, map[string]any{"scope": scope, "values": maskSecrets(scope, applied.Applied)})

	applied.Draft = maskSecrets(scope, applied.Draft)
	applied.Applied = maskSecrets(scope, applied.Applied)
	return applied, nil
}

// AgentConfig отдаёт агенту только применённые настройки оборудования.
func (s *SettingsExtraService) AgentConfig(ctx context.Context) (*dto.AgentSettingsResponse, error) {
	scopes, updatedAt, err := s.repo.AppliedProfiles(ctx, settingsAgentScopes)
	if err != nil {
		return nil, err
	}
	return &dto.AgentSettingsResponse{Scopes: scopes, UpdatedAt: updatedAt}, nil
}

func cloneValues(values map[string]any) map[string]any {
	result := make(map[string]any, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}

func mergeValues(base map[string]any, patch map[string]any) map[string]any {
	result := cloneValues(base)
	for key, value := range patch {
		result[key] = value
	}
	return result
}

func maskSecrets(scope string, values map[string]any) map[string]any {
	result := cloneValues(values)
	for field := range settingsSecretFields[scope] {
		if text, ok := result[field].(string); ok && text != "" {
			result[field] = secretMask
		}
	}
	return result
}

func restoreSecrets(scope string, values map[string]any, previous map[string]any) map[string]any {
	result := cloneValues(values)
	for field := range settingsSecretFields[scope] {
		if result[field] == secretMask {
			result[field] = previous[field]
		}
	}
	return result
}

func validateProfile(scope string, values map[string]any) error {
	for key, value := range values {
		switch typed := value.(type) {
		case float64:
			if typed < 0 {
				return &ValidationError{Message: fmt.Sprintf("Поле %s не может быть отрицательным", key)}
			}
			if (strings.Contains(strings.ToLower(key), "port") || strings.HasSuffix(key, "Port")) && typed > 65535 {
				return &ValidationError{Message: fmt.Sprintf("Порт %s вне диапазона 0–65535", key)}
			}
		case string:
			if len(typed) > 700 {
				return &ValidationError{Message: fmt.Sprintf("Поле %s слишком длинное", key)}
			}
		}
	}

	if volume, ok := values["volume"].(float64); ok && (volume < 0 || volume > 100) {
		return &ValidationError{Message: "Громкость должна быть в диапазоне 0–100"}
	}
	if scope == "email" {
		if enabled, _ := values["enable"].(bool); enabled {
			if strings.TrimSpace(stringValue(values["smtpHost"])) == "" {
				return &ValidationError{Message: "Укажите SMTP-сервер"}
			}
			if strings.TrimSpace(stringValue(values["userMail"])) == "" {
				return &ValidationError{Message: "Укажите почту отправителя"}
			}
		}
	}
	return nil
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}
