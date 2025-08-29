package postgres

import (
	"RSSHub/internal/domain"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

type CLIRepo interface {
	InsertFeed(ctx context.Context, body domain.Command) error
	CheckNameURL(ctx context.Context, name, URL string) (bool, error)
	CheckName(ctx context.Context, name string) (bool, error)
	GetFeeds(ctx context.Context, count int) ([]domain.Feed, error)
	InsertArticles(ctx context.Context, article domain.Article, name string) error
	BatchInsert(articles []*domain.Article) error
	UpdateFeed(feedID string) error
	DeleteFeed(ctx context.Context, name string) error
	ListFeeds(ctx context.Context, count int, limit bool) ([]domain.Feed, error)
	ListArticles(ctx context.Context, name string, count int) ([]domain.Article, error)
}

func Connect(cfgDB domain.DB) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfgDB.PostgresHost,
		cfgDB.PostgresPort,
		cfgDB.PostgresUser,
		cfgDB.PostgresPass,
		cfgDB.PostgresName,
	)

	fmt.Println(connStr)

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
