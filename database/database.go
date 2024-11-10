package database

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Handler interface {
	Insert(key, val string) error
	Get(key string) (string, error)
	Close() error
}

type SqliteHandler struct {
	db *sqlx.DB
}

type BadgerHandler struct {
	db *badger.DB
}

func (d BadgerHandler) Close() error {
	return d.db.Close()
}

func (d BadgerHandler) Insert(key, val string) error {
	prevVal, err := d.Get(key)
	var newVal string
	// previous value does not exist
	if err != nil {
		newVal = val
	} else {
		newVal = prevVal + "," + val
	}

	return d.db.Update(func(txn *badger.Txn) error {
		if len(key) == 0 {
			return nil
		}
		return txn.Set([]byte(key), []byte(newVal))
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
	row := h.db.QueryRow(getQuery, word)
	var meaning string
	err := row.Scan(&meaning)
	if err != nil {
		return "", fmt.Errorf("scanning row: %w", err)
	}
	return meaning, nil
}

func (h *SqliteHandler) Close() error {
	return h.db.Close()
}
