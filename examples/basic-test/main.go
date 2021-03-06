package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	startHTTPServing()
}

func startHTTPServing() {
	http.HandleFunc("/v1/users", users)
	http.HandleFunc("/v1/states", states)
	http.HandleFunc("/v1/users/1", user)
	http.HandleFunc("/health", health)
	http.HandleFunc("/shutdown", shutdown)

	port := "8080"
	log.Println("Starting server at port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK!")
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "exiting...")
	os.Exit(1)
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Weight  float32 `json:"weight"`
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if r.Method == http.MethodGet {
		user := User{ID: 13, Name: "Rashad", Surname: "Amirjanov", Weight: 81.25}

		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		log.Println("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if r.Method == http.MethodGet {
		users := []User{
			User{ID: 13, Name: "Rashad", Surname: "Amirjanov", Weight: 81.25},
			User{ID: 14, Name: "Pasha", Surname: "Amirjanov", Weight: 18.5},
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	} else if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var user User
		if json.Unmarshal(data, &user) != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.ID = 1
		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		log.Println("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func states(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if r.Method == http.MethodGet {
		states := []string{
			"OPEN", "CLOSED", "PENDING",
		}

		jsonData, err := json.Marshal(states)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	} else if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		createdState := string(data)

		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(createdState))
	} else {
		log.Println("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
