package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "vasco/models"
    "time"

)

type MoodInput struct {
	Mood int `json:"mood" binding:"required"`
}

func InputMood(c *gin.Context) {
    var input MoodInput 
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    mood := models.Mood{Mood: input.Mood, DateTime: time.Now()}

    models.DB.Create(&mood)

    c.JSON(http.StatusOK, gin.H{"data": mood})
}
func GetMoods(c *gin.Context) {
    var moods []models.Mood
    models.DB.Find(&moods)

    c.JSON(http.StatusOK, gin.H{"data": moods})
}

func DeleteMoods(c *gin.Context) {
    var moods []models.Mood
    models.DB.Find(&moods)
    models.DB.Delete(moods)

    c.JSON(http.StatusOK,gin.H{"data": "deleted"})
    
}

