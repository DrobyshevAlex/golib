package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type dbTxKey struct{}

func With(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, dbTxKey{}, tx)
}

func FromContext(ctx context.Context) (pgx.Tx, error) {
	if tx, ok := ctx.Value(dbTxKey{}).(pgx.Tx); ok {
		return tx, nil
	}
	return nil, errors.New("tx not found in context")
}
