package handlers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "home.page.html")
}

func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.html")
}
