package appbase

import (
	"net/http"
	"time"

	httpUtils "invoice-backend/pkg/http"
	openAPIUtils "invoice-backend/pkg/openapi"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	chiDDTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
)

func NewRouterMux(serviceName string, logger *zerolog.Logger, openAPIMiddleware *openAPIUtils.ValidationMiddleware, timeout time.Duration) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(chiMiddleware.RequestID)
	mux.Use(chiMiddleware.Recoverer)
	mux.Use(httpUtils.WithLogger(*logger))
	mux.Use(httpUtils.WithErrorLogs())
	mux.Use(chiMiddleware.Timeout(timeout))

	mux.Use(chiDDTrace.Middleware(
		chiDDTrace.WithServiceName(serviceName),
		chiDDTrace.WithIgnoreRequest(func(r *http.Request) bool {
			return r.URL.Path == "/" || r.URL.Path == "/healthz" || r.URL.Path == "/readyz"
		}),
	))

	mux.Use(chiMiddleware.Heartbeat("/"))
	mux.Use(chiMiddleware.Heartbeat("/healthz"))
	mux.Use(chiMiddleware.Heartbeat("/readyz"))

	mux.Use(openAPIMiddleware.Handler())

	return mux
}
