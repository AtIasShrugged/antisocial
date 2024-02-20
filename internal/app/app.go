package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/AtIasShrugged/antisocial/internal/app/router"
	"github.com/AtIasShrugged/antisocial/internal/config"
)

func Run(log *slog.Logger, cfg *config.Config) {
	app := router.Setup()
	fmt.Println(cfg)
	log.Info("Starting server on: " + cfg.Server.Host + ":" + cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, app)
	if err != nil {
		log.Error("Failed to start server: " + err.Error())
	}
}
