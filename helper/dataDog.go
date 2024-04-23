package helper

import (
	"reflect"
	"regexp"
	"runtime/debug"

	ddmodule "git-codecommit.eu-central-1.amazonaws.com/v1/repos/tfm-module-ddog-go.git"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/tfm-module-ddog-go.git/util"
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
	WrapHandler(handler interface{}, aPIKey, site string) interface{}
	InitLayer(environment, product, channel, serviceName, metricServiceName, sessionID, requestID string) error
	AddTracking(action, actionType, body, contentType, title string)
	LogInfoWithNamedJSON(name string, jsonObject map[string]interface{}, msg ...interface{})
	LogWarnWithJSON(jsonObject map[string]interface{}, msg ...interface{})
	LogError(msg ...interface{})
	LogErrorWithJSON(jsonObject map[string]interface{}, msg ...interface{})
	LogErrorWithInterface(errObject interface{}, msg ...interface{})
	LogInfo(msg ...interface{})
	LogWarn(msg ...interface{})
	LogDebug(msg ...interface{})
	LogDebugWithJSON(jsonObject map[string]interface{}, msg ...interface{})
	AddMetric(metric string, value float64)
	AddMetricWithTags(metric string, value float64, additionalTags ...string)
	AddThirdPartyMetric(metric string, value float64)
	AddThirdPartyMetricWithTags(metric string, value float64, additionalTags ...string)
}

type dataDogHelperImpl struct {
	IsVerbose bool
}

func NewDataDogHelperImpl(isVerbose bool) dataDogHelperImpl {
	return dataDogHelperImpl{IsVerbose: isVerbose}
}

func (dataDogHelperImpl) WrapHandler(handler interface{}, aPIKey, site string) interface{} {
	DDogConfig := ddmodule.DDogConfig{
		APIKey: aPIKey,
		Site:   site,
	}
	return ddmodule.WrapHandler(handler, DDogConfig)
}
func (dataDogHelperImpl) InitLayer(environment, product, channel, serviceName, metricServiceName,
	sessionID, requestID string) error {
	ddmodule.InitLayer(ddmodule.LayerConfig{
		Environment:       environment,
		Product:           product,
		Channel:           channel,
		ServiceName:       serviceName,
		MetricServiceName: metricServiceName,
		SessionId:         sessionID,
		RequestId:         requestID,
	})
	return nil
}
func (dataDogHelperImpl) AddTracking(action, actionType, body, contentType, title string) {
	ddmodule.AddTracking(action, actionType, body, contentType, title)
}
func (dataDogHelperImpl) LogInfoWithNamedJSON(name string, jsonObject map[string]interface{}, msg ...interface{}) {
	ddmodule.LogInfoWithNamedJson(name, jsonObject, msg...)
}
func (dataDogHelperImpl) LogWarnWithJSON(jsonObject map[string]interface{}, msg ...interface{}) {
	ddmodule.LogWarnWithJson(jsonObject, msg...)
}
func (dataDogHelperImpl) LogError(msg ...interface{}) {
	ddmodule.LogError(msg...)
}

func (di dataDogHelperImpl) LogErrorWithJSON(jsonObject map[string]interface{}, msg ...interface{}) {
	ddmodule.LogErrorWithJson(jsonObject, msg...)
}

func (dataDogHelperImpl) LogInfo(msg ...interface{}) {
	ddmodule.LogInfo(msg...)
}

func (dataDogHelperImpl) LogWarn(msg ...interface{}) {
	ddmodule.LogWarn(msg...)
}

func (di dataDogHelperImpl) LogDebug(msg ...interface{}) {
	if di.IsVerbose {
		ddmodule.LogInfo(msg...)
	}
}

func (di dataDogHelperImpl) LogDebugWithJSON(jsonObject map[string]interface{}, msg ...interface{}) {
	if di.IsVerbose {
		ddmodule.LogInfoWithJson(jsonObject, msg...)
	}
}

func (dataDogHelperImpl) AddMetric(metric string, value float64) {
	ddmodule.AddMetric(metric, value)
}
func (dataDogHelperImpl) AddMetricWithTags(metric string, value float64, additionalTags ...string) {
	ddmodule.AddMetricWithTags(metric, value, additionalTags...)
}
func (dataDogHelperImpl) AddThirdPartyMetric(metric string, value float64) {
	ddmodule.AddThirdPartyMetric(metric, value)
}
func (dataDogHelperImpl) AddThirdPartyMetricWithTags(metric string, value float64, additionalTags ...string) {
	ddmodule.AddThirdPartyMetricWithTags(metric, value, additionalTags...)
}
func (dataDogHelperImpl) LogErrorWithInterface(errObject interface{}, msg ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			ddmodule.LogError("something went wrong in LogErrorWithInterface", err, string(debug.Stack()))
		}
	}()
	errorObjectMap := make(map[string]interface{})
	errorObjectType := reflect.TypeOf(errObject)
	if errorObjectType == nil {
		return
	}
	var errorText string
	if errorMeothod, found := errorObjectType.MethodByName("Error"); found {
		result := errorMeothod.Func.Call([]reflect.Value{reflect.ValueOf(errObject)})
		errorText = stringiFy(result[0].String())
		errorObjectMap["errorMessage"] = errorText
	} else {
		errorObjectMap = util.JMap(errObject)
	}
	ddmodule.LogErrorWithJson(errorObjectMap, msg...)
}

func stringiFy(text string) (errR string) {
	defer func() {
		if err := recover(); err != nil {
			errR = "some thing went wrong while stringiFying error "
		}
	}()
	reg, _ := regexp.Compile("[^a-zA-Z0-9\\s\\t]+")
	processedString := reg.ReplaceAllString(text, "-")
	return processedString
}
