package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	"github.com/labstack/echo/v4"
)

type PostService interface {
	GetByID(ctx context.Context, id int) (models.Post, error)
	CreatePost(ctx context.Context, post models.Post) (int, error)
}

type PostHandler struct {
	log     *slog.Logger
	service PostService
}

func NewPostsController(log *slog.Logger, service PostService) *PostHandler {
	return &PostHandler{
		log:     log,
		service: service,
	}
}

func (p *PostHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		p.log.Error("PostHandler.GetByID", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	post, err := p.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == repo.ErrPostNotFound.Error() {
			return c.String(http.StatusNotFound, err.Error())
		}
		p.log.Error("PostHandler.GetByID", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	p.log.Info("PostHandler.GetByID", slog.Int("id", 1))
	return c.JSON(http.StatusOK, post)
}

func (p *PostHandler) CreatePost(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		p.log.Error("PostHandler.CreatePost", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(post)
	id, err := p.service.CreatePost(c.Request().Context(), post)
	if err != nil {
		p.log.Error("PostHandler.CreatePost", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	p.log.Info("PostHandler.CreatePost", slog.Int("id", id))
	return c.JSON(http.StatusOK, id)
}
