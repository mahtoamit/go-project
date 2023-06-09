package queries

import (
	
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
)

func DBGetBooks()(books []models.Book){
    db := database.Database
	db.Find(&books)
    
	return books
}

func DBGetSingleBook(name string)(books models.Book){
	db := database.Database
	db.Where("title", name).Find(&books)
	return books
}


func DBDeletetBook(title string,books models.Book)(bool){
	db := database.Database
	
	db.Where("title", title).Delete(&books)
	return true
}


func DBGetUpdateBook(title string)(books *models.Book){
	db := database.Database
	db.Where("title", title).Find(&books)
	return books
}

func DBUpdateBook(book *models.Book)(data *models.Book){
	db := database.Database
	db.Save(&book)
	return data
}


func DBCreate(book models.Book)(err error){
	db  := database.Database
	err = db.Create(&book).Error
	return err 
}