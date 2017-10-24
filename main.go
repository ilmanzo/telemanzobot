// telemanzobot project main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	telebot "gopkg.in/tucnak/telebot.v1"
)

// APITOKEN set the Telegram Token assigned to the Bot
const APITOKEN = "PUT HERE YOUR OWN TOKEN"

var bot *telebot.Bot

func main() {
	newBot, err := telebot.NewBot(APITOKEN)

	if err != nil {
		fmt.Printf("Error connecting to telegram server: %v", err.Error())
		return
	}

	bot = newBot

	bot.Messages = make(chan telebot.Message, 1000)
	bot.Queries = make(chan telebot.Query, 1000)

	go messages()
	go queries()
	fmt.Println("Welcome to TELEgram manzobot. Starting listening for messages...")
	bot.Start(1 * time.Second)
}

func temperatura() string {
	cmd := exec.Command("/usr/bin/sensors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v", err.Error())
		return "err"
	}
	return string(output)
}

func usage(bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		"Hi, this is Telemanzobot.\n"+
			"Available commands:\n"+
			"/hello\n"+
			"/help\n"+
			"/temperatura [execute /usr/bin/sensors]\n"+
			"/quit <password> [used for debugging]",
		nil)

}

// handle direct messages
func messages() {
	for message := range bot.Messages {
		switch message.Text {
		case "/ciao":
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!", nil)
		case "/help":
			usage(bot, message)
		case "/temperatura":
			msg := temperatura()
			if len(msg) > 0 {
				bot.SendMessage(message.Chat, msg, nil)
			}
		default:
			msg := strings.Split(message.Text, " ")
			if len(msg) < 2 {
				usage(bot, message)
				continue
			}
			if msg[0] == "/quit" {
				if msg[1] == "rossodisera" {
					bot.SendMessage(message.Chat, "Ok I'll leave the scene", nil)
					os.Exit(1)
				} else {
					bot.SendMessage(message.Chat, "sorry, wrong password", nil)
				}

			}
		}
	}
}

// inline query
func queries() {
	for query := range bot.Queries {
		fmt.Println("--- new query ---")
		fmt.Println("from:", query.From)
		fmt.Println("text:", query.Text)

		// https://core.telegram.org/bots/api#inlinequeryresult

		// per ora fissi
		results := []telebot.Result{
			telebot.ArticleResult{Title: "risultato1", Text: "testo risultato 1"},
			telebot.ArticleResult{Title: "risultato2", Text: "testo risultato 2"},
		}

		// respond to the query:
		if err := bot.Respond(query, results); err != nil {
			log.Println("ouch:", err)
		}
	}
}
