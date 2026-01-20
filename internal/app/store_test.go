package app

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	if err := Migrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestStoreCreateHabit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := &Store{db: db}
	ctx := context.Background()

	habit, err := store.CreateHabit(ctx, "Test Habit", HabitTypeBit, 1.0)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	if habit.Name != "Test Habit" {
		t.Errorf("Expected habit name 'Test Habit', got '%s'", habit.Name)
	}
	if habit.Type != HabitTypeBit {
		t.Errorf("Expected habit type Bit, got %v", habit.Type)
	}
	if habit.Goal != 1.0 {
		t.Errorf("Expected goal 1.0, got %f", habit.Goal)
	}
}

func TestStoreGetHabit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := &Store{db: db}
	ctx := context.Background()

	// Create habit
	habit, err := store.CreateHabit(ctx, "Test Habit", HabitTypeBit, 1.0)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// Get habit
	retrieved, err := store.GetHabit(ctx, habit.ID)
	if err != nil {
		t.Fatalf("Failed to get habit: %v", err)
	}

	if retrieved.ID != habit.ID {
		t.Errorf("Expected habit ID %d, got %d", habit.ID, retrieved.ID)
	}
	if retrieved.Name != habit.Name {
		t.Errorf("Expected habit name '%s', got '%s'", habit.Name, retrieved.Name)
	}
}

func TestStoreUpsertEntry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := &Store{db: db}
	ctx := context.Background()

	// Create habit
	habit, err := store.CreateHabit(ctx, "Test Habit", HabitTypeBit, 1.0)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// Upsert entry
	date := time.Now()
	err = store.UpsertEntry(ctx, habit.ID, date, 1.0)
	if err != nil {
		t.Fatalf("Failed to upsert entry: %v", err)
	}

	// Get entry
	entry, err := store.GetEntry(ctx, habit.ID, date)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if entry.Value != 1.0 {
		t.Errorf("Expected entry value 1.0, got %f", entry.Value)
	}
}

func TestStoreUpsertNote(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := &Store{db: db}
	ctx := context.Background()

	// Create habit
	habit, err := store.CreateHabit(ctx, "Test Habit", HabitTypeBit, 1.0)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// Create entry
	date := time.Now()
	err = store.UpsertEntry(ctx, habit.ID, date, 1.0)
	if err != nil {
		t.Fatalf("Failed to create entry: %v", err)
	}

	// Get entry to find entry ID
	entry, err := store.GetEntry(ctx, habit.ID, date)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	// Upsert note
	err = store.UpsertNote(ctx, entry.ID, "Test note")
	if err != nil {
		t.Fatalf("Failed to upsert note: %v", err)
	}

	// Get note
	note, err := store.GetNote(ctx, entry.ID)
	if err != nil {
		t.Fatalf("Failed to get note: %v", err)
	}

	if note.Note != "Test note" {
		t.Errorf("Expected note 'Test note', got '%s'", note.Note)
	}
}
