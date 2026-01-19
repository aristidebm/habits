-- +goose Up
CREATE TABLE IF NOT EXISTS habits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK(kind IN ('bit', 'count', 'float')),
    goal REAL NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS habit_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_id INTEGER NOT NULL,
    value REAL NOT NULL DEFAULT 0,
    date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
    UNIQUE(habit_id, date)
);

CREATE TABLE IF NOT EXISTS habit_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habit_entry_id INTEGER NOT NULL,
    note TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_entry_id) REFERENCES habit_entries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_habit_entries_date ON habit_entries(date);
-- +goose Down
DROP INDEX IF EXISTS idx_habit_entries_date;
DROP TABLE IF EXISTS habit_notes;
DROP TABLE IF EXISTS habit_entries;
DROP TABLE IF EXISTS habits;
