package coralogix

import "encoding/json"

// Log describe record format for Coralogix API
type Log struct {
	Timestamp  float64 `json:"timestamp"`  // Log record timestamp
	Severity   uint    `json:"severity"`   // Log record severity level
	Text       string  `json:"text"`       // Log record message
	Category   string  `json:"category"`   // Log record category
	ClassName  string  `json:"className"`  // Log record class name
	MethodName string  `json:"methodName"` // Log record method name
	ThreadID   string  `json:"threadId"`   // Thread ID

	size uint64
}

// Size calculate log record length in bytes
func (Record *Log) Size() uint64 {
	if Record.size == 0 {
		JSONRecord, err := json.Marshal(Record)
		if err != nil {
			return 0
		}
		Record.size = uint64(len(string(JSONRecord)))
	}

	return Record.size
}
