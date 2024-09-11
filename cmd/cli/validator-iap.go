package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/idtoken"
)

type Application struct {
	Audience string
}

func New() *Application {
	Audience := os.Getenv("GCP_IAP_JWT_AUDIENCE")

	return &Application{
		Audience,
	}
}

func logger(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		f.ServeHTTP(w, r)
		log.Printf("%s: %s %s %v\n", os.Args[0], r.Method, r.URL.String(), time.Since(t))
	}
}

func (app *Application) validatorAuth(w http.ResponseWriter, r *http.Request) {
	iapJWT := r.Header.Get("X-Goog-IAP-JWT-Assertion")

	if iapJWT == "" {
		log.Printf("missing IAP header\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var statusCode int

	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, iapJWT, app.Audience)
	if err != nil {
		statusCode = http.StatusUnauthorized
		log.Println(fmt.Errorf("idtoken.Validate: %w", err))
	} else {
		statusCode = http.StatusNoContent
		log.Printf("payload: %v", payload)
	}

	w.WriteHeader(statusCode)
}

func main() {
	port := flag.Int("p", 8081, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.validatorAuth)

	fmt.Printf("%s: listening on port %d (audience '%s')\n", os.Args[0], *port, app.Audience)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logger(mux)))
}
