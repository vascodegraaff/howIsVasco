package bot

import (
	// "errors"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vasco/models"

	// "github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
func StartBot(){

	b, err := tgbotapi.NewBotAPI("5366512490:AAFNjdosYKeQofgp4BdI0ehUissp7-sIGRM")
	b.Debug = true
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", b.Self.UserName)

	SetJobs(b)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
				// Construct a new message from the given chat ID and containing
				// the text that we received.
				// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

				// If the message was open, add a copy of our numeric keyboard.
				switch update.Message.Command() {
				case "r": {
					question_id,err := strconv.Atoi(strings.Split(strings.Split(update.Message.Text, " ")[1], ":")[1]);if err != nil {
						log.Printf("error converting question_id to int")
					}
					reply_value := strings.Split(update.Message.Text,"-")[1]

					var answer = models.Answer{
						QuestionID: question_id,
						Answer: reply_value,
						DateTime: time.Now(),
					}

					json, _ := json.Marshal(answer)

					resp, err := http.Post("http://localhost:8080/answer", "application/json", bytes.NewBuffer(json)); if err != nil {log.Printf("error posting answer")}
					log.Printf(resp.Status)
					log.Printf(string(json))

				}
				default: log.Printf("command not found")
				}

		} else if update.CallbackQuery != nil {
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := b.Request(callback); err != nil {
						panic(err)
				}

				log.Printf("callback registered")
				log.Printf("args: " + update.CallbackQuery.Message.CommandArguments())
				log.Printf("test: " + update.CallbackData())

				command := strings.Split(update.CallbackData(), " ")[0]
				args := strings.Join(strings.Split(update.CallbackData(), " ")[1:], " ")
				switch command {
					case "/mood": {
						mood, err := strconv.Atoi(args)
						if err != nil {
							panic(err)
						} else {
							log.Printf("Mood: %v",mood)
						}
						
					}
				}

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				if _, err := b.Send(msg); err != nil {
						panic(err)
				}
		}
	}
}

