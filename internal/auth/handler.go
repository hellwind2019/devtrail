package auth

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

/*
TODO:
ShowHomePage(w http.ResponseWriter, r *http.Request)

ShowLoginPage(w http.ResponseWriter, r *http.Request)

ShowRegisterPage(w http.ResponseWriter, r *http.Request)

HandleLogin(w http.ResponseWriter, r *http.Request)
*/
var tmpl *template.Template
var sessionKey []byte
var store *sessions.CookieStore

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

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == http.MethodGet {
		err := tmpl.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		// Спочатку парсимо форму
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
			return
		}
		user, err := parseLoginForm(r)
		hashedPassword, _ := storage.HashPassword(user.Password)
		user.Password = hashedPassword
		err = RegisterUser(user)
		if err != nil {
			http.Error(w, "Помилка реєстрації", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "%v успішно зареєстрований", user.Username)

	}
}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html", nil)
	} else if r.Method == http.MethodPost {
		user, err := parseLoginForm(r)
		if err != nil {
			http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
			return
		}
		valid, err := LoginUser(user)
		if err != nil {
			http.Error(w, "Помилка авторизації", http.StatusUnauthorized)
			return
		}
		if !valid {
			http.Error(w, "Неправильний логін або пароль", http.StatusUnauthorized)
			return
		}
		session, _ := store.Get(r, "auth-session")
		session.Values["username"] = user.Username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}

}
func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "home.html", nil)
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
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Помилка отримання сесії", http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = 7200
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

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
