// declare package
package fault

//import dependencies
import (
	"encoding/json"
	"errors"
	"runtime/debug"

	"github.com/rosalie12e/fault-injection/helper"
	"github.com/rosalie12e/fault-injection/utils"
)

var (
	params      *utils.FaultConfig
	initialised bool
	pErr        error
)

// function to act as injection point. used in main code to inject fault
func InjectFault(faultType string, value interface{}, requestConfig interface{}) interface{} {

	helper.DataDogHandle.LogInfo("Running InjectFault")

	//handle panic
	defer func() {
		if err := recover(); err != nil {
			helper.ValidateErrorCode("TFM_2009", "PANIC in InjectFault", string(debug.Stack()), false)
		}
	}()

	//memoisation - check if faultMap and params are initialised
	if !initialised {
		//get parameters
		params, pErr = getParams(requestConfig)
		helper.Config = params
		initialised = true
	}
	if pErr != nil {
		helper.ValidateErrorCode("TFM_2033", "Error in getParams - Fault Injection disabled", pErr.Error(), false)
	}

	//activate debug statements
	helper.DataDogHandle = helper.NewDataDogHelperImpl(params.IsVerbose)
	helper.DataDogHandle.LogDebug("Fault Injection module config: ", params)

	if params.FaultInjectionParams.IsEnabled && params.FaultInjectionParams.FailureMode == faultType {
		//fetch fault function
		helper.DataDogHandle.LogDebug("Fetching fault type: ", faultType)
		faultFunc := faultFactory(faultType)
		//check faultType matches a function.
		if faultFunc == nil {
			helper.ValidateErrorCode("TFM_2033", "can't match faultType to Fault Function", "", false)
			return value
		}

		//run fault function
		modifiedValue, err := faultFunc.Execute(params, value)
		if err != nil {
			helper.ValidateErrorCode("TFM_2033", "Error in fault function", err.Error(), false)
			return value
		}

		return modifiedValue
	}

	return value
}

func getParams(requestConfig interface{}) (*utils.FaultConfig, error) {
	helper.DataDogHandle.LogInfo("Fetching Fault Injection Parameters")

	//create new instance of FaultConfig with default values
	faultConfig := &utils.FaultConfig{
		FaultInjectionParams: utils.FIParamsMap{
			IsEnabled:   false,
			FailureMode: "",
		},
		WebserviceTimeout: "",
		WebServiceAPIErrorsMap: map[string]utils.ErrorTypeMap{
			"FUNCTIONAL": {
				Less_Critical: "TFM_2033", //add this here to make sure it always logs as a warning
			},
		},
		IsVerbose: true, //force this to be true
	}

	//convert requestConfig to map[string]interface{}
	rqConfigByte, _ := json.Marshal(requestConfig)
	rqConfigMap := make(map[string]interface{})
	json.Unmarshal([]byte(rqConfigByte), &rqConfigMap)

	//convert faultInjectionParams to FIParams struct and check results
	fipInt, exists := rqConfigMap["FAULT_INJECTION_PARAM"].(map[string]interface{})
	if !exists {
		return faultConfig, errors.New("can't find FAULT_INJECTION_PARAM")
	}
	//check that IS_ENABLED is bool
	_, ok := fipInt["IS_ENABLED"].(bool)
	if !ok {
		return faultConfig, errors.New("incorrect type for IS_ENABLED")
	}
	fipByte, _ := json.Marshal(fipInt)
	fipMap := utils.FIParamsMap{}
	json.Unmarshal([]byte(fipByte), &fipMap)

	//convert WebServiceAPIErrorsMap to ErrorTypeMap
	apiErrorsInt, exists := rqConfigMap["WS_API_ERRORS_MAP"].(map[string]interface{})
	if !exists {
		return faultConfig, errors.New("can't find WS_API_ERRORS_MAP")
	}
	apiErrorByte, _ := json.Marshal(apiErrorsInt)
	apiErrorMap := make(map[string]utils.ErrorTypeMap)
	json.Unmarshal([]byte(apiErrorByte), &apiErrorMap)

	//convert WebServiceTimeout to string. 2 options for JSON tag.
	timeout, exists := rqConfigMap["WS_API_CLIENT_TIME_OUT"].(string)
	if !exists {
		timeout, exists = rqConfigMap["WS_SESSION_TIMEOUT"].(string)
		if !exists {
			return faultConfig, errors.New("can't find value for API Timeout")
		}
	} //TODO for Pricing/Book - assert generic timeout as it's not always specified in the config

	//convert IsVerbose to bool
	isVerbose, exists := rqConfigMap["IS_VERBOSE"].(bool)
	if !exists {
		helper.DataDogHandle.LogInfo("can't find value for IS_VERBOSE - will not print debug statements.")
	}

	//map to faultConfig
	faultConfig = &utils.FaultConfig{
		FaultInjectionParams:   fipMap,
		WebserviceTimeout:      timeout,
		WebServiceAPIErrorsMap: apiErrorMap,
		IsVerbose:              isVerbose,
	}

	return faultConfig, nil
}
