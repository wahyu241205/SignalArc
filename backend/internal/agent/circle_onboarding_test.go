package agent

import (
	"context"
	"slices"
	"testing"
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
