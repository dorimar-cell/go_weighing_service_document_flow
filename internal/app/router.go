package app

import (
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"go_weighing_service/internal/handler"
	"go_weighing_service/internal/middleware"
)

type Handlers struct {
	Dictionary *handler.DictionaryHandler
	Session    *handler.SessionHandler
	Pass       *handler.PassHandler
	Photo      *handler.PhotoHandler
	Report     *handler.ReportHandler
	Settings   *handler.SettingsHandler
	RBAC       *handler.RBACHandler
	Users      *handler.UsersHandler
	Health     *handler.HealthHandler
	Station    *handler.StationHandler
	Metrics    *handler.MetricsHandler
}

func NewRouter(handlers Handlers, auth *middleware.Auth, corsOrigins []string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(middleware.RequestLogger())
	router.Use(middleware.Metrics)
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.Timeout(15 * time.Second))
	router.Use(middleware.CORS(corsOrigins))

	router.Get("/healthz", handlers.Health.Check)
	router.Get("/metrics", handlers.Metrics.Check)

	router.Route("/spa/api/weighing", func(r chi.Router) {
		r.Use(auth.Handler)

		r.Get("/dictionaries/{code}", handlers.Dictionary.List)
		r.Get("/dictionaries/{code}/export.csv", handlers.Dictionary.Export)
		r.Post("/dictionaries/{code}/import.csv", handlers.Dictionary.Import)
		r.Get("/dictionary-audit", handlers.Dictionary.AuditList)
		r.Get("/dictionary-audit/export.csv", handlers.Dictionary.AuditExport)
		r.Get("/dictionary-audit/export.xlsx", handlers.Dictionary.AuditExportXLSX)
		r.Post("/dictionaries/{code}", handlers.Dictionary.Create)
		r.Patch("/dictionaries/{code}/{id}", handlers.Dictionary.Update)
		r.Delete("/dictionaries/{code}/{id}", handlers.Dictionary.Delete)
		r.Post("/dictionaries/{code}/{id}/restore", handlers.Dictionary.Restore)

		r.Post("/sessions", handlers.Session.Create)
		r.Post("/sessions/flow", handlers.Session.Flow)
		r.Get("/sessions", handlers.Session.List)
		r.Get("/sessions/export", handlers.Session.Export)
		r.Get("/sessions/export.xlsx", handlers.Session.ExportXLSX)
		r.Get("/sessions/{id}", handlers.Session.Get)
		r.Patch("/sessions/{id}", handlers.Session.Update)
		r.Get("/sessions/{id}/measurements", handlers.Session.Measurements)
		r.Get("/sessions/{id}/audit", handlers.Session.Audit)
		r.Post("/sessions/{id}/second", handlers.Session.Close)
		r.Post("/sessions/{id}/cancel", handlers.Session.Cancel)
		r.Post("/sessions/{id}/restore", handlers.Session.Restore)
		r.Post("/sessions/bulk", handlers.Session.Bulk)
		r.Post("/sessions/{id}/photos", handlers.Photo.Upload)
		r.Get("/sessions/{id}/photos", handlers.Photo.List)
		r.Get("/sessions/{id}/photos/{photoId}/file", handlers.Photo.File)
		r.Get("/sessions/{id}/documents/{kind}", handlers.Report.GetDocument)
		r.Put("/sessions/{id}/documents/{kind}", handlers.Report.SaveDocument)

		r.Get("/permissions", handlers.Session.Permissions)
		r.Get("/audit", handlers.Session.AuditJournal)
		r.Get("/audit/export.csv", handlers.Session.AuditExport)
		r.Get("/audit/export.xlsx", handlers.Session.AuditExportXLSX)

		r.Get("/settings/app", handlers.Settings.GetApp)
		r.Put("/settings/app", handlers.Settings.SaveApp)

		r.Get("/rbac", handlers.RBAC.Matrix)
		r.Put("/rbac/groups/{id}", handlers.RBAC.SaveGroup)

		r.Get("/users", handlers.Users.List)
		r.Get("/weight-rooms", handlers.Users.WeightRooms)
		r.Put("/users/{id}/weight-room", handlers.Users.SetWeightRoom)
		r.Put("/users/{id}", handlers.Users.Update)
		r.Get("/journal/preferences", handlers.Users.GetPreferences)
		r.Put("/journal/preferences", handlers.Users.SavePreferences)

		r.Get("/stations", handlers.Station.List)
		r.Put("/stations/{id}", handlers.Station.Save)

		r.Get("/reports", handlers.Report.Catalog)
		r.Get("/reports/settings", handlers.Report.GetSettings)
		r.Put("/reports/settings", handlers.Report.SaveSettings)
		r.Get("/reports/{code}/run", handlers.Report.Run)

		r.Post("/passes", handlers.Pass.Create)
		r.Get("/passes", handlers.Pass.List)
		r.Get("/passes/{id}", handlers.Pass.Get)
		r.Patch("/passes/{id}/enabled", handlers.Pass.SetEnabled)
	})

	return router
}
