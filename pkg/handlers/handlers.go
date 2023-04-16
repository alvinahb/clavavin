package handlers

import (
	"net/http"

	"github.com/alvinahb/clavavin/pkg/config"
	"github.com/alvinahb/clavavin/pkg/models"
	"github.com/alvinahb/clavavin/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AddWine(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "add_wine.page.tmpl", &models.TemplateData{})
}
