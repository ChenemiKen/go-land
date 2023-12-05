package main

import (
	"testing"

	"github.com/chenemiken/goland/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	h := routes(&app)

	switch v := h.(type) {
	case *chi.Mux:

	default:
		t.Errorf("type is not chi.Mux, but is %T", v)
	}
}
