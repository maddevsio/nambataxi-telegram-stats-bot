package main

import (
	sc "github.com/maddevsio/simple-config"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"time"
	"encoding/json"
)

type logData struct {
    Target     string      `json:"target"`
    Datapoints [][]float64 `json:"datapoints"`
}

func GetDayBeforeInFormat(t time.Time) string {
    return t.AddDate(0, 0, -1).Format("20060102")
}

func GetMaxDataFromJSON(raw string) int {
    var data []logData
    _ = json.Unmarshal([]byte(raw), &data)

    var picked []int
    var max int
    for _, v := range data[0].Datapoints {
        if v[0] > 0 {
            if int(v[0]) > max {
                max = int(v[0])
            }
            picked = append(picked, int(v[0]))
        }
    }
    return max
}

func main() {
	config := sc.NewSimpleConfig("./config", "yml")
	token  := config.GetString("token")
	chatid := int64(config.Get("chatid").(int))

	bot, err := tgbotapi.NewBotAPI(token)
	checkErr(err)
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	msg := tgbotapi.NewMessage(chatid, "Hello")
	bot.Send(msg)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
