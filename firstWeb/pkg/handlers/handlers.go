package handlers

import (
	"github.com/Deo-Mugabe/GOLANG/pkg/config"
	"github.com/Deo-Mugabe/GOLANG/pkg/models"
	"github.com/Deo-Mugabe/GOLANG/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r

}

// Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again."

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

}
