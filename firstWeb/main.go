package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const portNumber = ":8080"

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.tmpl")
}

// takes a response, name of the template am a passing
func renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template: ", err)
		return
	}
}

// About page
func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.tmpl")
}

// main page
func main() {
	// url we are listening to and the annoymous function to give response to the given request
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Printf("Started application on port number: %s", portNumber)
	http.ListenAndServe(portNumber, nil)
}
