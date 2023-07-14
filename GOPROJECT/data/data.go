package data

import (
	"encoding/json"
	"fmt"

	"github.com/tutorialedge/go-fiber-tutorial/constants"

	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)

func DequeueEmployeeData() bool {

	// Calculate the response time
	// var books models.Book
	var cache models.CacheModel
	fmt.Println("yes")

	count := queries.GetCount()
	fmt.Println("count", count)

	intial := 0

	for intial <= count {

		// Retrieve the list values using LRANGE
		cachedData, err := queries.RetrieveData()
		if err == nil {

			if cachedData != constants.Index_value {
				err := json.Unmarshal([]byte(cachedData), &cache)
				if err != nil {
					return false
				}
				fmt.Println("cache", cache)

				// err = queries.DBCreate(books)
				// if err != nil {
				// 	return false
				// }

				result := queries.RedisDeleteBook(constants.Books_constant)
				if !result {
					utils.Log("Error", "book", constants.Url_add_book, "dequeuedata", "", "error for updating the redis cache")
					return false
				}

			}

		}
		intial++
	}

	return true
}
