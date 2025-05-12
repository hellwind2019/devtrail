package auth

import (
	"devtrail/internal/models"
	"fmt"
	"html/template"
	"net/http"
)

/*
TODO:
ShowHomePage(w http.ResponseWriter, r *http.Request)

ShowLoginPage(w http.ResponseWriter, r *http.Request)

ShowRegisterPage(w http.ResponseWriter, r *http.Request)

HandleLogin(w http.ResponseWriter, r *http.Request)
*/
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

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
		username := r.FormValue("login")
		password := r.FormValue("password")
		valid, err := LoginUser(models.User{Username: username, Password: password})
		if err != nil {
			http.Error(w, "Помилка авторизації", http.StatusUnauthorized)
			return
		}
		if valid {
			http.Redirect(w, r, "/dashboard?username="+username, http.StatusSeeOther)
		} else {
			http.Error(w, "Неправильний логін або пароль", http.StatusUnauthorized)
		}
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
	w.Header().Set("Content-Type", "text/html")

	// тут можна вставити перевірку сесії або токена
	username := r.URL.Query().Get("username")

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
