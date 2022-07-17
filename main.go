package main

import (
	"vasco/controllers"
	"vasco/db"
	"sync"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	db.ConnectDatabase()


	
	r.POST("/mood", controllers.InputMood)
	r.GET("/mood", controllers.GetMoods)

	r.GET("/delete", controllers.ClearQuestionSet)
	r.GET("/questions", controllers.GetAllQuestions)
	r.GET("/updateQuestions", controllers.UpdateQuestionSet)
	// Listen and Server in 0.0.0.0:8080

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		r.Run(":8080")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		controllers.StartBot()
		wg.Done()
	}()
	wg.Wait()
}

