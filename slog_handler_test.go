package coralogix

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlogHandler_Send(t *testing.T) {
	coralogixHandler := NewCoralogixHandler(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	defer func() { recover() }()
	defer coralogixHandler.Stop()

	log := slog.New(coralogixHandler)
	slog.SetDefault(log)

	attr := slog.Attr{Key: "extra", Value: slog.StringValue("additional")}

	testcases := []struct {
		name     string
		logfn    func(message string, args ...interface{})
		severity uint
	}{
		{
			name:     "test debug",
			severity: Level.DEBUG,
			logfn: func(message string, args ...interface{}) {
				log.Debug(message, args...)
			},
		},
		{
			name:     "test info",
			severity: Level.INFO,
			logfn: func(message string, args ...interface{}) {
				log.Info(message, args...)
			},
		},
		{
			name:     "test warn",
			severity: Level.WARNING,
			logfn: func(message string, args ...interface{}) {
				log.Warn(message, args...)
			},
		},
		{
			name:     "test error",
			severity: Level.ERROR,
			logfn: func(message string, args ...interface{}) {
				log.Error(message, args...)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			msg := fmt.Sprintf("%s (%s)", tc.name, t.Name())
			tc.logfn(msg, attr)
			time.Sleep(time.Duration(1) * time.Second)
			bulk, ok := mockHTTPServerMap[t.Name()]
			assert.True(t, ok, "%s key not found in mockHTTPServerMap", t.Name())

			var msgExists bool
			for _, entry := range bulk.LogEntries {
				if msgExists = strings.Contains(entry.Text, tc.name); msgExists {
					assert.Equal(t, tc.severity, entry.Severity)
					assert.True(t, strings.Contains(entry.Text, attr.Value.String()),
						"entry Text does not contain extra field", entry.Text, attr)
					break
				}
			}
			assert.True(t, msgExists, "no matching message found", string(bulk.ToJSON()))
		})
	}
}
