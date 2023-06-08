package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tutorialedge/go-fiber-tutorial/database"
	
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"github.com/tutorialedge/go-fiber-tutorial/routers"
	
)




// var startTime time.Time



func main() {
	utils.InitLogger()
	app := fiber.New()
	startTime := time.Now()
    utils.Log("INFO", "Main", "", "","main", "Application starting",startTime, time.Now())
	
	app.Use(logger.New())

	database.DBConnect()
	database.RedisConnect()

	utils.Log("INFO", "Main", "","", "main", "Database initialized",startTime, time.Now())

	routers.SetupRoutes(app)

    // Set up a channel to capture interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	serverShutdown := make(chan struct{})
    
	// Start a separate goroutine that will listen for the interrupt signal
	go func() {
		_ = <-c //wait for interrupt signal
		utils.Log("INFO", "Main", "","", "main", "Graceful shutdown",startTime, time.Now())
        fmt.Println("Gracefulshutdown")
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
