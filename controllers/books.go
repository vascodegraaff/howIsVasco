package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "vasco/models"

)

type BookInput struct {
	Title string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`

}

func CreateBook(c *gin.Context) {
    var input BookInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    book := models.Book{Title: input.Title, Author: input.Author}

    models.DB.Create(&book)

    c.JSON(http.StatusOK, gin.H{"data": book})
}
func FindBooks(c *gin.Context) {
    var books []models.Book
    models.DB.Find(&books)

    c.JSON(http.StatusOK, gin.H{"data": books})
}

