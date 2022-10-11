package repository

import (
	"context"
	"github.com/jackc/pgx"
)

type PostgresConfig struct {
	Host     string
	Port     uint16
	Username string
	Password string
	DBName   string
}

func NewPostgresDB(cfg PostgresConfig) (*pgx.Conn, error) {

	db, err := pgx.Connect(pgx.ConnConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Password: cfg.Password,
		User:     cfg.Username,
		Database: cfg.DBName,
	})
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
