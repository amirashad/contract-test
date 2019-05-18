package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"default" description:"Run profile"`
	File    string `short:"f" long:"file" default:"examples/basic-test/health.json" description:"Test file"`
}

func main() {
	flags.Parse(&opts)

	log.Println(opts.Profile)

	jsonFile, err := ioutil.ReadFile(opts.File)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(jsonFile))

	var test Test
	err = json.Unmarshal(jsonFile, &test)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(test)

	// var parsed map[string]interface{}
	// err = json.Unmarshal(jsonFile, &parsed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(parsed)

	actualResponse := sendRequest(test.Request)
	log.Println(actualResponse)
	checkCode(test.Response, actualResponse)
	checkHeaders(test.Response, actualResponse)
	checkBody(test.Response, actualResponse)
}

func sendRequest(request Request) Response {
	client := &http.Client{
		// CheckRedirect: redirectPolicyFunc,
	}

	url := strings.Replace(request.URL, "${env.USERS_ENDPOINT}", "http://localhost:80", 1)

	req, err := http.NewRequest(request.Method, url, nil)
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

	// var responseObj map[string]interface{}

	// http.Get(req.Method, req.URL)
	// http.Post(req.Method, req.URL)
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
