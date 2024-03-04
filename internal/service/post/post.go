package post

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	post_repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
)

type PostService struct {
	repo post_repo.PostRepository
	log  *slog.Logger
}

func New(repo post_repo.PostRepository, log *slog.Logger) *PostService {
	return &PostService{
		log:  log,
		repo: repo,
	}
}

func (p *PostService) GetByID(ctx context.Context, id int) (models.Post, error) {
	const op = "PostService.GetByID"

	post, err := p.repo.GetByID(ctx, id)
	if err != nil {
		p.log.Error(op+": "+err.Error(), err)
		return models.Post{}, err
	}
	return post, nil
}

func (p *PostService) Create(ctx context.Context, post models.Post) (int, error) {
	const op = "PostService.Create"

	id, err := p.repo.Create(ctx, post)
	if err != nil {
		p.log.Error(op+": "+err.Error(), err)
		return 0, fmt.Errorf("error from post_repository: %s", err.Error())
	}
	return id, nil
}
