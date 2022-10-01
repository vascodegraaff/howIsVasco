package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

func CreateFakeUpdate(msg string) tgbotapi.Update {
	return tgbotapi.Update{
		UpdateID: 123122213,
		Message: &tgbotapi.Message{
			MessageID: 23432,
			From: &tgbotapi.User{
				ID: 1231232121321,
			},
			Date: 1231232121321,
			Chat: &tgbotapi.Chat{
				ID:    5383565084,
				Title: "test chat",
			},
			Text: msg,
		},
	}
}

func commitData(data tm.Data) {
	for key, value := range data {
		if s, ok := value.(string); ok {
			fmt.Printf("%q is a string: %q\n", key, s)
		}
	}
}

func StartBot() {

	bot, err := tgbotapi.NewBotAPI("5366512490:AAFNjdosYKeQofgp4BdI0ehUissp7-sIGRM")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	c := cron.New()

	c.Start()

	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"MorningJournal",
			tm.NewLocalPersistence(),
			tm.StateMap{
				"": {
					tm.NewHandler(tm.IsCommandMessage("morningJournal"), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Good Morning. Lets start logging!",
						))
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"How many hours did you sleep last night",
						))
						u.PersistenceContext.SetState("sleepTime")
					}),
				},
				"sleepTime": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						n, err := strconv.Atoi(u.Message.Text)
						if err != nil {
							bot.Send(tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"Please input the hours you slept last night",
							))
							u.PersistenceContext.SetState("sleepTime")
						} else {
							data["sleepTime"] = n
							u.PersistenceContext.SetData(data)
							bot.Send(tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"How do you feel right now? 1-10",
							))
							u.PersistenceContext.SetState("howDoYouFeel")
						}
					}),
				},
				"howDoYouFeel": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						n, err := strconv.Atoi(u.Message.Text)
						if err != nil && n > 0 && n <= 10 {
							bot.Send(tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"Please input a value between 1-10",
							))
							u.PersistenceContext.SetState("howDoYouFeel")
						} else {
							data["howDoYouFeel"] = n
							u.PersistenceContext.SetData(data)
							bot.Send(tgbotapi.NewMessage(
								u.Message.Chat.ID,
								"What's something you're grateful for today",
							))
							u.PersistenceContext.SetState("gratitude")
						}
					}),
				},
				"gratitude": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["gratitude"] = u.Message.Text
						u.PersistenceContext.SetData(data)
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"What would make today great",
						))
						u.PersistenceContext.SetState("great")
					}),
				},
				"great": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["great"] = u.Message.Text
						u.PersistenceContext.SetData(data)
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"What's today daily affirmation",
						))
						u.PersistenceContext.SetState("affirmation")
					}),
				},
				"affirmation": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["affirmation"] = u.Message.Text
						u.PersistenceContext.SetData(data)
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Logging complete. Have a sick day!",
						))
						u.PersistenceContext.SetState("commit")
						commitData(data)
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
		))

	c.AddFunc("10 23 * * *", func() {
		mux.Dispatch(bot, CreateFakeUpdate("/morningJournal"))
	})

	updates := bot.GetUpdatesChan(u)

	log.Printf(time.Now().String())
	for update := range updates {
		log.Println(update)
		mux.Dispatch(bot, update)
	}
}
