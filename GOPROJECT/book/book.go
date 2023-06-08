package book


import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
)

// const (
// 	url_get_books = "/api/v1/books"
// 	url_get_single_book = "/api/v1/book/:title"
// 	url_add_book = "/api/v1/book/"
// 	url_update_book ="/api/v1/book/:id"
// 	unmarshal_error = "Error unmarshaling cached books data:"
// 	Error_caching_data = "Error caching books data:"
// 	Error_deleting_cached_data = "Error deleting employees cache: "
// )

var ctx = context.Background()


var bookDataChannel chan models.Book

func init() {
	bookDataChannel = make(chan models.Book, 100)
	go dequeueEmployeeData()
}
var startTime time.Time



func Getbooks( c *fiber.Ctx) error {
	utils.InitLogger()
	userId := c.Locals("userId").(string)
	startTime := time.Now() // Set the start time before processing the request
	utils.Log("INFO", "book", constants.Url_get_books,userId, "Getbooks", "started", startTime, time.Now())
	// Check if the data exists in the Redis cache
	
	key := "books"
	data:= queries.RedisCache(key,userId)
	
	if data != nil {
		endTime := time.Now()
	    utils.Log("INFO", "book", constants.Url_get_books,userId, "Getbooks", "ended", startTime,endTime)
        return c.Status(200).JSON(fiber.Map{"data":data})
	}
	
	// Data does not exist in the cache, fetch from the database
	books := queries.DBGetBooks()

	result , err := queries.RedisSetCache(books ,key,userId)
	
    if err !=nil && !result {
		utils.Log("DEBUG", "book", constants.Url_get_books,userId, "Getbooks", "ended")
		return c.Status(252).JSON(fiber.Map{"msg":err} )
	}
	endTime := time.Now()
    utils.Log("INFO", "book", constants.Url_get_books,userId, "Getbooks", "ended",startTime, endTime)
	return c.Status(200).JSON(fiber.Map{"data":&books})

}

func Getbook( c *fiber.Ctx) error {
	utils.InitLogger()
	startTime = time.Now()
	userId := c.Locals("userId").(string)
	utils.Log("INFO", "book", constants.Url_get_single_book,userId, "Getbook", "started",startTime, time.Now())
	name := c.Params("title")
	// Check if the data exists in the Redis cache
	titlekey := fmt.Sprintf(constants.Redis_book_const, name)
	data:= queries.RedisCacheGetBook(titlekey,userId)
	
	if data != nil {
		
		endTime := time.Now()
	    utils.Log("INFO", "book", constants.Url_get_single_book,userId, "Getbooks", "ended", startTime,endTime)
        return c.Status(200).JSON(fiber.Map{"data":data})
	}
    

	// Data does not exist in the cache, fetch from the database
	books := queries.DBGetSingleBook(name)
	

	if books.Title == "" {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book,userId, "Getbook", "ended",startTime, endTime)
		return c.Status(253).JSON("No book found for  this title")
	}

    
	result ,err := queries.RedisSetCacheBook(books ,titlekey,userId)
	
    if err !=nil && !result {
		utils.Log("DEBUG", "book", constants.Url_get_single_book,userId, "Getbook", "ended")
		return c.Status(252).JSON(fiber.Map{"msg":err})
	}


	endTime := time.Now()
    utils.Log("INFO", "book", constants.Url_get_single_book,userId, "Getbook", "ended",startTime, endTime)
	return c.Status(200).JSON(fiber.Map{"data":&books})

}

func Newbook(c *fiber.Ctx) error {
	
	utils.InitLogger()
	userId := c.Locals("userId").(string)
	
	// Set the start time before processing the request
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book,userId, "Newbook", "started", startTime, time.Now())
	var validate = validator.New()
	redisClient := database.RedisClient
	books := new(models.Book)
	//validate the request body
	if err := c.BodyParser(books); err != nil {
		utils.Log("ERROR", "book", constants.Url_add_book, userId,"Newbook", err.Error(), startTime, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//use the validator library to validate required fields
	if validationErr := validate.Struct(books); validationErr != nil {
		utils.Log("ERROR", "book", constants.Url_add_book, userId,"NewBook", validationErr.Error(),startTime, time.Now())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})

	}
    //to first store the data in the redis
	data, err := json.Marshal(books)
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_books,userId, "Newbook", "Error marshaling books data:"+err.Error(), startTime, time.Now())
	} else {
		err := redisClient.Set(ctx, "data", data, 1*time.Hour).Err()
		if err != nil {
			utils.Log("ERROR", "book", constants.Url_get_books,userId, "Newbook", constants.Error_caching_data + err.Error(), startTime, time.Now())
		}
	}

	// Invalidate the books cache in Redis
	
	del := redisClient.Del(ctx, "books").Err()
	if del!= nil {
		utils.Log("ERROR", "book", constants.Url_add_book ,userId, "Newbook", constants.Error_deleting_cached_data + err.Error(), startTime, time.Now())
	}

	// Check if the data exists in the Redis cache
	
	cachedData, err := redisClient.Get(ctx, "data").Result()
	fmt.Println("cacheddata",cachedData)
	fmt.Println("err",err)
	if err == nil {
		// Data exists in the cache, retrieve and return it
		
		if err := json.Unmarshal([]byte(cachedData),books); err != nil {
			utils.Log("ERROR", "book", constants.Url_add_book,userId, "Newbook", constants.Unmarshal_error + err.Error(), startTime, time.Now())
		}

		time.Sleep(5 *time.Second)
		endTime := time.Now() 
		utils.Log("INFO", "book", constants.Url_add_book,userId, "Newbook", "ended", startTime, endTime)
		
		
		//push the data to the channel to entry in the db
		bookDataChannel <- *books
		return c.JSON(fiber.Map{"data":books})
		
	}

	bookDataChannel <- *books
	
	endTime := time.Now() 
	utils.Log("INFO", "book", constants.Url_add_book,userId, "Newbook", "ended",startTime, endTime)
	return c.JSON(fiber.Map{"data":books})
}



func Deletebook(c *fiber.Ctx) error {
	
	
	utils.InitLogger()
    userId := c.Locals("userId").(string)
	// Set the start time before processing the request
	startTime := time.Now()

	utils.Log("INFO", "book", constants.Url_get_single_book, userId,"Deletebook", "started",startTime, time.Now())
	db := database.Database
	title := c.Params("title")
	var books models.Book
	db.Where("title= ?", title).Find(&books)

	if books.Title == "" {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book,userId, "Deletebook", "ended",startTime, endTime)
		return c.Status(253).JSON("No book found ")
	}
	db.Where("title", title).Delete(&books)

	// Invalidate the employees cache in Redis
	redisClient := database.RedisClient
	err := redisClient.Del(ctx, "books",fmt.Sprintf(constants.Redis_book_const, title)).Err()
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book,userId,"Deletebook", constants.Error_deleting_cached_data + err.Error(),startTime, time.Now())
	}
	endTime := time.Now() // Get the end time after processing the request
	utils.Log("INFO", "book", constants.Url_get_single_book,userId, "Deletebook", "ended",startTime, endTime)
	return c.JSON("book is deleted succesfully.")

}

func UpdateBook(c *fiber.Ctx) error {
	
	utils.InitLogger()
	userId := c.Locals("userId").(string)
	// Set the start time before processing the request
	startTime := time.Now()

	title := c.Params("title")
	utils.Log("INFO", "book", constants.Url_get_single_book,userId, "Updatebook", "started",startTime, time.Now())
	db := database.Database
	book := new(models.Book)
	db.Where("title= ?", title).Find(&book)
	if book.ID == 0 {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book,userId, "UpdateBook", "ended", startTime, endTime)
		return c.Status(253).JSON("No book found ")
	}

	if err := c.BodyParser(book); err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId,"UpdateBook", err.Error(),startTime, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	db.Save(&book)

	// Invalidate the employees cache in Redis
	redisClient := database.RedisClient
	err := redisClient.Del(ctx, "books",fmt.Sprintf(constants.Redis_book_const, title)).Err()
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId,"Deletebook", constants.Error_deleting_cached_data + err.Error(),startTime, time.Now())
	}
    endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_update_book,userId, "Updatebook", "ended", startTime, endTime)

	return c.JSON(book)
}

func dequeueEmployeeData() {
	utils.InitLogger()
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book,"dequequedata", "", "started",startTime, time.Now())
	
	for book := range bookDataChannel {
		utils.Log("INFO", "book", constants.Url_add_book,"dequequedata","", "Data received",startTime, time.Now())
		// Calculate the response time
		startTime := time.Now()

		db  := database.Database
	    err := db.Create(&book).Error
		if err !=nil {
			utils.Log("Error","book",constants.Url_add_book,"dequeuedata","",err.Error(),startTime, time.Now())
		}

		endTime := time.Now()

		utils.Log("INFO", "book", constants.Url_add_book,"dequequedata","","ended", startTime, endTime)
	
	}
	
	
}
