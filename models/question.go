package models

import (
	"time"
	"gorm.io/gorm"
)

type ReplyType string
const (
	RANGE ReplyType = "range"
	NUMBER ReplyType = "number"
	TEXT ReplyType = "text"
	YES_NO ReplyType = "yes_no"
)


type ScheduleType string
var (
	CRON ScheduleType = "cron"
	RANDOM ScheduleType = "random"
)

type QuestionSet struct {
  gorm.Model
	SetID uint `json:"set_id" gorm:"primaryKey"`
	QuestionSetName string `json:"set_name"`
	Description string
	Schedule ScheduleType
	ScheduleValue string
  Questions []*Question `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;unique"`
}

type Question struct {
  gorm.Model
  QuestionSetID uint
	QuestionID uint `json:"question_id" gorm:"primaryKey"`
  Question string
	ReplyType ReplyType
}

type Answer struct {
	gorm.Model
	QuestionID 		int `json:"question_id"`
	Answer 				string `json:"answer"` 
	DateTime 			time.Time `json:"date_time"`
}