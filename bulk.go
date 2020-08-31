package coralogix

import (
    "encoding/json"
    "os"
)

// Bulk describe logs batch format for Coralogix API
type Bulk struct {
    PrivateKey      string `json:"privateKey"`      // Coralogix private key
    ApplicationName string `json:"applicationName"` // Your application name
    SubsystemName   string `json:"subsystemName"`   // Subsystem name of your application
    ComputerName    string `json:"computerName"`    // Current machine hostname
    LogEntries      []Log  `json:"logEntries"`      // Log records list
}

// NewBulk initialize new logs bulk
func NewBulk(Credentials Credentials) *Bulk {
    Hostname, _ := os.Hostname()
    return &Bulk{
        Credentials.PrivateKey,
        Credentials.ApplicationName,
        Credentials.SubsystemName,
        Hostname,
        []Log{},
    }
}

// AddRecord add log record to the logs bulk
func (bulk *Bulk) AddRecord(Record Log) {
    bulk.LogEntries = append(bulk.LogEntries, Record)
}

// ToJSON convert logs bulk to JSON format
func (bulk *Bulk) ToJSON() []byte {
    data, err := json.Marshal(bulk)
    if err != nil {
        DebugLogger.Println("Can't convert logs bulk to JSON: ", err)
        return nil
    }
    return data
}
