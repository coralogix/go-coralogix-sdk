package coralogix

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// SendRequest send logs data to Coralogix server
func SendRequest(Bulk *Bulk) int {
	client := &http.Client{
		Timeout: time.Duration(time.Duration(HttpTimeout) * time.Second),
	}

	for Attempt := 1; uint(Attempt) <= HttpSendRetryCount; Attempt++ {
		DebugLogger.Println("About to send bulk to Coralogix server. Attempt number:", Attempt)

		response, err := client.Post(LogUrl, "application/json", bytes.NewBuffer(Bulk.ToJSON()))

		if err != nil {
			DebugLogger.Println("Can't execute HTTP request:", err)
			continue
		}

		if response.StatusCode != 200 {
			DebugLogger.Println("HTTP requests was failed with code:", response.StatusCode)
		} else {
			DebugLogger.Println("Successfully sent bulk to Coralogix server. Result is:", response.StatusCode)
			return response.StatusCode
		}

		time.Sleep(time.Duration(HttpSendRetryInterval) * time.Second)
	}

	return 0
}

// GetTimeSync synchronize logs time with Coralogix servers time
func GetTimeSync() (bool, float64) {
	DebugLogger.Println("Syncing time with Coralogix server...")

	client := &http.Client{
		Timeout: time.Duration(time.Duration(TimeDelayTimeout) * time.Second),
	}

	response, err := client.Get(TimeDeltaUrl)

	if err != nil {
		DebugLogger.Println("Can't execute HTTP request:", err)
		return false, 0
	}

	if response.StatusCode == 200 {
		response, _ := ioutil.ReadAll(response.Body)
		ServerTime, err := strconv.ParseFloat(string(response), 64)

		if err != nil {
			DebugLogger.Println("Can't parse HTTP response:", err)
			return false, 0
		}

		ServerTime = ServerTime / 1e4
		LocalTime := float64(time.Now().Unix() * 1e3)
		TimeDelta := ServerTime - LocalTime

		return true, TimeDelta
	}

	DebugLogger.Println("Can't get server time")
	return false, 0
}
