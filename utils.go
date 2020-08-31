package coralogix

import "os"

// Get environment variable or default value
func GetEnv(Variable string, Default string) string {
    if Value, Exists := os.LookupEnv(Variable); Exists {
        return Value
	}
	return Default
}
