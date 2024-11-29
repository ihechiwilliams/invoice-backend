package errors

import (
	"github.com/samber/lo"
	"github.com/samber/oops"
)

func IsNonRetryableError(err error) bool {
	oopsErr, isOopsErr := oops.AsOops(err)
	if !isOopsErr {
		return false
	}

	return lo.Contains(oopsErr.Tags(), NonRetryable)
}

func IsRetryableError(err error) bool {
	oopsErr, isOopsErr := oops.AsOops(err)
	if !isOopsErr {
		return false
	}

	return lo.Contains(oopsErr.Tags(), Retryable)
}

func IsNotFoundError(err error) bool {
	oopsErr, isOopsErr := oops.AsOops(err)
	if !isOopsErr {
		return false
	}

	return lo.Contains(oopsErr.Tags(), NotFound)
}

func IsAlertError(err error) bool {
	oopsErr, isOopsErr := oops.AsOops(err)
	if !isOopsErr {
		return false
	}

	return lo.Contains(oopsErr.Tags(), Alert)
}
