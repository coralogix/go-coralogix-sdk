package coralogix

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBulk(t *testing.T) {
	assert.IsType(t, &Bulk{}, CreateBulk(), "Bulk creation failed!")
}

func TestBulk_AddRecord(t *testing.T) {
	RecordsBulk := CreateBulk()
	RecordsBulk.AddRecord(Log{
		float64(time.Now().Unix()) * 1000.0,
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
		0,
	})

	assert.True(t, len(RecordsBulk.LogEntries) >= 1, "Adding new record to bulk failed!")
}

func TestBulk_ToJSON(t *testing.T) {
	RecordsBulkJSON := CreateBulk().ToJSON()
	assert.NotNil(t, RecordsBulkJSON)
	assert.IsType(t, []byte(""), RecordsBulkJSON, "Error while converting bulk to JSON!")

}

func TestBulk_ToJSONFail(t *testing.T) {
	RecordsBulk := CreateBulk()
	RecordsBulk.AddRecord(*InvalidLogMessage())
	assert.Nil(t, RecordsBulk.ToJSON(), "Error while catching JSON converting error!")
}

func CreateBulk() *Bulk {
	return NewBulk(
		Credentials{
			GetEnv(
				"PRIVATE_KEY",
				"7569303a-6269-4d2c-bf14-1aec9b1786a4",
			),
			"sdk-go",
			"test",
		},
	)
}
