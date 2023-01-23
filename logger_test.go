package coralogix

import (
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	if reflect.TypeOf(NewLoggerTestInstance) != reflect.TypeOf(&CoralogixLogger{}) {
		t.Error("CoralogixLogger creation failed!")
	}
}

func TestNewLoggerWithoutPrivateKey(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		"",
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	if NewLoggerTestInstance.LoggerManager.PrivateKey != FailedPrivateKey {
		t.Error("Invalid default value for incorrect private key!")
	}
}

func TestNewLoggerWithoutApplicationName(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	if NewLoggerTestInstance.LoggerManager.ApplicationName != NoAppName {
		t.Error("Invalid default value for empty application name!")
	}
}

func TestNewLoggerWithoutSubsystemName(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"",
	)
	defer NewLoggerTestInstance.Destroy()
	if NewLoggerTestInstance.LoggerManager.SubsystemName != NoSubSystem {
		t.Error("Invalid default value for empty subsystem name!")
	}
}

func TestNewLoggerWithCategory(t *testing.T) {
	TestCategory := "test"
	NewLoggerTestInstance := NewCoralogixLoggerWithCategory(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
		TestCategory,
	)
	defer NewLoggerTestInstance.Destroy()
	if NewLoggerTestInstance.Category != TestCategory {
		t.Error("Invalid logger category!")
	}
}

func TestNewLoggerWithEmptyCategory(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLoggerWithCategory(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
		"",
	)
	defer NewLoggerTestInstance.Destroy()
	if NewLoggerTestInstance.Category != LogCategory {
		t.Error("Category checking is not working!")
	}
}

func TestLogger_Destroy(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	NewLoggerTestInstance.Destroy()
}

func TestLogger_Log(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.Log(
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
	)
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 {
		t.Error("New log message add failed!")
	}
}

func TestLogger_Debug(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Debug("Test debug message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.DEBUG {
		t.Error("Debug log message add failed!")
	}
}

func TestLogger_Verbose(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Verbose("Test verbose message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.VERBOSE {
		t.Error("Verbose log message add failed!")
	}
}

func TestLogger_Info(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Info("Test info message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.INFO {
		t.Error("Info log message add failed!")
	}
}

func TestLogger_Warning(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Warning("Test warning message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.WARNING {
		t.Error("Warning log message add failed!")
	}
}

func TestLogger_Error(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Error("Test error message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.ERROR {
		t.Error("Error log message add failed!")
	}
}

func TestLogger_Critical(t *testing.T) {
	NewLoggerTestInstance := NewCoralogixLogger(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer NewLoggerTestInstance.Destroy()
	NewLoggerTestInstance.LoggerManager.Stop()
	NewLoggerTestInstance.Critical("Test critical message")
	if NewLoggerTestInstance.LoggerManager.LogsBuffer.Len() < 1 ||
		NewLoggerTestInstance.LoggerManager.LogsBuffer.Slice(1)[0].Severity != Level.CRITICAL {
		t.Error("Critical log message add failed!")
	}
}
