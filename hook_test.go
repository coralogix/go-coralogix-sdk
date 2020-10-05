package coralogix

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestHook_Send(t *testing.T) {
	CoralogixHook := NewCoralogixHook(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer func() { recover() }()
	defer CoralogixHook.Close()

	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)

	log.AddHook(CoralogixHook)

	log.WithFields(logrus.Fields{
		"Category":   "MyCategory",
		"ClassName":  "MyClassName",
		"MethodName": "MyMethodName",
		"ThreadId":   "MyThreadId",
		"extra":      "additional",
	}).Debug("Test message!")

	log.Trace("Test trace message!")
	log.Debug("Test debug message!")
	log.Info("Test info message!")
	log.Warn("Test warn message!")
	log.Error("Test error message!")
	log.Panic("Test panic message!")
}

func TestHook_SendWithCaller(t *testing.T) {
	CoralogixHook := NewCoralogixHook(
		GetEnv(
			"PRIVATE_KEY",
			"7569303a-6269-4d2c-bf14-1aec9b1786a4",
		),
		"sdk-go",
		"test",
	)
	defer func() { recover() }()
	defer CoralogixHook.Close()

	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.TraceLevel)

	log.AddHook(CoralogixHook)

	log.WithFields(logrus.Fields{
		"Category": "MyCategory",
		"ThreadId": "MyThreadId",
		"extra":    "additional",
	}).Debug("Test message!")

	log.Trace("Test trace message with caller!")
	log.Debug("Test debug message with caller!")
	log.Info("Test info message with caller!")
	log.Warn("Test warn message with caller!")
	log.Error("Test error message with caller!")
	log.Panic("Test panic message with caller!")
}
