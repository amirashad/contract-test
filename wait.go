package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func waitForSeconds(s int) {
	if s > 0 {
		log.Println("Waiting for", s, "seconds...")
		time.Sleep(time.Duration(s) * time.Second)
	}
}

func waitForEndpoint(u string) {
	if len(u) > 0 {
		count := 0
		for count < 180 && !recheck(u) {
			log.Println("Waiting for", u, "...")
			time.Sleep(time.Second)

			count++
		}
		if count == 180 { // 3 min timeout
			fmt.Println("Can't connect to server")
			os.Exit(1)
		}
	}
}

func recheck(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Debug("\t", err)
		return false
	}

	return resp.StatusCode == 200
}
