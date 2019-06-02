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

var appVersion = "v0.3.1"

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
	result := checkCode(TestData.Response, actualResponse)
	result = result && checkHeaders(TestData.Response, actualResponse)
	result = result && checkBody(TestData.Response, actualResponse)
	if !result {
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

func checkCode(expected Response, actual Response) bool {
	if expected.Code != actual.Code {
		log.Error("Status codes are different: actual: ", actual.Code, ", expected: ", expected.Code)
		return false
	}
	return true
}

func checkHeaders(expected Response, actual Response) bool {
	expectedHeaders := expected.Headers.(map[string]interface{})
	actualHeaders := actual.Headers.(http.Header)

	log.Println(expectedHeaders)
	result := true
	for k, v := range expectedHeaders {
		log.Println(k, v.(string))
		actualValue := actualHeaders.Get(k)
		expectedValue := v.(string)
		log.Println("values=expected:", expectedValue, ", actual: ", actualValue)
		if actualValue != expectedValue {
			log.Error("Headers are different: actual: ", actualValue, ", expected: ", expectedValue)
			result = false
		}
	}
	return result
}

func checkBody(expected Response, actual Response) bool {
	contentType := expected.Headers.(map[string]interface{})["Content-Type"].(string)
	log.Println("expected Content-Type:", contentType)
	result := true

	if strings.Contains(contentType, "text") {
		expectedBody := expected.Body.(string)
		actualBody := string(actual.Body.([]byte))

		log.Println(expectedBody)
		log.Println(actualBody)

		if expectedBody != actualBody {
			log.Error("Bodies are different: actual: ", actualBody, ", expected: ", expectedBody)
			result = false
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
				log.Error("Bodies are different: actual: ", actualBody, ", expected: ", expectedBody)
				result = false
			}
		case []interface{}:
			var actualBody []interface{}
			json.Unmarshal(actual.Body.([]byte), &actualBody)
			log.Println("array", expectedBody, "actualBody", actualBody)
			log.Println(expectedBody)
			equals := reflect.DeepEqual(expectedBody, actualBody)
			log.Println("are equal: ", equals)
			if !equals {
				log.Error("Bodies are different: actual: ", actualBody, ", expected: ", expectedBody)
				result = false
			}
		default:
			log.Error("not supported JSON object")
			result = false
		}
	} else {
		log.Error("Not supported Content-Type: ", contentType)
		result = false
	}

	return result
}

func deepEqual(m1, m2 map[string]interface{}) bool {
	if reflect.DeepEqual(m1, m2) {
		return true
	}

	equals := true
	for k1, v1 := range m1 {
		if v1 != m2[k1] {
			equals = false
			log.Error("Not equals: ", k1, " values: expected: ", v1, ", actual: ", m2[k1])
		}
	}
	return equals
}

func deepEqualArray(a1, a2 []map[string]interface{}) bool {
	return reflect.DeepEqual(a1, a2)
}
