package utils

// config for this module
type FaultConfig struct {
	WebserviceTimeout      string
	WebServiceAPIErrorsMap map[string]ErrorTypeMap
	ThirdPartyErrorsMap    map[string]string
	FaultInjectionParams   FIParamsMap
	IsVerbose              bool
}

type FIParamsMap struct {
	IsEnabled   bool   `json:"IS_ENABLED"`
	FailureMode string `json:"FAILURE_MODE"`
}

type ErrorTypeMap struct {
	Critical     string `json:"CRITICAL"`
	LessCritical string `json:"LESS_CRITICAL"`
}
