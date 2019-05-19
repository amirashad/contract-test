package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	Method  string `json:"method"`
	Headers struct {
		ContentType string `json:"Content-Type"`
	} `json:"headers"`
	URL  string      `json:"url"`
	Data interface{} `json:"data"`
}

type Response struct {
	Headers interface{} `json:"headers"`
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
}

type Test struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

var TestData Test

func loadTestFile(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &TestData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TestData:", TestData)
}
