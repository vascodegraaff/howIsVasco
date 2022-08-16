package controllers

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	// "strconv"
	"vasco/models"
	// "github.com/google/uuid"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

var moodKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("😡", "/mood 1"),
		tgbotapi.NewInlineKeyboardButtonData("😐", "/mood 2"),
		tgbotapi.NewInlineKeyboardButtonData("🙂", "/mood 3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("😃", "/mood 4"),
		tgbotapi.NewInlineKeyboardButtonData("🤩", "/mood 5"),
	),
)

// Send a message with the question and set a reply keyboard
// func SendMessage(bot *tgbotapi.BotAPI, question *models.Question) {
// 	message := tgbotapi.NewMessage(5383565084, question.Question)
// 	switch question.ReplyType {
// 	case models.RANGE:
// 		var replyKeyboard = tgbotapi.NewReplyKeyboard(
// 			tgbotapi.NewKeyboardButtonRow(
// 				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.QuestionID)+" - 1"),
// 				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.QuestionID)+" - 2"),
// 				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.QuestionID)+" - 3"),
// 			),
// 			tgbotapi.NewKeyboardButtonRow(
// 				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.QuestionID)+" - 4"),
// 				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.QuestionID)+" - 5"),
// 			),
// 		)
// 		message.ReplyMarkup = replyKeyboard
// 	}

// 	bot.Send(message)
// 	log.Printf("Message sent: " + question.Question)

// }
func SendMessage(bot *tgbotapi.BotAPI, question string) {
	message := tgbotapi.NewMessage(5383565084, question)
	message.ReplyMarkup = moodKeyboard
	bot.Send(message)
	log.Printf("Message sent: " + question)
}

// func HandleReply(bot *tgbotapi.BotAPI)

func SetJobs(bot *tgbotapi.BotAPI) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/vasco/question.json")
	if err != nil {
		panic("unable to read file")
	}

	Questions := make([]models.Question,0)
	
	questionSet := models.QuestionSet{
		Schedule: "cron",
		ScheduleValue: "0 0 0 * * *",
		Questions: []string{
			"test",
			"bruh",
		},
	}
	_ = json.Unmarshal([]byte(file), &Questions)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// x, err := json.MarshalIndent(Questions, "", "  ")
	if err != nil {
		log.Fatal("Error during MarshalIndent(): ", err)
	}
	// log.Printf("%s", x)
	c := cron.New()

	if questionSet.Schedule == models.CRON {
		c.AddFunc(questionSet.ScheduleValue, func() {
			for _, question := range questionSet.Questions {
				SendMessage(bot, question)
			}
		})
	}
	
	c.Start()
}
