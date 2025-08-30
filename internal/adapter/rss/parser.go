package rss

import (
	"RSSHub/internal/domain"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"
)

type Parser struct {
	client http.Client
}

func NewParser() *Parser {
	return &Parser{
		client: http.Client{Timeout: 15 * time.Second},
	}
}

func (p *Parser) ParseUrl(url string) (domain.RSSFeed, error) {
	response, err := p.client.Get(url)
	if err != nil {
		return domain.RSSFeed{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.RSSFeed{}, err
	}

	var result domain.RSSFeed

	if err := xml.Unmarshal(body, &result); err != nil {
		return domain.RSSFeed{}, err
	}

	return result, nil
}

func (p *Parser) ParseArticle(rssFeed domain.RSSFeed, feedID string) ([]*domain.Article, error) {
	articles := make([]*domain.Article, 0, len(rssFeed.Channel.Item))

	for _, item := range rssFeed.Channel.Item {
		pubDate, err := parseTimestamp(item.PubDate)
		if err != nil {
			return nil, err
		}

		article := &domain.Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     pubDate,
			FeedID:      feedID,
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func parseTimestamp(ts string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339,
		time.RFC3339Nano,
		"02 Jan 2006 15:04:05 MST",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, ts); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unsupported timestamp format: " + ts)
}
