package models

import "time"

type Journal struct {
	ID uint `json:"id" gorm:"primary_key"`
	Type string `json:"title"`
	DateTime time.Time `json:"dateTime"` 
	Text string `json:"text"`
}

