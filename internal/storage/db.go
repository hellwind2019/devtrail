package storage

import (
	"database/sql"
	"devtrail/internal/models"
	"encoding/json"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)

	}
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );
	
	CREATE TABLE IF NOT EXISTS projects (
   	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	userId int,
		ProjectName TEXT NOT NULL,
		ProjectDescription TEXT,
    	FOREIGN KEY (userId) REFERENCES users(id)
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

}
func SaveUserToDB(user models.User) error {
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

func DeleteUserFromDB(username string) error {
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

func CheckUserCredentialsDB(user models.User) (bool, error) {
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
func SaveUserToJson(user models.User) error {
	users := loadUsers()
	users = append(users, user)
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("users.json", data, 0644)
	return err
}

func CheckUserCredentialsJson(user models.User) (bool, error) {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	users := loadUsers()
	for _, savedUser := range users {
		if savedUser.Username == user.Username {
			if CheckPasswordHash(user.Password, savedUser.Password) {
				return true, nil
			}
		}
	}
	return false, nil
}

func loadUsers() []models.User {
	var users []models.User
	file, err := os.ReadFile("users.json")
	if err == nil {
		_ = json.Unmarshal(file, &users)
	}
	return users
}
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		}
	}
}
