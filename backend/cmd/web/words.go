package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *Application) Words(w http.ResponseWriter, r *http.Request) {
	app.render(w, "words-page", map[string]any{}, http.StatusOK)
}

func (app *Application) WordsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)

	http.Redirect(w, r, "/words", http.StatusSeeOther)
}

func (app *Application) WordsSearch(w http.ResponseWriter, r *http.Request) {
	// Randomly stolen from /usr/share/dict/words
	words := []string{
		"pterygophore", "viscousness", "unfascinating", "bedrid", "myoliposis",
		"Roxanne", "torpify", "opiner", "melon", "Bononian",
		"Telei", "crankcase", "prediluvian", "girlie", "pulmotracheary",
		"napoo", "Curt", "orchiocele", "shafty", "calcination",
		"sis", "telemechanism", "Dinotheriidae", "sexuparous", "paralepsis",
		"skal", "linamarin", "headily", "Aludra", "celioscope",
		"uninnocently", "haustellated", "might", "preaffection", "anthranoyl",
		"smorgasbord", "alcoholization", "inaccentuation", "cupay", "tintarron",
		"fassalite", "pantomimish", "unhesitatingness", "pretangible", "spillover",
		"individualizingly", "Hades", "outdress", "soupspoon", "pompless",
		"Melogrammataceae", "magnetogram", "Cardamine", "dermathemia", "ideologically",
		"symbolistic", "reverable", "datum", "foxtongue", "chasm",
		"prothrombogen", "prepense", "clientship", "zoozoo", "mandariness",
		"Samsonian", "ovum", "transigent", "invigoratingly", "thrill",
		"preconfuse", "heritability", "unstaunch", "cooter", "phagocytolytic",
		"aegagropile", "vanity", "prehesitation", "outinvent", "hyperspace",
		"farer", "becompliment", "sloka", "racketry", "pentadecylic",
		"refrigeratory", "brickset", "surrender", "apobiotic", "unquenchable",
		"enthraller", "grasswards", "peromelus", "flagellantism", "quotationist",
		"terebratulid", "Egeria", "colonopathy", "staphylinic", "excecate",
	}

	results := []string{}

	r.ParseForm()
	partial := r.FormValue("partial")
	chooser := r.PathValue("chooser")

	if len(partial) > 0 {
		for _, v := range words {
			if strings.Contains(v, partial) {
				results = append(results, v)
			}
		}
	}

	if chooser == "" {
		chooser = "single"
	}

	pageData := map[string]any{
		"Results": results,
	}

	block := "words-search-results-radio"
	if chooser == "multiple" {
		block = "words-search-results-checkbox"
	}

	app.render(w, block, pageData, http.StatusOK)
}
