package fault

import (
	"strconv"
	"strings"
	"time"

	"github.com/rosalie12e/fault-injection/helper"
	"github.com/rosalie12e/fault-injection/utils"
)

// DEFINE LIST OF FAULT TYPES HERE. This is what can be declared in the config store
const (
	Latency = "latency"
	Fault1  = "fault1"
)

type FaultFunc func(*utils.FaultConfig, interface{}) (interface{}, error)

type FaultMap struct {
	functions map[string]FaultFunc
}

func mapFaults() FaultMap { //TODO - fetch list from config store? hard to do.
	paramToFunc := FaultMap{
		functions: map[string]FaultFunc{
			Latency: latencyFunc,
		},
	}
	return paramToFunc
}

// Latency injection for TFM_5001/5002
func latencyFunc(requestConfig *utils.FaultConfig, value interface{}) (interface{}, error) {

	helper.DataDogHandle.LogDebug("Injecting fault: Latency")

	//parse webservice session timeout
	trimTimeout := strings.TrimSuffix(requestConfig.WebserviceTimeout, "ms")
	wsTimeout, err := strconv.Atoi(trimTimeout)
	if err != nil {
		return nil, err
	}

	//set timeout to more than webservice session timeout
	timeoutValue := wsTimeout + 1

	//force timeout
	helper.DataDogHandle.LogDebug("Sleeping...")
	time.Sleep(time.Duration(timeoutValue) * time.Millisecond)

	helper.DataDogHandle.LogInfo("Fault Injected: latency of ", timeoutValue, " milliseconds")
	return nil, nil
}
