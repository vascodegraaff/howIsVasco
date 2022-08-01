package controllers

import (
	// "errors"
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"

	// "net/http"
	"strconv"
	"strings"

	"time"
	// "vasco/models"

	// "github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// func (c *ConversationHandler) scheduleJob(){

// }

// type MessageHandler struct {
// 	filters []string
// 	callback func() int
// }

// func Filter()

// func (m *MessageHandler) MessageHandler(updates tgbotapi.UpdatesChannel,b *tgbotapi.BotAPI){
// 	log.Printf(m.Bot.Token)
// }

//TODO
func Contains[T comparable](arr []T, e T) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

// Write message handler
// Write Conversation handler

// Get the previous message sent to the user by the bot
// func GetPreviousMessage()

// take current updates and compare to ongoing conversations, then add values of updates to the correct conversation
func GlobalHandler(updates tgbotapi.UpdatesChannel, b *tgbotapi.BotAPI, handlers *ConversationHandler) {

	for update := range updates {
		if update.Message != nil {

			if Contains(handlers.EntryPoints, update.Message.Command()) {
				handlers.Active = true
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, handlers.States[handlers.CurrentState].description)
				msg.ReplyToMessageID = update.Message.MessageID
				// nRow := 3
				var replyKeyboard =	tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("test"),
						tgbotapi.NewKeyboardButton("test"),
					),
				)
				msg.ReplyMarkup = replyKeyboard
				b.Send(msg)
				handlers.CurrentState = "end"
			} else if handlers.Active {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, handlers.States[handlers.CurrentState].description)
				msg.ReplyToMessageID = update.Message.MessageID
				// nRow := 3
				var replyKeyboard =	tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("test"),
						tgbotapi.NewKeyboardButton("test"),
					),
				)
				msg.ReplyMarkup = replyKeyboard
				b.Send(msg)
				handlers.CurrentState = "end"	
			}
			// for finding value in dict
			// if messageHandler, ok := handlers.EntryPoints[update.Message.Command()]; ok {
			// 	b.Send(messageHandler.description)

			// 	//do something here
			// }
			// switch update.Message.Command(){
			// case handlers.EntryPoints[0]: {
			// 	b.Send()
			// }
			// 	case "r": {
			// 		log.Printf("received r")
			// 		log.Printf("%+v\n",update.Message.Chat.ID)
			// 		log.Printf("%+v\n",update.Message.MessageID)
			// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hello world")
			// 		msg.ReplyToMessageID = update.Message.MessageID
			// 		b.Send(msg)
			// 	}
			// 	default: {
			// 		log.Printf(update.Message.Text)
			// 		log.Printf("command not found")
			// 	}
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
			case "/mood":
				{
					mood, err := strconv.Atoi(args)
					if err != nil {
						panic(err)
					} else {
						log.Printf("Mood: %v", mood)
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

type MessageHandler struct {
	// filters []string
	description     string
	// replyTypes string
	keyboardOptions []string
	callbacks       func()
}

type Handler interface {
	scheduleJob()
	checkUpdate()
	handleUpdate()
	updateState()
	triggerTimeout()
}

type ConversationHandler struct {
	Active bool
	EntryPoints []string
	CurrentState string
	States      map[string]MessageHandler
	Fallbacks   []MessageHandler

}

func getInput(input chan string) {
	for {
		var data string
		fmt.Println("input a string")
		fmt.Scan(&data)
		input <- data
	}
}

func StartBot() {
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

	conversationHandler := ConversationHandler{

		// states should contain the possible fsm transitions
		// for simplicity, start with a single state
		Active: false,
		EntryPoints: []string{
			"journal",
		},
		CurrentState: "start",
		States: map[string]MessageHandler{
			"start": {
				"Please select a journal entry type",
				[]string{"Morning Routine", "Evening Routine", "Story", "Reflection"},
				func() {
					input := make(chan string, 1)
					go getInput(input)
					for {
						fmt.Println("input something")
						select {
						case i := <-input:
							fmt.Println("result")
							fmt.Println(i)
							switch i {
							case "Morning":
								{
									fmt.Println("Lets goooo")
								}
							}

						case <-time.After(4000 * time.Millisecond):
							fmt.Println("timed out")
						}
					}
				},
			},
			"end":{
				"Thanks for the entry",
				[]string{},
				func(){},
			},
		},
	}
	GlobalHandler(updates, b, &conversationHandler)

	// Loop through each update.
	// for update := range updates {
	// 	// Check if we've gotten a message update.
	// 	if update.Message != nil {
	// 			// Construct a new message from the given chat ID and containing
	// 			// the text that we received.
	// 			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	// 			// If the message was open, add a copy of our numeric keyboard.
	// 			switch update.Message.Command() {
	// 			case "r": {
	// 				question_id,err := strconv.Atoi(strings.Split(strings.Split(update.Message.Text, " ")[1], ":")[1]);if err != nil {
	// 					log.Printf("error converting question_id to int")
	// 				}
	// 				reply_value := strings.Split(update.Message.Text,"-")[1]

	// 				var answer = models.Answer{
	// 					QuestionID: question_id,
	// 					Answer: reply_value,
	// 					DateTime: time.Now(),
	// 				}

	// 				json, _ := json.Marshal(answer)

	// 				resp, err := http.Post("http://localhost:8080/answer", "application/json", bytes.NewBuffer(json)); if err != nil {log.Printf("error posting answer")}
	// 				log.Printf(resp.Status)
	// 				log.Printf(string(json))

	// 			}
	// 			default: log.Printf("command not found")
	// 			}

	// 	} else if update.CallbackQuery != nil {
	// 			// Respond to the callback query, telling Telegram to show the user
	// 			// a message with the data received.
	// 			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	// 			if _, err := b.Request(callback); err != nil {
	// 					panic(err)
	// 			}

	// 			log.Printf("callback registered")
	// 			log.Printf("args: " + update.CallbackQuery.Message.CommandArguments())
	// 			log.Printf("test: " + update.CallbackData())

	// 			command := strings.Split(update.CallbackData(), " ")[0]
	// 			args := strings.Join(strings.Split(update.CallbackData(), " ")[1:], " ")
	// 			switch command {
	// 				case "/mood": {
	// 					mood, err := strconv.Atoi(args)
	// 					if err != nil {
	// 						panic(err)
	// 					} else {
	// 						log.Printf("Mood: %v",mood)
	// 					}

	// 				}
	// 			}

	// 			// And finally, send a message containing the data received.
	// 			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	// 			if _, err := b.Send(msg); err != nil {
	// 					panic(err)
	// 			}
	// 	}
	// }
}
