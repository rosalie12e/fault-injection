package fault

import (
	"rosie/fault-injection/faulthelper"
	"rosie/fault-injection/utils"
	"strconv"
	"strings"
	"time"
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

	faulthelper.DataDogHandle.LogInfo("running LatencyFunc")
	//fmt.Print("\n running LatencyFunc")

	//parse webservice session timeout
	trimTimeout := strings.TrimSuffix(requestConfig.WebserviceTimeout, "ms")
	wsTimeout, err := strconv.Atoi(trimTimeout)
	if err != nil {
		return nil, err
	}

	//set timeout to more than webservice session timeout
	timeoutValue := wsTimeout + 1

	//force timeout
	faulthelper.DataDogHandle.LogInfo("sleeping...")
	//fmt.Print("\n sleeping...")
	time.Sleep(time.Duration(timeoutValue) * time.Millisecond)

	faulthelper.DataDogHandle.LogInfo("Fault Injected: latency of ", timeoutValue, " milliseconds")
	//fmt.Print("\n Fault Injected: latency of ", timeoutValue, " milliseconds")
	return nil, nil
}
