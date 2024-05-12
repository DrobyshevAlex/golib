package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnString string
	Tracer     pgx.QueryTracer
}

type Conn struct {
	*pgxpool.Pool
}

func NewConn(ctx context.Context, config Config) (*Conn, error) {
	cnf, err := pgxpool.ParseConfig(config.ConnString)
	if err != nil {
		return nil, err
	}

	cnf.ConnConfig.Tracer = config.Tracer

	poll, err := pgxpool.NewWithConfig(ctx, cnf)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Pool: poll,
	}, nil
}

func (c Conn) Close(ctx context.Context) {
	c.Pool.Close()
}

func (c Conn) Query(ctx context.Context, sql string, args ...interface{}) (*Rows, error) {
	rows, err := c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

func (c Conn) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return Row{row: c.Pool.QueryRow(ctx, sql, args...)}
}
