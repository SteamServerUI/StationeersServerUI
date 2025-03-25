package api

import (
	"StationeersServerUI/src/config"
	"net/http"
	"os/exec"
	"sync"
	"text/template"
)

var cmd *exec.Cmd
var mu sync.Mutex
var outputChannel chan string
var clients []chan string
var clientsMu sync.Mutex

type Config struct {
	Server struct {
		ExePath  string `xml:"exePath"`
		Settings string `xml:"settings"`
	} `xml:"server"`

	SaveFileName string `xml:"saveFileName"`
}

func StartAPI() {
	outputChannel = make(chan string, 100)
}

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeUI(w http.ResponseWriter, r *http.Request) {
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
