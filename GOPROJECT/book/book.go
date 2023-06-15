package book

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
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

var bookDataChannel chan models.Book

func init() {
	bookDataChannel = make(chan models.Book, 100)
	go dequeueEmployeeData()
}

var startTime time.Time

func Getbooks(c *fiber.Ctx) error {
	utils.InitLogger()
	userId := c.Locals("userId").(string)
	startTime := time.Now() // Set the start time before processing the request
	utils.Log("INFO", "book", constants.Url_get_books, userId, "Getbooks", "started", startTime, time.Now())
	// Check if the data exists in the Redis cache

	key := "books"
	data := queries.RedisCache(key, userId)

	if data != nil {
		endTime := time.Now()
		utils.Log("INFO", "book", constants.Url_get_books, userId, "Getbooks", "ended", startTime, endTime)
		return c.Status(200).JSON(fiber.Map{"data": data})
	}

	// Data does not exist in the cache, fetch from the database
	books := queries.DBGetBooks()

	result, err := queries.RedisSetCache(books, key, userId)

	if err != nil && !result {
		utils.Log("DEBUG", "book", constants.Url_get_books, userId, "Getbooks", "ended")
		return c.Status(252).JSON(fiber.Map{"msg": err})
	}
	endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_get_books, userId, "Getbooks", "ended", startTime, endTime)
	return c.Status(200).JSON(fiber.Map{"data": &books})

}

func Getbook(c *fiber.Ctx) error {
	utils.InitLogger()
	startTime = time.Now()
	userId := c.Locals("userId").(string)
	utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Getbook", "started", startTime, time.Now())
	name := c.Params("title")
	// Check if the data exists in the Redis cache
	titlekey := fmt.Sprintf(constants.Redis_book_const, name)
	data := queries.RedisCacheGetBook(titlekey, userId)

	if data != nil {

		endTime := time.Now()
		utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Getbooks", "ended", startTime, endTime)
		return c.Status(200).JSON(fiber.Map{"data": data})
	}

	// Data does not exist in the cache, fetch from the database
	books := queries.DBGetSingleBook(name)

	if books.Title == "" {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId, "Getbook", "ended", startTime, endTime)
		return c.Status(253).JSON("No book found for  this title")
	}

	result, err := queries.RedisSetCacheBook(books, titlekey, userId)

	if err != nil && !result {
		utils.Log("DEBUG", "book", constants.Url_get_single_book, userId, "Getbook", "ended")
		return c.Status(252).JSON(fiber.Map{"msg": err})
	}

	endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Getbook", "ended", startTime, endTime)
	return c.Status(200).JSON(fiber.Map{"data": &books})

}

func Newbook(c *fiber.Ctx) error {

	utils.InitLogger()
	userId := c.Locals("userId").(string)

	// Set the start time before processing the request
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, userId, "Newbook", "started", startTime, time.Now())
	var validate = validator.New()

	books := new(models.Book)
	//validate the request body
	if err := c.BodyParser(books); err != nil {
		utils.Log("ERROR", "book", constants.Url_add_book, userId, "Newbook", err.Error(), startTime, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//use the validator library to validate required fields
	if validationErr := validate.Struct(books); validationErr != nil {
		utils.Log("ERROR", "book", constants.Url_add_book, userId, "NewBook", validationErr.Error(), startTime, time.Now())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})

	}
	//to first store the data in the redis
	key := "data"
	data, err := queries.RedisSetNewCache(books, key, userId)

	if err != nil && !data {
		utils.Log("DEBUG", "book", constants.Url_add_book, userId, "NewBook", "ended")
		return c.Status(252).JSON(fiber.Map{"msg": err})
	}

	// Invalidate the books cache in Redis
	deletekey := "books"
	delete := queries.RedisDeleteBook(deletekey, userId)

	if !delete {
		utils.Log("DEBUG", "book", constants.Url_add_book, userId, "NewBook", "ended")
		return c.Status(253).JSON(fiber.Map{"msg": "cache not deleted"})
	}

	// Check if the data exists in the Redis cache

	book := queries.RedisNewCache(key, userId)

	if book.Title != "" {
		endTime := time.Now()
		utils.Log("INFO", "book", constants.Url_add_book, userId, "Newbook", "ended", startTime, endTime)
		//push the data to the channel to entry in the db
		bookDataChannel <- book
		return c.JSON(fiber.Map{"data": fiber.Map{"title": book.Title, "author": book.Author, "rating": book.Rating}})

	}

	endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, userId, "Newbook", "ended", startTime, endTime)
	return c.JSON(fiber.Map{"data": books})
}

func Deletebook(c *fiber.Ctx) error {

	utils.InitLogger()
	userId := c.Locals("userId").(string)
	// Set the start time before processing the request
	startTime := time.Now()

	utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Deletebook", "started", startTime, time.Now())

	title := c.Params("title")

	books := queries.DBGetSingleBook(title)

	if books.Title == "" {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId, "Deletebook", "ended", startTime, endTime)
		return c.Status(253).JSON("No book found ")
	}

	queries.DBDeletetBook(title, books)

	data := fmt.Sprintf(constants.Redis_book_const, title)

	err := queries.Deletebook(data, userId)

	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId, "Deletebook", constants.Error_deleting_cached_data+err.Error())
		return c.Status(252).JSON(fiber.Map{"msg": err})
	}

	endTime := time.Now() // Get the end time after processing the request
	utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Deletebook", "ended", startTime, endTime)
	return c.Status(200).JSON(fiber.Map{"msg": "book is deleted succesfully."})

}

func UpdateBook(c *fiber.Ctx) error {

	utils.InitLogger()
	userId := c.Locals("userId").(string)
	// Set the start time before processing the request
	startTime := time.Now()

	title := c.Params("title")
	utils.Log("INFO", "book", constants.Url_get_single_book, userId, "Updatebook", "started", startTime, time.Now())

	book := new(models.Book)
	books := queries.DBGetUpdateBook(title)
	if books.ID == 0 {
		endTime := time.Now()
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId, "UpdateBook", "ended", startTime, endTime)
		return c.Status(253).JSON("No book found ")
	}

	if err := c.BodyParser(book); err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book, userId, "UpdateBook", err.Error(), startTime, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	result := queries.DBUpdateBook(title, book)

	data := fmt.Sprintf(constants.Redis_book_const, title)

	err := queries.Deletebook(data, userId)

	if err != nil {
		utils.Log("ERROR", "book", constants.Url_update_book, userId, "Updatebook", constants.Error_deleting_cached_data+err.Error())
		return c.Status(252).JSON(fiber.Map{"msg": err})
	}

	endTime := time.Now()
	utils.Log("INFO", "book", constants.Url_update_book, userId, "Updatebook", "ended", startTime, endTime)

	return c.Status(200).JSON(fiber.Map{"data": result})
}

func dequeueEmployeeData() {
	utils.InitLogger()
	startTime := time.Now()
	utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "started", startTime, time.Now())

	for book := range bookDataChannel {
		utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "Data received", startTime, time.Now())
		// Calculate the response time
		startTime := time.Now()

		err := queries.DBCreate(book)
		if err != nil {
			utils.Log("Error", "book", constants.Url_add_book, "dequeuedata", "", err.Error(), startTime, time.Now())
		}

		endTime := time.Now()

		utils.Log("INFO", "book", constants.Url_add_book, "dequequedata", "", "ended", startTime, endTime)

	}

}
