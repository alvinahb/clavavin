package main

import (
	"net/http"

	"github.com/alvinahb/clavavin/internal/config"
	"github.com/alvinahb/clavavin/internal/handlers"
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

	mux.Get("/nouveau-vin", handlers.Repo.AddWine)
	mux.Post("/nouveau-vin", handlers.Repo.PostAddWine)

	mux.Get("/carte-des-vins", handlers.Repo.WineMap)

	mux.Get("/les-vins", handlers.Repo.WinesList)
	mux.Get("/les-vins/{id}", handlers.Repo.WinePage)

	mux.Get("/mon-compte/connexion", handlers.Repo.ShowLogin)
	mux.Post("/mon-compte/connexion", handlers.Repo.PostShowLogin)
	mux.Get("/mon-compte/deconnexion", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	return mux
}
