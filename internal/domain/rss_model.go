package domain

import (
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSWorkers struct {
	Name string
	URL  string
}

type Feed struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	URL       string
}

type Article struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	FeedID      string
}
