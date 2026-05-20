package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

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

type attachMarketContractRequest struct {
	MarketContractAddress  string `json:"market_contract_address"`
	MarketDeploymentTxHash string `json:"market_deployment_tx_hash"`
	MarketFactoryAddress   string `json:"market_factory_address"`
	ResolverAddress        string `json:"resolver_address"`
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
	_, hasMarketContractAddress := raw["market_contract_address"]
	_, hasMarketDeploymentTxHash := raw["market_deployment_tx_hash"]
	_, hasMarketFactoryAddress := raw["market_factory_address"]
	_, hasResolverAddress := raw["resolver_address"]
	_, hasOnchainDeploymentStatus := raw["onchain_deployment_status"]

	*request = createMarketRequest(alias)
	request.hasForbiddenKeys = hasStatus ||
		hasWinningOutcome ||
		hasResolvedAt ||
		hasSettledAt ||
		hasMarketContractAddress ||
		hasMarketDeploymentTxHash ||
		hasMarketFactoryAddress ||
		hasResolverAddress ||
		hasOnchainDeploymentStatus
	return nil
}

func (request createMarketRequest) toRepositoryInput(now time.Time) (repository.CreateMarketInput, error) {
	if request.hasForbiddenKeys {
		return repository.CreateMarketInput{}, errors.New("market lifecycle fields are server-owned")
	}

	creatorUserID := strings.TrimSpace(request.CreatorUserID)
	if creatorUserID == "" || !isUUIDShape(creatorUserID) {
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
		Status:           newMarketStatus(now, opensAt),
		OutcomeYesLabel:  defaultString(request.OutcomeYesLabel, "YES"),
		OutcomeNoLabel:   defaultString(request.OutcomeNoLabel, "NO"),
		CollateralAsset:  defaultString(request.CollateralAsset, "USDC"),
		Chain:            chain,
		ResolutionSource: optionalString(request.ResolutionSource),
		OpensAt:          opensAt,
		ClosesAt:         closesAt,
	}, nil
}

func newMarketStatus(now time.Time, opensAt sql.NullTime) string {
	if opensAt.Valid && opensAt.Time.After(now) {
		return "DRAFT"
	}

	return "OPEN"
}

func (request attachMarketContractRequest) toRepositoryInput() (repository.AttachMarketContractInput, error) {
	marketContractAddress := strings.TrimSpace(request.MarketContractAddress)
	marketDeploymentTxHash := strings.TrimSpace(request.MarketDeploymentTxHash)
	marketFactoryAddress := strings.TrimSpace(request.MarketFactoryAddress)
	resolverAddress := strings.TrimSpace(request.ResolverAddress)

	if !isEVMAddressShape(marketContractAddress) {
		return repository.AttachMarketContractInput{}, errors.New("market_contract_address must be an EVM address")
	}
	if !isEVMTxHashShape(marketDeploymentTxHash) {
		return repository.AttachMarketContractInput{}, errors.New("market_deployment_tx_hash must be an EVM transaction hash")
	}
	if !isEVMAddressShape(marketFactoryAddress) {
		return repository.AttachMarketContractInput{}, errors.New("market_factory_address must be an EVM address")
	}
	if !isEVMAddressShape(resolverAddress) {
		return repository.AttachMarketContractInput{}, errors.New("resolver_address must be an EVM address")
	}

	return repository.AttachMarketContractInput{
		MarketContractAddress:  strings.ToLower(marketContractAddress),
		MarketDeploymentTxHash: strings.ToLower(marketDeploymentTxHash),
		MarketFactoryAddress:   strings.ToLower(marketFactoryAddress),
		ResolverAddress:        strings.ToLower(resolverAddress),
	}, nil
}

func registerMarketRoutes(router chi.Router, marketsRepository *repository.MarketsRepository) {
	router.Get("/markets", func(w http.ResponseWriter, r *http.Request) {
		markets, err := marketsRepository.ListMarkets(r.Context(), defaultMarketsLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "markets_list_failed", "failed to list markets")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"markets": newMarketResponses(markets),
		})
	})

	router.Get("/agent/markets", func(w http.ResponseWriter, r *http.Request) {
		markets, err := marketsRepository.ListMarkets(r.Context(), defaultMarketsLimit)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "markets_list_failed", "failed to list markets")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"markets": newAgentMarketResponses(markets),
		})
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

		httpjson.WriteJSON(w, http.StatusCreated, map[string]any{
			"market": newMarketResponse(market),
		})
	})

	router.Patch("/markets/{id}/contract", func(w http.ResponseWriter, r *http.Request) {
		marketID := chi.URLParam(r, "id")
		if !isUUIDShape(marketID) {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_market_id", "market id is invalid")
			return
		}

		var request attachMarketContractRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
			return
		}

		input, err := request.toRepositoryInput()
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_contract_attachment", "invalid contract attachment")
			return
		}

		market, err := marketsRepository.AttachMarketContract(r.Context(), marketID, input)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				httpjson.WriteError(w, http.StatusNotFound, "market_not_found", "market not found")
				return
			}
			if isUniqueViolation(err) {
				httpjson.WriteError(w, http.StatusConflict, "market_contract_already_attached", "market contract address is already attached")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "market_contract_attach_failed", "failed to attach market contract")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"market": newMarketResponse(market),
		})
	})

	router.Get("/markets/{id}", func(w http.ResponseWriter, r *http.Request) {
		marketID := chi.URLParam(r, "id")
		market, err := marketsRepository.GetMarketByID(r.Context(), marketID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				httpjson.WriteError(w, http.StatusNotFound, "market_not_found", "market not found")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "market_get_failed", "failed to get market")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"market": newMarketResponse(market),
		})
	})
}
