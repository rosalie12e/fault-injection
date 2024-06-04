// Standard DataDog package for logging and metrics
package helper

import (
	ddmodule "git-codecommit.eu-central-1.amazonaws.com/v1/repos/tfm-module-ddog-go.git"
)

var (
	DataDogHandle DataDogHelper = dataDogHelperImpl{}
)

const (
	MetricRequestCount            = "request.count"
	MetricResponseTime            = "response.time"
	MetricTimeOut                 = "emergency_exit"
	MetricThirdPartyResponseError = "response.error"
	MetricThirdPartyResponseTime  = "response.time"
	MetricResults                 = "response.results"
	MetricError                   = "error"

	ActionMessage      = "Message"
	ActionError        = "Error"
	ActionAPICall      = "API-CALL"
	ActiontypeRequest  = "Request"
	ContenttypeJSON    = "JSON"
	ContenttypePLAIN   = "PLAIN"
	ContenttypeXML     = "XML"
	ActiontypeResponse = "Response"
)

type DataDogHelper interface {
	LogWarnWithJSON(jsonObject map[string]interface{}, msg ...interface{})
	LogErrorWithJSON(jsonObject map[string]interface{}, msg ...interface{})
	LogInfo(msg ...interface{})
	LogDebug(msg ...interface{})
	AddMetricWithTags(metric string, value float64, additionalTags ...string)
}

type dataDogHelperImpl struct {
	IsVerbose bool
}

// activates debug statements
func NewDataDogHelperImpl(isVerbose bool) dataDogHelperImpl {
	return dataDogHelperImpl{IsVerbose: isVerbose}
}

// logs a warning
func (dataDogHelperImpl) LogWarnWithJSON(jsonObject map[string]interface{}, msg ...interface{}) {
	ddmodule.LogWarnWithJson(jsonObject, msg...)
}

// logs an error
func (di dataDogHelperImpl) LogErrorWithJSON(jsonObject map[string]interface{}, msg ...interface{}) {
	ddmodule.LogErrorWithJson(jsonObject, msg...)
}

// logs an info statement
func (dataDogHelperImpl) LogInfo(msg ...interface{}) {
	ddmodule.LogInfo(msg...)
}

// logs a debug statement
func (di dataDogHelperImpl) LogDebug(msg ...interface{}) {
	if di.IsVerbose {
		ddmodule.LogInfo(msg...)
	}
}

// adds a metric (e.g responsetime / error / warning etc)
func (dataDogHelperImpl) AddMetricWithTags(metric string, value float64, additionalTags ...string) {
	ddmodule.AddMetricWithTags(metric, value, additionalTags...)
}
