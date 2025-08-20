package aggregator

import (
	"RSSHub/internal/apperrors"
	"RSSHub/internal/logger"
	"RSSHub/internal/models"
	"RSSHub/internal/storage"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	countWorker int
	interval    time.Duration
	cliRepo     storage.CLIRepo
	intervalCh  chan time.Duration
	mu          sync.Mutex
	cliLogger   logger.Logger
}

type Aggregator interface {
	Start()
	ChangeCountWorker(count int) error
	ChangeInterval(interval string) error
	Stop()
}

func InitAggregator(count int, inverval time.Duration, repo storage.CLIRepo, cliLogger logger.Logger) Aggregator {
	return &Manager{countWorker: count, interval: inverval, cliRepo: repo, intervalCh: make(chan time.Duration), cliLogger: cliLogger}
}

func (m *Manager) Start() {
	ticker := time.NewTicker(m.interval)

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			feeds, err := m.cliRepo.GetFeeds(ctx, m.countWorker)
			cancel()
			if err != nil {
				m.cliLogger.Error(err.Error())
				continue
			}

			if len(feeds) != m.countWorker {
				m.cliLogger.Info("Count of feed is not equal count of workers")
			}

			var wg sync.WaitGroup
			for _, feed := range feeds {
				wg.Add(1)

				go func(name, url string) {
					defer wg.Done()

					feed, err := parseUrl(url)
					if err != nil {
						m.cliLogger.Error(err.Error())
						return
					}

					for _, item := range feed.Channel.Item {
						err := m.cliRepo.InsertArticles(context.Background(), item, name)
						if err != nil {
							m.cliLogger.Error(err.Error())
							continue
						}
					}

				}(feed.Name, feed.URL)

			}

			wg.Wait()
		case m.interval = <-m.intervalCh:
			ticker.Stop()
			ticker = time.NewTicker(m.interval)
			m.cliLogger.Info("Interval has been changed")
		}
	}
}

func (m *Manager) ChangeCountWorker(count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if count > 15 || count < 1 {
		return apperrors.ErrCountWorker
	}
	m.countWorker = count

	return nil
}

func (m *Manager) ChangeInterval(interval string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	validInterval, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}

	m.intervalCh <- validInterval

	return nil
}

func (m *Manager) Stop() {}

func parseUrl(url string) (models.RSSFeed, error) {
	client := &http.Client{Timeout: time.Second * 15}
	response, err := client.Get(url)
	if err != nil {
		return models.RSSFeed{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return models.RSSFeed{}, err
	}

	var result models.RSSFeed

	if err := xml.Unmarshal(body, &result); err != nil {
		return models.RSSFeed{}, err
	}

	if result.Channel.Title == "" {
		return models.RSSFeed{}, errors.New("no feed found")
	}

	return result, nil
}
