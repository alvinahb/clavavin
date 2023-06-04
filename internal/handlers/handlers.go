package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alvinahb/clavavin/internal/config"
	"github.com/alvinahb/clavavin/internal/driver"
	"github.com/alvinahb/clavavin/internal/forms"
	"github.com/alvinahb/clavavin/internal/helpers"
	"github.com/alvinahb/clavavin/internal/models"
	"github.com/alvinahb/clavavin/internal/render"
	"github.com/alvinahb/clavavin/internal/repository"
	"github.com/alvinahb/clavavin/internal/repository/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
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

// AddWine renders the add wine page
func (m *Repository) AddWine(w http.ResponseWriter, r *http.Request) {
	var emptyWine models.Wine
	data := make(map[string]interface{})
	data["wine"] = emptyWine

	render.RenderTemplate(w, r, "add_wine.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAddWine
func (m *Repository) PostAddWine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	wine := models.Wine{
		Bottle:   r.Form.Get("bottle"),
		Domain:   r.Form.Get("domain"),
		Year:     r.Form.Get("year"),
		Location: r.Form.Get("location"),
		Color:    r.Form.Get("color"),
	}

	form := forms.New(r.PostForm)

	form.Required("bottle", "domain", "year", "location", "color")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["wine"] = wine

		m.App.Session.Put(r.Context(), "Error", "Ce vin n'a pas pu être ajouté...")
		render.RenderTemplate(w, r, "add_wine.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Nouveau vin ajouté !")
	http.Redirect(w, r, "/nouveau-vin", http.StatusSeeOther)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostAddWineJSON handles request and sends JSON response
func (m *Repository) PostAddWineJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Ajouté!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
