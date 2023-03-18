package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About page")
}

func AddValuees(x, y int) int {
	return x + y
}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	n, err := fmt.Fprintf(w, "Hello, go web!")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(fmt.Sprintf("Number of bytes written: %d", n))

	// })

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println((fmt.Sprintf("Starting application on port: %s", portNumber)))
	http.ListenAndServe(portNumber, nil)
}
