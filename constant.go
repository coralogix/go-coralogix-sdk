package coralogix

const (
	// MaxLogBufferSize is maximum log buffer size (default=128MiB)
	MaxLogBufferSize uint64 = 128 * (1024 * 1024)

	// MaxLogChunkSize is maximum chunk size (default=1.5MiB)
	MaxLogChunkSize uint64 = 1.5 * (1024 * 1024)

	// NormalSendSpeedInterval is a bulk send interval in normal mode
	NormalSendSpeedInterval float64 = 0.5

	// FastSendSpeedInterval is a bulk send interval in fast mode
	FastSendSpeedInterval float64 = 0.1

	// TimeDelayTimeout is a timeout for time-delay request
	TimeDelayTimeout uint = 5

	// FailedPrivateKey is a default private key
	FailedPrivateKey string = "no private key"

	// NoAppName is a default application name
	NoAppName string = "NO_APP_NAME"

	// NoSubSystem is a default subsystem name
	NoSubSystem string = "NO_SUB_NAME"

	// HttpTimeout is a default HTTP timeout
	HttpTimeout uint = 30

	// HttpSendRetryCount is a number of attempts to retry HTTP request
	HttpSendRetryCount uint = 5

	// HttpSendRetryInterval is a interval between failed http post requests
	HttpSendRetryInterval uint = 2

	// LogCategory is a default category for log record
	LogCategory string = "CORALOGIX"

	// SyncTimeUpdateInterval is a time synchronization interval (in minutes)
	SyncTimeUpdateInterval uint = 5
)

var (
	// LogUrl is the Coralogix logs url endpoint
	LogUrl string = GetEnv("CORALOGIX_LOG_URL", "https://api.coralogix.com:443/api/v1/logs")

	// TimeDeltaUrl is the Coralogix time delay url endpoint
	TimeDeltaUrl string = GetEnv("CORALOGIX_TIME_DELTA_URL", "https://api.coralogix.com:443/sdk/v1/time")
)
