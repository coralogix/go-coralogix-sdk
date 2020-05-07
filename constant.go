package coralogix

const (
	MaxLogBufferSize        uint64  = 128 * (1024 * 1024)                         // Maximum log buffer size (default=128MiB)
	MaxLogChunkSize         uint64  = 1.5 * (1024 * 1024)                         // Maximum chunk size (default=1.5MiB)
	NormalSendSpeedInterval float64 = 0.5                                         // Bulk send interval in normal mode
	FastSendSpeedInterval   float64 = 0.1                                         // Bulk send interval in fast mode
	LogUrl                  string  = "https://api.coralogix.com:443/api/v1/logs" // Coralogix logs url endpoint
	TimeDeltaUrl            string  = "https://api.coralogix.com:443/sdk/v1/time" // Coralogix time delay url endpoint
	TimeDelayTimeout        uint    = 5                                           // Timeout for time-delay request
	FailedPrivateKey        string  = "no private key"                            // Default private key
	NoAppName               string  = "NO_APP_NAME"                               // Default application name
	NoSubSystem             string  = "NO_SUB_NAME"                               // Default subsystem name
	HttpTimeout             uint    = 30                                          // Default HTTP timeout
	HttpSendRetryCount      uint    = 5                                           // Number of attempts to retry HTTP request
	HttpSendRetryInterval   uint    = 2                                           // Interval between failed http post requests
	LogCategory             string  = "CORALOGIX"                                 // Default category for log record
	SyncTimeUpdateInterval  uint    = 5                                           // Time synchronization interval (in minutes)
)
