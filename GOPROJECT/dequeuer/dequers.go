package dequeuer

import (
	"time"

	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)
var bookDataChannel chan models.Book



func DequeueEmployeeData() {
	utils.InitLogger()
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "started", startTime, time.Now())

	book := <- bookDataChannel
	
	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "Data received", startTime, time.Now())
	// Calculate the response time
	StartTime := time.Now()

	db :=  database.DB.Db
	err := db.Create(&book).Error
	if err != nil {
		utils.Log("Error", "book", constants.Url_add_book, "dequeuedata", "", err.Error(), StartTime, time.Now())
	}

	endTime := time.Now()

	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "ended", startTime, endTime)

	}

