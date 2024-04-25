// declare package
package fault

//import dependencies
import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rosalie12e/fault-injection/utils"
)

var (
	paramToFunc FaultMap
	params      *utils.FaultConfig
	pErr        error
	initialised bool
)

// function to act as injection point. used in main code to inject fault
func InjectFault(faultType string, value interface{}, requestConfig interface{}) interface{} {

	defer func() {
		if err := recover(); err != nil {
			//helper.ValidateErrorCode("TFM_2009", "PANIC in InjectFault", fmt.Sprintf("%v", err), false)
		}
	}()

	//memoisation - check if paramToFunc and params are initialised
	if !initialised {
		//map functions to failureModes
		paramToFunc = mapFaults()

		//get parameters
		params, pErr = getParams(requestConfig, paramToFunc)
		//helper.DataDogHandle.LogDebug("faultConfig: ", params)
		fmt.Print("\n params: ", params)
		initialised = true
	}

	//handle error from getParams
	if pErr != nil {
		//helper.DataDogHandle.LogError("TFM_2014", "Error in getParams - Fault Injection disabled", pErr.Error())
		fmt.Print("\n Error: ", pErr.Error())
	}
	//helper.DataDogHandle.LogDebug("requestConfig: ", requestConfig)

	if params.FaultInjectionParams.IsEnabled && params.FaultInjectionParams.FailureMode == faultType {
		//fetch correct fault function
		faultFunction := paramToFunc.functions[faultType]

		//handle nil faultFunction
		if faultFunction == nil {
			//helper.DataDogHandle.LogError("TFM_20XX", "Couldn't locate fault function", "", false)
			fmt.Print("\n Error: ", "couldn't locate fault function")
			return value
		}

		//run fault function
		modifiedValue, err := faultFunction(params, value)
		if err != nil {
			//helper.DataDogHandle.LogError("TFM_20XX", "Error in fault function", err.Error(), false)
			fmt.Print("\n Error: ", err.Error())
			return value
		}

		return modifiedValue
	}

	return value
}

func getParams(requestConfig interface{}, paramToFunc FaultMap) (*utils.FaultConfig, error) {
	fmt.Print("\n Internal")

	//create new instance of FaultConfig with default values
	defaultConfig := &utils.FaultConfig{
		FaultInjectionParams: utils.FIParamsMap{
			IsEnabled:   false,
			FailureMode: "",
		},
		WebserviceTimeout:      "",
		ThirdPartyErrorsMap:    make(map[string]string),
		WebServiceAPIErrorsMap: make(map[string]utils.ErrorTypeMap),
	}

	//convert requestConfig to map[string]interface{} //TODO error handling
	rqConfigByte, _ := json.Marshal(requestConfig)
	rqConfigMap := make(map[string]interface{})
	json.Unmarshal([]byte(rqConfigByte), &rqConfigMap)

	//convert faultInjectionParams to FIParams struct
	fipInt, ok := rqConfigMap["FAULT_INJECTION_PARAM"].(map[string]interface{})
	if !ok {
		return defaultConfig, errors.New("can't find FAULT_INJECTION_PARAM")
	}
	fipByte, _ := json.Marshal(fipInt)
	fipMap := utils.FIParamsMap{}
	json.Unmarshal([]byte(fipByte), &fipMap)

	//convert ThirdPartyErrorsMap to map[string]string
	tpErrorsInt, ok := rqConfigMap["THIRD_PARTY_ERRORS_MAP"].(map[string]interface{})
	if !ok {
		return defaultConfig, errors.New("can't find THIRD_PARTY_ERRORS_MAP")
	}
	tpErrorsByte, _ := json.Marshal(tpErrorsInt)
	tpErrorsMap := make(map[string]string)
	json.Unmarshal([]byte(tpErrorsByte), &tpErrorsMap)

	//convert WebServiceAPIErrorsMap to ErrorTypeMap
	apiErrorsInt, ok := rqConfigMap["WS_API_ERRORS_MAP"].(map[string]interface{})
	if !ok {
		return defaultConfig, errors.New("can't find WS_API_ERRORS_MAP")
	}
	apiErrorByte, _ := json.Marshal(apiErrorsInt)
	apiErrorMap := make(map[string]utils.ErrorTypeMap)
	json.Unmarshal([]byte(apiErrorByte), &apiErrorMap)

	//convert WebServiceTimeout to string. 2 options for JSON tag.
	timeout, ok := rqConfigMap["WS_API_CLIENT_TIME_OUT"].(string)
	if !ok {
		timeout, ok = rqConfigMap["WS_SESSION_TIMEOUT"].(string)
		if !ok {
			return defaultConfig, errors.New("can't find value for API Timeout")
		}
	}

	//map to faultConfig
	faultConfig := &utils.FaultConfig{
		FaultInjectionParams:   fipMap,
		WebserviceTimeout:      timeout,
		ThirdPartyErrorsMap:    tpErrorsMap,
		WebServiceAPIErrorsMap: apiErrorMap, //TODO check these have the same values each time.
	}

	//check fault type exists
	if faultConfig.FaultInjectionParams.IsEnabled {
		if _, ok := paramToFunc.functions[faultConfig.FaultInjectionParams.FailureMode]; !ok {
			return defaultConfig, errors.New("can't match FAILURE_MODE to Fault") //TODO make this a warning
		}
	}

	return faultConfig, nil
}
