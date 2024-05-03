package fault

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rosalie12e/fault-injection/utils"
)

type Config struct {
	WebserviceTimeout      string                        `json:"WS_SESSION_TIMEOUT"`
	WebServiceAPIErrorsMap map[string]utils.ErrorTypeMap `json:"WS_API_ERRORS_MAP"`
	ThirdPartyErrorsMap    map[string]string             `json:"THIRD_PARTY_ERRORS_MAP,omitempty"`
	FaultInjectionParams   map[string]interface{}        `json:"FAULT_INJECTION_PARAM,omitempty"`
	IsVerbose              bool                          `json:"IS_VERBOSE"`
	ExtraKey               string
}

func initConfig() *Config {
	//initialise basic config struct
	config := &Config{
		FaultInjectionParams: map[string]interface{}{
			"IS_ENABLED":   "true",
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
		IsVerbose: true,
		ExtraKey:  "string",
	}

	//set initialised to false to ensure getParams runs
	//initialised = false

	return config
}

// TestInjectLatency calls InjectFault with a valid config struct,
// with isEnabled = true and FailureMode = Latency
// checking for correct sleep time
/*func TestInjectLatency(t *testing.T) {
	config := initConfig()

	start := time.Now()

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(Latency, nil, &config)

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //repalce stdOut

	for _, log := range tfmLogToStr(r) {
		t.Log(log)
	}

	duration := time.Since(start)

	if duration < 3001*time.Millisecond {
		t.Errorf("Latency not injected: %s", duration)
	}
}*/

// TestInjectEmpty calls InjectFault with an invalid config struct,
// with isEnabled = true and FailureMode = ""
// checking for error log
func TestInjectEmpty(t *testing.T) {
	//assign empty value to FailureMode
	config := initConfig()
	config.FaultInjectionParams["FAILURE_MODE"] = ""

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(Latency, nil, &config) //run injectfault

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //repalce stdOut

	expectedLog := "can't match FAILURE_MODE to Fault" //assign desired fault
	logList := tfmLogToStr(r)                          //capture and decode tfmLog
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog) //compare logs to expected output

	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestInjectString calls InjectFault with an invalid config struct,
// with isEnabled = "true" and FailureMode = "latency"
// checking for log "incorrect type for IS_ENABLED"
func TestInjectString(t *testing.T) {
	//assign empty value to FailureMode
	config := initConfig()
	config.FaultInjectionParams["IS_ENABLED"] = "true"

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(Latency, nil, &config) //run injectfault

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //replace stdOut

	expectedLog := "incorrect type for IS_ENABLED" //assign desired fault
	logList := tfmLogToStr(r)                      //capture and decode tfmLog
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog) //compare logs to expected output

	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestInjectInt calls InjectFault with an invalid config struct,
// with isEnabled = true and FailureMode = 1
// checking for log "can't match FAILURE_MODE to Fault"
func TestInjectInt(t *testing.T) {
	//assign int value to FailureMode
	config := initConfig()
	config.FaultInjectionParams["FAILURE_MODE"] = 1

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(Latency, nil, &config) //run injectfault

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //replace stdOut

	expectedLog := "can't match FAILURE_MODE to Fault" //assign desired fault
	logList := tfmLogToStr(r)                          //capture and decode tfmLog
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog) //compare logs to expected output

	if match == false {
		t.Errorf("Expected error message not found in logs")
	}
}

// TestInjectMissingValue calls InjectFault with an invalid config struct,
// with a blank value for FaultInjectionParams
// checking for log "can't find FAULT_INJECTION_PARAM"
/*func TestInjectMissingValue(t *testing.T) {
	//assign blank map to FaultInjectionParams
	config := initConfig()
	config.FaultInjectionParams = make(map[string]interface{})

	catchStdout := os.Stdout //save original stdout
	r, w, _ := os.Pipe()     //create pipe to capture stdout
	os.Stdout = w

	InjectFault(Latency, nil, &config) //run injectfault

	w.Close()               //close the write end of the pipe
	os.Stdout = catchStdout //replace stdOut

	expectedLog := "can't find FAULT_INJECTION_PARAM" //assign desired fault
	logList := tfmLogToStr(r)                         //capture and decode tfmLog
	for _, log := range logList {
		t.Log(log)
	}
	match := compareLog(logList, expectedLog) //compare logs to expected output

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

	InjectFault(Latency, nil, &config) //run injectfault

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
}*/

// util func to convert tfmLog to string for testing purposes
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
