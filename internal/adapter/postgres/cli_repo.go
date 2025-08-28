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

func (r *Repo) CheckNameURL(ctx context.Context, name, URL string) (bool, error) {
	query := `SELECT 1 FROM feeds WHERE name = $1 OR url = $2;`

	var exists int
	err := r.db.QueryRowContext(ctx, query, name, URL).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil // не найдено
	}
	if err != nil {
		return false, err // ошибка реальная
	}

	return true, nil // найдено
}

func (r *Repo) CheckName(ctx context.Context, name string) (bool, error) {
	query := `SELECT name FROM feeds WHERE name = $1;`

	err := r.db.QueryRowContext(ctx, query, name).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return true, err
	}

	return true, nil
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
		if err = rows.Scan(&item.ID, &item.Name, &item.URL); err != nil {
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

func (r *Repo) DeleteFeed(ctx context.Context, name string) error {
	deleteFeedQuery := "DELETE FROM feeds WHERE name = $1 RETURNING id;"
	deleteArticleQuery := "DELETE FROM articles WHERE feed_id = $1;"

	var id string

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, deleteFeedQuery, name).Scan(&id)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, deleteArticleQuery, id).Err()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ListFeeds(ctx context.Context, count int, limit bool) ([]domain.Feed, error) {
	query := `SELECT name, url, created_at FROM feeds ORDER BY created_at DESC`

	arg := []interface{}{}

	if limit {
		query += " LIMIT $1;"
		arg = append(arg, count)
	}

	var result []domain.Feed

	rows, err := r.db.QueryContext(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Feed

		if err = rows.Scan(&item.Name, &item.URL, &item.CreatedAt); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
