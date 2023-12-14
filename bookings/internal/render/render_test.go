package render

import (
	"net/http"
	"testing"

	"github.com/chenemiken/goland/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "Flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value not 123")
	}
}

func TestRenderTemplate(t *testing.T) {
	t.Error("error")
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error("failed to create template cache")
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendered a non-existent package")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error()
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx, _ := session.Load(r.Context(), r.Header.Get("X-Session"))

	r = r.WithContext(ctx)
	return r, nil
}

type myWriter struct{}

func (w *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w *myWriter) Write([]byte) (int, error) {
	i := 1
	return i, nil
}

func (w *myWriter) WriteHeader(statusCode int) {}
