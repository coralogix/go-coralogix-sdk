package coralogix

import (
	"encoding/json"
	"sync"
)

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

type LogBuffer struct {
	buffer []Log
	size   uint64
	lock   sync.Mutex
}

func (lb *LogBuffer) Append(log Log) {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	lb.size += log.Size()
	lb.buffer = append(lb.buffer, log)
}

func (lb *LogBuffer) Size() uint64 {
	return lb.size
}

func (lb *LogBuffer) Len() int {
	return len(lb.buffer)
}

func (lb *LogBuffer) Slice(i int) []Log {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	if i > len(lb.buffer) {
		i = len(lb.buffer)
	}
	slice := lb.buffer[:i]
	lb.buffer = lb.buffer[i:]
	for _, l := range slice {
		lb.size -= l.Size()
	}

	return slice
}
