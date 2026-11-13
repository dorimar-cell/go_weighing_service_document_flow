package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/middleware"
)

// MetricsHandler — простые метрики Prometheus.
type MetricsHandler struct {
	pool      *pgxpool.Pool
	startedAt time.Time
}

func NewMetricsHandler(pool *pgxpool.Pool) *MetricsHandler {
	return &MetricsHandler{pool: pool, startedAt: time.Now()}
}

// Check: GET /metrics
func (h *MetricsHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbUp := 0
	totalConns := int32(0)
	idleConns := int32(0)

	if h.pool != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := h.pool.Ping(ctx); err == nil {
			dbUp = 1
			stat := h.pool.Stat()
			totalConns = stat.TotalConns()
			idleConns = stat.IdleConns()
		}
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_, _ = fmt.Fprintf(w, "# HELP weighing_up Service is running\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_up gauge\nweighing_up 1\n")
	_, _ = fmt.Fprintf(w, "# HELP weighing_db_up Database is reachable\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_db_up gauge\nweighing_db_up %d\n", dbUp)
	_, _ = fmt.Fprintf(w, "# HELP weighing_db_conns_total Total DB connections\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_db_conns_total gauge\nweighing_db_conns_total %d\n", totalConns)
	_, _ = fmt.Fprintf(w, "# HELP weighing_db_conns_idle Idle DB connections\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_db_conns_idle gauge\nweighing_db_conns_idle %d\n", idleConns)
	_, _ = fmt.Fprintf(w, "# HELP weighing_uptime_seconds Process uptime\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_uptime_seconds gauge\nweighing_uptime_seconds %.0f\n",
		time.Since(h.startedAt).Seconds())

	total, errors := middleware.RequestStats()
	_, _ = fmt.Fprintf(w, "# HELP weighing_http_requests_total Total HTTP requests\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_http_requests_total counter\nweighing_http_requests_total %d\n", total)
	_, _ = fmt.Fprintf(w, "# HELP weighing_http_errors_total HTTP 5xx responses\n")
	_, _ = fmt.Fprintf(w, "# TYPE weighing_http_errors_total counter\nweighing_http_errors_total %d\n", errors)
}
