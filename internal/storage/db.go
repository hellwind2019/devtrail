package storage

import (
	"database/sql"
	"devtrail/internal/models"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	schema, err := os.ReadFile("./internal/storage/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema file:", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to execute schema:", err)
	}

	log.Println("Database initialized successfully.")
}

func SaveUser(user models.User) error {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := db.Exec(query, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(username string) error {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `DELETE FROM users WHERE username = ?`
	result, err := db.Exec(query, username)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
func GetProjectsByUserID(userID int) ([]models.Project, error) {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `SELECT id, userId, ProjectName, ProjectDescription FROM projects WHERE userId = ? ORDER BY id DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ProjectID, &project.UserID, &project.Name, &project.Description); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func GetProjectByID(projectID int) (models.Project, error) {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `SELECT id, userId, ProjectName, ProjectDescription FROM projects WHERE id = ?`
	row := db.QueryRow(query, projectID)

	var project models.Project
	err := row.Scan(&project.ProjectID, &project.UserID, &project.Name, &project.Description)
	if err == sql.ErrNoRows {
		return models.Project{}, errors.New("project not found")
	} else if err != nil {
		return models.Project{}, err
	}
	return project, nil
}
func DeleteProjectByID(projectID int) error {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `DELETE FROM projects WHERE id = ?`
	result, err := db.Exec(query, projectID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}
func GetUserIDByUsername(username string) (int, error) {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `SELECT id FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	var userID int
	err := row.Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found")
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}

func AuthenticateUser(user models.User) (bool, error) {
	query := `SELECT password FROM users WHERE username = ?`
	row := db.QueryRow(query, user.Username)

	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return false, nil // User not found
	} else if err != nil {
		return false, err
	}

	// Compare the provided password with the hashed password
	if CheckPasswordHash(user.Password, hashedPassword) {
		return true, nil
	}
	return false, nil
}
func CreateProject(project models.Project) error {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	query := `INSERT INTO projects (userId, ProjectName, ProjectDescription) VALUES (?, ?, ?)`
	_, err := db.Exec(query, project.UserID, project.Name, project.Description)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		}
	}
}
