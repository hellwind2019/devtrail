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
	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)

	}
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

}
func SaveUserToDB(username, password string) error {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := db.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserFromDB(username string) error {
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

func CheckUserCredentials(username, password string) (bool, error) {
	query := `SELECT password FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return false, nil // User not found
	} else if err != nil {
		return false, err
	}

	// Compare the provided password with the hashed password
	if CheckPasswordHash(password, hashedPassword) {
		return true, nil
	}
	return false, nil
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
