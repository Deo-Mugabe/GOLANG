// package doesnt need to import anything but available to be used were it is needed
package config

import "html/template"

// AppConfig holds the application config
type AppConfig struct {
	TemplateCache map[string]*template.Template // variable name starts with capital letter public
}
