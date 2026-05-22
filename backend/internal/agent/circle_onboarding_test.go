package agent

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type fakeEnvCommandRunner struct {
	name string
	args []string
	env  []string
}

func (runner *fakeEnvCommandRunner) RunWithEnv(_ context.Context, name string, args []string, env []string) ([]byte, error) {
	runner.name = name
	runner.args = append([]string{}, args...)
	runner.env = append([]string{}, env...)
	return []byte(`{"ok":true}`), nil
}

type failingEnvCommandRunner struct {
	output []byte
	err    error
}

func (runner failingEnvCommandRunner) RunWithEnv(_ context.Context, _ string, _ []string, _ []string) ([]byte, error) {
	return runner.output, runner.err
}

func TestCircleCLIOnboardingRunnerVerifyOTPUsesDocumentedCommandShape(t *testing.T) {
	commandRunner := &fakeEnvCommandRunner{}
	runner := NewCircleCLIOnboardingRunner(CircleCLIOnboardingRunnerConfig{
		CLIPath:       "circle",
		Chain:         ChainArcTestnet,
		CommandRunner: commandRunner,
	})

	if err := runner.VerifyOTP(context.Background(), "request-123", "SIN-232794"); err != nil {
		t.Fatalf("verify otp: %v", err)
	}

	if commandRunner.name != "circle" {
		t.Fatalf("expected circle command, got %q", commandRunner.name)
	}
	expectedArgs := []string{"wallet", "login", "--request", "request-123", "--otp", "SIN-232794"}
	if !slices.Equal(commandRunner.args, expectedArgs) {
		t.Fatalf("unexpected args %#v", commandRunner.args)
	}
	if !slices.Contains(commandRunner.env, "CIRCLE_ACCEPT_TERMS=1") {
		t.Fatalf("expected CIRCLE_ACCEPT_TERMS=1 env, got %#v", commandRunner.env)
	}
}

func TestCircleCLIOnboardingRunnerVerifyOTPSanitizesFailureDiagnostics(t *testing.T) {
	requestID := "request-secret-123"
	otp := "SIN-232794"
	commandOutput := []byte("circle wallet login --request request-secret-123 --otp SIN-232794 failed: invalid or expired OTP")
	commandRunner := failingEnvCommandRunner{
		output: commandOutput,
		err:    errors.New("exit status 1 for request-secret-123 with SIN-232794"),
	}
	runner := NewCircleCLIOnboardingRunner(CircleCLIOnboardingRunnerConfig{
		CLIPath:       "circle",
		Chain:         ChainArcTestnet,
		CommandRunner: commandRunner,
	})

	var logs bytes.Buffer
	previousLogger := log.Logger
	log.Logger = zerolog.New(&logs)
	defer func() {
		log.Logger = previousLogger
	}()

	err := runner.VerifyOTP(context.Background(), requestID, otp)
	if err == nil {
		t.Fatal("expected verify error")
	}
	if !errors.Is(err, ErrCircleOnboardingCommandFailed) {
		t.Fatalf("expected command failed error, got %v", err)
	}

	errorText := err.Error()
	logText := logs.String()
	for _, text := range []string{errorText, logText} {
		if strings.Contains(text, requestID) {
			t.Fatalf("sanitized diagnostics exposed request ID: %s", text)
		}
		if strings.Contains(text, otp) {
			t.Fatalf("sanitized diagnostics exposed OTP: %s", text)
		}
		if !strings.Contains(text, "invalid or expired OTP") {
			t.Fatalf("expected sanitized diagnostic detail, got %s", text)
		}
	}
	if strings.Count(errorText, "[redacted]") < 2 {
		t.Fatalf("expected redacted request ID and OTP in error text, got %s", errorText)
	}
	if strings.Count(logText, "[redacted]") < 2 {
		t.Fatalf("expected redacted request ID and OTP in log text, got %s", logText)
	}
}
