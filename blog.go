package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	mainTemplate   = "main.html"
	submitTemplate = "submit.html"
)

type templateFunc func(http.ResponseWriter, *http.Request, *template.Template)

type blogEntry struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
}

type mainContent struct {
	Entries []blogEntry
}

var templateCache map[string]*template.Template
var blog mainContent
var funcMap template.FuncMap = template.FuncMap{
	"formatTime": FormatTime,
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC822Z)
}

func makeTemplateHandler(fn templateFunc, tmpl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tmp *template.Template

		if val, ok := templateCache[tmpl]; ok {
			tmp = val
		} else {
			var err error

			tmp, err = template.New(tmpl).Funcs(funcMap).ParseFiles(tmpl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			templateCache[tmpl] = tmp
		}

		fn(w, r, tmp)
	}
}

func readEntries() (data mainContent, err error) {
	f, err := os.Open("data/entries.json")
	if err != nil {
		return
	}
	defer f.Close()

	var e mainContent

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&e)
	if err != nil {
		return
	}

	return e, nil
}

func writeEntries(content mainContent) error {
	f, err := os.Create("data/entries.json")
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.Encode(content)

	return nil
}

func mainPage(w http.ResponseWriter, r *http.Request, t *template.Template) {
	t.Execute(w, blog)
}

func submitPage(w http.ResponseWriter, r *http.Request, t *template.Template) {
	if r.Method == "POST" {
		if r.FormValue("title") == "" || r.FormValue("content") == "" {
			http.Redirect(w, r, "/submit/", http.StatusFound)
			return
		}

		newEntry := blogEntry{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Date:    time.Now(),
		}

		if len(blog.Entries) > 0 {
			newEntry.ID = blog.Entries[0].ID + 1
			blog.Entries = append([]blogEntry{newEntry}, blog.Entries...)
		} else {
			newEntry.ID = 1
			blog.Entries = append(blog.Entries, newEntry)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func deferCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Ctrl-C (%s) caught, saving state...\n", sig)
			writeEntries(blog)
			os.Exit(0)
		}
	}()
}

func main() {
	var err error
	templateCache = make(map[string]*template.Template)

	deferCleanup()

	if blog, err = readEntries(); err != nil {
		panic("Blog entries could not be loaded")
	}

	http.HandleFunc("/", makeTemplateHandler(mainPage, mainTemplate))
	http.HandleFunc("/submit/", makeTemplateHandler(submitPage, submitTemplate))
	http.Handle("/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))),
	)

	http.ListenAndServe(":8080", nil)
}
