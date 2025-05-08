package internal

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"api/internal/api"
)

func StartApiServer(a App, apiUrl string) {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	humaApi := humachi.New(router, huma.DefaultConfig("Lead Locate API", "1.0.0"))

	nats := a.GetNatsConnection()

	api.RegisterGowler(humaApi, nats)
	// api.RegisterMaps(humaApi, nats)

	http.ListenAndServe(apiUrl, router)
}
