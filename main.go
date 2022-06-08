package main

import (
	"log"
	"vasco/controllers"
	"vasco/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv(){
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatalf("Error loading .env file")
  }
}

func main() {

	loadEnv()
	r := gin.Default()

	models.ConnectDatabase()

	controllers.StartBot()


	r.GET("/books", controllers.FindBooks)
	r.POST("/books", controllers.CreateBook)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
