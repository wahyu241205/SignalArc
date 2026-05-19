package api

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
)

type requestIDContextKey struct{}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (recorder *responseRecorder) WriteHeader(statusCode int) {
	recorder.statusCode = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func (recorder *responseRecorder) Write(data []byte) (int, error) {
	if recorder.statusCode == 0 {
		recorder.statusCode = http.StatusOK
	}

	return recorder.ResponseWriter.Write(data)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = newRequestID()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDContextKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		statusCode := recorder.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status_code", statusCode).
			Dur("duration", time.Since(startedAt)).
			Str("remote_addr", r.RemoteAddr).
			Str("request_id", requestIDFromContext(r.Context())).
			Msg("http request")
	})
}

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error().
					Interface("panic", recovered).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("request_id", requestIDFromContext(r.Context())).
					Msg("panic recovered")
				httpjson.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	var bytes [16]byte
	if _, err := cryptorand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}

	return time.Now().UTC().Format("20060102150405.000000000")
}

func requestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}
