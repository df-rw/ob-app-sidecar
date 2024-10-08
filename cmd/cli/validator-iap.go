package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
)

const headerIAP = "X-Goog-IAP-JWT-Assertion"
const headerUserEmail = "X-Goog-Authenticated-User-Email"

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
		fmt.Printf("%s: %s %s %v\n", os.Args[0], r.Method, r.URL.String(), time.Since(t))
	}
}

func (app *Application) validateIAP(w http.ResponseWriter, r *http.Request) {
	iapJWT := r.Header.Get(headerIAP)
	iapUserEmail := r.Header.Get(headerUserEmail)

	if iapJWT == "" {
		fmt.Printf("missing header %s\n", headerIAP)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// https://cloud.google.com/iap/docs/identity-howto#getting_the_users_identity_with_signed_headers
	if iapUserEmail == "" {
		fmt.Printf("missing header %s\n", headerUserEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, iapJWT, app.Audience)
	if err != nil {
		fmt.Println(fmt.Errorf("idtoken.Validate: %w", err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(iapUserEmail, "accounts.google.com:") {
		fmt.Printf("missing accounts.google.com prefix on email header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwtEmail := payload.Claims["email"].(string)
	if jwtEmail != strings.Replace(iapUserEmail, "accounts.google.com:", "", 1) {
		fmt.Println("email mismatch: JWT %s, header %s\n", jwtEmail, iapUserEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := flag.Int("p", 8081, "webserver port")
	flag.Parse()

	app := New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.validateIAP)

	fmt.Printf("%s: listening on port %d (audience '%s')\n", os.Args[0], *port, app.Audience)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logger(mux)))
}
