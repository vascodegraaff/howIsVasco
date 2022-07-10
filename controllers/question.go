package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	// "fmt"
	"encoding/json"
	"vasco/models"

	// "gorm.io/gorm"
	"github.com/gin-gonic/gin"
	// "vasco/models"
)


// Go through json of questions and construct the model in the database
func AddQuestionSet(c *gin.Context) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/vasco/question.json")
	if err != nil {
		panic("unable to read file")
	}
	
	var question_set []models.QuestionSet
	_ = json.Unmarshal([]byte(file), &question_set)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for i, set := range question_set {
		log.Printf("%v", set.SetID)
		log.Printf("%v", set.ID)
		log.Printf("%v", set.Description)
		log.Printf("%v", set.Schedule)
		log.Printf("%v", set.ScheduleValue)
		for _, question := range set.Questions {
			// log.Printf("%v", question.QuestionSetID)
			// values := map[string]int{"question":question}

			json, _ := json.Marshal(question)
			log.Printf("question: %v", string(json))
			log.Printf("%v", question.QuestionID)
			log.Printf("%v", question.Question)
			log.Printf("%v", question.ReplyType)
			// models.DB.Create(&question)
		}
		models.DB.Create(&question_set[i])
	}

	c.JSON(http.StatusOK, gin.H{"question_set": question_set})
}

func ClearQuestionSet(c *gin.Context) {
	var question_sets []models.QuestionSet
	models.DB.Find(&question_sets)
	for _, question_set := range question_sets {
		models.DB.Delete(&question_set)
	}
	// models.DB.Where("1 = 1").Delete(&models.Question{})
	// models.DB.Where("1 = 1").Delete(&models.QuestionSet{})
	c.Data(http.StatusAccepted, "application/json", []byte("{\"message\": \"question_set deleted\"}"))
}

func GetQuestionSets(c *gin.Context) {
	var questionSets []models.QuestionSet
	models.DB.Find(&questionSets)
	c.JSON(http.StatusOK, gin.H{"questionSets": questionSets})

}

func AnswerQuestion() {

}

