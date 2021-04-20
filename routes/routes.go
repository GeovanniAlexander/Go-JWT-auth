package routes

import (
	"github.com/GeovanniAlexander/02-authJWT/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.LogOut)

}
