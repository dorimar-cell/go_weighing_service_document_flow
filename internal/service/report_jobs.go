package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strings"
	"time"

	"go_weighing_service/internal/dto"
	"go_weighing_service/internal/reports"
)

const reportJobTimeout = 2 * time.Minute

// reportCacheKey считает ключ по коду, формату и всем параметрам запроса.
// Порядок параметров не влияет на результат.
func reportCacheKey(code string, format string, query url.Values) string {
	keys := make([]string, 0, len(query))
	for key := range query {
		if key == "nocache" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString(code)
	builder.WriteByte('\x00')
	builder.WriteString(format)
	for _, key := range keys {
		values := append([]string(nil), query[key]...)
		sort.Strings(values)
		builder.WriteByte('\x01')
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(strings.Join(values, "\x02"))
	}

	sum := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

// EnqueueReport ставит отчёт в очередь и сразу возвращает задание.
func (s *ReportService) EnqueueReport(
	ctx context.Context,
	code string,
	query url.Values,
	userID int64,
) (*dto.ReportJob, error) {
	if s.jobs == nil {
		return nil, &ValidationError{Message: "Очередь отчётов не настроена"}
	}
	if _, ok := reports.Definition(code); !ok {
		return nil, &ValidationError{Message: "Отчёт не найден"}
	}

	format := strings.ToLower(strings.TrimSpace(query.Get("format")))
	if format == "" {
		format = "html"
	}

	job, err := s.jobs.CreateJob(ctx, code, format, map[string][]string(query), userID)
	if err != nil {
		return nil, err
	}

	go s.runJob(job.ID, code, query, format)
	return job, nil
}

func (s *ReportService) runJob(id int64, code string, query url.Values, format string) {
	ctx, cancel := context.WithTimeout(context.Background(), reportJobTimeout)
	defer cancel()

	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Error("report job panic", "job", id, "panic", recovered)
			_ = s.jobs.FailJob(ctx, id, "Внутренняя ошибка генерации отчёта")
		}
	}()

	if err := s.jobs.MarkJobRunning(ctx, id); err != nil {
		slog.Error("report job mark running", "job", id, "error", err)
		return
	}

	file, err := s.generate(ctx, code, query, format)
	if err != nil {
		_ = s.jobs.FailJob(ctx, id, err.Error())
		return
	}

	if err := s.jobs.CompleteJob(ctx, id, file.Name, file.ContentType, file.Data); err != nil {
		slog.Error("report job complete", "job", id, "error", err)
	}
}

// ReportJob возвращает статус задания.
func (s *ReportService) ReportJob(ctx context.Context, id int64) (*dto.ReportJob, error) {
	if s.jobs == nil {
		return nil, &ValidationError{Message: "Очередь отчётов не настроена"}
	}
	return s.jobs.Job(ctx, id)
}

// ReportJobFile отдаёт готовый файл задания.
func (s *ReportService) ReportJobFile(ctx context.Context, id int64) (*ReportFile, error) {
	if s.jobs == nil {
		return nil, &ValidationError{Message: "Очередь отчётов не настроена"}
	}
	job, err := s.jobs.Job(ctx, id)
	if err != nil {
		return nil, err
	}
	if job.Status != "done" {
		return nil, &ValidationError{Message: "Отчёт ещё не готов: " + job.Status}
	}
	file, err := s.jobs.JobFile(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ReportFile{
		Name: file.FileName, ContentType: file.ContentType, Data: file.Data,
	}, nil
}

// ReportJobs возвращает последние задания.
func (s *ReportService) ReportJobs(ctx context.Context, limit int) ([]dto.ReportJob, error) {
	if s.jobs == nil {
		return nil, &ValidationError{Message: "Очередь отчётов не настроена"}
	}
	items, err := s.jobs.Jobs(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list report jobs: %w", err)
	}
	return items, nil
}
