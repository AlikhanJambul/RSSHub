package service

import "RSSHub/internal/storage"

type Service struct {
	cliRepo storage.CLIRepo
}

type CLIService interface{}

func NewService(repo storage.CLIRepo) CLIService {
	return &Service{cliRepo: repo}
}
