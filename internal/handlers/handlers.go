package handlers

import (
	"fmt"
	"log"
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

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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

	wine := models.Wine{
		Name:            r.Form.Get("name"),
		Domain:          r.Form.Get("domain"),
		Year:            r.Form.Get("year"),
		AppellationType: r.Form.Get("appellationType"),
		AppellationName: r.Form.Get("appellationName"),
		Location:        r.Form.Get("location"),
		Color:           r.Form.Get("color"),
		Culture:         r.Form.Get("culture"),
		Varieties:       r.Form.Get("varieties"),
		Robe:            r.Form.Get("robe"),
		Nose:            r.Form.Get("nose"),
		Taste:           r.Form.Get("taste"),
		Dishes:          r.Form.Get("dishes"),
		Season:          r.Form.Get("season"),
	}

	form := forms.New(r.PostForm)

	form.Required("name", "domain", "year", "location", "color")
	if r.Form.Get("appellationName") != "" {
		form.Has("appellationType")
		form.ContentIs("appellationType", []string{"AOC", "AOP"})
	} else {
		form.ContentIs("appellationType", []string{""})
	}
	// TODO: form.ContentIs("location", []string{})
	form.ContentIs("color", []string{"Rouge", "Blanc", "Orange", "Rosé", "Champagne"})
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

	htmlMessage := fmt.Sprintf(`
		<strong>New Wine Added !</strong></br>
		</br>
		A new bottle has been added :</br>
		Name : %s</br>
		Domain : %s</br>
		Year : %s</br>
		Location: %s</br>
		Color : %s</br>
	`, wine.Name, wine.Domain, wine.Year, wine.Location, wine.Color)

	msg := models.MailData{
		To:      "clavavin@gmail.com",
		From:    "noreply@clavavin.fr",
		Subject: "New Wine",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "flash", "Nouveau vin ajouté !")
	http.Redirect(w, r, "/nouveau-vin", http.StatusSeeOther)
}

// WineMap renders the wine map page
func (m *Repository) WineMap(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "wine_map.page.tmpl", &models.TemplateData{})
}

// WinesList renders the wines list page
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

	render.Template(w, r, "wines_list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// // WinePage renders the wine page
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

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail(email)
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Email ou mot de passe invalide(s)")
		http.Redirect(w, r, "/mon-compte/connexion", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Connexion réussie !")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/mon-compte/connexion", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin_dashboard.page.tmpl", &models.TemplateData{})
}
