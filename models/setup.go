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
		// db.DropTableIfExists(&QuestionSet{}, &Question{})
		//db.Migrator().DropTable(&QuestionSet{}, &Question{})
		
		db.AutoMigrate(&Mood{}, &QuestionSet{}, &Question{}, &Answer{})
		// db.AutoMigrate(&Mood{}, &Customer{}, &Contact{})
		// db.Model(&Question{}).AddForeignKey("set_id", "question(set_id)", "CASCADE", "CASCADE")

		//db.Model(&Contact{}).AddForeignKey("cust_id", "customers(cust_id)", "CASCADE", "CASCADE")
		DB = db
}