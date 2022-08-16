package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type StateType string

const (
	SELECT_ENTRY StateType = "selectEntry"
	INPUT_ENTRY  StateType = "inputEntry"
	TEXT         StateType = "text"
	NUMBER			 StateType = "number"

)

type MessageHandler struct {
	// filters []string
	description string
	// replyTypes string
	// isInput bool
	StateType       StateType
	keyboardOptions []string
	callback        func()
	nextState       string
}

type StateLogger struct {
	history []State
}

type State struct {
	MessageHandler
	field string
	input string
}

func (logger *StateLogger) addState(state State) {
	logger.history = append(logger.history, state)
	log.Printf("state logged: %v", state)
}

type Handler interface {
	scheduleJob()
	checkUpdate()
	handleUpdate()
	updateState()
	triggerTimeout()
}

type ConversationHandler struct {
	Active       bool
	EntryPoints  []string
	Cron 			 	 string
	CurrentState string
	States       map[string]MessageHandler
	Fallbacks    []MessageHandler
}

func getInput(input chan string) {
	for {
		var data string
		fmt.Println("input a string")
		fmt.Scan(&data)
		input <- data
	}
}

//TODO
func Contains[T comparable](arr []T, e T) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

func keyboardHelper(keyboardOptions []string) tgbotapi.ReplyKeyboardMarkup {
	ret := tgbotapi.NewReplyKeyboard()
	var mid int = len(keyboardOptions) / 2

	firstCol := make([]tgbotapi.KeyboardButton, mid+1)
	secondCol := make([]tgbotapi.KeyboardButton, mid+1)
	for i := 0; i < mid; i++ {
		firstCol[i].Text = keyboardOptions[i]
	}
	for i := mid; i < len(keyboardOptions); i++ {
		secondCol[i-mid].Text = keyboardOptions[i]
	}
	ret.Keyboard = append(ret.Keyboard, tgbotapi.NewKeyboardButtonRow(firstCol[:]...))
	ret.Keyboard = append(ret.Keyboard, tgbotapi.NewKeyboardButtonRow(secondCol[:]...))

	return ret
}

// Write message handler
// Write Conversation handler

// Get the previous message sent to the user by the bot
// func GetPreviousMessage()

func StartBot() {
	b, err := tgbotapi.NewBotAPI("5366512490:AAFNjdosYKeQofgp4BdI0ehUissp7-sIGRM")
	b.Debug = true
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", b.Self.UserName)

	// SetJobs(b)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	cronHandler := ConversationHandler{
		Active: false,
		Cron: "13 16 * * *",
		CurrentState: "sleep",
		States: map[string]MessageHandler{
			"sleep": {
				"How many hours did you sleep?",
				NUMBER,
				[]string{},
				func() {},
				"mood",
			},
			"mood": {
				"How are you feeling?",
				NUMBER,
				[]string{},
				func() {},
				"grateful",
			},
			"grateful": {
				"what are you grateful for?",
				TEXT,
				[]string{},
				func() {},
				"",
			},
		},
	}

	conversationHandler := ConversationHandler{

		// states should contain the possible fsm transitions
		// for simplicity, start with a single state
		Active: false,
		EntryPoints: []string{
			"journal",
		},
		CurrentState: "choosing",
		States: map[string]MessageHandler{
			"choosing": {
				"Please select a journal entry type",
				TEXT,
				[]string{"Morning Routine", "Evening Routine", "Story", "Reflection"},
				func() {},
				"selectEntry",
			},
			"selectEntry": {
				"Please enter the entry",
				SELECT_ENTRY,
				[]string{},
				func() {},
				"inputEntry",
			},
			"inputEntry": {
				"Thank you for the entry",
				INPUT_ENTRY,
				[]string{},
				func() {
					input := make(chan string, 1)
					go getInput(input)
					log.Println("input something")
					select {
					case i := <-input:
						log.Println("result")
						log.Println(i)
						switch i {
						case "Morning":
							{
								log.Println("Lets goooo")
							}
						}

					case <-time.After(5 * time.Minute):
						fmt.Println("timed out")
					}

				},
				"end",
			},
			"end": {
				"Thanks for the entry",
				TEXT,
				[]string{},
				func() {},
				"",
			},
		},
	}
	GlobalHandler(updates, b, &conversationHandler, &cronHandler)
}

// take current updates and compare to ongoing conversations, then add values of updates to the correct conversation
func GlobalHandler(updates tgbotapi.UpdatesChannel, b *tgbotapi.BotAPI, handlers *ConversationHandler, cronHandler *ConversationHandler) {
	// only one conversation handler should be active at a time
	var activeHandler ConversationHandler = *handlers
	var currentState MessageHandler = handlers.States[handlers.CurrentState]
	var stateLogger StateLogger = StateLogger{
		history: make([]State, 0),
	}
	c := cron.New()

	c.AddFunc(cronHandler.Cron, func(){
		log.Printf("cron job")
		currentState = cronHandler.States[cronHandler.CurrentState]
		activeHandler = *cronHandler
		cronHandler.Active = true
		msg := tgbotapi.NewMessage(5366512490, currentState.description)
		replyKeyboard := keyboardHelper(currentState.keyboardOptions)
		msg.ReplyMarkup = replyKeyboard
		b.Send(msg)
		stateLogger.addState(State{
			currentState,
			"start",
			currentState.description,
		})	
	})

	c.Start()

	for {

		for update := range updates {
			if update.Message != nil {
				// first message && checks if the the command is a valid entry point
				if Contains(handlers.EntryPoints, update.Message.Command()) {
					handlers.Active = true
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
					msg.ReplyToMessageID = update.Message.MessageID
					// replyKeyboard := keyboardHelper(currentState.keyboardOptions)
					// msg.ReplyMarkup = replyKeyboard
					b.Send(msg)

					stateLogger.addState(State{
						currentState,
						"start",
						update.Message.Text,
					})

					handlers.CurrentState = currentState.nextState
					currentState = handlers.States[handlers.CurrentState]
					// check if handler is active
				} else if activeHandler.Active {

					// check if the state is input
					if currentState.StateType == SELECT_ENTRY {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
						msg.ReplyToMessageID = update.Message.MessageID
						stateLogger.addState(State{
							currentState,
							stateLogger.history[len(stateLogger.history)-1].input,
							update.Message.Text,
						})
						b.Send(msg)
						log.Printf(update.Message.Text)
						handlers.CurrentState = currentState.nextState
						currentState = handlers.States[handlers.CurrentState]

					} else if currentState.StateType == INPUT_ENTRY {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
						msg.ReplyToMessageID = update.Message.MessageID
						stateLogger.addState(State{
							currentState,
							stateLogger.history[len(stateLogger.history)-1].input,
							update.Message.Text,
						})
						// parse state and add it to the db

						EnterJournalEntry(stateLogger.history[len(stateLogger.history)-1].field, stateLogger.history[len(stateLogger.history)-1].input)

						b.Send(msg)
						handlers.CurrentState = currentState.nextState
						currentState = handlers.States[handlers.CurrentState]
						log.Printf("input entry to db")
						log.Printf(stateLogger.history[len(stateLogger.history)-1].field)
						log.Printf(stateLogger.history[len(stateLogger.history)-1].input)

					} else {
						log.Println(handlers.CurrentState)
						log.Println(update.Message.Text)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
						b.Send(msg)
						handlers.CurrentState = currentState.nextState
						currentState = handlers.States[handlers.CurrentState]
						activeHandler = *handlers
					}
					// this is for inline keyboards, generally we can ignore this implementation if we dont use callbacks
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
	}
}
