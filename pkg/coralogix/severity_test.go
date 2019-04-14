package coralogix

import "testing"

func TestSeveritiesList(t *testing.T) {
	if Level.TRACE != 1 {
		t.Error("TRACE level should equals to 1!")
	}
	if Level.DEBUG != 1 {
		t.Error("DEBUG level should equals to 1!")
	}
	if Level.VERBOSE != 2 {
		t.Error("VERBOSE level should equals to 2!")
	}
	if Level.INFO != 3 {
		t.Error("INFO level should equals to 3!")
	}
	if Level.WARNING != 4 {
		t.Error("WARNING level should equals to 4!")
	}
	if Level.ERROR != 5 {
		t.Error("ERROR level should equals to 5!")
	}
	if Level.CRITICAL != 6 {
		t.Error("CRITICAL level should equals to 6!")
	}
	if Level.FATAL != 6 {
		t.Error("FATAL level should equals to 6!")
	}
	if Level.PANIC != 6 {
		t.Error("PANIC level should equals to 6!")
	}
}
