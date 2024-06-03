package fault

import (
	"strconv"
	"strings"
	"time"

	"github.com/rosalie12e/fault-injection/helper"
	"github.com/rosalie12e/fault-injection/utils"
)

// define generic Fault strategy interface with Execute method.
type Fault interface {
	Execute(*utils.FaultConfig, interface{}) (interface{}, error)
}

// define struct for each fault type
type latency struct{}

func faultFactory(faultType string) Fault {
	// Map fault type to a concrete Fault struct
	switch faultType {
	case utils.Latency:
		return &latency{}
	// Add cases for other fault types and their corresponding structs here
	default:
		return nil // Handle unknown fault types
	}
}

// Latency injection for TFM_5001/5002
func (f *latency) Execute(faultConfig *utils.FaultConfig, value interface{}) (interface{}, error) {

	helper.DataDogHandle.LogDebug("Injecting fault: Latency")

	//parse webservice session timeout
	trimTimeout := strings.TrimSuffix(faultConfig.WebserviceTimeout, "ms")
	wsTimeout, err := strconv.Atoi(trimTimeout)
	if err != nil {
		return value, err
	}

	//set timeout to more than webservice session timeout
	timeoutValue := wsTimeout + 1

	//force timeout
	helper.DataDogHandle.LogDebug("Sleeping...")
	time.Sleep(time.Duration(timeoutValue) * time.Millisecond)

	helper.DataDogHandle.LogInfo("Fault Injected: latency of ", timeoutValue, " milliseconds")
	return nil, nil
}
