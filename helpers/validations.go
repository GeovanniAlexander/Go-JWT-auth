package helpers

import (
	"github.com/GeovanniAlexander/02-authJWT/models"
	"gopkg.in/go-playground/validator.v9"
)

func ValidateStruct(user models.User) error {
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return err
	}
	return nil
}
