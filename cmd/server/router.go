package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"invoice-backend/internal/api"
	"invoice-backend/internal/appbase"
)

func buildRouter(app *appbase.AppBase) *chi.Mux {
	fmt.Println("hey")
	mux := do.MustInvokeNamed[*chi.Mux](app.Injector, appbase.InjectorApplicationRouter)
	routes := do.MustInvoke[*api.Routes](app.Injector)

	api.InitRoutes(mux, routes)

	return mux
}
