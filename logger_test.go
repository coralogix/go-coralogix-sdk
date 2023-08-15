package coralogix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPrivateKey = "7569303a-6269-4d2c-bf14-1aec9b1786a4"
	startupMessage = "The Application Name sdk-go and Subsystem Name test from the Go SDK, version 1.0.2 has started to send data"
)

// testRequestChan = a channel used to pass requests
// that have been sent to the mockHTTPServer
var (
	cxlogger *CoralogixLogger
	// mockHTTPServerMap - store Bulk requests sent to mock
	// endpoint
	mockHTTPServerMap map[string]*Bulk
	mockKeyRegex      = regexp.MustCompile(`\((.*)\)`)
)

// mockHTTPServer - stores all requests received in a map
// for validation
func mockHTTPServer(w http.ResponseWriter, r *http.Request) {
	var bulk Bulk
	if err := json.NewDecoder(r.Body).Decode(&bulk); err != nil {
		log.Println("unable to unmarshal test message:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(bulk.LogEntries) == 1 && bulk.LogEntries[0].Text == startupMessage {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, entry := range bulk.LogEntries {
		if matches := mockKeyRegex.FindStringSubmatch(entry.Text); len(matches) == 2 {
			k := matches[1]
			mockHTTPServerMap[k] = &bulk
		}
	}

	defer r.Body.Close()
	sc := http.StatusOK

	w.WriteHeader(sc)
}

func TestMain(m *testing.M) {
	mockHTTPServerMap = make(map[string]*Bulk)
	server := httptest.NewServer(http.HandlerFunc(mockHTTPServer))
	LogURL = server.URL
	cxlogger = NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)

	code := m.Run()

	cxlogger.Destroy()
	server.Close()
	os.Exit(code)
}

func TestNewLogger(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t,
		reflect.TypeOf(NewLoggerTestInstance),
		reflect.TypeOf(&CoralogixLogger{}))
}

func TestNewLoggerWithoutPrivateKey(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		"",
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t, FailedPrivateKey, NewLoggerTestInstance.LoggerManager.PrivateKey)
}

func TestNewLoggerWithoutApplicationName(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t, NoAppName, NewLoggerTestInstance.LoggerManager.ApplicationName)
}

func TestNewLoggerWithoutSubsystemName(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"",
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t, NoSubSystem, NewLoggerTestInstance.LoggerManager.SubsystemName)
}

func TestNewLoggerWithCategory(t *testing.T) {
	TestCategory := "test"
	NewLoggerTestInstance := NewCoralogixLoggerWithCategory(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
		TestCategory,
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t, TestCategory, NewLoggerTestInstance.Category)
}

func TestNewLoggerWithEmptyCategory(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLoggerWithCategory(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
		"",
	)
	defer NewLoggerTestInstance.Destroy()
	assert.Equal(t, LogCategory, NewLoggerTestInstance.Category)
}

func TestLoggerSeverity(t *testing.T) {
	testcases := []struct {
		name     string
		clogger  *CoralogixLogger
		match    string
		logfn    func(*CoralogixLogger)
		severity uint
	}{
		{
			name:  "debug",
			match: "debug",
			logfn: func(logger *CoralogixLogger) {
				logger.Debug(fmt.Sprintf("test message debug (%s/debug)", t.Name()))
			},
			severity: Level.DEBUG,
		},
		{
			name:  "verbose",
			match: "verbose",
			// clogger: newtestCoralogixLoggerFn(),
			logfn: func(logger *CoralogixLogger) {
				logger.Verbose(fmt.Sprintf("test message verbose (%s/verbose)", t.Name()))
			},
			severity: Level.VERBOSE,
		},
		{
			name: "info",
			// clogger: newtestCoralogixLoggerFn(),
			match: "info",
			logfn: func(logger *CoralogixLogger) {
				logger.Info(fmt.Sprintf("test message info (%s/info)", t.Name()))
			},
			severity: Level.INFO,
		},
		{
			name:  "warning",
			match: "warning",
			logfn: func(logger *CoralogixLogger) {
				logger.Warning(fmt.Sprintf("test message warning (%s/warning)", t.Name()))
			},
			severity: Level.WARNING,
		},
		{
			name:  "error",
			match: "error",
			logfn: func(logger *CoralogixLogger) {
				logger.Error(fmt.Sprintf("test message error (%s/error)", t.Name()))
			},
			severity: Level.ERROR,
		},
		{
			name:  "critical",
			match: "critical",
			// clogger: newtestCoralogixLoggerFn(),
			logfn: func(logger *CoralogixLogger) {
				logger.Critical(fmt.Sprintf("test message debug (%s/critical)", t.Name()))
			},
			severity: Level.CRITICAL,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logfn(cxlogger)
			// time.Sleep(time.Second * 1)
			cxlogger.LoggerManager.Flush()

			bulk, ok := mockHTTPServerMap[t.Name()]
			assert.True(t, ok, "%s key not found in mockHTTPServerMap", t.Name())

			var msg = true
			for _, entry := range bulk.LogEntries {
				if msg = strings.Contains(strings.ToLower(entry.Text), tc.name); msg {
					assert.Equal(t, tc.severity, entry.Severity)
					break
				}
			}

			assert.True(t, msg, "no matching message found")
		})
	}
}
