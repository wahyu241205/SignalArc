package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/wahyu241205/SignalArc/backend/internal/config"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateDatabaseURL(); err != nil {
		log.Fatal().Err(err).Msg("invalid backend config")
	}

	db, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("database startup check failed")
	}
	defer db.Close()

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/schema/validate", func(w http.ResponseWriter, r *http.Request) {
		result, err := db.ValidateSchema(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"status":            "error",
				"migration_version": 0,
				"dirty":             false,
				"missing_tables":    []string{},
			})
			return
		}

		statusCode := http.StatusOK
		if result.Status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		writeJSON(w, statusCode, result)
	})

	log.Info().Str("env", cfg.AppEnv).Str("port", cfg.AppPort).Msg("starting SignalArc API")

	if err := http.ListenAndServe(":"+cfg.AppPort, router); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
