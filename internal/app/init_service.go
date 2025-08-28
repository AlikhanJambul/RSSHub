package app

import (
	"RSSHub/internal/adapter/postgres"
	"RSSHub/internal/domain"
	"context"
)

type Service struct {
	cliRepo postgres.CLIRepo
}

type CLIService interface {
	AddService(ctx context.Context, body domain.Command) error
	DeleteService(ctx context.Context, body domain.Command) error
	ListService(ctx context.Context, count int) ([]domain.Feed, error)
}

func NewService(repo postgres.CLIRepo) CLIService {
	return &Service{cliRepo: repo}
}
