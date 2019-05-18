package main

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
