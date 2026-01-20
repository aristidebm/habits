package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	db *sql.DB
}

func (s *Store) CreateHabit(ctx context.Context, name string, habitType HabitType, goal float64) (*Habit, error) {
	result, err := s.db.ExecContext(ctx,
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

	slog.Info("Created new habit", "id", id, "name", name, "type", habitType.String(), "goal", goal)
	return s.GetHabit(ctx, int(id))
}

func (s *Store) CreateHabitsBulk(ctx context.Context, habits []struct {
	Name      string
	HabitType HabitType
	Goal      float64
}) ([]Habit, error) {
	if len(habits) == 0 {
		return []Habit{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	valuePlaceholders := make([]string, len(habits))
	args := make([]any, 0, len(habits)*3)

	for i, h := range habits {
		valuePlaceholders[i] = "(?, ?, ?)"
		args = append(args, h.Name, h.HabitType.String(), h.Goal)
	}

	query := "INSERT INTO habits (name, kind, goal) VALUES " + strings.Join(valuePlaceholders, ",")
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		err = fmt.Errorf("failed to bulk insert habits: %w", err)
		return nil, err
	}

	nameArgs := make([]any, len(habits))
	namePlaceholders := make([]string, len(habits))
	for i, h := range habits {
		nameArgs[i] = h.Name
		namePlaceholders[i] = "?"
	}

	selectQuery := "SELECT id, name, kind, goal, created_at FROM habits WHERE name IN (" +
		strings.Join(namePlaceholders, ",") + ") ORDER BY id DESC LIMIT " +
		strconv.Itoa(len(habits))

	rows, err := tx.QueryContext(ctx, selectQuery, nameArgs...)
	if err != nil {
		err = fmt.Errorf("failed to query created habits: %w", err)
		return nil, err
	}
	defer rows.Close()

	var createdHabits []Habit
	for rows.Next() {
		var h Habit
		var kindStr string
		if err := rows.Scan(&h.ID, &h.Name, &kindStr, &h.Goal, &h.CreatedAt); err != nil {
			err = fmt.Errorf("failed to scan habit: %w", err)
			return nil, err
		}
		h.Type, _ = HabitTypeFromString(kindStr)
		createdHabits = append(createdHabits, h)
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return nil, err
	}

	return createdHabits, nil
}

func (s *Store) GetHabit(ctx context.Context, id int) (*Habit, error) {
	h := &Habit{}
	var kindStr string
	err := s.db.QueryRowContext(ctx,
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

func (s *Store) GetHabitByName(ctx context.Context, name string) (*Habit, error) {
	h := &Habit{}
	var kindStr string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, kind, goal, created_at FROM habits WHERE name = ?",
		name,
	).Scan(&h.ID, &h.Name, &kindStr, &h.Goal, &h.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get habit by name: %w", err)
	}

	ht, err := HabitTypeFromString(kindStr)
	if err != nil {
		return nil, err
	}
	h.Type = ht

	return h, nil
}

func (s *Store) UpdateHabit(ctx context.Context, id int, name string, habitType HabitType, goal float64) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE habits SET name = ?, kind = ?, goal = ? WHERE id = ?",
		name, habitType.String(), goal, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}
	return nil
}

func (s *Store) ListHabits(ctx context.Context) ([]Habit, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, kind, goal, created_at FROM habits ORDER BY id")
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

func (s *Store) DeleteHabit(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM habits WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	return nil
}

func (s *Store) UpsertEntry(ctx context.Context, habitID int, date time.Time, value float64) error {
	dateStr := date.Format("2006-01-02")
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO habit_entries (habit_id, value, date)
		VALUES (?, ?, ?)
		ON CONFLICT(habit_id, date) DO UPDATE SET value = excluded.value
	`, habitID, value, dateStr)

	if err != nil {
		return fmt.Errorf("failed to upsert entry: %w", err)
	}
	slog.Debug("Upserted habit entry", "habit_id", habitID, "date", dateStr, "value", value)
	return nil
}

func (s *Store) GetEntry(ctx context.Context, habitID int, date time.Time) (*HabitEntry, error) {
	dateStr := date.Format("2006-01-02")
	e := &HabitEntry{}
	var dateStrFromDB string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, habit_id, value, date, created_at
		FROM habit_entries
		WHERE habit_id = ? AND date = ?
	`, habitID, dateStr).Scan(&e.ID, &e.HabitID, &e.Value, &dateStrFromDB, &e.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	e.Date, err = time.Parse("2006-01-02", dateStrFromDB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date from database: %w", err)
	}

	return e, nil
}

func (s *Store) ListEntries(ctx context.Context, habitID int, startDate, endDate time.Time) ([]HabitEntry, error) {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	rows, err := s.db.QueryContext(ctx, `
		SELECT e.id, e.habit_id, e.value, e.date, e.created_at,
		       EXISTS(SELECT 1 FROM habit_notes n WHERE n.habit_entry_id = e.id) as has_note
		FROM habit_entries e
		WHERE e.habit_id = ? AND e.date >= ? AND e.date <= ?
		ORDER BY e.date
	`, habitID, startStr, endStr)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}
	defer rows.Close()

	var entries []HabitEntry
	for rows.Next() {
		var e HabitEntry
		var dateStr string
		var hasNoteInt int
		if err := rows.Scan(&e.ID, &e.HabitID, &e.Value, &dateStr, &e.CreatedAt, &hasNoteInt); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		e.Date, _ = time.Parse("2006-01-02", dateStr)
		e.HasNote = hasNoteInt == 1
		entries = append(entries, e)
	}

	return entries, nil
}

func (s *Store) GetHabits(ctx context.Context) []Habit {
	habits, err := s.ListHabits(ctx)
	if err != nil {
		return []Habit{}
	}
	return habits
}

func (s *Store) ListNotes(ctx context.Context, habitEntryID int) ([]HabitNote, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, habit_entry_id, note, created_at
		FROM habit_notes
		WHERE habit_entry_id = ?
		ORDER BY created_at
	`, habitEntryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	var notes []HabitNote
	for rows.Next() {
		var n HabitNote
		if err := rows.Scan(&n.ID, &n.HabitEntryID, &n.Note, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, n)
	}

	return notes, nil
}

func (s *Store) GetNote(ctx context.Context, habitEntryID int) (*HabitNote, error) {
	note := &HabitNote{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, habit_entry_id, note, created_at
		FROM habit_notes
		WHERE habit_entry_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`, habitEntryID).Scan(&note.ID, &note.HabitEntryID, &note.Note, &note.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No note found
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}

func (s *Store) UpsertNote(ctx context.Context, habitEntryID int, note string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO habit_notes (habit_entry_id, note)
		VALUES (?, ?)
		ON CONFLICT(habit_entry_id) DO UPDATE SET
			note = excluded.note,
			created_at = CURRENT_TIMESTAMP
	`, habitEntryID, note)

	if err != nil {
		return fmt.Errorf("failed to upsert note: %w", err)
	}
	slog.Debug("Upserted note", "habit_entry_id", habitEntryID, "note_length", len(note))
	return nil
}

func (s *Store) DeleteNote(ctx context.Context, habitEntryID int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM habit_notes WHERE habit_entry_id = ?", habitEntryID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}

func (s *Store) HasNote(ctx context.Context, habitEntryID int) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM habit_notes WHERE habit_entry_id = ?", habitEntryID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for notes: %w", err)
	}
	return count > 0, nil
}

func (s *Store) DeleteEntriesByIDs(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM habit_entries WHERE id IN (%s)", strings.Join(placeholders, ","))
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete entries by IDs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	slog.Info("Deleted entries by IDs", "count", rowsAffected, "ids", ids)
	return nil
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
