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
    ThreadId   string  `json:"threadId"`   // Thread ID
}

// Size calculate log record length in bytes
func (Record *Log) Size() int {
    JSONRecord, err := json.Marshal(Record)
    if err != nil {
        return -1
    }
    return len(string(JSONRecord))
}
