package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginAuditor записывает успешные и неуспешные попытки входа.
type LoginAuditor interface {
	RecordLogin(
		ctx context.Context,
		userID int64,
		username string,
		ip string,
		userAgent string,
		result string,
		reason string,
	) error
}

// Auth проверяет JWT портала (RS256) и подставляет пользователя весовой.
// AUTH_MODE=dev использует сервисного администратора.
type Auth struct {
	mode      string
	publicKey *rsa.PublicKey
	users     *UserStore
	auditor   LoginAuditor
	seen      sync.Map
}

func NewAuth(
	mode string,
	publicKeyPath string,
	users *UserStore,
	auditor LoginAuditor,
) (*Auth, error) {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	if normalized != "jwt" {
		slog.Warn("AUTH_MODE!=jwt: используется сервисный пользователь (только для разработки)")
		return &Auth{mode: "dev", users: users, auditor: auditor}, nil
	}

	keyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}

	return &Auth{mode: "jwt", publicKey: publicKey, users: users, auditor: auditor}, nil
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	if a.mode != "jwt" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := a.users.DevUser(r.Context())
			if err != nil {
				slog.Error("dev user lookup failed", "error", err)
				a.auditLogin(r, 0, "", "failed", "dev_user_lookup")
				unauthorized(w)
				return
			}
			a.auditLogin(r, user.ID, user.Username, "success", "dev")
			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Агент станции не имеет JWT портала: для выдачи применённой конфигурации
		// достаточно X-Agent-Token. Неверный токен отклонит обработчик (401).
		if strings.HasSuffix(r.URL.Path, "/settings/agent-config") &&
			strings.TrimSpace(r.Header.Get("X-Agent-Token")) != "" {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			a.auditLogin(r, 0, "", "failed", "missing_token")
			unauthorized(w)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return a.publicKey, nil
		}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

		if err != nil || !token.Valid {
			a.auditLogin(r, 0, "", "failed", "invalid_token")
			unauthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			a.auditLogin(r, 0, "", "failed", "invalid_claims")
			unauthorized(w)
			return
		}

		portalID := claimInt64(claims, "id")
		if portalID == 0 {
			a.auditLogin(r, 0, "", "failed", "missing_portal_id")
			unauthorized(w)
			return
		}

		username := claimString(claims, "username")
		if username == "" {
			username = claimString(claims, "email")
		}

		user, err := a.users.Resolve(r.Context(), portalID, username)
		if err != nil {
			slog.Error("resolve portal user failed", "portal_id", portalID, "error", err)
			a.auditLogin(r, 0, username, "failed", "unknown_user")
			unauthorized(w)
			return
		}

		a.auditLogin(r, user.ID, firstValue(user.Username, username), "success", "jwt")
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}

// auditLogin пишет вход в БД не чаще одного раза в 5 минут на пользователя,
// IP и результат, чтобы каждый API-запрос не превращался в запись журнала.
func (a *Auth) auditLogin(
	r *http.Request,
	userID int64,
	username string,
	result string,
	reason string,
) {
	if a.auditor == nil {
		return
	}

	ip := clientIPHeader(r)
	key := fmt.Sprintf("%d|%s|%s", userID, ip, result)
	now := time.Now()
	if last, ok := a.seen.Load(key); ok {
		if lastTime, valid := last.(time.Time); valid && now.Sub(lastTime) < 5*time.Minute {
			return
		}
	}
	a.seen.Store(key, now)

	userAgent := r.UserAgent()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := a.auditor.RecordLogin(ctx, userID, username, ip, userAgent, result, reason); err != nil {
			slog.Error("login audit failed", "error", err)
		}
	}()
}

func clientIPHeader(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		first := strings.TrimSpace(strings.Split(forwarded, ",")[0])
		if first != "" {
			return first
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	return clientIP(r)
}

func claimInt64(claims jwt.MapClaims, key string) int64 {
	value, ok := claims[key]
	if !ok {
		return 0
	}

	switch typed := value.(type) {
	case float64:
		return int64(typed)
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}

func claimString(claims jwt.MapClaims, key string) string {
	value, ok := claims[key]
	if !ok {
		return ""
	}
	text, _ := value.(string)
	return text
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"unauthorized","code":"unauthorized"}`))
}

func firstValue(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
