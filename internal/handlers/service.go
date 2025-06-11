package handlers

import (
	"devtrail/internal/models"
	"devtrail/internal/storage"
)

func RegisterUser(user models.User) error {

	err := storage.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(user models.User) (bool, error) {
	valid, err := storage.AuthenticateUser(user)
	if err != nil {
		return false, err
	}
	return valid, nil
}
