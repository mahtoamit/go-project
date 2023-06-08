package queries

import (
	"context"
	"encoding/json"
	"time"

	
	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)
var ctx = context.Background()



func RedisCache(key string,userId string)(books []models.Book ){
	// Check if the data exists in the Redis cache
	var redisClient = database.RedisClient
	cachedData, err := redisClient.Get(ctx, key).Result()
	
	if err == nil {
		// Data exists in the cache, retrieve and return it
		
		if err := json.Unmarshal([]byte(cachedData), &books); err != nil {
			utils.Log("ERROR", "book", constants.Url_get_books,userId, "Getbooks", constants.Unmarshal_error +err.Error())
		}
		
		
	}
	return books
}

func RedisSetCache(books []models.Book,key string ,userId string)(flag bool,err error){
	// Store the data in the Redis cache for future use
	data, err := json.Marshal(books)
	var redisClient = database.RedisClient
	flag = false
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_books,userId, "Getbooks", constants.Unmarshal_error + err.Error())
		return flag ,err
	} else {
		err := redisClient.Set(ctx, key , data, 1*time.Hour).Err()
		if err != nil {
			utils.Log("ERROR", "book", constants.Url_get_books, userId,"Getbooks", constants.Error_caching_data + err.Error())
			return flag ,err
		}
	}
	flag =true
	return flag, err

}

func RedisSetCacheBook(books models.Book,key string ,userId string)(flag bool,err error){
	// Store the data in the Redis cache for future use
	data, err := json.Marshal(books)
	var redisClient = database.RedisClient
	flag = false
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_single_book,userId, "Getbooks", constants.Unmarshal_error + err.Error())
		return flag ,err
	} else {
		err := redisClient.Set(ctx, key , data, 1*time.Hour).Err()
		if err != nil {
			utils.Log("ERROR", "book", constants.Url_get_single_book, userId,"Getbooks", constants.Error_caching_data + err.Error())
			return flag ,err
		}
	}
	flag =true
	return flag ,err

}


func RedisCacheGetBook(userId string,key string)(books[]models.Book ){
	redisClient := database.RedisClient
	cachedData, err := redisClient.Get(ctx, key ).Result()
	if err == nil {
		// Data exists in the cache, retrieve and return it
		var books []models.Book
		if err := json.Unmarshal([]byte(cachedData), &books); err != nil {
			utils.Log("ERROR", "book" ,constants.Url_get_single_book,userId, "Getbook", constants.Unmarshal_error + err.Error())
		}	
	}
	return books
}