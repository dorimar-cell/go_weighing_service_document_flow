package app

import (
	"net/http"
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
}

func NewRouter(handlers Handlers, auth *middleware.Auth, corsOrigins []string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(middleware.RequestLogger())
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.Timeout(15 * time.Second))
	router.Use(middleware.CORS(corsOrigins))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok","service":"weighing"}`))
	})

	router.Route("/spa/api/weighing", func(r chi.Router) {
		r.Use(auth.Handler)

		r.Get("/dictionaries/{code}", handlers.Dictionary.List)
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
		r.Get("/sessions/{id}/measurements", handlers.Session.Measurements)
		r.Post("/sessions/{id}/second", handlers.Session.Close)
		r.Post("/sessions/{id}/cancel", handlers.Session.Cancel)
		r.Post("/sessions/{id}/photos", handlers.Photo.Upload)
		r.Get("/sessions/{id}/photos", handlers.Photo.List)
		r.Get("/sessions/{id}/photos/{photoId}/file", handlers.Photo.File)
		r.Get("/sessions/{id}/documents/{kind}", handlers.Report.GetDocument)
		r.Put("/sessions/{id}/documents/{kind}", handlers.Report.SaveDocument)

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
