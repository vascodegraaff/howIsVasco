package controllers

import (
	"io/ioutil"
	"log"
	"net/http"

	// "fmt"
	"encoding/json"
	"vasco/models"

	"github.com/google/uuid"
	// "gorm.io/gorm"

	"github.com/gin-gonic/gin"
	// "vasco/models"
)


// Go through json of questions and construct the model in the database
func AddQuestionSet(c *gin.Context) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/vasco/question_test.json")
	if err != nil {
		panic("unable to read file")
	}

	var questionSet []models.QuestionSet
	_ = json.Unmarshal([]byte(file), &questionSet)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for i, set := range questionSet {
		questionSet[i].Id = uuid.New()
		log.Printf("set_name: %s\n", set.Set_name)
		log.Printf("description: %s\n", set.Description)
		log.Printf("schedule type: %s\n", set.Schedule.T)
		for j, question := range set.Questions {
			questionSet[i].Questions[j].Id = uuid.New()
			log.Printf("question id: %s\n", questionSet[i].Questions[j].Id)
			log.Printf("question: %s\n", question.Question)
			log.Printf("reply type: %s\n", question.ReplyType)
		}
		
		models.DB.Create(&questionSet[i])
	}
	// models.DB.Create(&questionSet)

	c.JSON(http.StatusOK, gin.H{"question_set": questionSet})
}

func GetQuestionSets(c *gin.Context) {
	var questionSets []models.QuestionSet
	models.DB.Find(&questionSets)
	c.JSON(http.StatusOK, gin.H{"questionSets": questionSets})

}

func AnswerQuestion() {

}

