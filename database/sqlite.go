package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	db *sqlx.DB
}

func NewSqlite(dataSouceName string) (*SqliteHandler, error) {
	db, err := sqlx.Connect("sqlite3", dataSouceName)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, fmt.Errorf("setting journal mode: %w", err)
	}

	_, err = db.Exec("PRAGMA synchronous = normal;")
	if err != nil {
		return nil, fmt.Errorf("setting synchronous: %w", err)
	}

	_, err = db.Exec("PRAGMA busy_timeout = 3000;")
	if err != nil {
		return nil, fmt.Errorf("setting busy timeout: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS meaning(
		word VARCHAR NOT NULL PRIMARY KEY, 
		meaning VARCHAR NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	return &SqliteHandler{
		db: db,
	}, nil
}

func (h *SqliteHandler) Insert(word, meaning string) error {
	const insQuery = "INSERT OR IGNORE INTO meaning (word, meaning) VALUES (?, ?);"

	_, err := h.db.Exec(insQuery, word, meaning)

	return err
}

func (h *SqliteHandler) Get(word string) (string, error) {
	const getQuery = "SELECT meaning FROM meaning WHERE word=?;"

	var meaning string

	err := h.db.QueryRow(getQuery, word).Scan(&meaning)
	if err != nil {
		return "", fmt.Errorf("getting meaning: %w", err)
	}

	return meaning, nil
}

func (h *SqliteHandler) Close() error {
	return h.db.Close()
}
