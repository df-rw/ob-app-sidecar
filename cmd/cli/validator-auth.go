package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

func logger(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		f.ServeHTTP(w, r)
		log.Printf("%s %s %v\n", r.Method, r.URL.String(), time.Since(t))
	}
}

func (app *Application) validatorAuth(w http.ResponseWriter, r *http.Request) {
	var statusCode int

	// TODO Check for the GC IAP header
	// TODO Validate the JWT

	statusCode = http.StatusNoContent // TODO use correct response code

	w.WriteHeader(statusCode)
}

func main() {
	port := flag.Int("p", 8081, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.validatorAuth)

	fmt.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logger(mux)))
}
