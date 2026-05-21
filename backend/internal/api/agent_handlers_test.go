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
		"market_contract_address": "0x3333333333333333333333333333333333333333",
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
		"market_id": "market-2",
		"market_contract_address": "0x4444444444444444444444444444444444444444"
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
		"market_contract_address": "0x5555555555555555555555555555555555555555",
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
	if body.ExecutionPlan.TransactionRequest.BroadcastPerformed {
		t.Fatal("expected transaction_request broadcast_performed false")
	}
}

func TestConfirmCreateMarketReturnsFactoryTransactionRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-1",
		"question": "Will SignalArc create an agent market?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)

	executionPlan := confirmAgentIntent(t, router, intentID)

	if executionPlan.TransactionRequest.To != "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80" {
		t.Fatalf("expected factory address as to, got %q", executionPlan.TransactionRequest.To)
	}
	if executionPlan.TransactionRequest.Contract != "SignalArcAgentMarketFactory" {
		t.Fatalf("expected factory contract, got %q", executionPlan.TransactionRequest.Contract)
	}
	if executionPlan.TransactionRequest.Function != "createMarket" {
		t.Fatalf("expected createMarket function, got %q", executionPlan.TransactionRequest.Function)
	}
	expectedArgs := []string{
		"agent-market-1",
		"Will SignalArc create an agent market?",
		"1767225600",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	}
	assertStringSliceEqual(t, executionPlan.TransactionRequest.Args, expectedArgs)
	assertNoExecutionClaim(t, executionPlan)
}

func TestConfirmBuyYesReturnsMarketTransactionRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())
	marketAddress := "0x4444444444444444444444444444444444444444"

	intentID := createAgentIntent(t, router, `{
		"action": "buy_yes",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "`+marketAddress+`",
		"amount": "42.5"
	}`)

	executionPlan := confirmAgentIntent(t, router, intentID)

	if executionPlan.TransactionRequest.To != marketAddress {
		t.Fatalf("expected market address as to, got %q", executionPlan.TransactionRequest.To)
	}
	if executionPlan.TransactionRequest.Contract != "SignalArcAgentMarket" {
		t.Fatalf("expected agent market contract, got %q", executionPlan.TransactionRequest.Contract)
	}
	if executionPlan.TransactionRequest.Function != "buyYes" {
		t.Fatalf("expected buyYes function, got %q", executionPlan.TransactionRequest.Function)
	}
	assertStringSliceEqual(t, executionPlan.TransactionRequest.Args, []string{"42.5"})
	assertNoExecutionClaim(t, executionPlan)
}

func TestConfirmClaimRefundReturnsMarketTransactionRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore())
	marketAddress := "0x5555555555555555555555555555555555555555"

	intentID := createAgentIntent(t, router, `{
		"action": "claim_refund",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "`+marketAddress+`"
	}`)

	executionPlan := confirmAgentIntent(t, router, intentID)

	if executionPlan.TransactionRequest.To != marketAddress {
		t.Fatalf("expected market address as to, got %q", executionPlan.TransactionRequest.To)
	}
	if executionPlan.TransactionRequest.Contract != "SignalArcAgentMarket" {
		t.Fatalf("expected agent market contract, got %q", executionPlan.TransactionRequest.Contract)
	}
	if executionPlan.TransactionRequest.Function != "claimRefund" {
		t.Fatalf("expected claimRefund function, got %q", executionPlan.TransactionRequest.Function)
	}
	assertStringSliceEqual(t, executionPlan.TransactionRequest.Args, []string{})
	assertNoExecutionClaim(t, executionPlan)
}

func TestBuyYesMissingMarketContractAddressFailsValidation(t *testing.T) {
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
	if !containsString(body.Intent.ValidationResult.Errors, "market_contract_address is required for existing market contract actions") {
		t.Fatalf("expected market_contract_address validation error, got %#v", body.Intent.ValidationResult.Errors)
	}
}

func createValidAgentIntent(t *testing.T, router http.Handler) string {
	t.Helper()

	return createAgentIntent(t, router, `{
		"action": "buy_yes",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "0x3333333333333333333333333333333333333333",
		"amount": "12.5"
	}`)
}

func createAgentIntent(t *testing.T, router http.Handler, payload string) string {
	t.Helper()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(payload))
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

func confirmAgentIntent(t *testing.T, router http.Handler, intentID string) agentExecutionPlanResponse {
	t.Helper()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/confirm", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected confirm status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body struct {
		ExecutionPlan agentExecutionPlanResponse `json:"execution_plan"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode confirm response: %v", err)
	}

	return body.ExecutionPlan
}

func assertNoExecutionClaim(t *testing.T, executionPlan agentExecutionPlanResponse) {
	t.Helper()

	if executionPlan.BroadcastPerformed {
		t.Fatal("expected broadcast_performed false")
	}
	if executionPlan.TransactionHash != nil {
		t.Fatalf("expected transaction_hash null, got %q", *executionPlan.TransactionHash)
	}
	if executionPlan.TransactionRequest.BroadcastPerformed {
		t.Fatal("expected transaction_request broadcast_performed false")
	}
	if executionPlan.TransactionRequest.Value != "0" {
		t.Fatalf("expected transaction value 0, got %q", executionPlan.TransactionRequest.Value)
	}
	if executionPlan.TransactionRequest.Chain != "arc_testnet" {
		t.Fatalf("expected transaction chain arc_testnet, got %q", executionPlan.TransactionRequest.Chain)
	}
}

func assertStringSliceEqual(t *testing.T, actual []string, expected []string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected args %#v, got %#v", expected, actual)
	}
	for index := range expected {
		if actual[index] != expected[index] {
			t.Fatalf("expected args %#v, got %#v", expected, actual)
		}
	}
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}

	return false
}
