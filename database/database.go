package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Handler struct {
	db *sql.DB
}

func New(dataSouceName string) (*Handler, error) {
	db, err := sql.Open("sqlite", dataSouceName)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS meaning(
		word varchar primary key, 
		meaning varchar 
	)`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) Insert(word, meaning string) error {
	insQuery := `INSERT INTO meaning 
	values(?, ?)
	ON CONFLICT(word) DO UPDATE SET meaning=excluded.meaning;`
	_, err := h.db.Exec(insQuery, word, meaning)
	if err != nil {
		return fmt.Errorf("inserting values: %w", err)
	}
	return nil
}

func (h *Handler) Get(word string) (string, error) {
	getQuery := "SELECT meaning FROM meaning WHERE word=?"
	row := h.db.QueryRow(getQuery, word)
	var meaning string
	err := row.Scan(&meaning)
	if err != nil {
		return "", fmt.Errorf("scanning row: %w", err)
	}
	return meaning, nil
}
