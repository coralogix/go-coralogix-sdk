package go_coralogix_sdk

// SeveritiesList describe logs levels values
type SeveritiesList struct {
	TRACE    uint
	DEBUG    uint
	VERBOSE  uint
	INFO     uint
	WARNING  uint
	ERROR    uint
	CRITICAL uint
	FATAL    uint
	PANIC    uint
}

// Level is a list with default logs levels values
var Level = SeveritiesList{
	1,
	1,
	2,
	3,
	4,
	5,
	6,
	6,
	6,
}
