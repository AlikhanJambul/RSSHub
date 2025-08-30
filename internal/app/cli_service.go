package app

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"time"

	"RSSHub/internal/apperrors"
	"RSSHub/internal/domain"
)

func checkUrl(bodyUrl string) bool {
	_, err := url.ParseRequestURI(bodyUrl)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: time.Second * 5}

	resp, err := client.Get(bodyUrl)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var feed domain.RSSFeed

	if err = xml.Unmarshal(body, &feed); err != nil {
		return false
	}

	if feed.Channel.Title == "" {
		return false
	}

	return true
}

func (s *Service) AddService(ctx context.Context, body domain.Command) error {
	if body.NameArg == "" {
		return apperrors.ErrInvalidName
	}

	ok, err := s.cliRepo.CheckNameURL(ctx, body.NameArg, body.URL)
	if err != nil {
		return err
	}

	if ok {
		return apperrors.ErrNameExists
	}

	if !checkUrl(body.URL) {
		return apperrors.ErrInvalidURL
	}

	return s.cliRepo.InsertFeed(ctx, body)
}

func (s *Service) DeleteService(ctx context.Context, body domain.Command) error {
	ok, err := s.cliRepo.CheckName(ctx, body.NameArg)
	if err != nil || !ok {
		return apperrors.ErrNameNil
	}

	return s.cliRepo.DeleteFeed(ctx, body.NameArg)
}

func (s *Service) ListService(ctx context.Context, count int) ([]domain.Feed, error) {
	if count < 0 {
		return nil, apperrors.ErrListNum
	}

	limit := true

	if count == 0 {
		limit = false
	}

	return s.cliRepo.ListFeeds(ctx, count, limit)
}

func (s *Service) ListArticleService(ctx context.Context, name string, count int) ([]domain.Article, error) {
	ok, err := s.cliRepo.CheckName(ctx, name)
	if err != nil || !ok {
		return nil, apperrors.ErrNameNil
	}

	if count <= 0 {
		return nil, apperrors.ErrListNum
	}

	return s.cliRepo.ListArticles(ctx, name, count)
}
