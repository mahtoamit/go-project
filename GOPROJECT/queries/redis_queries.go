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
	cachedData, err := redisClient.Get(ctx,key).Result()
	if err == nil {
		// Data exists in the cache, retrieve and return it
		var books []models.Book
		if err := json.Unmarshal([]byte(cachedData), &books); err != nil {
			utils.Log("ERROR", "book" ,constants.Url_get_single_book,userId, "Getbook", constants.Unmarshal_error + err.Error())
		}	
	}
	return books
}


func RedisSetNewCache(books models.Book,key string)(flag1 bool,err error){
	// Store the data in the Redis cache for future use
	data, err := json.Marshal(books)
	print("data", data)
	var redisClient = database.RedisClient
	flag1 = false
	if err != nil {
		utils.Log("ERROR", "book", constants.Url_get_books ,"","Getbooks", constants.Unmarshal_error + err.Error())
		return flag1 ,err
	} else {
		err := redisClient.LPush(ctx,key , data, 1*time.Hour).Err()
		if err != nil {
			utils.Log("ERROR", "book", constants.Url_get_books,"","Getbooks", constants.Error_caching_data + err.Error())
			return flag1 ,err
		}
	}
	flag1 =true
	return flag1, err

}

func RedisDeleteBook(key string )(flag bool){
	var redisClient = database.RedisClient
	flag = false
	del := redisClient.Del(ctx, key ).Err()
	if del!= nil {
		return false  

	}
    return true
}


func RedisNewCache(key string,userId string)(books models.Book ){
	// Check if the data exists in the Redis cache
	var redisClient = database.RedisClient
	cachedData, err := redisClient.Get(ctx, key).Result()
	
	if err == nil {
		// Data exists in the cache, retrieve and return it
		
		if err := json.Unmarshal([]byte(cachedData), &books); err != nil {
			utils.Log("ERROR", "book", constants.Url_add_book,userId, "NewBook", constants.Unmarshal_error +err.Error())
		}
		
		
	}
	return books
}


func Deletebook(data string,userId string)(err error){
	// Invalidate the employees cache in Redis
	redisClient := database.RedisClient
	err = redisClient.Del(ctx, constants.Books_constant,data).Err()
	if err != nil {
		return err
	}
	return nil
}


func RetrieveData() (string, error){
	var redisClient = database.RedisClient
	// Retrieve the list values using LRANGE
	cachedData, err := redisClient.LPop(ctx, constants.Books_data).Result()

    return cachedData , err

}

func SetAuthenticationCache(key string,data string)( bool,error){
	var redisClient = database.RedisClient
	err := redisClient.Set(ctx, key , data, 1*time.Hour).Err()
		if err != nil {
			return false ,err
		}
		return true,err

}

func GetAuthenticationCache(key string)(string){
	var redisClient = database.RedisClient
	cachedData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		cachedData = ""
		return  cachedData
	}


	return cachedData
}

func RedisDeleteAuthentication(token string )(flag bool){
	var redisClient = database.RedisClient
	flag = false
	del := redisClient.Del(ctx, token).Err()
	if del!= nil {
		return false  

	}
    return true
}
