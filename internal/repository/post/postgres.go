package repo

import (
	"context"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
)

type PostRepository struct {
	log   *slog.Logger
	store map[int]models.Post
	idx   int
	// Conn *sql.DB
}

func NewPostRepository(log *slog.Logger) *PostRepository {
	return &PostRepository{
		log:   log,
		store: make(map[int]models.Post),
		idx:   0,
	}
}

func (p *PostRepository) GetByID(ctx context.Context, id int) (models.Post, error) {
	p.log.Info("PostRepository.GetByID", slog.Int("id", id))
	post, ok := p.store[id]
	if !ok {
		p.log.Info("PostRepository.PostNotFound", slog.Int("id", id))
		return models.Post{}, ErrPostNotFound
	}
	return post, nil
}

func (p *PostRepository) CreatePost(ctx context.Context, post models.Post) (int, error) {
	p.idx++
	post.ID = p.idx
	p.log.Info("PostRepository.Save", slog.Int("id", post.ID))
	p.store[post.ID] = post
	return post.ID, nil
}
