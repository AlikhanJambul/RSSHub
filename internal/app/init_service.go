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
}

func NewService(repo postgres.CLIRepo) CLIService {
	return &Service{cliRepo: repo}
}
