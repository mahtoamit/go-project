package dequeuers

import (
	"encoding/json"
	"fmt"

	// "github.com/tutorialedge/go-fiber-tutorial/data"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
)

// func DequeueEmployeeData() (bool) {

// 	// Calculate the response time
// 	var books models.Book
// 	fmt.Println("yes")

//     count := queries.GetCount()
// 	fmt.Println("count",count)

// 	intial := 0

// 	for intial <= count {

// 	// Retrieve the list values using LRANGE
// cachedData, err := queries.RetrieveData()
// 	if err == nil {

// 		        if cachedData != constants.Index_value {
// 				err := json.Unmarshal([]byte(cachedData), &books)
// 				if err != nil {
// 					return false
// 				}

// 				err = queries.DBCreate(books)
// 				if err != nil {
// 					return false
// 				}

// 				result := queries.RedisDeleteBook(constants.Books_constant)
// 				if !result {
// 					utils.Log("Error", "book", constants.Url_add_book, "dequeuedata", "","error for updating the redis cache")
// 					return false
// 				}

// 			}

// 		}
// 		intial++
// }

// return true
// }
// func main(){
// 		data.DequeueEmployeeData()
// 	}

func DequeueRedisQueue() error {

	var cache models.CacheModel
	count := queries.GetCount()
	fmt.Println("count", count)
	// Retrieve the list values using LRANGE
	cachedData, err := queries.RetrieveData()
	if err != nil {
		return err
	}
	fmt.Println("cachedData", cachedData)

	err = json.Unmarshal([]byte(cachedData), &cache)

	fmt.Println("cache", cache)
	// insert data in db
	tablename := cache.DatabaseName
	qtype := cache.QueryType

	if tablename == "books" && qtype == "insert" {
		err = queries.DBCreate(cache.Model)
		if err != nil {
			return err
		}
	}
	err = queries.DBCreate(cache.Model)
	fmt.Println("cache.Model", cache.Model)
	fmt.Println("err", err)
	return err

}
