package handlers

import (
	"devtrail/internal/storage"
	"net/http"
)

func HandleUserRepos(w http.ResponseWriter, r *http.Request) {
	username, ok := getSessionUsername(r)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	token, err := storage.GetGitHubTokenByUsername(username)
	if err != nil || token == "" {
		http.Error(w, `{"error":"GitHub token not found"}`, http.StatusUnauthorized)
		return
	}
	resp, err := MakeGitHubRequest("GET", "https://api.github.com/user/repos?sort=updated", token)

	if err != nil {
		http.Error(w, `{"error":"Failed to fetch repos"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
