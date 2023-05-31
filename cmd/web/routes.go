package main

import (
	"net/http"

	"github.com/alvinahb/clavavin/pkg/config"
	"github.com/alvinahb/clavavin/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/a-propos", handlers.Repo.About)
	mux.Get("/nouveau-vin", handlers.Repo.GetAddWine)
	mux.Post("/nouveau-vin", handlers.Repo.PostAddWineJSON)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
