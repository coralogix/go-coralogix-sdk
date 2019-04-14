package go_coralogix_sdk

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSetDebugEnable(t *testing.T) {
	SetDebug(true)
	defer SetDebug(false)
	if DebugLogger.Writer() != os.Stdout {
		t.Error("Can't enable internal logging!")
	}
}

func TestSetDebugDisable(t *testing.T) {
	SetDebug(false)
	if DebugLogger.Writer() != ioutil.Discard {
		t.Error("Can't disable internal logging!")
	}
}
