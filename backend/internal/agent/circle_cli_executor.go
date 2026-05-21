package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CommandRunner interface {
	Run(context.Context, string, []string) ([]byte, error)
}

type ExecCommandRunner struct{}

func (runner ExecCommandRunner) Run(ctx context.Context, name string, args []string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	output, err := command.Output()
	if err != nil {
		return nil, errors.New("Circle CLI command failed")
	}
	return output, nil
}

type CircleCLIExecutorConfig struct {
	Enabled       bool
	CLIPath       string
	Chain         string
	Timeout       time.Duration
	AgentFactory  string
	CommandRunner CommandRunner
}

type CircleCLIExecutor struct {
	cfg CircleCLIExecutorConfig
}

func NewCircleCLIExecutor(cfg CircleCLIExecutorConfig) *CircleCLIExecutor {
	if strings.TrimSpace(cfg.CLIPath) == "" {
		cfg.CLIPath = "circle"
	}
	if strings.TrimSpace(cfg.Chain) == "" {
		cfg.Chain = ChainArcTestnet
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 120 * time.Second
	}
	if strings.TrimSpace(cfg.AgentFactory) == "" {
		cfg.AgentFactory = AgentFactoryAddress
	}
	if cfg.CommandRunner == nil {
		cfg.CommandRunner = ExecCommandRunner{}
	}
	return &CircleCLIExecutor{cfg: cfg}
}

func (executor *CircleCLIExecutor) ExecuteCreateMarket(ctx context.Context, intent Intent) (ExecutionResult, error) {
	if err := executor.validateIntent(intent, ActionCreateMarket); err != nil {
		return ExecutionResult{}, err
	}

	resolver := intent.AgentWalletAddress
	collateralToken := strings.TrimSpace(intent.CollateralToken)
	if collateralToken == "" {
		collateralToken = ArcTestnetUSDCAddress
	}

	txHash, err := executor.walletExecute(ctx,
		"createMarket(string,string,uint256,address,address)",
		[]string{intent.MarketID, intent.Question, intent.CloseTimestamp, resolver, collateralToken},
		intent.AgentWalletAddress,
		executor.cfg.AgentFactory,
	)
	if err != nil {
		return ExecutionResult{}, err
	}

	marketCount, err := executor.contractQueryUint256(ctx, "marketCount()", nil, executor.cfg.AgentFactory)
	if err != nil {
		return ExecutionResult{}, err
	}
	lastIndex, err := previousIndex(marketCount)
	if err != nil {
		return ExecutionResult{}, err
	}
	createdMarket, err := executor.contractQueryAddress(ctx, "allMarkets(uint256)", []string{lastIndex}, executor.cfg.AgentFactory)
	if err != nil {
		return ExecutionResult{}, err
	}
	isMarket, err := executor.contractQueryBool(ctx, "isMarket(address)", []string{createdMarket}, executor.cfg.AgentFactory)
	if err != nil {
		return ExecutionResult{}, err
	}

	return ExecutionResult{
		IntentID:            intent.ID,
		AgentID:             intent.AgentID,
		AgentWalletAddress:  intent.AgentWalletAddress,
		WalletProvider:      intent.WalletProvider,
		Action:              intent.Action,
		Status:              StatusExecuted,
		ExecutionMode:       ExecutionModeCircleAgentWalletCLI,
		Network:             NetworkArcTestnet,
		AgentFactoryAddress: executor.cfg.AgentFactory,
		BroadcastPerformed:  true,
		TransactionHash:     txHash,
		Readback: ExecutionReadback{
			MarketCount:   marketCount,
			CreatedMarket: createdMarket,
			IsMarket:      &isMarket,
		},
	}, nil
}

func (executor *CircleCLIExecutor) ExecuteBuyYes(ctx context.Context, intent Intent) (ExecutionResult, error) {
	return executor.executeBuy(ctx, intent, ActionBuyYes, "buyYes(uint256)", "yesPositions(address)", "totalYes()")
}

func (executor *CircleCLIExecutor) ExecuteBuyNo(ctx context.Context, intent Intent) (ExecutionResult, error) {
	return executor.executeBuy(ctx, intent, ActionBuyNo, "buyNo(uint256)", "noPositions(address)", "totalNo()")
}

func (executor *CircleCLIExecutor) executeBuy(ctx context.Context, intent Intent, expectedAction string, buySignature string, positionSignature string, totalSignature string) (ExecutionResult, error) {
	if err := executor.validateIntent(intent, expectedAction); err != nil {
		return ExecutionResult{}, err
	}
	if strings.TrimSpace(intent.MarketContractAddress) == "" || strings.TrimSpace(intent.Amount) == "" {
		return ExecutionResult{}, ErrIntentInvalid
	}

	approveHash, err := executor.walletExecute(ctx,
		"approve(address,uint256)",
		[]string{intent.MarketContractAddress, intent.Amount},
		intent.AgentWalletAddress,
		ArcTestnetUSDCAddress,
	)
	if err != nil {
		return ExecutionResult{}, err
	}

	buyHash, err := executor.walletExecute(ctx,
		buySignature,
		[]string{intent.Amount},
		intent.AgentWalletAddress,
		intent.MarketContractAddress,
	)
	if err != nil {
		return ExecutionResult{}, err
	}

	positionValue, err := executor.contractQueryUint256(ctx, positionSignature, []string{intent.AgentWalletAddress}, intent.MarketContractAddress)
	if err != nil {
		return ExecutionResult{}, err
	}
	totalValue, err := executor.contractQueryUint256(ctx, totalSignature, nil, intent.MarketContractAddress)
	if err != nil {
		return ExecutionResult{}, err
	}
	totalCollateral, err := executor.contractQueryUint256(ctx, "totalCollateral()", nil, intent.MarketContractAddress)
	if err != nil {
		return ExecutionResult{}, err
	}
	usdcBalance, err := executor.contractQueryUint256(ctx, "balanceOf(address)", []string{intent.MarketContractAddress}, ArcTestnetUSDCAddress)
	if err != nil {
		return ExecutionResult{}, err
	}

	readback := ExecutionReadback{
		TotalCollateral: totalCollateral,
		USDCBalance:     usdcBalance,
	}
	if expectedAction == ActionBuyYes {
		readback.YesPositions = positionValue
		readback.TotalYes = totalValue
	} else {
		readback.NoPositions = positionValue
		readback.TotalNo = totalValue
	}

	return ExecutionResult{
		IntentID:               intent.ID,
		AgentID:                intent.AgentID,
		AgentWalletAddress:     intent.AgentWalletAddress,
		WalletProvider:         intent.WalletProvider,
		Action:                 intent.Action,
		Status:                 StatusExecuted,
		ExecutionMode:          ExecutionModeCircleAgentWalletCLI,
		Network:                NetworkArcTestnet,
		AgentFactoryAddress:    executor.cfg.AgentFactory,
		MarketContractAddress:  intent.MarketContractAddress,
		BroadcastPerformed:     true,
		ApproveTransactionHash: approveHash,
		TransactionHash:        buyHash,
		Readback:               readback,
	}, nil
}

func (executor *CircleCLIExecutor) validateIntent(intent Intent, expectedAction string) error {
	if executor == nil {
		return ErrExecutionConfigInvalid
	}
	if !executor.cfg.Enabled {
		return ErrExecutionProviderDisabled
	}
	if intent.Status != StatusConfirmed {
		return ErrIntentNotConfirmed
	}
	if !intent.ValidationResult.Valid {
		return ErrIntentInvalid
	}
	if intent.Action != expectedAction {
		return ErrExecutionNotImplemented
	}
	if intent.WalletProvider != WalletProviderCircleAgentWallet {
		return ErrExecutionConfigInvalid
	}
	if !AgentWalletAllowsAction(AgentWallet{AllowedActions: intent.AllowedActions}, expectedAction) {
		return ErrIntentInvalid
	}
	if intent.AgentWalletAddress == "" {
		return ErrIntentInvalid
	}
	if executor.cfg.Chain != ChainArcTestnet {
		return ErrExecutionConfigInvalid
	}
	return nil
}

func (executor *CircleCLIExecutor) walletExecute(ctx context.Context, signature string, params []string, walletAddress string, contractAddress string) (string, error) {
	args := []string{"wallet", "execute", signature}
	args = append(args, params...)
	args = append(args,
		"--address", walletAddress,
		"--contract", contractAddress,
		"--chain", executor.cfg.Chain,
		"--output", "json",
	)

	output, err := executor.run(ctx, args)
	if err != nil {
		return "", err
	}
	hash, ok := findJSONValue(output, "transactionHash", "transaction_hash", "txHash", "tx_hash", "hash")
	if !ok || !looksLikeHash(hash) {
		return "", errors.New("Circle CLI transaction hash not found in JSON output")
	}
	return hash, nil
}

func (executor *CircleCLIExecutor) contractQuery(ctx context.Context, signature string, params []string, contractAddress string) (string, error) {
	args := []string{"contract", "query", signature}
	args = append(args, params...)
	args = append(args,
		"--contract", contractAddress,
		"--chain", executor.cfg.Chain,
		"--output", "json",
	)

	output, err := executor.run(ctx, args)
	if err != nil {
		return "", err
	}
	value, ok := findJSONValue(output, "result", "value", "output", "data", "returnValue")
	if !ok {
		return "", errors.New("Circle CLI query result not found in JSON output")
	}
	return value, nil
}

func (executor *CircleCLIExecutor) contractQueryUint256(ctx context.Context, signature string, params []string, contractAddress string) (string, error) {
	value, err := executor.contractQuery(ctx, signature, params, contractAddress)
	if err != nil {
		return "", err
	}
	return decodeUint256Scalar(value)
}

func (executor *CircleCLIExecutor) contractQueryAddress(ctx context.Context, signature string, params []string, contractAddress string) (string, error) {
	value, err := executor.contractQuery(ctx, signature, params, contractAddress)
	if err != nil {
		return "", err
	}
	return decodeAddressScalar(value)
}

func (executor *CircleCLIExecutor) contractQueryBool(ctx context.Context, signature string, params []string, contractAddress string) (bool, error) {
	value, err := executor.contractQuery(ctx, signature, params, contractAddress)
	if err != nil {
		return false, err
	}
	return decodeBoolScalar(value)
}

func (executor *CircleCLIExecutor) run(parent context.Context, args []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(parent, executor.cfg.Timeout)
	defer cancel()

	output, err := executor.cfg.CommandRunner.Run(ctx, executor.cfg.CLIPath, args)
	if err != nil {
		return nil, errors.New("Circle CLI command failed")
	}
	if len(output) == 0 {
		return nil, errors.New("Circle CLI command returned empty output")
	}
	return output, nil
}

func findJSONValue(output []byte, keys ...string) (string, bool) {
	var decoded any
	if err := json.Unmarshal(output, &decoded); err != nil {
		return "", false
	}
	return findValue(decoded, keys)
}

func findValue(value any, keys []string) (string, bool) {
	switch typed := value.(type) {
	case map[string]any:
		for _, key := range keys {
			if child, ok := typed[key]; ok {
				if text, ok := scalarToString(child); ok {
					return text, true
				}
				if text, ok := findValue(child, keys); ok {
					return text, true
				}
			}
		}
		for _, child := range typed {
			if text, ok := findValue(child, keys); ok {
				return text, true
			}
		}
	case []any:
		if len(typed) == 1 {
			return scalarToString(typed[0])
		}
		for _, child := range typed {
			if text, ok := findValue(child, keys); ok {
				return text, true
			}
		}
	default:
		return scalarToString(typed)
	}
	return "", false
}

func scalarToString(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed), strings.TrimSpace(typed) != ""
	case float64:
		if typed == float64(int64(typed)) {
			return fmt.Sprintf("%.0f", typed), true
		}
		return fmt.Sprintf("%v", typed), true
	case bool:
		if typed {
			return "true", true
		}
		return "false", true
	default:
		return "", false
	}
}

func previousIndex(value string) (string, error) {
	decoded, err := decodeUint256Scalar(value)
	if err != nil {
		return "", err
	}
	count, ok := new(big.Int).SetString(decoded, 10)
	if !ok || count.Sign() <= 0 {
		return "", errors.New("invalid marketCount readback")
	}
	return new(big.Int).Sub(count, big.NewInt(1)).String(), nil
}

func decodeUint256Scalar(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("empty uint256 readback")
	}
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		integer, ok := new(big.Int).SetString(strings.TrimPrefix(strings.TrimPrefix(value, "0x"), "0X"), 16)
		if !ok {
			return "", errors.New("invalid hex uint256 readback")
		}
		return integer.String(), nil
	}
	integer, ok := new(big.Int).SetString(value, 10)
	if !ok || integer.Sign() < 0 {
		return "", errors.New("invalid decimal uint256 readback")
	}
	return integer.String(), nil
}

func decodeAddressScalar(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("empty address readback")
	}
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		hexValue := strings.TrimPrefix(strings.TrimPrefix(value, "0x"), "0X")
		if len(hexValue) == 64 {
			hexValue = hexValue[24:]
		}
		if len(hexValue) != 40 {
			return "", errors.New("invalid address readback")
		}
		return common.HexToAddress("0x" + hexValue).Hex(), nil
	}
	return "", errors.New("invalid address readback")
}

func decodeBoolScalar(value string) (bool, error) {
	value = strings.TrimSpace(value)
	if strings.EqualFold(value, "true") {
		return true, nil
	}
	if strings.EqualFold(value, "false") {
		return false, nil
	}
	decoded, err := decodeUint256Scalar(value)
	if err != nil {
		return false, err
	}
	switch decoded {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, errors.New("invalid bool readback")
	}
}

func looksLikeHash(value string) bool {
	value = strings.TrimSpace(value)
	return strings.HasPrefix(value, "0x") && len(value) >= 10
}
