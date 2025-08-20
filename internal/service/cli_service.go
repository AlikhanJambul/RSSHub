package service

import (
	"RSSHub/internal/apperrors"
	"RSSHub/internal/models"
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
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

	var feed models.RSSFeed

	if err = xml.Unmarshal(body, &feed); err != nil {
		return false
	}

	if feed.Channel.Title == "" {
		return false
	}

	for idx, item := range feed.Channel.Item {
		if idx <= 3 {
			continue
		}
		slog.Info(item.Title)
	}

	return true
}

func (s *Service) AddService(ctx context.Context, body models.Command) error {
	if body.Name == "" {
		return apperrors.ErrInvalidName
	}

	if !checkUrl(body.URL) {
		return apperrors.ErrInvalidURL
	}

	if s.cliRepo.CheckName(ctx, body.Name) {
		return apperrors.ErrNameExists
	}

	return s.cliRepo.InsertFeed(ctx, body)
}
