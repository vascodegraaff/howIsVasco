package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"vasco/db"
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
	
	var questions []models.Question
	_ = json.Unmarshal([]byte(file), &questions)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for _, set := range questions {
		json_set, _ := json.Marshal(set)
		log.Printf("%s", json_set)
		// for _, question := range set.Questions {
		// 	json_question, _ := json.Marshal(question)
		// 	log.Printf("%s", json_question)
		// 	// if(models.DB.Find(&models.Question{}, "id = ?", question.ID)!=nil) {
		// 		// models.DB.Model(&question).Where("id = ?", question.ID).Updates(question)
		// 	// } else {
		// 	models.DB.Create(&question)
		// 	// }
		// }
		log.Printf("%s", json_set)
		// models.DB.Create(&set)
	}

	c.JSON(http.StatusOK, gin.H{"question_set": questions})
}

func ClearQuestionSet(c *gin.Context) {
	var questions []models.Question
	for _, question := range questions {
		db.DB.Delete(&question)
	}
	// models.DB.Where("1 = 1").Delete(&models.Question{})
	// models.DB.Where("1 = 1").Delete(&models.QuestionSet{})
	c.Data(http.StatusAccepted, "application/json", []byte("{\"message\": \"question_set deleted\"}"))
}


func GetAllQuestions(c *gin.Context) {
	var questions []models.Question
	db.DB.Find(&questions)
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func AnswerQuestion() {

}

