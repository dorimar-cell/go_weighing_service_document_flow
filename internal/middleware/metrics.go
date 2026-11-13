package middleware

import (
	"net/http"
	"sync/atomic"
)

var (
	requestTotal  atomic.Int64
	requestErrors atomic.Int64
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(data)
}

// Metrics считает общее число запросов и ответов 5xx.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTotal.Add(1)
		recorder := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		if recorder.status >= 500 {
			requestErrors.Add(1)
		}
	})
}

// RequestStats возвращает счётчики HTTP-запросов.
func RequestStats() (total int64, errors int64) {
	return requestTotal.Load(), requestErrors.Load()
}
