package app

import (
	"RSSHub/internal/adapter/postgres"
	"RSSHub/internal/adapter/rss"
	"RSSHub/internal/apperrors"
	"RSSHub/internal/domain"
	"RSSHub/internal/logger"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type aggregator struct {
	countWorker    int32
	interval       time.Duration
	mu             sync.Mutex
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	jobs           chan domain.Feed
	running        bool
	ticker         *time.Ticker
	updateInterval chan time.Duration

	cliLogger logger.Logger
	cliRepo   postgres.CLIRepo
	cliParser *rss.Parser
}

type Aggregator interface {
	Start() error
	ChangeCountWorker(count int32) (string, error)
	ChangeInterval(interval string) (string, error)
	Stop()
}

func InitAggregator(count int32, inverval time.Duration, repo postgres.CLIRepo, cliLogger logger.Logger, cliParser *rss.Parser) Aggregator {
	return &aggregator{countWorker: count, interval: inverval, cliRepo: repo, cliLogger: cliLogger, cliParser: cliParser}
}

func (a *aggregator) Start() error {
	if a.running {
		return fmt.Errorf("fetch is already running")
	}

	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.wg.Add(1)
	a.ticker = time.NewTicker(a.interval)
	a.jobs = make(chan domain.Feed, 100)
	a.updateInterval = make(chan time.Duration, 1)
	a.running = true
	a.startWorkers()

	go a.runFetchLoop()

	a.cliLogger.Info("Aggregator has been started")
	return nil
}

func (a *aggregator) Stop() {
	if !a.running {
		a.cliLogger.Warn("aggregator is already stopped")
		return
	}

	a.ticker.Stop()
	a.cancel()
	close(a.jobs)
	a.wg.Wait()
	a.running = false

	a.cliLogger.Info("Aggregator has been stoped")

	return
}

func (a *aggregator) ChangeCountWorker(count int32) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if count > 15 || count < 1 {
		return "", apperrors.ErrCountWorker
	}

	oldCount := atomic.LoadInt32(&a.countWorker)
	atomic.StoreInt32(&a.countWorker, count)

	if a.running {
		dif := count - oldCount
		if dif > 0 {
			for i := int32(0); i < dif; i++ {
				a.wg.Add(1)
				go a.worker()
			}
		}
	}

	response := fmt.Sprintf("changed count from %d to %d", oldCount, count)

	a.cliLogger.Info(response)
	return response, nil
}

func (a *aggregator) ChangeInterval(interval string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return "", apperrors.ErrAggregatorStop
	}

	newInterval, err := time.ParseDuration(interval)
	if err != nil {
		return "", err
	}

	oldInterval := a.interval
	a.interval = newInterval

	a.updateInterval <- a.interval

	response := fmt.Sprintf("Interval changed from %s to %s", oldInterval, a.interval)
	a.cliLogger.Info(response)

	return response, nil
}

func (a *aggregator) runFetchLoop() {
	defer a.wg.Done()

	// a.fetchFeeds()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-a.ticker.C:
			a.fetchFeeds()
		case newInterval := <-a.updateInterval:
			a.ticker.Stop()
			a.ticker = time.NewTicker(newInterval)
		}
	}
}

func (a *aggregator) fetchFeeds() {
	numWorkers := atomic.LoadInt32(&a.countWorker)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	ctx := context.Background()

	feeds, err := a.cliRepo.GetFeeds(ctx, int(numWorkers))
	// cancel()
	if err != nil {
		a.cliLogger.Warn("Fetching feeds failed", "err", err)
		return
	}

	a.cliLogger.Info("GetFeeds finished")

	for _, feed := range feeds {
		select {
		case a.jobs <- feed:
		case <-ctx.Done():
			return
		}
	}
}

func (a *aggregator) startWorkers() {
	numWorkers := atomic.LoadInt32(&a.countWorker)

	for i := int32(0); i < numWorkers; i++ {
		a.wg.Add(1)
		go a.worker()
	}
}

func (a *aggregator) worker() {
	defer a.wg.Done()

	for {
		select {
		case feed, ok := <-a.jobs:
			if !ok {
				return
			}
			a.flush(feed)
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *aggregator) flush(feed domain.Feed) {
	parsedFeed, err := a.cliParser.ParseUrl(feed.URL)
	if err != nil {
		a.cliLogger.Warn(err.Error())
		return
	}

	a.cliLogger.Info("Feed parsed")

	articles, err := a.cliParser.ParseArticle(parsedFeed, feed.ID)
	if err != nil {
		a.cliLogger.Warn(err.Error())
		return
	}

	a.cliLogger.Info("article parsed")

	if err := a.cliRepo.BatchInsert(articles); err != nil {
		a.cliLogger.Warn(err.Error())
		return
	}

	a.cliLogger.Info("batch insert just finished")

	if err := a.cliRepo.UpdateFeed(feed.ID); err != nil {
		a.cliLogger.Warn(err.Error())
		return
	}

	a.cliLogger.Info("feed has been updated")
}
