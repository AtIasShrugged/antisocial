package post

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	post_repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	repoMock "github.com/AtIasShrugged/antisocial/internal/repository/post/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repoMock.NewMockPostRepository(ctrl)

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	in := 1

	mockResp := models.Post{
		ID:       1,
		AuthorID: 1,
		Body:     "test",
	}

	expected := models.Post{
		ID:       1,
		AuthorID: 1,
		Body:     "test",
	}
	repo.EXPECT().GetByID(ctx, in).Return(mockResp, nil).Times(1)

	service := New(repo, log)
	post, err := service.GetByID(ctx, in)
	require.NoError(t, err)
	require.Equal(t, expected, post)
}

func TestGetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repoMock.NewMockPostRepository(ctrl)

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repoErr := post_repo.ErrPostNotFound
	in := 1
	expected := models.Post{}
	repo.EXPECT().GetByID(ctx, in).Return(models.Post{}, repoErr).Times(1)

	service := New(repo, log)
	post, err := service.GetByID(ctx, in)
	require.Error(t, err)
	require.EqualError(t,
		repoErr,
		err.Error(),
	)
	require.Equal(t, expected, post)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repoMock.NewMockPostRepository(ctrl)

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	in := models.Post{
		AuthorID: 1,
		Body:     "test",
	}
	id := 1
	repo.EXPECT().Create(ctx, in).Return(id, nil).Times(1)

	service := New(repo, log)
	postID, err := service.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, id, postID)
}

func TestCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repoMock.NewMockPostRepository(ctrl)

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repoErr := errors.New("can't insert post: db is down")
	in := models.Post{
		AuthorID: 1,
		Body:     "test",
	}
	repo.EXPECT().Create(ctx, in).Return(0, repoErr).Times(1)

	service := New(repo, log)
	id, err := service.Create(ctx, in)
	require.Error(t, err)
	require.EqualError(t,
		repoErr,
		err.Error(),
	)
	require.Equal(t, 0, id)
}
