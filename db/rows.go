package db

import "github.com/jackc/pgx/v5"

type Rows struct {
	rows pgx.Rows
}

func (r Rows) Close() {
	r.rows.Close()
}

func (r Rows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r Rows) Next() bool {
	return r.rows.Next()
}

type Row struct {
	row pgx.Row
}

func (r Row) Close() {

}

func (r Row) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}
