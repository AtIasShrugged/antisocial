package server

import (
	"context"
	"log/slog"

	"github.com/AtIasShrugged/antisocial/internal/config"
	"github.com/AtIasShrugged/antisocial/internal/http/handler"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

func Run(log *slog.Logger, cfg *config.Config) {
	ctx := context.Background()
	router, err := handler.Router(ctx, log, cfg)
	if err != nil {
		log.Error("Failed to create router: "+err.Error(), err)
	}

	err = router.Start(":" + cfg.Server.Port)
	if err != nil {
		log.Error("Failed to start server: "+err.Error(), err)
	}
}
