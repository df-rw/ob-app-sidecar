// validator-dummy is a dummy validation application. Use this as a base for
// performing the appropriate validation steps for your use case.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
		log.Printf("%s: %s %s %v\n", os.Args[0], r.Method, r.URL.String(), time.Since(t))
	}
}

func (app *Application) validatorAuth(w http.ResponseWriter, r *http.Request) {
	// Validate the request however you like.
	statusCode := http.StatusNoContent

	w.WriteHeader(statusCode)
}

func main() {
	port := flag.Int("p", 8081, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.validatorAuth)

	fmt.Printf("%s: listening on port %d\n", os.Args[0], *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logger(mux)))
}
