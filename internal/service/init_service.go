package service

import (
	"RSSHub/internal/models"
	"RSSHub/internal/storage"
	"context"
)

type Service struct {
	cliRepo storage.CLIRepo
}

type CLIService interface {
	AddService(ctx context.Context, body models.Command) error
}

func NewService(repo storage.CLIRepo) CLIService {
	return &Service{cliRepo: repo}
}
