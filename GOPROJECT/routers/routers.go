package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/handler"
	"github.com/tutorialedge/go-fiber-tutorial/health"
	"github.com/tutorialedge/go-fiber-tutorial/auth"
	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/book"
	

)




func SetupRoutes(app *fiber.App) {
	app.Post("/api/v1/login", handler.Login)
	app.Post("/api/v1/signup", handler.Signup)
	app.Post("/api/v1/logout", handler.Logout)
	app.Get("/api/v1/health-check", health.Connect)
	app.Get("/api/v1/book", auth.Protected, book.Getbooks)
	app.Get(constants.Url_book, auth.Protected, book.Getbook)
	app.Post("/api/v1/book/", auth.Protected, book.Newbook)
	app.Put(constants.Url_book, auth.Protected, book.UpdateBook)
	app.Delete(constants.Url_book, auth.Protected, book.Deletebook)
}