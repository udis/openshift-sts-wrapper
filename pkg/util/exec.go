package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// CommandExecutor is an interface for executing commands (allows mocking in tests)
type CommandExecutor interface {
	Execute(name string, args ...string) (string, error)
	ExecuteWithEnv(name string, env []string, args ...string) (string, error)
	ExecuteInteractive(name string, args ...string) error
	ExecuteInteractiveWithEnv(name string, env []string, args ...string) error
}

// RealExecutor executes actual system commands
type RealExecutor struct{}

func (e *RealExecutor) Execute(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (e *RealExecutor) ExecuteWithEnv(name string, env []string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (e *RealExecutor) ExecuteInteractive(name string, args ...string) error {
	// For truly interactive commands, we need to ensure the TTY is properly connected
	binary, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("failed to find command %s: %w", name, err)
	}

	cmd := exec.Command(binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Don't create a new process group - inherit the parent's terminal
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: false,
		Pgid:    0,
	}

	return cmd.Run()
}

func (e *RealExecutor) ExecuteInteractiveWithEnv(name string, env []string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), env...)
	// Set process group to allow proper signal handling and terminal control
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: false,
	}
	return cmd.Run()
}

// MockExecutor is a mock executor for testing
type MockExecutor struct {
	Commands []string          // Records all executed commands
	Outputs  map[string]string // Map of command -> output
	Errors   map[string]error  // Map of command -> error
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		Commands: []string{},
		Outputs:  make(map[string]string),
		Errors:   make(map[string]error),
	}
}

func (e *MockExecutor) Execute(name string, args ...string) (string, error) {
	cmdStr := name + " " + strings.Join(args, " ")
	e.Commands = append(e.Commands, cmdStr)

	if err, ok := e.Errors[cmdStr]; ok {
		return "", err
	}

	if output, ok := e.Outputs[cmdStr]; ok {
		return output, nil
	}

	return "", nil
}

func (e *MockExecutor) ExecuteWithEnv(name string, env []string, args ...string) (string, error) {
	cmdStr := name + " " + strings.Join(args, " ")
	e.Commands = append(e.Commands, cmdStr)

	if err, ok := e.Errors[cmdStr]; ok {
		return "", err
	}

	if output, ok := e.Outputs[cmdStr]; ok {
		return output, nil
	}

	return "", nil
}

func (e *MockExecutor) SetOutput(cmd string, output string) {
	e.Outputs[cmd] = output
}

func (e *MockExecutor) SetError(cmd string, err error) {
	e.Errors[cmd] = err
}

func (e *MockExecutor) WasExecuted(cmd string) bool {
	for _, c := range e.Commands {
		if c == cmd {
			return true
		}
	}
	return false
}

func (e *MockExecutor) WasExecutedContaining(substring string) bool {
	for _, c := range e.Commands {
		if strings.Contains(c, substring) {
			return true
		}
	}
	return false
}

func (e *MockExecutor) ExecuteInteractive(name string, args ...string) error {
	cmdStr := name + " " + strings.Join(args, " ")
	e.Commands = append(e.Commands, cmdStr)

	if err, ok := e.Errors[cmdStr]; ok {
		return err
	}

	return nil
}

func (e *MockExecutor) ExecuteInteractiveWithEnv(name string, env []string, args ...string) error {
	cmdStr := name + " " + strings.Join(args, " ")
	e.Commands = append(e.Commands, cmdStr)

	if err, ok := e.Errors[cmdStr]; ok {
		return err
	}

	return nil
}

// RunCommand is a helper that uses the executor
func RunCommand(executor CommandExecutor, name string, args ...string) error {
	output, err := executor.Execute(name, args...)
	if err != nil {
		if output != "" {
			return fmt.Errorf("command failed: %s %v: %w\nOutput: %s", name, args, err, strings.TrimSpace(output))
		}
		return fmt.Errorf("command failed: %s %v: %w", name, args, err)
	}
	return nil
}

// RunCommandWithEnv runs a command with additional environment variables
func RunCommandWithEnv(executor CommandExecutor, env []string, name string, args ...string) error {
	output, err := executor.ExecuteWithEnv(name, env, args...)
	if err != nil {
		if output != "" {
			return fmt.Errorf("command failed: %s %v: %w\nOutput: %s", name, args, err, strings.TrimSpace(output))
		}
		return fmt.Errorf("command failed: %s %v: %w", name, args, err)
	}
	return nil
}

// RunInteractiveCommand runs a command with stdin/stdout/stderr connected to terminal
func RunInteractiveCommand(executor CommandExecutor, name string, args ...string) error {
	if err := executor.ExecuteInteractive(name, args...); err != nil {
		return fmt.Errorf("interactive command failed: %s %v: %w", name, args, err)
	}
	return nil
}
