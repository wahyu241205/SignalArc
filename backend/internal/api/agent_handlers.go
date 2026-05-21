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
	Action     string `json:"action"`
	UserWallet string `json:"user_wallet"`
	MarketID   string `json:"market_id"`
	Amount     string `json:"amount"`
	Outcome    string `json:"outcome"`
}

type agentIntentResponse struct {
	IntentID             string                 `json:"intent_id"`
	Action               string                 `json:"action"`
	Status               string                 `json:"status"`
	RequiresConfirmation bool                   `json:"requires_confirmation"`
	UserWallet           string                 `json:"user_wallet,omitempty"`
	Address              string                 `json:"address,omitempty"`
	MarketID             string                 `json:"market_id,omitempty"`
	Amount               string                 `json:"amount,omitempty"`
	Outcome              string                 `json:"outcome,omitempty"`
	ValidationResult     agent.ValidationResult `json:"validation_result"`
	Warnings             []string               `json:"warnings"`
	CreatedAt            string                 `json:"created_at"`
}

func registerAgentIntentRoutes(router chi.Router, store *agent.Store) {
	router.Post("/agent/intents", func(w http.ResponseWriter, r *http.Request) {
		var request createAgentIntentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
			return
		}

		intent, err := store.CreateIntent(agent.CreateIntentInput{
			Action:     request.Action,
			UserWallet: request.UserWallet,
			MarketID:   request.MarketID,
			Amount:     request.Amount,
			Outcome:    request.Outcome,
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
}

func newAgentIntentResponse(intent agent.Intent) agentIntentResponse {
	return agentIntentResponse{
		IntentID:             intent.ID,
		Action:               intent.Action,
		Status:               intent.Status,
		RequiresConfirmation: intent.RequiresConfirmation,
		UserWallet:           intent.UserWallet,
		Address:              intent.UserWallet,
		MarketID:             intent.MarketID,
		Amount:               intent.Amount,
		Outcome:              intent.Outcome,
		ValidationResult:     intent.ValidationResult,
		Warnings:             intent.Warnings,
		CreatedAt:            intent.CreatedAt.Format("2006-01-02T15:04:05.000000000Z07:00"),
	}
}
