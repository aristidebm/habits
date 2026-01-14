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
