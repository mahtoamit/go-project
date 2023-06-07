package database

import (
	
	"fmt"
	"time"

	
	"github.com/tutorialedge/go-fiber-tutorial/configs"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

var Database_URI string = configs.EnvDBURI("MYSQLURL")




func DBConnect() {
	utils.InitLogger()

	var err error

	Database, err = gorm.Open(mysql.Open(Database_URI), &gorm.Config{})
	if err != nil {
		fmt.Println("Not connected to database")
	}
	fmt.Println("Database connected Succesfully")

	Database.AutoMigrate(&models.Book{})
	Database.AutoMigrate(&models.SignupRequest{})

	
}

func CloseDatabase() {
	utils.InitLogger()
	startTime := time.Now()
	endTime := time.Now()

	db, err := Database.DB()
	if err != nil {

		utils.Log("ERROR", "database", "DB","", "CloseDatbase", "Error getting database connection",startTime, endTime)
	}
	db.Close()

	
}
