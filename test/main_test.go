package test

import (
	"os"
	"testing"

	"github.com/rizalgowandy/gdk/pkg/logx"
)

// TestMain is the entry point for the test suite.
// Start app and containers for testing.
func TestMain(m *testing.M) {
	ctx := logx.NewContext()
	err := logx.New(logx.Config{
		Debug:    true,
		AppName:  "e2e",
		Filename: "",
	})
	if err != nil {
		logx.FTL(ctx, err, "initialize logger")
	}

	err = SetupMailHog(ctx)
	if err != nil {
		logx.FTL(ctx, err, "set up mailhog container")
	}
	defer MailHogContainer.Terminate(ctx)

	code := m.Run()
	func() {
		os.Exit(code)
	}()
}
