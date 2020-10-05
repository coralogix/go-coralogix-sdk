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
	LogsBuffer          []Log          // Logs buffer
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
		[]Log{},
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

// LogsBufferLength calculate buffer length in bytes
func (manager *LoggerManager) LogsBufferLength() uint64 {
	JSONBuffer, err := json.Marshal(manager.LogsBuffer)
	if err != nil {
		DebugLogger.Println("Can't convert to JSON: ", err)
		return 0
	}
	return uint64(len(JSONBuffer))
}

// LogsBufferSize return buffer entries count
func (manager *LoggerManager) LogsBufferSize() int {
	return len(manager.LogsBuffer)
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
	if manager.LogsBufferLength() < MaxLogBufferSize {
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
		}

		if MaxLogChunkSize <= uint64(NewLogRecord.Size()) {
			DebugLogger.Printf(
				"AddLogLine(): received log message too big of size= %d MB, bigger than max_log_chunk_size= %d; throwing...\n",
				NewLogRecord.Size()/(1024*1024),
				MaxLogChunkSize,
			)
			return
		}

		manager.LogsBuffer = append(manager.LogsBuffer, NewLogRecord)
	}
}

// SendBulk send logs bulk to Coralogix
func (manager *LoggerManager) SendBulk(SyncTime bool) bool {
	if SyncTime {
		manager.UpdateTimeDeltaInterval()
	}

	BufferSizeToSend := manager.LogsBufferSize()
	if BufferSizeToSend < 1 {
		DebugLogger.Println("Buffer is empty, there is nothing to send!")
		return false
	}

	for manager.LogsBufferLength() > MaxLogChunkSize && BufferSizeToSend > 1 {
		BufferSizeToSend = int(BufferSizeToSend / 2)
	}

	if BufferSizeToSend < 1 {
		BufferSizeToSend = 1
	}

	DebugLogger.Println("Checking buffer size. Total log entries is:", BufferSizeToSend)
	LogsBulk := NewBulk(manager.Credentials)
	for _, Record := range manager.LogsBuffer[:BufferSizeToSend] {
		LogsBulk.AddRecord(Record)
	}
	manager.LogsBuffer = manager.LogsBuffer[BufferSizeToSend:]
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

		if manager.LogsBufferLength() > (MaxLogChunkSize / 2) {
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
