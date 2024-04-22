// declare package
package fault

//import dependencies
import (
	"encoding/json"
	"errors"
	"fmt"
	"rosie/fault-injection/faulthelper"
	"rosie/fault-injection/utils"
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
			faulthelper.ValidateErrorCode("TFM_2009", "PANIC in InjectFault", fmt.Sprintf("%v", err), false)
		}
	}()

	//memoisation - check if paramToFunc and params are initialised
	if !initialised {
		//map functions to failureModes
		paramToFunc = mapFaults()
		faulthelper.DataDogHandle.LogInfo("paramToFunc: ", paramToFunc)
		//fmt.Print("\n paramToFunc: ", paramToFunc)

		//get parameters
		params, pErr = getParams(requestConfig, paramToFunc)
		faulthelper.DataDogHandle.LogInfo("faultConfig: ", params)
		//fmt.Print("\n params: ", params)
		initialised = true
	}

	//handle error from getParams
	if pErr != nil {
		faulthelper.DataDogHandle.LogError("TFM_2014", "Error in getParams - Fault Injection disabled", pErr.Error())
		//fmt.Print("\n Error: ", pErr.Error())
	}
	faulthelper.DataDogHandle.LogInfo("requestConfig: ", requestConfig)

	if params.FaultInjectionParams.IsEnabled && params.FaultInjectionParams.FailureMode == faultType {
		//fetch correct fault function
		faultFunction := paramToFunc.functions[faultType]

		//handle nil faultFunction
		if faultFunction == nil {
			faulthelper.DataDogHandle.LogError("TFM_20XX", "Couldn't locate fault function", "", false)
			//fmt.Print("\n Error: ", "couldn't locate fault function")
			return value
		}

		//run fault function
		modifiedValue, err := faultFunction(params, value)
		if err != nil {
			faulthelper.DataDogHandle.LogError("TFM_20XX", "Error in fault function", err.Error(), false)
			//fmt.Print("\n Error: ", err.Error())
			return value
		}

		return modifiedValue
	}

	return value
}

func getParams(requestConfig interface{}, paramToFunc FaultMap) (*utils.FaultConfig, error) {
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

	//convert requestConfig to map[string]interface{} //TODO tidy this up somehow. error handling?
	rqConfigByte, _ := json.Marshal(requestConfig)
	rqConfigMap := make(map[string]interface{})
	json.Unmarshal([]byte(rqConfigByte), &rqConfigMap)

	//convert faultInjectionParams to FIParams struct
	fipMap := utils.FIParamsMap{}
	fipByte, _ := json.Marshal(rqConfigMap["FAULT_INJECTION_PARAM"])
	json.Unmarshal([]byte(fipByte), &fipMap)

	tpErrorsByte, _ := json.Marshal(rqConfigMap["THIRD_PARTY_ERRORS_MAP"])
	tpErrorsMap := make(map[string]string)
	json.Unmarshal([]byte(tpErrorsByte), &tpErrorsMap)

	//convert WebServiceAPIErrorsMap to ErrorTypeMap
	errorMap := make(map[string]utils.ErrorTypeMap)
	errorByte, _ := json.Marshal(rqConfigMap["WS_API_ERRORS_MAP"])
	json.Unmarshal([]byte(errorByte), &errorMap)

	//map to faultConfig
	faultConfig := &utils.FaultConfig{
		FaultInjectionParams:   fipMap,
		WebserviceTimeout:      rqConfigMap["WS_API_CLIENT_TIME_OUT"].(string),
		ThirdPartyErrorsMap:    tpErrorsMap,
		WebServiceAPIErrorsMap: errorMap, //TODO check these have the same values each time.
	} //TODO check these exist

	//check fault type exists
	if faultConfig.FaultInjectionParams.IsEnabled {
		if _, ok := paramToFunc.functions[faultConfig.FaultInjectionParams.FailureMode]; !ok {
			return defaultConfig, errors.New("can't match FAILURE_MODE to Fault") //TODO make this a warning
		}
	}

	return faultConfig, nil
}
