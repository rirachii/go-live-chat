package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	utils "github.com/rirachii/golivechat/util"
)

const (
	
	// Host name of the HTTP Server
	Host = "localhost"
	// Port of the HTTP Server
	Port = "8080"

	// File folder for Frontend
	clientFolder = "client"
)

func main() {
	address := fmt.Sprintf("%s:%s", Host, Port)
	utils.ConsoleLog()

	var serverMuxRouter *http.ServeMux = http.NewServeMux()
	svr := &http.Server{
		Addr:        ":8080",
		Handler:     serverMuxRouter,
		ReadTimeout: 5 * time.Second,
	}

	// set up client directories
	pagesDir := http.Dir(filepath.Join(clientFolder, "pages"))
	cssDir := http.Dir(filepath.Join(clientFolder, "css"))
	javaScriptDir := http.Dir(filepath.Join(clientFolder, "js"))

	// Set up client-accessible routes
	// serverMuxRouter.Handle("GET /", http.StripPrefix("/", http.FileServer(pagesDir)))
	_ = pagesDir // ignore unused

	serverMuxRouter.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(cssDir)))
	serverMuxRouter.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(javaScriptDir)))

	// Page routes
	serverMuxRouter.HandleFunc("GET /", handleRootPageRequest)
	serverMuxRouter.HandleFunc("GET /landing", handleLanding)
	serverMuxRouter.HandleFunc("GET /register", handleRegister)

	// API routes
	serverMuxRouter.HandleFunc("GET /random-msgs", getRandomMsgs)
	serverMuxRouter.HandleFunc("POST /register", handleRegister)

	// Open server
	log.Println("Listening on:", fmt.Sprintf("http://%s", address))
	err := svr.ListenAndServe()
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}

func handleRootPageRequest(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/landing", http.StatusPermanentRedirect)

}

func handleLanding(w http.ResponseWriter, r *http.Request) {

	landingPage := "client/templates/landing.html"

	data := make(map[string]string)
	data["Title"] = "LIVE CHAT SERVERRR!"

	t, err := template.ParseFiles(landingPage)

	if err != nil {
		log.Fatal("error parsing landing template")
	}

	t.Execute(w, data)

}

func handleRegister(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL, r.Method)

	switch requestMethod := r.Method; requestMethod {
	case "GET":
		// Serve html
		registerPage := "client/templates/register.html"

		t, err := template.ParseFiles(registerPage)

		if err != nil {
			log.Fatal("error parsing register template")
		}

		t.Execute(w, nil)

	case "POST":

		// Handle request to register
		log.Println("Register POST data received!")

		type RegisterForm struct {
			Username string
			Password string
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		postData := RegisterForm{
			Username: username,
			Password: password,
		}

		log.Println("Received Username: ", postData.Username)
		log.Println("Received Password: ", postData.Password)

		// http.Redirect(w, r, "/landing", http.StatusFound)

		w.Header().Set("HX-Location", "/landing")
		w.WriteHeader(http.StatusFound)

	}

}

func getRandomMsgs(w http.ResponseWriter, r *http.Request) {

	msgs := []string{
		"random 1",
		"welcome to the unknown",
		"im a random messsage",
		"KKB on toppp",
		"akjsdhiuandi",
	}

	randomIndex := rand.Intn(len(msgs))
	randomMsg := msgs[randomIndex]

	jsonData, err := json.Marshal(randomMsg)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
