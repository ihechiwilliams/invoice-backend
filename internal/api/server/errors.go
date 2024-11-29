package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/render"
	"github.com/joomcode/errorx"
	"github.com/samber/lo"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/temporal"

	"invoice-backend/internal/constants"
	sentryUtils "invoice-backend/pkg/sentry"
)

const (
	badRequestErrorTitle = "BAD_REQUEST"
	processingErrorTitle = "PROCESSING_ERROR"
	timeoutErrorTitle    = "TIMEOUT"
	notFoundErrorTitle   = "NOT_FOUND"

	notFoundErrorDetail = "record not found"
)

func BadRequestError(badRequestErr error, w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusBadRequest

	errs := make([]Error, 0)

	err := Error{
		Code:   http.StatusText(statusCode),
		Detail: badRequestErr.Error(),
		Meta:   lo.ToPtr(map[string]interface{}{}),
		Status: statusCode,
		Title:  badRequestErrorTitle,
	}

	errs = append(errs, err)

	errResponse := ErrorResponse{Errors: errs}

	render.Status(r, statusCode)
	render.JSON(w, r, errResponse)
}

func ProcessingError(processingErr error, w http.ResponseWriter, r *http.Request) {
	errorCode := mapProcessingErrorToErrorCode(processingErr)
	ProcessingErrorWithCode(processingErr, errorCode, w, r)
}

func ProcessingErrorWithCode(processingErr error, errorCode string, w http.ResponseWriter, r *http.Request) {
	statusCode, title := mapErrorToStatusCodeAndTitle(processingErr)
	errs := make([]Error, 0)
	err := Error{
		Code:   errorCode,
		Detail: processingErr.Error(),
		Meta:   lo.ToPtr(map[string]interface{}{}),
		Status: statusCode,
		Title:  title,
	}

	errs = append(errs, err)

	errResponse := ErrorResponse{Errors: errs}

	render.Status(r, statusCode)
	render.JSON(w, r, errResponse)
}

func mapErrorToStatusCodeAndTitle(processingErr error) (int, string) {
	var (
		temporalActivityErr                *temporal.ActivityError
		temporalApplicationErr             *temporal.ApplicationError
		temporalServiceCanceledErr         *serviceerror.Canceled
		temporalServiceDeadlineExceededErr *serviceerror.DeadlineExceeded
	)

	switch {
	case errors.As(processingErr, &temporalApplicationErr),
		errors.As(processingErr, &temporalActivityErr):
		return http.StatusUnprocessableEntity, processingErrorTitle
	case errors.As(processingErr, &temporalServiceDeadlineExceededErr),
		errors.As(processingErr, &temporalServiceCanceledErr),
		errors.Is(processingErr, context.DeadlineExceeded),
		errors.Is(processingErr, context.Canceled):
		return http.StatusGatewayTimeout, timeoutErrorTitle
	case errorx.IsNotFound(processingErr):
		return http.StatusNotFound, notFoundErrorTitle
	default:
		return http.StatusUnprocessableEntity, processingErrorTitle
	}
}

func mapProcessingErrorToErrorCode(processingErr error) string {
	var temporalApplicationErr *temporal.ApplicationError

	if !errors.As(processingErr, &temporalApplicationErr) {
		return processingErrorTitle
	}

	errCode, err := constants.ParseErrorCode(temporalApplicationErr.Type())
	if err != nil {
		return processingErrorTitle
	}

	return errCode.String()
}

func (e ErrorResponse) Error() string {
	errs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		errs[i] = err.Detail
	}

	return strings.Join(errs, ";")
}

func ErrorRenderer(schemaErrors []openapi3.SchemaError, w http.ResponseWriter, r *http.Request) {
	errs := make([]Error, 0)

	for _, schemaError := range schemaErrors {
		err := Error{
			Code:   http.StatusText(http.StatusBadRequest),
			Detail: schemaError.Error(),
			Meta: lo.ToPtr(map[string]interface{}{
				"caused_by": "openapi3.SchemaError",
				"pointer":   strings.Join(schemaError.JSONPointer(), "/"),
			}),
			Status: http.StatusBadRequest,
			Title:  badRequestErrorTitle,
		}

		errs = append(errs, err)
	}

	errResponse := ErrorResponse{Errors: errs}

	sentryUtils.CaptureSentry(r.Context(), "OpenAPI validation error", errResponse)

	render.JSON(w, r, errResponse)
}

func NotFoundError(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusNotFound

	errs := make([]Error, 0)

	err := Error{
		Code:   http.StatusText(statusCode),
		Detail: notFoundErrorDetail,
		Meta:   lo.ToPtr(map[string]interface{}{}),
		Status: statusCode,
		Title:  notFoundErrorTitle,
	}

	errs = append(errs, err)

	errResponse := ErrorResponse{Errors: errs}

	render.Status(r, statusCode)
	render.JSON(w, r, errResponse)
}

func ConflictError(conflictErr error, pendingAccounts []string, w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusConflict

	errs := make([]Error, 0)

	err := Error{
		Code:   http.StatusText(statusCode),
		Detail: conflictErr.Error(),
		Meta: lo.ToPtr(map[string]interface{}{
			"pending_accounts": pendingAccounts,
		}),
		Status: statusCode,
		Title:  "Conflict Error",
	}

	errs = append(errs, err)

	errResponse := ErrorResponse{Errors: errs}

	render.Status(r, statusCode)
	render.JSON(w, r, errResponse)
}
