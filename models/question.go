package models

import (
	"github.com/google/uuid"
)


type QuestionSet struct {
	Id uuid.UUID `json:"set_id" gorm:"primaryKey"`
	Set_name string `json:"set_name" gorm:"set_name"`
	Description string `json:"description" gorm:"description"`
	// fixedSchedule bool 
	Schedule Schedule `json:"schedule" gorm:"foreignKey:Id"`
	Questions []QuestionItem `json:"questions" gorm:"foreignKey:Id"`
}

type Schedule struct {
	Id uuid.UUID `json:"id" gorm:"id"`
	T string `json: "t" gorm: "t"`
	Value string `json: "value" gorm:"value"`
}

type QuestionItem struct {
	Id uuid.UUID `json:"id" gorm:"primaryKey"`
	Question string `json:"question" gorm: "question"`
	ReplyType string `json:"reply_type" gorm:"reply_type"`
}