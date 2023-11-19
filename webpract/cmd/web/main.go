package main

import (
	"fmt"
	"net/http"

	"github.com/chenemiken/goland/webpract/pkg/handlers"
)

const portNumber = ":8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println((fmt.Sprintf("Starting applicatin on port %s", portNumber)))
	_ = http.ListenAndServe(portNumber, nil)
}
