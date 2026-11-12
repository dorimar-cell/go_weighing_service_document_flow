package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"go_weighing_service/internal/config"
	"go_weighing_service/internal/dto"
)

// MeasurementVerifier проверяет HMAC-подпись Station Agent.
type MeasurementVerifier struct {
	key []byte
}

func NewMeasurementVerifier(cfg *config.Config) *MeasurementVerifier {
	return &MeasurementVerifier{key: []byte(cfg.AgentSigningKey)}
}

func (v *MeasurementVerifier) Verify(measurement dto.MeasurementDTO) error {
	if strings.TrimSpace(measurement.OperationID) == "" {
		return &ValidationError{Message: "Не указан operationId измерения"}
	}
	if measurement.StationID == "" {
		return &ValidationError{Message: "Не указан stationId"}
	}
	if measurement.Platform <= 0 {
		return &ValidationError{Message: "Не указана платформа"}
	}
	if !measurement.Stable {
		return &ValidationError{Message: "Вес не стабилен"}
	}
	if measurement.WeightKg <= 0 {
		return &ValidationError{Message: "Некорректный вес"}
	}
	if measurement.Checksum == "" {
		return &ValidationError{Message: "Не указана подпись измерения"}
	}

	weight := strconv.FormatFloat(measurement.WeightKg, 'f', 3, 64)
	stable := "0"
	if measurement.Stable {
		stable = "1"
	}

	payload := strings.Join([]string{
		measurement.OperationID,
		measurement.StationID,
		strconv.Itoa(measurement.Platform),
		measurement.Terminal,
		weight,
		stable,
		strconv.FormatInt(measurement.MeasuredAt.UnixMilli(), 10),
		strconv.FormatInt(measurement.Sequence, 10),
	}, "|")

	mac := hmac.New(sha256.New, v.key)
	_, _ = mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(strings.ToLower(measurement.Checksum)), []byte(expected)) {
		return &ValidationError{Message: "Подпись измерения неверна"}
	}

	now := time.Now()
	if measurement.MeasuredAt.After(now.Add(2 * time.Minute)) {
		return &ValidationError{Message: "Время измерения в будущем"}
	}
	if now.Sub(measurement.MeasuredAt) > 15*time.Minute {
		return &ValidationError{Message: "Измерение устарело"}
	}

	return nil
}
