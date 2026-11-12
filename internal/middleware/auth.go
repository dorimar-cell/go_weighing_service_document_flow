package middleware

import (
	"crypto/rsa"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Auth проверяет JWT портала (RS256) и подставляет пользователя весовой.
// AUTH_MODE=dev использует сервисного администратора.
type Auth struct {
	mode      string
	publicKey *rsa.PublicKey
	users     *UserStore
}

func NewAuth(mode string, publicKeyPath string, users *UserStore) (*Auth, error) {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	if normalized != "jwt" {
		slog.Warn("AUTH_MODE!=jwt: используется сервисный пользователь (только для разработки)")
		return &Auth{mode: "dev", users: users}, nil
	}

	keyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}

	return &Auth{mode: "jwt", publicKey: publicKey, users: users}, nil
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	if a.mode != "jwt" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := a.users.DevUser(r.Context())
			if err != nil {
				slog.Error("dev user lookup failed", "error", err)
				unauthorized(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
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
			unauthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(w)
			return
		}

		portalID := claimInt64(claims, "id")
		if portalID == 0 {
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
			unauthorized(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
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
