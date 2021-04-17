package database

import (
	"github.com/GeovanniAlexander/02-authJWT/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	connection, err := gorm.Open(mysql.Open("root:Jumanji.1997*@/gocrud"), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})

}
