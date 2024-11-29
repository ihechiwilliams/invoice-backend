package errors

import (
	"invoice-backend/internal/constants"

	"github.com/samber/oops"
	"go.temporal.io/sdk/temporal"
)

const (
	APIErrorFmt = "API error service_name=%s http_status=%d body=%s"
)

var (
	DefaultRetryableError         = oops.Tags(Retryable)
	DefaultNonRetryableError      = oops.Tags(NonRetryable)
	DefaultAlertNonRetryableError = oops.Tags(Alert, NonRetryable)
)

func NewAPIRetryableError(serviceName constants.ServiceName, status int, body []byte) error {
	strBody := string(body)

	return DefaultRetryableError.
		With("service", serviceName).
		With("status_code", status).
		Errorf(APIErrorFmt, serviceName, status, strBody)
}

func NewRetryableError(errMsg, errType string, details ...interface{}) error {
	return temporal.NewApplicationErrorWithCause(errMsg, errType, DefaultRetryableError.Errorf("%s", errMsg), details)
}

func NewNonRetryableError(errMsg, errType string, details ...interface{}) error {
	return temporal.NewNonRetryableApplicationError(errMsg, errType, DefaultNonRetryableError.Errorf("%s", errMsg), details)
}

func NewAPINonRetryableError(serviceName constants.ServiceName, status int, body []byte) error {
	strBody := string(body)

	return DefaultNonRetryableError.
		With("service", serviceName).
		With("status_code", status).
		Errorf(APIErrorFmt, serviceName, status, strBody)
}

func NewAPIAlertNonRetryableError(serviceName constants.ServiceName, status int, body []byte) error {
	strBody := string(body)

	return DefaultAlertNonRetryableError.
		With("service", serviceName).
		With("status_code", status).
		Errorf(APIErrorFmt, serviceName, status, strBody)
}

func NewNetworkError(err error, serviceName constants.ServiceName) error {
	return DefaultRetryableError.
		With("service", serviceName).
		Wrap(err)
}

func APIErrorToWorkflowError(err error, details ...any) error {
	switch {
	case IsNonRetryableError(err):
		return temporal.NewNonRetryableApplicationError(
			err.Error(),
			constants.ErrorCodeApiError.String(),
			err,
			details,
		)
	case IsRetryableError(err):
		return temporal.NewApplicationErrorWithCause(
			err.Error(),
			constants.ErrorCodeApiError.String(),
			err,
			details,
		)
	default:
		return err
	}
}
