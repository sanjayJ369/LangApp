package database

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type BadgerHandler struct {
	db *badger.DB
}

func NewBadger(dataSourceName string) (*BadgerHandler, error) {
	db, err := badger.Open(
		badger.DefaultOptions(dataSourceName).
			WithLoggingLevel(badger.WARNING),
	)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	handler := &BadgerHandler{
		db: db,
	}
	return handler, nil
}

func (h BadgerHandler) Insert(key, val string) error {
	return h.db.Update(func(txn *badger.Txn) error {
		if len(key) == 0 {
			return nil
		}
		return txn.Set([]byte(key), []byte(val))
	})
}

func (h BadgerHandler) Get(key string) (string, error) {
	var value []byte
	err := h.db.View(func(txn *badger.Txn) error {
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

func (d BadgerHandler) Close() error {
	return d.db.Close()
}
