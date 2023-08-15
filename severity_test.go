package coralogix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeveritiesList(t *testing.T) {
	assert.Equal(t, uint(1), Level.TRACE, "TRACE level should equals to 1!")
	assert.Equal(t, uint(1), Level.DEBUG, "DEBUG level should equals to 1!")
	assert.Equal(t, uint(2), Level.VERBOSE, "VERBOSE level should equals to 2!")
	assert.Equal(t, uint(3), Level.INFO, "INFO level should equals to 3!")
	assert.Equal(t, uint(4), Level.WARNING, "WARNING level should equals to 4!")
	assert.Equal(t, uint(5), Level.ERROR, "ERROR level should equals to 5!")
	assert.Equal(t, uint(6), Level.CRITICAL, "CRITICAL level should equals to 6!")
	assert.Equal(t, uint(6), Level.FATAL, "FATAL level should equals to 6!")
	assert.Equal(t, uint(6), Level.PANIC, "PANIC level should equals to 6!")
}
