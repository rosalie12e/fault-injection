package helper

import (
	"encoding/json"
	"runtime/debug"
	"strings"

	"github.com/rosalie12e/fault-injection/utils"
)

var (
	ConfigParam utils.FaultConfig
)

const (
	TECHNICAL  = "Technical"
	FUNCTIONAL = "Functional"
)

// Method to log the error according to severity.
// To display the TFM custom error code for context.title and error in context.message.
func ValidateErrorCode(code, title, description string, isThirdPartyError bool) {

	defer func() {
		if err := recover(); err != nil {
			contextMessage := BuildContextMessage("TFM_2009", "PANIC in ValidateErrorCode", string(debug.Stack()), TECHNICAL)
			DataDogHandle.LogErrorWithJSON(contextMessage, "TFM_2009")
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
		}
	}()

	var tfmCustomCode string

	if strings.HasPrefix(code, "TFM_") {
		//Usually these will be internal TFM errors.
		tfmCustomCode = code
	} else {
		//Iterate through THIRD_PARTY_ERRORS_MAP keys, which are substring of description.
		for key, value := range ConfigParam.ThirdPartyErrorsMap {
			if strings.Contains(utils.FormatKey(description), key) {
				tfmCustomCode = value
				break
			}
		}
	}

	if tfmCustomCode == "" {
		//New error. Code and title needs to be registered and mapped in config server
		//Default airline error parameter are used till its registered
		if code != "" {
			tfmCustomCode = code
		} else if title != "" {
			tfmCustomCode = title
		} else if description != "" {
			tfmCustomCode = description
		} else {
			tfmCustomCode = "TFM_XXXX"
		}
		DataDogHandle.LogDebug("New API error. Key and code needs to be registered")
	}

	criticalFunctionalAPIErrors, lessCriticalFunctionalAPIErrors, criticalTechnicalAPIErrors, lessCriticalTechnicalAPIErrors := []string{}, []string{}, []string{}, []string{}
	if ConfigParam.WebServiceAPIErrorsMap["FUNCTIONAL"].Critical != "" {
		criticalFunctionalAPIErrors = strings.Split(ConfigParam.WebServiceAPIErrorsMap["FUNCTIONAL"].Critical, ",")
	}
	/*if ConfigParam.WebServiceAPIErrorsMap["FUNCTIONAL"].Less_Critical != "" {
		lessCriticalFunctionalAPIErrors = strings.Split(ConfigParam.WebServiceAPIErrorsMap["FUNCTIONAL"].Less_Critical, ",")
	}*/
	if ConfigParam.WebServiceAPIErrorsMap["TECHNICAL"].Critical != "" {
		criticalTechnicalAPIErrors = strings.Split(ConfigParam.WebServiceAPIErrorsMap["TECHNICAL"].Critical, ",")
	}
	/*if ConfigParam.WebServiceAPIErrorsMap["TECHNICAL"].Less_Critical != "" {
		lessCriticalTechnicalAPIErrors = strings.Split(ConfigParam.WebServiceAPIErrorsMap["TECHNICAL"].Less_Critical, ",")
	}*/

	if utils.CheckSliceHasTheElement(criticalFunctionalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, FUNCTIONAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		if isThirdPartyError {
			DataDogHandle.AddThirdPartyMetricWithTags(MetricThirdPartyResponseError, 1, "errortype:FUNCTIONAL", "errorseverity:CRITICAL")
		} else {
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:FUNCTIONAL", "errorseverity:CRITICAL")
		}
	} else if utils.CheckSliceHasTheElement(lessCriticalFunctionalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, FUNCTIONAL)
		DataDogHandle.LogWarnWithJSON(contextMessage, tfmCustomCode)
		if isThirdPartyError {
			DataDogHandle.AddThirdPartyMetricWithTags(MetricThirdPartyResponseError, 1, "errortype:FUNCTIONAL", "errorseverity:LESSCRITICAL")
		} else {
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:FUNCTIONAL", "errorseverity:LESSCRITICAL")
		}
	} else if utils.CheckSliceHasTheElement(criticalTechnicalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		if isThirdPartyError {
			DataDogHandle.AddThirdPartyMetricWithTags(MetricThirdPartyResponseError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
		} else {
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
		}
	} else if utils.CheckSliceHasTheElement(lessCriticalTechnicalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		if isThirdPartyError {
			DataDogHandle.AddThirdPartyMetricWithTags(MetricThirdPartyResponseError, 1, "errortype:TECHNICAL", "errorseverity:LESSCRITICAL")
		} else {
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:LESSCRITICAL")
		}
	} else {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		if isThirdPartyError {
			DataDogHandle.AddThirdPartyMetricWithTags(MetricThirdPartyResponseError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
		} else {
			DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
		}
	}
}

// Method to build context error message
func BuildContextMessage(code, title, description, errorType string) map[string]interface{} {

	var errorMap map[string]interface{}
	errorStr := `{"errorType":"` + errorType + `"}`
	_ = json.Unmarshal([]byte(errorStr), &errorMap)

	errorDescription := strings.Builder{}
	if code != "" {
		errorDescription.WriteString(utils.FormatMessage(code))
		errorDescription.WriteString(" : ")
	}
	if title != "" {
		errorDescription.WriteString(utils.FormatMessage(title))
		errorDescription.WriteString(" : ")
	}
	if description != "" {
		errorDescription.WriteString(utils.FormatMessage(description))
	}

	errorMap["errorMessage"] = errorDescription.String()

	return errorMap
}
