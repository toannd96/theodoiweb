package log

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// CaptureException calls Sentry's CaptureException function with hub from context
func CaptureException(ctx context.Context, err error) {
	if ctx != nil {
		hub := sentry.GetHubFromContext(ctx)
		if hub != nil {
			hub.CaptureException(err)
			return
		}
	}
	sentry.CurrentHub().Clone().CaptureException(err)
}
