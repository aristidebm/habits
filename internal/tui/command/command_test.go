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

func TestCommandLineShowHide(t *testing.T) {
	cl := NewCommandLine()

	// Initially not visible
	if cl.IsVisible() {
		t.Error("Expected command line to be hidden initially")
	}

	// Show command line
	cmd := cl.Show()
	if cmd == nil {
		t.Error("Expected Show() to return a command")
	}

	if !cl.IsVisible() {
		t.Error("Expected command line to be visible after Show()")
	}

	// Hide command line
	cl.Hide()
	if cl.IsVisible() {
		t.Error("Expected command line to be hidden after Hide()")
	}
}

func TestCommandLineQuitCommands(t *testing.T) {
	cl := NewCommandLine()

	// Register quit command with alias
	cl.RegisterCommand(Command{
		Name:    "quit",
		Aliases: []string{"q"},
		Handler: func(args []string) (Result, tea.Cmd) {
			return Quit(), tea.Quit
		},
	})

	// Test quit command execution
	cl.input.SetValue("quit")
	cl.input.SetCursor(4)
	cmd := cl.executeCommand()

	if cmd == nil {
		t.Error("Expected executeCommand to return tea.Quit for quit command")
	}

	// Test q alias
	cl.input.SetValue("q")
	cl.input.SetCursor(1)
	cmd = cl.executeCommand()

	if cmd == nil {
		t.Error("Expected executeCommand to return tea.Quit for q alias")
	}
}

func TestCommandLineUnknownCommand(t *testing.T) {
	cl := NewCommandLine()

	// Try to execute unknown command
	cl.input.SetValue("unknown-command")
	cl.input.SetCursor(15)
	cmd := cl.executeCommand()

	if cmd != nil {
		t.Error("Expected executeCommand to return nil for unknown command")
	}

	// Check that error message was set
	if cl.lastError == "" {
		t.Error("Expected error message to be set for unknown command")
	}

	if cl.lastError != "Unknown command: unknown-command" {
		t.Errorf("Expected specific error message, got: %s", cl.lastError)
	}
}

func TestCommandLineEmptyInput(t *testing.T) {
	cl := NewCommandLine()
	cl.Show() // Make it visible

	// Try to execute empty command
	cl.input.SetValue("")
	cl.input.SetCursor(0)
	cmd := cl.executeCommand()

	if cmd != nil {
		t.Error("Expected executeCommand to return nil for empty input")
	}

	// Should hide the command line
	if cl.IsVisible() {
		t.Error("Expected command line to be hidden after empty input")
	}
}

func TestCommandLineWithArguments(t *testing.T) {
	cl := NewCommandLine()

	executed := false
	var receivedArgs []string

	cl.RegisterCommand(Command{
		Name: "test-args",
		Handler: func(args []string) (Result, tea.Cmd) {
			executed = true
			receivedArgs = args
			return Success("args received"), nil
		},
	})

	// Execute command with arguments
	cl.input.SetValue("test-args arg1 arg2 arg3")
	cl.input.SetCursor(23)
	cl.executeCommand()

	if !executed {
		t.Error("Command should have been executed")
	}

	if len(receivedArgs) != 3 {
		t.Errorf("Expected 3 arguments, got %d", len(receivedArgs))
	}

	expectedArgs := []string{"arg1", "arg2", "arg3"}
	for i, expected := range expectedArgs {
		if i >= len(receivedArgs) || receivedArgs[i] != expected {
			t.Errorf("Expected arg[%d] to be '%s', got '%s'", i, expected, receivedArgs[i])
		}
	}
}

func TestCommandLineSuccessMessage(t *testing.T) {
	cl := NewCommandLine()

	cl.RegisterCommand(Command{
		Name: "success-test",
		Handler: func(args []string) (Result, tea.Cmd) {
			return Success("Operation successful"), nil
		},
	})

	// Execute command
	cl.input.SetValue("success-test")
	cl.input.SetCursor(12)
	cl.executeCommand()

	// Check success message
	if cl.lastSuccess != "Operation successful" {
		t.Errorf("Expected success message 'Operation successful', got '%s'", cl.lastSuccess)
	}

	if cl.lastError != "" {
		t.Errorf("Expected no error message, got '%s'", cl.lastError)
	}
}

func TestCommandLineErrorMessage(t *testing.T) {
	cl := NewCommandLine()

	cl.RegisterCommand(Command{
		Name: "error-test",
		Handler: func(args []string) (Result, tea.Cmd) {
			return Error("Something went wrong"), nil
		},
	})

	// Execute command
	cl.input.SetValue("error-test")
	cl.input.SetCursor(10)
	cl.executeCommand()

	// Check error message
	if cl.lastError != "Something went wrong" {
		t.Errorf("Expected error message 'Something went wrong', got '%s'", cl.lastError)
	}

	if cl.lastSuccess != "" {
		t.Errorf("Expected no success message, got '%s'", cl.lastSuccess)
	}
}

func TestCommandLineAliasRegistration(t *testing.T) {
	cl := NewCommandLine()

	// Register command with aliases
	// Register command with aliases
	cl.RegisterCommand(Command{
		Name:    "test-cmd",
		Aliases: []string{"tc", "test"},
		Handler: func(args []string) (Result, tea.Cmd) {
			return Success("alias worked"), nil
		},
	})

	// Test main command
	if _, exists := cl.commands["test-cmd"]; !exists {
		t.Error("Main command not registered")
	}

	// Test aliases
	if _, exists := cl.commands["tc"]; !exists {
		t.Error("Alias 'tc' not registered")
	}

	if _, exists := cl.commands["test"]; !exists {
		t.Error("Alias 'test' not registered")
	}

	// Test that aliases point to same command
	mainCmd := cl.commands["test-cmd"]
	aliasCmd := cl.commands["tc"]

	if mainCmd.Name != aliasCmd.Name {
		t.Error("Alias should point to same command")
	}
}
