package http

import (
	"bytes"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func WithLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

type responseRecorder struct {
	http.ResponseWriter

	body   bytes.Buffer
	status int

	loggedStatusHeader bool
	loggedBody         bool
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	if r.loggedStatusHeader {
		return
	}

	r.status = statusCode
	r.loggedStatusHeader = true

	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	if r.loggedBody {
		return 0, nil
	}

	r.body.Write(body)
	r.loggedBody = true

	return r.ResponseWriter.Write(body)
}

func WithErrorLogs() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			recorder := &responseRecorder{ResponseWriter: w}

			next.ServeHTTP(recorder, r)

			if recorder.status >= http.StatusBadRequest && recorder.status <= http.StatusNetworkAuthenticationRequired {
				log.Ctx(r.Context()).Error().
					RawJSON("error_body", compactJSON(recorder.body.Bytes())).
					Str("http_method", r.Method).
					Str("uri", r.RequestURI).
					Int("http_status", recorder.status).
					Str("http_status_text", http.StatusText(recorder.status)).
					Msg("http error")
			}
		}

		return http.HandlerFunc(fn)
	}
}

func WithAccessLogs() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestLogger := log.Ctx(r.Context()).With().
				Timestamp().
				Str("type", "access").
				Str("host", r.URL.Host).
				Str("path", r.URL.Path).
				Str("protocol", r.Proto).
				Str("method", r.Method).
				Str("headers-host", r.Header.Get("Host")).
				Str("headers-content-length", r.Header.Get("Content-Length")).
				Str("headers-x-app-version", r.Header.Get("X-App-Version")).
				Str("headers-x-b3-trace-id", r.Header.Get("X-B3-Traceid")).
				Str("headers-x-user-id", r.Header.Get("X-User-Id")).
				Str("params-version-code", r.URL.Query().Get("version_code")).
				Str("headers-user-agent", r.Header.Get("User-Agent")).Logger()

			requestLogger.Debug().Msgf("%s request_started", r.URL.Path)

			startedAt := time.Now()
			responseWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				finishedAt := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					requestLogger.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")

					http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				requestLogger.Debug().
					Int("bytes_out", responseWriter.BytesWritten()).
					Int("status_code", responseWriter.Status()).
					Dur("duration", finishedAt.Sub(startedAt)).
					Msgf("%s request_finished", r.URL.Path)
			}()

			next.ServeHTTP(responseWriter, r)
		}

		return http.HandlerFunc(fn)
	}
}
