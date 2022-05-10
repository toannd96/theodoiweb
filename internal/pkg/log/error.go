package log

import (
	"context"
	"fmt"
	"os"

	goerrors "github.com/go-errors/errors"
)

// LogError prints error stack trace in dev environment and captures the error with Sentry
func LogError(ctx context.Context, err error) {
	if os.Getenv("MODE") != "production" && os.Getenv("MODE") != "staging" {
		PrintStackTrace(err)
	}
	CaptureException(ctx, err)
}

// PrintStackTrace prints out error stack trace
func PrintStackTrace(err interface{}) {
	fmt.Println("Traceback:")
	fmt.Println(goerrors.Wrap(err, 1).ErrorStack())
	fmt.Println(err)
}
