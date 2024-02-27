package post_repo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
)

type PostRepository interface {
	GetByID(ctx context.Context, id int) (models.Post, error)
	Create(ctx context.Context, post models.Post) (int, error)
}

type Repository struct {
	conn *sql.DB
	log  *slog.Logger
}

func New(conn *sql.DB, log *slog.Logger) *Repository {
	return &Repository{
		conn: conn,
		log:  log,
	}
}

func (p *Repository) GetByID(ctx context.Context, id int) (models.Post, error) {
	const op = "PostRepository.GetByID"

	query := `SELECT * FROM posts WHERE id = $1`

	post, err := p.fetch(ctx, query, id)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return models.Post{}, fmt.Errorf("can't fetch post: %s", err.Error())
	}

	if len(post) == 0 {
		p.log.Error(op, ErrPostNotFound)
		return models.Post{}, ErrPostNotFound
	}

	return post[0], nil
}

func (p *Repository) Create(ctx context.Context, post models.Post) (int, error) {
	const op = "PostRepository.CreatePost"

	query := `INSERT INTO posts (author_id, body) VALUES ($1, $2) RETURNING id`
	id, err := p.insertAndGetId(ctx, query, post.AuthorID, post.Body)
	if err != nil {
		p.log.Error(op + ": " + err.Error())
		return 0, fmt.Errorf("can't insert post: %s", err.Error())
	}

	return id, nil
}

func (p *Repository) fetch(ctx context.Context, query string, args ...any) ([]models.Post, error) {
	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.Error(err.Error())
		return nil, fmt.Errorf("can't exec query: select posts: %s", err.Error())
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			p.log.Error(err.Error())
		}
	}()

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Body); err != nil {
			return nil, fmt.Errorf("can't scan post: %s", err.Error())
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Repository) insertAndGetId(ctx context.Context, query string, args ...any) (int, error) {
	const op = "PostRepository.insertAndGetId"

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return 0, fmt.Errorf("can't prepare query: %s", err.Error())
	}

	res, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return 0, fmt.Errorf("can't exec query: insert post and get id: %s", err.Error())
	}

	defer func() {
		err = res.Close()
		if err != nil {
			p.log.Error(op + ":" + err.Error())
		}
	}()

	var id int
	for res.Next() {
		if err := res.Scan(&id); err != nil {
			p.log.Error(op + ":" + err.Error())
			return 0, fmt.Errorf("can't scan post id: %s", err.Error())
		}
	}
	return id, nil
}
