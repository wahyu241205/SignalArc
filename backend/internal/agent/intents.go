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

	StatusPreview = "preview"
)

var ErrIntentNotFound = errors.New("agent intent not found")

type CreateIntentInput struct {
	Action     string
	UserWallet string
	MarketID   string
	Amount     string
	Outcome    string
}

type Intent struct {
	ID                   string
	Action               string
	Status               string
	RequiresConfirmation bool
	UserWallet           string
	MarketID             string
	Amount               string
	Outcome              string
	ValidationResult     ValidationResult
	Warnings             []string
	CreatedAt            time.Time
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
		ID:                   intentID,
		Action:               normalized.Action,
		Status:               StatusPreview,
		RequiresConfirmation: true,
		UserWallet:           normalized.UserWallet,
		MarketID:             normalized.MarketID,
		Amount:               normalized.Amount,
		Outcome:              normalized.Outcome,
		ValidationResult:     validationResult,
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

func normalizeInput(input CreateIntentInput) CreateIntentInput {
	return CreateIntentInput{
		Action:     strings.TrimSpace(input.Action),
		UserWallet: strings.TrimSpace(input.UserWallet),
		MarketID:   strings.TrimSpace(input.MarketID),
		Amount:     strings.TrimSpace(input.Amount),
		Outcome:    strings.TrimSpace(input.Outcome),
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
