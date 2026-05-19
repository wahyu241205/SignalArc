package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

const defaultMarketsLimit = 50

func NewRouter(db *database.DB) http.Handler {
	router := chi.NewRouter()
	marketsRepository := repository.NewMarketsRepository(db)

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

	router.Get("/markets", func(w http.ResponseWriter, r *http.Request) {
		markets, err := marketsRepository.ListMarkets(r.Context(), defaultMarketsLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "markets_list_failed", "failed to list markets")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"markets": newMarketResponses(markets)})
	})

	router.Get("/markets/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		market, err := marketsRepository.GetMarketByID(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			httpjson.WriteError(w, http.StatusNotFound, "market_not_found", "market not found")
			return
		}
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "market_get_failed", "failed to get market")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"market": newMarketResponse(market)})
	})

	return router
}

type marketResponse struct {
	ID               string     `json:"id"`
	CreatorUserID    string     `json:"creator_user_id"`
	Title            string     `json:"title"`
	Description      *string    `json:"description"`
	Category         *string    `json:"category"`
	Status           string     `json:"status"`
	OutcomeYesLabel  string     `json:"outcome_yes_label"`
	OutcomeNoLabel   string     `json:"outcome_no_label"`
	CollateralAsset  string     `json:"collateral_asset"`
	Chain            string     `json:"chain"`
	ResolutionSource *string    `json:"resolution_source"`
	OpensAt          *time.Time `json:"opens_at"`
	ClosesAt         time.Time  `json:"closes_at"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	SettledAt        *time.Time `json:"settled_at"`
	WinningOutcome   *string    `json:"winning_outcome"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func newMarketResponses(markets []repository.Market) []marketResponse {
	responses := make([]marketResponse, 0, len(markets))
	for _, market := range markets {
		responses = append(responses, newMarketResponse(market))
	}

	return responses
}

func newMarketResponse(market repository.Market) marketResponse {
	return marketResponse{
		ID:               market.ID,
		CreatorUserID:    market.CreatorUserID,
		Title:            market.Title,
		Description:      nullStringPtr(market.Description),
		Category:         nullStringPtr(market.Category),
		Status:           market.Status,
		OutcomeYesLabel:  market.OutcomeYesLabel,
		OutcomeNoLabel:   market.OutcomeNoLabel,
		CollateralAsset:  market.CollateralAsset,
		Chain:            market.Chain,
		ResolutionSource: nullStringPtr(market.ResolutionSource),
		OpensAt:          nullTimePtr(market.OpensAt),
		ClosesAt:         market.ClosesAt,
		ResolvedAt:       nullTimePtr(market.ResolvedAt),
		SettledAt:        nullTimePtr(market.SettledAt),
		WinningOutcome:   nullStringPtr(market.WinningOutcome),
		CreatedAt:        market.CreatedAt,
		UpdatedAt:        market.UpdatedAt,
	}
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}
