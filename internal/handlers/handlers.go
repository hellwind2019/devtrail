package handlers

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var tmpl *template.Template
var sessionKey []byte
var store *sessions.CookieStore
var AuthSessionName = "auth-session"

func init() {
	// Завантажуємо змінні середовища з файлу .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Помилка завантаження .env файлу")
	}
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	sessionKey = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
	}
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := struct {
			OAuthInfo models.OAuthInfo
		}{
			OAuthInfo: models.OAuthInfo{
				ClientID:     os.Getenv("GH_BASIC_CLIENT_ID"),
				ClientSecret: os.Getenv("GH_BASIC_CLIENT_SECRET"),
			},
		}
		renderTemplate(w, "home.html", data)
	}
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	username, ok := getSessionUsername(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, err := storage.GetUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		return
	}

	projects, err := storage.GetProjectsByUserID(userID)

	if err != nil {
		http.Error(w, "Error retrieving projects", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Projects []models.Project
	}{
		Username: username,
		Projects: projects,
	}

	renderTemplate(w, "dashboard.html", data)
}
