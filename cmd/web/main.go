package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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
		log.Printf("%s: %s %s %v\n", os.Args[0], r.Method, r.URL.String(), time.Since(t))
	}
}

func whoami(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iAm := r.Header.Get("x-goog-authenticated-user-email")

		if iAm == "" {
			log.Println("missing x-goog-authenticated-user-email header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(iAm, "accounts.google.com:") {
			log.Printf("invalid x-goog-authenticated-user-email: %s\n", iAm)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		iAm = strings.Replace(iAm, "accounts.google.com:", "", 1)
		f.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userEmail", iAm)))
	}
}

func main() {
	port := flag.Int("p", 8082, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/now", app.Now)
	mux.HandleFunc("/api/then", app.Then)
	mux.HandleFunc("/api/whoami", app.WhoAmI)

	mux.HandleFunc("GET /api/todos", app.Todos)
	mux.HandleFunc("POST /api/todos/add", app.TodosAdd)
	mux.HandleFunc("POST /api/todos/toggle/{id}", app.TodosToggle)

	fmt.Printf("%s: listening on port %d\n", os.Args[0], *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), whoami(logger(mux))))
}
