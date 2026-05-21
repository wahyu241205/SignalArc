package agent

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fakeCommandRunner struct {
	outputs [][]byte
	errAt   int
	calls   []fakeCommandCall
}

type fakeCommandCall struct {
	name string
	args []string
}

func (runner *fakeCommandRunner) Run(_ context.Context, name string, args []string) ([]byte, error) {
	runner.calls = append(runner.calls, fakeCommandCall{name: name, args: append([]string{}, args...)})
	if runner.errAt > 0 && len(runner.calls) == runner.errAt {
		return nil, errors.New("raw sensitive CLI error")
	}
	if len(runner.outputs) == 0 {
		return nil, errors.New("missing fake output")
	}
	output := runner.outputs[0]
	runner.outputs = runner.outputs[1:]
	return output, nil
}

func TestCircleCLIExecutorCreateMarketBuildsCommandsAndReadback(t *testing.T) {
	runner := &fakeCommandRunner{outputs: [][]byte{
		[]byte(`{"transactionHash":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`),
		jsonResult(abiUint256("9")),
		jsonResult(abiAddress("0xabcf081e456c1a11106def590666a07b76d456f8")),
		jsonResult(abiUint256("1")),
	}}
	executor := newTestCircleCLIExecutor(runner, true)
	intent := confirmedIntent(ActionCreateMarket)
	intent.Resolver = "0x2222222222222222222222222222222222222222"

	result, err := executor.ExecuteCreateMarket(context.Background(), intent)
	if err != nil {
		t.Fatalf("execute create market: %v", err)
	}

	expectedFirstArgs := []string{
		"wallet", "execute", "createMarket(string,string,uint256,address,address)",
		"market-1", "Will SignalArc execute through Circle?", "1770000000",
		"0x9999999999999999999999999999999999999999", ArcTestnetUSDCAddress,
		"--address", "0x9999999999999999999999999999999999999999",
		"--contract", AgentFactoryAddress,
		"--chain", ChainArcTestnet,
		"--output", "json",
	}
	assertArgs(t, runner.calls[0].args, expectedFirstArgs)
	assertArgs(t, runner.calls[2].args, []string{
		"contract", "query", "allMarkets(uint256)", "8",
		"--contract", AgentFactoryAddress,
		"--chain", ChainArcTestnet,
		"--output", "json",
	})
	if result.ExecutionMode != ExecutionModeCircleAgentWalletCLI {
		t.Fatalf("expected circle execution mode, got %q", result.ExecutionMode)
	}
	if result.TransactionHash != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("unexpected tx hash %q", result.TransactionHash)
	}
	if result.Readback.MarketCount != "9" || result.Readback.CreatedMarket != "0xabcF081E456c1A11106DeF590666A07B76d456f8" {
		t.Fatalf("unexpected readback %#v", result.Readback)
	}
	if result.Readback.IsMarket == nil || !*result.Readback.IsMarket {
		t.Fatalf("expected is market true, got %#v", result.Readback.IsMarket)
	}
}

func TestCircleCLIExecutorBuyYesBuildsApproveAndBuyCommands(t *testing.T) {
	runner := &fakeCommandRunner{outputs: buyOutputs("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "1000000", "1000000", "1000000", "1000000")}
	executor := newTestCircleCLIExecutor(runner, true)

	result, err := executor.ExecuteBuyYes(context.Background(), confirmedIntent(ActionBuyYes))
	if err != nil {
		t.Fatalf("execute buy yes: %v", err)
	}

	assertArgs(t, runner.calls[0].args, []string{
		"wallet", "execute", "approve(address,uint256)",
		"0x3333333333333333333333333333333333333333", "1000000",
		"--address", "0x9999999999999999999999999999999999999999",
		"--contract", ArcTestnetUSDCAddress,
		"--chain", ChainArcTestnet,
		"--output", "json",
	})
	assertArgs(t, runner.calls[1].args, []string{
		"wallet", "execute", "buyYes(uint256)", "1000000",
		"--address", "0x9999999999999999999999999999999999999999",
		"--contract", "0x3333333333333333333333333333333333333333",
		"--chain", ChainArcTestnet,
		"--output", "json",
	})
	if result.ApproveTransactionHash != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" || result.TransactionHash != "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Fatalf("unexpected hashes %#v", result)
	}
	if result.Readback.YesPositions != "1000000" || result.Readback.TotalYes != "1000000" {
		t.Fatalf("unexpected readback %#v", result.Readback)
	}
}

func TestCircleCLIExecutorBuyNoBuildsApproveAndBuyCommands(t *testing.T) {
	runner := &fakeCommandRunner{outputs: buyOutputs("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "1000000", "1000000", "1000000", "1000000")}
	executor := newTestCircleCLIExecutor(runner, true)

	result, err := executor.ExecuteBuyNo(context.Background(), confirmedIntent(ActionBuyNo))
	if err != nil {
		t.Fatalf("execute buy no: %v", err)
	}

	assertArgs(t, runner.calls[1].args, []string{
		"wallet", "execute", "buyNo(uint256)", "1000000",
		"--address", "0x9999999999999999999999999999999999999999",
		"--contract", "0x3333333333333333333333333333333333333333",
		"--chain", ChainArcTestnet,
		"--output", "json",
	})
	if result.Readback.NoPositions != "1000000" || result.Readback.TotalNo != "1000000" {
		t.Fatalf("unexpected readback %#v", result.Readback)
	}
}

func TestCircleCLIExecutorDisabledFailsClosed(t *testing.T) {
	executor := newTestCircleCLIExecutor(&fakeCommandRunner{}, false)

	_, err := executor.ExecuteCreateMarket(context.Background(), confirmedIntent(ActionCreateMarket))
	if !errors.Is(err, ErrExecutionProviderDisabled) {
		t.Fatalf("expected disabled error, got %v", err)
	}
}

func TestCircleCLIExecutorLifecycleActionsRemainNotImplemented(t *testing.T) {
	executor := newTestCircleCLIExecutor(&fakeCommandRunner{}, true)

	_, err := executor.ExecuteCreateMarket(context.Background(), confirmedIntent(ActionCloseMarket))
	if !errors.Is(err, ErrExecutionNotImplemented) {
		t.Fatalf("expected not implemented, got %v", err)
	}
}

func TestCircleCLIExecutorDisallowedActionFailsClosed(t *testing.T) {
	executor := newTestCircleCLIExecutor(&fakeCommandRunner{}, true)
	intent := confirmedIntent(ActionBuyYes)
	intent.AllowedActions = []string{ActionBuyNo}

	_, err := executor.ExecuteBuyYes(context.Background(), intent)
	if !errors.Is(err, ErrIntentInvalid) {
		t.Fatalf("expected invalid intent, got %v", err)
	}
}

func TestCircleCLIExecutorReadbackErrorIsSanitized(t *testing.T) {
	runner := &fakeCommandRunner{
		outputs: [][]byte{[]byte(`{"transactionHash":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)},
		errAt:   2,
	}
	executor := newTestCircleCLIExecutor(runner, true)

	_, err := executor.ExecuteCreateMarket(context.Background(), confirmedIntent(ActionCreateMarket))
	if err == nil {
		t.Fatal("expected readback error")
	}
	if err.Error() == "raw sensitive CLI error" {
		t.Fatal("expected sanitized error")
	}
}

func TestCircleCLIExecutorDecodesABIHexScalars(t *testing.T) {
	uintValue, err := decodeUint256Scalar("0x0000000000000000000000000000000000000000000000000000000000000008")
	if err != nil {
		t.Fatalf("decode uint256: %v", err)
	}
	if uintValue != "8" {
		t.Fatalf("expected 8, got %q", uintValue)
	}

	index, err := previousIndex("0x0000000000000000000000000000000000000000000000000000000000000008")
	if err != nil {
		t.Fatalf("previous index: %v", err)
	}
	if index != "7" {
		t.Fatalf("expected previous index 7, got %q", index)
	}

	address, err := decodeAddressScalar("0x000000000000000000000000abcf081e456c1a11106def590666a07b76d456f8")
	if err != nil {
		t.Fatalf("decode address: %v", err)
	}
	if address != "0xabcF081E456c1A11106DeF590666A07B76d456f8" {
		t.Fatalf("unexpected address %q", address)
	}

	boolean, err := decodeBoolScalar("0x0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatalf("decode bool: %v", err)
	}
	if !boolean {
		t.Fatal("expected true")
	}
}

func newTestCircleCLIExecutor(runner *fakeCommandRunner, enabled bool) *CircleCLIExecutor {
	return NewCircleCLIExecutor(CircleCLIExecutorConfig{
		Enabled:       enabled,
		CLIPath:       "circle",
		Chain:         ChainArcTestnet,
		Timeout:       time.Second,
		AgentFactory:  AgentFactoryAddress,
		CommandRunner: runner,
	})
}

func confirmedIntent(action string) Intent {
	return Intent{
		ID:                    "agent_intent_test",
		AgentID:               "agent_test",
		AgentWalletAddress:    "0x9999999999999999999999999999999999999999",
		WalletProvider:        WalletProviderCircleAgentWallet,
		AllowedActions:        []string{ActionCreateMarket, ActionBuyYes, ActionBuyNo},
		Action:                action,
		Status:                StatusConfirmed,
		UserWallet:            "0x1111111111111111111111111111111111111111",
		MarketID:              "market-1",
		MarketContractAddress: "0x3333333333333333333333333333333333333333",
		Amount:                "1000000",
		Resolver:              "",
		CollateralToken:       ArcTestnetUSDCAddress,
		CloseTimestamp:        "1770000000",
		Question:              "Will SignalArc execute through Circle?",
		ValidationResult:      ValidationResult{Valid: true, Errors: []string{}},
	}
}

func buyOutputs(approveHash string, buyHash string, position string, total string, collateral string, balance string) [][]byte {
	return [][]byte{
		[]byte(`{"transactionHash":"` + approveHash + `"}`),
		[]byte(`{"transactionHash":"` + buyHash + `"}`),
		jsonResult(abiUint256(position)),
		jsonResult(abiUint256(total)),
		jsonResult(abiUint256(collateral)),
		jsonResult(abiUint256(balance)),
	}
}

func jsonResult(value string) []byte {
	return []byte(`{"result":"` + value + `"}`)
}

func abiUint256(decimal string) string {
	value, ok := new(big.Int).SetString(decimal, 10)
	if !ok {
		panic("invalid test decimal")
	}
	return "0x" + fmt.Sprintf("%064x", value)
}

func abiAddress(address string) string {
	trimmed := strings.TrimPrefix(strings.TrimPrefix(address, "0x"), "0X")
	return "0x" + strings.Repeat("0", 64-len(trimmed)) + trimmed
}

func assertArgs(t *testing.T, actual []string, expected []string) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected args %#v, got %#v", expected, actual)
	}
}
