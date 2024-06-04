// declare package
package fault

import (
	"strconv" //converting data to strings
	"strings" //handle UTF-8 strings
	"time"    //manipulate time

	"github.com/rosalie12e/fault-injection/helper" //logging
	"github.com/rosalie12e/fault-injection/utils"  //utility
)

// define generic Fault strategy interface with Execute method.
type Fault interface {
	Execute(*utils.FaultConfig, interface{}) (interface{}, error)
}

// define concrete implementation for each fault type
type latency struct{}
type demoFault struct{} //demo fault

func faultFactory(faultType string) Fault {
	// Map fault type to a concrete Fault struct
	switch faultType {
	case utils.Latency: //TFM_5001/5002 - Context Deadline Exceeded
		return &latency{}
	case utils.Fault1: //demo fault
		return &demoFault{}
	// Add cases for other fault types and their corresponding structs here
	default:
		return nil // Handle unknown fault types
	}
}

// Latency injection for TFM_5001/5002: Context Deadline Exceeded
func (f *latency) Execute(faultConfig *utils.FaultConfig, value interface{}) (interface{}, error) {

	helper.DataDogHandle.LogDebug("Injecting fault: Latency") //debug statement

	//parse webservice session timeout and standardise format
	trimTimeout := strings.TrimSuffix(faultConfig.WebserviceTimeout, "ms")
	wsTimeout, err := strconv.Atoi(trimTimeout)
	//handle error in strconv.Atoi
	if err != nil {
		return value, err
	}

	//set timeout to more than webservice session timeout
	timeoutValue := wsTimeout + 1

	//force timeout
	helper.DataDogHandle.LogDebug("Sleeping...") //debug statement
	time.Sleep(time.Duration(timeoutValue) * time.Millisecond)

	helper.DataDogHandle.LogInfo("Fault Injected: latency of ", timeoutValue, " milliseconds") //log statement
	return nil, nil
}

// generic fault for demo purposes
func (f *demoFault) Execute(faultConfig *utils.FaultConfig, value interface{}) (interface{}, error) {

	//debug statement
	helper.DataDogHandle.LogDebug("Injecting fault: demo")

	//assignvalue to modifiedValue
	modifiedValue := value

	//log statement
	helper.DataDogHandle.LogInfo("Fault Injected: demo fault.")

	return modifiedValue, nil
}
