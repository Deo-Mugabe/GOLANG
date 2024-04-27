package main

import (
	"fmt"
	"github.com/Deo-Mugabe/GOLANG/pkg/config"
	"github.com/Deo-Mugabe/GOLANG/pkg/handlers"
	"github.com/Deo-Mugabe/GOLANG/pkg/render"
	"log"
	"net/http"
)

const portNumber = ":8080"

// main page
func main() {

	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	// url we are listening to and the annoymous function to give response to the given request
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("Started application on port number: %s", portNumber)
	http.ListenAndServe(portNumber, nil)
}
