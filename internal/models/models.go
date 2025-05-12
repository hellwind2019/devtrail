package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Repository struct {
	ID       int
	UserID   int
	Owner    string
	RepoName string
}

type CommitReport struct {
	ID         int
	Repository int
	CommitHash string
	Message    string
	Date       string
}
