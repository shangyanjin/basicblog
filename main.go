package main

import (
	"./blog"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Template files
const (
	mainTemplate   = "main.html"
	submitTemplate = "submit.html"
)

// templateFunc is a wrapped Handler function associated with a loaded template
type templateFunc func(http.ResponseWriter, *http.Request, *template.Template)

// Cached templates to save disk I/O
var templateCache map[string]*template.Template

// Current state
var blogState *blog.Blog

// Functions exported into templates
var funcMap template.FuncMap = template.FuncMap{
	"formatTime": formatTime,
}

// formatTime formats a Time object into the default RFC822Z representation.
func formatTime(t time.Time) string {
	return t.Format(time.RFC822Z)
}

// makeTemplateHandler loads from disk or from cache the template passed by
// the filename tmpl and creates a new functions that executes fn with the
// loaded and validated template.
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

// mainPage is the main page served on "/"
func mainPage(w http.ResponseWriter, r *http.Request, t *template.Template) {
	t.Execute(w, blogState)
}

// submitPage is the submission page served on "/submit/"
func submitPage(w http.ResponseWriter, r *http.Request, t *template.Template) {
	if r.Method == "POST" {
		if r.FormValue("title") == "" || r.FormValue("content") == "" {
			http.Redirect(w, r, "/submit/", http.StatusFound)
			return
		}

		newEntry := &blog.BlogEntry{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Date:    time.Now(),
		}

		blogState.AddEntry(newEntry)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// deferCleanup listens for SIGIT (Ctrl-C) and saves the state on disk before
// exiting.
func deferCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Ctrl-C (%s) caught, saving state...\n", sig)
			blogState.Save("data/entries.json")
			os.Exit(0)
		}
	}()
}

func main() {
	var err error
	templateCache = make(map[string]*template.Template)

	blogState, err = blog.NewFromFile("data/entries.json")
	if err != nil {
		panic("Blog entries could not be loaded")
	}

	deferCleanup()

	http.HandleFunc("/", makeTemplateHandler(mainPage, mainTemplate))
	http.HandleFunc("/submit/", makeTemplateHandler(submitPage, submitTemplate))
	http.Handle("/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))),
	)

	http.ListenAndServe(":8080", nil)
}
