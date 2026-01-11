package command

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommandFunc is the function signature for command handlers
type CommandFunc func(args []string) tea.Cmd

// Command represents a registerable command
type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     CommandFunc
}

// CommandLine is a vim-style command line component
type CommandLine struct {
	input       textinput.Model
	commands    map[string]Command
	visible     bool
	lastError   string
	lastSuccess string
	width       int

	// Styles
	promptStyle  lipgloss.Style
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style
}

// NewCommandLine creates a new command line component
func NewCommandLine() *CommandLine {
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Prompt = ":"
	ti.CharLimit = 256

	cl := &CommandLine{
		input:    ti,
		commands: make(map[string]Command),
		width:    80,

		promptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true),
		successStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")),
	}

	// Register built-in help command
	cl.RegisterCommand(Command{
		Name:        "help",
		Description: "Show available commands",
		Usage:       "help [command]",
		Handler:     cl.handleHelp,
	})

	return cl
}

// RegisterCommand registers a new command
func (c *CommandLine) RegisterCommand(cmd Command) {
	c.commands[cmd.Name] = cmd
}

// Show shows the command line
func (c *CommandLine) Show() tea.Cmd {
	c.visible = true
	c.lastError = ""
	c.lastSuccess = ""
	c.input.Focus()
	c.input.SetValue("")
	return textinput.Blink
}

// Hide hides the command line
func (c *CommandLine) Hide() {
	c.visible = false
	c.input.Blur()
}

// IsVisible returns whether the command line is visible
func (c *CommandLine) IsVisible() bool {
	return c.visible
}

// SetWidth sets the width of the command line
func (c *CommandLine) SetWidth(width int) {
	c.width = width
	c.input.Width = width - 2
}

// Update handles messages
func (c *CommandLine) Update(msg tea.Msg) tea.Cmd {
	if !c.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			c.Hide()
			return nil
		case tea.KeyEnter:
			return c.executeCommand()
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return cmd
}

// View renders the command line
func (c *CommandLine) View() string {
	if !c.visible {
		// Show status messages even when command line is hidden
		if c.lastError != "" {
			return c.errorStyle.Render(c.lastError)
		}
		if c.lastSuccess != "" {
			return c.successStyle.Render(c.lastSuccess)
		}
		return ""
	}

	return c.input.View()
}

// executeCommand parses and executes the entered command
func (c *CommandLine) executeCommand() tea.Cmd {
	input := strings.TrimSpace(c.input.Value())
	if input == "" {
		c.Hide()
		return nil
	}

	// Parse command and arguments
	parts := strings.Fields(input)
	cmdName := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	// Find and execute command
	cmd, exists := c.commands[cmdName]
	if !exists {
		c.lastError = fmt.Sprintf("Unknown command: %s", cmdName)
		c.Hide()
		return nil
	}

	c.Hide()
	return cmd.Handler(args)
}

// handleHelp handles the help command
func (c *CommandLine) handleHelp(args []string) tea.Cmd {
	if len(args) > 0 {
		// Show help for specific command
		cmdName := args[0]
		cmd, exists := c.commands[cmdName]
		if !exists {
			c.lastError = fmt.Sprintf("Unknown command: %s", cmdName)
			return nil
		}
		c.lastSuccess = fmt.Sprintf("%s: %s\nUsage: %s", cmd.Name, cmd.Description, cmd.Usage)
		return nil
	}

	// Show all commands
	var help strings.Builder
	help.WriteString("Available commands:\n")
	for _, cmd := range c.commands {
		help.WriteString(fmt.Sprintf("  %-15s %s\n", cmd.Name, cmd.Description))
	}
	c.lastSuccess = help.String()
	return nil
}

// SetError sets an error message to display
func (c *CommandLine) SetError(msg string) {
	c.lastError = msg
	c.lastSuccess = ""
}

// SetSuccess sets a success message to display
func (c *CommandLine) SetSuccess(msg string) {
	c.lastSuccess = msg
	c.lastError = ""
}

// ClearMessages clears all messages
func (c *CommandLine) ClearMessages() {
	c.lastError = ""
	c.lastSuccess = ""
}
