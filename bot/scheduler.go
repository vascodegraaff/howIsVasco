package bot

import (
	"encoding/json"
	"fmt"
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

var rangeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
	),
)

// Send a message with the question and set a reply keyboard
func SendMessage(bot *tgbotapi.BotAPI, question *models.Question) {
	message := tgbotapi.NewMessage(5383565084, question.Question)
	switch question.ReplyType {
	case models.RANGE:
		var replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.ID)+" - 1"),
				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.ID)+" - 2"),
				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.ID)+" - 3"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.ID)+" - 4"),
				tgbotapi.NewKeyboardButton("/r q_id:"+fmt.Sprintf("%v", question.ID)+" - 5"),
			),
		)
		message.ReplyMarkup = replyKeyboard
	}

	bot.Send(message)
	log.Printf("Message sent: " + question.Question)

}

func SetJobs(bot *tgbotapi.BotAPI) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/bot/question.json")
	if err != nil {
		panic("unable to read file")
	}

	questionSet := make([]models.QuestionSet,0)
	_ = json.Unmarshal([]byte(file), &questionSet)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	c := cron.New()

	for i, set := range questionSet {
		// questionSet[i].Id = uuid.New()
		log.Printf("set_name: %s\n", set.QuestionSetName)
		log.Printf("description: %s\n", set.Description)
		log.Printf("schedule type: %v\n", set.Schedule)
		for _, question := range set.Questions {
			switch questionSet[i].Schedule {
			case models.CRON:
				c.AddFunc(questionSet[i].ScheduleValue, func() {
					SendMessage(bot, question)
					log.Printf("cron job executed")
				})
			case models.RANDOM:


			}
			log.Printf("question id: %v\n", question.ID)
			log.Printf("question: %s\n", question.Question)
			log.Printf("reply type: %v\n", question.ReplyType)
		}
	}
	c.Start()
}
