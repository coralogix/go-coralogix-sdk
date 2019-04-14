package go_coralogix_sdk

import "testing"

func TestNormalSendSpeedInterval(t *testing.T) {
	if NormalSendSpeedInterval <= 0 {
		t.Error("Bulk send interval in normal mode should be greater than 0!")
	}
}

func TestNormalFastSendSpeedInterval(t *testing.T) {
	if FastSendSpeedInterval <= 0 {
		t.Error("Bulk send interval in fast mode should be greater than 0!")
	}
}
