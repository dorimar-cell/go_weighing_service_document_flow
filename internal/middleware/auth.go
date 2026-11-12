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

// Auth проверяет JWT портала (RS256, публичный ключ Symfony).
//
// AUTH_MODE=dev отключает проверку — только для локальной разработки.
type Auth struct {
	mode      string
	publicKey *rsa.PublicKey
}

func NewAuth(mode string, publicKeyPath string) (*Auth, error) {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	if normalized != "jwt" {
		slog.Warn("AUTH_MODE!=jwt: проверка токена отключена (только для разработки)")
		return &Auth{mode: "dev"}, nil
	}

	keyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}

	return &Auth{mode: "jwt", publicKey: publicKey}, nil
}

func (a *Auth) Handler(next http.Handler) http.Handler {
	if a.mode != "jwt" {
		return next
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

		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"unauthorized","code":"unauthorized"}`))
}
