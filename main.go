package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"errors"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type bnResp struct {
	Price float64 `json:"price,string"`
	Code int64 `json:"code"`
}
type bnRespRub struct {
	Price float64 `json:"price,string"`
	Code int64 `json:"code"`
}

type wallet map[string]float64
var db = map[int64]wallet{}

func main() {
	bot, err := tgbotapi.NewBotAPI("1862395535:AAEWEa3cAVuaX7eEH37zhpz6mbBtt5Lggco")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Split(update.Message.Text, " ")
		switch command[0] {
		case "ADD": 
		if len(command) != 3 {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
		}
		amount, err := strconv.ParseFloat(command[2], 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
		}
		if _, ok := db[update.Message.Chat.ID]; !ok {
			db[update.Message.Chat.ID] = wallet{}
		}
		db[update.Message.Chat.ID][command[1]] += amount
		balanceTExt := fmt.Sprintf("%f", db[update.Message.Chat.ID][command[1]])
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceTExt))
		case "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			}
			if _, ok := db[update.Message.Chat.ID]; !ok {
				continue
			}
			db[update.Message.Chat.ID][command[1]] -= amount
			balanceTExt := fmt.Sprintf("%f", db[update.Message.Chat.ID][command[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceTExt))
		case "DEL":
			if len(command) != 2 {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			}
			delete(db[update.Message.Chat.ID], command[1])
		case "SHOW":
			msg:= ""
			var sum float64 = 0
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getPrice(key)
				priceRub, _ := getRoubles(price)
				sum += value*price*priceRub
				msg += fmt.Sprintf("%s: %f [%.2f]\n", key, value, value*price*priceRub)
			}
			msg += fmt.Sprintf("Total: %.2f\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default: 
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не найдена"))
		}
	}

	}
	func getPrice(symbol string) (price float64, err error) {
		resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol))
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var jsonResp bnResp

		err = json.NewDecoder(resp.Body).Decode(&jsonResp)
		if err != nil {
			return
		}
		if jsonResp.Code != 0 {
			err = errors.New("Неверный символ")
		}
		price = jsonResp.Price
		return
	}
	func getRoubles(price float64) (priceRub float64, err error) {
		resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=USDTRUB")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var jsonRespRub bnRespRub
		err = json.NewDecoder(resp.Body).Decode(&jsonRespRub)
		if err != nil {
			return
		}
		if jsonRespRub.Code != 0 {
			err = errors.New("Неверный символ")
		}
		priceRub = jsonRespRub.Price
		return
	} 


