package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

const (
	// Host name of the HTTP Server
	Host = "localhost"
	// Port of the HTTP Server
	Port = "8080"
)

// func handleHome(w http.ResponseWriter, r *http.Request) {

	
// 	fmt.Fprintf(w, "LIVE CHAT SERVERRR!")

// }

func random_msgs(w http.ResponseWriter, r *http.Request) {


	msg := "welcome to the unknown"

	msgData := map[string]string {
		"msg": msg,
		"msg2": "I am a random message",
	}

	jsonData, err := json.Marshal(msgData)
    if err != nil {
        http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
        return
    }



	w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}



func main() {
	address := fmt.Sprintf("%s:%s", Host, Port)

	http.Handle("/", http.FileServer(http.Dir("client")))

	http.HandleFunc("/random-msgs", random_msgs)

	fmt.Println("Listening on:", fmt.Sprintf("http://%s",address))

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}
