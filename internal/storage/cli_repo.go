package storage

import (
	"RSSHub/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

func (r *Repo) InsertFeed(ctx context.Context, body models.Command) error {
	query := `INSERT INTO feeds (name, url) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, body.NameArg, body.URL)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CheckName(ctx context.Context, name string) bool {
	query := `SELECT name FROM feeds WHERE name = $1;`

	var n string
	err := r.db.QueryRowContext(ctx, query, name).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	if err != nil {
		return true
	}

	return true
}

func (r *Repo) GetFeeds(ctx context.Context, count int) ([]models.RSSWorkers, error) {
	query := `SELECT name, url FROM feeds LIMIT $1;`
	rows, err := r.db.QueryContext(ctx, query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.RSSWorkers

	for rows.Next() {
		var item models.RSSWorkers
		if err = rows.Scan(&item.Name, &item.URL); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *Repo) InsertArticles(ctx context.Context, feed models.RSSItem, name string) error {
	updateQuery := "UPDATE feeds SET updated_at = $1 WHERE name = $2 RETURNING id;"

	insertQuery := `INSERT INTO articles (title, link, description, published_at, feed_id)
			  VALUES ($1, $2, $3, $4, $5);`

	pubData, err := time.Parse(time.RFC1123Z, feed.PubDate)
	if err != nil {
		pubData = time.Now()
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id string

	err = tx.QueryRowContext(ctx, updateQuery, time.Now(), name).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, insertQuery, feed.Title, feed.Link, feed.Description, pubData, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
