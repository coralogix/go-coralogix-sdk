package coralogix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHook_Send(t *testing.T) {
	CoralogixHook := NewCoralogixHook(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)
	defer func() { recover() }()
	defer CoralogixHook.Close()

	log := logrus.New()

	log.SetLevel(logrus.TraceLevel)
	log.AddHook(CoralogixHook)
	fields := logrus.Fields{
		"Category": "MyCategory",
		"ThreadId": "MyThreadId",
		"extra":    "additional",
	}

	testcases := []struct {
		name     string
		logfn    func(args ...interface{})
		severity uint
	}{
		{
			name:     "test trace",
			severity: Level.TRACE,
			logfn: func(args ...interface{}) {
				log.WithFields(fields).Trace(args...)
			},
		},
		{
			name:     "test debug",
			severity: Level.DEBUG,
			logfn: func(args ...interface{}) {
				log.WithFields(fields).Debug(args...)
			},
		},
		{
			name:     "test error",
			severity: Level.ERROR,
			logfn: func(args ...interface{}) {
				log.WithFields(fields).Error(args...)
			},
		},
		{
			name:     "test info",
			severity: Level.INFO,
			logfn: func(args ...interface{}) {
				log.WithFields(fields).Info(args...)
			},
		},
		{
			name:     "test warn",
			severity: Level.WARNING,
			logfn: func(args ...interface{}) {
				log.WithFields(fields).Warn(args...)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			msg := fmt.Sprintf("%s (%s)", tc.name, t.Name())
			tc.logfn(msg)
			time.Sleep(time.Duration(1) * time.Second)
			bulk, ok := mockHTTPServerMap[t.Name()]
			assert.True(t, ok, "%s key not found in mockHTTPServerMap", t.Name())

			var msgExists bool
			for _, entry := range bulk.LogEntries {
				if msgExists = strings.Contains(entry.Text, tc.name); msgExists {
					assert.Equal(t, tc.severity, entry.Severity)
					assert.Equal(t, fields["Category"], entry.Category)
					assert.Equal(t, fields["ThreadId"], entry.ThreadID,
						"expected %v, got %v", fields["MyThreadId"], entry.ThreadID)
					assert.True(t, strings.Contains(entry.Text, fields["extra"].(string)),
						"entry Text does not contain extra field", entry.Text, fields["extra"])
					break
				}
			}
			assert.True(t, msgExists, "no matching message found", string(bulk.ToJSON()))
		})
	}

	// test with caller
	log.SetReportCaller(true)
	for _, tc := range testcases {
		tc.name = fmt.Sprintf("%s_withReportCaller", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			msg := fmt.Sprintf("%s (%s)", tc.name, t.Name())
			tc.logfn(msg)
			time.Sleep(time.Duration(1) * time.Second)
			bulk, ok := mockHTTPServerMap[t.Name()]
			assert.True(t, ok, "%s key not found in mockHTTPServerMap", t.Name())

			var msgExists bool
			for _, entry := range bulk.LogEntries {
				if msgExists = strings.Contains(entry.Text, tc.name); msgExists {
					assert.Equal(t, tc.severity, entry.Severity)
					assert.Equal(t, fields["Category"], entry.Category)
					assert.Equal(t, fields["ThreadId"], entry.ThreadID,
						"expected %v, got %v", fields["ThreadId"], entry.ThreadID)
					assert.True(t, strings.Contains(entry.Text, fields["extra"].(string)),
						"entry Text does not contain extra field", entry.Text, fields["extra"])
					break
				}
			}
			assert.True(t, msgExists, "no matching message found", string(bulk.ToJSON()))
		})
	}
}
