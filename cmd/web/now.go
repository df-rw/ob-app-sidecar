package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *Application) Now(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	w.Header().Set("Content-type", "text/html")
	fmt.Fprintf(w, "%s", t)
}
