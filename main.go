package main

import (
	sc "github.com/maddevsio/simple-config"
	"gopkg.in/resty.v1"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/carlescere/scheduler"

	"log"
	"time"
	"encoding/json"
	"fmt"
	"strconv"
	"runtime"
)

type logData struct {
	Target     string      `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

type Config struct {
    Url    string
    Token  string
    ChatID int64

	TimeForYesterdayData string
	TimeForDriversData   string
}

var config = Config{}

const (
	TARGET_ORDERS_TOTAL    = "taxi.orders.total"
	TARGET_ORDERS_FINISHED = "taxi.orders.finished"
	TARGET_ORDERS_REJECTED = "taxi.orders.rejected"
)

func (cs *Config) Fill(configFile string, configExt string) {
    c         := sc.NewSimpleConfig(configFile, configExt)
    cs.Url    = c.GetString("url")
    cs.Token  = c.GetString("token")
    cs.TimeForYesterdayData  = c.GetString("timeforyesterdaydata")
    cs.ChatID = int64(c.Get("chatid").(int))
}

func GetDayBeforeInFormat(t time.Time) string {
    return t.AddDate(0, 0, -1).Format("20060102")
}

func GetMaxForDateAndTarget(date string, target string, config Config) string {
	url       := fmt.Sprintf(config.Url, date, target)
	resp, err := resty.R().Get(url)
	checkErr(err)
	return strconv.Itoa(GetMaxDataFromJSON(resp.String()))
}

func GetMaxDataFromJSON(raw string) int {
    var data []logData
    err := json.Unmarshal([]byte(raw), &data)
	checkErr(err)
    if len(data) == 0 {
		return 0
	}

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

func ConnectTelegramAndSendMessage(message string, config Config) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	checkErr(err)
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true
	msg := tgbotapi.NewMessage(config.ChatID, message)
	bot.Send(msg)
}

func GetRejectPercent(maxTotal string, maxRejected string) string {
	maxRejectedInt, err := strconv.Atoi(maxRejected)
	checkErr(err)
	maxTotalInt, err := strconv.Atoi(maxTotal)
	checkErr(err)
	return strconv.Itoa(maxRejectedInt*100/maxTotalInt) + "%"
}

func CreateMessageForYesterday() string {
	message := "СТАТИСТИКА ЗА ВЧЕРА: \n"

	maxTotal := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()), TARGET_ORDERS_TOTAL, config)
	message  += "Всего заказов: " + maxTotal + "\n"

	maxRejected := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()), TARGET_ORDERS_REJECTED, config)
	message     += "Всего отмененных заказов: " + maxRejected + "\n"

	maxFinished := GetMaxForDateAndTarget(GetDayBeforeInFormat(time.Now()), TARGET_ORDERS_FINISHED, config)
	message     += "Всего выполненных заказов: " + maxFinished + "\n"

	rejectedPercent := GetRejectPercent(maxTotal, maxRejected)
	message         += "Процент отмент: " + rejectedPercent

	return message
}

func main() {
	config.Fill("./config", "yml")
	log.Printf("scheduler for TFYD: %s", config.TimeForYesterdayData)
	job := func() {
		message := CreateMessageForYesterday()
		ConnectTelegramAndSendMessage(message, config)
	}
	scheduler.Every().Day().At(config.TimeForYesterdayData).Run(job)
	runtime.Goexit()
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
