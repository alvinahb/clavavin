package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// GetAddWine renders the add wine page
func (m *Repository) GetAddWine(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "add_wine.page.tmpl", &models.TemplateData{})
}

// PostAddWine
func (m *Repository) PostAddWine(w http.ResponseWriter, r *http.Request) {
	bottle := r.Form.Get("bottle")
	domain := r.Form.Get("domain")
	year := r.Form.Get("year")

	w.Write([]byte(fmt.Sprintf("Bootle: %s - Domain: %s - Year: %s", bottle, domain, year)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) PostAddWineJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Ajouté!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
