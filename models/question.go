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

// type QuestionSet struct {
//   gorm.Model
// 	SetID uint `json:"set_id" gorm:"primaryKey"`
// 	// QuestionSetName string `json:"set_name" gorm:"unique"`
// 	Description string
// }

type Question struct {
  gorm.Model
	QuestionID uint `json:"question_id"` 
  Question string 
	Schedule ScheduleType
	ScheduleValue string `json:"schedule_value"`
	ReplyType ReplyType
}

type Answer struct {
	gorm.Model
	QuestionID 		int `json:"question_id"`
	Answer 				string `json:"answer"` 
	DateTime 			time.Time `json:"date_time"`
}