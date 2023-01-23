package coralogix

import (
	"reflect"
	"testing"
	"time"
)

func TestNewBulk(t *testing.T) {
	if reflect.TypeOf(CreateBulk()) != reflect.TypeOf(&Bulk{}) {
		t.Error("Bulk creation failed!")
	}
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
	if len(RecordsBulk.LogEntries) < 1 {
		t.Error("Adding new record to bulk failed!")
	}
}

func TestBulk_ToJSON(t *testing.T) {
	RecordsBulkJSON := CreateBulk().ToJSON()
	if RecordsBulkJSON == nil || reflect.TypeOf(RecordsBulkJSON) != reflect.TypeOf([]byte{}) {
		t.Error("Error while converting bulk to JSON!")
	}
}

func TestBulk_ToJSONFail(t *testing.T) {
	RecordsBulk := CreateBulk()
	RecordsBulk.AddRecord(*InvalidLogMessage())
	if RecordsBulk.ToJSON() != nil {
		t.Error("Error while catching JSON converting error!")
	}
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
