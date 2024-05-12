package db

import (
	"context"
)

type DBRole string

const (
	RoleMaster = "master"
	RoleRead   = "read"
)

type DB struct {
	conn map[DBRole]*Conn
}

func Open(
	ctx context.Context,
	configs map[DBRole]Config,
) (db *DB, err error) {
	db = &DB{
		make(map[DBRole]*Conn, len(configs)),
	}

	defer func() {
		if err != nil {
			db.Close(ctx)
		}
	}()

	for role, cnf := range configs {
		conn, err := NewConn(ctx, cnf)
		if err != nil {
			return nil, err
		}
		db.conn[role] = conn
	}
	return db, nil
}

func (db DB) Close(ctx context.Context) (err error) {
	for _, conn := range db.conn {
		conn.Close(ctx)
	}
	return
}

func (db DB) For(role DBRole) *Conn {
	if c, ok := db.conn[role]; ok {
		return c
	}
	return nil
}

func (db DB) Write() *Conn {
	if c, ok := db.conn[RoleMaster]; ok {
		return c
	}
	return nil
}

func (db DB) Read() *Conn {
	if c, ok := db.conn[RoleRead]; ok {
		return c
	}
	return db.Write()
}
