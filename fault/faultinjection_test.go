// declare package
package fault

import (
	"bytes"           //manipulate bytes
	"compress/gzip"   //converting gzip data
	"encoding/base64" //converting base64 data
	"io"              //basic I/O primitive interfacing
	"os"              //basic OS interfacing
	"strings"         //handling UTF-8 strings
	"testing"         //automated unit testing for Go packages
	"time"            //manipulating time

	"github.com/rosalie12e/fault-injection/utils" //utility
)

// define mock requestConfig struct
type Config struct {
	WebserviceTimeout      string                        `json:"WS_SESSION_TIMEOUT"`
	WebServiceAPIErrorsMap map[string]utils.ErrorTypeMap `json:"WS_API_ERRORS_MAP"`
	ThirdPartyErrorsMap    map[string]string             `json:"THIRD_PARTY_ERRORS_MAP,omitempty"`
	FaultInjectionParams   map[string]interface{}        `json:"FAULT_INJECTION_PARAM,omitempty"`
	IsVerbose              bool                          `json:"IS_VERBOSE"`
	ExtraKey               string                        //extra key to check mapping
}

// initialise mock requestConfig struct
func initConfig() *Config {
	config := &Config{
		FaultInjectionParams: map[string]interface{}{
			"IS_ENABLED":   true,
			"FAILURE_MODE": "latency",
		},
		WebserviceTimeout: "3000ms",
		WebServiceAPIErrorsMap: map[string]utils.ErrorTypeMap{
			"FUNCTIONAL": {
				Critical:      "",
				Less_Critical: "TFM_2033",
			},
			"TECHNICAL": {
				Critical:      "",
				Less_Critical: "",
			},
		},
		IsVerbose: false,
		ExtraKey:  "string",
	}

	//set initialised to false to ensure getParams runs
	initialised = false

	return config
}

// TestInjectLatency calls InjectFault with a valid config struct,
// with isEnabled = true and FailureMode = Latency
// checking for correct sleep time
func TestInjectLatency(t *testing.T) {
	//create requestConfig
	config := initConfig()

	//mock API POST request time start
	start := time.Now()

	//save original stdout
	catchStdout := os.Stdout
	//create pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//call Inject Fault
	InjectFault(utils.Latency, nil, &config)

	//close the write end of the pipe
	w.Close()
	//replace stdOut
	os.Stdout = catchStdout

	//print saved logs to console
	for _, log := range tfmLogToStr(r) {
		t.Log(log)
	}

	//mock API POST request time end
	duration := time.Since(start)

	//check duration and compare to expected timeout
	if duration < 3001*time.Millisecond {
		//fail test
		t.Errorf("Latency not injected: %s", duration)
	}
}

// TestInjectString calls InjectFault with an invalid config struct,
// with isEnabled = "true" and FailureMode = "latency"
// checking for log "incorrect type for IS_ENABLED"
func TestInjectString(t *testing.T) {
	config := initConfig()
	//assign empty value to FailureMode
	config.FaultInjectionParams["IS_ENABLED"] = "true"

	catchStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	InjectFault(utils.Latency, nil, &config)

	w.Close()
	os.Stdout = catchStdout

	//assign desired fault description
	expectedLog := "incorrect type for IS_ENABLED"

	//capture and decode tfmLog
	logList := tfmLogToStr(r)

	for _, log := range logList {
		t.Log(log)
	}
	//compare logs to expected output
	match := compareLog(logList, expectedLog)
	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestInjectMissingFunc calls InjectFault with faultType = "fault",
// isEnabled = true and FailureMode = "fault", with no "fault" function defined.
// checking for log "can't match FAILURE_MODE to Fault"
func TestInjectMissingFunc(t *testing.T) {
	config := initConfig()
	//assign value to FailureMode
	config.FaultInjectionParams["FAILURE_MODE"] = "fault"

	catchStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	InjectFault("fault", nil, &config)

	w.Close()
	os.Stdout = catchStdout

	expectedLog := "can't match faultType to Fault Function"
	logList := tfmLogToStr(r)
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog)

	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestInjectMissingValue calls InjectFault with an invalid config struct,
// with a blank value for FaultInjectionParams
// checking for log "can't find FAULT_INJECTION_PARAM"
func TestInjectMissingValue(t *testing.T) {
	//assign blank map to FaultInjectionParams
	config := initConfig()
	config.FaultInjectionParams = make(map[string]interface{})

	catchStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	InjectFault(utils.Latency, nil, &config)

	w.Close()
	os.Stdout = catchStdout

	expectedLog := "can't find FAULT_INJECTION_PARAM"
	logList := tfmLogToStr(r)
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog)

	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestIsVerbose calls InjectFault with isVerbose = true,
// checking for log "Fault Injection module config: "
func TestIsVerbose(t *testing.T) {
	//assign true to IsVerbose
	config := initConfig()
	config.IsVerbose = true

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(utils.Latency, nil, &config) //run injectfault

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //repalce stdOut

	expectedLog := "Fault Injection module config:" //assign desired log
	logList := tfmLogToStr(r)                       //capture and decode tfmLog
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog) //compare logs to expected output

	if match == false {
		t.Errorf("Debug statements not printed")
	}
}

// utility func to convert tfmLog to string for testing purposes
func tfmLogToStr(r *os.File) []string {
	// Read all the log content from the file
	logenziBytes, _ := io.ReadAll(r) //encoded, zipped, bytes
	logenziStrList := strings.Split(string(logenziBytes), "\n")

	logList := []string{}
	for _, log := range logenziStrList {
		logenziStr := string(log)                                      //encoded, zipped, string
		logdeziBytes, _ := base64.StdEncoding.DecodeString(logenziStr) //decoded, zipped, bytes

		//create new GZip reader
		bufz := bytes.NewBuffer(logdeziBytes)
		gzipReader, err := gzip.NewReader(bufz)
		if err != nil || err == io.EOF {
			break
		}

		var logBytes bytes.Buffer
		logBytes.ReadFrom(gzipReader) //decoded, unzipped, bytes

		logStr := logBytes.String() //decoded, unzipped, string
		logList = append(logList, logStr)
		gzipReader.Close()
	}
	// Return match
	return logList
}

// util func to compare tfmLog to string for testing purposes
func compareLog(logList []string, expectedLog string) bool {
	match := false
	for _, log := range logList {
		//compare actual log and expected error message
		if strings.Contains(log, expectedLog) {
			match = true
			break
		}
	}
	return match
}
