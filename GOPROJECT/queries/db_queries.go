package queries

import (
	"time"
	
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
)

func DBGetBooks()(books []models.Book){
    db :=  database.DB.Db
	db.Find(&books)
    
	return books
}

func DBGetSingleBook(name string)(books models.Book){
	db :=  database.DB.Db
	db.Where("title", name).Find(&books)
	return books
}


func DBDeletetBook(title string,books models.Book)(bool){
	db :=  database.DB.Db
	
	db.Where("title", title).Delete(&books)
	return true
}


func DBGetUpdateBook(title string)(books *models.Book){
	db :=  database.DB.Db
	db.Where("title", title).Find(&books)
	return books
}

func DBUpdateBook(title string,books *models.Book)(*models.Book){
	db :=  database.DB.Db
	db.Exec("UPDATE books SET title = ?,author = ?,deleted_at =?,rating = ? Where title = ?",books.Title,books.Author,time.Now().Format("2006-01-02T15:04:05.000"),books.Rating,title)
	return books
}


func DBCreate(book models.Book)(err error){
	db  :=  database.DB.Db
	err = db.Create(&book).Error
	return err 
}