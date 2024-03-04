package post_repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository interface {
	GetByID(ctx context.Context, id int) (models.Post, error)
	Create(ctx context.Context, post models.Post) (int, error)
}

type Repository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func New(pool *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{
		db:  pool,
		log: log,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int) (models.Post, error) {
	const op = "PostRepository.GetByID"

	query := `SELECT * FROM posts WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.AuthorID, &post.Body)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.log.Error(op + ":" + err.Error())
			return models.Post{}, ErrPostNotFound
		}
		r.log.Error(op + ":" + err.Error())
		return models.Post{}, fmt.Errorf("can't scan post: %s", err.Error())
	}

	return post, nil
}

func (r *Repository) Create(ctx context.Context, post models.Post) (int, error) {
	const op = "PostRepository.Create"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error(op + ":" + err.Error())
		return 0, fmt.Errorf("can't create transaction: %s", err.Error())
	}

	query := `INSERT INTO posts (author_id, body) VALUES ($1, $2) RETURNING id`
	var id int
	err = tx.QueryRow(ctx, query, post.AuthorID, post.Body).Scan(&id)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			r.log.Error(op + ":" + rollbackErr.Error())
			return 0, fmt.Errorf("can't rollback transaction: %s", rollbackErr.Error())
		}
		r.log.Error(op + ":" + err.Error())
		return 0, fmt.Errorf("can't insert post: %s", err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error(op + ":" + err.Error())
		return 0, fmt.Errorf("can't commit transaction: %s", err.Error())
	}

	return id, nil
}
