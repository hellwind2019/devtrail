package handlers

import (
	"devtrail/internal/models"
	"path"
	"strconv"

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
func GetCurrentProjectId(r *http.Request, w http.ResponseWriter) (int, bool) {
	projectID := path.Base(r.URL.Path)
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return 0, true
	}

	projectIDint, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return 0, true
	}
	return projectIDint, false
}
