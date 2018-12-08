package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
)

type Contact struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone"`
	ContactReason string `json:"contact_reason"`
	Message       string `json:"message"`
}

func contactEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var contact Contact
	_ = json.NewDecoder(r.Body).Decode(&contact)

	user := os.Getenv("USER_EMAIL")
	password := os.Getenv("PASSWORD")

	s := fmt.Sprintf("Name: %s \nPhone: %s\n Purpose of Contact: %s \n\n%s",
		contact.Name, contact.PhoneNumber, contact.ContactReason, contact.Message)

	m := gomail.NewMessage()
	m.SetHeader("From", contact.Email)
	m.SetHeader("To", user)
	m.SetHeader("Subject", contact.ContactReason+" - "+contact.Name)
	m.SetBody("text/plain", s)

	fmt.Println(user)

	d := gomail.NewPlainDialer("smtp.gmail.com", 465, user, password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return http.HandlerFunc(fn)
}

func home(w http.ResponseWriter, r *http.Request) {
	response := time.Now().String()
	fmt.Fprintln(w, response)
}

func main() {

	port := flag.String("port", "8080", "the port application is listening on")
	// staticAssets := flag.String("staticAssets", "client/dist/", "Serves static assets")
	// app := flag.String("app", "client/dist/index.html", "Serve JavaScript application's entry-point (index.html)")

	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/contact", contactEmailHandler).Methods("POST")
	api.HandleFunc("/home", home).Methods("GET")
	// Static assets directly.
	// r.PathPrefix("/static").Handler(http.FileServer(http.Dir(*staticAssets)))
	// JavaScript application  entry-point
	// r.PathPrefix("/").HandlerFunc(IndexHandler(*app))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + *port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Server starting on")
	log.Fatal(srv.ListenAndServe())
}
