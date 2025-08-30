package app

import (
	"context"

	"RSSHub/internal/adapter/postgres"
	"RSSHub/internal/domain"
)

type Service struct {
	cliRepo postgres.CLIRepo
}

type CLIService interface {
	AddService(ctx context.Context, body domain.Command) error
	DeleteService(ctx context.Context, body domain.Command) error
	ListService(ctx context.Context, count int) ([]domain.Feed, error)
	ListArticleService(ctx context.Context, name string, count int) ([]domain.Article, error)
}

func NewService(repo postgres.CLIRepo) CLIService {
	return &Service{cliRepo: repo}
}
