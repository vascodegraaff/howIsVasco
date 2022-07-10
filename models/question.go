package models

import (
	"time"

	"gorm.io/gorm"
)

type ReplyType int32
const (
	RANGE ReplyType = 0
	NUMBER ReplyType = 1
	TEXT ReplyType = 2
	YES_NO ReplyType = 3
)

type ScheduleType int32
const (
	CRON ScheduleType = 0
	RANDOM ScheduleType = 1
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
	QuestionID 		int `json:"question_id"`
	Answer 				string `json:"answer"` 
	DateTime 			time.Time `json:"date_time"`
}