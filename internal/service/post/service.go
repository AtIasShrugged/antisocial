package post

import (
	"context"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
)

type PostsRepository interface {
	GetByID(ctx context.Context, id int) (models.Post, error)
	CreatePost(ctx context.Context, post models.Post) (int, error)
}

type PostService struct {
	repo PostsRepository
	log  *slog.Logger
}

func NewPostService(repo PostsRepository, log *slog.Logger) *PostService {
	return &PostService{
		log:  log,
		repo: repo,
	}
}

func (p *PostService) GetByID(ctx context.Context, id int) (models.Post, error) {
	return p.repo.GetByID(ctx, id)
}

func (p *PostService) CreatePost(ctx context.Context, post models.Post) (int, error) {
	return p.repo.CreatePost(ctx, post)
}
