package server

import (
	"devtrail/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", auth.HandleHome)
	mux.HandleFunc("/register", auth.HandleRegister)
	mux.HandleFunc("/login", auth.HandleLogin)
	mux.HandleFunc("/dashboard", auth.HandleDashboard)
	mux.HandleFunc("/logout", auth.HandleLogout)
	mux.HandleFunc("/create-project", auth.HandleCreateProject)
	mux.HandleFunc("/delete-project/", auth.HandleDeleteProject)
}
