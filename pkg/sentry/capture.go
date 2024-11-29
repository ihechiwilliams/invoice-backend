package sentry

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"
)

func CaptureSentry(ctx context.Context, msg string, err error) {
	hub := GetCurrentHub(ctx)

	evt := sentry.NewEvent()
	evt.Message = msg
	evt.Extra["error"] = err.Error()

	hub.CaptureEvent(evt)
}

func CaptureSentryWithContext(ctx context.Context, msg, contextKey string, contextValue map[string]interface{}, userID *string) {
	hub := GetCurrentHub(ctx)

	hub.ConfigureScope(func(scope *sentry.Scope) {
		if userID != nil {
			scope.SetUser(sentry.User{
				ID: *userID,
			})
		}

		scope.SetFingerprint([]string{msg, lo.FromPtrOr(userID, "unknown")})
		scope.SetContext(contextKey, contextValue)
	})

	evt := sentry.NewEvent()
	evt.Message = msg

	hub.CaptureEvent(evt)
}

func CaptureSQSEventError(ctx context.Context, eventType string, msg *types.Message, err error) {
	hub := GetCurrentHub(ctx)

	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("sqs", map[string]interface{}{
			"event_type":   eventType,
			"message_body": msg.Body,
			"message_id":   msg.MessageId,
		})
	})

	evt := sentry.NewEvent()
	evt.Message = fmt.Sprintf("events %s %s", eventType, err.Error())
	evt.Transaction = err.Error()
	evt.Extra["error"] = err.Error()

	hub.CaptureEvent(evt)
}
