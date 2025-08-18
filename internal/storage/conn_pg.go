package storage

import (
	"RSSHub/internal/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

type CLIRepo interface{}

func Connect(cfgDB models.DB) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfgDB.PostgresHost,
		cfgDB.PostgresPort,
		cfgDB.PostgresUser,
		cfgDB.PostgresPass,
		cfgDB.PostgresName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepo(db *sql.DB) CLIRepo {
	return &Repo{db: db}
}
