package health

import (
	"context"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/configs"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB
var RedisClient *redis.Client
var Database_URI string = configs.EnvDBURI("MYSQLURL")
var RedisURL string = configs.EnvDBURI("REDISURL")

var ctx = context.Background()

func Connect(c * fiber.Ctx) error {
	utils.InitLogger()
	startTime := time.Now()

	var err error

	Database, err = gorm.Open(mysql.Open(Database_URI), &gorm.Config{})
	if err != nil {
		utils.Log("ERROR", "database", "Database","", "Connect", "Failed to connect to the Database server",startTime, time.Now())
		return c.Status(500).JSON(fiber.Map{"error":"Internal server error for database"})
	}
	fmt.Println("Database connected Succesfully")

	Database.AutoMigrate(&models.Book{})
	Database.AutoMigrate(&models.SignupRequest{})

	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisURL,
		Password: "", // Add Redis password if required
		DB:       0,  // Use default Redis database
	})

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		utils.Log("ERROR", "database", "redis","", "Connect", "Failed to connect to the Redis server",startTime, time.Now())
		return c.Status(500).JSON(fiber.Map{"error":"Internal server error for redis"})
	}
	utils.Log("INFO", "database", "redis","", "Connect", "Connection opened to the Redis server",startTime, time.Now())
	utils.Log("INFO", "database", "Database", "","Connect", "Database connected Succesfully",startTime, time.Now())
	return c.SendStatus(200)
}

