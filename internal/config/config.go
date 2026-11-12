package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	Env             string
	Port            string
	DatabaseURL     string
	AuthMode        string
	JWTPublicKeyPath string
	CORSOrigins     []string
	DBMaxConns      int32
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("файл .env не найден, используются переменные окружения")
	}

	cfg := &Config{
		Env:              getEnv("ENV", "local"),
		Port:             getEnv("SERVER_PORT", "8089"),
		DatabaseURL:      getEnv("WEIGHING_DB_URL", "postgresql://weighing:weighing@localhost:5436/weighing_db?sslmode=disable"),
		AuthMode:         getEnv("AUTH_MODE", "dev"),
		JWTPublicKeyPath: getEnv("JWT_PUBLIC_KEY_PATH", "../symfony_documents_flow/config/jwt/public.pem"),
		CORSOrigins:      splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3001")),
		DBMaxConns:       int32(getEnvAsInt("DB_MAX_CONNS", 10)),
	}

	return cfg
}

// ConnectDB создаёт пул подключений и проверяет доступность PostgreSQL.
func ConnectDB(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	poolCfg.MaxConns = cfg.DBMaxConns
	poolCfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value, err := strconv.Atoi(getEnv(key, ""))
	if err != nil {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
