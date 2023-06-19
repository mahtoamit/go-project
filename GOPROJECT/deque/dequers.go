package deque

import (
	"encoding/json"
	"time"
	
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/constants"
	
	"github.com/tutorialedge/go-fiber-tutorial/queries"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)



func DequeueEmployeeData(c *fiber.Ctx) error {
	utils.InitLogger()
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "started", startTime, time.Now())

	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "Data received", startTime, time.Now())
	// Calculate the response time
	var books models.Book


	
	// Retrieve the list values using LRANGE
	cachedData, err := queries.RetrieveData()
	if err == nil {
		
		
		       if cachedData != constants.Index_value {
				err := json.Unmarshal([]byte(cachedData), &books)
				if err != nil {
					utils.Log("ERROR", "book", constants.Url_add_book, "", "dequeuedata", constants.Unmarshal_error+err.Error())
				}
				

				
				err = queries.DBCreate(books)
				if err != nil {
					utils.Log("Error", "book", constants.Url_add_book, "dequeuedata", "", err.Error())
				}

			}else{
				return c.Status(253).JSON(fiber.Map{"error":"The Data consist no values"})
			}

		}
	endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "ended", startTime, endTime)
	return c.JSON("yes")
	}

	

