package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

var EnvVars struct {
	Env map[string]string
}

func loadEnvFile(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &EnvVars)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("EnvVars:", EnvVars)
}
