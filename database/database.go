package database

import (
	"database/sql"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	_ "modernc.org/sqlite"
)

type Handler interface {
	Insert(key, val string) error
	Get(key string) (string, error)
}

type SqliteHandler struct {
	db *sql.DB
}

type BadgerHandler struct {
	db *badger.DB
}

func (d BadgerHandler) Insert(key, val string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		if len(key) == 0 {
			return nil
		}
		return txn.Set([]byte(key), []byte(val))
	})
}

func (d BadgerHandler) Get(key string) (string, error) {
	var value []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("getting value: %w", err)
		}
		item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			fmt.Println(string(val))
			return nil
		})
		return nil
	})
	if err != nil {
		return string(value), err
	}
	return string(value), nil
}

func NewBadger(dataSourceName string) (*BadgerHandler, error) {
	db, err := badger.Open(badger.DefaultOptions(dataSourceName))
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	handler := &BadgerHandler{
		db: db,
	}
	return handler, nil
}

func NewSqlite(dataSouceName string) (*SqliteHandler, error) {
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

	return &SqliteHandler{
		db: db,
	}, nil
}

func (h *SqliteHandler) Insert(word, meaning string) error {
	insQuery := `INSERT INTO meaning 
	values(?, ?)
	ON CONFLICT(word) DO UPDATE SET meaning=excluded.meaning;`
	_, err := h.db.Exec(insQuery, word, meaning)
	if err != nil {
		return fmt.Errorf("inserting values: %w", err)
	}
	return nil
}

func (h *SqliteHandler) Get(word string) (string, error) {
	getQuery := "SELECT meaning FROM meaning WHERE word=?"
	row := h.db.QueryRow(getQuery, word)
	var meaning string
	err := row.Scan(&meaning)
	if err != nil {
		return "", fmt.Errorf("scanning row: %w", err)
	}
	return meaning, nil
}
