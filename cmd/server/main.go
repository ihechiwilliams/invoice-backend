package main

import (
	"context"
	"fmt"
	"net/http"

	"invoice-backend/internal/appbase"
	"invoice-backend/pkg/signals"

	"github.com/rs/zerolog/log"
)

const (
	serviceName = "invoice-backend.server"
)

func main() {
	ctx, mainCtxStop := context.WithCancel(context.Background())

	app := appbase.New(
		appbase.Init(serviceName),
		appbase.WithDependencyInjector(),
	)
	defer app.Shutdown()
	fmt.Println(serviceName)

	router := buildRouter(app)

	httpServer := &http.Server{
		Addr:              app.Config.ServerAddress,
		Handler:           router,
		ReadHeaderTimeout: app.Config.HTTPServerTimeout(),
	}

	signals.HandleSignals(ctx, mainCtxStop, func() {
		shutdownErr := httpServer.Shutdown(ctx)
		if shutdownErr != nil {
			log.Fatal().Err(shutdownErr).Msg("server shutdown failed")
		}
	})

	log.Info().Msgf("started server on %s", app.Config.ServerAddress)

	serverErr := httpServer.ListenAndServe()
	if serverErr != nil {
		log.Err(serverErr).Msg("server stopped")
	}

	<-ctx.Done()
}
