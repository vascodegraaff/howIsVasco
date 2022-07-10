package main

import (
	"vasco/controllers"
	"vasco/models"
	"sync"
	"github.com/gin-gonic/gin"
	"vasco/bot"
)

func main() {

	r := gin.Default()

	models.ConnectDatabase()


	
	r.POST("/mood", controllers.InputMood)
	r.GET("/mood", controllers.GetMoods)

	r.GET("/delete", controllers.ClearQuestionSet)
	r.GET("/question", controllers.GetQuestionSets)
	r.GET("/updateQuestions", controllers.AddQuestionSet)
	// Listen and Server in 0.0.0.0:8080

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		r.Run(":8080")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		bot.StartBot()
		wg.Done()
	}()
	wg.Wait()
}

