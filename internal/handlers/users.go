package handlers

import (
	"devtrail/internal/storage"
	"fmt"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case http.MethodGet:
		err := tmpl.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form ", http.StatusBadRequest)
			return
		}
		user, _ := parseLoginForm(r)

		hashedPassword, _ := storage.HashPassword(user.Password)
		user.Password = hashedPassword

		err = RegisterUser(user)
		if err != nil {
			http.Error(w, "Registration error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "%v succesfuly registered", user.Username)

	}
}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, "login.html", nil)
	case http.MethodPost:
		user, err := parseLoginForm(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		valid, err := LoginUser(user)
		if err != nil {
			http.Error(w, "Authorization error", http.StatusUnauthorized)
			return
		}
		if !valid {
			http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
			return
		}
		session, _ := store.Get(r, AuthSessionName)
		session.Values["username"] = user.Username
		session.Options.MaxAge = 3600 // 1 hour
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}

}
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, AuthSessionName)
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
