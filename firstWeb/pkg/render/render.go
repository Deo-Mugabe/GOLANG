package render

import (
	"bytes"
	"github.com/Deo-Mugabe/GOLANG/pkg/config"
	"github.com/Deo-Mugabe/GOLANG/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the templates package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders a template
/*
func RenderTemplate(w http.ResponseWriter, tmpl string) {: This line declares a function named RenderTemplate
with two parameters:
w http.ResponseWriter: This parameter represents the HTTP response writer. It's used to write the rendered
template to the HTTP response.
tmpl string: This parameter is a string representing the name of the template to render.
*/

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	// get the template cache from the app config
	tc := app.TemplateCache
	// create a template cache
	/*
		tc, err := createTemplateCache(): This line calls the createTemplateCache function to create a cache of
		templates. It returns a map of template names to template objects and an error.
	*/
	//tc, err := CreateTemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// get requested template from cache
	/*
		t, ok := tc[tmpl]: This line attempts to retrieve the requested template from the cache (tc).
		It uses the tmpl string as the key to access the template from the map.
	*/
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}
	/*
		buf := new(bytes.Buffer): This line creates a new buffer (bytes.Buffer) to store the rendered
		template content.
	*/
	buf := new(bytes.Buffer)
	/*
		err = t.Execute(buf, nil): This line executes the template (t) and writes the result to the buffer (buf).
		   The nil value indicates that no data is being passed to the template for rendering.
	*/
	_ = t.Execute(buf, td)
	//if err != nil {
	//	log.Println(err)
	//}

	// render the template
	/*
		_, err = buf.WriteTo(w): This line writes the content of the buffer (buf), which contains the rendered
		template, to the HTTP response writer (w). It returns the number of bytes written and any error encountered.
	*/
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

/*
func createTemplateCache() (map[string]*template.Template, error) {: This line declares a function named
createTemplateCache that returns a map of strings to pointers of template.Template and an error.
*/
func CreateTemplateCache() (map[string]*template.Template, error) {
	/*
		myCache := map[string]*template.Template{}: This line initializes an empty map called myCache to
		store templates. The keys are strings representing template names, and the values are pointers to template.
		Template objects.
	*/
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{} //empty map

	// get all of the files named *.page.tmpl from ./templates

	/*
		pages, err := filepath.Glob("./templates/*.page.tmpl"): This line uses filepath.Glob to find all files
		matching the pattern *.page.tmpl in the ./templates directory. These files represent page-specific
		templates.
		if err != nil { return myCache, err }: This line checks if there was an error while finding the template
		files. If there was, it returns the empty cache and the error.
	*/
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	/*
		This part of the code iterates over each page template file found:
		for _, page := range pages {: This line iterates over each file found in pages.
		name := filepath.Base(page): This line extracts the base name of the file, which will be used as the
		name of the template.
		ts, err := template.New(name).ParseFiles(page): This line creates a new template with the specified name
		and parses the content of the template file using ParseFiles. The parsed template is stored in the variable ts.
		if err != nil { return myCache, err }: This line checks if there was an error while parsing the template file.
		If there was, it returns the cache and the error.
	*/
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		/*
			matches, err := filepath.Glob("./templates/*.layout.tmpl"): This line finds all files matching the
			pattern *.layout.tmpl in the ./templates directory. These files represent layout templates.
			if err != nil { return myCache, err }: This line checks if there was an error while finding the layout
			template files. If there was, it returns the cache and the error.
		*/
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		/*
			if len(matches) > 0 {: This line checks if there are any layout template files found.
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl"): This line parses all layout template files using
			ParseGlob and adds them to the existing template (ts).
			if err != nil { return myCache, err }: This line checks if there was an error while parsing the layout
			template files. If there was, it returns the cache and the error.
		*/
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		/*
			myCache[name] = ts: This line adds the parsed template (ts) to the cache (myCache) using the template
			name (name) as the key.
		*/
		myCache[name] = ts
	}
	/*
	      return myCache, nil: This line returns the populated cache and nil error, indicating that there were
	   no errors during the caching process.
	*/
	return myCache, nil
}
