package api

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

type createTradeIntentRequest struct {
	UserID   string `json:"user_id"`
	MarketID string `json:"market_id"`
	Outcome  string `json:"outcome"`
	Side     string `json:"side"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

func (request createTradeIntentRequest) toRepositoryInput() (repository.CreateTradeIntentInput, error) {
	userID := strings.TrimSpace(request.UserID)
	marketID := strings.TrimSpace(request.MarketID)
	outcome := strings.TrimSpace(request.Outcome)
	side := strings.TrimSpace(request.Side)
	quantityValue := strings.TrimSpace(request.Quantity)
	priceValue := strings.TrimSpace(request.Price)

	if userID == "" || !isUUIDShape(userID) {
		return repository.CreateTradeIntentInput{}, errors.New("user_id is required")
	}
	if marketID == "" || !isUUIDShape(marketID) {
		return repository.CreateTradeIntentInput{}, errors.New("market_id is required")
	}
	if outcome != "YES" && outcome != "NO" {
		return repository.CreateTradeIntentInput{}, errors.New("outcome must be YES or NO")
	}
	if side != "BUY" && side != "SELL" {
		return repository.CreateTradeIntentInput{}, errors.New("side must be BUY or SELL")
	}

	quantity, ok := parseDecimal(quantityValue)
	if !ok || quantity.Sign() <= 0 {
		return repository.CreateTradeIntentInput{}, errors.New("quantity must be greater than zero")
	}

	price, ok := parseDecimal(priceValue)
	if !ok || price.Sign() < 0 || price.Cmp(big.NewRat(1, 1)) > 0 {
		return repository.CreateTradeIntentInput{}, errors.New("price must be between zero and one")
	}

	collateral, ok := decimalString(new(big.Rat).Mul(quantity, price))
	if !ok {
		return repository.CreateTradeIntentInput{}, errors.New("decimal requires more than 18 fractional digits")
	}

	return repository.CreateTradeIntentInput{
		UserID:           userID,
		MarketID:         marketID,
		Outcome:          outcome,
		Side:             side,
		Quantity:         quantityValue,
		Price:            priceValue,
		CollateralAmount: collateral,
	}, nil
}

func registerTradeRoutes(router chi.Router, tradesRepository *repository.TradesRepository, marketsRepository *repository.MarketsRepository) {
	router.Post("/trade-intents", func(w http.ResponseWriter, r *http.Request) {
		var request createTradeIntentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
			return
		}

		input, err := request.toRepositoryInput()
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_trade_intent", "invalid trade intent")
			return
		}

		market, err := marketsRepository.GetMarketByID(r.Context(), input.MarketID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				httpjson.WriteError(w, http.StatusNotFound, "market_not_found", "market not found")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "trade_intent_create_failed", "failed to create trade intent")
			return
		}

		if market.Status != "OPEN" {
			httpjson.WriteError(w, http.StatusBadRequest, "market_not_open", "market is not open")
			return
		}

		trade, err := tradesRepository.CreateTradeIntent(r.Context(), input)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "trade_intent_create_failed", "failed to create trade intent")
			return
		}

		httpjson.WriteJSON(w, http.StatusCreated, map[string]any{
			"trade": newTradeResponse(trade),
			"execution": map[string]string{
				"status": "not_executed",
				"reason": "onchain execution is not implemented in Phase 3 backend MVP",
			},
		})
	})
}
