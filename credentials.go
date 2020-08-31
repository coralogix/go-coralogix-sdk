package coralogix

// Credentials is the container for Coralogix data
type Credentials struct {
    PrivateKey      string // Coralogix private key
    ApplicationName string // Your application name
    SubsystemName   string // Subsystem name of your application
}
