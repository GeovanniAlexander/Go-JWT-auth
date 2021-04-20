package controllers

import (
	"strconv"
	"time"

	"github.com/GeovanniAlexander/02-authJWT/database"
	"github.com/GeovanniAlexander/02-authJWT/helpers"
	"github.com/GeovanniAlexander/02-authJWT/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Data struct {
	Success bool          `json:"success"`
	Users   []models.User `json:"users"`
	Errors  []string      `json:"errors"`
}

const SecretKey = "Secret"

func Register(c *fiber.Ctx) error {
	var dataInp map[string]string
	if err := c.BodyParser(&dataInp); err != nil {
		return err
	}
	var data Data = Data{Success: true, Errors: make([]string, 0)}
	password, _ := bcrypt.GenerateFromPassword([]byte(dataInp["password"]), 10)
	user := models.User{
		Name:     dataInp["name"],
		Email:    dataInp["email"],
		Password: password,
	}
	if err := helpers.ValidateStruct(user); err != nil {
		data.Success = false
		data.Errors = append(data.Errors, "The data is required")
		return c.JSON(data)
	}
	if result := database.DB.Create(&user); result.Error != nil {
		data.Success = false
		data.Errors = append(data.Errors, result.Error.Error())
		return c.JSON(data)
	}
	data.Users = append(data.Users, user)

	return c.JSON(data)
}

func Login(c *fiber.Ctx) error {
	var dataInp map[string]string
	if err := c.BodyParser(&dataInp); err != nil {
		return err
	}
	var data Data = Data{Success: true, Errors: make([]string, 0)}
	var user models.User
	database.DB.Where("email = ?", dataInp["email"]).First(&user)
	if user.Id == 0 {
		data.Success = false
		data.Errors = append(data.Errors, "User not found")
		c.Status(fiber.StatusNotFound)
		return c.JSON(data)
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(dataInp["password"])); err != nil {
		data.Success = false
		data.Errors = append(data.Errors, "Incorrect password")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(data)
	}
	data.Users = append(data.Users, user)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		data.Success = false
		data.Errors = append(data.Errors, "Could not logIn")
		return c.JSON(data)
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	data := Data{Success: true, Errors: make([]string, 0)}
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		data.Success = false
		data.Errors = append(data.Errors, err.Error())
		return c.JSON(data)
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)
	data.Users = append(data.Users, user)
	return c.JSON(data)
}

func LogOut(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	data := Data{Success: true, Errors: make([]string, 0)}
	return c.JSON(data)
}
