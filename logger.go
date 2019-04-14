package go_coralogix_sdk

import "runtime"

// CoralogixLogger is interface for using SDK
type CoralogixLogger struct {
	Category      string        // Current logger logs records category
	LoggerManager LoggerManager // CoralogixLogger manager instance
}

// NewCoralogixLogger initialize new SDK interface instance
func NewCoralogixLogger(PrivateKey string, ApplicationName string, SubsystemName string) *CoralogixLogger {
	return NewCoralogixLoggerWithCategory(
		PrivateKey,
		ApplicationName,
		SubsystemName,
		LogCategory,
	)
}

// NewCoralogixLoggerWithCategory initialize new SDK interface instance with custom category
func NewCoralogixLoggerWithCategory(PrivateKey string, ApplicationName string, SubsystemName string, Category string) *CoralogixLogger {
	if PrivateKey == "" {
		PrivateKey = FailedPrivateKey
	}

	if ApplicationName == "" {
		ApplicationName = NoAppName
	}

	if SubsystemName == "" {
		SubsystemName = NoSubSystem
	}

	if Category == "" {
		Category = LogCategory
	}

	NewLoggerInstance := &CoralogixLogger{
		Category,
		*NewLoggerManager(
			PrivateKey,
			ApplicationName,
			SubsystemName,
			true,
		),
	}

	go NewLoggerInstance.LoggerManager.Run()

	runtime.SetFinalizer(NewLoggerInstance, func(logger *CoralogixLogger) { logger.Destroy() })

	return NewLoggerInstance
}

// Destroy stop logger manager and cleanup logs buffer before exit
func (logger *CoralogixLogger) Destroy() {
	logger.LoggerManager.Stop()
}

// Log send record message to logger manager
func (logger *CoralogixLogger) Log(Severity uint, Text interface{}, Category string, ClassName string, MethodName string, ThreadId string) {
	if Category == "" {
		Category = logger.Category
	}
	logger.LoggerManager.AddLogLine(Severity, Text, Category, ClassName, MethodName, ThreadId)
}

// Debug send log message with DEBUG severity level
func (logger *CoralogixLogger) Debug(Text interface{}) {
	logger.Log(Level.DEBUG, Text, "", "", "", "")
}

// Verbose send log message with VERBOSE severity level
func (logger *CoralogixLogger) Verbose(Text interface{}) {
	logger.Log(Level.VERBOSE, Text, "", "", "", "")
}

// Info send log message with INFO severity level
func (logger *CoralogixLogger) Info(Text interface{}) {
	logger.Log(Level.INFO, Text, "", "", "", "")
}

// Warning send log message with WARNING severity level
func (logger *CoralogixLogger) Warning(Text interface{}) {
	logger.Log(Level.WARNING, Text, "", "", "", "")
}

// Error send log message with ERROR severity level
func (logger *CoralogixLogger) Error(Text interface{}) {
	logger.Log(Level.ERROR, Text, "", "", "", "")
}

// Critical send log message with CRITICAL severity level
func (logger *CoralogixLogger) Critical(Text interface{}) {
	logger.Log(Level.CRITICAL, Text, "", "", "", "")
}
