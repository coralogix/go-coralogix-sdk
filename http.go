package coralogix

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

// SendRequest send logs data to Coralogix server
func SendRequest(Bulk *Bulk) int {
	client := &http.Client{
		Timeout: time.Duration(HTTPTimeout) * time.Second,
	}

	for {
		statusCode, err := func() (int, error) {
			DebugLogger.Println("About to send bulk to Coralogix server. Attempt number:", Retry.Attempt())
			request, err := http.NewRequest(http.MethodPost, LogURL, bytes.NewBuffer(Bulk.ToJSON()))
			if err != nil {
				DebugLogger.Println("Can't create HTTP request:", err)
				return 0, err
			}
			request.Header = Headers

			response, err := client.Do(request)
			if err != nil {
				DebugLogger.Println("Can't execute HTTP request:", err)
				return 0, err
			}

			if response.StatusCode != 200 {
				DebugLogger.Println("HTTP requests was failed with code:", response.StatusCode)
			} else {
				DebugLogger.Println("Successfully sent bulk to Coralogix server. Result is:", response.StatusCode)
			}
			return response.StatusCode, nil
		}()
		if err == nil {
			return statusCode
		}

		if duration, retry := Retry.RetryDelay(); retry {
			time.Sleep(duration)
		} else {
			return 0
		}
	}
}

// GetTimeSync synchronize logs time with Coralogix servers time
func GetTimeSync() (bool, float64) {
	DebugLogger.Println("Syncing time with Coralogix server...")

	client := &http.Client{
		Timeout: time.Duration(TimeDelayTimeout) * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, TimeDeltaURL, nil)
	if err != nil {
		DebugLogger.Println("Can't create HTTP request:", err)
		return false, 0
	}
	request.Header = Headers

	response, err := client.Do(request)
	if err != nil {
		DebugLogger.Println("Can't execute HTTP request:", err)
		return false, 0
	}

	if response.StatusCode == 200 {
		response, _ := io.ReadAll(response.Body)
		ServerTime, err := strconv.ParseFloat(string(response), 64)

		if err != nil {
			DebugLogger.Println("Can't parse HTTP response:", err)
			return false, 0
		}

		ServerTime = ServerTime / 1e4
		TimeDelta := ServerTime - float64(time.Now().UnixMilli())

		return true, TimeDelta * 1e3 // convert to microseconds, because log timestamp is in microseconds.
	}

	DebugLogger.Println("Can't get server time")
	return false, 0
}
