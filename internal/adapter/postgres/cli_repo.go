package postgres

import (
	"RSSHub/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (r *Repo) InsertFeed(ctx context.Context, body domain.Command) error {
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

func (r *Repo) GetFeeds(ctx context.Context, count int) ([]domain.Feed, error) {
	query := `SELECT id, name, url FROM feeds WHERE created_at = updated_at LIMIT $1;`
	rows, err := r.db.QueryContext(ctx, query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Feed

	for rows.Next() {
		var item domain.Feed
		if err = rows.Scan(item.ID, &item.Name, &item.URL); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *Repo) InsertArticles(ctx context.Context, article domain.Article, id string) error {
	//	updateQuery := "UPDATE feeds SET updated_at = $1 WHERE id = $2;"
	//
	//	insertQuery := `INSERT INTO articles (title, link, description, published_at, feed_id)
	//			  VALUES ($1, $2, $3, $4, $5);`
	//
	//	pubData, err := time.Parse(time.RFC1123Z, article.PubDate)
	//	if err != nil {
	//		pubData = time.Now()
	//	}
	//
	//	tx, err := r.db.Begin()
	//	if err != nil {
	//		return err
	//	}
	//	defer tx.Rollback()
	//
	//	err = tx.QueryRowContext(ctx, updateQuery, time.Now(), id).Err()
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = tx.ExecContext(ctx, insertQuery, article.Title, article.Link, article.Description, pubData, id)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = tx.Commit()
	//	if err != nil {
	//		return err
	//	}
	//
	return nil
}

func (r *Repo) BatchInsert(articles []*domain.Article) error {
	query := `INSERT INTO articles (created_at, updated_at, title, link, description, published_at, feed_id)
			  VALUES `

	valueStrings := []string{}
	args := []interface{}{}
	i := 1

	for _, article := range articles {
		article.CreatedAt = time.Now()
		article.UpdatedAt = time.Now()

		// формируем часть вида ($1,$2,$3,...)
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i, i+1, i+2, i+3, i+4, i+5, i+6))

		args = append(args,
			article.CreatedAt,
			article.UpdatedAt,
			article.Title,
			article.Link,
			article.Description,
			article.PubDate,
			article.FeedID,
		)

		i += 7
	}

	query = query + strings.Join(valueStrings, ",") + " ON CONFLICT (link) DO NOTHING"

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repo) UpdateFeed(feedID string) error {
	updateQuery := "UPDATE feeds SET updated_at = $1 WHERE id = $2;"

	err := r.db.QueryRow(updateQuery, time.Now(), feedID).Err()
	return err
}
