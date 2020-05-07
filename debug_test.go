package coralogix

import (
	"testing"
)

func TestSetDebugEnable(t *testing.T) {
	SetDebug(true)
}

func TestSetDebugDisable(t *testing.T) {
	SetDebug(false)
}
