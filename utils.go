package coralogix

import "os"

// GetEnv receives environment variable or default value
func GetEnv(Variable string, Default string) string {
	if Value, Exists := os.LookupEnv(Variable); Exists {
		return Value
	}
	return Default
}
