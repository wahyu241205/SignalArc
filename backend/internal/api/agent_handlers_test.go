package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/agent"
)

func TestCreateAgentIntentPreview(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "buy_yes",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"amount": "12.5"
	}`))

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	var body struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Intent.IntentID == "" {
		t.Fatal("expected intent_id")
	}
	if body.Intent.Action != "buy_yes" {
		t.Fatalf("expected action buy_yes, got %q", body.Intent.Action)
	}
	if body.Intent.Status != "preview" {
		t.Fatalf("expected preview status, got %q", body.Intent.Status)
	}
	if !body.Intent.RequiresConfirmation {
		t.Fatal("expected requires_confirmation true")
	}
	if body.Intent.UserWallet != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("expected user wallet to be echoed, got %q", body.Intent.UserWallet)
	}
	if !body.Intent.ValidationResult.Valid {
		t.Fatalf("expected valid validation result, got %#v", body.Intent.ValidationResult)
	}
	if len(body.Intent.Warnings) == 0 {
		t.Fatal("expected preview warnings")
	}
}

func TestGetAgentIntentPreview(t *testing.T) {
	store := agent.NewStore()
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store)

	createResponse := httptest.NewRecorder()
	createRequest := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "claim_payout",
		"user_wallet": "0x2222222222222222222222222222222222222222",
		"market_id": "market-2"
	}`))
	router.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createResponse.Code)
	}

	var createBody struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(createResponse.Body).Decode(&createBody); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	getResponse := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodGet, "/agent/intents/"+createBody.Intent.IntentID, nil)
	router.ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, getResponse.Code)
	}

	var getBody struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(getResponse.Body).Decode(&getBody); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if getBody.Intent.IntentID != createBody.Intent.IntentID {
		t.Fatalf("expected intent id %q, got %q", createBody.Intent.IntentID, getBody.Intent.IntentID)
	}
}

func TestCreateAgentIntentValidationErrors(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "buy_no",
		"market_id": "market-1",
		"amount": "0"
	}`))

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	var body struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Intent.ValidationResult.Valid {
		t.Fatal("expected invalid validation result")
	}
	if len(body.Intent.ValidationResult.Errors) != 2 {
		t.Fatalf("expected two validation errors, got %#v", body.Intent.ValidationResult.Errors)
	}
}

func TestGetAgentIntentNotFound(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/agent/intents/missing", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestConfirmValidAgentIntent(t *testing.T) {
	store := agent.NewStore()
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store)
	intentID := createValidAgentIntent(t, router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/confirm", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body struct {
		ExecutionPlan agentExecutionPlanResponse `json:"execution_plan"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.ExecutionPlan.IntentID != intentID {
		t.Fatalf("expected intent id %q, got %q", intentID, body.ExecutionPlan.IntentID)
	}
	if body.ExecutionPlan.Status != "confirmed" {
		t.Fatalf("expected confirmed status, got %q", body.ExecutionPlan.Status)
	}
	if body.ExecutionPlan.ExecutionMode != "agent_contract" {
		t.Fatalf("expected agent_contract execution mode, got %q", body.ExecutionPlan.ExecutionMode)
	}
	if body.ExecutionPlan.Network != "arc_testnet" {
		t.Fatalf("expected arc_testnet network, got %q", body.ExecutionPlan.Network)
	}
	if body.ExecutionPlan.AgentFactoryAddress != "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80" {
		t.Fatalf("unexpected agent factory address %q", body.ExecutionPlan.AgentFactoryAddress)
	}
	if !body.ExecutionPlan.RequiresSignature {
		t.Fatal("expected requires_signature true")
	}
	if body.ExecutionPlan.BroadcastPerformed {
		t.Fatal("expected broadcast_performed false")
	}
	if body.ExecutionPlan.TransactionHash != nil {
		t.Fatalf("expected nil transaction hash, got %q", *body.ExecutionPlan.TransactionHash)
	}
	if len(body.ExecutionPlan.Warnings) == 0 {
		t.Fatal("expected warnings")
	}
}

func TestConfirmMissingAgentIntentReturnsNotFound(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/missing/confirm", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestConfirmInvalidAgentIntentReturnsBadRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	createResponse := httptest.NewRecorder()
	createRequest := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "buy_yes",
		"market_id": "market-1",
		"amount": "0"
	}`))
	router.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusBadRequest {
		t.Fatalf("expected create status %d, got %d", http.StatusBadRequest, createResponse.Code)
	}

	var createBody struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(createResponse.Body).Decode(&createBody); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+createBody.Intent.IntentID+"/confirm", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected confirm status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestConfirmAgentIntentIsIdempotent(t *testing.T) {
	store := agent.NewStore()
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store)
	intentID := createValidAgentIntent(t, router)

	firstResponse := httptest.NewRecorder()
	firstRequest := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/confirm", nil)
	router.ServeHTTP(firstResponse, firstRequest)

	secondResponse := httptest.NewRecorder()
	secondRequest := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/confirm", nil)
	router.ServeHTTP(secondResponse, secondRequest)

	if firstResponse.Code != http.StatusOK {
		t.Fatalf("expected first status %d, got %d", http.StatusOK, firstResponse.Code)
	}
	if secondResponse.Code != http.StatusOK {
		t.Fatalf("expected second status %d, got %d", http.StatusOK, secondResponse.Code)
	}

	var firstBody struct {
		ExecutionPlan agentExecutionPlanResponse `json:"execution_plan"`
	}
	if err := json.NewDecoder(firstResponse.Body).Decode(&firstBody); err != nil {
		t.Fatalf("decode first response: %v", err)
	}

	var secondBody struct {
		ExecutionPlan agentExecutionPlanResponse `json:"execution_plan"`
	}
	if err := json.NewDecoder(secondResponse.Body).Decode(&secondBody); err != nil {
		t.Fatalf("decode second response: %v", err)
	}

	if firstBody.ExecutionPlan.IntentID != secondBody.ExecutionPlan.IntentID {
		t.Fatalf("expected same intent id, got %q and %q", firstBody.ExecutionPlan.IntentID, secondBody.ExecutionPlan.IntentID)
	}
	if secondBody.ExecutionPlan.Status != "confirmed" {
		t.Fatalf("expected confirmed status, got %q", secondBody.ExecutionPlan.Status)
	}
}

func TestConfirmAgentIntentResponseSaysBroadcastNotPerformed(t *testing.T) {
	store := agent.NewStore()
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store)
	intentID := createValidAgentIntent(t, router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/confirm", nil)
	router.ServeHTTP(response, request)

	var body struct {
		ExecutionPlan agentExecutionPlanResponse `json:"execution_plan"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.ExecutionPlan.BroadcastPerformed {
		t.Fatal("expected broadcast_performed false")
	}
	if body.ExecutionPlan.TransactionHash != nil {
		t.Fatalf("expected transaction_hash null, got %q", *body.ExecutionPlan.TransactionHash)
	}
}

func createValidAgentIntent(t *testing.T, router http.Handler) string {
	t.Helper()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "buy_yes",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"amount": "12.5"
	}`))
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	var body struct {
		Intent agentIntentResponse `json:"intent"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	return body.Intent.IntentID
}
