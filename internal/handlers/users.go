package handlers

import (
	"devtrail/internal/storage"
	"fmt"
	"net/http"
)

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
		user, _ := parseLoginForm(r)

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
		session.Options.MaxAge = 3600 // 1 година
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}

}
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Помилка отримання сесії", http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
