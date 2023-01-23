package coralogix

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// LoggerManager is a logs buffer operations manager
type LoggerManager struct {
	SyncTime            bool           // Synchronize time with Coralogix servers
	TimeDelta           float64        // Time difference between local machine and Coralogix servers
	TimeDeltaLastUpdate int            // Last time-delta update time
	Stopped             bool           // Is current logger manager stopped
	LogsBuffer          LogBuffer      // Logs buffer
	Credentials                        // Credentials for Coralogix account
	Lock                sync.WaitGroup // CoralogixLogger manager locker
}

// NewLoggerManager configure new logger manager instance
func NewLoggerManager(PrivateKey string, ApplicationName string, SubsystemName string, SyncTime bool) *LoggerManager {
	LoggerManagerInstance := &LoggerManager{
		SyncTime,
		0,
		0,
		false,
		LogBuffer{},
		Credentials{
			PrivateKey,
			ApplicationName,
			SubsystemName,
		},
		sync.WaitGroup{},
	}

	LoggerManagerInstance.Lock.Add(1)

	LoggerManagerInstance.SendInitMessage()
	return LoggerManagerInstance
}

// SendInitMessage send initialization message to Coralogix for connection verify
func (manager *LoggerManager) SendInitMessage() {
	manager.AddLogLine(
		Level.INFO,
		fmt.Sprintf(
			"The Application Name %s and Subsystem Name %s from the Go SDK, version %s has started to send data",
			manager.ApplicationName,
			manager.SubsystemName,
			sdkVersion,
		),
		LogCategory,
		"",
		"",
		"",
	)
}

// AddLogLine push log record to buffer
func (manager *LoggerManager) AddLogLine(Severity uint, Text interface{}, Category string, ClassName string, MethodName string, ThreadID string) {
	if manager.LogsBuffer.Size() < MaxLogBufferSize {
		if Severity < Level.DEBUG || Severity > Level.CRITICAL {
			Severity = Level.INFO
		}
		if Text == nil || Text == "" {
			Text = "EMPTY_STRING"
		}
		if Category == "" {
			Category = LogCategory
		}

		NewLogRecord := Log{
			float64(time.Now().Unix())*1000.0 + manager.TimeDelta,
			Severity,
			MessageToString(Text),
			Category,
			ClassName,
			MethodName,
			ThreadID,
			0,
		}

		if MaxLogChunkSize <= NewLogRecord.Size() {
			DebugLogger.Printf(
				"AddLogLine(): received log message too big of size= %d MB, bigger than max_log_chunk_size= %d; throwing...\n",
				NewLogRecord.Size()/(1024*1024),
				MaxLogChunkSize,
			)
			return
		}

		manager.LogsBuffer.Append(NewLogRecord)
	}
}

// SendBulk send logs bulk to Coralogix
func (manager *LoggerManager) SendBulk(SyncTime bool) bool {
	if SyncTime {
		manager.UpdateTimeDeltaInterval()
	}

	BufferLenToSend := manager.LogsBuffer.Len()
	if BufferLenToSend < 1 {
		DebugLogger.Println("buffer is empty, there is nothing to send!")
		return false
	}

	for manager.LogsBuffer.Size() > MaxLogChunkSize && BufferLenToSend > 1 {
		BufferLenToSend = BufferLenToSend / 2
	}

	if BufferLenToSend < 1 {
		BufferLenToSend = 1
	}

	DebugLogger.Println("Checking buffer size. Total log entries is:", BufferLenToSend)
	LogsBulk := NewBulk(manager.Credentials)
	for _, Record := range manager.LogsBuffer.Slice(BufferLenToSend) {
		LogsBulk.AddRecord(Record)
	}

	SendRequest(LogsBulk)
	return true
}

// Run should work in separate thread and asynchronously operate with logs
func (manager *LoggerManager) Run() {
	var NextSendInterval float64

	defer manager.Lock.Done()

	for {
		if manager.Stopped {
			manager.Flush()
			return
		}

		manager.SendBulk(manager.SyncTime)

		if manager.LogsBuffer.Size() > (MaxLogChunkSize / 2) {
			NextSendInterval = FastSendSpeedInterval
		} else {
			NextSendInterval = NormalSendSpeedInterval
		}

		DebugLogger.Printf("Next buffer check is scheduled in %.1f seconds\n", NextSendInterval)
		time.Sleep(time.Duration(NextSendInterval) * time.Second)
	}
}

// Flush clean buffer and send logs to Coralogix
func (manager *LoggerManager) Flush() {
	manager.SendBulk(false)
}

// Stop logger manager and kill threaded agent
func (manager *LoggerManager) Stop() {
	manager.Stopped = true
	manager.Lock.Wait()
}

// MessageToString convert log content to simple string
func MessageToString(Message interface{}) string {
	switch Text := Message.(type) {
	case []byte:
		return string(Text)
	case string:
		return Text
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", Text)
	case float32, float64:
		return fmt.Sprintf("%f", Text)
	default:
		JSONMessage, err := json.Marshal(Text)
		if err != nil {
			return ""
		}
		return string(JSONMessage)
	}
}

// UpdateTimeDeltaInterval get time difference between local machine and Coralogix servers
func (manager *LoggerManager) UpdateTimeDeltaInterval() {
	if (uint(int(time.Now().Unix()) - manager.TimeDeltaLastUpdate)) >= 60*SyncTimeUpdateInterval {
		Result, TimeDelta := GetTimeSync()
		if Result {
			manager.TimeDelta = TimeDelta
			manager.TimeDeltaLastUpdate = int(time.Now().Unix())
		} else {
			DebugLogger.Println("Time synchronization was unsuccessful!")
		}
	}
}
