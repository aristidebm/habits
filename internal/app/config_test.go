package app

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test DB config
	if config.DB.Path != StatePath() {
		t.Errorf("Expected DB path '%s', got '%s'", StatePath(), config.DB.Path)
	}

	// Test theme config
	if config.Theme.Name != "default" {
		t.Errorf("Expected theme name 'default', got '%s'", config.Theme.Name)
	}

	// Test view configs
	if config.Views.Weekly.Missed != "⛌" {
		t.Errorf("Expected weekly missed '⛌', got '%s'", config.Views.Weekly.Missed)
	}
	if config.Views.Weekly.Completed != "🗸" {
		t.Errorf("Expected weekly completed '🗸', got '%s'", config.Views.Weekly.Completed)
	}
	if config.Views.Weekly.Untracked != "○" {
		t.Errorf("Expected weekly untracked '○', got '%s'", config.Views.Weekly.Untracked)
	}

	// Monthly should match weekly
	if config.Views.Monthly.Missed != config.Views.Weekly.Missed {
		t.Errorf("Expected monthly missed to match weekly")
	}
	if config.Views.Monthly.Completed != config.Views.Weekly.Completed {
		t.Errorf("Expected monthly completed to match weekly")
	}
	if config.Views.Monthly.Untracked != config.Views.Weekly.Untracked {
		t.Errorf("Expected monthly untracked to match weekly")
	}
}
