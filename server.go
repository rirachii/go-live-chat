package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// Host name of the HTTP Server
	Host = "localhost"
	// Port of the HTTP Server
	Port = "8080"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "LIVE CHAT SERVERRR!")
}

func main() {
	http.HandleFunc("/", home)
	err := http.ListenAndServe(Host+":"+Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}