package post_handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	"github.com/labstack/echo/v4"
)

type PostService interface {
	GetByID(ctx context.Context, id int) (models.Post, error)
	Create(ctx context.Context, post models.Post) (int, error)
}

type PostHandler struct {
	service PostService
	log     *slog.Logger
}

func New(service PostService, log *slog.Logger) *PostHandler {
	return &PostHandler{
		service: service,
		log:     log,
	}
}

func (p *PostHandler) GetByID(c echo.Context) error {
	const op = "PostHandler.GetByID"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	post, err := p.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrPostNotFound) {
			p.log.Error(op + ":" + err.Error())
			return c.String(http.StatusNotFound, err.Error())
		}
		p.log.Error(op + ":" + err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, post)
}

func (p *PostHandler) Create(c echo.Context) error {
	const op = "PostHandler.Create"

	var post models.Post
	if err := c.Bind(&post); err != nil {
		p.log.Error(op + ":" + err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := p.service.Create(c.Request().Context(), post)
	if err != nil {
		p.log.Error(op + ":" + err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, id)
}
