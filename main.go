package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	Version  bool   `short:"v" long:"version" description:"Show version"`
	EnvFile  string `short:"e" long:"env" default:"examples/basic-test/env.json" description:"Environment file"`
	TestFile string `short:"t" long:"testfile" default:"examples/basic-test/health.json" description:"Test file"`
}

func main() {
	flags.Parse(&opts)
	if opts.Version {
		fmt.Println("v0.0.2")
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

	var response Response
	response.Code = resp.StatusCode
	response.Headers = resp.Header

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	response.Body = string(body)

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
	expectedBody := expected.Body.(string)
	actualBody := actual.Body.(string)

	log.Println(expectedBody)
	log.Println(actualBody)

	if expectedBody != actualBody {
		log.Error("Bodies are different: actual: ", actualBody, ", expected: ", expectedBody)
	}
}
