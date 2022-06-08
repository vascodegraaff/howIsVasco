package controllers

import (
	"log"
	"time"
	tele "gopkg.in/telebot.v3"
	"os"

)
func StartBot(){
	pref := tele.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	var (
		// Universal markup builders.
		menu     = &tele.ReplyMarkup{ResizeKeyboard: true}
		selector = &tele.ReplyMarkup{}
	
		// Reply buttons.
		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")
	
		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind
		// since it's required for callback routing to work.
		//
		btnPrev = selector.Data("⬅", "prev")
		btnNext = selector.Data("➡", "next")
	)
	
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)
	
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!", menu)
	})
	
	// On reply button pressed (message)
	bot.Handle(&btnHelp, func(c tele.Context) error {
		return c.Edit("Here is some help: ...")
	})
	
	// On inline button pressed (callback)
	bot.Handle(&btnPrev, func(c tele.Context) error {
		return c.Respond()
	})

	bot.Start()
}