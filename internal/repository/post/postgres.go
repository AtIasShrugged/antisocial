package repo

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
)

type PostRepository struct {
	conn  *sql.DB
	log   *slog.Logger
	store map[int]models.Post
	idx   int
}

func NewPostRepository(conn *sql.DB, log *slog.Logger) *PostRepository {
	return &PostRepository{
		conn:  conn,
		log:   log,
		store: make(map[int]models.Post),
		idx:   0,
	}
}

func (p *PostRepository) GetByID(ctx context.Context, id int) (models.Post, error) {
	const op = "PostRepository.GetByID"

	query := `SELECT * FROM posts WHERE id = $1`

	post, err := p.fetch(ctx, query, id)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return models.Post{}, err
	}

	if len(post) == 0 {
		p.log.Error(op, ErrPostNotFound)
		return models.Post{}, ErrPostNotFound
	}

	return post[0], nil
}

func (p *PostRepository) CreatePost(ctx context.Context, post models.Post) (int, error) {
	const op = "PostRepository.CreatePost"

	query := `INSERT INTO posts (author_id, body) VALUES ($1, $2) RETURNING id`
	id, err := p.insertAndGetId(ctx, query, post.AuthorID, post.Body)
	if err != nil {
		p.log.Error(op + ": " + err.Error())
		return 0, err
	}

	return id, nil
}

func (p *PostRepository) fetch(ctx context.Context, query string, args ...any) ([]models.Post, error) {
	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostRepository) insertAndGetId(ctx context.Context, query string, args ...any) (int, error) {
	const op = "PostRepository.insertAndGetId"

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return 0, err
	}

	res, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return 0, err
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
			return 0, err
		}
	}
	return id, nil
}
