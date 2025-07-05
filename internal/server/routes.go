package server

import (
	"devtrail/internal/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.HandleHome)
	mux.HandleFunc("/register", handlers.HandleRegister)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/github-auth", handlers.HandleGitHubAuth)
	mux.HandleFunc("/dashboard", handlers.HandleDashboard)
	mux.HandleFunc("/logout", handlers.HandleLogout)
	mux.HandleFunc("/create-project", handlers.HandleCreateProject)
	mux.HandleFunc("/delete-project/", handlers.HandleDeleteProject)
	mux.HandleFunc("/projects/", handlers.HandleProjectPage)
	mux.HandleFunc("/add-commit", handlers.HandleAddCommit)
	mux.HandleFunc("/user/repos", handlers.HandleUserRepos)

}
