package coralogix

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VARIABLE", "1")
	if GetEnv("TEST_VARIABLE", "2") != "1" {
		t.Error("Function should return value of existing variable")
	}
}

func TestGetEnvDefault(t *testing.T) {
	os.Unsetenv("TEST_VARIABLE")
	if GetEnv("TEST_VARIABLE", "2") != "2" {
		t.Error("Function should return default value for variable")
	}
}
