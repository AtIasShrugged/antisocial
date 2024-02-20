package router

import (
	"net/http"

	"github.com/AtIasShrugged/antisocial/internal/app/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /ping", http.HandlerFunc(controllers.Ping))

	return mux
}
