package handlers

import (
	"github.com/Deo-Mugabe/GOLANG/pkg/render"
	"net/http"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

// About page
func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}
