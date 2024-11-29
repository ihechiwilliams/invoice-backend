package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

const (
	formatEmail = "email"
	formatUUID  = "uuid"
	// excludeUUIDFallbackUUID is a fallback UUID that will be excluded from the regex.
	excludeUUIDFallbackUUID = `|00000000-0000-0000-0000-000000000000`
)

// formatOfStringForUUIDOfRFC4122WithoutFallbackUUID is regex for UUID v1-v5 as specified by RFC4122 without fallback UUID.
var formatOfStringForUUIDOfRFC4122WithoutFallbackUUID = strings.ReplaceAll(openapi3.FormatOfStringForUUIDOfRFC4122, excludeUUIDFallbackUUID, "")

type errorRendererFunc func([]openapi3.SchemaError, http.ResponseWriter, *http.Request)

type ValidationMiddleware struct {
	doc            *openapi3.T
	errorRenderer  errorRendererFunc
	openAPIOptions *openapi3filter.Options
}

func NewValidationMiddleware(options ...func(middleware *ValidationMiddleware)) *ValidationMiddleware {
	m := &ValidationMiddleware{}

	for _, o := range options {
		o(m)
	}

	if m.openAPIOptions == nil {
		m.openAPIOptions = &openapi3filter.Options{MultiError: true}
	}

	return m
}

func WithDoc(doc *openapi3.T) func(m *ValidationMiddleware) {
	return func(m *ValidationMiddleware) {
		doc.Servers = nil
		m.doc = doc
	}
}

func WithErrorRenderer(fn errorRendererFunc) func(m *ValidationMiddleware) {
	return func(m *ValidationMiddleware) {
		m.errorRenderer = fn
	}
}

func WithKinOpenAPIDefaults() func(m *ValidationMiddleware) {
	return func(m *ValidationMiddleware) {
		if !openapi3.SchemaErrorDetailsDisabled {
			openapi3.SchemaErrorDetailsDisabled = true
		}

		_, emailOk := openapi3.SchemaStringFormats[formatEmail]
		if !emailOk {
			openapi3.DefineStringFormat(formatEmail, openapi3.FormatOfStringForEmail)
		}

		_, uuidOk := openapi3.SchemaStringFormats[formatUUID]
		if !uuidOk {
			openapi3.DefineStringFormat(formatUUID, formatOfStringForUUIDOfRFC4122WithoutFallbackUUID)
		}
	}
}

func WithOpenAPIOptions(opts *openapi3filter.Options) func(m *ValidationMiddleware) {
	return func(m *ValidationMiddleware) {
		m.openAPIOptions = opts
	}
}

func (m *ValidationMiddleware) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := m.validateRequest(w, r)
			if err == nil {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusBadRequest)

			schemaErrs := make([]openapi3.SchemaError, 0)

			schemaErrs = collectErrors(&schemaErrs, err)

			if m.errorRenderer != nil {
				m.errorRenderer(schemaErrs, w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			err = json.NewEncoder(w).Encode(schemaErrs)
			if err != nil {
				panic(err)
			}
		})
	}
}

func (m *ValidationMiddleware) validateRequest(w http.ResponseWriter, r *http.Request) error {
	validator, err := gorillamux.NewRouter(m.doc)
	if err != nil {
		panic(err)
	}

	route, pathParams, err := validator.FindRoute(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return nil
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    m.openAPIOptions,
	}

	return openapi3filter.ValidateRequest(r.Context(), requestValidationInput)
}

func collectErrors(accumulator *[]openapi3.SchemaError, originalErr error) []openapi3.SchemaError {
	switch rootErr := originalErr.(type) {
	case openapi3.MultiError:
		for _, err := range rootErr {
			switch innerErr1 := err.(type) {
			case *openapi3filter.RequestError:
				switch innerErr2 := innerErr1.Err.(type) {
				case openapi3.MultiError:
					collectErrors(accumulator, innerErr2)
				default:
					*accumulator = append(*accumulator, openapi3.SchemaError{Reason: innerErr2.Error()})
				}
			case openapi3.MultiError:
				collectErrors(accumulator, innerErr1)
			case *openapi3.SchemaError:
				*accumulator = append(*accumulator, *innerErr1)
			default:
				*accumulator = append(*accumulator, openapi3.SchemaError{Reason: innerErr1.Error()})
			}
		}
	default:
		*accumulator = append(*accumulator, openapi3.SchemaError{Reason: rootErr.Error()})
	}

	return *accumulator
}
