package coralogix

import (
	"reflect"
	"testing"
	"time"
)

func TestSendRequestSuccess(t *testing.T) {
	BulkToSend := CreateBulk()
	BulkToSend.AddRecord(Log{
		float64(time.Now().Unix()) * 1000.0,
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
		0,
	})
	HTTPStatus := SendRequest(BulkToSend)
	if HTTPStatus != 200 {
		t.Error("Logs bulk sending failed!")
	}
}

func TestSendRequestPostFail(t *testing.T) {
	SetDebug(true)
	BulkToSend := CreateBulk()
	BulkToSend.AddRecord(*InvalidLogMessage())
	HTTPStatus := SendRequest(BulkToSend)
	if HTTPStatus > 0 {
		t.Error("Sending of invalid request should be failed!")
	}
}

func TestSendRequestErrorResponseStatus(t *testing.T) {
	BulkToSend := CreateBulk()
	BulkToSend.AddRecord(Log{
		1,
		Level.DEBUG,
		"Test message",
		LogCategory,
		"",
		"",
		"",
		0,
	})
	HTTPStatus := SendRequest(BulkToSend)
	if HTTPStatus == 200 {
		t.Error("Logs bulk was successful!")
	}
}

func TestGetTimeSync(t *testing.T) {
	Status, TimeDelta := GetTimeSync()
	if Status == false || reflect.TypeOf(TimeDelta).Kind() != reflect.Float64 {
		t.Error("Time synchronization failed!")
	}
}
