package health

import (
	"context"
	"fmt"
	"time"
	"strconv"
	
	"gorm.io/driver/postgres"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/configs"
	
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var Database *gorm.DB
var RedisClient *redis.Client
// var Database_URI string = configs.EnvDBURI("MYSQLURL")
var RedisURL string = configs.EnvDBURI("REDISURL")
type Dbinstance struct {
	Db *gorm.DB
   }
var DB Dbinstance

var ctx = context.Background()

func Connect(c * fiber.Ctx) error {
	utils.InitLogger()
	startTime := time.Now()

	var err error

	p := configs.EnvDBURI("DB_PORT")
	// because our config function returns a string, we are parsing our      str to int here
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
	 fmt.Println("Error parsing str to int")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", configs.EnvDBURI("DB_HOST"), configs.EnvDBURI("DB_USER"), configs.EnvDBURI("DB_PASSWORD"),  configs.EnvDBURI("DB_NAME"), port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	})
	if err != nil {
		utils.Log("ERROR", "database", "Database","", "Connect", "Failed to connect to the Database server",startTime, time.Now())
		return c.Status(500).JSON(fiber.Map{"error":"Internal server error for database"})
	 
	}
	fmt.Println("Database connected Succesfully")
	DB = Dbinstance{
	 Db: db,
	}
		
	
	




	

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

