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
	log  *slog.Logger
	repo PostsRepository
}

func NewPostService(log *slog.Logger, repo PostsRepository) *PostService {
	return &PostService{
		log:  log,
		repo: repo,
	}
}

func (p *PostService) GetByID(ctx context.Context, id int) (models.Post, error) {
	p.log.Info("PostService.GetByID", slog.Int("id", id))
	return p.repo.GetByID(ctx, id)
}

func (p *PostService) CreatePost(ctx context.Context, post models.Post) (int, error) {
	p.log.Info("PostService.CreatePost", slog.Int("id", post.ID))
	return p.repo.CreatePost(ctx, post)
}
