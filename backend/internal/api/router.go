package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
)

func NewRouter(db *database.DB) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpjson.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			httpjson.WriteError(w, http.StatusServiceUnavailable, "database_unavailable", "database is not reachable")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/schema/validate", func(w http.ResponseWriter, r *http.Request) {
		result, err := db.ValidateSchema(r.Context())
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "schema_validation_failed", "schema validation query failed")
			return
		}

		statusCode := http.StatusOK
		if result.Status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		httpjson.WriteJSON(w, statusCode, result)
	})

	return router
}
