package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func Postgres(config PostgresConfig) (DB, error) {
	const op = "repositories.friends.PostgresRepository"
	psqlConnection := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	)
	db, err := sql.Open("postgres", psqlConnection)
	if err != nil {
		log.Printf("%s: open db failed err: %v", op, err)
		return nil, err
	}
	return db, nil
}
