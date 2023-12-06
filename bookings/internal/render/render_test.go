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

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx, _ := session.Load(r.Context(), r.Header.Get("X-Session"))

	r = r.WithContext(ctx)
	return r, nil
}
