package database

import (
	"fmt"
	"strings"

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
		return nil, fmt.Errorf("setting WAL mode: %s", err)
	}

	_, err = db.Exec("PRAGMA wal_autocheckpoint=1000;;")
	if err != nil {
		return nil, fmt.Errorf("setting WAL mode: %s", err)
	}

	_, err = db.Exec("PRAGMA timeout = 10000;")
	if err != nil {
		return nil, fmt.Errorf("setting WAL mode: %s", err)
	}

	query := `CREATE TABLE IF NOT EXISTS meaning(
		word varchar , 
		meaning varchar ,
		primary key (word, meaning)
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
	// retry failed insertes
	for i := 0; i < 3; i++ {
		tx, err := h.db.Begin()
		if err != nil {
			return fmt.Errorf("creating transcation")
		}

		insQuery := `INSERT OR IGNORE INTO meaning (word, meaning)
			VALUES (?, ?) `
		_, err = tx.Exec(insQuery, word, meaning)
		if err != nil {
			tx.Rollback()
			if err.Error() == "database is locked" {
				continue
			}
			return fmt.Errorf("inserting values: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
	}
	return nil
}

func (h *SqliteHandler) Get(word string) (string, error) {
	getQuery := "SELECT meaning FROM meaning WHERE word=?"
	rows, err := h.db.Query(getQuery, word)
	defer rows.Close()
	if err != nil {
		return "", fmt.Errorf("getting meaning: %w", err)
	}
	senses := make([]string, 0)
	var meaning string
	for rows.Next() {
		err := rows.Scan(&meaning)
		if err != nil {
			return "", fmt.Errorf("scanning row: %w", err)
		}
		senses = append(senses, meaning)
	}
	meaning = strings.Join(senses, ",")
	return meaning, nil
}

func (h *SqliteHandler) Close() error {
	return h.db.Close()
}
