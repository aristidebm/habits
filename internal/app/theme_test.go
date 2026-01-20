package app

import (
	"testing"
)

func TestGetBuiltInTheme(t *testing.T) {
	tm := NewThemeManager()

	// Test default theme
	theme, err := tm.LoadTheme("default")
	if err != nil {
		t.Fatalf("Failed to load default theme: %v", err)
	}
	if theme.Name != "Default" {
		t.Errorf("Expected theme name 'Default', got '%s'", theme.Name)
	}
	if theme.Colors.Header != "#87CEEB" {
		t.Errorf("Expected header color '#87CEEB', got '%s'", theme.Colors.Header)
	}

	// Test dark theme
	theme, err = tm.LoadTheme("dark")
	if err != nil {
		t.Fatalf("Failed to load dark theme: %v", err)
	}
	if theme.Name != "Dark" {
		t.Errorf("Expected theme name 'Dark', got '%s'", theme.Name)
	}

	// Test light theme
	theme, err = tm.LoadTheme("light")
	if err != nil {
		t.Fatalf("Failed to load light theme: %v", err)
	}
	if theme.Name != "Light" {
		t.Errorf("Expected theme name 'Light', got '%s'", theme.Name)
	}

	// Test mono theme
	theme, err = tm.LoadTheme("mono")
	if err != nil {
		t.Fatalf("Failed to load mono theme: %v", err)
	}
	if theme.Name != "Monochrome" {
		t.Errorf("Expected theme name 'Monochrome', got '%s'", theme.Name)
	}

	// Test invalid theme
	_, err = tm.LoadTheme("invalid")
	if err == nil {
		t.Error("Expected error for invalid theme, got nil")
	}
}

func TestThemeManagerCaching(t *testing.T) {
	tm := NewThemeManager()

	// Load theme twice
	theme1, err := tm.LoadTheme("default")
	if err != nil {
		t.Fatalf("Failed to load default theme first time: %v", err)
	}

	theme2, err := tm.LoadTheme("default")
	if err != nil {
		t.Fatalf("Failed to load default theme second time: %v", err)
	}

	// Should be the same instance
	if theme1 != theme2 {
		t.Error("Expected cached theme to be the same instance")
	}
}
