package app

import (
	"fmt"
	"time"
)

type HabitType int

const (
	HabitTypeBit HabitType = iota
	HabitTypeCount
	HabitTypeFloat
)

func (ht HabitType) String() string {
	switch ht {
	case HabitTypeBit:
		return "bit"
	case HabitTypeCount:
		return "count"
	case HabitTypeFloat:
		return "float"
	default:
		return ""
	}
}

func HabitTypeFromString(s string) (HabitType, error) {
	switch s {
	case "bit":
		return HabitTypeBit, nil
	case "count":
		return HabitTypeCount, nil
	case "float":
		return HabitTypeFloat, nil
	default:
		return 0, fmt.Errorf("invalid habit type: %s", s)
	}
}

type Habit struct {
	ID        int
	Name      string
	Type      HabitType
	Goal      float64
	CreatedAt time.Time
}

type HabitEntry struct {
	ID        int
	HabitID   int
	Value     float64
	Date      time.Time
	CreatedAt time.Time
}

type HabitNote struct {
	ID           int
	HabitEntryID int
	Note         string
	CreatedAt    time.Time
}

// Export structures for JSON output

type ExportNote struct {
	Note string `json:"note"`
}

type ExportEntry struct {
	ID    int          `json:"id"`
	Date  string       `json:"date"`
	Value float64      `json:"value"`
	Notes []ExportNote `json:"notes"`
}

type ExportHabit struct {
	ID      int           `json:"id"`
	Name    string        `json:"name"`
	Goal    float64       `json:"goal"`
	Kind    string        `json:"kind"`
	Entries []ExportEntry `json:"entries"`
}
