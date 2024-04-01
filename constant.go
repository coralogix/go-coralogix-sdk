package coralogix

import (
	"bufio"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

const (
	// MaxLogBufferSize is maximum log buffer size (default=128MiB)
	MaxLogBufferSize uint64 = 128 * (1024 * 1024)

	// MaxLogChunkSize is maximum chunk size (default=1.5MiB)
	MaxLogChunkSize uint64 = 1.5 * (1024 * 1024)

	// NormalSendSpeedInterval is a bulk send interval in normal mode
	NormalSendSpeedInterval = 1 * time.Second

	// FastSendSpeedInterval is a bulk send interval in fast mode
	FastSendSpeedInterval = 500 * time.Millisecond

	// TimeDelayTimeout is a timeout for time-delay request
	TimeDelayTimeout uint = 5

	// FailedPrivateKey is a default private key
	FailedPrivateKey string = "no private key"

	// NoAppName is a default application name
	NoAppName string = "NO_APP_NAME"

	// NoSubSystem is a default subsystem name
	NoSubSystem string = "NO_SUB_NAME"

	// HTTPTimeout is a default HTTP timeout
	HTTPTimeout uint = 30

	// HTTPSendRetryCount is a number of attempts to retry HTTP request
	HTTPSendRetryCount uint = 5

	// HTTPSendRetryInterval is a interval between failed http post requests
	HTTPSendRetryInterval uint = 2

	// LogCategory is a default category for log record
	LogCategory string = "CORALOGIX"

	// SyncTimeUpdateInterval is a time synchronization interval (in minutes)
	SyncTimeUpdateInterval uint = 5
)

var (
	// LogURL is the Coralogix logs url endpoint
	LogURL string = GetEnv("CORALOGIX_LOG_URL", "https://api.coralogix.com:443/api/v1/logs")

	// TimeDeltaURL is the Coralogix time delay url endpoint
	TimeDeltaURL string = GetEnv("CORALOGIX_TIME_DELTA_URL", "https://api.coralogix.com:443/sdk/v1/time")

	// Headers is the list of headers added to each send logs request
	Headers http.Header = func() http.Header {
		headers := GetEnv("CORALOGIX_HEADERS", "")
		tp := textproto.NewReader(bufio.NewReader(strings.NewReader(headers)))
		mimeHeader, err := tp.ReadMIMEHeader()
		if err != nil {
			mimeHeader = map[string][]string{}
		}
		mimeHeader.Set("Content-Type", "application/json")
		return http.Header(mimeHeader)
	}()
)
