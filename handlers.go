package main

import (
	"ascii-art-web/banner"
	"ascii-art-web/parser"
	"ascii-art-web/render"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type PageData struct {
	Text string
	Banner string
	Result string
	Error string
}

var (
	ErrInvalidBanner = errors.New("invalid banner name")
	ErrBannerLoad = errors.New("could not load banner")
	ErrEmptyText = errors.New("empty text")
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "The page doesn't exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, PageData{Result: ""})
}

func generate(text, bannerName string) (string, error) {
	// validate, load, parse, render — return the art or an error
	var bannerFilePath string

	if text == "" {
		return "", ErrEmptyText
	}

	switch bannerName {
	case "standard", "shadow", "thinkertoy":
		bannerFilePath = "banners/" + bannerName + ".txt"
	default:
		return "", ErrInvalidBanner
	}

	LoadedBanner, err := banner.LoadBanner(bannerFilePath)
	if err != nil {
		return "", ErrBannerLoad
	}

	ParsedInput, err := parser.ParseInput(text)
	if err != nil {
		return "", ErrEmptyText
	}

	Output := render.RenderToString(LoadedBanner, ParsedInput)
	return Output, nil
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	bannerName := strings.ToLower(r.FormValue("banner"))

	art, err := generate(text, bannerName)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidBanner):
			http.Error(w, "Invalid banner name", http.StatusBadRequest)
			return
		case errors.Is(err, ErrBannerLoad):
			http.Error(w, "Banner Not Found", http.StatusNotFound)
			return
		case errors.Is(err, ErrEmptyText):
			tmpl.Execute(w, PageData{Error: "Empty input - please type something."})
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    		return
		}
	}
	tmpl.Execute(w, PageData{Result: art, Text: text, Banner: bannerName})
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.URL.Query().Get("text")
	bannerName := strings.ToLower(r.URL.Query().Get("banner"))

	art, err := generate(text, bannerName)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidBanner):
			http.Error(w, "Invalid banner name", http.StatusBadRequest)
			return
		case errors.Is(err, ErrBannerLoad):
			http.Error(w, "Banner Not Found", http.StatusNotFound)
			return
		case errors.Is(err, ErrEmptyText):
			http.Error(w, "Missing text or banner name parameter", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    		return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(art)))
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")

	fmt.Fprint(w, art)
}