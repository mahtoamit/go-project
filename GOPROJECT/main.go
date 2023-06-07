package main

import (
	
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tutorialedge/go-fiber-tutorial/auth"
	"github.com/tutorialedge/go-fiber-tutorial/book"
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/handler"
	"github.com/tutorialedge/go-fiber-tutorial/health"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"github.com/tutorialedge/go-fiber-tutorial/configs"
)




// var startTime time.Time
var bookDataChannel chan models.Book

func init() {
	bookDataChannel = make(chan models.Book, 100)
}

func SetupRoutes(app *fiber.App) {
	app.Post("/api/v1/login", handler.Login)
	app.Post("/api/v1/signup", handler.Signup)
	app.Post("/api/v1/logout", handler.Logout)
	app.Get("/api/v1/health-check", health.Connect)
	app.Get("/api/v1/book", auth.Protected(), book.Getbooks)
	app.Get(configs.Url_book, auth.Protected(), book.Getbook)
	app.Post("/api/v1/book/", auth.Protected(), book.Newbook(bookDataChannel))
	app.Put(configs.Url_book, auth.Protected(), book.UpdateBook(bookDataChannel))
	app.Delete(configs.Url_book, auth.Protected(), book.Deletebook(bookDataChannel))
}

func main() {
	utils.InitLogger()
	app := fiber.New()
	startTime := time.Now()
    utils.Log("INFO", "Main", "", "","main", "Application starting",startTime, time.Now())
	
	app.Use(logger.New())

	database.DBConnect()
	database.RedisConnect()

	utils.Log("INFO", "Main", "","", "main", "Database initialized",startTime, time.Now())

	SetupRoutes(app)

    // Set up a channel to capture interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	serverShutdown := make(chan struct{})
    
	// Start a separate goroutine that will listen for the interrupt signal
	go func() {
		_ = <-c //wait for interrupt signal
		utils.Log("INFO", "Main", "","", "main", "Graceful shutdown",startTime, time.Now())

		_ = app.Shutdown()
		serverShutdown <- struct{}{}
	   
		// Close the database and Redis connections
		database.CloseDatabase()
		database.ReddisClose()

		utils.Log("INFO", "Main", "","", "main", "Application has been stopped",startTime, time.Now())
	
	}()

	// ...

	if err := app.Listen(":3000"); err != nil {
		log.Panic(err)
	}

	<-serverShutdown

	

}
