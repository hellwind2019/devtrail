package auth

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
)

func RegisterUser(user models.User) error {

	err := storage.SaveUserToJson(user)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(user models.User) (bool, error) {
	valid, err := storage.CheckUserCredentialsJson(user)
	if err != nil {
		return false, err
	}
	return valid, nil
}
