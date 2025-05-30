package auth

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

func HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	userId, err := storage.GetUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Не вдалося отримати ID користувача", http.StatusInternalServerError)
		return
	}
	err = storage.CreateProject(models.Project{
		UserID:      userId,
		Name:        name,
		Description: description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
func HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	// Перевірка сесії користувача
	session, _ := store.Get(r, "auth-session")
	username, ok := session.Values["username"].(string)
	fmt.Println("Проект з ID успішно видалено")
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Отримання ID проекту з URL
	projectID := r.URL.Path[len("/delete-project/"):]
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	// Видалення проекту з бази даних
	projectIDint, error := strconv.Atoi(projectID)
	if error != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}
	err := storage.DeleteProjectByID(projectIDint)
	if err != nil {
		http.Error(w, "Error deleting project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func HandleProjectPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Отримання ID проекту з URL
	projectID := path.Base(r.URL.Path)
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	projectIDint, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Отримання ID користувача з бази даних
	userID, err := storage.GetUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Error retrieving user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перевірка, чи проект належить користувачу
	project, err := storage.GetProjectByID(projectIDint)
	if err != nil {
		http.Error(w, "Error retrieving project: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if project.UserID != userID {
		http.Error(w, "Access denied: You do not own this project", http.StatusForbidden)
		return
	}

	// Передача даних у шаблон
	data := struct {
		Project models.Project
	}{
		Project: project,
	}

	renderTemplate(w, "project.html", data)
}
