package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	tm "github.com/and3rson/telemux/v2"
)

type StateType string

const (
	SELECT_ENTRY StateType = "selectEntry"
	INPUT_ENTRY  StateType = "inputEntry"
	TEXT         StateType = "text"
	NUMBER       StateType = "number"
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
	Cron         string
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

// TODO
func Contains[T comparable](arr []T, e T) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

func isValidEntrypoint(entrypoint string, handlers []ConversationHandler) (bool, ConversationHandler) {
	for _, handler := range handlers {
		if Contains(handler.EntryPoints, entrypoint) {
			return true, handler
		}
	}
	return false, ConversationHandler{}
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

// Photo describes a submitted photo
type Photo struct {
	ID          int
	FileID      string
	Description string
}

type Journal struct {
	ID       int
	Title    string
	Text     string
	DateTime time.Time
}

func StartBot() {

	bot, err := tgbotapi.NewBotAPI("5366512490:AAFNjdosYKeQofgp4BdI0ehUissp7-sIGRM")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var photos []Photo
    var journals []Journal

	var lastID = 0

	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"journal",
			tm.NewLocalPersistence(), // we could also use `tm.NewFilePersistence("db.json")` or `&gormpersistence.GORMPersistence(db)` to keep data across bot restarts
			tm.StateMap{
				"": {
					tm.NewHandler(tm.IsCommandMessage("start"), func(u *tm.Update) {
						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Select a journal entry type",
						)
                        msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("Morning Entry"),
								tgbotapi.NewKeyboardButton("Evening Entry"),
							),
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("Random Story"),
								tgbotapi.NewKeyboardButton("Quotes/Thoughts"),
							),
                        )
						bot.Send(msg)
						u.PersistenceContext.SetState("enter_entry")
					}),
				},
				"enter_entry": {
					tm.NewHandler(tm.Not(tm.IsCommandMessage("cancel")), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["journal_entry"] = u.Message.Text
						u.PersistenceContext.SetData(data)

						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Please enter the text for the journal entry",
						)
 						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

						bot.Send(msg)
						u.PersistenceContext.SetState("enter_text")
					}),
				},
				"enter_text": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["journal_text"] = u.Message.Text
						u.PersistenceContext.SetData(data)
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Thanks for the submission, type /list to view previous submissions")
						bot.Send(msg)
                        journals = append(journals, Journal{
                            ID: lastID,
                            Title: data["journal_entry"].(string),
                            Text: data["journal_text"].(string),
                            DateTime: time.Now(),
                        })
						u.PersistenceContext.SetState("confirm_submission")
					}),
					tm.NewHandler(tm.Not(tm.IsCommandMessage("cancel")), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Sorry, I did not understand that. Please enter some text!",
						))
					}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
				}),
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("list"),
			func(u *tm.Update) {
				var lines []string
				for _, photo := range photos {
					lines = append(lines, fmt.Sprintf("- %s (/view_%d)", photo.Description, photo.ID))
				}
				if len(lines) == 0 {
					lines = append(lines, "No photos yet.")
				}
				message := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Photos:\n"+strings.Join(lines, "\n"),
				)
				message.ReplyToMessageID = u.Message.MessageID
				bot.Send(message)
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasRegex(`^/view_(\d+)$`),
			func(u *tm.Update) {
				photoID := strings.Split(u.Message.Text, "_")[1]
				var match *Photo
				for i, photo := range photos {
					if fmt.Sprint(photo.ID) == photoID {
						match = &photos[i]
					}
				}
				if match == nil {
					bot.Send(tgbotapi.NewMessage(
						u.Message.Chat.ID,
						"Photo not found!",
					))
				} else {
					share := tgbotapi.NewPhoto(u.Message.Chat.ID, tgbotapi.FileID(match.FileID))
					share.Caption = fmt.Sprintf("Description: %s", match.Description)
					_, err := bot.Send(share)
					if err != nil {
						log.Printf("main: send file: %s", err)
					}
				}
			},
		)).
		AddHandler(tm.NewMessageHandler(
			nil,
			func(u *tm.Update) {
				message := tgbotapi.NewMessage(
					u.Message.Chat.ID,
                    "Hello! I'm a Vasco's bot.\n\nI allow Vasco to update his database and track EVERYTHING!\n\nAvailable commands:\n/start - Start a journal entry \n/list - list the past 10 journal entries",

				)
				message.ReplyToMessageID = u.Message.MessageID
				bot.Send(message)
			},
		))
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}

func StartBotOld() {
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
		Active:       false,
		Cron:         "* * * * *",
		CurrentState: "sleep",
		States: map[string]MessageHandler{
			"sleep": {
				description: "How many hours did you sleep?",
				StateType:   INPUT_ENTRY,
				nextState:   "mood",
			},
			"mood": {
				description: "How are you feeling?",
				StateType:   INPUT_ENTRY,
				nextState:   "grateful",
			},
			"grateful": {
				description: "what are you grateful for?",
				StateType:   INPUT_ENTRY,
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

	var handlers []ConversationHandler = []ConversationHandler{cronHandler, conversationHandler}
	GlobalHandler(updates, b, handlers)
}

// take current updates and compare to ongoing conversations, then add values of updates to the correct conversation
func GlobalHandler(updates tgbotapi.UpdatesChannel, b *tgbotapi.BotAPI, handlers []ConversationHandler) {
	// There are two types of conversations Handlers, cronHandlers and regular ConversationsHandlers
	// There can only be one active conversation at any time, if a cron job is being executed, no other cron jobs can be ran and a queue should be built up
	// eg. If a user is in the middle of a conversation, don't start a new one

	// var currentState MessageHandler = handlers.States[handlers.CurrentState]
	// var currentHandler ConversationHandler = *handlers
	var currentState MessageHandler
	var currentConversationHandler ConversationHandler
	var stateLogger StateLogger = StateLogger{
		history: make([]State, 0),
	}
	c := cron.New()

	for _, handler := range handlers {
		_, err := cron.ParseStandard(handler.Cron)
		if err != nil {
			log.Printf("Non cron Schedule")
		} else {
			log.Printf("Cron Schedule added")
			cronHandler := handler

			c.AddFunc(handler.Cron, func() {
				log.Printf("cron job")
				cronHandler.Active = true
				currentState = cronHandler.States[cronHandler.CurrentState]
				currentConversationHandler = cronHandler
				msg := tgbotapi.NewMessage(5383565084, currentState.description)
				replyKeyboard := keyboardHelper(currentState.keyboardOptions)
				msg.ReplyMarkup = replyKeyboard
				b.Send(msg)
				stateLogger.addState(State{
					currentState,
					"start",
					currentState.description,
				})
			})
		}
	}

	c.Start()

	for update := range updates {
		if update.Message != nil {
			// first message && checks if the the command is a valid entry point
			isValid, entrypointHandler := isValidEntrypoint(update.Message.Text, handlers)
			if isValid {
				entrypointHandler.Active = true
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

				entrypointHandler.CurrentState = currentState.nextState
				currentState = entrypointHandler.States[entrypointHandler.CurrentState]
				// check if handler is active
			} else if currentConversationHandler.Active {

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
					currentConversationHandler.CurrentState = currentState.nextState
					currentState = currentConversationHandler.States[currentConversationHandler.CurrentState]

				} else if currentState.StateType == INPUT_ENTRY {
					stateLogger.addState(State{
						currentState,
						stateLogger.history[len(stateLogger.history)-1].input,
						update.Message.Text,
					})
					// parse state and add it to the db
					EnterJournalEntry(stateLogger.history[len(stateLogger.history)-1].field, stateLogger.history[len(stateLogger.history)-1].input)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
					msg.ReplyToMessageID = update.Message.MessageID
					b.Send(msg)

					currentConversationHandler.CurrentState = currentState.nextState
					currentState = currentConversationHandler.States[currentConversationHandler.CurrentState]
				} else {
					log.Println(currentConversationHandler.CurrentState)
					log.Println(update.Message.Text)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, currentState.description)
					b.Send(msg)
					currentConversationHandler.CurrentState = currentState.nextState
					currentState = currentConversationHandler.States[currentConversationHandler.CurrentState]
				}
				if currentState.nextState == "" {
					currentConversationHandler.Active = false
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
