package coralogix

import (
	"math"
	"testing"
)

func TestLog_Size(t *testing.T) {
	LogRecord := Log{
		1554616361,
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
		0,
	}
	if LogRecord.Size() == 0 {
		t.Error("Invalid log record size calculation!")
	}
}

func TestLog_SizeFail(t *testing.T) {
	if InvalidLogMessage().Size() > 0 {
		t.Error("Log size should be equals to 0 due to JSON parsing error!")
	}
}

func InvalidLogMessage() *Log {
	return &Log{
		math.Inf(1),
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
		0,
	}
}
