package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chenemiken/goland/webpract/pkg/config"
	"github.com/chenemiken/goland/webpract/pkg/handlers"
	"github.com/chenemiken/goland/webpract/pkg/render"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println((fmt.Sprintf("Starting application on port %s", portNumber)))
	// _ = http.ListenAndServe(portNumber, nil)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
