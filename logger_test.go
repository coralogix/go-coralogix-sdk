package coralogix

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPrivateKey = "7569303a-6269-4d2c-bf14-1aec9b1786a4"
)

// testRequestChan = a channel used to pass requests
// that have been sent to the mockHTTPServer
var testRequestChan chan *Bulk

// mockHTTPServer - stores all requests received in a map
// for validation
func mockHTTPServer(w http.ResponseWriter, r *http.Request) {
	var bulk Bulk
	if err := json.NewDecoder(r.Body).Decode(&bulk); err != nil {
		log.Fatalf("unable to unmarshal test message: %s", err)
	}

	testRequestChan <- &bulk
	defer r.Body.Close()
	sc := http.StatusOK
	w.WriteHeader(sc)
}

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(mockHTTPServer))
	testRequestChan = make(chan *Bulk, 1)
	LogURL = server.URL
	code := m.Run()
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

func TestLogger_Destroy(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)
	NewLoggerTestInstance.Destroy()
}

func TestLogger_Log(t *testing.T) {
	clogger := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)
	defer clogger.Destroy()
	clogger.Log(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)

	clogger.LoggerManager.Flush()
	bulk := <-testRequestChan

	var msg bool
	for _, entry := range bulk.LogEntries {
		msg = entry.Text == "Test message"
	}
	assert.True(t, msg, "no matching message found")
}

func TestLoggerSeverity(t *testing.T) {
	var cxlogger = NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			testPrivateKey,
		),
		"sdk-go",
		"test",
	)

	defer cxlogger.Destroy()

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
				logger.Debug("Test message DEBUG")
			},
			severity: Level.DEBUG,
		},
		{
			name:  "verbose",
			match: "verbose",
			// clogger: newtestCoralogixLoggerFn(),
			logfn: func(logger *CoralogixLogger) {
				logger.Verbose("Test message VERBOSE")
			},
			severity: Level.VERBOSE,
		},
		{
			name: "info",
			// clogger: newtestCoralogixLoggerFn(),
			match: "info",
			logfn: func(logger *CoralogixLogger) {
				logger.Info("Test message INFO")
			},
			severity: Level.INFO,
		},
		{
			name:  "warning",
			match: "warning",
			logfn: func(logger *CoralogixLogger) {
				logger.Warning("Test message WARNING")
			},
			severity: Level.WARNING,
		},
		{
			name:  "error",
			match: "error",
			logfn: func(logger *CoralogixLogger) {
				logger.Error("Test message ERROR")
			},
			severity: Level.ERROR,
		},
		{
			name:  "critical",
			match: "critical",
			// clogger: newtestCoralogixLoggerFn(),
			logfn: func(logger *CoralogixLogger) {
				logger.Critical("Test message CRITICAL")
			},
			severity: Level.CRITICAL,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logfn(cxlogger)
			cxlogger.LoggerManager.Flush()
			bulkReq := <-testRequestChan

			var msg = true
			for _, entry := range bulkReq.LogEntries {
				// fmt.Println(i, entry.Text)
				if msg = strings.Contains(strings.ToLower(entry.Text), tc.name); msg {
					assert.Equal(t, tc.severity, entry.Severity)
				}
			}

			assert.True(t, msg, "no matching message found")
		})
	}
}
