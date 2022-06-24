package models

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)
var DB *gorm.DB
func ConnectDatabase() {
    dsn := "host=localhost user=postgres password=bayesian dbname=postgres port=6776 sslmode=disable TimeZone=Europe/Amsterdam"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

		// Migrate the schema
		db.AutoMigrate(&Mood{})

		DB = db
}