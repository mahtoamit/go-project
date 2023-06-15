package queries


import (
	
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
)



func Login(Email string,PassWord string)(user models.SignupRequest){
	db := database.DB.Db
	db.Where("email = ? AND password = ?", Email, PassWord).First(&user)
	return user
}


func CreateUser(user models.SignupRequest)(err error){
	db :=  database.DB.Db
	err = db.Create(&user).Error
	return err
}