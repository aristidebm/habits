package app

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/pressly/goose/v3"
)

type App struct {
	db *sql.DB
}

func NewApp(dbPath string) (*App, error) {
	db, err := OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{db: db}, nil
}

func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *App) Migrate() error {
	return Migrate(a.db, "migrations")
}

func (a *App) CreateHabit(name string, habitType HabitType, goal float64) (*Habit, error) {
	result, err := a.db.Exec(
		"INSERT INTO habits (name, kind, goal) VALUES (?, ?, ?)",
		name, habitType.String(), goal,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get habit id: %w", err)
	}

	return a.GetHabit(int(id))
}

func (a *App) GetHabit(id int) (*Habit, error) {
	h := &Habit{}
	var kindStr string
	err := a.db.QueryRow(
		"SELECT id, name, kind, goal, created_at FROM habits WHERE id = ?",
		id,
	).Scan(&h.ID, &h.Name, &kindStr, &h.Goal, &h.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	ht, err := HabitTypeFromString(kindStr)
	if err != nil {
		return nil, err
	}
	h.Type = ht

	return h, nil
}

func (a *App) ListHabits() ([]Habit, error) {
	rows, err := a.db.Query("SELECT id, name, kind, goal, created_at FROM habits ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var h Habit
		var kindStr string
		if err := rows.Scan(&h.ID, &h.Name, &kindStr, &h.Goal, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		h.Type, _ = HabitTypeFromString(kindStr)
		habits = append(habits, h)
	}

	return habits, nil
}

func (a *App) DeleteHabit(id int) error {
	_, err := a.db.Exec("DELETE FROM habits WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	return nil
}

func (a *App) UpsertEntry(habitID int, date time.Time, value float64) error {
	dateStr := date.Format("2006-01-02")
	_, err := a.db.Exec(`
		INSERT INTO habit_entries (habit_id, value, date)
		VALUES (?, ?, ?)
		ON CONFLICT(habit_id, date) DO UPDATE SET value = excluded.value
	`, habitID, value, dateStr)

	if err != nil {
		return fmt.Errorf("failed to upsert entry: %w", err)
	}
	return nil
}

func (a *App) GetEntry(habitID int, date time.Time) (*HabitEntry, error) {
	dateStr := date.Format("2006-01-02")
	e := &HabitEntry{}
	err := a.db.QueryRow(`
		SELECT id, habit_id, value, date, created_at
		FROM habit_entries
		WHERE habit_id = ? AND date = ?
	`, habitID, dateStr).Scan(&e.ID, &e.HabitID, &e.Value, &e.Date, &e.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	return e, nil
}

func (a *App) ListEntries(habitID int, startDate, endDate time.Time) ([]HabitEntry, error) {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	rows, err := a.db.Query(`
		SELECT id, habit_id, value, date, created_at
		FROM habit_entries
		WHERE habit_id = ? AND date >= ? AND date <= ?
		ORDER BY date
	`, habitID, startStr, endStr)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}
	defer rows.Close()

	var entries []HabitEntry
	for rows.Next() {
		var e HabitEntry
		var dateStr string
		if err := rows.Scan(&e.ID, &e.HabitID, &e.Value, &dateStr, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		e.Date, _ = time.Parse("2006-01-02", dateStr)
		entries = append(entries, e)
	}

	return entries, nil
}

func (a *App) GetHabits() []Habit {
	habits, err := a.ListHabits()
	if err != nil {
		return []Habit{}
	}
	return habits
}

func (a *App) GetEntryByID(habitID int, date time.Time) (HabitEntry, bool) {
	entry, err := a.GetEntry(habitID, date)
	if err != nil || entry == nil {
		return HabitEntry{}, false
	}
	return *entry, true
}
