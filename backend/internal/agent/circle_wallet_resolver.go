package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ErrCircleAgentWalletNotFound            = errors.New("Circle Agent Wallet not found")
	ErrCircleAgentWalletResolutionAmbiguous = errors.New("Circle Agent Wallet resolution ambiguous")
	ErrCircleAgentWalletResolutionFailed    = errors.New("Circle Agent Wallet resolution failed")
	ErrCircleAgentWalletBalanceFailed       = errors.New("Circle Agent Wallet balance lookup failed")
)

var evmAddressPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

type CircleAgentWallet struct {
	Address string
	Chain   string
}

type CircleAgentWalletBalances struct {
	Balances []any
}

type CircleWalletResolver interface {
	ResolveAgentWallet(context.Context, string) (CircleAgentWallet, error)
	GetAgentWalletBalances(context.Context, string) (CircleAgentWalletBalances, error)
}

type CircleCLIWalletResolverConfig struct {
	CLIPath       string
	Chain         string
	Timeout       time.Duration
	CommandRunner EnvCommandRunner
}

type CircleCLIWalletResolver struct {
	cfg CircleCLIWalletResolverConfig
}

func NewCircleCLIWalletResolver(cfg CircleCLIWalletResolverConfig) *CircleCLIWalletResolver {
	if strings.TrimSpace(cfg.CLIPath) == "" {
		cfg.CLIPath = "circle"
	}
	if strings.TrimSpace(cfg.Chain) == "" {
		cfg.Chain = ChainArcTestnet
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 120 * time.Second
	}
	if cfg.CommandRunner == nil {
		cfg.CommandRunner = ExecEnvCommandRunner{}
	}
	return &CircleCLIWalletResolver{cfg: cfg}
}

func (resolver *CircleCLIWalletResolver) ResolveAgentWallet(parent context.Context, email string) (CircleAgentWallet, error) {
	if resolver == nil || resolver.cfg.Chain != ChainArcTestnet {
		return CircleAgentWallet{}, ErrCircleAgentWalletResolutionFailed
	}
	ctx, cancel := context.WithTimeout(parent, resolver.cfg.Timeout)
	defer cancel()

	args := []string{"wallet", "list", "--type", "agent", "--chain", ChainArcTestnet, "--output", "json"}
	output, err := resolver.cfg.CommandRunner.RunWithEnv(ctx, resolver.cfg.CLIPath, args, []string{"CIRCLE_ACCEPT_TERMS=1"})
	if err != nil {
		resolver.logResolutionFailure("circle_agent_wallet_list", string(output), err.Error(), email, "")
		return CircleAgentWallet{}, ErrCircleAgentWalletResolutionFailed
	}

	wallets, parseErr := parseCircleAgentWallets(output)
	if parseErr != nil {
		resolver.logResolutionFailure("circle_agent_wallet_list", string(output), parseErr.Error(), email, "")
		return CircleAgentWallet{}, ErrCircleAgentWalletResolutionFailed
	}
	switch len(wallets) {
	case 0:
		resolver.logResolutionFailure("circle_agent_wallet_list", string(output), "", email, "")
		return CircleAgentWallet{}, ErrCircleAgentWalletNotFound
	case 1:
		return wallets[0], nil
	default:
		resolver.logResolutionFailure("circle_agent_wallet_list", string(output), "multiple ARC-TESTNET agent wallets returned; exact onboarding match is unknown / not documented", email, "")
		return CircleAgentWallet{}, ErrCircleAgentWalletResolutionAmbiguous
	}
}

func (resolver *CircleCLIWalletResolver) GetAgentWalletBalances(parent context.Context, address string) (CircleAgentWalletBalances, error) {
	if resolver == nil || resolver.cfg.Chain != ChainArcTestnet || strings.TrimSpace(address) == "" {
		return CircleAgentWalletBalances{}, ErrCircleAgentWalletBalanceFailed
	}
	ctx, cancel := context.WithTimeout(parent, resolver.cfg.Timeout)
	defer cancel()

	address = strings.TrimSpace(address)
	args := []string{"wallet", "balance", "--address", address, "--chain", ChainArcTestnet, "--output", "json"}
	output, err := resolver.cfg.CommandRunner.RunWithEnv(ctx, resolver.cfg.CLIPath, args, []string{"CIRCLE_ACCEPT_TERMS=1"})
	if err != nil {
		resolver.logResolutionFailure("circle_agent_wallet_balance", string(output), err.Error(), "", address)
		return CircleAgentWalletBalances{}, ErrCircleAgentWalletBalanceFailed
	}

	balances, parseErr := parseCircleAgentWalletBalances(output)
	if parseErr != nil {
		resolver.logResolutionFailure("circle_agent_wallet_balance", string(output), parseErr.Error(), "", address)
		return CircleAgentWalletBalances{}, ErrCircleAgentWalletBalanceFailed
	}
	return CircleAgentWalletBalances{Balances: balances}, nil
}

func (resolver *CircleCLIWalletResolver) logResolutionFailure(operation string, output string, errText string, email string, address string) {
	sanitizedOutput := sanitizeCircleOnboardingText(output, email, address)
	sanitizedError := sanitizeCircleOnboardingText(errText, email, address)
	log.Error().
		Str("operation", operation).
		Str("output", sanitizedOutput).
		Str("error", sanitizedError).
		Msg("Circle CLI wallet read failed")
}

// extractJSONFromCLIOutput finds the first JSON object or array in raw CLI
// output that may be prefixed with non-JSON text such as Node.js deprecation
// warnings. It returns the byte slice starting at the first '{' or '[' that
// successfully parses as JSON. If no valid JSON envelope is found, it returns
// nil and an error.
func extractJSONFromCLIOutput(raw []byte) ([]byte, error) {
	// Fast path: if the output starts with JSON, return as-is.
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		var probe json.RawMessage
		if json.Unmarshal(trimmed, &probe) == nil {
			return trimmed, nil
		}
	}

	// Scan for each '{' or '[' and try to parse from that position.
	for i := 0; i < len(raw); i++ {
		if raw[i] != '{' && raw[i] != '[' {
			continue
		}
		candidate := raw[i:]
		var probe json.RawMessage
		if json.Unmarshal(candidate, &probe) == nil {
			return candidate, nil
		}
	}

	return nil, errors.New("no valid JSON object or array found in CLI output")
}

func parseCircleAgentWallets(output []byte) ([]CircleAgentWallet, error) {
	cleaned, extractErr := extractJSONFromCLIOutput(output)
	if extractErr != nil {
		return nil, extractErr
	}
	var decoded any
	if err := json.Unmarshal(cleaned, &decoded); err != nil {
		return nil, err
	}
	seen := map[string]CircleAgentWallet{}
	collectCircleAgentWallets(decoded, seen)
	wallets := make([]CircleAgentWallet, 0, len(seen))
	for _, wallet := range seen {
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func collectCircleAgentWallets(value any, seen map[string]CircleAgentWallet) {
	switch typed := value.(type) {
	case map[string]any:
		if address := circleWalletAddressFromMap(typed); address != "" {
			chain := circleWalletChainFromMap(typed)
			if chain == "" || strings.EqualFold(chain, ChainArcTestnet) {
				seen[strings.ToLower(address)] = CircleAgentWallet{Address: address, Chain: ChainArcTestnet}
			}
		}
		for _, child := range typed {
			collectCircleAgentWallets(child, seen)
		}
	case []any:
		for _, child := range typed {
			collectCircleAgentWallets(child, seen)
		}
	}
}

func circleWalletAddressFromMap(value map[string]any) string {
	for _, key := range []string{"address", "walletAddress", "wallet_address"} {
		if text, ok := scalarToString(value[key]); ok && evmAddressPattern.MatchString(text) {
			return text
		}
	}
	return ""
}

func circleWalletChainFromMap(value map[string]any) string {
	for _, key := range []string{"chain", "network", "blockchain"} {
		if text, ok := scalarToString(value[key]); ok {
			return text
		}
	}
	return ""
}

func parseCircleAgentWalletBalances(output []byte) ([]any, error) {
	cleaned, extractErr := extractJSONFromCLIOutput(output)
	if extractErr != nil {
		return nil, extractErr
	}
	var decoded any
	if err := json.Unmarshal(cleaned, &decoded); err != nil {
		return nil, err
	}
	balances, ok := findBalancesArray(decoded)
	if !ok {
		return nil, errors.New("balances array not found")
	}
	return balances, nil
}

func findBalancesArray(value any) ([]any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		if balances, ok := typed["balances"].([]any); ok {
			return balances, true
		}
		for _, child := range typed {
			if balances, ok := findBalancesArray(child); ok {
				return balances, true
			}
		}
	case []any:
		for _, child := range typed {
			if balances, ok := findBalancesArray(child); ok {
				return balances, true
			}
		}
	}
	return nil, false
}
