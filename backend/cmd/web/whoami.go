package main

import "net/http"

func (app *Application) WhoAmI(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]any{
		"IAm": r.Context().Value("userEmail"),
	}

	app.render(w, "whoami", pageData, http.StatusOK)
}
