package database

import (
	
	"fmt"
	// "time"
	"strconv"
	"os"

	
	"github.com/tutorialedge/go-fiber-tutorial/configs"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	// "github.com/tutorialedge/go-fiber-tutorial/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var Database *gorm.DB

// var Database_URI string = configs.EnvDBURI("MYSQLURL")






// Database instance
type Dbinstance struct {
	Db *gorm.DB
   }
var DB Dbinstance
   // Connect function
func Connect() {
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
	 os.Exit(2)
	}
	db.AutoMigrate(&models.Book{},&models.SignupRequest{})
	DB = Dbinstance{
	 Db: db,
	}
}

  


	

