package handlers

import (
	"net/http"

	"github.com/chenemiken/goland/webpract/models"
	"github.com/chenemiken/goland/webpract/pkg/config"
	"github.com/chenemiken/goland/webpract/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remoteIp", remoteIp)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	var stringMap = make(map[string]string)
	stringMap["test"] = "Hello dean"

	stringMap["remoteIp"] = m.App.Session.GetString(r.Context(), "remoteIp")

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
