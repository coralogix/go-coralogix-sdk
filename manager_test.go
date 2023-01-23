package coralogix

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

type TestDummyStruct struct {
	Foo string
}

func TestNewLoggerManager(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	if reflect.TypeOf(NewLoggerManagerTestInstance) != reflect.TypeOf(&LoggerManager{}) ||
		!strings.Contains(NewLoggerManagerTestInstance.LogsBuffer.Slice(1)[0].Text, sdkVersion) {
		t.Error("CoralogixLogger manager creation failed!")
	}
}

func TestLoggerManager_LogsBufferSize(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Len() != 2 {
		t.Error("Failed to check logs buffer size!")
	}
}

func TestLoggerManager_LogsBufferLength(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Len() == 0 {
		t.Error("Failed to check logs buffer length!")
	}
}

func TestLoggerManager_LogsBufferLengthFail(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.LogsBuffer = LogBuffer{}
	NewLoggerManagerTestInstance.LogsBuffer.Append(*InvalidLogMessage())
	if NewLoggerManagerTestInstance.LogsBuffer.Size() > 0 {
		t.Error("Buffer length should fail due to incorrect content!")
	}
}

func TestLoggerManager_SendInitMessage(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.SendInitMessage()
	if !strings.Contains(NewLoggerManagerTestInstance.LogsBuffer.Slice(2)[1].Text, sdkVersion) {
		t.Error("Initial message sending failed!")
	}
}

func TestLoggerManager_AddLogLine(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Len() != 2 {
		t.Error("Failed to add log to buffer!")
	}
}

func TestLoggerManager_AddLogLineWithInvalidSeverity(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		13,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Slice(2)[1].Severity != Level.INFO {
		t.Error("Severity checking failed!")
	}
}

func TestLoggerManager_AddLogLineWithEmptyText(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Slice(2)[1].Text != "EMPTY_STRING" {
		t.Error("Log text checking failed!")
	}
}

func TestLoggerManager_AddLogLineWithNil(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		nil,
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Slice(2)[1].Text != "EMPTY_STRING" {
		t.Error("Log text checking failed!")
	}
}

func TestLoggerManager_AddLogLineWithEmptyCategory(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		"",
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Slice(2)[1].Category != LogCategory {
		t.Error("Log category checking failed!")
	}
}

func TestLoggerManager_AddLogLineOverflow(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		RandStringBytes(MaxLogChunkSize),
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.LogsBuffer.Len() > 1 {
		t.Error("Failed to check log record max length!")
	}
}

func TestLoggerManager_SendBulk(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerManagerTestInstance.SendBulk(true) != true ||
		NewLoggerManagerTestInstance.LogsBuffer.Len() > 0 {
		t.Error("Failed to check log record max length!")
	}
}

func TestLoggerManager_SendBulkWithEmptyBuffer(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.LogsBuffer = LogBuffer{}
	if NewLoggerManagerTestInstance.SendBulk(true) != false {
		t.Error("Unexpected behavior when sending empty buffer!")
	}
}

func TestLoggerManager_SendBulkWithBigBuffer(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	for i := 0; i < 3; i++ {
		NewLoggerManagerTestInstance.AddLogLine(
			Level.DEBUG,
			RandStringBytes(MaxLogChunkSize/2),
			LogCategory,
			"",
			"",
			"",
		)
	}
	NewLoggerManagerTestInstance.SendBulk(true)
	if NewLoggerManagerTestInstance.LogsBuffer.Len() == 0 {
		t.Error("Incorrect buffer chunk process!")
	}
}

func TestLoggerManager_Flush(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.AddLogLine(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	NewLoggerManagerTestInstance.Flush()
	if NewLoggerManagerTestInstance.LogsBuffer.Len() > 0 {
		t.Error("Logs buffer flush failed!")
	}
}

func TestLoggerManager_Stop(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	go NewLoggerManagerTestInstance.Run()
	NewLoggerManagerTestInstance.Stop()
	if NewLoggerManagerTestInstance.Stopped != true {
		t.Error("CoralogixLogger manager process stopping failed!")
	}
}

func TestMessageToString(t *testing.T) {
	Message := "Test message"
	if MessageToString(Message) != Message {
		t.Error("Failed to convert string message!")
	}
}

func TestMessageToStringFail(t *testing.T) {
	if MessageToString(make(chan string)) != "" {
		t.Error("Failed to catch error when converting invalid message!")
	}
}

func TestMessageToStringWithBytes(t *testing.T) {
	Message := []byte("Test message")
	if MessageToString(Message) != string(Message) {
		t.Error("Failed to convert bytes array message!")
	}
}

func TestMessageToStringWithInteger(t *testing.T) {
	Message := rand.Intn(100)
	if MessageToString(Message) != fmt.Sprintf("%d", Message) {
		t.Error("Failed to convert integer message!")
	}
}

func TestMessageToStringWithFloat(t *testing.T) {
	Message := math.Pi
	if MessageToString(Message) != fmt.Sprintf("%f", Message) {
		t.Error("Failed to convert float message!")
	}
}

func TestMessageToStringWithStruct(t *testing.T) {
	Message := TestDummyStruct{"Bar"}
	JSONMessage, _ := json.Marshal(Message)
	if MessageToString(Message) != string(JSONMessage) {
		t.Error("Failed to convert float message!")
	}
}

func TestLoggerManager_UpdateTimeDeltaInterval(t *testing.T) {
	NewLoggerManagerTestInstance := CreateLoggerManager()
	NewLoggerManagerTestInstance.UpdateTimeDeltaInterval()
	if NewLoggerManagerTestInstance.TimeDelta == 0 {
		t.Error("Time synchronization failed!")
	}
}

func CreateLoggerManager() *LoggerManager {
	return NewLoggerManager(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
		true,
	)
}

func RandStringBytes(Length uint64) string {
	const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RandomString := make([]byte, Length)
	for i := range RandomString {
		RandomString[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(RandomString)
}
