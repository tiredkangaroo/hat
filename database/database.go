package database

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/tiredkangaroo/hat/proxy/config"
)

//go:embed initialize.sql
var initialize_sql string

type DB struct {
	conn *pgx.Conn
}

func (db *DB) initialize() error {
	var err error
	db.conn, err = pgx.Connect(context.Background(), config.DefaultConfig.Database.PostgresURL)
	if err != nil {
		return fmt.Errorf("pgx connect: %w", err)
	}
	_, err = db.conn.Exec(context.Background(), initialize_sql)
	if err != nil {
		return fmt.Errorf("create tables: %w", err)
	}
	return nil
}

// GetDB returns a database instance. It requires that configuration be initialized.
func GetDB() (*DB, error) {
	db := &DB{}
	if err := db.initialize(); err != nil {
		return nil, err
	}
	return db, nil
}
