package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_weighing_service/internal/app"
	"go_weighing_service/internal/config"
	"go_weighing_service/internal/database"
	"go_weighing_service/internal/handler"
	"go_weighing_service/internal/logger"
	"go_weighing_service/internal/middleware"
	"go_weighing_service/internal/reports"
	"go_weighing_service/internal/repository"
	"go_weighing_service/internal/service"
)

func main() {
	cfg := config.Load()
	logger.Setup(cfg.Env)

	if err := reports.SetTimezone(cfg.ReportTimezone); err != nil {
		slog.Warn("report timezone not applied", "timezone", cfg.ReportTimezone, "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := config.ConnectDB(ctx, cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	loginRepo := repository.NewLoginAuditRepository(pool)

	auth, err := middleware.NewAuth(
		cfg.AuthMode, cfg.JWTPublicKeyPath, middleware.NewUserStore(pool), loginRepo,
	)
	if err != nil {
		slog.Error("auth init failed", "error", err)
		os.Exit(1)
	}

	dictionaryRepo := repository.NewDictionaryRepository(pool)
	dictionaryService := service.NewDictionaryService(dictionaryRepo)
	dictionaryHandler := handler.NewDictionaryHandler(dictionaryService)

	measurementVerifier := service.NewMeasurementVerifier(cfg)
	passRepo := repository.NewPassRepository(pool)
	passService := service.NewPassService(passRepo)
	passHandler := handler.NewPassHandler(passService)

	sessionRepo := repository.NewSessionRepository(pool)
	settingsRepo := repository.NewSettingsRepository(pool)
	sessionService := service.NewSessionService(sessionRepo, passRepo, settingsRepo, measurementVerifier)
	sessionHandler := handler.NewSessionHandler(sessionService)

	settingsService := service.NewSettingsService(settingsRepo)
	extraRepo := repository.NewSettingsExtraRepository(pool)
	extraService := service.NewSettingsExtraService(extraRepo)
	settingsHandler := handler.NewSettingsHandler(settingsService, extraService)
	extraHandler := handler.NewSettingsExtraHandler(extraService, cfg.AgentToken)

	rbacRepo := repository.NewRBACRepository(pool)
	rbacService := service.NewRBACService(rbacRepo)
	rbacHandler := handler.NewRBACHandler(rbacService)

	usersRepo := repository.NewUsersRepository(pool)
	usersService := service.NewUsersService(usersRepo)
	usersHandler := handler.NewUsersHandler(usersService)

	healthHandler := handler.NewHealthHandler(pool)
	metricsHandler := handler.NewMetricsHandler(pool)

	stationRepo := repository.NewStationRepository(pool)
	stationService := service.NewStationService(stationRepo)
	stationHandler := handler.NewStationHandler(stationService, extraService)

	photoRepo := repository.NewPhotoRepository(pool)
	photoService := service.NewPhotoService(photoRepo, cfg.PhotoDir)
	photoHandler := handler.NewPhotoHandler(photoService)

	documentRepo := repository.NewDocumentRepository(pool)
	reportJobRepo := repository.NewReportJobRepository(pool)
	reportService := service.NewReportService(
		sessionRepo, documentRepo, photoRepo, reportJobRepo, cfg.PhotoDir,
		time.Duration(cfg.ReportCacheTTL)*time.Second,
	)
	reportHandler := handler.NewReportHandler(reportService)
	loginAuditHandler := handler.NewLoginAuditHandler(loginRepo)

	router := app.NewRouter(
		app.Handlers{
			Dictionary: dictionaryHandler,
			Session:    sessionHandler,
			Pass:       passHandler,
			Photo:      photoHandler,
			Report:     reportHandler,
			Settings:   settingsHandler,
			RBAC:       rbacHandler,
			Users:      usersHandler,
			Health:     healthHandler,
			Station:    stationHandler,
			Extra:      extraHandler,
			LoginAudit: loginAuditHandler,
			Metrics:    metricsHandler,
		},
		auth,
		cfg.CORSOrigins,
	)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("weighing service started", "port", cfg.Port, "env", cfg.Env, "auth_mode", cfg.AuthMode)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
	slog.Info("weighing service stopped")
}
