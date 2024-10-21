package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *Application) Then(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	t = t.Add(10 * time.Minute)

	w.Header().Set("Content-type", "text/html")
	fmt.Fprintf(w, "%s", t)
}
