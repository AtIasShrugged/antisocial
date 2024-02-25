package server

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/config"
	handlers "github.com/AtIasShrugged/antisocial/internal/http/handlers"
	repo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	"github.com/AtIasShrugged/antisocial/internal/service/post"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func Run(log *slog.Logger, cfg *config.Config) {
	e := echo.New()
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.Driver, cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	dbConn, err := sql.Open(cfg.DB.Driver, dsn)
	fmt.Println(cfg)
	if err != nil {
		log.Error("Failed to open DB connection: "+err.Error(), err)
	}
	log.Info(fmt.Sprintf("Connected to %s on port %s", cfg.DB.Driver, cfg.DB.Port))

	postRepo := repo.NewPostRepository(dbConn, log)

	postService := post.NewPostService(postRepo, log)

	postHandler := handlers.NewPostsController(postService, log)

	e.GET("/posts/:id", postHandler.GetByID)
	e.POST("/posts/create", postHandler.CreatePost)

	log.Info("Starting server on: " + cfg.Server.Host + ":" + cfg.Server.Port)
	err = e.Start(":" + cfg.Server.Port)
	if err != nil {
		log.Error("Failed to start server: "+err.Error(), err)
	}
}
