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