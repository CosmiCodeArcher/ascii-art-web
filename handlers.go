package main

import (
	"net/http"
	"html/template"
	"ascii-art-web/banner"
	"ascii-art-web/parser"
	"ascii-art-web/render"
)

type PageData struct {
	Text string
	Banner string
	Result string
	Error string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, PageData{Result: ""})
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	LoadedBanner, err := banner.LoadBanner("banners/" + bannerName + ".txt")
	if err != nil {
		http.Error(w, "Banner Not Found", http.StatusNotFound)
		return
	}

	ParsedInput, err := parser.ParseInput(text)
	if err != nil {
		tmpl.Execute(w, PageData{Error: "Empty input - please type something."})
		return
	}

	RenderedOutput := render.RenderToString(LoadedBanner, ParsedInput)
	tmpl.Execute(w, PageData{Result: RenderedOutput, Text: text, Banner: bannerName})
}