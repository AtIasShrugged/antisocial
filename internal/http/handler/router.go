package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/config"
	post_handler "github.com/AtIasShrugged/antisocial/internal/http/handler/post"
	post_repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	"github.com/AtIasShrugged/antisocial/internal/service/post"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router(ctx context.Context, log *slog.Logger, cfg *config.Config) (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.Driver, cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error("Failed to open DB connection: "+err.Error(), err)
		return nil, err
	}
	log.Info(fmt.Sprintf("Connected to %s on port %s", cfg.DB.Driver, cfg.DB.Port))

	postRepo := post_repo.New(pool, log)

	postService := post.New(postRepo, log)

	postHandler := post_handler.New(postService, log)

	e.GET("/posts/:id", postHandler.GetByID)
	e.POST("/posts/create", postHandler.Create)

	return e, nil
}
