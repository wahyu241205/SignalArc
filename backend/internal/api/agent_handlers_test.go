package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/agent"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

type stubAgentExecutor struct {
	result agent.ExecutionResult
	err    error
	intent agent.Intent
	called bool
}

func (executor *stubAgentExecutor) ExecuteCreateMarket(_ context.Context, intent agent.Intent) (agent.ExecutionResult, error) {
	executor.called = true
	executor.intent = intent
	if executor.err != nil {
		return agent.ExecutionResult{}, executor.err
	}
	return executor.result, nil
}

func (executor *stubAgentExecutor) ExecuteBuyYes(_ context.Context, intent agent.Intent) (agent.ExecutionResult, error) {
	executor.called = true
	executor.intent = intent
	if executor.err != nil {
		return agent.ExecutionResult{}, executor.err
	}
	return executor.result, nil
}

func (executor *stubAgentExecutor) ExecuteBuyNo(_ context.Context, intent agent.Intent) (agent.ExecutionResult, error) {
	executor.called = true
	executor.intent = intent
	if executor.err != nil {
		return agent.ExecutionResult{}, executor.err
	}
	return executor.result, nil
}

type testAgentWalletRegistry struct {
	wallets map[string]repository.AgentWallet
}

func newTestAgentWalletRegistry() *testAgentWalletRegistry {
	return &testAgentWalletRegistry{wallets: map[string]repository.AgentWallet{}}
}

func (registry *testAgentWalletRegistry) RegisterAgentWallet(_ context.Context, input repository.UpsertAgentWalletInput) (repository.AgentWallet, error) {
	now := time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC)
	wallet := repository.AgentWallet{
		ID:                 "agent_wallet_test_1",
		AgentID:            input.AgentID,
		UserWallet:         input.UserWallet,
		UserEmail:          input.UserEmail,
		AgentWalletAddress: input.AgentWalletAddress,
		WalletProvider:     input.WalletProvider,
		Chain:              input.Chain,
		Status:             input.Status,
		AllowedActions:     input.AllowedActions,
		PolicyMetadata:     input.PolicyMetadata,
		SourceClient:       input.SourceClient,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	registry.wallets[wallet.AgentID] = wallet
	return wallet, nil
}

func (registry *testAgentWalletRegistry) GetAgentWalletByAgentID(_ context.Context, agentID string) (repository.AgentWallet, error) {
	wallet, ok := registry.wallets[agentID]
	if !ok {
		return repository.AgentWallet{}, sql.ErrNoRows
	}
	return wallet, nil
}

func (registry *testAgentWalletRegistry) DisableAgentWallet(_ context.Context, agentID string) (repository.AgentWallet, error) {
	wallet, ok := registry.wallets[agentID]
	if !ok {
		return repository.AgentWallet{}, sql.ErrNoRows
	}
	wallet.Status = "disabled"
	wallet.UpdatedAt = wallet.UpdatedAt.Add(time.Minute)
	registry.wallets[agentID] = wallet
	return wallet, nil
}

func TestCreateAgentIntentPreview(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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

func TestRegisterAgentWallet(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/wallets", bytes.NewBufferString(`{
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"agent_wallet_address": "0x9999999999999999999999999999999999999999",
		"wallet_provider": "circle_agent_wallet",
		"chain": "ARC-TESTNET",
		"allowed_actions": ["create_market", "buy_yes"],
		"status": "active",
		"policy_metadata": {
			"per_tx_usdc_cap": "unknown / not documented until Circle policy is configured"
		}
	}`))

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	var body struct {
		AgentWallet agentWalletResponse `json:"agent_wallet"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AgentWallet.AgentID != "agent_test_1" {
		t.Fatalf("expected agent id, got %q", body.AgentWallet.AgentID)
	}
	if body.AgentWallet.AgentWalletAddress != "0x9999999999999999999999999999999999999999" {
		t.Fatalf("unexpected agent wallet address %q", body.AgentWallet.AgentWalletAddress)
	}
	if body.AgentWallet.WalletProvider != agent.WalletProviderCircleAgentWallet {
		t.Fatalf("unexpected wallet provider %q", body.AgentWallet.WalletProvider)
	}
}

func TestGetAgentWalletByAgentID(t *testing.T) {
	registry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, registry, agent.ActionCreateMarket, agent.ActionBuyYes)
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), registry, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/agent/wallets/agent_test_1", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body struct {
		AgentWallet agentWalletResponse `json:"agent_wallet"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AgentWallet.AgentID != "agent_test_1" {
		t.Fatalf("expected agent id, got %q", body.AgentWallet.AgentID)
	}
	if !containsString(body.AgentWallet.AllowedActions, agent.ActionBuyYes) {
		t.Fatalf("expected buy_yes in allowed actions, got %#v", body.AgentWallet.AllowedActions)
	}
}

func TestDisableAgentWallet(t *testing.T) {
	registry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, registry, agent.ActionBuyYes)
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), registry, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/wallets/agent_test_1/disable", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body struct {
		AgentWallet agentWalletResponse `json:"agent_wallet"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AgentWallet.Status != "disabled" {
		t.Fatalf("expected disabled status, got %q", body.AgentWallet.Status)
	}
}

func TestRegisterAgentWalletRejectsDeployerWallet(t *testing.T) {
	assertAgentWalletRegistrationFails(t, `{
		"agent_id": "agent_bad",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"agent_wallet_address": "0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC",
		"wallet_provider": "circle_agent_wallet",
		"chain": "ARC-TESTNET",
		"allowed_actions": ["buy_yes"]
	}`)
}

func TestRegisterAgentWalletRejectsUserWalletReuse(t *testing.T) {
	assertAgentWalletRegistrationFails(t, `{
		"agent_id": "agent_bad",
		"user_wallet": "0x9999999999999999999999999999999999999999",
		"agent_wallet_address": "0x9999999999999999999999999999999999999999",
		"wallet_provider": "circle_agent_wallet",
		"chain": "ARC-TESTNET",
		"allowed_actions": ["buy_yes"]
	}`)
}

func TestRegisterAgentWalletRejectsUnsupportedProvider(t *testing.T) {
	assertAgentWalletRegistrationFails(t, `{
		"agent_id": "agent_bad",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"agent_wallet_address": "0x9999999999999999999999999999999999999999",
		"wallet_provider": "temporary_testnet_agent_eoa",
		"chain": "ARC-TESTNET",
		"allowed_actions": ["buy_yes"]
	}`)
}

func TestRegisterAgentWalletRejectsWrongChain(t *testing.T) {
	assertAgentWalletRegistrationFails(t, `{
		"agent_id": "agent_bad",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"agent_wallet_address": "0x9999999999999999999999999999999999999999",
		"wallet_provider": "circle_agent_wallet",
		"chain": "BASE",
		"allowed_actions": ["buy_yes"]
	}`)
}

func TestGetAgentIntentPreview(t *testing.T) {
	store := agent.NewStore()
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, newTestAgentWalletRegistry(), nil)

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
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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

func TestCreateAgentIntentIncludesRegisteredAgentMetadata(t *testing.T) {
	registry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, registry, agent.ActionBuyYes)
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), registry, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents", bytes.NewBufferString(`{
		"action": "buy_yes",
		"agent_id": "agent_test_1",
		"source_client": "chatgpt_custom_action",
		"client_request_id": "client_req_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "0x3333333333333333333333333333333333333333",
		"amount": "1000000"
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
	if body.Intent.AgentID != "agent_test_1" {
		t.Fatalf("expected agent metadata, got %q", body.Intent.AgentID)
	}
	if body.Intent.AgentWalletAddress != "0x9999999999999999999999999999999999999999" {
		t.Fatalf("unexpected agent wallet address %q", body.Intent.AgentWalletAddress)
	}
	if body.Intent.WalletProvider != agent.WalletProviderCircleAgentWallet {
		t.Fatalf("unexpected wallet provider %q", body.Intent.WalletProvider)
	}
	if body.Intent.SourceClient != "chatgpt_custom_action" {
		t.Fatalf("unexpected source client %q", body.Intent.SourceClient)
	}
	if body.Intent.ClientRequestID != "client_req_1" {
		t.Fatalf("unexpected client request id %q", body.Intent.ClientRequestID)
	}
}

func TestGetAgentIntentNotFound(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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
	registerAgentIntentRoutes(router, store, newTestAgentWalletRegistry(), nil)
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
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/missing/confirm", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestConfirmInvalidAgentIntentReturnsBadRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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
	registerAgentIntentRoutes(router, store, newTestAgentWalletRegistry(), nil)
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
	registerAgentIntentRoutes(router, store, newTestAgentWalletRegistry(), nil)
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
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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

func TestExecuteConfirmedCreateMarketReturnsRealExecutionShape(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCreateMarket)
	isMarket := true
	executor := &stubAgentExecutor{
		result: agent.ExecutionResult{
			IntentID:            "set-by-test",
			AgentID:             "agent_test_1",
			AgentWalletAddress:  "0x9999999999999999999999999999999999999999",
			WalletProvider:      agent.WalletProviderCircleAgentWallet,
			Action:              agent.ActionCreateMarket,
			Status:              agent.StatusExecuted,
			ExecutionMode:       agent.ExecutionModeAgentContract,
			Network:             agent.NetworkArcTestnet,
			AgentFactoryAddress: "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80",
			BroadcastPerformed:  true,
			TransactionHash:     "0x1111111111111111111111111111111111111111111111111111111111111111",
			Readback: agent.ExecutionReadback{
				MarketCount:   "5",
				CreatedMarket: "0x2222222222222222222222222222222222222222",
				IsMarket:      &isMarket,
			},
		},
	}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-execute-1",
		"question": "Will SignalArc execute an agent market?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)
	confirmAgentIntent(t, router, intentID)
	executor.result.IntentID = intentID

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if !executor.called {
		t.Fatal("expected executor to be called")
	}
	if executor.intent.ID != intentID {
		t.Fatalf("expected executor intent %q, got %q", intentID, executor.intent.ID)
	}

	var body struct {
		Execution agentExecutionResponse `json:"execution"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode execute response: %v", err)
	}

	if body.Execution.Status != agent.StatusExecuted {
		t.Fatalf("expected executed status, got %q", body.Execution.Status)
	}
	if !body.Execution.BroadcastPerformed {
		t.Fatal("expected broadcast_performed true")
	}
	if body.Execution.TransactionHash != executor.result.TransactionHash {
		t.Fatalf("unexpected transaction hash %q", body.Execution.TransactionHash)
	}
	if body.Execution.Network != agent.NetworkArcTestnet {
		t.Fatalf("expected arc_testnet network, got %q", body.Execution.Network)
	}
	if body.Execution.Readback.MarketCount != "5" {
		t.Fatalf("expected market count 5, got %q", body.Execution.Readback.MarketCount)
	}
	if body.Execution.Readback.CreatedMarket != "0x2222222222222222222222222222222222222222" {
		t.Fatalf("unexpected created market %q", body.Execution.Readback.CreatedMarket)
	}
	if body.Execution.Readback.IsMarket == nil || !*body.Execution.Readback.IsMarket {
		t.Fatalf("expected is_market true, got %#v", body.Execution.Readback.IsMarket)
	}
}

func TestExecuteUnconfirmedIntentReturnsConflict(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCreateMarket)
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-execute-2",
		"question": "Will SignalArc reject unconfirmed execution?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusConflict, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called for unconfirmed intent")
	}
}

func TestExecuteMissingAgentWalletReturnsBadRequest(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "buy_yes",
		"agent_id": "agent_missing",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "0x3333333333333333333333333333333333333333",
		"amount": "1000000"
	}`)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusBadRequest, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called when the agent wallet is missing")
	}
}

func TestExecuteRejectsDeployerWalletAsAgentWallet(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	_, err := walletRegistry.RegisterAgentWallet(context.Background(), repository.UpsertAgentWalletInput{
		AgentID:            "agent_test_1",
		UserWallet:         "0x1111111111111111111111111111111111111111",
		AgentWalletAddress: knownDeployerResolverWallet(),
		WalletProvider:     agent.WalletProviderCircleAgentWallet,
		Chain:              agent.ChainArcTestnet,
		AllowedActions:     []string{agent.ActionBuyYes},
		Status:             agent.WalletStatusActive,
	})
	if err != nil {
		t.Fatalf("register agent wallet: %v", err)
	}
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createValidAgentIntent(t, router)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusForbidden, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called for deployer wallet")
	}
}

func TestExecuteRejectsInactiveAgentWallet(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionBuyYes)
	wallet, err := walletRegistry.DisableAgentWallet(context.Background(), "agent_test_1")
	if err != nil {
		t.Fatalf("disable test wallet: %v", err)
	}
	if wallet.Status != "disabled" {
		t.Fatalf("expected disabled wallet, got %q", wallet.Status)
	}
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createValidAgentIntent(t, router)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusForbidden, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called for inactive wallet")
	}
}

func TestExecuteRejectsDisallowedAction(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionBuyNo)
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createValidAgentIntent(t, router)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusForbidden, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called for disallowed action")
	}
}

func TestExecuteCircleProviderDisabledReturnsServiceUnavailable(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionBuyYes)
	executor := agent.NewCircleCLIExecutor(agent.CircleCLIExecutorConfig{
		Enabled: false,
	})
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createValidAgentIntent(t, router)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusServiceUnavailable, response.Code, response.Body.String())
	}
	if !bytes.Contains(response.Body.Bytes(), []byte("agent_execution_provider_disabled")) {
		t.Fatalf("expected provider disabled code, got %s", response.Body.String())
	}
}

func TestExecuteConfirmedBuyYesReturnsRealExecutionShape(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionBuyYes)
	executor := &stubAgentExecutor{
		result: agent.ExecutionResult{
			IntentID:               "set-by-test",
			AgentID:                "agent_test_1",
			AgentWalletAddress:     "0x9999999999999999999999999999999999999999",
			WalletProvider:         agent.WalletProviderCircleAgentWallet,
			Action:                 agent.ActionBuyYes,
			Status:                 agent.StatusExecuted,
			ExecutionMode:          agent.ExecutionModeAgentContract,
			Network:                agent.NetworkArcTestnet,
			MarketContractAddress:  "0x3333333333333333333333333333333333333333",
			BroadcastPerformed:     true,
			ApproveTransactionHash: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			TransactionHash:        "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			Readback: agent.ExecutionReadback{
				YesPositions:    "1000000",
				TotalYes:        "1000000",
				TotalCollateral: "1000000",
				USDCBalance:     "1000000",
			},
		},
	}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createValidAgentIntent(t, router)
	confirmAgentIntent(t, router, intentID)
	executor.result.IntentID = intentID

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if !executor.called {
		t.Fatal("expected executor to be called")
	}
	if executor.intent.Action != agent.ActionBuyYes {
		t.Fatalf("expected buy_yes intent, got %q", executor.intent.Action)
	}

	var body struct {
		Execution agentExecutionResponse `json:"execution"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode execute response: %v", err)
	}

	if body.Execution.Action != agent.ActionBuyYes {
		t.Fatalf("expected buy_yes action, got %q", body.Execution.Action)
	}
	if body.Execution.Status != agent.StatusExecuted {
		t.Fatalf("expected executed status, got %q", body.Execution.Status)
	}
	if !body.Execution.BroadcastPerformed {
		t.Fatal("expected broadcast_performed true")
	}
	if body.Execution.ApproveTransactionHash != executor.result.ApproveTransactionHash {
		t.Fatalf("unexpected approve transaction hash %q", body.Execution.ApproveTransactionHash)
	}
	if body.Execution.TransactionHash != executor.result.TransactionHash {
		t.Fatalf("unexpected buyYes transaction hash %q", body.Execution.TransactionHash)
	}
	if body.Execution.MarketContractAddress != executor.result.MarketContractAddress {
		t.Fatalf("unexpected market address %q", body.Execution.MarketContractAddress)
	}
	if body.Execution.Readback.YesPositions != "1000000" {
		t.Fatalf("expected yes positions 1000000, got %q", body.Execution.Readback.YesPositions)
	}
	if body.Execution.Readback.TotalYes != "1000000" {
		t.Fatalf("expected total yes 1000000, got %q", body.Execution.Readback.TotalYes)
	}
	if body.Execution.Readback.TotalCollateral != "1000000" {
		t.Fatalf("expected total collateral 1000000, got %q", body.Execution.Readback.TotalCollateral)
	}
	if body.Execution.Readback.USDCBalance != "1000000" {
		t.Fatalf("expected usdc balance 1000000, got %q", body.Execution.Readback.USDCBalance)
	}
}

func TestExecuteConfirmedBuyNoReturnsRealExecutionShape(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionBuyNo)
	executor := &stubAgentExecutor{
		result: agent.ExecutionResult{
			IntentID:               "set-by-test",
			AgentID:                "agent_test_1",
			AgentWalletAddress:     "0x9999999999999999999999999999999999999999",
			WalletProvider:         agent.WalletProviderCircleAgentWallet,
			Action:                 agent.ActionBuyNo,
			Status:                 agent.StatusExecuted,
			ExecutionMode:          agent.ExecutionModeAgentContract,
			Network:                agent.NetworkArcTestnet,
			MarketContractAddress:  "0x3333333333333333333333333333333333333333",
			BroadcastPerformed:     true,
			ApproveTransactionHash: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			TransactionHash:        "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
			Readback: agent.ExecutionReadback{
				NoPositions:     "1000000",
				TotalNo:         "1000000",
				TotalCollateral: "1000000",
				USDCBalance:     "1000000",
				USDCAllowance:   "0",
			},
		},
	}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "buy_no",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "0x3333333333333333333333333333333333333333",
		"amount": "1000000"
	}`)
	confirmAgentIntent(t, router, intentID)
	executor.result.IntentID = intentID

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if !executor.called {
		t.Fatal("expected executor to be called")
	}
	if executor.intent.Action != agent.ActionBuyNo {
		t.Fatalf("expected buy_no intent, got %q", executor.intent.Action)
	}

	var body struct {
		Execution agentExecutionResponse `json:"execution"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode execute response: %v", err)
	}

	if body.Execution.Action != agent.ActionBuyNo {
		t.Fatalf("expected buy_no action, got %q", body.Execution.Action)
	}
	if body.Execution.Status != agent.StatusExecuted {
		t.Fatalf("expected executed status, got %q", body.Execution.Status)
	}
	if !body.Execution.BroadcastPerformed {
		t.Fatal("expected broadcast_performed true")
	}
	if body.Execution.ApproveTransactionHash != executor.result.ApproveTransactionHash {
		t.Fatalf("unexpected approve transaction hash %q", body.Execution.ApproveTransactionHash)
	}
	if body.Execution.TransactionHash != executor.result.TransactionHash {
		t.Fatalf("unexpected buyNo transaction hash %q", body.Execution.TransactionHash)
	}
	if body.Execution.MarketContractAddress != executor.result.MarketContractAddress {
		t.Fatalf("unexpected market address %q", body.Execution.MarketContractAddress)
	}
	if body.Execution.Readback.NoPositions != "1000000" {
		t.Fatalf("expected no positions 1000000, got %q", body.Execution.Readback.NoPositions)
	}
	if body.Execution.Readback.TotalNo != "1000000" {
		t.Fatalf("expected total no 1000000, got %q", body.Execution.Readback.TotalNo)
	}
	if body.Execution.Readback.TotalCollateral != "1000000" {
		t.Fatalf("expected total collateral 1000000, got %q", body.Execution.Readback.TotalCollateral)
	}
	if body.Execution.Readback.USDCBalance != "1000000" {
		t.Fatalf("expected usdc balance 1000000, got %q", body.Execution.Readback.USDCBalance)
	}
	if body.Execution.Readback.USDCAllowance != "0" {
		t.Fatalf("expected usdc allowance 0, got %q", body.Execution.Readback.USDCAllowance)
	}
}

func TestExecuteUnsupportedActionReturnsNotImplemented(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCancelMarket)
	executor := &stubAgentExecutor{}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "cancel_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "market-1",
		"market_contract_address": "0x3333333333333333333333333333333333333333"
	}`)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotImplemented {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusNotImplemented, response.Code, response.Body.String())
	}
	if executor.called {
		t.Fatal("executor should not be called for unsupported action")
	}
}

func TestExecuteReportsConfigErrorWithoutSecretDetails(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCreateMarket)
	executor := &stubAgentExecutor{err: agent.ErrExecutionConfigInvalid}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-execute-3",
		"question": "Will SignalArc hide execution config details?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusServiceUnavailable, response.Code, response.Body.String())
	}
	if !executor.called {
		t.Fatal("expected executor to be called")
	}
	if bytes.Contains(response.Body.Bytes(), []byte("private")) {
		t.Fatalf("response should not expose private-key details: %s", response.Body.String())
	}
}

func TestExecuteMapsExecutorNotImplemented(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCreateMarket)
	executor := &stubAgentExecutor{err: agent.ErrExecutionNotImplemented}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-execute-4",
		"question": "Will SignalArc map executor not implemented?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotImplemented {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusNotImplemented, response.Code, response.Body.String())
	}
}

func TestExecuteMapsUnexpectedExecutorError(t *testing.T) {
	store := agent.NewStore()
	walletRegistry := newTestAgentWalletRegistry()
	registerTestAgentWallet(t, walletRegistry, agent.ActionCreateMarket)
	executor := &stubAgentExecutor{err: errors.New("rpc unavailable")}
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, store, walletRegistry, executor)

	intentID := createAgentIntent(t, router, `{
		"action": "create_market",
		"agent_id": "agent_test_1",
		"user_wallet": "0x1111111111111111111111111111111111111111",
		"market_id": "agent-market-execute-5",
		"question": "Will SignalArc map executor errors?",
		"close_timestamp": "1767225600",
		"resolver": "0x2222222222222222222222222222222222222222",
		"collateral_token": "0x3333333333333333333333333333333333333333"
	}`)
	confirmAgentIntent(t, router, intentID)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/intents/"+intentID+"/execute", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadGateway {
		t.Fatalf("expected execute status %d, got %d: %s", http.StatusBadGateway, response.Code, response.Body.String())
	}
}

func TestConfirmBuyYesReturnsMarketTransactionRequest(t *testing.T) {
	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)
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
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)
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
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

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
		"agent_id": "agent_test_1",
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

func assertAgentWalletRegistrationFails(t *testing.T, payload string) {
	t.Helper()

	router := chi.NewRouter()
	registerAgentIntentRoutes(router, agent.NewStore(), newTestAgentWalletRegistry(), nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/agent/wallets", bytes.NewBufferString(payload))
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, response.Code, response.Body.String())
	}
}

func registerTestAgentWallet(t *testing.T, registry *testAgentWalletRegistry, allowedActions ...string) {
	t.Helper()

	_, err := registry.RegisterAgentWallet(context.Background(), repository.UpsertAgentWalletInput{
		AgentID:            "agent_test_1",
		UserWallet:         "0x1111111111111111111111111111111111111111",
		AgentWalletAddress: "0x9999999999999999999999999999999999999999",
		WalletProvider:     agent.WalletProviderCircleAgentWallet,
		Chain:              agent.ChainArcTestnet,
		AllowedActions:     allowedActions,
		Status:             agent.WalletStatusActive,
		PolicyMetadata:     json.RawMessage(`{"source":"test"}`),
	})
	if err != nil {
		t.Fatalf("register test agent wallet: %v", err)
	}
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
