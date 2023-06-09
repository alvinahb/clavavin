package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alvinahb/clavavin/internal/config"
	"github.com/alvinahb/clavavin/internal/driver"
	"github.com/alvinahb/clavavin/internal/forms"
	"github.com/alvinahb/clavavin/internal/helpers"
	"github.com/alvinahb/clavavin/internal/models"
	"github.com/alvinahb/clavavin/internal/render"
	"github.com/alvinahb/clavavin/internal/repository"
	"github.com/alvinahb/clavavin/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
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
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// AddWine renders the add wine page
func (m *Repository) AddWine(w http.ResponseWriter, r *http.Request) {
	var emptyWine models.Wine
	data := make(map[string]interface{})
	data["wine"] = emptyWine

	render.Template(w, r, "add_wine.page.tmpl", &models.TemplateData{
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

	var appellation = ""
	if r.Form.Get("appellationType") != "" || r.Form.Get("appellationName") != "" {
		appellation = fmt.Sprintf("%s - %s",
			r.Form.Get("appellationType"), r.Form.Get("appellationName"))
	}

	wine := models.Wine{
		Name:        r.Form.Get("name"),
		Domain:      r.Form.Get("domain"),
		Year:        r.Form.Get("year"),
		Appellation: appellation,
		Location:    r.Form.Get("location"),
		Color:       r.Form.Get("color"),
		Culture:     r.Form.Get("culture"),
		Varieties:   r.Form.Get("varieties"),
		Robe:        r.Form.Get("robe"),
		Nose:        r.Form.Get("nose"),
		Taste:       r.Form.Get("taste"),
		Dishes:      r.Form.Get("dishes"),
		Season:      r.Form.Get("season"),
	}

	form := forms.New(r.PostForm)

	form.Required("name", "domain", "year", "location", "color")
	if r.Form.Get("appellationName") != "" {
		form.ContentIs("appellationType", []string{"AOC", "AOP"})
	} else {
		form.ContentIs("appellationType", []string{""})
	}
	// TODO: form.ContentIs("location", []string{})
	form.ContentIs("color", []string{"Rouge", "Blanc", "Orange", "Rosé"})
	form.ContentIs("season", []string{"", "Printemps", "Eté", "Automne", "Hiver"})

	if !form.Valid() {
		data := make(map[string]interface{})
		data["wine"] = wine

		m.App.Session.Put(r.Context(), "error", "Ce vin n'a pas pu être ajouté...")
		render.Template(w, r, "add_wine.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertWine(wine)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "flash", "Nouveau vin ajouté !")
	http.Redirect(w, r, "/nouveau-vin", http.StatusSeeOther)
}

// WineMap renders the wine map page
func (m *Repository) WineMap(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "wine_map.page.tmpl", &models.TemplateData{})
}

// WinesList renders the wines page
func (m *Repository) WinesList(w http.ResponseWriter, r *http.Request) {
	wines, err := m.DB.AllWinesSummary()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(wines) == 0 {
		// TODO : No wine
	}

	data := make(map[string]interface{})
	data["wines"] = wines

	render.Template(w, r, "wines-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) WinePage(w http.ResponseWriter, r *http.Request) {
	wineID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	wine, err := m.DB.WineByID(wineID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["wine"] = wine

	render.Template(w, r, "wine.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
