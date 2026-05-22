package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	ErrCircleOnboardingDisabled                = errors.New("Circle Agent Wallet OTP onboarding start is disabled")
	ErrCircleOnboardingRequestIDNotDocumented  = errors.New("Circle CLI OTP start request_id field is unknown / not documented")
	ErrCircleOnboardingRequestIDMissing        = errors.New("Circle CLI OTP start request_id not found")
	ErrCircleOnboardingRequestIDEmpty          = errors.New("Circle CLI OTP start request_id is empty")
	ErrCircleOnboardingCommandFailed           = errors.New("Circle CLI OTP start command failed")
	ErrCircleOnboardingCommandReturnedNoOutput = errors.New("Circle CLI OTP start command returned empty output")
	ErrCircleOnboardingUnsupportedChain        = errors.New("Circle Agent Wallet onboarding chain must be ARC-TESTNET")
	ErrCircleOnboardingEmailRequired           = errors.New("user_email is required for Circle Agent Wallet OTP onboarding")
)

type CircleOTPStartResult struct {
	RequestID string
	ExpiresAt time.Time
}

type CircleOnboardingRunner interface {
	StartOTP(context.Context, string) (CircleOTPStartResult, error)
}

type CircleOnboardingStarter struct {
	Enabled      bool
	Runner       CircleOnboardingRunner
	RequestStore *CircleOTPRequestStore
}

func (starter CircleOnboardingStarter) StartOTP(ctx context.Context, onboardingID string, email string) (CircleOTPStartResult, string, bool, error) {
	if !starter.Enabled {
		return CircleOTPStartResult{}, "", false, ErrCircleOnboardingDisabled
	}
	if starter.Runner == nil {
		return CircleOTPStartResult{}, "", true, ErrCircleOnboardingCommandFailed
	}

	result, err := starter.Runner.StartOTP(ctx, email)
	if err != nil {
		return CircleOTPStartResult{}, "", true, err
	}
	requestID := strings.TrimSpace(result.RequestID)
	if requestID == "" {
		return CircleOTPStartResult{}, "", true, ErrCircleOnboardingRequestIDEmpty
	}
	if result.ExpiresAt.IsZero() {
		result.ExpiresAt = time.Now().UTC().Add(10 * time.Minute)
	}
	result.RequestID = requestID
	requestIDHash := HashCircleRequestID(requestID)
	if starter.RequestStore != nil {
		starter.RequestStore.Save(onboardingID, requestID)
	}
	return result, requestIDHash, true, nil
}

type CircleOTPRequestStore struct {
	mu                      sync.RWMutex
	requestIDByOnboardingID map[string]string
}

func NewCircleOTPRequestStore() *CircleOTPRequestStore {
	return &CircleOTPRequestStore{requestIDByOnboardingID: map[string]string{}}
}

func (store *CircleOTPRequestStore) Save(onboardingID string, requestID string) {
	if store == nil {
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	store.requestIDByOnboardingID[strings.TrimSpace(onboardingID)] = strings.TrimSpace(requestID)
}

func (store *CircleOTPRequestStore) Get(onboardingID string) (string, bool) {
	if store == nil {
		return "", false
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	requestID, ok := store.requestIDByOnboardingID[strings.TrimSpace(onboardingID)]
	return requestID, ok
}

func HashCircleRequestID(requestID string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(requestID)))
	return hex.EncodeToString(sum[:])
}

type CircleCLIOnboardingRunnerConfig struct {
	CLIPath       string
	Chain         string
	Timeout       time.Duration
	CommandRunner EnvCommandRunner
}

type EnvCommandRunner interface {
	RunWithEnv(context.Context, string, []string, []string) ([]byte, error)
}

type ExecEnvCommandRunner struct{}

func (runner ExecEnvCommandRunner) RunWithEnv(ctx context.Context, name string, args []string, env []string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Env = append(os.Environ(), env...)
	output, err := command.Output()
	if err != nil {
		return nil, ErrCircleOnboardingCommandFailed
	}
	return output, nil
}

type CircleCLIOnboardingRunner struct {
	cfg CircleCLIOnboardingRunnerConfig
}

func NewCircleCLIOnboardingRunner(cfg CircleCLIOnboardingRunnerConfig) *CircleCLIOnboardingRunner {
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
	return &CircleCLIOnboardingRunner{cfg: cfg}
}

func (runner *CircleCLIOnboardingRunner) StartOTP(parent context.Context, email string) (CircleOTPStartResult, error) {
	if runner == nil {
		return CircleOTPStartResult{}, ErrCircleOnboardingCommandFailed
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return CircleOTPStartResult{}, ErrCircleOnboardingEmailRequired
	}
	if runner.cfg.Chain != ChainArcTestnet {
		return CircleOTPStartResult{}, ErrCircleOnboardingUnsupportedChain
	}

	ctx, cancel := context.WithTimeout(parent, runner.cfg.Timeout)
	defer cancel()

	args := []string{"wallet", "login", email, "--init", "--type", "agent", "--testnet", "--output", "json"}
	output, err := runner.cfg.CommandRunner.RunWithEnv(ctx, runner.cfg.CLIPath, args, []string{"CIRCLE_ACCEPT_TERMS=1"})
	if err != nil {
		return CircleOTPStartResult{}, ErrCircleOnboardingCommandFailed
	}
	if len(output) == 0 {
		return CircleOTPStartResult{}, ErrCircleOnboardingCommandReturnedNoOutput
	}

	requestID, ok := findCircleOTPRequestID(output)
	if !ok {
		return CircleOTPStartResult{}, ErrCircleOnboardingRequestIDNotDocumented
	}
	if strings.TrimSpace(requestID) == "" {
		return CircleOTPStartResult{}, ErrCircleOnboardingRequestIDMissing
	}

	return CircleOTPStartResult{
		RequestID: strings.TrimSpace(requestID),
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}, nil
}

func findCircleOTPRequestID(output []byte) (string, bool) {
	var decoded any
	if err := json.Unmarshal(output, &decoded); err != nil {
		return "", false
	}
	// Circle documents that --init returns a request ID, but the exact JSON
	// field name is unknown / not documented in the official pages reviewed.
	return findValue(decoded, []string{"request_id", "requestId", "requestID", "id"})
}
