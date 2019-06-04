package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var appVersion = "v0.3.2"

var opts struct {
	Version  bool   `long:"version" description:"Show version"`
	EnvFile  string `short:"e" long:"envfile" default:"examples/basic-test/env.json" description:"Environment file"`
	TestFile string `short:"t" long:"testfile" default:"examples/basic-test/health.json" description:"Test file"`
}

func main() {
	flags.Parse(&opts)
	if opts.Version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	loadEnvFile(opts.EnvFile)
	loadTestFile(opts.TestFile)
	configureEndpoints()

	actualResponse := sendRequest(TestData.Request)
	log.Println(actualResponse)

	if !check(TestData.Response, actualResponse) {
		os.Exit(1)
	}
}

func configureEndpoints() {
	for k, v := range EnvVars.Env {
		envVar := "${env." + k + "}"
		TestData.Request.URL = strings.Replace(TestData.Request.URL, envVar, v, 1)
	}
}

func sendRequest(request Request) Response {
	client := &http.Client{}

	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{
		resp.Header,
		resp.StatusCode,
		body,
	}
	return response
}

func check(expected Response, actual Response) bool {
	result := true
	if err := checkCode(expected, actual); err != nil {
		log.Error(err)
		result = false
	}
	if err := checkHeaders(expected, actual); err != nil {
		log.Error(err)
		result = false
	}
	if err := checkBody(expected, actual); err != nil {
		log.Error(err)
		result = false
	}
	return result
}

func checkCode(expected Response, actual Response) error {
	if expected.Code != actual.Code {
		return fmt.Errorf("Status codes are different: actual: %d, expected: %d", actual.Code, expected.Code)
	}
	return nil
}

func checkHeaders(expected Response, actual Response) error {
	if expected.Headers == nil {
		return nil // there is no any expected headers
	}

	expectedHeaders := expected.Headers.(map[string]interface{})
	if len(expectedHeaders) > 0 && actual.Headers == nil {
		return fmt.Errorf("Headers are different: actual: %v, expected: %v", actual.Headers, expectedHeaders)
	}
	actualHeaders := actual.Headers.(http.Header)

	log.Println(expectedHeaders)
	for k, v := range expectedHeaders {
		log.Println(k, v.(string))
		actualValue := actualHeaders.Get(k)
		expectedValue := v.(string)
		log.Println("values=expected:", expectedValue, ", actual: ", actualValue)
		if actualValue != expectedValue {
			return fmt.Errorf("Headers are different: actual: %v, expected: %v", actualValue, expectedValue)
		}
	}
	return nil
}

func checkBody(expected Response, actual Response) error {
	if expected.Body == nil {
		return nil // there is no any expected body
	}

	contentType := "text/plain"
	if expected.Headers != nil {
		contentType = expected.Headers.(map[string]interface{})["Content-Type"].(string)
	}
	log.Println("expected Content-Type:", contentType)

	if actual.Body == nil {
		return fmt.Errorf("Bodies are different: actual: %v, expected: %v", nil, expected.Body)
	}

	if strings.Contains(contentType, "text") {
		expectedBody := expected.Body.(string)
		actualBody := string(actual.Body.([]byte))

		log.Println(expectedBody)
		log.Println(actualBody)

		if expectedBody != actualBody {
			return fmt.Errorf("Bodies are different: actual: %v, expected: %v", actualBody, expectedBody)
		}
	} else if strings.Contains(contentType, "json") {
		switch expectedBody := expected.Body.(type) {
		case map[string]interface{}:
			var actualBody map[string]interface{}
			json.Unmarshal(actual.Body.([]byte), &actualBody)
			log.Println(expectedBody)
			equals := deepEqual(expectedBody, actualBody)
			log.Println("are equal: ", equals)
			if !equals {
				return fmt.Errorf("Bodies are different: actual: %v, expected: %v", actualBody, expectedBody)
			}
		case []interface{}:
			var actualBody []interface{}
			json.Unmarshal(actual.Body.([]byte), &actualBody)
			log.Println("array", expectedBody, "actualBody", actualBody)
			log.Println(expectedBody)
			equals := reflect.DeepEqual(expectedBody, actualBody)
			log.Println("are equal: ", equals)
			if !equals {
				return fmt.Errorf("Bodies are different: actual: %v, expected: %v", actualBody, expectedBody)
			}
		default:
			return fmt.Errorf("not supported JSON object type")
		}
	} else {
		return fmt.Errorf("Not supported Content-Type: %s", contentType)
	}

	return nil
}

func deepEqual(m1, m2 map[string]interface{}) bool {
	if reflect.DeepEqual(m1, m2) {
		return true
	}

	equals := true
	for k1, v1 := range m1 {
		if v1 != m2[k1] {
			equals = false
			log.Error("Not equals: ", k1, " values: expected: ", v1, ", actual: ", m2[k1], ", types:", reflect.TypeOf(v1), reflect.TypeOf(m2[k1]))
			if b1, a1 := isNumber(v1); b1 == true {
				if b2, a2 := isNumber(m2[k1]); b2 == true {
					log.Info("Both of them are number, so checking again...")
					if a1 == a2 {
						equals = true
					}
				}
			}
		}
	}
	return equals
}

func deepEqualArray(a1, a2 []map[string]interface{}) bool {
	return reflect.DeepEqual(a1, a2)
}

func isNumber(a interface{}) (bool, float64) {
	switch a.(type) {
	case float64:
		return true, a.(float64)
	case float32:
		return true, float64(a.(float32))
	case int:
		return true, float64(a.(int))
	case int8:
		return true, float64(a.(int8))
	case int16:
		return true, float64(a.(int16))
	case int32:
		return true, float64(a.(int32))
	case int64:
		return true, float64(a.(int64))
	case uint:
		return true, float64(a.(uint))
	case uint8:
		return true, float64(a.(uint8))
	case uint16:
		return true, float64(a.(uint16))
	case uint32:
		return true, float64(a.(uint32))
	case uint64:
		return true, float64(a.(uint64))
	default:
		log.Println("not found")
		return false, 0
	}
}
