package storage

import (
	"devtrail/internal/models"
	"encoding/json"
	"os"
)

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

func CheckUserCredentials(user models.User) (bool, error) {
	users := loadUsers()
	for _, savedUser := range users {
		if compareUsers(savedUser, user) {
			return true, nil
		}
	}
	return false, nil
}
func compareUsers(user1, user2 models.User) bool {
	if user1.Username != user2.Username {
		return false
	}
	if user1.Password != user2.Password {
		return false
	}
	return true
}
func loadUsers() []models.User {
	var users []models.User
	file, err := os.ReadFile("users.json")
	if err == nil {
		_ = json.Unmarshal(file, &users)
	}
	return users
}
