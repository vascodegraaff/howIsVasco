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
// if the set id already exists, we update the values, otherwise we create a new set
func UpdateQuestionSet(c *gin.Context) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/vasco/question.json")
	if err != nil {
		panic("unable to read file")
	}
	
	var question_set []models.QuestionSet
	_ = json.Unmarshal([]byte(file), &question_set)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for _, set := range question_set {
		json_set, _ := json.Marshal(set)
		log.Printf("%s", json_set)
		for _, question := range set.Questions {
			json_question, _ := json.Marshal(question)
			log.Printf("%s", json_question)
			// if(models.DB.Find(&models.Question{}, "id = ?", question.ID)!=nil) {
				// models.DB.Model(&question).Where("id = ?", question.ID).Updates(question)
			// } else {
			models.DB.Create(&question)
			// }
		}
		log.Printf("%s", json_set)
		// models.DB.Create(&set)
	}

	c.JSON(http.StatusOK, gin.H{"question_set": question_set})
}

func ClearQuestionSet(c *gin.Context) {
	var question_sets []models.QuestionSet
	models.DB.Find(&question_sets)
	for _, question_set := range question_sets {
		models.DB.Delete(&question_set)
	}
	var questions []models.Question
	for _, question := range questions {
		models.DB.Delete(&question)
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

func GetAllQuestions(c *gin.Context) {
	var questions []models.Question
	models.DB.Find(&questions)
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func AnswerQuestion() {

}

