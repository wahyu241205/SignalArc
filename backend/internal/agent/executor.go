package agent

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const agentFactoryABIJSON = `[
	{"type":"function","name":"createMarket","inputs":[{"name":"marketId","type":"string"},{"name":"question","type":"string"},{"name":"closeTimestamp","type":"uint256"},{"name":"resolver","type":"address"},{"name":"collateralToken","type":"address"}],"outputs":[{"name":"market","type":"address"}],"stateMutability":"nonpayable"},
	{"type":"function","name":"marketCount","inputs":[],"outputs":[{"name":"","type":"uint256"}],"stateMutability":"view"},
	{"type":"function","name":"isMarket","inputs":[{"name":"","type":"address"}],"outputs":[{"name":"","type":"bool"}],"stateMutability":"view"},
	{"type":"event","name":"AgentMarketDeployed","inputs":[{"name":"marketId","type":"string","indexed":true},{"name":"market","type":"address","indexed":true},{"name":"admin","type":"address","indexed":true},{"name":"resolver","type":"address","indexed":false},{"name":"collateralToken","type":"address","indexed":false},{"name":"closeTimestamp","type":"uint256","indexed":false},{"name":"question","type":"string","indexed":false}],"anonymous":false}
]`

var (
	ErrExecutionNotImplemented = errors.New("agent execution action is not implemented")
	ErrIntentNotConfirmed      = errors.New("agent intent is not confirmed")
	ErrExecutionConfigInvalid  = errors.New("agent execution config is invalid")
)

type Executor interface {
	ExecuteCreateMarket(context.Context, Intent) (ExecutionResult, error)
}

type ExecutionResult struct {
	IntentID            string
	Action              string
	Status              string
	ExecutionMode       string
	Network             string
	AgentFactoryAddress string
	BroadcastPerformed  bool
	TransactionHash     string
	Readback            ExecutionReadback
}

type ExecutionReadback struct {
	MarketCount   string
	CreatedMarket string
	IsMarket      *bool
}

type ArcExecutorConfig struct {
	RPCURL             string
	ExecutorPrivateKey string
	AgentFactory       string
}

type ArcExecutor struct {
	cfg ArcExecutorConfig
}

func NewArcExecutorFromEnv() (*ArcExecutor, error) {
	cfg := ArcExecutorConfig{
		RPCURL:             strings.TrimSpace(os.Getenv("ARC_TESTNET_RPC_URL")),
		ExecutorPrivateKey: strings.TrimSpace(os.Getenv("AGENT_EXECUTOR_PRIVATE_KEY")),
		AgentFactory:       strings.TrimSpace(os.Getenv("AGENT_FACTORY_ADDRESS")),
	}
	if cfg.AgentFactory == "" {
		cfg.AgentFactory = AgentFactoryAddress
	}

	executor := &ArcExecutor{cfg: cfg}
	if err := executor.validateConfig(); err != nil {
		return nil, err
	}

	return executor, nil
}

func (executor *ArcExecutor) ExecuteCreateMarket(ctx context.Context, intent Intent) (ExecutionResult, error) {
	if intent.Status != StatusConfirmed {
		return ExecutionResult{}, ErrIntentNotConfirmed
	}
	if !intent.ValidationResult.Valid {
		return ExecutionResult{}, ErrIntentInvalid
	}
	if intent.Action != ActionCreateMarket {
		return ExecutionResult{}, ErrExecutionNotImplemented
	}
	if err := executor.validateConfig(); err != nil {
		return ExecutionResult{}, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(agentFactoryABIJSON))
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("parse agent factory abi: %w", err)
	}

	client, err := ethclient.DialContext(ctx, executor.cfg.RPCURL)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("connect arc testnet rpc: %w", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(executor.cfg.ExecutorPrivateKey, "0x"))
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("%w: AGENT_EXECUTOR_PRIVATE_KEY is invalid", ErrExecutionConfigInvalid)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("read chain id: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("create transactor: %w", err)
	}
	auth.Context = ctx

	closeTimestamp, ok := new(big.Int).SetString(intent.CloseTimestamp, 10)
	if !ok {
		return ExecutionResult{}, fmt.Errorf("%w: close_timestamp is invalid", ErrIntentInvalid)
	}

	factoryAddress := common.HexToAddress(executor.cfg.AgentFactory)
	factory := bind.NewBoundContract(factoryAddress, parsedABI, client, client, client)
	tx, err := factory.Transact(auth, "createMarket", intent.MarketID, intent.Question, closeTimestamp, common.HexToAddress(intent.Resolver), common.HexToAddress(intent.CollateralToken))
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("broadcast createMarket: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("wait for createMarket receipt: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return ExecutionResult{}, fmt.Errorf("createMarket transaction failed with receipt status %d", receipt.Status)
	}

	createdMarket := createdMarketFromReceipt(parsedABI, factoryAddress, receipt)
	marketCount, err := readMarketCount(ctx, client, parsedABI, factoryAddress)
	if err != nil {
		return ExecutionResult{}, err
	}

	readback := ExecutionReadback{
		MarketCount: marketCount.String(),
	}
	if createdMarket != (common.Address{}) {
		readback.CreatedMarket = createdMarket.Hex()
		isMarket, err := readIsMarket(ctx, client, parsedABI, factoryAddress, createdMarket)
		if err != nil {
			return ExecutionResult{}, err
		}
		readback.IsMarket = &isMarket
	}

	return ExecutionResult{
		IntentID:            intent.ID,
		Action:              intent.Action,
		Status:              StatusExecuted,
		ExecutionMode:       ExecutionModeAgentContract,
		Network:             NetworkArcTestnet,
		AgentFactoryAddress: factoryAddress.Hex(),
		BroadcastPerformed:  true,
		TransactionHash:     tx.Hash().Hex(),
		Readback:            readback,
	}, nil
}

func (executor *ArcExecutor) validateConfig() error {
	if executor == nil {
		return ErrExecutionConfigInvalid
	}
	if executor.cfg.RPCURL == "" {
		return fmt.Errorf("%w: ARC_TESTNET_RPC_URL is required", ErrExecutionConfigInvalid)
	}
	if executor.cfg.ExecutorPrivateKey == "" {
		return fmt.Errorf("%w: AGENT_EXECUTOR_PRIVATE_KEY is required", ErrExecutionConfigInvalid)
	}
	if !common.IsHexAddress(executor.cfg.AgentFactory) {
		return fmt.Errorf("%w: AGENT_FACTORY_ADDRESS is invalid", ErrExecutionConfigInvalid)
	}
	return nil
}

func createdMarketFromReceipt(parsedABI abi.ABI, factoryAddress common.Address, receipt *types.Receipt) common.Address {
	event := parsedABI.Events["AgentMarketDeployed"]
	for _, log := range receipt.Logs {
		if log.Address != factoryAddress || len(log.Topics) < 3 || log.Topics[0] != event.ID {
			continue
		}
		return common.BytesToAddress(log.Topics[2].Bytes())
	}
	return common.Address{}
}

func readMarketCount(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, factoryAddress common.Address) (*big.Int, error) {
	data, err := parsedABI.Pack("marketCount")
	if err != nil {
		return nil, fmt.Errorf("pack marketCount call: %w", err)
	}

	output, err := client.CallContract(ctx, ethereum.CallMsg{To: &factoryAddress, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("read marketCount: %w", err)
	}

	values, err := parsedABI.Unpack("marketCount", output)
	if err != nil {
		return nil, fmt.Errorf("unpack marketCount: %w", err)
	}
	if len(values) != 1 {
		return nil, errors.New("marketCount returned unexpected value count")
	}

	marketCount, ok := values[0].(*big.Int)
	if !ok {
		return nil, errors.New("marketCount returned unexpected value type")
	}

	return marketCount, nil
}

func readIsMarket(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, factoryAddress common.Address, market common.Address) (bool, error) {
	data, err := parsedABI.Pack("isMarket", market)
	if err != nil {
		return false, fmt.Errorf("pack isMarket call: %w", err)
	}

	output, err := client.CallContract(ctx, ethereum.CallMsg{To: &factoryAddress, Data: data}, nil)
	if err != nil {
		return false, fmt.Errorf("read isMarket: %w", err)
	}

	values, err := parsedABI.Unpack("isMarket", output)
	if err != nil {
		return false, fmt.Errorf("unpack isMarket: %w", err)
	}
	if len(values) != 1 {
		return false, errors.New("isMarket returned unexpected value count")
	}

	isMarket, ok := values[0].(bool)
	if !ok {
		return false, errors.New("isMarket returned unexpected value type")
	}

	return isMarket, nil
}
