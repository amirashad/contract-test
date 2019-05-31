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

var appVersion = "v0.1.0"

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
	checkCode(TestData.Response, actualResponse)
	checkHeaders(TestData.Response, actualResponse)
	checkBody(TestData.Response, actualResponse)
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

func checkCode(expected Response, actual Response) {
	if expected.Code != actual.Code {
		log.Error("Status codes are different: actual: ", actual.Code, ", expected: ", expected.Code)
	}
}

func checkHeaders(expected Response, actual Response) {
	expectedHeaders := expected.Headers.(map[string]interface{})
	actualHeaders := actual.Headers.(http.Header)

	log.Println(expectedHeaders)
	for k, v := range expectedHeaders {
		log.Println(k, v.(string))
		actualValue := actualHeaders.Get(k)
		expectedValue := v.(string)
		log.Println("values=expected:", expectedValue, ", actual: ", actualValue)
		if actualValue != expectedValue {
			log.Error("Headers are different: actual: ", actualValue, ", expected: ", expectedValue)
		}
	}
}

func checkBody(expected Response, actual Response) {
	contentType := expected.Headers.(map[string]interface{})["Content-Type"].(string)
	log.Println("expected Content-Type:", contentType)

	if strings.Contains(contentType, "text") {
		expectedBody := expected.Body.(string)
		actualBody := string(actual.Body.([]byte))

		log.Println(expectedBody)
		log.Println(actualBody)

		if expectedBody != actualBody {
			log.Error("Bodies are different: actual: ", actualBody, ", expected: ", expectedBody)
		}
	} else if strings.Contains(contentType, "json") {
		switch expectedBody := expected.Body.(type) {
		case map[string]interface{}:
			var actualBody map[string]interface{}
			json.Unmarshal(actual.Body.([]byte), &actualBody)
			log.Println(expectedBody)
			log.Println("are equal: ", deepEqual(expectedBody, actualBody))
		}
	} else {
		log.Error("Not supported Content-Type: ", contentType)
	}
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
