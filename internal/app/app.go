package app

import (
	_ "github.com/pressly/goose/v3"
)

type App struct {
	*Store
}

func NewApp(dbPath string) (*App, error) {
	db, err := OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{
		Store: &Store{db: db},
	}, nil
}

func (a *App) Migrate() error {
	return Migrate(a.Store.db, "migrations")
}

func (a *App) Close() error {
	return a.Store.Close()
}
