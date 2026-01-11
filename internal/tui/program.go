package tui

import (
	"fmt"
	"os"
	"time"

	"example.com/habits/internal/tui/calendar"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	calendar *calendar.Calendar
}

func initialModel() model {
	// Create some sample habits
	habits := []calendar.Habit{
		{Name: "Morning Run", Type: calendar.HabitTypeBit},
		{Name: "Read Pages", Type: calendar.HabitTypeCount},
		{Name: "Water (L)", Type: calendar.HabitTypeFloat},
		{Name: "Meditation", Type: calendar.HabitTypeBit},
	}

	calendar := calendar.NewCalendar(habits)

	// Add sample data for the past 2 weeks and current month
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	// Fill data from start of month to today
	for d := monthStart; !d.After(now); d = d.AddDate(0, 0, 1) {
		dayOfMonth := d.Day()
		
		// Morning Run - alternating pattern
		if dayOfMonth%2 == 0 {
			calendar.SetEntry("Morning Run", d, true, "")
		} else {
			calendar.SetEntry("Morning Run", d, false, "-")
		}

		// Read Pages - varying numbers
		if dayOfMonth%3 != 0 {
			value := fmt.Sprintf("%d", 10+(dayOfMonth%20))
			calendar.SetEntry("Read Pages", d, false, value)
		} else {
			calendar.SetEntry("Read Pages", d, false, "-")
		}

		// Water - float values
		if dayOfMonth%4 != 1 {
			value := fmt.Sprintf("%.1f", 2.0+(float64(dayOfMonth%10)*0.1))
			calendar.SetEntry("Water (L)", d, false, value)
		} else {
			calendar.SetEntry("Water (L)", d, false, "-")
		}

		// Meditation - mostly consistent
		if dayOfMonth%5 != 0 {
			calendar.SetEntry("Meditation", d, true, "")
		} else {
			calendar.SetEntry("Meditation", d, false, "-")
		}
	}

	return model{
		calendar: calendar,
	}
}

func (m model) Init() tea.Cmd {
	return m.calendar.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Update calendar
	updated, cmd := m.calendar.Update(msg)
	m.calendar = updated.(*calendar.Calendar)
	return m, cmd
}

func (m model) View() string {
	help := "\nKeys: [h/l] prev/next day | [H/L] week/month jump | [j/k] habit (weekly) | [n/p] habit | [TAB] switch view | [t] today | [q] quit\n"
	return m.calendar.View() + help
}

func Execute() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
