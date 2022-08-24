package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	Tx(ctx context.Context, cb func(*Queries, *interface{}) (interface{}, error)) (interface{}, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SQLStore) Tx(ctx context.Context, cb func(*Queries, *interface{}) (interface{}, error)) (interface{}, error) {
	var results interface{}

	err := store.execTx(ctx, func(q *Queries) error {
		return nil
	})

	return results, err
}
