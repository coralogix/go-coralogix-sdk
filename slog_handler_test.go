package coralogix

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlogHandler_WithAttrs(t *testing.T) {
	onTest1 := &CoralogixHandler{}
	assert.Empty(t, onTest1.data)

	onTest2 := onTest1.WithAttrs([]slog.Attr{
		slog.String("1", "1"),
		slog.Int("2", 2),
	})

	assert.Equal(t, map[string]interface{}{
		"1": "1",
		"2": int64(2),
	}, onTest2.(*CoralogixHandler).data)
}

func TestSlogHandler_AttrToMap(t *testing.T) {
	m1 := map[string]any{}
	attrToMap(m1, slog.String("1", "1"))
	attrToMap(m1, slog.Int("2", 2))

	assert.Equal(t, map[string]any{
		"1": "1",
		"2": int64(2),
	}, m1)

	m2 := map[string]any{}
	attrToMap(m2, slog.String("1", "1"))
	attrToMap(m2, slog.Any("arr", []string{"1", "2"}))
	assert.Equal(t, map[string]any{
		"1":   "1",
		"arr": []string{"1", "2"},
	}, m2)

	m3 := map[string]any{}
	attrToMap(m3, slog.String("1", "1"))
	attrToMap(m3, slog.Group("group",
		slog.String("2", "2"),
		slog.Any("arr", []string{"1", "2"}),
	))
	assert.Equal(t, map[string]any{
		"1": "1",
		"group": map[string]any{
			"2":   "2",
			"arr": []string{"1", "2"},
		},
	}, m3)
}

func TestSlogHandler_Send(t *testing.T) {
	coralogixHandler := NewCoralogixHandler(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
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
			time.Sleep(1 * time.Second)
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
