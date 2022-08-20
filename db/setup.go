package db

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vasco/models"
)
var DB *gorm.DB
func ConnectDatabase() {
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatalf("Error loading .env file")
  }
	host := os.Getenv("HOST")
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	port := os.Getenv("DBPORT")

	dsn := "host="+ host+ " user="+user+" password="+password+" dbname="+dbname+" port="+port+" sslmode=disable TimeZone=Europe/Amsterdam"
	print(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
			panic("failed to connect database")
	}

	// Migrate the schema
	// db.DropTableIfExists(&QuestionSet{}, &Question{})
	
	db.AutoMigrate(&models.Mood{}, &models.Question{}, &models.Answer{}, &models.Journal{})
	// db.AutoMigrate(&Mood{}, &Customer{}, &Contact{})
	// db.Model(&Question{}).AddForeignKey("set_id", "question(set_id)", "CASCADE", "CASCADE")

	//db.Model(&Contact{}).AddForeignKey("cust_id", "customers(cust_id)", "CASCADE", "CASCADE")
	DB = db
}