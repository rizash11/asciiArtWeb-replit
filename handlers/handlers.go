package handlers

import (
	"html/template"
	"net/http"

	"asciiArtWeb/asciiArt"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// w.WriteHeader(http.StatusNotFound)
		http.Error(w, "404 Error Status Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Error 405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("ui/index.html")
	if err != nil {
		return
	}

	tmpl.Execute(w, nil)
}

func AsciiArtWeb(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		http.Error(w, "Error 404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		if r.Method == http.MethodGet {
			http.Error(w, "Error 404 Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error 405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	banner := r.FormValue("banner")
	text := r.FormValue("input")

	res, status := asciiArt.AsciiArt(banner, text)
	switch status {
	case http.StatusBadRequest:
		http.Error(w, "400 Error Bad Request", status)
		return
	case http.StatusInternalServerError:
		http.Error(w, "500 Internal Server Error", status)
		return
	case http.StatusNotFound:
		http.Error(w, "404 Not Found", status)
		return
	}

	tmpl, err := template.ParseFiles("ui/index.html")
	if err != nil {
		return
	}

	tmpl.Execute(w, res)
}
