package main

import (
	"fmt"
	"strings"
	"strconv"
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

	var bannerFilePath string
	bannerName = strings.ToLower(bannerName)

	switch bannerName {
	case "standard", "shadow", "thinkertoy":
		bannerFilePath = "banners/" + bannerName + ".txt"
	default:
		http.Error(w, "Invalid banner name", http.StatusBadRequest)
		return
	}

	LoadedBanner, err := banner.LoadBanner(bannerFilePath)
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

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.URL.Query().Get("text")
	bannerName := r.URL.Query().Get("banner")

	var bannerFilePath string
	bannerName = strings.ToLower(bannerName)

	switch bannerName {
	case "standard", "shadow", "thinkertoy":
		bannerFilePath = "banners/" + bannerName + ".txt"
	default:
		http.Error(w, "Invalid banner name", http.StatusBadRequest)
		return
	}

	LoadedBanner, err := banner.LoadBanner(bannerFilePath)
	if err != nil {
		http.Error(w, "Banner Not Found", http.StatusNotFound)
		return
	}

	ParsedInput, err := parser.ParseInput(text)
	if err != nil {
		http.Error(w, "Empty text", http.StatusBadRequest)
		return
	}

	Output := render.RenderToString(LoadedBanner, ParsedInput)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(Output)))
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")

	fmt.Fprintln(w, Output)
}