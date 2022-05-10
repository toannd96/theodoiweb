package log

import (
	"context"
	"os"

	"github.com/getsentry/sentry-go"
)

// HandlePanic calls recover function and log it with Sentry
func HandlePanic(ctx context.Context) {
	err := recover()
	if err != nil {
		LogPanic(ctx, err)
	}
}

// LogPanic logs error using Sentry recover function with hub from context
func LogPanic(ctx context.Context, err interface{}) {
	if os.Getenv("MODE") != "production" && os.Getenv("MODE") != "staging" {
		PrintStackTrace(err)
	}
	if ctx != nil {
		hub := sentry.GetHubFromContext(ctx)
		if hub != nil {
			hub.RecoverWithContext(ctx, err)
			return
		}
	}
	sentry.CurrentHub().Clone().Recover(err)
}
