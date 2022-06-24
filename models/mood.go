package models

import ("time")

type Mood struct {
	ID uint `json:"id" gorm:"primary_key"`
	Mood int `json:"mood"`
	DateTime time.Time `json:"time"`
}
