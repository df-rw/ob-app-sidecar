package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Application struct {
	tpl *template.Template
}

func New() *Application {
	tpl := template.Must(template.ParseGlob("templates/*.tmpl"))

	return &Application{tpl}
}

func (app *Application) render(w http.ResponseWriter, name string, data map[string]any, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	err := app.tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}

func logger(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		f.ServeHTTP(w, r)
		log.Printf("Backend: %s %s %v\n", r.Method, r.URL.String(), time.Since(t))
	}
}

func main() {
	port := flag.Int("p", 8082, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/now", app.Now)
	mux.HandleFunc("/api/then", app.Then)

	mux.HandleFunc("GET /api/todos", app.Todos)
	mux.HandleFunc("POST /api/todos/add", app.TodosAdd)
	mux.HandleFunc("POST /api/todos/toggle/{id}", app.TodosToggle)

	/*
		mux.HandleFunc("GET /api/dog/name", app.dogName)
		mux.HandleFunc("GET /api/dog/breed", app.dogBreed)
		mux.HandleFunc("GET /api/dog/colour", app.dogColour)
		mux.HandleFunc("GET /api/dog/save", app.dogSave)
	*/

	fmt.Println("Backend: listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logger(mux)))
}
