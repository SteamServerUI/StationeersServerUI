package core

import (
	"StationeersServerUI/src/config"
	"net/http"
	"text/template"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./UIMod/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Version: config.Version,
		Branch:  config.Branch,
	}
	if data.Version == "" {
		data.Version = "unknown"
	}
	if data.Branch == "" {
		data.Branch = "unknown"
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
