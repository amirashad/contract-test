package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	startHTTPServing()
}

func startHTTPServing() {
	http.HandleFunc("/v1/users", users)
	http.HandleFunc("/v1/users/1", user)
	http.HandleFunc("/readiness", health)
	http.HandleFunc("/health", health)

	port := "8080"
	log.Println("Starting server at port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK!")
}

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if r.Method != http.MethodGet {
		log.Println("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := User{ID: 13, Name: "Rashad"}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func users(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if r.Method != http.MethodGet {
		log.Println("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users := []User{
		User{ID: 13, Name: "Rashad"},
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}
