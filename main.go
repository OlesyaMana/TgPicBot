package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pic/db"
	"pic/keyboard"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BotToken  = "6356948192:AAH-OAe0VLuTbwfcOooh81-L9xVm-jqrj0s"
	DogUrl    = "https://dog.ceo/api/breeds/image/random"
	SleepTime = 30
)

var botUsersMap map[int64]db.BotUser

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 3

	updates := bot.GetUpdatesChan(u)

	database, err := db.StartDB()
	if err != nil {
		log.Panic(err)
	}

	botUsersMap, err = db.GetBotUsers(database)
	if err != nil {
		log.Panic(err)
	}

	go startSpamming(bot)
	for update := range updates {
		if update.Message != nil {
			user, contains := botUsersMap[update.Message.Chat.ID]
			if contains {
				switch update.Message.Text {
				case "/start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "keyboard")
					msg.ReplyMarkup = keyboard.GetKeyboard()
					bot.Send(msg)

				case "start":
					err = db.UpdateUser(database, user.Id, db.DogsColumn, 1)
					if err != nil {
						panic(err)
					}

					user.ReceiveDogs = 1
					botUsersMap[user.Id] = user
				case "stop":
					err = db.UpdateUser(database, user.Id, db.DogsColumn, 0)
					if err != nil {
						panic(err)
					}

					user.ReceiveDogs = 0
					botUsersMap[user.Id] = user

				default:
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "unknown command"))
				}
			} else {
				err = db.AddNewUser(database, update.Message.Chat.ID, update.Message.Chat.UserName)
				if err != nil {
					log.Panic(err)
				}

				botUsersMap[update.Message.Chat.ID] = db.BotUser{
					Id:    update.Message.Chat.ID,
					Login: update.Message.Chat.UserName,
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы добалены в бот, для получения собак напишите dogs, а для остановки stop dogs"))
			}
		}
	}
}

func startSpamming(bot *tgbotapi.BotAPI) {
	for {
		for id, user := range botUsersMap {
			if user.ReceiveDogs == 1 {
				startSendDogs(id, bot)
			}
		}

		time.Sleep(SleepTime * time.Second)
	}
}
func startSendDogs(chatId int64, bot *tgbotapi.BotAPI) {
	client := http.Client{}
	resp, err := client.Get(DogUrl)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	var result DogPicture
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Panic(err)
	}

	bot.Send(tgbotapi.NewMessage(chatId, result.Message))
	log.Printf("sended to %d", chatId)
}

type DogPicture struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
