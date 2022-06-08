package models

import "time"

type Journal struct {
	ID uint `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	dateTime time.Time `json:"dateTime"` 
	Text string `json:"text"`
}