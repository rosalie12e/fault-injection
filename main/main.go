package main

import (
	"github.com/rosalie12e/fault-injection/fault"
)

type Config struct {
	WebserviceTimeout      string                  `json:"WS_API_CLIENT_TIME_OUT"`
	WebServiceAPIErrorsMap map[string]ErrorTypeMap `json:"WS_API_ERRORS_MAP,omitempty"`
	ThirdPartyErrorsMap    map[string]string
	FaultInjectionParams   map[string]interface{} `json:"FAULT_INJECTION_PARAM,omitempty"`
	ExtraNonsense          string
}

type ErrorTypeMap struct {
	Critical string `json:"CRITICAL"`
}

func initConfig() Config {

	config := Config{
		FaultInjectionParams: map[string]interface{}{
			"IS_ENABLED":   true,
			"FAILURE_MODE": "latency",
		},
		WebserviceTimeout:   "3000",
		ThirdPartyErrorsMap: make(map[string]string),
		ExtraNonsense:       "nonsense",
	} //default

	return config
}

func main() {
	config := initConfig()

	_ = fault.InjectFault(fault.Latency, nil, config)
}
