package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/agent"
	"github.com/wahyu241205/SignalArc/backend/internal/httpjson"
)

type createAgentIntentRequest struct {
	Action                string `json:"action"`
	UserWallet            string `json:"user_wallet"`
	MarketID              string `json:"market_id"`
	MarketContractAddress string `json:"market_contract_address"`
	Amount                string `json:"amount"`
	Outcome               string `json:"outcome"`
	Resolver              string `json:"resolver"`
	CollateralToken       string `json:"collateral_token"`
	CloseTimestamp        string `json:"close_timestamp"`
	Question              string `json:"question"`
}

type agentIntentResponse struct {
	IntentID              string                 `json:"intent_id"`
	Action                string                 `json:"action"`
	Status                string                 `json:"status"`
	RequiresConfirmation  bool                   `json:"requires_confirmation"`
	UserWallet            string                 `json:"user_wallet,omitempty"`
	Address               string                 `json:"address,omitempty"`
	MarketID              string                 `json:"market_id,omitempty"`
	MarketContractAddress string                 `json:"market_contract_address,omitempty"`
	Amount                string                 `json:"amount,omitempty"`
	Outcome               string                 `json:"outcome,omitempty"`
	Resolver              string                 `json:"resolver,omitempty"`
	CollateralToken       string                 `json:"collateral_token,omitempty"`
	CloseTimestamp        string                 `json:"close_timestamp,omitempty"`
	Question              string                 `json:"question,omitempty"`
	ValidationResult      agent.ValidationResult `json:"validation_result"`
	Warnings              []string               `json:"warnings"`
	CreatedAt             string                 `json:"created_at"`
}

type agentExecutionPlanResponse struct {
	IntentID            string                     `json:"intent_id"`
	Action              string                     `json:"action"`
	Status              string                     `json:"status"`
	ExecutionMode       string                     `json:"execution_mode"`
	Network             string                     `json:"network"`
	AgentFactoryAddress string                     `json:"agent_factory_address"`
	RequiresSignature   bool                       `json:"requires_signature"`
	BroadcastPerformed  bool                       `json:"broadcast_performed"`
	TransactionHash     *string                    `json:"transaction_hash"`
	TransactionRequest  transactionRequestResponse `json:"transaction_request"`
	Warnings            []string                   `json:"warnings"`
}

type agentExecutionResponse struct {
	IntentID               string                `json:"intent_id"`
	Action                 string                `json:"action"`
	Status                 string                `json:"status"`
	ExecutionMode          string                `json:"execution_mode"`
	Network                string                `json:"network"`
	AgentFactoryAddress    string                `json:"agent_factory_address"`
	MarketContractAddress  string                `json:"market_contract_address,omitempty"`
	BroadcastPerformed     bool                  `json:"broadcast_performed"`
	ApproveTransactionHash string                `json:"approve_transaction_hash,omitempty"`
	TransactionHash        string                `json:"transaction_hash"`
	Readback               agentReadbackResponse `json:"readback"`
}

type agentReadbackResponse struct {
	MarketCount     string `json:"market_count"`
	CreatedMarket   string `json:"created_market,omitempty"`
	IsMarket        *bool  `json:"is_market,omitempty"`
	YesPositions    string `json:"yes_positions,omitempty"`
	NoPositions     string `json:"no_positions,omitempty"`
	TotalYes        string `json:"total_yes,omitempty"`
	TotalNo         string `json:"total_no,omitempty"`
	TotalCollateral string `json:"total_collateral,omitempty"`
	USDCBalance     string `json:"usdc_balance,omitempty"`
	USDCAllowance   string `json:"usdc_allowance,omitempty"`
}

type transactionRequestResponse struct {
	To                 string   `json:"to"`
	Contract           string   `json:"contract"`
	Function           string   `json:"function"`
	Args               []string `json:"args"`
	Value              string   `json:"value"`
	Chain              string   `json:"chain"`
	BroadcastPerformed bool     `json:"broadcast_performed"`
}

func registerAgentIntentRoutes(router chi.Router, store *agent.Store, executor agent.Executor) {
	router.Post("/agent/intents", func(w http.ResponseWriter, r *http.Request) {
		var request createAgentIntentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
			return
		}

		intent, err := store.CreateIntent(agent.CreateIntentInput{
			Action:                request.Action,
			UserWallet:            request.UserWallet,
			MarketID:              request.MarketID,
			MarketContractAddress: request.MarketContractAddress,
			Amount:                request.Amount,
			Outcome:               request.Outcome,
			Resolver:              request.Resolver,
			CollateralToken:       request.CollateralToken,
			CloseTimestamp:        request.CloseTimestamp,
			Question:              request.Question,
		})
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "agent_intent_create_failed", "failed to create agent intent preview")
			return
		}

		if !intent.ValidationResult.Valid {
			httpjson.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"intent": newAgentIntentResponse(intent),
			})
			return
		}

		httpjson.WriteJSON(w, http.StatusCreated, map[string]any{
			"intent": newAgentIntentResponse(intent),
		})
	})

	router.Get("/agent/intents/{id}", func(w http.ResponseWriter, r *http.Request) {
		intentID := chi.URLParam(r, "id")
		intent, err := store.GetIntent(intentID)
		if err != nil {
			if errors.Is(err, agent.ErrIntentNotFound) {
				httpjson.WriteError(w, http.StatusNotFound, "agent_intent_not_found", "agent intent not found")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "agent_intent_get_failed", "failed to get agent intent")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"intent": newAgentIntentResponse(intent),
		})
	})

	router.Post("/agent/intents/{id}/confirm", func(w http.ResponseWriter, r *http.Request) {
		intentID := chi.URLParam(r, "id")
		executionPlan, err := store.ConfirmIntent(intentID)
		if err != nil {
			if errors.Is(err, agent.ErrIntentNotFound) {
				httpjson.WriteError(w, http.StatusNotFound, "agent_intent_not_found", "agent intent not found")
				return
			}
			if errors.Is(err, agent.ErrIntentInvalid) {
				httpjson.WriteError(w, http.StatusBadRequest, "agent_intent_invalid", "agent intent validation failed")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "agent_intent_confirm_failed", "failed to confirm agent intent")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"execution_plan": newAgentExecutionPlanResponse(executionPlan),
		})
	})

	router.Post("/agent/intents/{id}/execute", func(w http.ResponseWriter, r *http.Request) {
		intentID := chi.URLParam(r, "id")
		intent, err := store.GetIntent(intentID)
		if err != nil {
			if errors.Is(err, agent.ErrIntentNotFound) {
				httpjson.WriteError(w, http.StatusNotFound, "agent_intent_not_found", "agent intent not found")
				return
			}

			httpjson.WriteError(w, http.StatusInternalServerError, "agent_intent_get_failed", "failed to get agent intent")
			return
		}

		if intent.Status != agent.StatusConfirmed {
			httpjson.WriteError(w, http.StatusConflict, "agent_intent_not_confirmed", "agent intent must be confirmed before execution")
			return
		}
		if !intent.ValidationResult.Valid {
			httpjson.WriteError(w, http.StatusBadRequest, "agent_intent_invalid", "agent intent validation failed")
			return
		}
		if intent.Action != agent.ActionCreateMarket && intent.Action != agent.ActionBuyYes && intent.Action != agent.ActionBuyNo {
			httpjson.WriteError(w, http.StatusNotImplemented, "not_implemented", "only create_market, buy_yes, and buy_no execution are implemented")
			return
		}

		activeExecutor := executor
		if activeExecutor == nil {
			activeExecutor, err = agent.NewArcExecutorFromEnv()
			if err != nil {
				httpjson.WriteError(w, http.StatusServiceUnavailable, "agent_execution_config_invalid", "agent execution environment is not configured")
				return
			}
		}

		var result agent.ExecutionResult
		switch intent.Action {
		case agent.ActionCreateMarket:
			result, err = activeExecutor.ExecuteCreateMarket(r.Context(), intent)
		case agent.ActionBuyYes:
			result, err = activeExecutor.ExecuteBuyYes(r.Context(), intent)
		case agent.ActionBuyNo:
			result, err = activeExecutor.ExecuteBuyNo(r.Context(), intent)
		default:
			err = agent.ErrExecutionNotImplemented
		}
		if err != nil {
			if errors.Is(err, agent.ErrExecutionNotImplemented) {
				httpjson.WriteError(w, http.StatusNotImplemented, "not_implemented", "only create_market, buy_yes, and buy_no execution are implemented")
				return
			}
			if errors.Is(err, agent.ErrIntentInvalid) {
				httpjson.WriteError(w, http.StatusBadRequest, "agent_intent_invalid", "agent intent validation failed")
				return
			}
			if errors.Is(err, agent.ErrIntentNotConfirmed) {
				httpjson.WriteError(w, http.StatusConflict, "agent_intent_not_confirmed", "agent intent must be confirmed before execution")
				return
			}
			if errors.Is(err, agent.ErrExecutionConfigInvalid) {
				httpjson.WriteError(w, http.StatusServiceUnavailable, "agent_execution_config_invalid", "agent execution environment is not configured")
				return
			}

			httpjson.WriteError(w, http.StatusBadGateway, "agent_execution_failed", "agent execution failed")
			return
		}

		httpjson.WriteJSON(w, http.StatusOK, map[string]any{
			"execution": newAgentExecutionResponse(result),
		})
	})
}

func newAgentIntentResponse(intent agent.Intent) agentIntentResponse {
	return agentIntentResponse{
		IntentID:              intent.ID,
		Action:                intent.Action,
		Status:                intent.Status,
		RequiresConfirmation:  intent.RequiresConfirmation,
		UserWallet:            intent.UserWallet,
		Address:               intent.UserWallet,
		MarketID:              intent.MarketID,
		MarketContractAddress: intent.MarketContractAddress,
		Amount:                intent.Amount,
		Outcome:               intent.Outcome,
		Resolver:              intent.Resolver,
		CollateralToken:       intent.CollateralToken,
		CloseTimestamp:        intent.CloseTimestamp,
		Question:              intent.Question,
		ValidationResult:      intent.ValidationResult,
		Warnings:              intent.Warnings,
		CreatedAt:             intent.CreatedAt.Format("2006-01-02T15:04:05.000000000Z07:00"),
	}
}

func newAgentExecutionPlanResponse(executionPlan agent.ExecutionPlan) agentExecutionPlanResponse {
	return agentExecutionPlanResponse{
		IntentID:            executionPlan.IntentID,
		Action:              executionPlan.Action,
		Status:              executionPlan.Status,
		ExecutionMode:       executionPlan.ExecutionMode,
		Network:             executionPlan.Network,
		AgentFactoryAddress: executionPlan.AgentFactoryAddress,
		RequiresSignature:   executionPlan.RequiresSignature,
		BroadcastPerformed:  executionPlan.BroadcastPerformed,
		TransactionHash:     executionPlan.TransactionHash,
		TransactionRequest:  newTransactionRequestResponse(executionPlan.TransactionRequest),
		Warnings:            executionPlan.Warnings,
	}
}

func newTransactionRequestResponse(transactionRequest agent.TransactionRequest) transactionRequestResponse {
	return transactionRequestResponse{
		To:                 transactionRequest.To,
		Contract:           transactionRequest.Contract,
		Function:           transactionRequest.Function,
		Args:               transactionRequest.Args,
		Value:              transactionRequest.Value,
		Chain:              transactionRequest.Chain,
		BroadcastPerformed: transactionRequest.BroadcastPerformed,
	}
}

func newAgentExecutionResponse(result agent.ExecutionResult) agentExecutionResponse {
	return agentExecutionResponse{
		IntentID:               result.IntentID,
		Action:                 result.Action,
		Status:                 result.Status,
		ExecutionMode:          result.ExecutionMode,
		Network:                result.Network,
		AgentFactoryAddress:    result.AgentFactoryAddress,
		MarketContractAddress:  result.MarketContractAddress,
		BroadcastPerformed:     result.BroadcastPerformed,
		ApproveTransactionHash: result.ApproveTransactionHash,
		TransactionHash:        result.TransactionHash,
		Readback: agentReadbackResponse{
			MarketCount:     result.Readback.MarketCount,
			CreatedMarket:   result.Readback.CreatedMarket,
			IsMarket:        result.Readback.IsMarket,
			YesPositions:    result.Readback.YesPositions,
			NoPositions:     result.Readback.NoPositions,
			TotalYes:        result.Readback.TotalYes,
			TotalNo:         result.Readback.TotalNo,
			TotalCollateral: result.Readback.TotalCollateral,
			USDCBalance:     result.Readback.USDCBalance,
			USDCAllowance:   result.Readback.USDCAllowance,
		},
	}
}
