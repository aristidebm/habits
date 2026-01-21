package command

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandFuncSignature(t *testing.T) {
	// Test that CommandFunc returns both Result and Cmd
	cmd := Command{
		Name: "test",
		Handler: func(args []string) (Result, tea.Cmd) {
			return Success("test"), nil
		},
	}

	result, cmdResult := cmd.Handler([]string{})
	if result.Type != ResultSuccess {
		t.Errorf("Expected success result, got %v", result.Type)
	}
	if result.Message != "test" {
		t.Errorf("Expected message 'test', got '%s'", result.Message)
	}
	if cmdResult != nil {
		t.Errorf("Expected nil command, got %v", cmdResult)
	}
}

func TestCommandLineExecuteWithCommands(t *testing.T) {
	cl := NewCommandLine()
	executed := false

	// Register a command that returns a command
	cl.RegisterCommand(Command{
		Name: "test-cmd",
		Handler: func(args []string) (Result, tea.Cmd) {
			executed = true
			return Success("command executed"), tea.Quit
		},
	})

	// Mock the input to trigger command execution (without the : prefix)
	cl.input.SetValue("test-cmd")
	cl.input.SetCursor(len("test-cmd"))

	// Execute the command
	cmd := cl.executeCommand()

	// Check that the command was executed
	if !executed {
		t.Error("Command handler was not executed")
	}

	// Check that a command was returned (tea.Quit in this case)
	if cmd == nil {
		t.Error("Expected a command to be returned, got nil")
	}
}

func TestResultTypes(t *testing.T) {
	// Test Success result
	success := Success("operation successful")
	if success.Type != ResultSuccess {
		t.Errorf("Expected ResultSuccess, got %v", success.Type)
	}
	if success.Message != "operation successful" {
		t.Errorf("Expected message 'operation successful', got '%s'", success.Message)
	}
	if success.Cmd != nil {
		t.Errorf("Expected nil Cmd for success, got %v", success.Cmd)
	}

	// Test Error result
	err := Error("something went wrong")
	if err.Type != ResultError {
		t.Errorf("Expected ResultError, got %v", err.Type)
	}
	if err.Message != "something went wrong" {
		t.Errorf("Expected message 'something went wrong', got '%s'", err.Message)
	}
	if err.Cmd != nil {
		t.Errorf("Expected nil Cmd for error, got %v", err.Cmd)
	}

	// Test Quit result
	quit := Quit()
	if quit.Type != ResultQuit {
		t.Errorf("Expected ResultQuit, got %v", quit.Type)
	}
	if quit.Cmd == nil {
		t.Error("Expected non-nil Cmd for quit")
	}
}

func TestCommandRegistration(t *testing.T) {
	cl := NewCommandLine()

	// Test registering a command
	initialCount := len(cl.commands)
	cl.RegisterCommand(Command{
		Name: "new-command",
		Handler: func(args []string) (Result, tea.Cmd) {
			return Success("new command"), nil
		},
	})

	if len(cl.commands) != initialCount+1 {
		t.Errorf("Expected command count to increase by 1, got %d", len(cl.commands))
	}

	// Test that the command is registered
	if _, exists := cl.commands["new-command"]; !exists {
		t.Error("Expected 'new-command' to be registered")
	}
}
