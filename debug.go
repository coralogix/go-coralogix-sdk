package coralogix

import (
	"io/ioutil"
	"log"
	"os"
)

// DebugLogger is the internal logger (disabled by default)
var DebugLogger = log.New(ioutil.Discard, "CORALOGIX: ", log.Ldate|log.Ltime)

// SetDebug enable/disable internal logger
func SetDebug(Status bool) {
	if Status {
		DebugLogger.SetOutput(os.Stdout)
	} else {
		DebugLogger.SetOutput(ioutil.Discard)
	}
}
