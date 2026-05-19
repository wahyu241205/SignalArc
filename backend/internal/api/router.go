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

const (
	defaultListLimit    = 50
	defaultMarketsLimit = defaultListLimit
)

func NewRouter(db *database.DB) http.Handler {
	router := chi.NewRouter()
	marketsRepository := repository.NewMarketsRepository(db)
	positionsRepository := repository.NewPositionsRepository(db)
	resolutionsRepository := repository.NewResolutionsRepository(db)
	settlementsRepository := repository.NewSettlementsRepository(db)

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

	router.Get("/users/{user_id}/positions", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")

		positions, err := positionsRepository.ListPositionsByUserID(r.Context(), userID, defaultListLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "positions_list_failed", "failed to list positions")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"positions": newPositionResponses(positions)})
	})

	router.Get("/markets/{market_id}/positions", func(w http.ResponseWriter, r *http.Request) {
		marketID := chi.URLParam(r, "market_id")

		positions, err := positionsRepository.ListPositionsByMarketID(r.Context(), marketID, defaultListLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "positions_list_failed", "failed to list positions")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"positions": newPositionResponses(positions)})
	})

	router.Get("/markets/{market_id}/resolution", func(w http.ResponseWriter, r *http.Request) {
		marketID := chi.URLParam(r, "market_id")

		resolution, err := resolutionsRepository.GetResolutionByMarketID(r.Context(), marketID)
		if errors.Is(err, pgx.ErrNoRows) {
			httpjson.WriteError(w, http.StatusNotFound, "resolution_not_found", "resolution not found")
			return
		}
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "resolution_get_failed", "failed to get resolution")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"resolution": newResolutionResponse(resolution)})
	})

	router.Get("/users/{user_id}/settlements", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")

		settlements, err := settlementsRepository.ListSettlementsByUserID(r.Context(), userID, defaultListLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "settlements_list_failed", "failed to list settlements")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"settlements": newSettlementResponses(settlements)})
	})

	router.Get("/markets/{market_id}/settlements", func(w http.ResponseWriter, r *http.Request) {
		marketID := chi.URLParam(r, "market_id")

		settlements, err := settlementsRepository.ListSettlementsByMarketID(r.Context(), marketID, defaultListLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "settlements_list_failed", "failed to list settlements")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"settlements": newSettlementResponses(settlements)})
	})

	return router
}

type positionResponse struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	MarketID          string    `json:"market_id"`
	Outcome           string    `json:"outcome"`
	Quantity          string    `json:"quantity"`
	AverageEntryPrice string    `json:"average_entry_price"`
	RealizedPnL       string    `json:"realized_pnl"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func newPositionResponses(positions []repository.Position) []positionResponse {
	responses := make([]positionResponse, 0, len(positions))
	for _, position := range positions {
		responses = append(responses, positionResponse{
			ID:                position.ID,
			UserID:            position.UserID,
			MarketID:          position.MarketID,
			Outcome:           position.Outcome,
			Quantity:          position.Quantity,
			AverageEntryPrice: position.AverageEntryPrice,
			RealizedPnL:       position.RealizedPnL,
			CreatedAt:         position.CreatedAt,
			UpdatedAt:         position.UpdatedAt,
		})
	}

	return responses
}

type resolutionResponse struct {
	ID                string     `json:"id"`
	MarketID          string     `json:"market_id"`
	WinningOutcome    *string    `json:"winning_outcome"`
	Status            string     `json:"status"`
	ResolverType      *string    `json:"resolver_type"`
	EvidenceReference *string    `json:"evidence_reference"`
	ResolvedAt        *time.Time `json:"resolved_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func newResolutionResponse(resolution repository.Resolution) resolutionResponse {
	return resolutionResponse{
		ID:                resolution.ID,
		MarketID:          resolution.MarketID,
		WinningOutcome:    nullStringPtr(resolution.WinningOutcome),
		Status:            resolution.Status,
		ResolverType:      nullStringPtr(resolution.ResolverType),
		EvidenceReference: nullStringPtr(resolution.EvidenceReference),
		ResolvedAt:        nullTimePtr(resolution.ResolvedAt),
		CreatedAt:         resolution.CreatedAt,
		UpdatedAt:         resolution.UpdatedAt,
	}
}

type settlementResponse struct {
	ID           string     `json:"id"`
	MarketID     string     `json:"market_id"`
	UserID       *string    `json:"user_id"`
	ResolutionID *string    `json:"resolution_id"`
	Outcome      *string    `json:"outcome"`
	Amount       string     `json:"amount"`
	Status       string     `json:"status"`
	TxHash       *string    `json:"tx_hash"`
	SettledAt    *time.Time `json:"settled_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func newSettlementResponses(settlements []repository.Settlement) []settlementResponse {
	responses := make([]settlementResponse, 0, len(settlements))
	for _, settlement := range settlements {
		responses = append(responses, settlementResponse{
			ID:           settlement.ID,
			MarketID:     settlement.MarketID,
			UserID:       nullStringPtr(settlement.UserID),
			ResolutionID: nullStringPtr(settlement.ResolutionID),
			Outcome:      nullStringPtr(settlement.Outcome),
			Amount:       settlement.Amount,
			Status:       settlement.Status,
			TxHash:       nullStringPtr(settlement.TxHash),
			SettledAt:    nullTimePtr(settlement.SettledAt),
			CreatedAt:    settlement.CreatedAt,
			UpdatedAt:    settlement.UpdatedAt,
		})
	}

	return responses
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
