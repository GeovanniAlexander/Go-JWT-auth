package models

type User struct {
	Id       uint
	Name     string
	Email    string `json:"email" validate:"required" gorm:"unique"`
	Password []byte
}
