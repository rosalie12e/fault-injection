package utils

// DEFINE LIST OF FAULT TYPES HERE. This is what can be declared in the config store
const (
	Latency = "latency"
	Fault1  = "fault1"
)

type FaultFunc func(*FaultConfig, interface{}) (interface{}, error)

type FaultMap struct {
	Functions map[string]FaultFunc
}

// config for this module
type FaultConfig struct {
	WebserviceTimeout      string
	WebServiceAPIErrorsMap map[string]ErrorTypeMap
	FaultInjectionParams   FIParamsMap
	IsVerbose              bool
}

type FIParamsMap struct {
	IsEnabled   bool   `json:"IS_ENABLED"`
	FailureMode string `json:"FAILURE_MODE"`
}

type ErrorTypeMap struct {
	Critical      string `json:"CRITICAL"`
	Less_Critical string `json:"LESS_CRITICAL"`
}
