package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type rateBucket struct {
	tokens   float64
	lastSeen time.Time
}

// RateLimiter ограничивает число запросов с одного IP (token bucket).
// Без внешних зависимостей: подходит для одного инстанса сервиса.
func RateLimiter(perSecond float64, burst float64) func(http.Handler) http.Handler {
	var mu sync.Mutex
	buckets := make(map[string]*rateBucket)
	lastCleanup := time.Now()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			now := time.Now()

			mu.Lock()
			if now.Sub(lastCleanup) > 5*time.Minute {
				for key, item := range buckets {
					if now.Sub(item.lastSeen) > 10*time.Minute {
						delete(buckets, key)
					}
				}
				lastCleanup = now
			}

			bucket, ok := buckets[ip]
			if !ok {
				bucket = &rateBucket{tokens: burst, lastSeen: now}
				buckets[ip] = bucket
			}

			elapsed := now.Sub(bucket.lastSeen).Seconds()
			bucket.lastSeen = now
			bucket.tokens += elapsed * perSecond
			if bucket.tokens > burst {
				bucket.tokens = burst
			}

			allowed := bucket.tokens >= 1
			if allowed {
				bucket.tokens--
			}
			mu.Unlock()

			if !allowed {
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"error":"rate_limited","message":"Слишком много запросов"}`,
					http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if index := indexByte(forwarded, ','); index >= 0 {
			forwarded = forwarded[:index]
		}
		if trimmed := trimSpace(forwarded); trimmed != "" {
			return trimmed
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func indexByte(value string, target byte) int {
	for index := 0; index < len(value); index++ {
		if value[index] == target {
			return index
		}
	}
	return -1
}

func trimSpace(value string) string {
	start, end := 0, len(value)
	for start < end && (value[start] == ' ' || value[start] == '\t') {
		start++
	}
	for end > start && (value[end-1] == ' ' || value[end-1] == '\t') {
		end--
	}
	return value[start:end]
}
