package controllers

import (
	"log"
	"net/http"
	"time"
	"vasco/db"
	"vasco/models"
	"github.com/gin-gonic/gin"
)

type JournalInput struct {
	Type string `json:"mood" binding:"required"`
	Text string `json:"text" binding:"required"`
}

func EnterJournalEntry(Type string, Text string) {
	journal := models.Journal{Type: Type, Text: Text, DateTime: time.Now(),}
	db.DB.Create(&journal)
	log.Printf("Journal entry created: %v", journal)
}
			
func GetJournalEntries(c *gin.Context) {
	var journals []models.Journal
	db.DB.Find(&journals)
	c.JSON(http.StatusOK, gin.H{"data": journals})
}

func GetJournalEntriesByType(c *gin.Context, Type string) {
	var journals []models.Journal
	db.DB.Find(&journals, "type = ?", Type)
	c.JSON(http.StatusOK, gin.H{"data": journals})
}

func DeleteAllJournals(c *gin.Context) {
		var journals []models.Journal
		db.DB.Find(&journals)
    db.DB.Delete(journals)
    c.JSON(http.StatusOK,gin.H{"data": "deleted"})
}

