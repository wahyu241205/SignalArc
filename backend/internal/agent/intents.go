package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"
)

const (
	ActionCreateMarket  = "create_market"
	ActionBuyYes        = "buy_yes"
	ActionBuyNo         = "buy_no"
	ActionCancelMarket  = "cancel_market"
	ActionCloseMarket   = "close_market"
	ActionResolveMarket = "resolve_market"
	ActionClaimRefund   = "claim_refund"
	ActionClaimPayout   = "claim_payout"

	StatusPreview   = "preview"
	StatusConfirmed = "confirmed"
	StatusExecuted  = "executed"

	ExecutionModeAgentContract = "agent_contract"
	NetworkArcTestnet          = "arc_testnet"
	AgentFactoryAddress        = "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80"
)

var (
	ErrIntentNotFound = errors.New("agent intent not found")
	ErrIntentInvalid  = errors.New("agent intent is invalid")
)

type CreateIntentInput struct {
	Action                string
	UserWallet            string
	MarketID              string
	MarketContractAddress string
	Amount                string
	Outcome               string
	Resolver              string
	CollateralToken       string
	CloseTimestamp        string
	Question              string
}

type Intent struct {
	ID                    string
	Action                string
	Status                string
	RequiresConfirmation  bool
	UserWallet            string
	MarketID              string
	MarketContractAddress string
	Amount                string
	Outcome               string
	Resolver              string
	CollateralToken       string
	CloseTimestamp        string
	Question              string
	ValidationResult      ValidationResult
	Warnings              []string
	CreatedAt             time.Time
}

type ExecutionPlan struct {
	IntentID            string
	Action              string
	Status              string
	ExecutionMode       string
	Network             string
	AgentFactoryAddress string
	RequiresSignature   bool
	BroadcastPerformed  bool
	TransactionHash     *string
	TransactionRequest  TransactionRequest
	Warnings            []string
}

type TransactionRequest struct {
	To                 string
	Contract           string
	Function           string
	Args               []string
	Value              string
	Chain              string
	BroadcastPerformed bool
}

type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

type Store struct {
	mu      sync.RWMutex
	intents map[string]Intent
	now     func() time.Time
	newID   func() (string, error)
}

func NewStore() *Store {
	return &Store{
		intents: make(map[string]Intent),
		now:     time.Now,
		newID:   newIntentID,
	}
}

func (store *Store) CreateIntent(input CreateIntentInput) (Intent, error) {
	normalized := normalizeInput(input)
	validationResult := validateIntent(normalized)
	intentID, err := store.newID()
	if err != nil {
		return Intent{}, err
	}

	intent := Intent{
		ID:                    intentID,
		Action:                normalized.Action,
		Status:                StatusPreview,
		RequiresConfirmation:  true,
		UserWallet:            normalized.UserWallet,
		MarketID:              normalized.MarketID,
		MarketContractAddress: normalized.MarketContractAddress,
		Amount:                normalized.Amount,
		Outcome:               normalized.Outcome,
		Resolver:              normalized.Resolver,
		CollateralToken:       normalized.CollateralToken,
		CloseTimestamp:        normalized.CloseTimestamp,
		Question:              normalized.Question,
		ValidationResult:      validationResult,
		Warnings: []string{
			"preview only; no transaction has been executed",
			"Circle Agent Wallet integration is not enabled",
			"contract execution wiring is not enabled",
		},
		CreatedAt: store.now().UTC(),
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	store.intents[intent.ID] = intent

	return intent, nil
}

func (store *Store) GetIntent(id string) (Intent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	intent, ok := store.intents[strings.TrimSpace(id)]
	if !ok {
		return Intent{}, ErrIntentNotFound
	}

	return intent, nil
}

func (store *Store) ConfirmIntent(id string) (ExecutionPlan, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	intent, ok := store.intents[strings.TrimSpace(id)]
	if !ok {
		return ExecutionPlan{}, ErrIntentNotFound
	}

	if !intent.ValidationResult.Valid {
		return ExecutionPlan{}, ErrIntentInvalid
	}

	intent.Status = StatusConfirmed
	store.intents[intent.ID] = intent

	return NewExecutionPlan(intent), nil
}

func NewExecutionPlan(intent Intent) ExecutionPlan {
	return ExecutionPlan{
		IntentID:            intent.ID,
		Action:              intent.Action,
		Status:              StatusConfirmed,
		ExecutionMode:       ExecutionModeAgentContract,
		Network:             NetworkArcTestnet,
		AgentFactoryAddress: AgentFactoryAddress,
		RequiresSignature:   true,
		BroadcastPerformed:  false,
		TransactionHash:     nil,
		TransactionRequest:  NewTransactionRequest(intent),
		Warnings: []string{
			"confirmation produced an execution plan only; no transaction has been executed",
			"no private key, signing, RPC call, or broadcast was performed",
			"Circle Agent Wallet integration is not enabled",
		},
	}
}

func NewTransactionRequest(intent Intent) TransactionRequest {
	transactionRequest := TransactionRequest{
		Value:              "0",
		Chain:              NetworkArcTestnet,
		BroadcastPerformed: false,
	}

	switch intent.Action {
	case ActionCreateMarket:
		transactionRequest.To = AgentFactoryAddress
		transactionRequest.Contract = "SignalArcAgentMarketFactory"
		transactionRequest.Function = "createMarket"
		transactionRequest.Args = []string{
			intent.MarketID,
			intent.Question,
			intent.CloseTimestamp,
			intent.Resolver,
			intent.CollateralToken,
		}
	case ActionBuyYes:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "buyYes"
		transactionRequest.Args = []string{intent.Amount}
	case ActionBuyNo:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "buyNo"
		transactionRequest.Args = []string{intent.Amount}
	case ActionCancelMarket:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "cancelMarket"
		transactionRequest.Args = []string{}
	case ActionCloseMarket:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "closeMarket"
		transactionRequest.Args = []string{}
	case ActionResolveMarket:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "resolve"
		transactionRequest.Args = []string{intent.Outcome}
	case ActionClaimRefund:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "claimRefund"
		transactionRequest.Args = []string{}
	case ActionClaimPayout:
		transactionRequest.To = intent.MarketContractAddress
		transactionRequest.Contract = "SignalArcAgentMarket"
		transactionRequest.Function = "claimPayout"
		transactionRequest.Args = []string{}
	}

	return transactionRequest
}

func normalizeInput(input CreateIntentInput) CreateIntentInput {
	return CreateIntentInput{
		Action:                strings.TrimSpace(input.Action),
		UserWallet:            strings.TrimSpace(input.UserWallet),
		MarketID:              strings.TrimSpace(input.MarketID),
		MarketContractAddress: strings.TrimSpace(input.MarketContractAddress),
		Amount:                strings.TrimSpace(input.Amount),
		Outcome:               strings.TrimSpace(input.Outcome),
		Resolver:              strings.TrimSpace(input.Resolver),
		CollateralToken:       strings.TrimSpace(input.CollateralToken),
		CloseTimestamp:        strings.TrimSpace(input.CloseTimestamp),
		Question:              strings.TrimSpace(input.Question),
	}
}

func validateIntent(input CreateIntentInput) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []string{}}

	if input.Action == "" {
		result.Errors = append(result.Errors, "action is required")
	} else if !isSupportedAction(input.Action) {
		result.Errors = append(result.Errors, "action must be one of the supported agent intent actions")
	}

	if isTransactionAction(input.Action) && input.UserWallet == "" {
		result.Errors = append(result.Errors, "user_wallet is required for transaction actions")
	}

	if isMarketSpecificAction(input.Action) && input.MarketID == "" {
		result.Errors = append(result.Errors, "market_id is required for market-specific actions")
	}

	if input.Action == ActionCreateMarket {
		if input.Question == "" {
			result.Errors = append(result.Errors, "question is required for create_market")
		}
		if input.CloseTimestamp == "" {
			result.Errors = append(result.Errors, "close_timestamp is required for create_market")
		}
		if input.Resolver == "" {
			result.Errors = append(result.Errors, "resolver is required for create_market")
		}
		if input.CollateralToken == "" {
			result.Errors = append(result.Errors, "collateral_token is required for create_market")
		}
	}

	if requiresMarketContractAddress(input.Action) && input.MarketContractAddress == "" {
		result.Errors = append(result.Errors, "market_contract_address is required for existing market contract actions")
	}

	if input.Action == ActionBuyYes || input.Action == ActionBuyNo {
		if input.Amount == "" {
			result.Errors = append(result.Errors, "amount is required for buy_yes and buy_no")
		} else if !isPositiveDecimal(input.Amount) {
			result.Errors = append(result.Errors, "amount must be positive")
		}
	}

	if input.Action == ActionResolveMarket && input.Outcome == "" {
		result.Errors = append(result.Errors, "outcome is required for resolve_market")
	}

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

func isSupportedAction(action string) bool {
	switch action {
	case ActionCreateMarket,
		ActionBuyYes,
		ActionBuyNo,
		ActionCancelMarket,
		ActionCloseMarket,
		ActionResolveMarket,
		ActionClaimRefund,
		ActionClaimPayout:
		return true
	default:
		return false
	}
}

func isTransactionAction(action string) bool {
	return isSupportedAction(action)
}

func isMarketSpecificAction(action string) bool {
	switch action {
	case ActionBuyYes,
		ActionBuyNo,
		ActionCancelMarket,
		ActionCloseMarket,
		ActionResolveMarket,
		ActionClaimRefund,
		ActionClaimPayout:
		return true
	default:
		return false
	}
}

func requiresMarketContractAddress(action string) bool {
	switch action {
	case ActionBuyYes,
		ActionBuyNo,
		ActionCancelMarket,
		ActionCloseMarket,
		ActionResolveMarket,
		ActionClaimRefund,
		ActionClaimPayout:
		return true
	default:
		return false
	}
}

func isPositiveDecimal(value string) bool {
	decimal, ok := new(big.Rat).SetString(value)
	return ok && decimal.Sign() > 0
}

func newIntentID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "agent_intent_" + hex.EncodeToString(bytes), nil
}
