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
