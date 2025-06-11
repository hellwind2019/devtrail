package handlers

import (
	"devtrail/internal/models"

	"net/http"
)

func renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getSessionUsername(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "auth-session")
	username, ok := session.Values["username"].(string)
	return username, ok && username != ""
}

func parseLoginForm(r *http.Request) (models.User, error) {
	if err := r.ParseForm(); err != nil {
		return models.User{}, err
	}
	return models.User{
		Username: r.FormValue("login"),
		Password: r.FormValue("password"),
	}, nil
}
