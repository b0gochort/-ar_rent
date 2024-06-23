package postgres

import (
	"fmt"
	"github.com/b0gochort/car-rent/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(config model.Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Database, config.Postgres.Password, config.Postgres.SSL)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
