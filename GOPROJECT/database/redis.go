package database


import (
	"context"
	"github.com/tutorialedge/go-fiber-tutorial/configs"
	"github.com/go-redis/redis/v8"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"time"
)

var RedisClient *redis.Client
var RedisURL string = configs.EnvDBURI("REDISURL")

var Ctx = context.Background()

func RedisConnect(){
	utils.InitLogger()
	startTime := time.Now()
	endTime := time.Now()

	var err error

	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisURL,
		Password: "", // Add Redis password if required
		DB:       0,  // Use default Redis database
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		utils.Log("ERROR", "database", "redis","", "Connect", "Failed to connect to the Redis server",startTime, endTime)
		panic("failed to connect to the Redis server")
	}
	utils.Log("INFO", "database", "redis","", "Connect", "Connection opened to the Redis server",startTime, time.Now())
}

func ReddisClose(){
	utils.InitLogger()
	startTime := time.Now()
	endTime := time.Now()

	err := RedisClient.Close()
	if err != nil {
		utils.Log("ERROR", "database", "redis","", "CloseDatbase", "Error closing Redis client connection",startTime, endTime)
	}
	utils.Log("INFO", "database", "redis","", "CloseDatabse", "Connection closed to the Redis server",startTime, time.Now())
}