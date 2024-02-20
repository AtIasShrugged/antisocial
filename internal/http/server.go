package server

import (
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/config"
	handlers "github.com/AtIasShrugged/antisocial/internal/http/handlers"
	repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	"github.com/AtIasShrugged/antisocial/internal/service/post"
	"github.com/labstack/echo/v4"
)

func Run(log *slog.Logger, cfg *config.Config) {
	e := echo.New()

	postRepo := repo.NewPostRepository(log)

	postService := post.NewPostService(log, postRepo)

	postHandler := handlers.NewPostsController(log, postService)

	e.GET("/posts/:id", postHandler.GetByID)
	e.POST("/posts/create", postHandler.CreatePost)

	log.Info("Starting server on: " + cfg.Server.Host + ":" + cfg.Server.Port)
	err := e.Start(":" + cfg.Server.Port)
	if err != nil {
		log.Error("Failed to start server: "+err.Error(), err)
	}
}
