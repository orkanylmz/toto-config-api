package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"toto-config-api/internal/common/logs"
	"toto-config-api/internal/common/server"
	"toto-config-api/internal/skuconfig/ports"
	"toto-config-api/internal/skuconfig/service"
)

func main() {
	logs.Init()

	ctx := context.Background()

	app := service.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

}
