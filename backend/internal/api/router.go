package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

const (
	defaultListLimit    = 50
	defaultMarketsLimit = defaultListLimit
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var decimalPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]{1,18})?$`)

func NewRouter(db *database.DB) http.Handler {
	router := chi.NewRouter()
	marketsRepository := repository.NewMarketsRepository(db)
	positionsRepository := repository.NewPositionsRepository(db)
	resolutionsRepository := repository.NewResolutionsRepository(db)
	settlementsRepository := repository.NewSettlementsRepository(db)
	tradesRepository := repository.NewTradesRepository(db)

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

	router.Get("/agent/markets", func(w http.ResponseWriter, r *http.Request) {
		markets, err := marketsRepository.ListMarkets(r.Context(), defaultMarketsLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "markets_list_failed", "failed to list markets")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"markets": newAgentMarketResponses(markets)})
	})

	router.Post("/markets", func(w http.ResponseWriter, r *http.Request) {
		var request createMarketRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
			return
		}

		input, err := request.toRepositoryInput(time.Now())
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_market_request", "invalid market request")
			return
		}

		market, err := marketsRepository.CreateMarket(r.Context(), input)
		if isForeignKeyViolation(err) {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_creator_user", "creator user is invalid")
			return
		}
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "market_create_failed", "failed to create market")
			return
		}

		httpjson.WriteJSON(w, http.StatusCreated, map[string]any{"market": newMarketResponse(market)})
	})

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
		if errors.Is(err, pgx.ErrNoRows) {
			httpjson.WriteError(w, http.StatusNotFound, "market_not_found", "market not found")
			return
		}
		if err != nil {
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
	if userID == "" || !uuidPattern.MatchString(userID) {
		return repository.CreateTradeIntentInput{}, errors.New("user_id is required")
	}

	marketID := strings.TrimSpace(request.MarketID)
	if marketID == "" || !uuidPattern.MatchString(marketID) {
		return repository.CreateTradeIntentInput{}, errors.New("market_id is required")
	}

	outcome := strings.TrimSpace(request.Outcome)
	if outcome != "YES" && outcome != "NO" {
		return repository.CreateTradeIntentInput{}, errors.New("outcome must be YES or NO")
	}

	side := strings.TrimSpace(request.Side)
	if side != "BUY" && side != "SELL" {
		return repository.CreateTradeIntentInput{}, errors.New("side must be BUY or SELL")
	}

	quantity, quantityRat, err := parseDecimal(request.Quantity)
	if err != nil || quantityRat.Sign() <= 0 {
		return repository.CreateTradeIntentInput{}, errors.New("quantity must be greater than zero")
	}

	price, priceRat, err := parseDecimal(request.Price)
	if err != nil || priceRat.Sign() < 0 || priceRat.Cmp(big.NewRat(1, 1)) > 0 {
		return repository.CreateTradeIntentInput{}, errors.New("price must be between zero and one")
	}

	collateralAmount, err := decimalString(new(big.Rat).Mul(quantityRat, priceRat))
	if err != nil {
		return repository.CreateTradeIntentInput{}, err
	}

	return repository.CreateTradeIntentInput{
		UserID:           userID,
		MarketID:         marketID,
		Outcome:          outcome,
		Side:             side,
		Quantity:         quantity,
		Price:            price,
		CollateralAmount: collateralAmount,
	}, nil
}

func parseDecimal(value string) (string, *big.Rat, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || !decimalPattern.MatchString(trimmed) {
		return "", nil, errors.New("invalid decimal")
	}
	integerPart := strings.SplitN(trimmed, ".", 2)[0]
	if significantIntegerDigits(integerPart) > 18 {
		return "", nil, errors.New("decimal integer part exceeds 18 digits")
	}

	rat, ok := new(big.Rat).SetString(trimmed)
	if !ok {
		return "", nil, errors.New("invalid decimal")
	}

	return trimmed, rat, nil
}

func decimalString(value *big.Rat) (string, error) {
	numerator := new(big.Int).Set(value.Num())
	denominator := new(big.Int).Set(value.Denom())
	integer := new(big.Int)
	remainder := new(big.Int)
	integer.QuoRem(numerator, denominator, remainder)

	if remainder.Sign() == 0 {
		return integer.String(), nil
	}

	digits := strings.Builder{}
	ten := big.NewInt(10)
	for i := 0; i < 18 && remainder.Sign() != 0; i++ {
		remainder.Mul(remainder, ten)
		digit := new(big.Int)
		digit.QuoRem(remainder, denominator, remainder)
		digits.WriteString(digit.String())
	}
	if remainder.Sign() != 0 {
		return "", errors.New("decimal requires more than 18 fractional digits")
	}

	if significantIntegerDigits(integer.String()) > 18 {
		return "", errors.New("decimal integer part exceeds 18 digits")
	}

	fractional := strings.TrimRight(digits.String(), "0")
	if fractional == "" {
		return integer.String(), nil
	}

	return integer.String() + "." + fractional, nil
}

func significantIntegerDigits(value string) int {
	trimmed := strings.TrimLeft(value, "0")
	if trimmed == "" {
		return 0
	}

	return len(trimmed)
}

type createMarketRequest struct {
	CreatorUserID    string  `json:"creator_user_id"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	Category         *string `json:"category"`
	OutcomeYesLabel  *string `json:"outcome_yes_label"`
	OutcomeNoLabel   *string `json:"outcome_no_label"`
	CollateralAsset  *string `json:"collateral_asset"`
	Chain            string  `json:"chain"`
	ResolutionSource *string `json:"resolution_source"`
	OpensAt          *string `json:"opens_at"`
	ClosesAt         string  `json:"closes_at"`
	hasForbiddenKeys bool
}

type tradeResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	MarketID         string    `json:"market_id"`
	Outcome          string    `json:"outcome"`
	Side             string    `json:"side"`
	Quantity         string    `json:"quantity"`
	Price            string    `json:"price"`
	CollateralAmount string    `json:"collateral_amount"`
	FeeAmount        string    `json:"fee_amount"`
	Status           string    `json:"status"`
	TxHash           *string   `json:"tx_hash"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func newTradeResponse(trade repository.Trade) tradeResponse {
	return tradeResponse{
		ID:               trade.ID,
		UserID:           trade.UserID,
		MarketID:         trade.MarketID,
		Outcome:          trade.Outcome,
		Side:             trade.Side,
		Quantity:         trade.Quantity,
		Price:            trade.Price,
		CollateralAmount: trade.CollateralAmount,
		FeeAmount:        trade.FeeAmount,
		Status:           trade.Status,
		TxHash:           nullStringPtr(trade.TxHash),
		CreatedAt:        trade.CreatedAt,
		UpdatedAt:        trade.UpdatedAt,
	}
}

func (request *createMarketRequest) UnmarshalJSON(data []byte) error {
	type createMarketRequestAlias createMarketRequest

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var alias createMarketRequestAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	_, hasStatus := raw["status"]
	_, hasWinningOutcome := raw["winning_outcome"]
	_, hasResolvedAt := raw["resolved_at"]
	_, hasSettledAt := raw["settled_at"]

	*request = createMarketRequest(alias)
	request.hasForbiddenKeys = hasStatus || hasWinningOutcome || hasResolvedAt || hasSettledAt
	return nil
}

func (request createMarketRequest) toRepositoryInput(now time.Time) (repository.CreateMarketInput, error) {
	if request.hasForbiddenKeys {
		return repository.CreateMarketInput{}, errors.New("market lifecycle fields are server-owned")
	}

	creatorUserID := strings.TrimSpace(request.CreatorUserID)
	if creatorUserID == "" || !uuidPattern.MatchString(creatorUserID) {
		return repository.CreateMarketInput{}, errors.New("creator_user_id is required")
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		return repository.CreateMarketInput{}, errors.New("title is required")
	}

	chain := strings.TrimSpace(request.Chain)
	if chain == "" {
		return repository.CreateMarketInput{}, errors.New("chain is required")
	}

	if strings.TrimSpace(request.ClosesAt) == "" {
		return repository.CreateMarketInput{}, errors.New("closes_at is required")
	}
	closesAt, err := time.Parse(time.RFC3339, strings.TrimSpace(request.ClosesAt))
	if err != nil {
		return repository.CreateMarketInput{}, err
	}
	if !closesAt.After(now) {
		return repository.CreateMarketInput{}, errors.New("closes_at must be in the future")
	}

	opensAt := sql.NullTime{}
	if request.OpensAt != nil && strings.TrimSpace(*request.OpensAt) != "" {
		parsedOpensAt, err := time.Parse(time.RFC3339, strings.TrimSpace(*request.OpensAt))
		if err != nil {
			return repository.CreateMarketInput{}, err
		}
		if !parsedOpensAt.Before(closesAt) {
			return repository.CreateMarketInput{}, errors.New("opens_at must be before closes_at")
		}
		opensAt = sql.NullTime{Time: parsedOpensAt, Valid: true}
	}

	return repository.CreateMarketInput{
		CreatorUserID:    creatorUserID,
		Title:            title,
		Description:      optionalString(request.Description),
		Category:         optionalString(request.Category),
		OutcomeYesLabel:  defaultString(request.OutcomeYesLabel, "YES"),
		OutcomeNoLabel:   defaultString(request.OutcomeNoLabel, "NO"),
		CollateralAsset:  defaultString(request.CollateralAsset, "USDC"),
		Chain:            chain,
		ResolutionSource: optionalString(request.ResolutionSource),
		OpensAt:          opensAt,
		ClosesAt:         closesAt,
	}, nil
}

func optionalString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: trimmed, Valid: true}
}

func defaultString(value *string, fallback string) string {
	if value == nil {
		return fallback
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
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

type agentMarketResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Status           string    `json:"status"`
	Category         *string   `json:"category"`
	CollateralAsset  string    `json:"collateral_asset"`
	Chain            string    `json:"chain"`
	ClosesAt         time.Time `json:"closes_at"`
	ResolutionSource *string   `json:"resolution_source"`
}

func newAgentMarketResponses(markets []repository.Market) []agentMarketResponse {
	responses := make([]agentMarketResponse, 0, len(markets))
	for _, market := range markets {
		responses = append(responses, agentMarketResponse{
			ID:               market.ID,
			Title:            market.Title,
			Status:           market.Status,
			Category:         nullStringPtr(market.Category),
			CollateralAsset:  market.CollateralAsset,
			Chain:            market.Chain,
			ClosesAt:         market.ClosesAt,
			ResolutionSource: nullStringPtr(market.ResolutionSource),
		})
	}

	return responses
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
