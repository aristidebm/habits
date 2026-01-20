-- +goose Up
-- Remove duplicate notes, keeping the most recent one for each habit_entry_id
DELETE FROM habit_notes
WHERE id NOT IN (
    SELECT MAX(id)
    FROM habit_notes
    GROUP BY habit_entry_id
);

-- Add unique constraint on habit_entry_id
CREATE UNIQUE INDEX idx_habit_notes_entry_id ON habit_notes(habit_entry_id);

-- +goose Down
DROP INDEX IF EXISTS idx_habit_notes_entry_id;
