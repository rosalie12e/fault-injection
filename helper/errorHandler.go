package helper

import (
	"encoding/json"
	"runtime/debug"
	"strings"

	"github.com/rosalie12e/fault-injection/utils"
)

var Config *utils.FaultConfig

const (
	TECHNICAL  = "Technical"
	FUNCTIONAL = "Functional"
)

// Note  - this module only handles internal TFM errors. No 3P error handling is necessary.
// Note - for ease of use, failure of this module does not cause failure of NDC connector.
// Therefore, all errors in this connector are logged as warnings.

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

	DataDogHandle.LogInfo("errorHandler setup: ", Config)

	var tfmCustomCode string

	if strings.HasPrefix(code, "TFM_") {
		tfmCustomCode = code
	} else {
		tfmCustomCode = "TFM_2033"
	}

	criticalFunctionalAPIErrors, lessCriticalFunctionalAPIErrors, criticalTechnicalAPIErrors, lessCriticalTechnicalAPIErrors := []string{}, []string{}, []string{}, []string{}
	if Config.WebServiceAPIErrorsMap["FUNCTIONAL"].Critical != "" {
		criticalFunctionalAPIErrors = strings.Split(Config.WebServiceAPIErrorsMap["FUNCTIONAL"].Critical, ",")
	}
	if Config.WebServiceAPIErrorsMap["FUNCTIONAL"].Less_Critical != "" {
		lessCriticalFunctionalAPIErrors = strings.Split(Config.WebServiceAPIErrorsMap["FUNCTIONAL"].Less_Critical, ",")
	}
	if Config.WebServiceAPIErrorsMap["TECHNICAL"].Critical != "" {
		criticalTechnicalAPIErrors = strings.Split(Config.WebServiceAPIErrorsMap["TECHNICAL"].Critical, ",")
	}
	if Config.WebServiceAPIErrorsMap["TECHNICAL"].Less_Critical != "" {
		lessCriticalTechnicalAPIErrors = strings.Split(Config.WebServiceAPIErrorsMap["TECHNICAL"].Less_Critical, ",")
	}

	if utils.CheckSliceHasTheElement(criticalFunctionalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, FUNCTIONAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:FUNCTIONAL", "errorseverity:CRITICAL")
	} else if utils.CheckSliceHasTheElement(lessCriticalFunctionalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, FUNCTIONAL)
		DataDogHandle.LogWarnWithJSON(contextMessage, tfmCustomCode)
		DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:FUNCTIONAL", "errorseverity:LESSCRITICAL")
	} else if utils.CheckSliceHasTheElement(criticalTechnicalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
	} else if utils.CheckSliceHasTheElement(lessCriticalTechnicalAPIErrors, tfmCustomCode) {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:LESSCRITICAL")

	} else {
		contextMessage := BuildContextMessage(code, title, description, TECHNICAL)
		DataDogHandle.LogErrorWithJSON(contextMessage, tfmCustomCode)
		DataDogHandle.AddMetricWithTags(MetricError, 1, "errortype:TECHNICAL", "errorseverity:CRITICAL")
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
