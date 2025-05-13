package auth

import (
	"devtrail/internal/models"
	"fmt"
	"html/template"
	"net/http"
	"os"

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
		username := r.FormValue("login")
		password := r.FormValue("password")

		err = RegisterUser(models.User{Username: username, Password: password})
		if err != nil {
			http.Error(w, "Помилка реєстрації", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "%v успішно зареєстрований", username)

	}
}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
	}
	fmt.Println("Session key: ", sessionKey)
	w.Header().Set("Content-Type", "text/html")
	if r.Method == http.MethodGet {
		err := tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
			return
		}
		user := models.User{
			Username: r.FormValue("login"),
			Password: r.FormValue("password"),
		}
		valid, err := LoginUser(user)

		if err != nil {
			http.Error(w, "Помилка авторизації", http.StatusUnauthorized)
			return
		}
		if valid {
			session, _ := store.Get(r, "auth-session")
			fmt.Printf("IsNew: %v \n", session.IsNew)
			session.Values["username"] = user.Username
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(session.Values["username"])
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Неправильний логін або пароль", http.StatusUnauthorized)
		}
		// if valid {
		// 	http.Redirect(w, r, "/dashboard?username="+user.Username, http.StatusSeeOther)
		// } else {
		// 	http.Error(w, "Неправильний логін або пароль", http.StatusUnauthorized)
		// }
	}

}
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == http.MethodGet {
		err := tmpl.ExecuteTemplate(w, "home.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, ok := session.Values["username"].(string)
	fmt.Printf("Ok: %v username: %v\n", ok, username)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := struct {
		Username string
	}{
		Username: username,
	}

	err := tmpl.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
