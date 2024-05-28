package fault

import (
	"strconv"
	"strings"
	"time"

	"github.com/rosalie12e/fault-injection/helper"
	"github.com/rosalie12e/fault-injection/utils"
)

// TODO - can this be made a method of FaultMap?
func mapFaults() utils.FaultMap { //TODO - fetch list from config store? hard to do.
	faultMap := utils.FaultMap{
		Functions: map[string]utils.FaultFunc{
			utils.Latency: latencyFunc,
		},
	}
	return faultMap
}

// Latency injection for TFM_5001/5002
func latencyFunc(faultConfig *utils.FaultConfig, value interface{}) (interface{}, error) {

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
