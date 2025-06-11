package handlers

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
	"fmt"
	"net/http"
	"strconv"
)

func HandleAddCommit(w http.ResponseWriter, r *http.Request) {
	_, err := getSessionUser(r, w)
	if err {
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	commitMessage := r.FormValue("message")
	projectIdStr := r.FormValue("project_id")
	ratingStr := r.FormValue("rating")
	projectId, _ := strconv.Atoi(projectIdStr)
	rating, _ := strconv.Atoi(ratingStr)

	if commitMessage == "" {
		http.Error(w, "Commit message cannot be empty", http.StatusBadRequest)
		return
	}

	commit := models.Commit{
		ProjectId: projectId,
		Message:   commitMessage,
		Rating:    rating,
	}
	storage.AddCommit(commit)
	http.Redirect(w, r, fmt.Sprintf("/projects/%d", projectId), http.StatusSeeOther)

}
