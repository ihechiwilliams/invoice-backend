package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
)

func GetCurrentHub(ctx context.Context) *sentry.Hub {
	hub := sentry.GetHubFromContext(ctx)
	if hub != nil {
		return hub
	}

	return sentry.CurrentHub().Clone()
}
